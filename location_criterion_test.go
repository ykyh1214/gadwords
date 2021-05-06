package gadwords

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
	//  "encoding/xml"
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)

}

func testLocationCriterionService(t *testing.T) (service *LocationCriterionService) {
	return &LocationCriterionService{Auth: testAuthSetup(t)}
}

func TestLocationCriterion(t *testing.T) {
	lcs := testLocationCriterionService(t)
	locationCriterions, err := lcs.Get(
		Selector{
			Fields: []string{
				"Id",
				"LocationName",
				"CanonicalName",
				"DisplayType",
				"CountryCode",
				"Reach",
				//"Locale",
			},
			Predicates: []Predicate{
				{"Id", "EQUALS", []string{"2156"}},
				//{"CountryCode", "EQUALS", []string{"CN"}},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	for _, locationCriterion := range locationCriterions {
		location := locationCriterion.Location
		js, _ := json.MarshalIndent(location, "", " ")
		log.Println(string(js))
		fmt.Printf("%d. %s, %s\n", location.Id, location.LocationName, location.DisplayType)
	}
}

func TestLink(t *testing.T) {
	lcs := &CustomerService{Auth: testAuthSetup(t)}
	links, err := lcs.GetServiceLinks(
		Selector{
			Fields: []string{
				"ServiceType",
				"ServiceLinkId",
				"LinkStatus",
				"Name",
			},
			Predicates: []Predicate{
				{"ServiceType", "IN", []string{"MERCHANT_CENTER"}},
			},
		},
	)
	log.Println(links, err)
	if err != nil {
		t.Error(err)
	}
}

func TestCustomers(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	l, _ := time.LoadLocation("America/Chicago")
	log.Println(time.Now().In(l).Zone())
	return

	lcs := &CustomerService{Auth: testAuthSetup(t)}
	locationCriterions, err := lcs.GetCustomers()
	if err != nil {
		t.Error(err)
	}
	for _, locationCriterion := range locationCriterions {

		log.Println(locationCriterion)
	}
}
