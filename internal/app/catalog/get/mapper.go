package get

import (
	"encoding/json"
	"strconv"
	"strings"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

// ResponseMapper convierte el resultado del dominio al DTO de respuesta JSON.
type ResponseMapper = func(result interface{}) (interface{}, error)

// NewResponseMapper devuelve un mapper que transforma entidades de dominio a DTOs de respuesta.
func NewResponseMapper() ResponseMapper {
	return func(result interface{}) (interface{}, error) {
		switch data := result.(type) {
		case []domain.CatalogItem:
			dtos := make([]CatalogItemResponseDto, len(data))
			for i, item := range data {
				dtos[i] = CatalogItemResponseDto{
					CodigoAdm:   item.Code,
					Descripcion: item.Description,
				}
			}
			return dtos, nil

		case []domain.ComunaCatalogItem:
			dtos := make([]ComunaCatalogItemResponseDto, len(data))
			for i, item := range data {
				dtos[i] = ComunaCatalogItemResponseDto{
					CodigoAdm:   item.Code,
					Descripcion: item.Description,
					RegionID:    item.RegionID,
				}
			}
			return dtos, nil

		case []domain.TipoDocumentoCatalogItem:
			dtos := make([]TipoDocumentoCatalogItemResponseDto, len(data))
			for i, item := range data {
				dtos[i] = TipoDocumentoCatalogItemResponseDto{
					CodigoAdm:   item.Code,
					Descripcion: item.Description,
					GrupoID:     item.GroupID,
				}
			}
			return dtos, nil

		case []domain.NoResultsCatalogItem:
			dtos := make([]NoResultsCatalogItemResponseDto, len(data))
			for i, item := range data {
				dtos[i] = NoResultsCatalogItemResponseDto{
					CodRespuesta: item.CodRespuesta,
					Mensaje:      item.Mensaje,
					Excepcion:    item.Excepcion,
				}
			}
			return dtos, nil

		case []domain.InsuranceCatalogItem:
			dtos := make([]InsuranceCatalogItemResponseDto, len(data))
			for i, item := range data {
				dtos[i] = InsuranceCatalogItemResponseDto{
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

// formatFloatLikeJackson formatea un float64 replicando el comportamiento de Jackson (Java).
func formatFloatLikeJackson(f float64) string {
	s := strconv.FormatFloat(f, 'g', -1, 64)
	if f != 0 && f > -1e-3 && f < 1e-3 {
		s = strconv.FormatFloat(f, 'E', -1, 64)
		s = strings.Replace(s, "E-0", "E-", 1)
		s = strings.Replace(s, "E+0", "E", 1)
	}
	return s
}
