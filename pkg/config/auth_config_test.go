package config

import (
	"reflect"
	"testing"
)

var (
	missingConfigLocation    = "./testdata/missing.yaml"
	invalidConfigLocation    = "./testdata/invalid.yaml"
	singleUserConfigLocation = "./testdata/single-user.yaml"
	multiUserConfigLocation  = "./testdata/multi-user.yaml"
	expectedSingleUserConfig = AuthenticationConfig{
		[]User{
			{
				"Grafana",
				"Loki",
				"tenant-1",
			},
		},
	}
	expectedMultipleUserAuth = AuthenticationConfig{
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
)

func TestReadAuthConfigFile(t *testing.T) {
	tests := []struct {
		name         string
		fileLocation string
		want         *AuthenticationConfig
		wantErr      bool
	}{
		{
			"Single user",
			singleUserConfigLocation,
			&expectedSingleUserConfig,
			false,
		}, {
			"Multiples users",
			multiUserConfigLocation,
			&expectedMultipleUserAuth,
			false,
		}, {
			"Missing config",
			missingConfigLocation,
			nil,
			true,
		}, {
			"Invalid config",
			invalidConfigLocation,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readAuthConfigFile(tt.fileLocation)
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
