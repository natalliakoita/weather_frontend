package apiclient

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HTTPClientMock struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (H HTTPClientMock) Do(r *http.Request) (*http.Response, error) {
	return H.DoFunc(r)
}

func TestApiWeather_GetWheater(t *testing.T) {
	type fields struct {
		Host string
	}
	type args struct {
		city     string
	}
	tests := []struct {
		name       string
		fields     fields
		Body       string
		StatusCode int
		args       args
		want       *GetWeatherResponse
		wantErr    bool
	}{
		{
			name:   "succses call",
			fields: fields{Host: "test"},
			Body: `{
						"city": "Slonim",
						"dt": "2021-08-23T21:17:26+03:00",
						"temperature": 289.11
					}`,
			StatusCode: 200,
			want:       &GetWeatherResponse{},
			wantErr:    false,
		},
		{
			name: "some error",
			fields: fields{
				Host: "test",
			},
			StatusCode: 500,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &HTTPClientMock{}
			client.DoFunc = func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					Body:       io.NopCloser(strings.NewReader(tt.Body)),
					StatusCode: tt.StatusCode,
				}, nil
			}

			a := ApiWeather{
				Client: client,
				Host:   tt.fields.Host,
			}

			got, err := a.GetWheater(tt.args.city)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got.TimeStamp)
			}
		})
	}
}

func TestApiWeather_GetListWeatherRequest(t *testing.T) {
	type fields struct {
		Host string
	}
	tests := []struct {
		name       string
		fields     fields
		Body       string
		StatusCode int
		want       *WeatherListResponse
		wantErr    bool
	}{
		{
			name: "no error",
			fields: fields{
				Host: "test"},
			Body: `{
				"cities": [
					{
					  "city": "city",
					  "dt": "2021-08-15T23:31:55.047518Z",
					  "temperature": 2
					},
					{
					  "city": "Barcelona",
					  "dt": "2021-08-16T23:30:04Z",
					  "temperature": 297.88
					},
					{
					  "city": "Slonim",
					  "dt": "2021-08-23T21:17:26Z",
					  "temperature": 289.11
					}
				]
					}`,
			StatusCode: 200,
			// args:       args{},
			want:       &WeatherListResponse{},
			wantErr:    false,
		},
		{
			name: "some error",
			fields: fields{
				Host: "test",
			},
			StatusCode: 500,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &HTTPClientMock{}
			client.DoFunc = func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					Body:       io.NopCloser(strings.NewReader(tt.Body)),
					StatusCode: tt.StatusCode,
				}, nil
			}

			a := ApiWeather{
				Client: client,
				Host:   tt.fields.Host,
			}

			got, err := a.GetListWeatherRequest()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got.Cities)
			}
		})
	}
}
