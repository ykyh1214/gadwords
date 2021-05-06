package gadwords

import (
	"encoding/json"
	"log"
	"time"
)

// A campaignService holds the connection information for the
// campaign service.
func CreateAppCampaign(budgetId, targetCpa int64, name, appId string) {
	authConf, err := NewCredentialsFromFile("./config.json")
	cs := NewCampaignService(&authConf.Auth)
	campaigns, err := cs.Mutate(
		CampaignOperations{
			"ADD": {
				Campaign{
					Name:      name,
					Status:    "PAUSED",
					StartDate: time.Now().Format("20060102"),
					BudgetId:  budgetId,
					UniversalAppCampaignInfo: &CampaignAppSetting{
						AppVendor:                           "VENDOR_GOOGLE_MARKET",
						AppId:                               appId,
						UniversalAppBiddingStrategyGoalType: "OPTIMIZE_FOR_INSTALL_CONVERSION_VOLUME",
					},
					Settings: []CampaignSetting{
						NewGeoTargetTypeSetting("DONT_CARE", "LOCATION_OF_PRESENCE"),
					},
					AdvertisingChannelType:    "MULTI_CHANNEL",
					AdvertisingChannelSubType: "UNIVERSAL_APP_CAMPAIGN",
					BiddingStrategyConfiguration: &BiddingStrategyConfiguration{
						StrategyType: "TARGET_CPA",
						Scheme:       &BiddingScheme{Type: "TargetCpaBiddingScheme", TargetCpa: &TargetCpa{Amount: targetCpa}},
					},
				},
			},
		},
	)
	js, _ := json.MarshalIndent(campaigns, "", " ")
	log.Println(string(js))
	log.Println(err)
	return
}
