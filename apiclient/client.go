package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type ApiWeather struct {
	Client HTTPClient
	Host   string
}

func NewApiWeather(conn *http.Client) ApiWeather {
	host := os.Getenv("HOST")
	return ApiWeather{
		Client: conn,
		Host:   host,
	}
}

func (a ApiWeather) GetWheater(city string) (*GetWeatherResponse, error) {
	baseUrl := fmt.Sprintf("http://%s:8080/api/v0/%s/weather", a.Host, city)
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	w := new(GetWeatherResponse)
	if err = json.Unmarshal(bodyData, w); err != nil {
		return nil, err
	}
	return w, nil
}

type ApiWeatherInterface interface {
	GetWheater(city string) (*GetWeatherResponse, error)
	GetListWeatherRequest() (*WeatherListResponse, error)
}

type GetWeatherResponse struct {
	ID          int       `json:"id,omitempty"`
	City        string    `json:"city"`
	TimeStamp   time.Time `json:"dt"`
	Temperature float32   `json:"temperature"`
}

type WeatherListResponse struct {
	Cities []GetWeatherResponse
}

func (a ApiWeather) GetListWeatherRequest() (*WeatherListResponse, error) {
	baseUrl := fmt.Sprintf("http://%s:8080/api/v0/weather", a.Host)
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	w := new(WeatherListResponse)
	if err = json.Unmarshal(bodyData, w); err != nil {
		return nil, err
	}
	return w, nil
}
