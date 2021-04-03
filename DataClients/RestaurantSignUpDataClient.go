package dataclients

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type RestaurantBasicInfo struct {
	Name     string `json:"restaurantname"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	State    string `json:"state"`
	ZipCode  string `json:"zipcode"`
	Country  string `json:"country"`
	Image []byte `json:"image"`
	FloorPlan []byte `json:"floorplan"`
}

type DetailRestaurant struct {
	Name     string `json:"restaurantname"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	State    string `json:"state"`
	ZipCode  string `json:"zipcode"`
	Country  string `json:"country"`
	Image []byte `json:"image"`
	Items []MenuItemModel `json:"items"`
	Cusines []string `json:"cusines"`
}
type AdditionalInfo struct {
	MainImage []byte   `json:"mainImage"`
	Cusines   []string `json:"cusines"`
}

type BusinessHoursModel struct {
	Day    string `json:"day"`
	Opens  string `json:"opens"`
	Closes string `json:"closes"`
}
type BusinessHoursModels struct {
	Hours []BusinessHoursModel `json:"hours"`
}

type MenuItemsModel struct {
	Items []MenuItemModel `json:"items"`
}

type MenuItemModel struct {
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Image       []byte `json:"image"`
}

func (d *Database) RestaurantBasicInfoGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	command := `SELECT address1, address2, city, state, zipcode, country, restaurantname FROM restaurants where restaurantid = $1`
	var bs RestaurantBasicInfo
	err := d.pg.QueryRow(context.Background(), command, id).Scan(&bs.Address1, &bs.Address2, &bs.City, &bs.State, &bs.ZipCode, &bs.Country, &bs.Name)
	if err != nil {
		log.Fatal(err, id)
		w.WriteHeader(http.StatusPreconditionFailed)
	} else {
		json.NewEncoder(w).Encode(bs)
	}
}

func (d *Database) PostRestaurantBasicInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var br RestaurantBasicInfo
	id := r.URL.Query().Get("id")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &br)
	if err != nil {
		log.Fatal(err)
	} else {
		command :=
			`update restaurants set 
		address1 = $1, 
		address2 = $2, 
		city = $3, 
		state = $4, 
		zipcode = $5, 
		country = $6
		where restaurantid = $7`
		_, err = d.pg.Exec(context.Background(), command, br.Address1, br.Address2, br.City, br.State, br.ZipCode, br.Country, id)

		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusPreconditionFailed)
		}
	}
}

func (d *Database) PostRestaurantAdditionInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var AI AdditionalInfo
	id := r.URL.Query().Get("id")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &AI)
	if err != nil {
		log.Fatal(err)
	} else {
		command :=
			`update restaurants set 
			image = $1 
			where restaurantid = $2`
		_, err = d.pg.Exec(context.Background(), command, AI.MainImage, id)

		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusPreconditionFailed)
		}

		command =
			`delete from cusines where restaurantid = $1`
		_, err = d.pg.Exec(context.Background(), command, id)

		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusPreconditionFailed)
		}

		for _, cusine := range AI.Cusines {
			command =
				`insert into cusines  (restaurantid, name)
				values ($1, $2) ON CONFLICT do NOTHING`
			_, err = d.pg.Exec(context.Background(), command, id, cusine)

			if err != nil {
				log.Fatal(err)
				w.WriteHeader(http.StatusPreconditionFailed)
			}
		}
	}
}

func (d *Database) PostRestaurantHours(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bhs BusinessHoursModels
	id := r.URL.Query().Get("id")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &bhs)
	if err != nil {
		log.Fatal(err)
	} else {

		command :=
			`delete from business_hours  where restaurantid = $1`
		_, err = d.pg.Exec(context.Background(), command, id)

		if err != nil {
			log.Fatal(err, bhs)
			w.WriteHeader(http.StatusPreconditionFailed)
		}

		for _, bh := range bhs.Hours {
			command :=
				`insert into business_hours  (restaurantid, day, open_time, close_time)
			values ($1, $2, $3, $4) ON CONFLICT do NOTHING`
			_, err = d.pg.Exec(context.Background(), command, id, bh.Day, bh.Opens, bh.Closes)

			if err != nil {
				log.Fatal(err, bhs)
				w.WriteHeader(http.StatusPreconditionFailed)
			}
		}
	}
}

func (d *Database) PostRestaurantMenuItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var menu MenuItemsModel
	id := r.URL.Query().Get("id")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &menu)
	if err != nil {
		log.Fatal(err)
	} else {
		for _, item := range menu.Items {
			command :=
				`insert into restaurant_menu_items  (restaurantid, name, price, description, category, image)
			values ($1, $2, $3, $4, $5, $6) ON CONFLICT do NOTHING`
			_, err = d.pg.Exec(context.Background(), command, id, item.Name, item.Price, item.Description, item.Category, item.Image)

			if err != nil {
				log.Fatal(err)
				w.WriteHeader(http.StatusPreconditionFailed)
			}
		}
	}
}

func (d *Database) GetDetailRestaurant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	var dr DetailRestaurant

	command := `SELECT address1, address2, city, state, zipcode, country, restaurantname, image FROM restaurants where restaurantid = $1`
	err := d.pg.QueryRow(context.Background(), command, id).Scan(&dr.Address1, &dr.Address2, &dr.City, &dr.State, &dr.ZipCode, &dr.Country, &dr.Name, &dr.Image)
	if err != nil {
		log.Fatal(err, id)
		w.WriteHeader(http.StatusPreconditionFailed)
	}

	command = `select name
	from cusines 
	where restaurantid = $1`
	rows, err := d.pg.Query(context.Background(), command, id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusPreconditionFailed)
	}
	defer rows.Close()

	for rows.Next() {
		var cusine string
		rows.Scan(&cusine)
		dr.Cusines = append(dr.Cusines, cusine)
	}

	command = `select name, price, description, category, image 
	from restaurant_menu_items 
	where restaurantid = $1`
	rows, err = d.pg.Query(context.Background(), command, id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusPreconditionFailed)
	}
	defer rows.Close()

	for rows.Next() {
		var mi MenuItemModel
		rows.Scan(&mi.Name, &mi.Price, &mi.Description, &mi.Category, &mi.Image)
		dr.Items = append(dr.Items, mi)
	}

	json.NewEncoder(w).Encode(dr)
}
