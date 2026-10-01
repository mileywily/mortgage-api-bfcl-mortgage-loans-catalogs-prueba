package catalog_service

// ——— Finnflow upstream DTOs ———

// FinnflowInsuranceItem es el struct de respuesta de Finnflow para catálogos de seguros.
type FinnflowInsuranceItem struct {
	IdentificadorSeguro       string  `json:"IdentificadorSeguro"`
	Descripcion               string  `json:"Descripcion"`
	Poliza                    string  `json:"Poliza"`
	NombreCompania            string  `json:"NombreCompania"`
	CodigoCompania            int     `json:"CodigoCompania"`
	CodigoTipoSeguro          int     `json:"CodigoTipoSeguro"`
	CorrelativoPoliza         int     `json:"CorrelativoPoliza"`
	Tasa                      float64 `json:"Tasa"`
	Factor                    float64 `json:"Factor"`
	IndicadorPolizaIndividual int     `json:"IndicadorPolizaIndividual"`
	PorValorCuota             int     `json:"PorValorCuota"`
	IndicadorPolizaExterna    int     `json:"IndicadorPolizaExterna"`
}

// FinnflowTipoDocumentoItem es el struct de respuesta de Finnflow para tipos de documento.
type FinnflowTipoDocumentoItem struct {
	CodigoAdm   int    `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
	GrupoID     int    `json:"grupo_id"`
}

// FinnflowComunaItem es el struct de respuesta de Finnflow para comunas.
type FinnflowComunaItem struct {
	CodigoAdm   interface{} `json:"codigo_adm"`
	Descripcion string      `json:"descripcion"`
	RegionID    int         `json:"region_id"`
}

// FinnflowGenericItem es el struct de respuesta de Finnflow para catálogos estándar.
type FinnflowGenericItem struct {
	CodigoAdm    interface{} `json:"codigo_adm"`
	Descripcion  string      `json:"descripcion"`
	CodRespuesta *int        `json:"codRespuesta"`
}
