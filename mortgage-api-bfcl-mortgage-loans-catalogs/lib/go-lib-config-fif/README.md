# Config-Fif

Es una libreria que te permite configurar una lista de variables de entrada, obteniendolas desde las variables de entorno.

## Cómo usar

Primero descargamos la libreria y sus dependencias utilizando el comando:
```
go get "github.com/falabella-regulado/go-lib-config-fif"
```

Para utilizarla se debe importar la librería:

```
import (
    "github.com/falabella-regulado/go-lib-config-fif"
)
```

Crear una lista de ConfigEntry con el nombre de la variable a leer y su valor por defecto

```
type ConfigEntry struct {
    VariableName string
    DefaultValue interface{}
    Description  string
}
```

Llamar al metodo LoadConfig enviandole la lista configurada previamente y este devolvera un mapa map[string]interface{}
```Ejemplo

configVar := []configFif.ConfigEntry{
		{
			VariableName: "FIELD1",
			DefaultValue: 1,
			Description: "",
		},
		{
			VariableName: "FIELD2",
			DefaultValue: "hello",
			Description: ""
		},
        }
        
mapVar:= configFif.LoadConfig(configVar)
```