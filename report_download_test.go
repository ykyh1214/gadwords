package gadwords

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type TestClient struct {
	res *http.Response
}

func (s *TestClient) Do(req *http.Request) (*http.Response, error) {
	return s.res, nil
}

func TestReportDownloadAuthError(t *testing.T) {
	body := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?> <reportDownloadError> <ApiError> <type>AuthorizationError.USER_PERMISSION_DENIED</type> <trigger>&lt;null&gt;</trigger> <fieldPath></fieldPath> </ApiError> </reportDownloadError>`)
	client := &TestClient{
		res: &http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
			StatusCode: 400,
		},
	}

	auth := &Auth{
		Client: client,
	}

	rs := NewReportDownloadService(auth)
	def := ReportDefinition{}
	_, err := rs.Get(def)
	if err == nil {
		t.Fatalf("expected api error")
	}

	type ErrorCode interface {
		Code() string
	}

	expectedCode := "USER_PERMISSION_DENIED"
	if ec, ok := err.(ErrorCode); ok {
		if ec.Code() != expectedCode {
			t.Errorf("got %s, expected %s\n", ec.Code(), expectedCode)
		}
	} else {
		t.Errorf("error expected to satisfy ErrorCode interface")
	}

}

func TestReportQueryStream(t *testing.T) {
	query := `SELECT  AccountDescriptiveName, AdvertisingChannelType, Clicks, ConversionValue, Cost, Impressions, Device, ExternalCustomerId, DayOfWeek, CampaignId  FROM CAMPAIGN_PERFORMANCE_REPORT  DURING YESTERDAY`
	config := getTestConfig()
	svc := NewReportDownloadService(&config.Auth)

	report, err := svc.StreamAWQL(query, "CSV")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(report)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}

func TestReportStreamDownloadAuthError(t *testing.T) {
	body := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?> <reportDownloadError> <ApiError> <type>AuthorizationError.USER_PERMISSION_DENIED</type> <trigger>&lt;null&gt;</trigger> <fieldPath></fieldPath> </ApiError> </reportDownloadError>`)
	client := &TestClient{
		res: &http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
			StatusCode: 400,
		},
	}

	auth := &Auth{
		Client: client,
	}

	rs := NewReportDownloadService(auth)
	_, err := rs.StreamAWQL("", "")
	if err == nil {
		t.Fatalf("expected api error")
	}

	type ErrorCode interface {
		Code() string
	}

	expectedCode := "USER_PERMISSION_DENIED"
	if ec, ok := err.(ErrorCode); ok {
		if ec.Code() != expectedCode {
			t.Errorf("got %s, expected %s\n", ec.Code(), expectedCode)
		}
	} else {
		t.Errorf("error expected to satisfy ErrorCode interface")
	}

}

func TestReportQuery2(t *testing.T) {
	query := `SELECT  CampaignId,CampaignName,GoogleClickId FROM CAMPAIGN_PERFORMANCE_REPORT where CampaignId=10116433597 `
	//config := getTestConfig()
	svc := &ReportDownloadService{Auth: testAuthSetup(t)}

	report, err := svc.StreamAWQL(query, "CSV")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(report)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}

//WHERE CampaignStatus IN [ENABLED]
func TestReportQuery(t *testing.T) {
	query := `SELECT  AdvertisingChannelSubType,AdvertisingChannelType,CampaignId,BiddingStrategyType,Amount , CampaignName, AverageCpc, Ctr, Cost, Conversions, CostPerConversion, Clicks, CampaignStatus  FROM CAMPAIGN_PERFORMANCE_REPORT DURING 20200917,20200917`
	//config := getTestConfig()
	svc := &ReportDownloadService{Auth: testAuthSetup(t)}

	report, err := svc.StreamAWQL(query, "CSV")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(report)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}

//keyword 的数据只能用adwords来拿,因为GA上没id
func TestKeywordReportQuery(t *testing.T) {
	query := `SELECT  Criteria,BiddingStrategyType,SearchImpressionShare,HistoricalQualityScore,HistoricalLandingPageQualityScore,HistoricalSearchPredictedCtr,HistoricalCreativeQualityScore FROM KEYWORDS_PERFORMANCE_REPORT WHERE AdGroupId IN [101412673665] DURING TODAY`
	//config := getTestConfig()
	svc := &ReportDownloadService{Auth: testAuthSetup(t)}

	report, err := svc.StreamAWQL(query, "CSV")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(report)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}

//search term 数据只能用adwords来拿
func TestSearchTermReportQuery(t *testing.T) {
	query := `SELECT   CampaignName,Cost,Conversions,ConversionValue, Impressions, Ctr  FROM CAMPAIGN_PERFORMANCE_REPORT WHERE Cost>0 DURING 20201026,20201122`
	//config := getTestConfig()
	svc := &ReportDownloadService{Auth: testAuthSetup(t)}

	report, err := svc.StreamAWQL(query, "CSV")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(report)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}

//走不通 clickshare只能从python ads中拿
//SearchClickShare
func TestAdGroupReportQuery(t *testing.T) {
	query := `SELECT  AdGroupId,CampaignId,BiddingStrategyType,SearchImpressionShare  FROM ADGROUP_PERFORMANCE_REPORT WHERE AdGroupId IN [100292335125,100292334925] DURING 20200913,20200913`
	//config := getTestConfig()
	svc := &ReportDownloadService{Auth: testAuthSetup(t)}

	report, err := svc.StreamAWQL(query, "CSV")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(report)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}

func TestBidReportQuery(t *testing.T) {
	query := `SELECT  Id,Name,TargetCpa  FROM BID_GOAL_PERFORMANCE_REPORT  DURING YESTERDAY`
	//config := getTestConfig()
	svc := &ReportDownloadService{Auth: testAuthSetup(t)}

	report, err := svc.StreamAWQL(query, "CSV")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(report)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}
func TestAdsReportQuery(t *testing.T) {
	s := &ReportDownloadService{Auth: testAuthSetup(t)}
	resp, err := s.makeadsRequest()
	fmt.Println(resp, err)
}

func TestKeyword(t *testing.T) {
	auth := testAuthSetup(t)
	cs := NewAdGroupCriterionService(&auth)
	fmt.Printf("Keywords\n")
	foundKeywords, totalCount, err := cs.Query("SELECT Id, KeywordText, KeywordMatchType WHERE AdGroupId = '104366013084'")
	fmt.Println(totalCount)
	if err != nil {
		log.Fatal(err)
	}
	for _, keyword := range foundKeywords {
		keywordJSON, _ := json.MarshalIndent(keyword, "", "  ")
		fmt.Printf("%s\n", keywordJSON)
	}
}

func TestSearchItems(t *testing.T) {
	query := `SELECT 
	Query,
	KeywordTextMatchingQuery, 
	Cost,
	Conversions,
	Clicks,
	Impressions,
	QueryTargetingStatus,
	InteractionRate,
	QueryMatchTypeWithVariant
	  FROM SEARCH_QUERY_PERFORMANCE_REPORT` //config := getTestConfig()
	svc := &ReportDownloadService{Auth: testAuthSetup22(t)}

	report, err := svc.StreamAWQL(query, "CSV")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(report)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}

func TestKeywords(t *testing.T) {
	query := `SELECT  CampaignId, CampaignName, Cost, Conversions, Clicks  FROM CAMPAIGN_PERFORMANCE_REPORT WHERE AdvertisingChannelSubType IN [UNIVERSAL_APP_CAMPAIGN] AND Cost>0  DURING 20210225,20210226 `
	//config := getTestConfig()
	svc := &ReportDownloadService{Auth: testAuthSetup22(t)}

	report, err := svc.StreamAWQL(query, "CSV")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(report)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}
