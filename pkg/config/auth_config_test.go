package config

import (
	"reflect"
	"testing"
)

func TestReadAuthConfigFile(t *testing.T) {
	configInvalidLocation := "../../configs/no.config.yaml"
	configInvalidConfigFileLocation := "../../configs/bad.yaml"
	configSampleLocation := "../../configs/sample.yaml"
	configMultipleUserLocation := "../../configs/multiple.user.yaml"
	expectedSampleAuth := AuthenticationConfig{
		[]User{
			{
				"Grafana",
				"Loki",
				"tenant-1",
			},
		},
	}
	expectedMultipleUserAuth := AuthenticationConfig{
		[]User{
			{
				"User-a",
				"pass-a",
				"tenant-a",
			},
			{
				"User-b",
				"pass-b",
				"tenant-b",
			},
		},
	}
	type args struct {
		location string
	}
	tests := []struct {
		name    string
		args    args
		want    *AuthenticationConfig
		wantErr bool
	}{
		{
			"Basic",
			args{
				configSampleLocation,
			},
			&expectedSampleAuth,
			false,
		}, {
			"Multiples users",
			args{
				configMultipleUserLocation,
			},
			&expectedMultipleUserAuth,
			false,
		}, {
			"Invalid location",
			args{
				configInvalidLocation,
			},
			nil,
			true,
		}, {
			"Invalid yaml file",
			args{
				configInvalidConfigFileLocation,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readAuthConfigFile(tt.args.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAuthConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readAuthConfigFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
