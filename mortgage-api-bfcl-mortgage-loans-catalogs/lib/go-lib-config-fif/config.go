package config_fif

import (
	"os"
	"strconv"
	"strings"
)

type ConfigEntry struct {
	VariableName string      // Nombre de la variable de entorno
	DefaultValue interface{} // Valor por defecto si no está definida
	Description  string      // Descripción de la variable
}

func LoadConfig(vars []ConfigEntry) map[string]interface{} {
	config := make(map[string]interface{})

	for _, v := range vars {
		value, exists := os.LookupEnv(strings.ToUpper(v.VariableName))
		if !exists {
			config[v.VariableName] = v.DefaultValue
			continue
		}

		// Intentar convertir el valor según el tipo del valor por defecto
		switch v.DefaultValue.(type) {
		case int:
			if intValue, err := strconv.Atoi(value); err == nil {
				config[v.VariableName] = intValue
			} else {
				config[v.VariableName] = v.DefaultValue
			}
		case bool:
			if boolValue, err := strconv.ParseBool(value); err == nil {
				config[v.VariableName] = boolValue
			} else {
				config[v.VariableName] = v.DefaultValue
			}
		case float64:
			if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
				config[v.VariableName] = floatValue
			} else {
				config[v.VariableName] = v.DefaultValue
			}
		default:
			config[v.VariableName] = value
		}
	}

	return config
}
