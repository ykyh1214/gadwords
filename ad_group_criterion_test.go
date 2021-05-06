package gadwords

import (
	"encoding/json"
	"log"
	"testing"
	//	"encoding/xml"
)

func testAdGroupCriterionService(t *testing.T) (service *AdGroupCriterionService) {
	return &AdGroupCriterionService{Auth: testAuthSetup22(t)}
}

type nkw struct {
	AdGroupId    string
	CriterionUse string
	Criterion    struct {
		Id        int64
		Text      string
		MatchType string
	}
}

func TestGetAdGroupCriterion2(t *testing.T) {

	agcs := testAdGroupCriterionService(t)
	ad, _, _ := agcs.Get(
		Selector{
			Fields: []string{
				"Id",
				"CriteriaType",
				"Status",
				"KeywordText",
				//"BaseCampaignId",
				//"Bids",
				//"BiddingScheme",
			},
			Predicates: []Predicate{
				{"CriteriaType", "EQUALS", []string{"KEYWORD"}},
				{"CriterionUse", "EQUALS", []string{"NEGATIVE"}},
				{"BaseCampaignId", "EQUALS", []string{"11679056322"}},
				//{"KeywordText", "IN", []string{"+dog +bowls", "+cat +bowls"}},
			},
		},
	)
	js, _ := json.MarshalIndent(ad, "", " ")

	log.Println(string(js))

	a := ad[0].(NegativeAdGroupCriterion)
	kw := a.Criterion.(KeywordCriterion)

	log.Println(a, kw.Text)
	/* 	var ads AdGroupCriterions
	   	for _, a := range ad {
	   		b := a.(BiddableAdGroupCriterion)
	   		b.UserStatus = "ENABLED"
	   		ads = append(ads, b)
	   	}

	   	log.Println(agcs.Mutate(
	   		AdGroupCriterionOperations{
	   			"SET": ads,
	   		},
	   	)) */
}

func TestAdGroupCriterion(t *testing.T) {
	adGroup, cleanupAdGroup := testAdGroup(nil)
	defer cleanupAdGroup()

	agcs := testAdGroupCriterionService(t)
	adGroupCriterions, err := agcs.Mutate(
		AdGroupCriterionOperations{
			"ADD": {
				/*
									NegativeAdGroupCriterion{
					          AdGroupId: adGroup.Id,
					          Criterion: AgeRangeCriterion{AgeRangeType:"AGE_RANGE_25_34"},
					        },
									NegativeAdGroupCriterion{
					          AdGroupId: adGroup.Id,
					          Criterion: GenderCriterion{},
					        },
									NewBiddableAdGroupCriterion{
					          AdGroupId: adGroup.Id,
					          Criterion: MobileAppCategoryCriterion{
					            60000,"My Google Play Android Apps"
					          }
					        },
				*/
				BiddableAdGroupCriterion{
					AdGroupId:  adGroup.Id,
					Criterion:  KeywordCriterion{Text: "test1", MatchType: "EXACT"},
					UserStatus: "PAUSED",
				},
				BiddableAdGroupCriterion{
					AdGroupId:  adGroup.Id,
					Criterion:  KeywordCriterion{Text: "test2", MatchType: "PHRASE"},
					UserStatus: "PAUSED",
				},
				BiddableAdGroupCriterion{
					AdGroupId:  adGroup.Id,
					Criterion:  KeywordCriterion{Text: "test3", MatchType: "BROAD"},
					UserStatus: "PAUSED",
				},
				NegativeAdGroupCriterion{
					AdGroupId: adGroup.Id,
					Criterion: KeywordCriterion{Text: "test4", MatchType: "BROAD"},
				},
				BiddableAdGroupCriterion{
					AdGroupId:  adGroup.Id,
					Criterion:  PlacementCriterion{Url: "https://classdo.com"},
					UserStatus: "PAUSED",
				},
				// NewBiddableAdGroupCriterion(adGroup.Id, NewUserInterestCriterion()),
				// NewBiddableAdGroupCriterion(adGroup.Id, NewUserListCriterion()),
				// NewBiddableAdGroupCriterion(adGroup.Id, NewVerticalCriterion(0, 0, []string{"Pets & Anamals","Pets","Dogs"})),
				BiddableAdGroupCriterion{
					AdGroupId: adGroup.Id,
					Criterion: WebpageCriterion{
						Parameter: WebpageParameter{
							CriterionName: "test criterion",
							Conditions: []WebpageCondition{
								WebpageCondition{
									Operand:  "URL",
									Argument: "example.com",
								},
							},
						},
					},
					UserStatus: "PAUSED",
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", adGroupCriterions)

	defer func() {
		_, err = agcs.Mutate(AdGroupCriterionOperations{"REMOVE": adGroupCriterions})
		if err != nil {
			t.Error(err)
		}
	}()
	/*
	   reqBody, err := xml.MarshalIndent(adGroupCriterions,"  ", "  ")
	   t.Fatalf("%s\n",reqBody)
	*/
}

func TestGetAdGroupCriterion(t *testing.T) {

	agcs := testAdGroupCriterionService(t)
	agcs.Get(
		Selector{
			Fields: []string{
				"Id",
				"CriteriaType",
				"PartitionType",
				"ParentCriterionId",
				"CaseValue",
				"Status",
				"CpcBid",
				"KeywordText",
				"BiddingStrategyType",
				//"Bids",
				//"BiddingScheme",
			},
			Predicates: []Predicate{
				{"AdGroupId", "EQUALS", []string{"105430503550"}},
			},
		},
	)
}

func TestAdGroupCriterionSet(t *testing.T) {
	agcs := testAdGroupCriterionService(t)
	adGroupCriterions, err := agcs.Mutate(
		AdGroupCriterionOperations{
			"SET": {
				BiddableAdGroupCriterion{
					AdGroupId: 101412673385,
					Criterion: KeywordCriterion{Id: 606804959267},
					BiddingStrategyConfiguration: &BiddingStrategyConfiguration{
						Bids: []Bid{
							Bid{
								Amount: 250000,
							},
						},
					},
				},
			},
		},
	)
	log.Println(err)
	log.Println(adGroupCriterions, err)
	log.Println(err)
}

//否词操作
func TestAdGroupCriterionNegative(t *testing.T) {
	agcs := testAdGroupCriterionService(t)
	adGroupCriterions, err := agcs.Mutate(
		AdGroupCriterionOperations{
			"SET": {
				NegativeAdGroupCriterion{
					AdGroupId: 101412673385,
					Criterion: KeywordCriterion{Id: 606804959267},
				},
			},
		},
	)
	log.Println(err)
	log.Println(adGroupCriterions, err)
	log.Println(err)
}
