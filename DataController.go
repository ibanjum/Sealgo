package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (d *Database) IdsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ids []string
	var response string
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &ids)
	if err != nil {
		log.Fatal(err)
	}
	err = d.InsertIds(ids)
	if err != nil {
		response = "Error while inserterting ids"
		json.NewEncoder(w).Encode(response)
		return
	}
	response = "Error while inserterting restaurant"
	json.NewEncoder(w).Encode(response)
	return
}
func (d *Database) GetRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	requests, err := d.SelectAllRequests()
	if err != nil {
		e := "Error while selecting requests"
		json.NewEncoder(w).Encode(e)
		return
	}
	json.NewEncoder(w).Encode(requests)
}
func (d *Database) fzPostRestaurant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var br BusinessResponse
	var response string
	vars := mux.Vars(r)
	userid := vars["userid"]
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &br)
	if err != nil {
		log.Fatal(err)
	}
	err = d.InsertRestaurant(br, userid)
	if err != nil {
		response = "Error while inserterting restaurant"
		json.NewEncoder(w).Encode(response)
		return
	}
	response = "Successful"
	json.NewEncoder(w).Encode(response)
	return
}
