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

type ApiWeather struct {
	Client *http.Client
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

	resp, err := a.Client.Get(u.String())
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
}

type GetWeatherResponse struct {
	ID          int       `json:"id,omitempty"`
	City        string    `json:"city"`
	TimeStamp   time.Time `json:"dt"`
	Temperature float32   `json:"temperature"`
}
