// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package migrate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

const (
	configKeySqlUsername = "spring.datasource.username"
)

var (
	logger = log.NewLogger("msx.repository.sql.migrate")
)

type Migrator struct {
	ctx       context.Context
	manifest  *Manifest
	db        *sqlx.DB
	versioner Versioner
}

func (m *Migrator) ValidateMigration(n int, migration Migration, appliedMigration AppliedMigration) error {
	// Ensure manifest migration matches applied migration
	if appliedMigration.InstalledRank != n+1 {
		return errors.Errorf("Incorrect installed rank: %+v", appliedMigration)
	}

	if !appliedMigration.Success {
		return errors.Errorf("Failed migration recorded: %+v", appliedMigration)
	}

	if appliedMigration.Type == "JDBC" {
		logger.WithContext(m.ctx).Warnf(
			"Skipping verification of existing JDBC migration %s checksum and description",
			migration.Version)
		return nil
	}

	// skip description and checksum validation if checksum is equal to 0
	// used when we patch a released migration step and need to skip the checksum validation
	if migration.Checksum == nil || *migration.Checksum != 0 {
		if appliedMigration.Description != migration.Description {
			return errors.Errorf("Mismatched description: %+v", appliedMigration)
		}

		if appliedMigration.Checksum == nil && migration.Checksum != nil ||
			appliedMigration.Checksum != nil && migration.Checksum == nil {
			return errors.Errorf("Mismatched checksum: %+v vs %+v", appliedMigration.Checksum, migration.Checksum)
		} else if appliedMigration.Checksum != nil &&
			migration.Checksum != nil &&
			*appliedMigration.Checksum != *migration.Checksum {
			return errors.Errorf("Mismatched checksum: %d vs %d", *appliedMigration.Checksum, *migration.Checksum)
		}
	}

	return nil
}

func (m *Migrator) ValidateManifest(appliedMigrations []AppliedMigration, preUpgrade bool) error {
	logger.WithContext(m.ctx).Info("Validating previously applied migrations")

	postUpgradeVersion, err := m.manifest.PostUpgradeVersion()
	if err != nil {
		return errors.Wrap(err, "Failed to parse Post-Upgrade Version")
	}

	if postUpgradeVersion == nil {
		// No post-upgrade version set, run all migrations
		preUpgrade = false
	}

	// Version where issue that skipped validation was fixed.  See usage below.
	minValidationVersion, _ := types.NewVersion("4.3.0")
	n := 0
	var migration *Migration
	for n, migration = range m.manifest.Migrations() {
		if n == len(appliedMigrations) {
			break
		}
		appliedMigration := appliedMigrations[n]

		if preUpgrade && !migration.Version.Lt(postUpgradeVersion) {
			logger.WithContext(m.ctx).
				WithField("version", migration.Version).
				Infof("Skipping verification of post-upgrade %s migration %s: %s",
					migration.Type,
					migration.Version,
					migration.Description)
			continue
		}

		// Preserve earlier behavior of skipping validation/ignoring errors with older migrations (before the error
		// check was fixed).  Enabling validation on the older migrations is risky - if a migration was
		// modified after release, a failure may not appear until a production upgrade.
		if migration.Version.Lt(minValidationVersion) {
			continue
		}

		err = m.ValidateMigration(n, *migration, appliedMigration)
		if err != nil {
			return err
		}

		logger.WithContext(m.ctx).
			WithField("version", migration.Version).
			Infof("Validated %s migration %s: %s",
				migration.Type,
				migration.Version,
				migration.Description)

	}

	return nil
}

func (m *Migrator) ApplyMigrations(lastAppliedMigration int, userName string, preUpgrade bool) (err error) {
	logger.WithContext(m.ctx).Info("Applying new migrations")

	postUpgradeVersion, err := m.manifest.PostUpgradeVersion()
	if err != nil {
		return errors.Wrap(err, "Failed to parse Post-Upgrade Version")
	}

	if postUpgradeVersion == nil {
		// No post-upgrade version set, run all migrations
		preUpgrade = false
	}

	migrations := m.manifest.Migrations()
	for n := lastAppliedMigration; n < len(migrations); n++ {
		migration := migrations[n]

		if preUpgrade && !migration.Version.Lt(postUpgradeVersion) {
			logger.WithContext(m.ctx).
				WithField("version", migration.Version).
				Infof("Skipping post-upgrade %s migration %s: %s",
					migration.Type,
					migration.Version,
					migration.Description)
			continue
		}

		appliedMigration := AppliedMigration{
			Version:       migration.Version.String(),
			Description:   migration.Description,
			Script:        migration.Script,
			Type:          migration.Type,
			Checksum:      migration.Checksum,
			InstalledOn:   time.Now(),
			InstalledBy:   userName,
			InstalledRank: n + 1,
		}

		logger.WithContext(m.ctx).
			WithField("version", migration.Version).
			Infof("Applying %s migration %s: %s",
				migration.Type,
				migration.Version,
				migration.Description)

		if err := migration.Func(m.ctx, m.db); err != nil {
			return errors.Wrap(err, "Migration failed")
		}

		appliedMigration.Success = true
		appliedMigration.ExecutionTime = int(time.Since(appliedMigration.InstalledOn) / time.Millisecond)

		if err := m.versioner.RecordAppliedMigration(appliedMigration); err != nil {
			return errors.Wrap(err, "Failed to record applied migration")
		}
	}
	return nil
}

func (m *Migrator) Migrate(preUpgrade bool) error {
	if err := m.versioner.CreateVersionTables(); err != nil {
		return err
	}

	appliedMigrations, err := m.versioner.GetAppliedMigrations()
	if err != nil {
		return err
	}

	if err = m.ValidateManifest(appliedMigrations, preUpgrade); err != nil {
		return err
	}

	userName, err := config.FromContext(m.ctx).StringOr(configKeySqlUsername, "root")
	if err != nil {
		return err
	}

	if err = m.ApplyMigrations(len(appliedMigrations), userName, preUpgrade); err != nil {
		return err
	}

	logger.WithContext(m.ctx).Info("Database migration completed successfully.")

	return nil
}

func NewMigrator(ctx context.Context, db *sqlx.DB) (*Migrator, error) {
	versioner, err := NewVersioner(ctx, db)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		ctx:       ctx,
		manifest:  ManifestFromContext(ctx),
		db:        db,
		versioner: versioner,
	}, nil
}

func Migrate(ctx context.Context) error {
	logger.WithContext(ctx).Info("Executing SQL db migrate")

	preUpgrade, _ := config.FromContext(ctx).BoolOr("cli.flag.preupgrade", false)

	sqlPool, err := sqldb.PoolFromContext(ctx)
	if err == sqldb.ErrDisabled {
		logger.WithContext(ctx).WithError(err).Warn("Skipping SQL db migration.")
		return nil
	} else if err != nil {
		return err
	}

	return sqlPool.WithSqlxConnection(ctx, func(ctx context.Context, db *sqlx.DB) error {
		migrator, err := NewMigrator(ctx, db)
		if err != nil {
			return err
		}
		return migrator.Migrate(preUpgrade)
	})
}
