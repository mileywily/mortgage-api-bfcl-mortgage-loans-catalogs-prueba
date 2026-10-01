package get

import "encoding/json"

type RequestDto struct {
	Header HeaderDto
	Uri    UriDto
}

type HeaderDto struct {
	Channel       string
	Commerce      string
	TransactionID string
	Backend       string
}

type UriDto struct {
	CatalogName string
}

type catalogItemDTO struct {
	CodigoAdm   string `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
}

type comunaCatalogItemDTO struct {
	CodigoAdm   string `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
	RegionID    int    `json:"region_id"`
}

type tipoDocumentoCatalogItemDTO struct {
	CodigoAdm   int    `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
	GrupoID     int    `json:"grupo_id"`
}

type noResultsCatalogItemDTO struct {
	CodRespuesta int    `json:"codRespuesta"`
	Mensaje      string `json:"Mensaje"`
	Excepcion    string `json:"Excepcion"`
}

type insuranceCatalogItemDTO struct {
	IdentificadorSeguro       string      `json:"IdentificadorSeguro"`
	Descripcion               string      `json:"Descripcion"`
	Poliza                    string      `json:"Poliza"`
	NombreCompania            string      `json:"NombreCompania"`
	CodigoCompania            int         `json:"CodigoCompania"`
	CodigoTipoSeguro          int         `json:"CodigoTipoSeguro"`
	CorrelativoPoliza         int         `json:"CorrelativoPoliza"`
	Tasa                      json.Number `json:"Tasa"`
	Factor                    json.Number `json:"Factor"`
	IndicadorPolizaIndividual int         `json:"IndicadorPolizaIndividual"`
	PorValorCuota             int         `json:"PorValorCuota"`
	IndicadorPolizaExterna    int         `json:"IndicadorPolizaExterna"`
}
