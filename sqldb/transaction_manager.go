// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/jmoiron/sqlx"
)

func TransactionDecorator(action types.ActionFunc) types.ActionFunc {
	return func(ctx context.Context) error {
		pool, err := PoolFromContext(ctx)
		if err != nil {
			return err
		}

		err = pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
			tx, err := conn.Beginx()
			if err != nil {
				return err
			}

			ctx = ContextSqlExecutor().Set(ctx, tx)
			actionErr := action(ctx)

			if actionErr == nil {
				err = tx.Commit()
			} else {
				err = tx.Rollback()
			}

			if err == nil && actionErr != nil {
				err = actionErr
			}

			return err
		})

		return err
	}
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, action types.ActionFunc) error
}

func NewTransactionManager() TransactionManager {
	return new(SqlTransactionManager)
}

type SqlTransactionManager struct{}

func (t SqlTransactionManager) WithTransaction(ctx context.Context, action types.ActionFunc) error {
	wrappedAction := TransactionDecorator(action)
	return wrappedAction(ctx)
}
