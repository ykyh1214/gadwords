package gadwords

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestSearchCampaign(t *testing.T) {
	CreateSearchCampaign("yykserach test2")
}

func TestSearchAdgroup(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	//cs := testCampaignService(t)
	auth := testAuthSetup(t)
	AddSearchAdGroup(&auth, 11585961744)
}

func TestAddSearchKeywords(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	//cs := testCampaignService(t)
	auth := testAuthSetup(t)
	AddSearchKeywords(&auth, 112993761796)
}
func TestAddSearchAd(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	headlines := []AssetLink{
		{Asset: Asset{Type: "TextAsset", AssetText: "headline1"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline2"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline3"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline4"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline5"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline6"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline7"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline8"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline9"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline10"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline11"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline12"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline13"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline14"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "headline15"}},
	}
	descriptions := []AssetLink{
		{Asset: Asset{Type: "TextAsset", AssetText: "descriptions1"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "description2"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "description3"}},
		{Asset: Asset{Type: "TextAsset", AssetText: "description4"}},
	}
	agas := testAdGroupAdService(t)
	adGroupAds, err := agas.Mutate(
		AdGroupAdOperations{
			"ADD": {
				ResponsiveSearchAd{
					AdGroupId:    112993761796,
					FinalUrls:    []string{"https://www.worth2own.com/us/en/products/22-alcohol-wet-wipes-80"},
					Headlines:    headlines,
					Descriptions: descriptions,
					Path1:        "Hand Sanitizer",
					Path2:        "Alcohol Wipes",
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(adGroupAds, err)

}
func TestAccountAuottargeting(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	//cs := testCampaignService(t)
	/* layout := fmt.Sprintf("2006-01-02 15:%02d:05", 22)
	log.Println(layout) */
	AmericaLocation := time.FixedZone("America", -5*3600)
	location := time.FixedZone("ccc", int(8*3600))
	layout := fmt.Sprintf("2006-01-02 %02d:04:05", 11)

	startT, _ := time.ParseInLocation("2006-01-02 15:04:05", layout, AmericaLocation)
	log.Println(startT.In(location).Hour())
	return
	auth := testAuthSetup(t)
	agcs := NewCustomerService(&auth)
	coustoms, err := agcs.Mutate(
		CustomerOperations{
			"SET": Customer{AutoTaggingEnabled: true},
		},
	)
	log.Println(coustoms, err)
}
