package domain

type CatalogItem struct {
	Code        string
	Description string
}

type ComunaCatalogItem struct {
	Code        string
	Description string
	RegionID    int
}

type TipoDocumentoCatalogItem struct {
	Code        int
	Description string
	GroupID     int
}

type NoResultsCatalogItem struct {
	CodRespuesta int
	Mensaje      string
	Excepcion    string
}

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

type CatalogBackend string

const (
	CatalogBackendDefault  CatalogBackend = ""
	CatalogBackendDummy    CatalogBackend = "dummy"
	CatalogBackendFinnflow CatalogBackend = "real"
	CatalogBackendJava     CatalogBackend = "java"
)

type GetCatalogRequest struct {
	CatalogName   string
	Channel       string
	Commerce      string
	TransactionID string
	Backend       CatalogBackend
}
