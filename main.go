package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"math/rand"
	"net/http"
	"os"
	"time"
)

//no of balls faced or bowled in test cricket - used to calculate the
const (
	DRAVID int = 31258
	SACHIN int = 29437
	KALLIS int = 28903
	MURALI int = 44039
	KUMBLE int = 40850
	WARNE  int = 40705
)

var total_balls = []int{DRAVID, SACHIN, KALLIS, MURALI, KUMBLE, WARNE}

var falcones map[string]int = make(map[string]int)

var cities = Cities{City{"Mysore", 200}, City{"Chennai", 400}, City{"Hyderabad", 600}, City{"Pune", 800}, City{"Mumbai", 1000}, City{"Ahmedabad", 1200}}

var vehicles = TransportVehicles{TransportVehicle{"Car", 2, 400, 50}, TransportVehicle{"Bus", 1, 600, 100}, TransportVehicle{"Train", 1, 800, 200}, TransportVehicle{"Plane", 2, 1200, 400}}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

///returns a token for an N integer
func randSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//returns a random no out of 0-6
func where_is_falcone() int {
	rand.Seed(time.Now().UTC().UnixNano())
	var no = rand.Intn(6)
	return (total_balls[no] * rand.Intn(10)) % 6
}

//returns a token for a user who is trying to find falcone.
func Init(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	rw.WriteHeader(http.StatusOK)
	var random_str = randSeq(32)
	var falcone_city = where_is_falcone()
	falcones[random_str] = falcone_city
	var token = map[string]string{"token": random_str}
	if err := json.NewEncoder(rw).Encode(token); err != nil {
		panic(err)
	}
}

//returns all the cities
func CitiesHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(cities); err != nil {
		panic(err)
	}
}

//returns all the vehicles
func VehicleHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(vehicles); err != nil {
		panic(err)
	}
}

//API to find the falcone
func FindFalcone(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	decoder := json.NewDecoder(req.Body)
	var find_falcone FindFalconeReq
	err := decoder.Decode(&find_falcone)
	if err != nil {
		errorHandler(rw, req, http.StatusBadRequest, err.Error())
		return
	}
	var cityNames = find_falcone.CityNames
	if len(cityNames) != 4 {
		errorHandler(rw, req, http.StatusBadRequest, "No of City names has to be 4")
		return
	}
	var vehicleNames = find_falcone.TransportVehicleNames
	if len(vehicleNames) != 4 {
		errorHandler(rw, req, http.StatusBadRequest, "No of Vehicle names has to be 4")
		return
	}
	if len(falcones) == 0 {
		errorHandler(rw, req, http.StatusBadRequest, "Token not initialized. Please get a new token with the /token API")
		return
	}

	// var falconeCityIndex = falcones[find_falcone.Token]
	if falconeCityIndex, ok := falcones[find_falcone.Token]; ok {
		rw.WriteHeader(http.StatusOK)
		var falconeCity = cities[falconeCityIndex]
		for _, name := range cityNames {
			if name == falconeCity.Name {
				var status = map[string]string{"status": "success", "city_name": name}
				if err := json.NewEncoder(rw).Encode(status); err != nil {
					panic(err)
				}
				return
			}
		}
	} else {
		errorHandler(rw, req, http.StatusBadRequest, "Token not initialized. Please get a new token with the /token API")
		return
	}
	var status = map[string]string{"status": "false"}
	if err := json.NewEncoder(rw).Encode(status); err != nil {
		panic(err)
	}
}

func errorHandler(rw http.ResponseWriter, req *http.Request, status int, message string) {
	rw.WriteHeader(status)
	var error = map[string]string{"error": message}
	if err := json.NewEncoder(rw).Encode(error); err != nil {
		panic(err)
	}
}

func main() {
	r := mux.NewRouter()
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("Starting server on " + port)
	r.HandleFunc("/token", Init).Methods("POST").Headers("Accept", "application/json")
	r.HandleFunc("/cities", CitiesHandler).Methods("GET")
	r.HandleFunc("/vehicles", VehicleHandler).Methods("GET")
	r.HandleFunc("/find", FindFalcone).Methods("POST").Headers("Accept", "application/json")
	// c := cors.New(cors.Options{
	// 	AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
	// })
	handler := cors.Default().Handler(r)
	// handler := c.Handler(r)
	http.ListenAndServe(":"+port, handler)
}
