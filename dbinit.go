package main

import (
	"context"

	"github.com/pkg/errors"
)

func (d *Database) Init() error {

	var exists bool
	err := d.pg.QueryRow(context.Background(),
		`SELECT EXISTS (
			SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'extra_data'
		)`).Scan(&exists)
	if err != nil {
		return errors.Wrap(err, "check for existing schema")
	}
	if exists {
		return nil
	}

	d.firstRun = true

	queries := []string{

		`CREATE TABLE restaurants (
			restaurantid SERIAL PRIMARY KEY,
			Address1 VARCHAR(120) NOT null default '',
			Address2 VARCHAR(120) not null default '',
			Address3 VARCHAR(120)not null default '',
			City VARCHAR(100) NOT null default '',
			State VARCHAR(120) NOT null default '',
			Country VARCHAR(120) NOT null default '',
			Zipcode VARCHAR(16) NOT null default '',
			restaurantname VARCHAR(16) NOT null default ''`,

		`CREATE TABLE "business_hours" (
				"restaurantid" integer REFERENCES restaurants (restaurantid),
				"day" VARCHAR(16) NOT NULL,
				"open_time" time,
				"close_time" time
		   )`,

		`CREATE TABLE "restaurant_menu_items" (
			"restaurantid" integer REFERENCES restaurants (restaurantid),
			name VARCHAR(120) NOT NULL,
			price VARCHAR(16) NOT NULL,
			description VARCHAR(120) NOT NULL,
			category VARCHAR(120) NOT NULL,
			image bytea NULL
	   ))`,
	}
	return d.batchExecute(queries)
}

func (d *Database) Prepare() error {
	return nil
}

func (d *Database) commit() error {
	if tx := d.tx; tx != nil {
		return tx.Commit(context.Background())
	}
	return nil
}

func (d *Database) Close() error {
	defer d.Rollback()

	if err := d.commit(); err != nil {
		return errors.Wrap(err, "commit")
	}

	if !d.firstRun {
		return nil
	}
	return nil
}

func (d *Database) batchExecute(queries []string) error {
	for _, q := range queries {
		if _, err := d.pg.Exec(context.Background(), q); err != nil {
			return errors.Wrapf(err, "query: %s", q)
		}
	}
	return nil
}
