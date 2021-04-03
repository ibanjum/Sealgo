package main

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	pg            *pgxpool.Pool
	propertyCache sync.Map
	propertyMu    sync.Mutex
	kindMu        sync.Mutex
	firstRun      bool
	tx            pgx.Tx
}

func NewDatabase(db *pgxpool.Pool) *Database {
	d := &Database{
		pg: db,
	}
	return d
}

func (d *Database) BeginTx() (err error) {
	d.tx, err = d.pg.Begin(context.Background())
	return
}

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
func (d *Database) SelectUser(username string) (User, error) {
	command := `SELECT * FROM users where username = $1`
	var u User
	err := d.pg.QueryRow(context.Background(), command, username).Scan(&u.ID, &u.Level, &u.Username, &u.FirstName, &u.LastName, &u.Email, &u.PhoneNumber, &u.Password, &u.Token)
	if err != nil {
		return u, err
	}

	return u, nil
}
func (d *Database) UserExists(username string) bool {
	command := `SELECT username FROM users where username = $1`
	var u string
	_ = d.pg.QueryRow(context.Background(), command, username).Scan(&u)
	if u != "" {
		return true
	}
	return false
}
func (d *Database) SelectRequest(email string) (RestaurantRequest, error) {
	command := `SELECT * FROM restaurantrequests where email = $1 and enabled = true`
	var r RestaurantRequest
	err := d.pg.QueryRow(context.Background(), command, email).Scan(&r.ID, &r.FirstName, &r.LastName, &r.RestaurantName, &r.Email, &r.PhoneNumber, &r.City, &r.State, &r.ZipCode, &r.Country)
	if err != nil {
		return r, err
	}
	return r, nil
}
func (d *Database) SelectAllRequests() ([]RestaurantRequest, error) {
	command := `SELECT * FROM restaurantrequests where active = true`
	var list []RestaurantRequest
	rows, err := d.pg.Query(context.Background(), command)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var r RestaurantRequest
		rows.Scan(&r.ID, &r.FirstName, &r.LastName, &r.RestaurantName, &r.Email, &r.PhoneNumber, &r.City, &r.State, &r.ZipCode, &r.Country, &r.Active, &r.Enabled)
		list = append(list, r)
	}
	return list, nil
}

func (d *Database) InsertRestaurant(r BusinessResponse, username string) error {
	command := `insert into restaurants (id, managerid, name, phone, price, rating,  url, reviewcount, imageurl) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := d.pg.Exec(context.Background(), command, r.ID, username, r.Name, r.Phone, r.Price, r.Rating, r.URL, r.ReviewCount, r.ImageURL)
	if err != nil {
		return err
	}
	l := r.Location
	command = `insert into locations (restaurantid, address1, address2, address3, city, state, zipcode, country) values ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = d.pg.Exec(context.Background(), command, r.ID, l.Address1, l.Address2, l.Address3, l.City, l.State, l.ZipCode, l.Country)
	if err != nil {
		return err
	}
	for _, o := range r.Hours[0].Open {
		command = `insert into hours (restaurantid, isovernight, day, starttime, endtime) values ($1, $2, $3, $4, $5)`
		_, err = d.pg.Exec(context.Background(), command, r.ID, o.IsOvernight, o.Day, o.Start, o.End)
		if err != nil {
			return err
		}
	}
	for _, p := range r.Photos {
		command = `insert into photos (restaurantid, imageurl) values ($1, $2)`
		_, err = d.pg.Exec(context.Background(), command, r.ID, p)
		if err != nil {
			return err
		}
	}
	for _, p := range r.Categories {
		command = `insert into categories (restaurantid, title) values ($1, $2)`
		_, err = d.pg.Exec(context.Background(), command, r.ID, p.Title)
		if err != nil {
			return err
		}
	}
	if username != "admin" {
		command = `update users set level = $1 where email = $2`
		_, err = d.pg.Exec(context.Background(), command, 2, username)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Database) SelectRestaurant(email string) (BusinessResponse, error) {
	var e string
	var r BusinessResponse
	command := `select * from restaurants where managerid = $1`
	err := d.pg.QueryRow(context.Background(), command, email).Scan(&r.ID, &e, &r.Name, &r.Phone, &r.Price, &r.Rating, &r.URL, &r.ReviewCount, &r.ImageURL)
	if err != nil {
		return r, err
	}
	command = `select * from locations where restaurantid = $1`
	var n int
	err = d.pg.QueryRow(context.Background(), command, r.ID).Scan(&n, &r.ID, &r.Location.Address1, &r.Location.Address2, &r.Location.Address3, &r.Location.City, &r.Location.State, &r.Location.Country, &r.Location.ZipCode)
	if err != nil {
		return r, err
	}
	command = `select title from categories where restaurantid = $1`
	rows, err := d.pg.Query(context.Background(), command, r.ID)
	if err != nil {
		return r, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Category
		rows.Scan(&c.Title, &c.Alias)
		r.Categories = append(r.Categories, c)
	}
	command = `select imageurl from photos where restaurantid = $1`
	rows, err = d.pg.Query(context.Background(), command, r.ID)
	if err != nil {
		return r, err
	}
	defer rows.Close()

	for rows.Next() {
		var p string
		rows.Scan(&p)
		r.Photos = append(r.Photos, p)
	}
	command = `select isovernight, day, starttime, endtime from hours where restaurantid = $1`
	rows, err = d.pg.Query(context.Background(), command, r.ID)
	if err != nil {
		return r, err
	}
	defer rows.Close()

	var opens []Open
	for rows.Next() {
		var o Open
		rows.Scan(&o.IsOvernight, &o.Day, &o.Start, &o.End)
		opens = append(opens, o)
	}
	var hours []Hour
	var hour Hour
	hour.Open = opens
	hours = append(hours, hour)
	r.Hours = hours

	return r, nil
}

/*func (d *Database) InsertCategory(title string) error {
	command := `SELECT column_name FROM information_schema.columns WHERE table_name='categories' and column_name=$1`
	_, err := d.pg.Exec(context.Background(), command, title)
	if err.Error() == "no rows in result set" {
		command = `alter table categories add columns ` + title + ` boolean`
		_, err := d.pg.Exec(context.Background(), command, title)
		if err != nil {
			return err
		}
		command = `insert into categories (restaurantid, ` + title + `) values ($1, $2)`
		_, err = d.pg.Exec(context.Background(), command, r.ID, p)
		if err != nil {
			return err
		}
	}
	return nil
}*/
func (d *Database) DeleteRequest(firstname string, email string) error {
	command := `delete from restaurantrequests where firstname = $1 and email = $2`
	_, err := d.pg.Exec(context.Background(), command, firstname, email)
	if err != nil {
		return err
	}
	return nil
}
func (d *Database) InsertUser(user User) error {
	command := `INSERT INTO users (Username, level, Email, Password, firstname, lastname, phone) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := d.pg.Exec(context.Background(), command, user.Username, user.Level, user.Email, user.Password, user.FirstName, user.LastName, user.PhoneNumber)
	if err != nil {
		return err
	}
	return nil
}
func (d *Database) InsertRequest(r RestaurantRequest) error {
	command := `INSERT INTO restaurantrequests (FirstName, LastName, RestaurantName, Email, Phone, City, State, ZipCode, Country, Active, Enabled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := d.pg.Exec(context.Background(), command, r.FirstName, r.LastName, r.RestaurantName, r.Email, r.PhoneNumber, r.City, r.State, r.ZipCode, r.Country, r.Active, r.Enabled)
	if err != nil {
		return err
	}
	return nil
}
func (d *Database) InsertIds(ids []string) error {
	command := `INSERT INTO businessids (id) VALUES ($1) ON CONFLICT (id) do NOTHING`
	for _, id := range ids {
		_, err := d.pg.Exec(context.Background(), command, id)
		if err != nil {
			return err
		}
	}
	return nil
}
