package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/natalliakoita/weather_frontend/service"
	"github.com/stretchr/testify/assert"
)

type FakeApiSVC struct {
	MockGetWheaterFn            func(string) (service.Weather, error)
	MockGetListWeatherRequestFn func() ([]service.Weather, error)
}

func (fake *FakeApiSVC) GetWheater(city string) (service.Weather, error) {
	return fake.MockGetWheaterFn(city)
}

func (fake *FakeApiSVC) GetListWeatherRequest() ([]service.Weather, error) {
	return fake.MockGetListWeatherRequestFn()
}

func MockWeatherModel() service.Weather {
	var w service.Weather
	w.ID = 1
	w.City = "Minsk"
	w.TimeStamp = time.Now()
	w.Temperature = 100
	return w
}

func TestApiHandler_GetWeatherByCity(t *testing.T) {
	type args struct {
		city        string
		testWeather service.Weather
		apiErr      error
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		statusCode int
	}{
		{
			name: "succses call",
			args: args{
				city:        "Minsk",
				testWeather: MockWeatherModel(),
				apiErr:      nil,
			},
			wantErr:    false,
			statusCode: http.StatusOK,
		},
		{
			name: "some error api",
			args: args{
				city:   "Minsk",
				apiErr: errors.New("a some err"),
			},
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "bad request",
			args: args{
				city: "",
			},
			wantErr:    true,
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &FakeApiSVC{
				MockGetWheaterFn: func(string) (service.Weather, error) { return tt.args.testWeather, tt.args.apiErr },
			}
			u := &ApiHandler{
				apiSvc: fake,
			}

			path := fmt.Sprintf("/api/v0/%s/weather", tt.args.city)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			vars := map[string]string{}
			if tt.args.city != "" {
				vars["city"] = tt.args.city
			}

			req = mux.SetURLVars(req, vars)

			handler := http.HandlerFunc(u.GetWeatherByCity)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, tt.statusCode)
			b, err := getBody(tt.args.testWeather)
			assert.NoError(t, err)
			if !tt.wantErr {
				assert.Equal(t, b, rr.Body.Bytes())
			}
		})
	}
}

func getBody(w service.Weather) ([]byte, error) {
	result := GetWeatherResponse{
		City:        w.City,
		TimeStamp:   w.TimeStamp,
		Temperature: w.Temperature,
	}
	b, err := json.Marshal(result)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func TestApiHandler_WeatherListRequest(t *testing.T) {
	type args struct {
		weatherList []service.Weather
		apiErr      error
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		statusCode int
	}{
		{
			name: "succses call",
			args: args{
				weatherList: []service.Weather{
					MockWeatherModel(),
				},
				apiErr: nil,
			},
			wantErr:    false,
			statusCode: http.StatusOK,
		},
		{
			name: "some error api",
			args: args{
				apiErr: errors.New("a some err"),
			},
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &FakeApiSVC{
				MockGetListWeatherRequestFn: func() ([]service.Weather, error) { return tt.args.weatherList, tt.args.apiErr },
			}
			u := &ApiHandler{
				apiSvc: fake,
			}

			path := "/api/v0/weather"
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()


			handler := http.HandlerFunc(u.WeatherListRequest)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, tt.statusCode)
			b, err := getArrayBody(tt.args.weatherList)
			assert.NoError(t, err)
			if !tt.wantErr {
				assert.Equal(t, b, rr.Body.Bytes())
			}
		})
	}
}

func getArrayBody(w []service.Weather) ([]byte, error) {
	response := []GetWeatherResponse{}
	for _, city := range w {
		q := GetWeatherResponse{}
		q.City = city.City
		q.TimeStamp = city.TimeStamp
		q.Temperature = city.Temperature

		response = append(response, q)
	}

	b, err := json.Marshal(response)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
