package gadwords

import (
	"encoding/xml"
	"log"
)

type CustomerService struct {
	Auth
}

type Customer struct {
	CustomerId      int64  `xml:"customerId,omitempty"`
	CurrencyCode    string `xml:"currencyCode,omitempty"`
	DateTimeZone    string `xml:"dateTimeZone,omitempty"`
	DescriptiveName string `xml:"descriptiveName,omitempty"`
	/* 	CanManageClients           bool                       `xml:"canManageClients"`
	   	TestAccount                bool                       `xml:"testAccount"` */
	AutoTaggingEnabled         bool                        `xml:"autoTaggingEnabled"`
	TrackingUrlTemplate        string                      `xml:"trackingUrlTemplate,omitempty"`
	ConversionTrackingSettings *ConversionTrackingSettings `xml:"conversionTrackingSettings,omitempty"`
	RemarketingSettings        *RemarketingSettings        `xml:"remarketingSettings,omitempty"`
}
type ServiceLink struct {
	ServiceType   string `xml:"serviceType"`
	ServiceLinkId int64  `xml:"serviceLinkId"`
	LinkStatus    string `xml:"linkStatus"`
	Name          string `xml:"name"`
}

type RemarketingSettings struct {
	Snippet string `xml:"snippet"`
}

func NewCustomerService(auth *Auth) *CustomerService {
	return &CustomerService{Auth: *auth}
}

func (s *CustomerService) GetServiceLinks(selector Selector) (links []ServiceLink, err error) {
	selector.XMLName = xml.Name{baseMcmUrl, "selector"}
	respBody, err := s.Auth.request(
		customerServiceUrl,
		//"getCustomers",
		"getServiceLinks",
		struct {
			XMLName xml.Name
			Sel     Selector
		}{
			XMLName: xml.Name{
				Space: baseMcmUrl,
				Local: "getServiceLinks",
			},
			Sel: selector,
		},
	)
	if err != nil {
		return links, err
	}
	getResp := struct {
		Links []ServiceLink `xml:"rval"`
	}{}
	log.Println(string(respBody))
	err = xml.Unmarshal([]byte(respBody), &getResp)
	log.Println(err)
	if err != nil {
		return links, err
	}
	return getResp.Links, nil
}

func (s *CustomerService) GetCustomers() (customers []Customer, err error) {
	respBody, err := s.Auth.request(
		customerServiceUrl,
		//"getCustomers",
		"getCustomers",
		struct {
			XMLName xml.Name
		}{
			XMLName: xml.Name{
				Space: baseMcmUrl,
				Local: "getCustomers",
			},
		},
	)
	if err != nil {
		return customers, err
	}
	getResp := struct {
		Customers []Customer `xml:"rval"`
	}{}
	log.Println(string(respBody))
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return customers, err
	}
	return getResp.Customers, nil
}

type CustomerOperations map[string]Customer

func (s *CustomerService) Mutate(customerOperations CustomerOperations) (customers []Customer, err error) {
	/* 	type customerOperation struct {
	   		Action   string   `xml:"operator"`
	   		Customer Customer `xml:"operand"`
	   	}
	   	operations := []customerOperation{}
	   	for action, customers := range customerOperations {
	   		for _, customer := range customers {
	   			operations = append(operations,
	   				customerOperation{
	   					Action:   action,
	   					Customer: customer,
	   				},
	   			)
	   		}
	   	} */
	mutation := struct {
		XMLName xml.Name
		Ops     Customer `xml:"customer"`
	}{
		XMLName: xml.Name{
			Space: baseMcmUrl,
			Local: "mutate",
		},
		Ops: customerOperations["SET"]}

	respBody, err := s.Auth.request(customerServiceUrl, "mutate", mutation)
	if err != nil {
		return customers, err
	}
	getResp := struct {
		Customers []Customer `xml:"rval"`
	}{}
	log.Println(string(respBody))
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return customers, err
	}
	return getResp.Customers, nil
}

/* type adGroupAdOperation struct {
	Action    string     `xml:"operator"`
	AdGroupAd AdGroupAds `xml:"operand"`
}
operations := []adGroupAdOperation{}
for action, adGroupAds := range adGroupAdOperations {
	for _, adGroupAd := range adGroupAds {
		ad := []interface{}{adGroupAd}
		operations = append(operations,
			adGroupAdOperation{
				Action:    action,
				AdGroupAd: ad,
			},
		)
	}
}
mutation := struct {
	XMLName xml.Name
	Ops     []adGroupAdOperation `xml:"operations"`
}{
	XMLName: xml.Name{
		Space: baseUrl,
		Local: "mutate",
	},
	Ops: operations,
}

respBody, err := s.Auth.request(adGroupAdServiceUrl, "mutate", mutation)
if err != nil {
	return adGroupAds, err
}
mutateResp := struct {
	AdGroupAds AdGroupAds `xml:"rval>value"`
}{}
log.Println(string(respBody))
err = xml.Unmarshal([]byte(respBody), &mutateResp)
if err != nil {
	return adGroupAds, err
}
return mutateResp.AdGroupAds, err */
