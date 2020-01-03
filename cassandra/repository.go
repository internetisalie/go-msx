package cassandra

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	gocqlxqb "github.com/scylladb/gocqlx/qb"
)

type CrudRepository struct {
	Table ddl.Table
}

func (r *CrudRepository) FindAll(ctx context.Context, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSession(func(session *gocql.Session) error {
		stmt, names := gocqlxqb.
			Select(r.Table.Name).
			Columns(r.Table.ColumnNames()...).
			ToCql()

		return gocqlx.
			Query(session.Query(stmt), names).
			SelectRelease(dest)
	})

	return
}

func (r *CrudRepository) FindAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSession(func(session *gocql.Session) error {
		var cmps []gocqlxqb.Cmp
		for k, _ := range where {
			cmps = append(cmps, gocqlxqb.EqNamed(k, k))
		}

		stmt, names := gocqlxqb.
			Select(r.Table.Name).
			Columns(r.Table.ColumnNames()...).
			Where(cmps...).
			ToCql()

		return gocqlx.
			Query(session.Query(stmt), names).
			BindMap(where).
			SelectRelease(dest)
	})

	return
}

func (r *CrudRepository) FindOneBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSession(func(session *gocql.Session) error {
		var cmps []gocqlxqb.Cmp
		for k, _ := range where {
			cmps = append(cmps, gocqlxqb.EqNamed(k, k))
		}

		stmt, names := gocqlxqb.
			Select(r.Table.Name).
			Columns(r.Table.ColumnNames()...).
			Where(cmps...).
			ToCql()

		return gocqlx.
			Query(session.Query(stmt), names).
			BindMap(where).
			GetRelease(dest)
	})

	return
}

func (r *CrudRepository) Save(ctx context.Context, value interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSessionRetry(ctx, func(session *gocql.Session) error {
		stmt, names := gocqlxqb.
			Insert(r.Table.Name).
			Columns(r.Table.ColumnNames()...).
			ToCql()

		return gocqlx.Query(session.Query(stmt), names).
			BindStruct(value).
			ExecRelease()
	})

	return
}

func (r *CrudRepository) UpdateBy(ctx context.Context, where map[string]interface{}, values map[string]interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSessionRetry(ctx, func(session *gocql.Session) error {
		var bind = make(map[string]interface{})
		var cmps []gocqlxqb.Cmp
		for k, v := range where {
			cmps = append(cmps, gocqlxqb.EqNamed(k, k + "Where"))
			bind[k + "Where"] = v
		}

		var sets []string
		for k, v := range values {
			sets = append(sets, k)
			bind[k] = v
		}

		stmt, names := gocqlxqb.
			Update(r.Table.Name).
			Set(sets...).
			Where(cmps...).
			ToCql()

		return gocqlx.Query(session.Query(stmt), names).
			BindMap(bind).
			ExecRelease()
	})

	return
}

func (r *CrudRepository) DeleteBy(ctx context.Context, where map[string]interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSessionRetry(ctx, func(session *gocql.Session) error {
		var cmps []gocqlxqb.Cmp
		for k, _ := range where {
			cmps = append(cmps, gocqlxqb.EqNamed(k, k))
		}

		stmt, names := gocqlxqb.
			Delete(r.Table.Name).
			Where(cmps...).
			ToCql()

		return gocqlx.Query(session.Query(stmt), names).
			BindMap(where).
			ExecRelease()
	})

	return
}
