package gadwords

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
	//  "encoding/xml"
)

func testCampaignCriterionService(t *testing.T) (service *CampaignCriterionService) {
	return &CampaignCriterionService{Auth: testAuthSetup22(t)}
}

func TestCampaignCriterion(t *testing.T) {
	/* campaign, cleanupCampaign := testCampaign(t)
	defer cleanupCampaign() */

	ccs := testCampaignCriterionService(t)
	_, err := ccs.Mutate(
		CampaignCriterionOperations{
			"ADD": {
				//CampaignCriterion{CampaignId: campaign.Id, Criterion: AdScheduleCriterion{DayOfWeek: "MONDAY", StartHour: "10", StartMinute: "ZERO", EndHour: "13", EndMinute: "ZERO"}},
				CampaignCriterion{CampaignId: 11023075160, Criterion: Location{Id: 2764}},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	//t.Fatalf("%#v\n",campaignCriterions)

	/* 	defer func() {
		_, err = ccs.Mutate(CampaignCriterionOperations{"REMOVE": campaignCriterions})
		if err != nil {
			t.Error(err)
		}
	}() */
}

func TestCampaignCriterion_Get(t *testing.T) {
	// load credentials from
	cs := testCampaignCriterionService(t)
	// This example illustrates how to retrieve all the campaigns for an account.
	var pageSize int64 = 500
	var offset int64 = 0
	paging := Paging{
		Offset: offset,
		Limit:  pageSize,
	}
	totalCount := 0
	for {
		campaignCriterions, totalCount, err := cs.Get(
			Selector{
				Fields: []string{
					"CampaignId",
					"CriteriaType",
					"DayOfWeek",
					"StartHour",
					"EndHour",
				},
				Paging: &paging,
				Predicates: []Predicate{{"CampaignName", "EQUALS", []string{"yykserach test"}},
					{"CriteriaType", "EQUALS", []string{"AD_SCHEDULE"}},
				},
			},
		)
		if err != nil {
			log.Println(err)
			fmt.Printf("Error occured finding campaigns")
		}
		js, _ := json.MarshalIndent(campaignCriterions, "", " ")
		log.Println(string(js))
		// Increment values to request the next page.
		offset += pageSize
		paging.Offset = offset
		if totalCount < offset {
			break
		}
	}
	fmt.Printf("\tTotal number of campaigns found: %d.", totalCount)
}

func TestCampaignCriterionSet(t *testing.T) {
	/* campaign, cleanupCampaign := testCampaign(t)
	defer cleanupCampaign() */

	ccs := testCampaignCriterionService(t)

	log.Println(ccs.SetSchedule(11585961744, 12, 18))
	return
	_, err := ccs.Mutate(
		CampaignCriterionOperations{
			"SET": {
				//CampaignCriterion{CampaignId: campaign.Id, Criterion: AdScheduleCriterion{DayOfWeek: "MONDAY", StartHour: "10", StartMinute: "ZERO", EndHour: "13", EndMinute: "ZERO"}},
				CampaignCriterion{CampaignId: 11585961744, Criterion: AdScheduleCriterion{Id: 303236, EndHour: "11"}},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	//t.Fatalf("%#v\n",campaignCriterions)

	/* 	defer func() {
		_, err = ccs.Mutate(CampaignCriterionOperations{"REMOVE": campaignCriterions})
		if err != nil {
			t.Error(err)
		}
	}() */
}

func TestCampaignCriterion_Get2(t *testing.T) {
	// load credentials from
	cs := testCampaignCriterionService(t)

	campaignCriterions, totalCount, err := cs.Get(
		Selector{
			Fields: []string{
				"CampaignId",
				"CriteriaType",
				"KeywordText",
				"IsNegative",
				"CampaignCriterionStatus",
			},
			Predicates: []Predicate{
				{"IsNegative", "EQUALS", []string{"true"}},
				{"CriteriaType", "EQUALS", []string{"KEYWORD"}},
				{"CampaignId", "EQUALS", []string{"11679056322"}},
			},
		},
	)

	js, _ := json.MarshalIndent(campaignCriterions, "", " ")
	log.Println(string(js), err)

	fmt.Printf("\tTotal number of campaigns found: %d.", totalCount)
}

var AmericaLocation = time.FixedZone("America", -5*3600)

func TestCampaignCriterionsss(t *testing.T) {
	log.Println(parseStartEnd(12, 24))

}

func getStartEnd(start, end int) (int, int) {
	if start == 24 {
		start = 0
	}
	if end == 24 {
		end = 0
	}
	location := time.FixedZone("", int(-3*3600))
	startLayout := fmt.Sprintf("2006-01-02 %02d:04:05", start)
	startT, _ := time.ParseInLocation("2006-01-02 15:04:05", startLayout, AmericaLocation)

	endLayout := fmt.Sprintf("2006-01-02 %02d:04:05", end)
	endT, _ := time.ParseInLocation("2006-01-02 15:04:05", endLayout, AmericaLocation)

	startHour := startT.In(location).Hour()
	endHour := endT.In(location).Hour()
	if endHour == 0 && startHour != 0 {
		endHour = 24
	}
	return startHour, endHour
}

func parseStartEnd(start, end int) (int, int) {
	if end == 24 {
		end = 0
	}

	location := time.FixedZone("", int(-5*3600))
	startLayout := fmt.Sprintf("2006-01-02 %02d:04:05", start)
	startT, _ := time.ParseInLocation("2006-01-02 15:04:05", startLayout, location)
	endLayout := fmt.Sprintf("2006-01-02 %02d:04:05", end)
	endT, _ := time.ParseInLocation("2006-01-02 15:04:05", endLayout, location)

	return startT.In(AmericaLocation).Hour(), endT.In(AmericaLocation).Hour()
}
