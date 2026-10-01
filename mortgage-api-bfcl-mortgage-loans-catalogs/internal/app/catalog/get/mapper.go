package get

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

func NewRequestMapper() RequestMapper {
	return func(req RequestDto) (domain.GetCatalogRequest, error) {
		return domain.GetCatalogRequest{
			CatalogName:   req.Uri.CatalogName,
			Channel:       req.Header.Channel,
			Commerce:      req.Header.Commerce,
			TransactionID: req.Header.TransactionID,
			Backend:       domain.CatalogBackend(strings.ToLower(req.Header.Backend)),
		}, nil
	}
}

func NewResponseMapper() ResponseMapper {
	return func(result interface{}) (interface{}, error) {
		switch data := result.(type) {
		case []domain.CatalogItem:
			dtos := make([]catalogItemDTO, len(data))
			for i, item := range data {
				dtos[i] = catalogItemDTO{CodigoAdm: item.Code, Descripcion: item.Description}
			}
			return dtos, nil
		case []domain.ComunaCatalogItem:
			dtos := make([]comunaCatalogItemDTO, len(data))
			for i, item := range data {
				dtos[i] = comunaCatalogItemDTO{CodigoAdm: item.Code, Descripcion: item.Description, RegionID: item.RegionID}
			}
			return dtos, nil
		case []domain.TipoDocumentoCatalogItem:
			dtos := make([]tipoDocumentoCatalogItemDTO, len(data))
			for i, item := range data {
				dtos[i] = tipoDocumentoCatalogItemDTO{CodigoAdm: item.Code, Descripcion: item.Description, GrupoID: item.GroupID}
			}
			return dtos, nil
		case []domain.NoResultsCatalogItem:
			dtos := make([]noResultsCatalogItemDTO, len(data))
			for i, item := range data {
				dtos[i] = noResultsCatalogItemDTO{CodRespuesta: item.CodRespuesta, Mensaje: item.Mensaje, Excepcion: item.Excepcion}
			}
			return dtos, nil
		case []domain.InsuranceCatalogItem:
			dtos := make([]insuranceCatalogItemDTO, len(data))
			for i, item := range data {
				dtos[i] = insuranceCatalogItemDTO{
					IdentificadorSeguro:       item.InsuranceIdentifier,
					Descripcion:               item.Description,
					Poliza:                    item.Policy,
					NombreCompania:            item.CompanyName,
					CodigoCompania:            item.CompanyCode,
					CodigoTipoSeguro:          item.InsuranceTypeCode,
					CorrelativoPoliza:         item.PolicyCorrelative,
					Tasa:                      json.Number(formatFloatLikeJackson(item.Rate)),
					Factor:                    json.Number(formatFloatLikeJackson(item.Factor)),
					IndicadorPolizaIndividual: item.IndividualPolicyFlag,
					PorValorCuota:             item.PerQuotaValueFlag,
					IndicadorPolizaExterna:    item.ExternalPolicyFlag,
				}
			}
			return dtos, nil
		default:
			return result, nil
		}
	}
}

func formatFloatLikeJackson(value float64) string {
	formatted := strconv.FormatFloat(value, 'g', -1, 64)
	if value != 0 && value > -1e-3 && value < 1e-3 {
		formatted = strconv.FormatFloat(value, 'E', -1, 64)
		formatted = strings.Replace(formatted, "E-0", "E-", 1)
		formatted = strings.Replace(formatted, "E+0", "E", 1)
	}
	return fmt.Sprint(formatted)
}
