package configuration

import (
	configFif "github.com/falabella-regulado/go-lib-config-fif"
	"time"
)

const (
	port           = "port"
	waitTolerance  = "wait_tolerance_server"
	writeTolerance = "write_tolerance_server"
	readTolerance  = "read_tolerance_server"
	idleTolerance  = "idle_tolerance_server"
)

var vars = []configFif.ConfigEntry{
	{VariableName: port, DefaultValue: ":8080", Description: "Puerto del servidor"},
	{VariableName: waitTolerance, DefaultValue: 30, Description: "Tolerancia de espera al cerrar el servidor"},
	{VariableName: writeTolerance, DefaultValue: 15, Description: "Tiempo máximo permitido para escribir una respuesta al cliente"},
	{VariableName: readTolerance, DefaultValue: 15, Description: "Tiempo máximo permitido para leer la solicitud del cliente"},
	{VariableName: idleTolerance, DefaultValue: 60, Description: "Tiempo máximo permitido para que una conexión permanezca inactiva"},
}

type Configuration struct {
	Port           string
	WaitTolerance  time.Duration
	WriteTolerance time.Duration
	ReadTolerance  time.Duration
	IdleTolerance  time.Duration
}

func GetConfiguration() Configuration {
	config := configFif.LoadConfig(vars)
	return Configuration{
		WaitTolerance:  time.Duration(config[waitTolerance].(int)) * time.Second,
		WriteTolerance: time.Duration(config[writeTolerance].(int)) * time.Second,
		ReadTolerance:  time.Duration(config[readTolerance].(int)) * time.Second,
		IdleTolerance:  time.Duration(config[idleTolerance].(int)) * time.Second,
		Port:           config[port].(string),
	}
}
