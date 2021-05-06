package gadwords

import (
	"encoding/xml"
	"log"
	"testing"
)

func testAdGroupAdService(t *testing.T) (service *AdGroupAdService) {
	return &AdGroupAdService{Auth: testAuthSetup(t)}
}

func TestAdGroupAdGET(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	agas := testAdGroupAdService(t)
	ads, totalCount, err := agas.Get(
		Selector{
			Fields: []string{
				"AdGroupId",
				"Id",
			},
			Predicates: []Predicate{
				{"AdType", "EQUALS", []string{"RESPONSIVE_SEARCH_AD"}},
			},
		},
	)
	log.Println(err, totalCount, ads)
}

/* func TestAdGroupAd(t *testing.T) {
	adGroup, cleanupAdGroup := testAdGroup(t)
	defer cleanupAdGroup()

	agas := testAdGroupAdService(t)
	adGroupAds, err := agas.Mutate(
		AdGroupAdOperations{
			"ADD": {
				NewTextAd(adGroup.Id, "https://classdo.com/en", "classdo.com", "test headline "+rand_word(10), "test line one", "test line two", "PAUSED"),
				NewTextAd(adGroup.Id, "https://classdo.com/en", "classdo.com", "test   teStTo "+rand_word(10), "test line one", "test line two", "PAUSED"),
				NewTextAd(adGroup.Id, "https://classdo.com/en", "classdo.com", "test headline "+rand_word(10), "test line one", "test line two", "PAUSED"),
				NewTextAd(adGroup.Id, "https://classdo.com/en", "classdo.com", "test headline "+rand_word(10), "test line one", "test line two", "PAUSED"),
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
		defer func() {
		_, err = agas.Mutate(AdGroupAdOperations{"REMOVE": adGroupAds})
		if err != nil {
			t.Error(err)
		}
	}()

	adGroupIdStr := fmt.Sprintf("%d", adGroup.Id)
	_, _, err = agas.Get(
		Selector{
			Fields: []string{
				"AdGroupId",
				"Id",
				"Status",
				"AdGroupCreativeApprovalStatus",
				"AdGroupAdDisapprovalReasons",
				"AdGroupAdTrademarkDisapproved",
			},
			Predicates: []Predicate{
				{"AdGroupId", "EQUALS", []string{adGroupIdStr}},
			},
			Ordering: []OrderBy{
				{"AdGroupId", "ASCENDING"},
				{"Id", "ASCENDING"},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

}
*/
func TestAppAd(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	headlines := []AssetLink{{Asset: Asset{Type: "TextAsset", AssetText: "headline1"}}, {Asset: Asset{Type: "TextAsset", AssetText: "headline2"}}}
	descriptions := []AssetLink{{Asset: Asset{Type: "TextAsset", AssetText: "descriptions1"}}, {Asset: Asset{Type: "TextAsset", AssetText: "description2"}}}
	images := []AssetLink{{Asset: Asset{Type: "ImageAsset", Id: 9006120670}}, {Asset: Asset{Type: "ImageAsset", Id: 9006120673}}}
	agas := testAdGroupAdService(t)
	adGroupAds, err := agas.Mutate(
		AdGroupAdOperations{
			"ADD": {
				NewAppAd(109054157675, headlines, descriptions, images, nil),
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(adGroupAds, err)

}
func TestResponsiveSearchAd(t *testing.T) {
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
				ExpandedTextAd{
					AdGroupId:     1234567890,
					FinalUrls:     []string{"https://classdo.com/en"},
					Path1:         "path1",
					Path2:         "path2",
					HeadlinePart1: "test headline",
					HeadlinePart2: "test headline2",
					HeadlinePart3: "test headline3",
					Description:   "test line one",
					Description2:  "test line one2",
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(adGroupAds, err)

}

var xmldata = `      <mutateResponse
xmlns="https://adwords.google.com/api/adwords/cm/v201809">
<rval>
	<ListReturnValue.Type>AdGroupAdReturnValue</ListReturnValue.Type>
	<value>
		<adGroupId>105449240142</adGroupId>
		<ad
			xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="UniversalAppAd">
			<id>461697449754</id>
			<type>UNIVERSAL_APP_AD</type>
			<Ad.Type>UniversalAppAd</Ad.Type>
			<headlines>
				<asset xsi:type="TextAsset">
					<assetId>10391090888</assetId>
					<assetSubtype>TEXT</assetSubtype>
					<Asset.Type>TextAsset</Asset.Type>
					<assetText>headline1</assetText>
				</asset>
			</headlines>
			<headlines>
				<asset xsi:type="TextAsset">
					<assetId>10391090885</assetId>
					<assetSubtype>TEXT</assetSubtype>
					<Asset.Type>TextAsset</Asset.Type>
					<assetText>headline2</assetText>
				</asset>
			</headlines>
			<descriptions>
				<asset xsi:type="TextAsset">
					<assetId>10391090891</assetId>
					<assetSubtype>TEXT</assetSubtype>
					<Asset.Type>TextAsset</Asset.Type>
					<assetText>descriptions1</assetText>
				</asset>
			</descriptions>
			<descriptions>
				<asset xsi:type="TextAsset">
					<assetId>10391090894</assetId>
					<assetSubtype>TEXT</assetSubtype>
					<Asset.Type>TextAsset</Asset.Type>
					<assetText>description2</assetText>
				</asset>
			</descriptions>
			<images>
				<asset xsi:type="ImageAsset">
					<assetId>9006120670</assetId>
					<assetName>11182136_MOMO FILTER (1).jpg</assetName>
					<assetSubtype>IMAGE</assetSubtype>
					<Asset.Type>ImageAsset</Asset.Type>
					<imageFileSize>143392</imageFileSize>
					<imageMimeType>IMAGE_JPEG</imageMimeType>
					<fullSizeInfo>
						<imageHeight>628</imageHeight>
						<imageWidth>1200</imageWidth>
					</fullSizeInfo>
				</asset>
			</images>
			<images>
				<asset xsi:type="ImageAsset">
					<assetId>9006120673</assetId>
					<assetName>11182137_MOMO FILTER (8).jpg</assetName>
					<assetSubtype>IMAGE</assetSubtype>
					<Asset.Type>ImageAsset</Asset.Type>
					<imageFileSize>128951</imageFileSize>
					<imageMimeType>IMAGE_JPEG</imageMimeType>
					<fullSizeInfo>
						<imageHeight>628</imageHeight>
						<imageWidth>1200</imageWidth>
					</fullSizeInfo>
				</asset>
			</images>
		</ad>
		<status>ENABLED</status>
	</value>
</rval>
</mutateResponse>`

func TestXml(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	unmalxml()
}

func unmalxml() {
	mutateResp := struct {
		AdGroupAds AdGroupAds `xml:"rval>value"`
	}{}
	err := xml.Unmarshal([]byte(xmldata), &mutateResp)
	log.Println(err, mutateResp.AdGroupAds)
}

func TestDisableAd(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	agas := testAdGroupAdService(t)
	adGroupAds, err := agas.Mutate(
		AdGroupAdOperations{
			"SET": {
				OptAd{Id: 4378231191538, AdGroupId: 105299758481, Status: "PAUSED"},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(adGroupAds, err)

}
