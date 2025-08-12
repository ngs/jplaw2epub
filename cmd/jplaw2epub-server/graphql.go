package main

import (
	"fmt"
	"time"
	
	"github.com/graphql-go/graphql"
	jplaw "go.ngs.io/jplaw-api-v2"
)

// GraphQL Enums
var categoryCodeEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "CategoryCode",
	Description: "Law category codes",
	Values: graphql.EnumValueConfigMap{
		"CONSTITUTION": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdConstitution,
			Description: "憲法",
		},
		"CRIMINAL": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdCriminal,
			Description: "刑事",
		},
		"FINANCE_GENERAL": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdFinanceGeneral,
			Description: "財務通則",
		},
		"FISHERIES": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdFisheries,
			Description: "水産業",
		},
		"TOURISM": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdTourism,
			Description: "観光",
		},
		"PARLIAMENT": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdParliament,
			Description: "国会",
		},
		"POLICE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdPolice,
			Description: "警察",
		},
		"NATIONAL_PROPERTY": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdNationalProperty,
			Description: "国有財産",
		},
		"MINING": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdMining,
			Description: "鉱業",
		},
		"POSTAL_SERVICE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdPostalService,
			Description: "郵務",
		},
		"ADMINISTRATIVE_ORG": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdAdministrativeOrg,
			Description: "行政組織",
		},
		"FIRE_SERVICE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdFireService,
			Description: "消防",
		},
		"NATIONAL_TAX": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdNationalTax,
			Description: "国税",
		},
		"INDUSTRY": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdIndustry,
			Description: "工業",
		},
		"TELECOMMUNICATIONS": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdTelecommunications,
			Description: "電気通信",
		},
		"CIVIL_SERVICE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdCivilService,
			Description: "公務員",
		},
		"NATIONAL_DEVELOPMENT": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdNationalDevelopment,
			Description: "国土開発",
		},
		"BUSINESS": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdBusiness,
			Description: "事業",
		},
		"COMMERCE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdCommerce,
			Description: "商事",
		},
		"LABOR": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdLabor,
			Description: "労働",
		},
		"ADMINISTRATIVE_PROC": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdAdministrativeProc,
			Description: "行政手続",
		},
		"LAND": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdLand,
			Description: "土地",
		},
		"NATIONAL_BONDS": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdNationalBonds,
			Description: "国債",
		},
		"FINANCE_INSURANCE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdFinanceInsurance,
			Description: "金融・保険",
		},
		"ENVIRONMENTAL_PROTECT": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdEnvironmentalProtect,
			Description: "環境保全",
		},
		"STATISTICS": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdStatistics,
			Description: "統計",
		},
		"CITY_PLANNING": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdCityPlanning,
			Description: "都市計画",
		},
		"EDUCATION": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdEducation,
			Description: "教育",
		},
		"FOREIGN_EXCHANGE_TRADE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdForeignExchangeTrade,
			Description: "外国為替・外国貿易",
		},
		"PUBLIC_HEALTH": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdPublicHealth,
			Description: "保健・衛生",
		},
		"LOCAL_GOVERNMENT": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdLocalGovernment,
			Description: "地方自治",
		},
		"ROADS": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdRoads,
			Description: "道路",
		},
		"CULTURE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdCulture,
			Description: "文教",
		},
		"LAND_TRANSPORT": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdLandTransport,
			Description: "陸運",
		},
		"SOCIAL_WELFARE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdSocialWelfare,
			Description: "社会福祉",
		},
		"LOCAL_FINANCE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdLocalFinance,
			Description: "地方財政",
		},
		"RIVERS": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdRivers,
			Description: "河川",
		},
		"INDUSTRY_GENERAL": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdIndustryGeneral,
			Description: "産業通則",
		},
		"MARITIME_TRANSPORT": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdMaritimeTransport,
			Description: "海運",
		},
		"SOCIAL_INSURANCE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdSocialInsurance,
			Description: "社会保険",
		},
		"JUDICIARY": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdJudiciary,
			Description: "司法",
		},
		"DISASTER_MANAGEMENT": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdDisasterManagement,
			Description: "災害対策",
		},
		"AGRICULTURE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdAgriculture,
			Description: "農業",
		},
		"AVIATION": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdAviation,
			Description: "航空",
		},
		"DEFENSE": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdDefense,
			Description: "防衛",
		},
		"CIVIL": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdCivil,
			Description: "民事",
		},
		"BUILDING_HOUSING": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdBuildingHousing,
			Description: "建築・住宅",
		},
		"FORESTRY": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdForestry,
			Description: "林業",
		},
		"FREIGHT_TRANSPORT": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdFreightTransport,
			Description: "貨物運送",
		},
		"FOREIGN_AFFAIRS": &graphql.EnumValueConfig{
			Value: jplaw.CategoryCdForeignAffairs,
			Description: "外事",
		},
	},
})

var lawTypeEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "LawType",
	Description: "Types of laws",
	Values: graphql.EnumValueConfigMap{
		"CONSTITUTION": &graphql.EnumValueConfig{
			Value: jplaw.LawTypeConstitution,
			Description: "憲法",
		},
		"ACT": &graphql.EnumValueConfig{
			Value: jplaw.LawTypeAct,
			Description: "法律",
		},
		"CABINET_ORDER": &graphql.EnumValueConfig{
			Value: jplaw.LawTypeCabinetorder,
			Description: "政令",
		},
		"IMPERIAL_ORDER": &graphql.EnumValueConfig{
			Value: jplaw.LawTypeImperialorder,
			Description: "勅令",
		},
		"MINISTERIAL_ORDINANCE": &graphql.EnumValueConfig{
			Value: jplaw.LawTypeMinisterialordinance,
			Description: "府省令",
		},
		"RULE": &graphql.EnumValueConfig{
			Value: jplaw.LawTypeRule,
			Description: "規則",
		},
		"MISC": &graphql.EnumValueConfig{
			Value: jplaw.LawTypeMisc,
			Description: "その他",
		},
	},
})

var lawNumEraEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "LawNumEra",
	Description: "Era names for law numbers",
	Values: graphql.EnumValueConfigMap{
		"MEIJI": &graphql.EnumValueConfig{
			Value: jplaw.LawNumEraMeiji,
			Description: "明治",
		},
		"TAISHO": &graphql.EnumValueConfig{
			Value: jplaw.LawNumEraTaisho,
			Description: "大正",
		},
		"SHOWA": &graphql.EnumValueConfig{
			Value: jplaw.LawNumEraShowa,
			Description: "昭和",
		},
		"HEISEI": &graphql.EnumValueConfig{
			Value: jplaw.LawNumEraHeisei,
			Description: "平成",
		},
		"REIWA": &graphql.EnumValueConfig{
			Value: jplaw.LawNumEraReiwa,
			Description: "令和",
		},
	},
})

var lawNumTypeEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "LawNumType",
	Description: "Law number types",
	Values: graphql.EnumValueConfigMap{
		"CONSTITUTION": &graphql.EnumValueConfig{
			Value: jplaw.LawNumTypeConstitution,
			Description: "憲法",
		},
		"ACT": &graphql.EnumValueConfig{
			Value: jplaw.LawNumTypeAct,
			Description: "法律",
		},
		"CABINET_ORDER": &graphql.EnumValueConfig{
			Value: jplaw.LawNumTypeCabinetorder,
			Description: "政令",
		},
		"IMPERIAL_ORDER": &graphql.EnumValueConfig{
			Value: jplaw.LawNumTypeImperialorder,
			Description: "勅令",
		},
		"MINISTERIAL_ORDINANCE": &graphql.EnumValueConfig{
			Value: jplaw.LawNumTypeMinisterialordinance,
			Description: "府省令",
		},
		"RULE": &graphql.EnumValueConfig{
			Value: jplaw.LawNumTypeRule,
			Description: "規則",
		},
		"MISC": &graphql.EnumValueConfig{
			Value: jplaw.LawNumTypeMisc,
			Description: "その他",
		},
	},
})

var currentRevisionStatusEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "CurrentRevisionStatus",
	Description: "History status",
	Values: graphql.EnumValueConfigMap{
		"CURRENT_ENFORCED": &graphql.EnumValueConfig{
			Value: jplaw.CurrentRevisionStatusCurrentenforced,
			Description: "現施行法令",
		},
		"UNENFORCED": &graphql.EnumValueConfig{
			Value: jplaw.CurrentRevisionStatusUnenforced,
			Description: "未施行法令",
		},
		"PREVIOUS_ENFORCED": &graphql.EnumValueConfig{
			Value: jplaw.CurrentRevisionStatusPreviousenforced,
			Description: "過去施行法令",
		},
		"REPEAL": &graphql.EnumValueConfig{
			Value: jplaw.CurrentRevisionStatusRepeal,
			Description: "廃止法令",
		},
	},
})

var repealStatusEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "RepealStatus",
	Description: "Repeal status",
	Values: graphql.EnumValueConfigMap{
		"NONE": &graphql.EnumValueConfig{
			Value: jplaw.RepealStatusNone,
			Description: "廃止・失効等のステータスなし",
		},
		"REPEAL": &graphql.EnumValueConfig{
			Value: jplaw.RepealStatusRepeal,
			Description: "廃止",
		},
		"EXPIRE": &graphql.EnumValueConfig{
			Value: jplaw.RepealStatusExpire,
			Description: "失効",
		},
		"SUSPEND": &graphql.EnumValueConfig{
			Value: jplaw.RepealStatusSuspend,
			Description: "停止",
		},
		"LOSS_OF_EFFECTIVENESS": &graphql.EnumValueConfig{
			Value: jplaw.RepealStatusLossofeffectiveness,
			Description: "実効性喪失",
		},
	},
})

var missionEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "Mission",
	Description: "Law mission type",
	Values: graphql.EnumValueConfigMap{
		"NEW": &graphql.EnumValueConfig{
			Value: jplaw.MissionNew,
			Description: "新規制定",
		},
		"PARTIAL": &graphql.EnumValueConfig{
			Value: jplaw.MissionPartial,
			Description: "一部改正",
		},
	},
})

// GraphQL Types - Fixed to match jplaw-api-v2 types

var lawInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LawInfo",
	Fields: graphql.Fields{
		"lawId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if law, ok := p.Source.(*jplaw.LawInfo); ok {
					return law.LawId, nil
				}
				return nil, nil
			},
		},
		"lawNum": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if law, ok := p.Source.(*jplaw.LawInfo); ok {
					return law.LawNum, nil
				}
				return nil, nil
			},
		},
		"lawNumEra": &graphql.Field{
			Type: lawNumEraEnum,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if law, ok := p.Source.(*jplaw.LawInfo); ok && law.LawNumEra != nil {
					return *law.LawNumEra, nil
				}
				return nil, nil
			},
		},
		"lawNumYear": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if law, ok := p.Source.(*jplaw.LawInfo); ok {
					return law.LawNumYear, nil
				}
				return nil, nil
			},
		},
		"lawNumNum": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if law, ok := p.Source.(*jplaw.LawInfo); ok {
					return law.LawNumNum, nil
				}
				return nil, nil
			},
		},
		"lawNumType": &graphql.Field{
			Type: lawNumTypeEnum,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if law, ok := p.Source.(*jplaw.LawInfo); ok && law.LawNumType != nil {
					return *law.LawNumType, nil
				}
				return nil, nil
			},
		},
		"lawType": &graphql.Field{
			Type: lawTypeEnum,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if law, ok := p.Source.(*jplaw.LawInfo); ok && law.LawType != nil {
					return *law.LawType, nil
				}
				return nil, nil
			},
		},
		"promulgationDate": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if law, ok := p.Source.(*jplaw.LawInfo); ok {
					return law.PromulgationDate.String(), nil
				}
				return nil, nil
			},
		},
	},
})

var revisionInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RevisionInfo",
	Fields: graphql.Fields{
		"lawRevisionId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.LawRevisionId, nil
				}
				return nil, nil
			},
		},
		"lawTitle": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.LawTitle, nil
				}
				return nil, nil
			},
		},
		"lawTitleKana": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.LawTitleKana, nil
				}
				return nil, nil
			},
		},
		"abbrev": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.Abbrev, nil
				}
				return nil, nil
			},
		},
		"lawType": &graphql.Field{
			Type: lawTypeEnum,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok && rev.LawType != nil {
					return *rev.LawType, nil
				}
				return nil, nil
			},
		},
		"amendmentLawId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.AmendmentLawId, nil
				}
				return nil, nil
			},
		},
		"amendmentLawTitle": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.AmendmentLawTitle, nil
				}
				return nil, nil
			},
		},
		"amendmentLawNum": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.AmendmentLawNum, nil
				}
				return nil, nil
			},
		},
		"amendmentPromulgateDate": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.AmendmentPromulgateDate.String(), nil
				}
				return nil, nil
			},
		},
		"amendmentEnforcementDate": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.AmendmentEnforcementDate.String(), nil
				}
				return nil, nil
			},
		},
		"repealDate": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.RepealDate.String(), nil
				}
				return nil, nil
			},
		},
		"remainInForce": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.RemainInForce, nil
				}
				return nil, nil
			},
		},
		"updated": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok {
					return rev.Updated.String(), nil
				}
				return nil, nil
			},
		},
		"currentRevisionStatus": &graphql.Field{
			Type: currentRevisionStatusEnum,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok && rev.CurrentRevisionStatus != nil {
					return *rev.CurrentRevisionStatus, nil
				}
				return nil, nil
			},
		},
		"repealStatus": &graphql.Field{
			Type: repealStatusEnum,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok && rev.RepealStatus != nil {
					return *rev.RepealStatus, nil
				}
				return nil, nil
			},
		},
		"mission": &graphql.Field{
			Type: missionEnum,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok && rev.Mission != nil {
					return *rev.Mission, nil
				}
				return nil, nil
			},
		},
	},
})

var lawItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LawItem",
	Fields: graphql.Fields{
		"lawInfo": &graphql.Field{
			Type: lawInfoType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if item, ok := p.Source.(jplaw.LawItem); ok {
					return item.LawInfo, nil
				}
				return nil, nil
			},
		},
		"revisionInfo": &graphql.Field{
			Type: revisionInfoType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if item, ok := p.Source.(jplaw.LawItem); ok {
					return item.RevisionInfo, nil
				}
				return nil, nil
			},
		},
		"currentRevisionInfo": &graphql.Field{
			Type: revisionInfoType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if item, ok := p.Source.(jplaw.LawItem); ok {
					return item.CurrentRevisionInfo, nil
				}
				return nil, nil
			},
		},
	},
})

var lawsResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LawsResponse",
	Fields: graphql.Fields{
		"count": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if resp, ok := p.Source.(*jplaw.LawsResponse); ok {
					return resp.Count, nil
				}
				return nil, nil
			},
		},
		"totalCount": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if resp, ok := p.Source.(*jplaw.LawsResponse); ok {
					return resp.TotalCount, nil
				}
				return nil, nil
			},
		},
		"nextOffset": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if resp, ok := p.Source.(*jplaw.LawsResponse); ok {
					return resp.NextOffset, nil
				}
				return nil, nil
			},
		},
		"laws": &graphql.Field{
			Type: graphql.NewList(lawItemType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if resp, ok := p.Source.(*jplaw.LawsResponse); ok {
					return resp.Laws, nil
				}
				return nil, nil
			},
		},
	},
})

var revisionsResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RevisionsResponse",
	Fields: graphql.Fields{
		"lawInfo": &graphql.Field{
			Type: lawInfoType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if resp, ok := p.Source.(*jplaw.LawRevisionsResponse); ok {
					return &resp.LawInfo, nil
				}
				return nil, nil
			},
		},
		"revisions": &graphql.Field{
			Type: graphql.NewList(revisionInfoType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if resp, ok := p.Source.(*jplaw.LawRevisionsResponse); ok {
					// Convert to interface slice for GraphQL
					result := make([]interface{}, len(resp.Revisions))
					for i, rev := range resp.Revisions {
						revCopy := rev
						result[i] = &revCopy
					}
					return result, nil
				}
				return nil, nil
			},
		},
	},
})

var keywordSentenceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "KeywordSentence",
	Fields: graphql.Fields{
		"text": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if sentence, ok := p.Source.(jplaw.KeywordSentence); ok {
					return sentence.Text, nil
				}
				return nil, nil
			},
		},
		"position": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if sentence, ok := p.Source.(jplaw.KeywordSentence); ok {
					return sentence.Position, nil
				}
				return nil, nil
			},
		},
	},
})

var keywordItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "KeywordItem",
	Fields: graphql.Fields{
		"lawInfo": &graphql.Field{
			Type: lawInfoType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if item, ok := p.Source.(jplaw.KeywordItem); ok {
					return item.LawInfo, nil
				}
				return nil, nil
			},
		},
		"revisionInfo": &graphql.Field{
			Type: revisionInfoType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if item, ok := p.Source.(jplaw.KeywordItem); ok {
					return item.RevisionInfo, nil
				}
				return nil, nil
			},
		},
		"sentences": &graphql.Field{
			Type: graphql.NewList(keywordSentenceType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if item, ok := p.Source.(jplaw.KeywordItem); ok {
					return item.Sentences, nil
				}
				return nil, nil
			},
		},
	},
})

var keywordResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "KeywordResponse",
	Fields: graphql.Fields{
		"totalCount": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if resp, ok := p.Source.(*jplaw.KeywordResponse); ok {
					return resp.TotalCount, nil
				}
				return nil, nil
			},
		},
		"sentenceCount": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if resp, ok := p.Source.(*jplaw.KeywordResponse); ok {
					return resp.SentenceCount, nil
				}
				return nil, nil
			},
		},
		"nextOffset": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if resp, ok := p.Source.(*jplaw.KeywordResponse); ok {
					return resp.NextOffset, nil
				}
				return nil, nil
			},
		},
		"items": &graphql.Field{
			Type: graphql.NewList(keywordItemType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if resp, ok := p.Source.(*jplaw.KeywordResponse); ok {
					return resp.Items, nil
				}
				return nil, nil
			},
		},
	},
})

// Query Type
var queryTypeFixed = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"laws": &graphql.Field{
			Type: lawsResponseType,
			Args: graphql.FieldConfigArgument{
				"lawId": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"lawNum": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"lawTitle": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"lawTitleKana": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"lawType": &graphql.ArgumentConfig{
					Type: graphql.NewList(lawTypeEnum),
				},
				"asof": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"categoryCode": &graphql.ArgumentConfig{
					Type: graphql.NewList(categoryCodeEnum),
				},
				"promulgateDateFrom": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"promulgateDateTo": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"limit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 100,
				},
				"offset": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				client := jplaw.NewClient()
				params := &jplaw.GetLawsParams{}

				// Map GraphQL arguments to API parameters
				if lawId, ok := p.Args["lawId"].(string); ok {
					params.LawId = &lawId
				}
				if lawNum, ok := p.Args["lawNum"].(string); ok {
					params.LawNum = &lawNum
				}
				if lawTitle, ok := p.Args["lawTitle"].(string); ok {
					params.LawTitle = &lawTitle
				}
				if lawTitleKana, ok := p.Args["lawTitleKana"].(string); ok {
					params.LawTitleKana = &lawTitleKana
				}
				if lawTypes, ok := p.Args["lawType"].([]interface{}); ok && len(lawTypes) > 0 {
					var types []jplaw.LawType
					for _, t := range lawTypes {
						if str, ok := t.(string); ok {
							types = append(types, jplaw.LawType(str))
						}
					}
					params.LawType = &types
				}
				if asof, ok := p.Args["asof"].(string); ok {
					if t, err := time.Parse("2006-01-02", asof); err == nil {
						date := jplaw.Date(t)
						params.Asof = &date
					}
				}
				if categoryCodes, ok := p.Args["categoryCode"].([]interface{}); ok && len(categoryCodes) > 0 {
					var codes []jplaw.CategoryCd
					for _, c := range categoryCodes {
						if str, ok := c.(string); ok {
							codes = append(codes, jplaw.CategoryCd(str))
						}
					}
					params.CategoryCd = &codes
				}
				if dateFrom, ok := p.Args["promulgateDateFrom"].(string); ok {
					if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
						date := jplaw.Date(t)
						params.PromulgationDateFrom = &date
					}
				}
				if dateTo, ok := p.Args["promulgateDateTo"].(string); ok {
					if t, err := time.Parse("2006-01-02", dateTo); err == nil {
						date := jplaw.Date(t)
						params.PromulgationDateTo = &date
					}
				}
				if limit, ok := p.Args["limit"].(int); ok {
					limit32 := int32(limit)
					params.Limit = &limit32
				}
				if offset, ok := p.Args["offset"].(int); ok {
					offset32 := int32(offset)
					params.Offset = &offset32
				}

				response, err := client.GetLaws(params)
				if err != nil {
					return nil, err
				}

				return response, nil
			},
		},
		"revisions": &graphql.Field{
			Type: revisionsResponseType,
			Args: graphql.FieldConfigArgument{
				"lawId": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Law ID or Law Number or Law Revision ID",
				},
				"lawTitle": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"lawTitleKana": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"amendmentLawId": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"amendmentDateFrom": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"amendmentDateTo": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"categoryCode": &graphql.ArgumentConfig{
					Type: graphql.NewList(categoryCodeEnum),
				},
				"updatedFrom": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"updatedTo": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				client := jplaw.NewClient()
				
				// lawId is required for GetRevisions
				lawId, ok := p.Args["lawId"].(string)
				if !ok || lawId == "" {
					return nil, fmt.Errorf("lawId is required")
				}

				params := &jplaw.GetRevisionsParams{}

				// Map optional parameters
				if lawTitle, ok := p.Args["lawTitle"].(string); ok {
					params.LawTitle = &lawTitle
				}
				if lawTitleKana, ok := p.Args["lawTitleKana"].(string); ok {
					params.LawTitleKana = &lawTitleKana
				}
				if amendmentLawId, ok := p.Args["amendmentLawId"].(string); ok {
					params.AmendmentLawId = &amendmentLawId
				}
				if dateFrom, ok := p.Args["amendmentDateFrom"].(string); ok {
					if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
						date := jplaw.Date(t)
						params.AmendmentDateFrom = &date
					}
				}
				if dateTo, ok := p.Args["amendmentDateTo"].(string); ok {
					if t, err := time.Parse("2006-01-02", dateTo); err == nil {
						date := jplaw.Date(t)
						params.AmendmentDateTo = &date
					}
				}
				if categoryCodes, ok := p.Args["categoryCode"].([]interface{}); ok && len(categoryCodes) > 0 {
					var codes []jplaw.CategoryCd
					for _, c := range categoryCodes {
						if str, ok := c.(string); ok {
							codes = append(codes, jplaw.CategoryCd(str))
						}
					}
					params.CategoryCd = &codes
				}
				if updatedFrom, ok := p.Args["updatedFrom"].(string); ok {
					if t, err := time.Parse("2006-01-02", updatedFrom); err == nil {
						date := jplaw.Date(t)
						params.UpdatedFrom = &date
					}
				}
				if updatedTo, ok := p.Args["updatedTo"].(string); ok {
					if t, err := time.Parse("2006-01-02", updatedTo); err == nil {
						date := jplaw.Date(t)
						params.UpdatedTo = &date
					}
				}

				response, err := client.GetRevisions(lawId, params)
				if err != nil {
					return nil, err
				}

				return response, nil
			},
		},
		"keyword": &graphql.Field{
			Type: keywordResponseType,
			Args: graphql.FieldConfigArgument{
				"keyword": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"lawNum": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"lawType": &graphql.ArgumentConfig{
					Type: graphql.NewList(lawTypeEnum),
				},
				"asof": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"categoryCode": &graphql.ArgumentConfig{
					Type: graphql.NewList(categoryCodeEnum),
				},
				"promulgateDateFrom": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"promulgateDateTo": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"limit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 100,
				},
				"offset": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"sentencesLimit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 10,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				client := jplaw.NewClient()
				
				keyword, ok := p.Args["keyword"].(string)
				if !ok || keyword == "" {
					return nil, fmt.Errorf("keyword is required")
				}

				params := &jplaw.GetKeywordParams{
					Keyword: keyword,
				}

				// Map optional parameters
				if lawNum, ok := p.Args["lawNum"].(string); ok {
					params.LawNum = &lawNum
				}
				if lawTypes, ok := p.Args["lawType"].([]interface{}); ok && len(lawTypes) > 0 {
					var types []jplaw.LawType
					for _, t := range lawTypes {
						if str, ok := t.(string); ok {
							types = append(types, jplaw.LawType(str))
						}
					}
					params.LawType = &types
				}
				if asof, ok := p.Args["asof"].(string); ok {
					if t, err := time.Parse("2006-01-02", asof); err == nil {
						date := jplaw.Date(t)
						params.Asof = &date
					}
				}
				if categoryCodes, ok := p.Args["categoryCode"].([]interface{}); ok && len(categoryCodes) > 0 {
					var codes []jplaw.CategoryCd
					for _, c := range categoryCodes {
						if str, ok := c.(string); ok {
							codes = append(codes, jplaw.CategoryCd(str))
						}
					}
					params.CategoryCd = &codes
				}
				if dateFrom, ok := p.Args["promulgateDateFrom"].(string); ok {
					if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
						date := jplaw.Date(t)
						params.PromulgationDateFrom = &date
					}
				}
				if dateTo, ok := p.Args["promulgateDateTo"].(string); ok {
					if t, err := time.Parse("2006-01-02", dateTo); err == nil {
						date := jplaw.Date(t)
						params.PromulgationDateTo = &date
					}
				}
				if limit, ok := p.Args["limit"].(int); ok {
					limit32 := int32(limit)
					params.Limit = &limit32
				}
				if offset, ok := p.Args["offset"].(int); ok {
					offset32 := int32(offset)
					params.Offset = &offset32
				}
				if sentencesLimit, ok := p.Args["sentencesLimit"].(int); ok {
					limit32 := int32(sentencesLimit)
					params.SentencesLimit = &limit32
				}

				response, err := client.GetKeyword(params)
				if err != nil {
					return nil, err
				}

				return response, nil
			},
		},
	},
})

// SchemaFixed - New schema with fixed types
var SchemaFixed, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryTypeFixed,
})