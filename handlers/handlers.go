package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/natalliakoita/weather_frontend/service"
)

type ApiHandler struct {
	apiSvc service.ApiServiceInterface
}

func NewApiHandler(a service.ApiServiceInterface) ApiHandler {
	cont := ApiHandler{
		apiSvc: a,
	}
	return cont
}

func (u *ApiHandler) GetWeatherByCity(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	city, ok := vars["city"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := u.apiSvc.GetWheater(city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := GetWeatherResponse{
		ID:          resp.ID,
		City:        resp.City,
		TimeStamp:   resp.TimeStamp,
		Temperature: resp.Temperature,
	}

	tmpl := template.Must(template.ParseFiles("form.html"))
	tmpl.Execute(w, result)

	err = result.writeToWeb(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type GetWeatherResponse struct {
	ID          int
	City        string
	TimeStamp   time.Time
	Temperature float32
}

func (c GetWeatherResponse) writeToWeb(w http.ResponseWriter) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}
