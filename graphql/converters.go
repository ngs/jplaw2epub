package graphql

import (
	jplaw "go.ngs.io/jplaw-api-v2"
	"go.ngs.io/jplaw2epub/graphql/model"
)

// convertCategoryCode converts GraphQL CategoryCode to jplaw CategoryCd
func convertCategoryCode(codes []model.CategoryCode) []jplaw.CategoryCd {
	if len(codes) == 0 {
		return nil
	}

	result := make([]jplaw.CategoryCd, 0, len(codes))
	for _, code := range codes {
		switch code {
		case model.CategoryCodeConstitution:
			result = append(result, jplaw.CategoryCdConstitution)
		case model.CategoryCodeCriminal:
			result = append(result, jplaw.CategoryCdCriminal)
		case model.CategoryCodeFinanceGeneral:
			result = append(result, jplaw.CategoryCdFinanceGeneral)
		case model.CategoryCodeFisheries:
			result = append(result, jplaw.CategoryCdFisheries)
		case model.CategoryCodeTourism:
			result = append(result, jplaw.CategoryCdTourism)
		case model.CategoryCodeParliament:
			result = append(result, jplaw.CategoryCdParliament)
		case model.CategoryCodePolice:
			result = append(result, jplaw.CategoryCdPolice)
		case model.CategoryCodeNationalProperty:
			result = append(result, jplaw.CategoryCdNationalProperty)
		case model.CategoryCodeMining:
			result = append(result, jplaw.CategoryCdMining)
		case model.CategoryCodePostalService:
			result = append(result, jplaw.CategoryCdPostalService)
		case model.CategoryCodeAdministrativeOrg:
			result = append(result, jplaw.CategoryCdAdministrativeOrg)
		case model.CategoryCodeFireService:
			result = append(result, jplaw.CategoryCdFireService)
		case model.CategoryCodeNationalTax:
			result = append(result, jplaw.CategoryCdNationalTax)
		case model.CategoryCodeIndustry:
			result = append(result, jplaw.CategoryCdIndustry)
		case model.CategoryCodeTelecommunications:
			result = append(result, jplaw.CategoryCdTelecommunications)
		case model.CategoryCodeCivilService:
			result = append(result, jplaw.CategoryCdCivilService)
		case model.CategoryCodeNationalDevelopment:
			result = append(result, jplaw.CategoryCdNationalDevelopment)
		case model.CategoryCodeBusiness:
			result = append(result, jplaw.CategoryCdBusiness)
		case model.CategoryCodeCommerce:
			result = append(result, jplaw.CategoryCdCommerce)
		case model.CategoryCodeLabor:
			result = append(result, jplaw.CategoryCdLabor)
		case model.CategoryCodeAdministrativeProc:
			result = append(result, jplaw.CategoryCdAdministrativeProc)
		case model.CategoryCodeLand:
			result = append(result, jplaw.CategoryCdLand)
		case model.CategoryCodeNationalBonds:
			result = append(result, jplaw.CategoryCdNationalBonds)
		case model.CategoryCodeFinanceInsurance:
			result = append(result, jplaw.CategoryCdFinanceInsurance)
		case model.CategoryCodeEnvironmentalProtect:
			result = append(result, jplaw.CategoryCdEnvironmentalProtect)
		case model.CategoryCodeStatistics:
			result = append(result, jplaw.CategoryCdStatistics)
		case model.CategoryCodeCityPlanning:
			result = append(result, jplaw.CategoryCdCityPlanning)
		case model.CategoryCodeEducation:
			result = append(result, jplaw.CategoryCdEducation)
		case model.CategoryCodeForeignExchangeTrade:
			result = append(result, jplaw.CategoryCdForeignExchangeTrade)
		case model.CategoryCodePublicHealth:
			result = append(result, jplaw.CategoryCdPublicHealth)
		case model.CategoryCodeLocalGovernment:
			result = append(result, jplaw.CategoryCdLocalGovernment)
		case model.CategoryCodeRoads:
			result = append(result, jplaw.CategoryCdRoads)
		case model.CategoryCodeCulture:
			result = append(result, jplaw.CategoryCdCulture)
		case model.CategoryCodeLandTransport:
			result = append(result, jplaw.CategoryCdLandTransport)
		case model.CategoryCodeSocialWelfare:
			result = append(result, jplaw.CategoryCdSocialWelfare)
		case model.CategoryCodeLocalFinance:
			result = append(result, jplaw.CategoryCdLocalFinance)
		case model.CategoryCodeRivers:
			result = append(result, jplaw.CategoryCdRivers)
		case model.CategoryCodeIndustryGeneral:
			result = append(result, jplaw.CategoryCdIndustryGeneral)
		case model.CategoryCodeMaritimeTransport:
			result = append(result, jplaw.CategoryCdMaritimeTransport)
		case model.CategoryCodeSocialInsurance:
			result = append(result, jplaw.CategoryCdSocialInsurance)
		case model.CategoryCodeJudiciary:
			result = append(result, jplaw.CategoryCdJudiciary)
		case model.CategoryCodeDisasterManagement:
			result = append(result, jplaw.CategoryCdDisasterManagement)
		case model.CategoryCodeAgriculture:
			result = append(result, jplaw.CategoryCdAgriculture)
		case model.CategoryCodeAviation:
			result = append(result, jplaw.CategoryCdAviation)
		case model.CategoryCodeDefense:
			result = append(result, jplaw.CategoryCdDefense)
		case model.CategoryCodeCivil:
			result = append(result, jplaw.CategoryCdCivil)
		case model.CategoryCodeBuildingHousing:
			result = append(result, jplaw.CategoryCdBuildingHousing)
		case model.CategoryCodeForestry:
			result = append(result, jplaw.CategoryCdForestry)
		case model.CategoryCodeFreightTransport:
			result = append(result, jplaw.CategoryCdFreightTransport)
		case model.CategoryCodeForeignAffairs:
			result = append(result, jplaw.CategoryCdForeignAffairs)
		}
	}
	return result
}

// convertLawType converts GraphQL LawType to jplaw LawType
func convertLawType(types []model.LawType) []jplaw.LawType {
	if len(types) == 0 {
		return nil
	}

	result := make([]jplaw.LawType, 0, len(types))
	for _, t := range types {
		switch t {
		case model.LawTypeConstitution:
			result = append(result, jplaw.LawTypeConstitution)
		case model.LawTypeAct:
			result = append(result, jplaw.LawTypeAct)
		case model.LawTypeCabinetOrder:
			result = append(result, jplaw.LawTypeCabinetorder)
		case model.LawTypeImperialOrder:
			result = append(result, jplaw.LawTypeImperialorder)
		case model.LawTypeMinisterialOrdinance:
			result = append(result, jplaw.LawTypeMinisterialordinance)
		case model.LawTypeRule:
			result = append(result, jplaw.LawTypeRule)
		case model.LawTypeMisc:
			result = append(result, jplaw.LawTypeMisc)
		}
	}
	return result
}

// Reverse conversions for output

func convertLawTypeToModel(t *jplaw.LawType) *model.LawType {
	if t == nil {
		return nil
	}

	var result model.LawType
	switch *t {
	case jplaw.LawTypeConstitution:
		result = model.LawTypeConstitution
	case jplaw.LawTypeAct:
		result = model.LawTypeAct
	case jplaw.LawTypeCabinetorder:
		result = model.LawTypeCabinetOrder
	case jplaw.LawTypeImperialorder:
		result = model.LawTypeImperialOrder
	case jplaw.LawTypeMinisterialordinance:
		result = model.LawTypeMinisterialOrdinance
	case jplaw.LawTypeRule:
		result = model.LawTypeRule
	case jplaw.LawTypeMisc:
		result = model.LawTypeMisc
	default:
		return nil
	}
	return &result
}

func convertLawNumEraToModel(e *jplaw.LawNumEra) *model.LawNumEra {
	if e == nil {
		return nil
	}

	var result model.LawNumEra
	switch *e {
	case jplaw.LawNumEraMeiji:
		result = model.LawNumEraMeiji
	case jplaw.LawNumEraTaisho:
		result = model.LawNumEraTaisho
	case jplaw.LawNumEraShowa:
		result = model.LawNumEraShowa
	case jplaw.LawNumEraHeisei:
		result = model.LawNumEraHeisei
	case jplaw.LawNumEraReiwa:
		result = model.LawNumEraReiwa
	default:
		return nil
	}
	return &result
}

func convertLawNumTypeToModel(t *jplaw.LawNumType) *model.LawNumType {
	if t == nil {
		return nil
	}

	var result model.LawNumType
	switch *t {
	case jplaw.LawNumTypeConstitution:
		result = model.LawNumTypeConstitution
	case jplaw.LawNumTypeAct:
		result = model.LawNumTypeAct
	case jplaw.LawNumTypeCabinetorder:
		result = model.LawNumTypeCabinetOrder
	case jplaw.LawNumTypeImperialorder:
		result = model.LawNumTypeImperialOrder
	case jplaw.LawNumTypeMinisterialordinance:
		result = model.LawNumTypeMinisterialOrdinance
	case jplaw.LawNumTypeRule:
		result = model.LawNumTypeRule
	case jplaw.LawNumTypeMisc:
		result = model.LawNumTypeMisc
	default:
		return nil
	}
	return &result
}

func convertCurrentRevisionStatusToModel(s *jplaw.CurrentRevisionStatus) *model.CurrentRevisionStatus {
	if s == nil {
		return nil
	}

	var result model.CurrentRevisionStatus
	switch *s {
	case jplaw.CurrentRevisionStatusCurrentenforced:
		result = model.CurrentRevisionStatusCurrentEnforced
	case jplaw.CurrentRevisionStatusUnenforced:
		result = model.CurrentRevisionStatusUnenforced
	case jplaw.CurrentRevisionStatusPreviousenforced:
		result = model.CurrentRevisionStatusPreviousEnforced
	case jplaw.CurrentRevisionStatusRepeal:
		result = model.CurrentRevisionStatusRepeal
	default:
		return nil
	}
	return &result
}

func convertRepealStatusToModel(s *jplaw.RepealStatus) *model.RepealStatus {
	if s == nil {
		return nil
	}

	var result model.RepealStatus
	switch *s {
	case jplaw.RepealStatusNone:
		result = model.RepealStatusNone
	case jplaw.RepealStatusRepeal:
		result = model.RepealStatusRepeal
	case jplaw.RepealStatusExpire:
		result = model.RepealStatusExpire
	case jplaw.RepealStatusSuspend:
		result = model.RepealStatusSuspend
	case jplaw.RepealStatusLossofeffectiveness:
		result = model.RepealStatusLossOfEffectiveness
	default:
		return nil
	}
	return &result
}

func convertMissionToModel(m *jplaw.Mission) *model.Mission {
	if m == nil {
		return nil
	}

	var result model.Mission
	switch *m {
	case jplaw.MissionNew:
		result = model.MissionNew
	case jplaw.MissionPartial:
		result = model.MissionPartial
	default:
		return nil
	}
	return &result
}
