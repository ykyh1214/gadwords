package gadwords

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"time"
)

func rand_str2(str_size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

// A campaignService holds the connection information for the
// campaign service.
func CreateSearchCampaign(name string) {
	authConf, err := NewCredentialsFromFile("./config.json")
	bs := NewBudgetService(&authConf.Auth)
	budgets, err := bs.Mutate(
		BudgetOperations{
			"ADD": {
				Budget{
					Name: "testbudget " + rand_str2(10),
					//Period:   "DAILY",
					Amount:   60000000,
					Delivery: "STANDARD",
					Shared:   false,
				},
			},
		},
	)
	cs := NewCampaignService(&authConf.Auth)
	campaigns, err := cs.Mutate(
		CampaignOperations{
			"ADD": {
				Campaign{
					Name:      name,
					Status:    "PAUSED",
					StartDate: time.Now().Format("20060102"),
					EndDate:   time.Now().Format("20060102"),
					BudgetId:  budgets[0].Id,
					Settings: []CampaignSetting{
						NewGeoTargetTypeSetting("LOCATION_OF_PRESENCE", "LOCATION_OF_PRESENCE"),
						CampaignSetting{Type: "TargetingSetting", Details: []TargetDetails{{CriterionTypeGroup: "USER_INTEREST_AND_LIST", TargetAll: true}}},
					},
					NetworkSetting: &NetworkSetting{
						TargetGoogleSearch:         true,
						TargetSearchNetwork:        true, //只有这个能改
						TargetContentNetwork:       false,
						TargetPartnerSearchNetwork: false,
					},
					AdvertisingChannelType: "SEARCH",
					BiddingStrategyConfiguration: &BiddingStrategyConfiguration{
						StrategyType: "TARGET_SPEND",
						Scheme:       &BiddingScheme{Type: "TargetSpendBiddingScheme", BidCeiling: &TargetCpa{1200000}},
					},
				},
			},
		},
	)
	js, _ := json.MarshalIndent(campaigns, "", " ")
	log.Println(string(js))
	log.Println(err)
	AddSearchCampaignCriterion(&authConf.Auth, campaigns[0].Id)
	return
}
func AddSearchCampaignCriterion(auth *Auth, CampaignId int64) {
	ccs := NewCampaignCriterionService(auth)
	_, err := ccs.Mutate(
		CampaignCriterionOperations{
			"ADD": {
				//AdSchedule
				CampaignCriterion{CampaignId: CampaignId, Criterion: AdScheduleCriterion{DayOfWeek: "MONDAY", StartHour: "10", StartMinute: "ZERO", EndHour: "13", EndMinute: "ZERO"}},
				//Location
				CampaignCriterion{CampaignId: CampaignId, Criterion: Location{Id: 2764}},
				//language
				CampaignCriterion{CampaignId: CampaignId, Criterion: LanguageCriterion{Id: 1000}},
				CampaignCriterion{CampaignId: CampaignId, Criterion: LanguageCriterion{Id: 1017}},
			},
		},
	)
	_ = err
}

func AddSearchAdGroup(auth *Auth, CampaignId int64) {
	ags := NewAdGroupService(auth)
	adGroups, err := ags.Mutate(
		AdGroupOperations{
			"ADD": {
				AdGroup{
					Name:        "test ad group " + rand_str2(10),
					Status:      "ENABLED",
					CampaignId:  CampaignId,
					AdGroupType: "SEARCH_STANDARD",
				},
			},
		},
	)
	log.Println(adGroups, err)
}

func AddSearchKeywords(auth *Auth, adGroupId int64) {
	ags := NewAdGroupCriterionService(auth)
	adGroups, err := ags.Mutate(AdGroupCriterionOperations{
		"ADD": {
			BiddableAdGroupCriterion{
				AdGroupId: adGroupId,
				Criterion: KeywordCriterion{Text: "shoe2", MatchType: "PHRASE"},
			},
		},
	})

	log.Println(adGroups, err)
}

//
//112993761796
