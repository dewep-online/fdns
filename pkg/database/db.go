package database

import (
	"context"
	"database/sql"

	"github.com/deweppro/go-app/application/ctx"
	"github.com/deweppro/go-logger"
	"github.com/deweppro/go-orm"
	"github.com/deweppro/go-orm/schema"
	"github.com/deweppro/go-orm/schema/sqlite"
)

type Database struct {
	conn schema.Connector
	pool *orm.Stmt
	log  logger.Logger
}

func New(log logger.Logger, conf *sqlite.Config) (*Database, error) {
	conn := sqlite.New(conf)
	pool := orm.NewDB(conn, orm.Plugins{Logger: log})
	db := &Database{
		conn: conn,
		pool: pool.Pool(""),
		log:  log,
	}
	return db, nil
}

func (v *Database) Up(ctx ctx.Context) error {
	if err := v.conn.Reconnect(); err != nil {
		return err
	}
	if err := v.pool.Ping(); err != nil {
		return err
	}
	return v.pool.CallContext("check_tables", ctx.Context(), func(ctx context.Context, conn *sql.DB) error {
		row := conn.QueryRowContext(ctx, checkTables)
		var count int
		if err := row.Scan(&count); err != nil {
			return err
		}
		if err := row.Err(); err != nil {
			return err
		}
		if count == 0 {
			v.log.Infof("Creating SQLite database")
			for _, migration := range migrations {
				if _, err := conn.ExecContext(ctx, migration); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (v *Database) Down(_ ctx.Context) error {
	return v.conn.Close()
}
