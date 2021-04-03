package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	dc "main/dataclients"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

var (
	dbHost     = flag.String("db-host", "localhost", "Postgres host")
	dbPort     = flag.Int64("db-port", 5431, "Postgres port")
	dbUser     = flag.String("db-user", "postgres", "Postgres user")
	dbPass     = flag.String("db-pass", "", "Postgres password")
	dbDatabase = flag.String("db-database", "seal", "Postgres database")
	listen     = flag.String("listen", ":8885", "Listen host:port")
)

func main() {

	flag.Parse()

	pg, err := pgxpool.Connect(context.Background(), fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s pool_max_conns=10", *dbUser, *dbPass, *dbHost, *dbPort, *dbDatabase))

	if err != nil {
		log.Fatal("failed to connect to db", zap.Error(err))
		return
	}
	defer pg.Close()

	pgdb := dc.NewDatabase(pg)

	/*err = pgdb.Init()
	if err != nil {
		log.Fatal("failed to initialize db", zap.Error(err))
		return
	}*/

	r := mux.NewRouter()
	/*r.HandleFunc("/register", pgdb.RegisterHandler).
		Methods("POST")
	r.HandleFunc("/login", pgdb.LoginHandler).
		Methods("POST")
	r.HandleFunc("/profile", ProfileHandler).
		Methods("GET")
	r.HandleFunc("/postids", pgdb.IdsHandler).
		Methods("POST")
	r.HandleFunc("/restaurantignup", pgdb.RegisterRestaurantHandler).
		Methods("POST")
	r.HandleFunc("/requests", pgdb.GetRequests).
		Methods("GET")
	r.HandleFunc("/addrestaurant/{userid}", pgdb.PostRestaurant).
		Methods("POST")*/

	r.HandleFunc("/GetRestaurantBasicInfo", pgdb.RestaurantBasicInfoGet).
		Methods("GET")
	r.HandleFunc("/PostRestaurantBasicInfo", pgdb.PostRestaurantBasicInfo).
		Methods("POST")
	r.HandleFunc("/PostRestaurantHours", pgdb.PostRestaurantHours).
		Methods("POST")
	r.HandleFunc("/PostRestaurantMenuItems", pgdb.PostRestaurantMenuItems).
		Methods("POST")
	r.HandleFunc("/GetDetailRestaurant", pgdb.GetDetailRestaurant).
		Methods("GET")
	r.HandleFunc("/PostRestaurantAdditionalInfo", pgdb.PostRestaurantAdditionInfo).
		Methods("POST")

	r.HandleFunc("/PostFloorPlan", pgdb.PostFloorPlan).
		Methods("POST")
	r.HandleFunc("/GetFloorPlan", pgdb.GetFloorPlan).
		Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
