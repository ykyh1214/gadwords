package gadwords

import (
	"testing"
)

func testGAService(t *testing.T) (service *GAService) {
	return &GAService{Auth: testAuthSetup2(t)}
}

/* func TestGA(t *testing.T) {

	cs := testGAService(t)
	ctx := GAContext{
		ViewId:     viewid,
		StartDay:   start,
		EndDay:     end,
		Dimensions: []string{"ga:keyword", "adwordsAdGroupID", "ga:adwordsCampaignID", "ga:campaign"},
		Metrics:    []string{"ga:CPC", "ga:CTR", "ga:quantityAddedToCart", "ga:goal1Completions", "ga:goal1Value", "ga:adCost", "ga:adClicks"},
		Actid:      auth.CustomerId,
	}
	cs.Query("217869142", "2020-09-13")

} */
