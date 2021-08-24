package service

import (
	"errors"
	"testing"
	"time"

	"github.com/natalliakoita/weather_frontend/apiclient"
	"github.com/stretchr/testify/assert"
)

type FakeApiclient struct {
	MockGetWheaterFn            func(city string) (*apiclient.GetWeatherResponse, error)
	MockGetListWeatherRequestFn func() (*apiclient.WeatherListResponse, error)
}

func (fake *FakeApiclient) GetWheater(city string) (*apiclient.GetWeatherResponse, error) {
	return fake.MockGetWheaterFn(city)
}

func (fake *FakeApiclient) GetListWeatherRequest() (*apiclient.WeatherListResponse, error) {
	return fake.MockGetListWeatherRequestFn()
}

func TestApiService_GetWheater(t *testing.T) {
	type args struct {
		w         *apiclient.GetWeatherResponse
		city      string
		respError error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "succses call",
			args: args{
				city: "Minsk",
				w: &apiclient.GetWeatherResponse{
					ID:          1,
					City:        "Minsk",
					TimeStamp:   time.Now(),
					Temperature: 298,
				},
			},
			wantErr: false,
		},
		{
			name: "call with a some error",
			args: args{
				respError: errors.New("error"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &FakeApiclient{
				MockGetWheaterFn: func(string) (*apiclient.GetWeatherResponse, error) { return tt.args.w, tt.args.respError },
			}
			dbs := &ApiService{
				DS: fake,
			}
			got, err := dbs.GetWheater(tt.args.city)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.IsType(t, Weather{}, got)
				assert.NoError(t, err)
			}
		})
	}
}

func TestApiService_GetListWeatherRequest(t *testing.T) {
	type args struct {
		w         *apiclient.WeatherListResponse
		city      string
		respError error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		lenResp int
	}{
		{
			name: "call when len(resp) > 0",
			args: args{
				city: "Minsk",
				w: &apiclient.WeatherListResponse{
					Cities: []apiclient.GetWeatherResponse{
						{
							ID:          1,
							City:        "Minsk",
							TimeStamp:   time.Now(),
							Temperature: 298,
						},
					},
				},
			},
			wantErr: false,
			lenResp: 1,
		},
		{
			name: "call with a some error",
			args: args{
				respError: errors.New("error"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &FakeApiclient{
				MockGetListWeatherRequestFn: func() (*apiclient.WeatherListResponse, error) { return tt.args.w, tt.args.respError },
			}
			dbs := &ApiService{
				DS: fake,
			}
			got, err := dbs.GetListWeatherRequest()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Len(t, got, tt.lenResp)
				assert.NoError(t, err)
			}
		})
	}
}
