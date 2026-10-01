package catalog_service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	httpFIF "github.com/falabella-regulado/go-lib-http-fif"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/out"
)

// finnflowRepository implementa CatalogRepository usando go-lib-http-fif para GET a Finnflow.
type finnflowRepository struct {
	client      httpFIF.RestClient
	routeSuffix string // e.g. "/api/catalogo_detail/"
	basicAuth   string // "Basic <base64(key:secret)>" — precomputado
}

// NewFinnflowRepository crea un repositorio para Finnflow (GET + Basic Auth).
// La autenticación Basic se inyecta como header Authorization en cada request.
func NewFinnflowRepository(client httpFIF.RestClient, routeSuffix, key, secret string) out.CatalogRepository {
	encoded := base64.StdEncoding.EncodeToString([]byte(key + ":" + secret))
	return &finnflowRepository{
		client:      client,
		routeSuffix: routeSuffix,
		basicAuth:   "Basic " + encoded,
	}
}

func (r *finnflowRepository) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	path := r.routeSuffix + req.CatalogName

	headers := http.Header{}
	if r.basicAuth != "Basic " {
		headers.Set("Authorization", r.basicAuth)
	}

	res, err := r.client.Get(ctx, path, headers)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "timeout") || strings.Contains(msg, "connection") || strings.Contains(msg, "refused") {
			return noResults(), nil
		}
		return nil, &domain.ServiceError{Message: "request to Finnflow failed", Cause: err}
	}
	defer res.Body.Close()

	return handleCatalogResponse(res, req.CatalogName)
}

// javaProxyRepository implementa CatalogRepository usando go-lib-http-fif para POST al proxy Java.
type javaProxyRepository struct {
	client      httpFIF.RestClient
	routeSuffix string // e.g. "/v1/bfcl/mortgage-loan/catalogs/"
}

// NewJavaProxyRepository crea un repositorio para el proxy Java legado (POST + X-Headers).
func NewJavaProxyRepository(client httpFIF.RestClient, routeSuffix string) out.CatalogRepository {
	return &javaProxyRepository{
		client:      client,
		routeSuffix: routeSuffix,
	}
}

func (r *javaProxyRepository) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	path := r.routeSuffix + req.CatalogName

	headers := http.Header{}
	if req.Channel != "" {
		headers.Set("X-Channel", req.Channel)
	}
	if req.Commerce != "" {
		headers.Set("X-Commerce", req.Commerce)
	}
	if req.TransactionID != "" {
		headers.Set("X-Transaction-ID", req.TransactionID)
	}

	res, err := r.client.Post(ctx, path, nil, headers)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "timeout") || strings.Contains(msg, "connection") || strings.Contains(msg, "refused") {
			return noResults(), nil
		}
		return nil, &domain.ServiceError{Message: "request to Java proxy failed", Cause: err}
	}
	defer res.Body.Close()

	return handleCatalogResponse(res, req.CatalogName)
}

// handleCatalogResponse interpreta el status HTTP y decodifica el body.
func handleCatalogResponse(res *http.Response, catalogName string) (interface{}, error) {
	if res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusUnauthorized:
			return nil, &domain.AuthenticationError{StatusCode: res.StatusCode}
		case http.StatusPaymentRequired, http.StatusNotFound:
			return nil, domain.ErrCatalogNotFound
		default:
			return nil, &domain.ServiceError{Message: fmt.Sprintf("catalog upstream returned status %d", res.StatusCode)}
		}
	}
	return decodeCatalog(catalogName, res)
}

// decodeCatalog parsea el body JSON según el tipo de catálogo.
func decodeCatalog(catalogName string, res *http.Response) (interface{}, error) {
	catalogLower := strings.ToLower(catalogName)

	// ——— Seguros ———
	if strings.HasPrefix(catalogLower, "seguros") {
		var items []FinnflowInsuranceItem
		if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
			if strings.Contains(err.Error(), "EOF") {
				return noResults(), nil
			}
			return nil, &domain.ServiceError{Message: "can't decode insurance catalog response", Cause: err}
		}
		out := make([]domain.InsuranceCatalogItem, len(items))
		for i, item := range items {
			out[i] = domain.InsuranceCatalogItem{
				InsuranceIdentifier:  item.IdentificadorSeguro,
				Description:          item.Descripcion,
				Policy:               item.Poliza,
				CompanyName:          item.NombreCompania,
				CompanyCode:          item.CodigoCompania,
				InsuranceTypeCode:    item.CodigoTipoSeguro,
				PolicyCorrelative:    item.CorrelativoPoliza,
				Rate:                 item.Tasa,
				Factor:               item.Factor,
				IndividualPolicyFlag: item.IndicadorPolizaIndividual,
				PerQuotaValueFlag:    item.PorValorCuota,
				ExternalPolicyFlag:   item.IndicadorPolizaExterna,
			}
		}
		return out, nil
	}

	// ——— TiposDocumentos ———
	if catalogName == "TiposDocumentos" {
		var items []FinnflowTipoDocumentoItem
		if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
			if strings.Contains(err.Error(), "EOF") {
				return noResults(), nil
			}
			return nil, &domain.ServiceError{Message: "can't decode TiposDocumentos catalog response", Cause: err}
		}
		out := make([]domain.TipoDocumentoCatalogItem, len(items))
		for i, item := range items {
			out[i] = domain.TipoDocumentoCatalogItem{
				Code:        item.CodigoAdm,
				Description: item.Descripcion,
				GroupID:     item.GrupoID,
			}
		}
		return out, nil
	}

	// ——— Comunas ———
	if catalogName == "Comunas" {
		var items []FinnflowComunaItem
		if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
			if strings.Contains(err.Error(), "EOF") {
				return noResults(), nil
			}
			return nil, &domain.ServiceError{Message: "can't decode Comunas catalog response", Cause: err}
		}
		out := make([]domain.ComunaCatalogItem, len(items))
		for i, item := range items {
			out[i] = domain.ComunaCatalogItem{
				Code:        fmt.Sprintf("%v", item.CodigoAdm),
				Description: item.Descripcion,
				RegionID:    item.RegionID,
			}
		}
		return out, nil
	}

	// ——— Catálogo genérico ———
	var items []FinnflowGenericItem
	if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
		if strings.Contains(err.Error(), "EOF") {
			return noResults(), nil
		}
		return nil, &domain.ServiceError{Message: "can't decode catalog response", Cause: err}
	}

	if len(items) == 0 || (len(items) > 0 && items[0].CodRespuesta != nil) {
		return noResults(), nil
	}

	result := make([]domain.CatalogItem, len(items))
	for i, item := range items {
		result[i] = domain.CatalogItem{
			Code:        fmt.Sprintf("%v", item.CodigoAdm),
			Description: item.Descripcion,
		}
	}
	return result, nil
}

func noResults() []domain.NoResultsCatalogItem {
	return []domain.NoResultsCatalogItem{{CodRespuesta: 3, Mensaje: "Sin resultados.", Excepcion: "Ninguna"}}
}
