package gadwords

import (
	"encoding/xml"
	"fmt"
	"log"
	"testing"
	"time"
)

func testCampaignService(t *testing.T) (service *CampaignService) {
	return &CampaignService{Auth: testAuthSetup22(t)}
}

func testCampaign(t *testing.T) (Campaign, func()) {
	budget, cleanupBudget := testBudget(t)
	cs := testCampaignService(t)
	campaigns, err := cs.Mutate(
		CampaignOperations{
			"ADD": {
				Campaign{
					Name:      "test campaign " + rand_str(10),
					Status:    "PAUSED",
					StartDate: time.Now().Format("20060102"),
					BudgetId:  budget.Id,
					Settings: []CampaignSetting{
						NewGeoTargetTypeSetting("DONT_CARE", "LOCATION_OF_PRESENCE"),
					},
					AdvertisingChannelType: "SEARCH",
					BiddingStrategyConfiguration: &BiddingStrategyConfiguration{
						StrategyType: "MANUAL_CPC",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	cleanupCampaign := func() {
		campaigns[0].Status = "REMOVED"
		_, err = cs.Mutate(CampaignOperations{"SET": campaigns})
		if err != nil {
			t.Error(err)
		}
		cleanupBudget()
	}
	return campaigns[0], cleanupCampaign
}

func TestCampaign(t *testing.T) {
	budget, cleanupBudget := testBudget(t)
	defer cleanupBudget()

	cs := testCampaignService(t)
	campaigns, err := cs.Mutate(
		CampaignOperations{
			"ADD": {
				Campaign{
					Name:      "test campaign " + rand_str(10),
					Status:    "PAUSED",
					StartDate: time.Now().Format("20060102"),
					BudgetId:  budget.Id,
					Settings: []CampaignSetting{
						NewGeoTargetTypeSetting("DONT_CARE", "LOCATION_OF_PRESENCE"),
					},
					AdvertisingChannelType: "SEARCH",
					NetworkSetting: &NetworkSetting{
						TargetGoogleSearch:         true,
						TargetSearchNetwork:        true,
						TargetContentNetwork:       false,
						TargetPartnerSearchNetwork: false,
					},
					BiddingStrategyConfiguration: &BiddingStrategyConfiguration{
						StrategyType: "MANUAL_CPC",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(campaigns)
	log.Println(err)
	return

	defer func(campaigns []Campaign) {
		campaigns[0].Status = "REMOVED"
		_, err = cs.Mutate(CampaignOperations{"SET": campaigns})
		if err != nil {
			t.Error(err)
		}
	}(campaigns)

	label, labelCleanup := testLabel(t)
	defer labelCleanup()

	campaignLabels, err := cs.MutateLabel(
		CampaignLabelOperations{
			"ADD": {
				CampaignLabel{CampaignId: campaigns[0].Id, LabelId: label.Id},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		campaignLabels, err = cs.MutateLabel(CampaignLabelOperations{"REMOVE": campaignLabels})
		if err != nil {
			t.Fatal(err)
		}
	}()

	foundCampaigns, _, err := cs.Get(
		Selector{
			Fields: []string{
				"Id",
				"Name",
				"Status",
				"ServingStatus",
				"StartDate",
				"EndDate",
				"Settings",
				"Labels",
			},
			Predicates: []Predicate{
				{"Status", "EQUALS", []string{"PAUSED"}},
			},
			Ordering: []OrderBy{
				{"Id", "ASCENDING"},
			},
			Paging: &Paging{
				Offset: 0,
				Limit:  100,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("found %d campaigns\n", len(foundCampaigns))
	for _, c := range campaigns {
		func(campaign Campaign) {
			for _, foundCampaign := range foundCampaigns {
				if foundCampaign.Id == campaign.Id {
					fmt.Printf("%#v", foundCampaign)
					return
				}
			}
			t.Errorf("campaign %d not found in \n%#v\n", campaign.Id, foundCampaigns)
		}(c)
	}

}

func TestAppCampaign(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	budget, _ := testBudget(t)

	CreateAppCampaign(budget.Id, 160000, "YKTEST0902", "com.camera.momo.cam")
	return

}
func TestCampaignGET(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	//cs := testCampaignService(t)
	cs := testCampaignService(t)
	ads, totalCount, err := cs.Get(
		Selector{
			Fields: []string{
				"Name",
				"AdvertisingChannelSubType",
				/* 				"BiddingStrategyId",
				   				"BiddingStrategyName",
				   				"BiddingStrategyType",
				   				"TargetCpa",
				   				"Status",
				   				"ServingStatus", */
				"AdvertisingChannelType",
				//"Amount",
				"Id",
				"Settings",
				"TargetGoogleSearch",
				"TargetSearchNetwork",
				"TargetContentNetwork",
				"TargetPartnerSearchNetwork",
				"SelectiveOptimization",
				"BiddingStrategyId",
				"BiddingStrategyName",
				"BiddingStrategyType",
				"EnhancedCpcEnabled",
				"AdServingOptimizationStatus",
				"Eligible",
				"TargetSpendBidCeiling",
				"TargetSpendSpendTarget",
				//"BiddingStrategyConfiguration",
				//"AverageCpc",
				//"CpcBid",
			},
			Predicates: []Predicate{{"Name", "EQUALS", []string{"yyktest 1120 testupdate"}}},
		},
	)
	log.Println(err, totalCount, ads)
}
func TestCampaignGET2(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	//cs := testCampaignService(t)
	cs := testCampaignService(t)
	ads, totalCount, err := cs.Get(
		Selector{
			Fields: []string{
				"Name",
				"BudgetId",
				"Status",
				"Id",
			},
		},
	)
	log.Println(err, totalCount, ads[0].BudgetId)
}

//APP调价
func TestCampaignSet(t *testing.T) {

	cs := testCampaignService(t)
	_, err := cs.Mutate(
		CampaignOperations{
			"SET": {
				Campaign{
					Id: 10113627246,
					//Name:
					Status:                       "ENABLED",
					BiddingStrategyConfiguration: &BiddingStrategyConfiguration{Scheme: &BiddingScheme{Type: "TargetCpaBiddingScheme", TargetCpa: &TargetCpa{Amount: 190000}}},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	return
}

//feed
func TestCampaignFeed(t *testing.T) {
	var selector Selector
	selector = Selector{
		Fields: []string{
			"FeedId",
			"CampaignId",
			/* 				"BiddingStrategyId",
			"BiddingStrategyName",
			"BiddingStrategyType",
			"TargetCpa",
			"Status",
			"ServingStatus", */
			"MatchingFunction",
			//"Amount",
			"PlaceholderTypes",
			"Status",
			"BaseCampaignId",
		},
		//Predicates: []Predicate{{"Name", "EQUALS", []string{"Search-Clean wipes US-US-E-CPA-200513-PY"}}},
	}
	selector.XMLName = xml.Name{baseUrl, "selector"}
	s := testCampaignService(t)
	respBody, err := s.Auth.request(
		campaignFeedServiceUrl,
		"get",
		struct {
			XMLName xml.Name
			Sel     Selector
		}{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "get",
			},
			Sel: selector,
		},
	)
	log.Println(string(respBody))
	if err != nil {
		t.Fatal(err)
	}
	return
}

func TestCampaignConversionType(t *testing.T) {
	var selector Selector
	selector = Selector{
		Fields: []string{
			"Name",
			"Status",
			"Category",
			"CountingType",
		},
		//Predicates: []Predicate{{"Name", "EQUALS", []string{"Search-Clean wipes US-US-E-CPA-200513-PY"}}},
	}
	selector.XMLName = xml.Name{baseUrl, "selector"}
	s := &ConversionTrackerService{Auth: testAuthSetup(t)}
	s.Get(selector)
	/* 	log.Println(string(respBody))
	   	if err != nil {
	   		t.Fatal(err)
	   	} */
	return
}

/* 435205089 */
/* 434712695 */
