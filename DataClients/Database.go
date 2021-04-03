package dataclients

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Database
type Database struct {
	pg            *pgxpool.Pool
	propertyCache sync.Map
	propertyMu    sync.Mutex
	kindMu        sync.Mutex
	firstRun      bool
	tx            pgx.Tx
}

//New Instance
func NewDatabase(db *pgxpool.Pool) *Database {
	d := &Database{
		pg: db,
	}
	return d
}

//tx
func (d *Database) BeginTx() (err error) {
	d.tx, err = d.pg.Begin(context.Background())
	return
}

//rollback
func (d *Database) Rollback() error {
	if tx := d.tx; tx != nil {
		err := d.tx.Rollback(context.Background())
		if err == pgx.ErrTxClosed {
			return nil
		}
		return err
	}
	return nil
}
