package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (d *Database) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	exists := d.UserExists(user.Username)

	if !exists {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
		if err != nil {
			res.Error = "Error While Hashing Password"
			json.NewEncoder(w).Encode(res)
			return
		}
		user.Password = string(hash)
		if user.Username == "admin" {
			user.Level = 3
		}
		err = d.InsertUser(user)
		if err != nil {
			res.Error = "Error While Creating User"
			json.NewEncoder(w).Encode(res)
			return
		}
		res.Result = "Registration Successful"
		json.NewEncoder(w).Encode(res)
		return
	} else {
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
	}
	res.Error = "Username already exists!"
	json.NewEncoder(w).Encode(res)
	return
}
func (d *Database) LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	var res ResponseResult

	result, err := d.SelectUser(user.Username)

	if err != nil && err.Error() == "no rows in result set" {
		res.Error = "Username does not exists"
		json.NewEncoder(w).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		res.Error = "Invalid password"
		json.NewEncoder(w).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": result.Username,
		"email":    result.Email,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		res.Error = "Error while generating token"
		json.NewEncoder(w).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""
	if result.Level > 0 {
		if result.Username != "admin" {
			request, err := d.SelectRequest(result.Email)
			if err != nil {
				res.Error = "Error while selecting request"
				json.NewEncoder(w).Encode(res)
				return
			}
			result.Requests = append(result.Requests, request)
		} else {
			allrequests, err := d.SelectAllRequests()
			if err != nil {
				res.Error = "Error while selecting all requests"
				json.NewEncoder(w).Encode(res)
				return
			}
			result.Requests = allrequests
		}
	}
	if result.Level > 1 {
		r, err := d.SelectRestaurant(result.Email)
		if err != nil {
			if err.Error() == "no rows in result set" {
				json.NewEncoder(w).Encode(result)
				return
			} else {
				res.Error = "Error while selecting restaurant"
				json.NewEncoder(w).Encode(res)
				return
			}
		}
		result.Restaurant = r
	}
	json.NewEncoder(w).Encode(result)
	return
}
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var result User
	var res ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Username = claims["username"].(string)
		result.Email = claims["email"].(string)

		json.NewEncoder(w).Encode(result)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}
func (d *Database) RegisterRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var restaurant RestaurantRequest
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &restaurant)
	var res ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	err = d.InsertRequest(restaurant)
	if err != nil {
		res.Error = "Error while adding restaurant request"
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Result = "Registration Successful"
	json.NewEncoder(w).Encode(res)
	return
}
