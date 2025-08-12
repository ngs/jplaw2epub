package main

import (
	"fmt"
	"time"
	
	"github.com/graphql-go/graphql"
	jplaw "go.ngs.io/jplaw-api-v2"
)

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
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if law, ok := p.Source.(*jplaw.LawInfo); ok && law.LawNumEra != nil {
					return string(*law.LawNumEra), nil
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
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if law, ok := p.Source.(*jplaw.LawInfo); ok && law.LawNumType != nil {
					return string(*law.LawNumType), nil
				}
				return nil, nil
			},
		},
		"lawType": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if law, ok := p.Source.(*jplaw.LawInfo); ok && law.LawType != nil {
					return string(*law.LawType), nil
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
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok && rev.LawType != nil {
					return string(*rev.LawType), nil
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
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok && rev.CurrentRevisionStatus != nil {
					return string(*rev.CurrentRevisionStatus), nil
				}
				return nil, nil
			},
		},
		"repealStatus": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok && rev.RepealStatus != nil {
					return string(*rev.RepealStatus), nil
				}
				return nil, nil
			},
		},
		"mission": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if rev, ok := p.Source.(*jplaw.RevisionInfo); ok && rev.Mission != nil {
					return string(*rev.Mission), nil
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
					Type: graphql.NewList(graphql.String),
				},
				"asof": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"categoryCode": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
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
					Type: graphql.NewList(graphql.String),
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
					Type: graphql.NewList(graphql.String),
				},
				"asof": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"categoryCode": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
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