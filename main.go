package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/natalliakoita/weather_frontend/apiclient"
	"github.com/natalliakoita/weather_frontend/handlers"
	"github.com/natalliakoita/weather_frontend/service"
)

func main() {
	c := http.Client{Timeout: time.Duration(40) * time.Second}
	a := apiclient.NewApiWeather(&c)

	apiSvc := service.NewApiService(a)

	h := handlers.NewApiHandler(&apiSvc)

	router := mux.NewRouter()
	router.HandleFunc("/api/v0/{city}/weather", h.GetWeatherByCity).Methods(http.MethodGet)
	router.HandleFunc("/api/v0/weather", h.WeatherListRequest).Methods(http.MethodGet)

	log.Println("Starting API server on 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatal(err)
	}
}
