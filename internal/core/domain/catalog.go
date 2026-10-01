package domain

// CatalogItem represents a generic catalog item.
type CatalogItem struct {
	Code        string
	Description string
}

// ComunaCatalogItem represents a catalog item for comunas (includes RegionID).
type ComunaCatalogItem struct {
	Code        string
	Description string
	RegionID    int
}

// TipoDocumentoCatalogItem represents a catalog item for tipos de documento.
type TipoDocumentoCatalogItem struct {
	Code        int
	Description string
	GroupID     int
}

// NoResultsCatalogItem represents the legacy empty/network-error fallback.
type NoResultsCatalogItem struct {
	CodRespuesta int
	Mensaje      string
	Excepcion    string
}

// InsuranceCatalogItem represents an insurance catalog item (seguros).
type InsuranceCatalogItem struct {
	InsuranceIdentifier  string
	Description          string
	Policy               string
	CompanyName          string
	CompanyCode          int
	InsuranceTypeCode    int
	PolicyCorrelative    int
	Rate                 float64
	Factor               float64
	IndividualPolicyFlag int
	PerQuotaValueFlag    int
	ExternalPolicyFlag   int
}

// GetCatalogRequest represents the pure domain request to get a catalog.
type GetCatalogRequest struct {
	CatalogName   string
	Channel       string
	Commerce      string
	TransactionID string
}
