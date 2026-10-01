package config_test

import (
	"testing"
	"time"

	"mortgage-api-bfcl-mortgage-loans-catalogs/cmd/config"

	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    func(t *testing.T, c *config.APIConfig)
		wantErr bool
	}{
		{
			name: "defaults con external_service_url definido",
			env:  map[string]string{"EXTERNAL_SERVICE_URL": "http://localhost:9999"},
			want: func(t *testing.T, c *config.APIConfig) {
				assert.Equal(t, "/fifcl/v1", c.URIPrefix)
				assert.Equal(t, "mortgage-api-bfcl-mortgage-loans-catalogs", c.AppName)
				assert.Equal(t, "dev", c.Env)
				assert.Equal(t, "0.0.0", c.Version)
				assert.Equal(t, "CL", c.Country)
				assert.Equal(t, "debug", c.GinMode)
				assert.Equal(t, loggerFif.InfoLevel, c.LoggingLevel)
				assert.Equal(t, 10*time.Second, c.Timeout)
				assert.Equal(t, "localhost:8126", c.DataDog.AgentAddr)
				assert.False(t, c.DataDog.ProfilerEnabled)
			},
		},
		{
			name: "override de variables y tipos",
			env: map[string]string{
				"EXTERNAL_SERVICE_URL": "http://localhost:9999",
				"URI_PREFIX":           "/fifpe/v2",
				"TIMEOUT":              "5",
				"DD_PROFILE_ENABLED":   "true",
				"GIN_MODE":             "RELEASE",
			},
			want: func(t *testing.T, c *config.APIConfig) {
				assert.Equal(t, "/fifpe/v2", c.URIPrefix)
				assert.Equal(t, 5*time.Second, c.Timeout)
				assert.True(t, c.DataDog.ProfilerEnabled)
				assert.Equal(t, "release", c.GinMode)
			},
		},
		{
			name:    "external_service_url vacío retorna error",
			env:     map[string]string{},
			wantErr: true,
		},
		{
			name: "logging_level inválido retorna error",
			env: map[string]string{
				"EXTERNAL_SERVICE_URL": "http://localhost:9999",
				"LOGGING_LEVEL":        "VERBOSE",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got, err := config.Load()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.want(t, got)
		})
	}
}
