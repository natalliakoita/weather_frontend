package service

import (
	"time"

	"github.com/natalliakoita/weather_frontend/apiclient"
)

type ApiService struct {
	DS apiclient.ApiWeatherInterface
}

func NewApiService(d apiclient.ApiWeatherInterface) ApiService {
	return ApiService{
		DS: d,
	}
}

type Weather struct {
	ID          int
	City        string
	TimeStamp   time.Time
	Temperature float32
}

func (dbs *ApiService) GetWheater(city string) (Weather, error) {
	resp, err := dbs.DS.GetWheater(city)
	if err != nil {
		return Weather{}, err
	}
	var m Weather
	m.City = resp.City
	m.TimeStamp = resp.TimeStamp
	m.Temperature = resp.Temperature

	return m, nil
}

type ApiServiceInterface interface {
	GetWheater(city string) (Weather, error)
}
