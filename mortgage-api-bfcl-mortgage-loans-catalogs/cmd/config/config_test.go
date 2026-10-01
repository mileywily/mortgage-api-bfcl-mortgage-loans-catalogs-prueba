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
			name: "defaults con Finnflow configurado",
			env:  map[string]string{"FINNFLOW_URL": "http://localhost:9999"},
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
				assert.Equal(t, "real", c.DefaultBackend)
				assert.Equal(t, "http://localhost:9999", c.FinnflowURL)
			},
		},
		{
			name: "override de variables y tipos",
			env: map[string]string{
				"FINNFLOW_URL":       "http://localhost:9999",
				"FINNFLOW_KEY":       "client",
				"FINNFLOW_SECRET":    "secret",
				"DEFAULT_BACKEND":    "java",
				"JAVA_LEGACY_URL":    "http://localhost:9998",
				"URI_PREFIX":         "/fifpe/v2",
				"TIMEOUT":            "5",
				"DD_PROFILE_ENABLED": "true",
				"GIN_MODE":           "RELEASE",
			},
			want: func(t *testing.T, c *config.APIConfig) {
				assert.Equal(t, "/fifpe/v2", c.URIPrefix)
				assert.Equal(t, 5*time.Second, c.Timeout)
				assert.True(t, c.DataDog.ProfilerEnabled)
				assert.Equal(t, "release", c.GinMode)
				assert.Equal(t, "java", c.DefaultBackend)
				assert.Equal(t, "client", c.FinnflowKey)
				assert.Equal(t, "secret", c.FinnflowSecret)
				assert.Equal(t, "http://localhost:9998", c.JavaLegacyURL)
			},
		},
		{
			name:    "Finnflow URL requerida para backend real",
			env:     map[string]string{"DEFAULT_BACKEND": "real"},
			wantErr: true,
		},
		{
			name: "logging_level inválido retorna error",
			env: map[string]string{
				"FINNFLOW_URL":  "http://localhost:9999",
				"LOGGING_LEVEL": "VERBOSE",
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
