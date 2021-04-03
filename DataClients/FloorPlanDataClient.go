package dataclients

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type FloorPlan struct {
	SceneByteArray []byte `json:"scenebytearray"`
}

func (d *Database) PostFloorPlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var fl FloorPlan
	id := r.URL.Query().Get("id")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &fl)
	if err != nil {
		log.Fatal(err)
	} else {
		command :=
			`update restaurants set 
		floorplan = $1
		where restaurantid = $2`
		_, err = d.pg.Exec(context.Background(), command, fl.SceneByteArray, id)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusPreconditionFailed)
		}
	}
}

func (d *Database) GetFloorPlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	command := `SELECT floorplan
	from restaurants where restaurantid = $1`
	var fl FloorPlan
	err := d.pg.QueryRow(context.Background(), command, id).Scan(&fl.SceneByteArray)
	if err != nil {
		log.Fatal(err, id)
		w.WriteHeader(http.StatusPreconditionFailed)
	} else {
		json.NewEncoder(w).Encode(fl)
	}
}
