package get

import (
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

// RequestDto agrupa todos los datos de entrada del handler GET /catalogs/:catalog.
type RequestDto struct {
	Header HeaderDto
	Uri    UriDto
}

// HeaderDto contiene los headers legados obligatorios.
type HeaderDto struct {
	Channel       string `header:"X-Channel"`
	Commerce      string `header:"X-Commerce"`
	TransactionID string `header:"X-Transaction-ID"`
	BackendEnv    string `header:"X-Backend-Env"`
}

// UriDto contiene los parámetros de ruta.
type UriDto struct {
	Catalog string `uri:"catalog" binding:"required"`
}

// toDomainRequest convierte el DTO a la solicitud de dominio.
func (r RequestDto) toDomainRequest() domain.GetCatalogRequest {
	return domain.GetCatalogRequest{
		CatalogName:   r.Uri.Catalog,
		Channel:       r.Header.Channel,
		Commerce:      r.Header.Commerce,
		TransactionID: r.Header.TransactionID,
	}
}

// ——— DTOs de respuesta ———

// CatalogItemResponseDto es el payload JSON para catálogos estándar.
type CatalogItemResponseDto struct {
	CodigoAdm   string `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
}

// ComunaCatalogItemResponseDto es el payload JSON para comunas.
type ComunaCatalogItemResponseDto struct {
	CodigoAdm   string `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
	RegionID    int    `json:"region_id"`
}

// TipoDocumentoCatalogItemResponseDto es el payload JSON para tipos de documento.
type TipoDocumentoCatalogItemResponseDto struct {
	CodigoAdm   int    `json:"codigo_adm"`
	Descripcion string `json:"descripcion"`
	GrupoID     int    `json:"grupo_id"`
}

// NoResultsCatalogItemResponseDto es el fallback legacy (sin resultados / error de red).
type NoResultsCatalogItemResponseDto struct {
	CodRespuesta int    `json:"codRespuesta"`
	Mensaje      string `json:"Mensaje"`
	Excepcion    string `json:"Excepcion"`
}

// InsuranceCatalogItemResponseDto es el payload JSON para catálogos de seguros.
type InsuranceCatalogItemResponseDto struct {
	IdentificadorSeguro       string      `json:"IdentificadorSeguro"`
	Descripcion               string      `json:"Descripcion"`
	Poliza                    string      `json:"Poliza"`
	NombreCompania            string      `json:"NombreCompania"`
	CodigoCompania            int         `json:"CodigoCompania"`
	CodigoTipoSeguro          int         `json:"CodigoTipoSeguro"`
	CorrelativoPoliza         int         `json:"CorrelativoPoliza"`
	Tasa                      interface{} `json:"Tasa"`
	Factor                    interface{} `json:"Factor"`
	IndicadorPolizaIndividual int         `json:"IndicadorPolizaIndividual"`
	PorValorCuota             int         `json:"PorValorCuota"`
	IndicadorPolizaExterna    int         `json:"IndicadorPolizaExterna"`
}

// ErrorResponseDto asegura el orden exacto de claves JSON como el legado Java.
type ErrorResponseDto struct {
	Detail       *string       `json:"detail"`
	Code         *string       `json:"code"`
	Messages     []interface{} `json:"messages"`
	ErrorsDetail *string       `json:"errors_detail"`
}
