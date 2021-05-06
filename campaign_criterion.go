package gadwords

import (
	"encoding/xml"
	"fmt"

	"github.com/pkg/errors"
)

type CampaignCriterionService struct {
	Auth
}

func NewCampaignCriterionService(auth *Auth) *CampaignCriterionService {
	return &CampaignCriterionService{Auth: *auth}
}

type CampaignCriterion struct {
	CampaignId  int64     `xml:"campaignId"`
	IsNegative  bool      `xml:"isNegative,omitempty"`
	Criterion   Criterion `xml:"criterion"`
	BidModifier *float64  `xml:"bidModifier,omitempty"`
	Errors      []error   `xml:"-"`
}

func (cc CampaignCriterion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	isNegative := false
	//fmt.Printf("processing -> %#v\n",ncc)
	start.Attr = append(
		start.Attr,
		xml.Attr{
			xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"},
			"CampaignCriterion",
		},
	)
	e.EncodeToken(start)
	e.EncodeElement(&cc.CampaignId, xml.StartElement{Name: xml.Name{"", "campaignId"}})
	e.EncodeElement(&isNegative, xml.StartElement{Name: xml.Name{"", "isNegative"}})
	if err := criterionMarshalXML(cc.Criterion, e); err != nil {
		return err
	}
	if cc.BidModifier != nil {
		e.EncodeElement(&cc.BidModifier, xml.StartElement{Name: xml.Name{"", "bidModifier"}})
	}

	e.EncodeToken(start.End())
	return nil
}

type NegativeCampaignCriterion struct {
	CampaignId  int64     `xml:"campaignId"`
	IsNegative  bool      `xml:"isNegative,omitempty"`
	Criterion   Criterion `xml:"criterion"`
	BidModifier *float64  `xml:"bidModifier,omitempty"`
	Errors      []error   `xml:"-"`
}

type CampaignCriterions []interface{}
type CampaignCriterionOperations map[string]CampaignCriterions

func (ncc NegativeCampaignCriterion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	isNegative := true
	//fmt.Printf("processing -> %#v\n",ncc)
	start.Attr = append(
		start.Attr,
		xml.Attr{
			xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"},
			"NegativeCampaignCriterion",
		},
	)
	e.EncodeToken(start)
	e.EncodeElement(&ncc.CampaignId, xml.StartElement{Name: xml.Name{"", "campaignId"}})
	e.EncodeElement(&isNegative, xml.StartElement{Name: xml.Name{"", "isNegative"}})
	criterionMarshalXML(ncc.Criterion, e)
	e.EncodeToken(start.End())
	return nil
}

func (ccs *CampaignCriterions) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	cc := NegativeCampaignCriterion{}
	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		if err != nil {
			return err
		}
		switch start := token.(type) {
		case xml.StartElement:
			switch start.Name.Local {
			case "campaignId":
				if err := dec.DecodeElement(&cc.CampaignId, &start); err != nil {
					return err
				}
			case "criterion":
				criterion, err := criterionUnmarshalXML(dec, start)
				if err != nil {
					return err
				}
				cc.Criterion = criterion
			case "bidModifier":
				if err := dec.DecodeElement(&cc.BidModifier, &start); err != nil {
					return err
				}
			case "isNegative":
				if err := dec.DecodeElement(&cc.IsNegative, &start); err != nil {
					return err
				}
			}
		}
	}
	*ccs = append(*ccs, cc)
	return nil
}

/*
func NewNegativeCampaignCriterion(campaignId int64, bidModifier float64, criterion interface{}) CampaignCriterion {
  return CampaignCriterion{
    CampaignId: campaignId,
    Criterion: criterion,
    BidModifier: bidModifier
  }
  switch c := criterion.(type) {
  case AdScheduleCriterion:
  case AgeRangeCriterion:
  case ContentLabelCriterion:
  case GenderCriterion:
  case KeywordCriterion:
  case LanguageCriterion:
  case LocationCriterion:
  case MobileAppCategoryCriterion:
  case MobileApplicationCriterion:
  case MobileDeviceCriterion:
  case OperatingSystemVersionCriterion:
  case PlacementCriterion:
  case PlatformCriterion:
  case ProductCriterion:
  case ProximityCriterion:
  case UserInterestCriterion:
    cc.Criterion = criterion
  case UserListCriterion:
    cc.Criterion = criterion
  case VerticalCriterion:
  }
}
*/

func (s *CampaignCriterionService) Get(selector Selector) (campaignCriterions CampaignCriterions, totalCount int64, err error) {
	selector.XMLName = xml.Name{baseUrl, "serviceSelector"}
	respBody, err := s.Auth.request(
		campaignCriterionServiceUrl,
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
	if err != nil {
		return campaignCriterions, totalCount, err
	}
	getResp := struct {
		Size               int64              `xml:"rval>totalNumEntries"`
		CampaignCriterions CampaignCriterions `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return campaignCriterions, totalCount, err
	}
	return getResp.CampaignCriterions, getResp.Size, err
}

func (s *CampaignCriterionService) Mutate(campaignCriterionOperations CampaignCriterionOperations) (campaignCriterions CampaignCriterions, err error) {
	type campaignCriterionOperation struct {
		Action            string      `xml:"operator"`
		CampaignCriterion interface{} `xml:"operand"`
	}
	operations := []campaignCriterionOperation{}
	for action, campaignCriterions := range campaignCriterionOperations {
		for _, campaignCriterion := range campaignCriterions {
			operations = append(operations,
				campaignCriterionOperation{
					Action:            action,
					CampaignCriterion: campaignCriterion,
				},
			)
		}
	}
	mutation := struct {
		XMLName xml.Name
		Ops     []campaignCriterionOperation `xml:"operations"`
	}{
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "mutate",
		},
		Ops: operations,
	}
	respBody, err := s.Auth.request(campaignCriterionServiceUrl, "mutate", mutation)
	if err != nil {
		/*
			    switch t := err.(type) {
			    case *ErrorsType:
				    for action, campaignCriterions := range campaignCriterionOperations {
					    for _, campaignCriterion := range campaignCriterions {
			          campaignCriterions = append(campaignCriterions,campaignCriterion)
			        }
			      }
			      for _, aef := range t.ApiExceptionFaults {
			        for _,e := range aef.Errors {
			          switch et := e.(type) {
			          case CriterionError:
			            offset, err := strconv.ParseInt(strings.Trim(et.FieldPath,"abcdefghijklmnop.]["),10,64)
			            if err != nil {
			              return CampaignCriterions{}, err
			            }
			            cc := campaignCriterions[offset]
			            switch c := cc.(type) {
			            case CampaignCriterion:
			              CampaignCriterion(campaignCriterions[offset]).Errors = append(campaignCriterions[offset].(CampaignCriterion).Errors,fmt.Errorf(et.Reason))
			            case NegativeCampaignCriterion:
			              NegativeCampaignCriterion(campaignCriterions[offset]).Errors = append(NegativeCampaignCriterion(campaignCriterions[offset].Errors),fmt.Errorf(et.Reason))
			            }
			          }
			        }
			      }
			    default:
		*/
		return campaignCriterions, err
		//}
	}
	mutateResp := struct {
		CampaignCriterions CampaignCriterions `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return campaignCriterions, err
	}
	return mutateResp.CampaignCriterions, err
}

func (s *CampaignCriterionService) Query(query string) (campaignCriterions CampaignCriterions, totalCount int64, err error) {
	respBody, err := s.Auth.request(
		campaignCriterionServiceUrl,
		"query",
		AWQLQuery{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "query",
			},
			Query: query,
		},
	)

	if err != nil {
		return campaignCriterions, totalCount, err
	}

	getResp := struct {
		Size               int64              `xml:"rval>totalNumEntries"`
		CampaignCriterions CampaignCriterions `xml:"rval>entries"`
	}{}

	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return campaignCriterions, totalCount, err
	}
	return getResp.CampaignCriterions, getResp.Size, err
}

type CampaignScheduleCriterion struct {
	CampaignId  int64               `xml:"campaignId"`
	IsNegative  bool                `xml:"isNegative,omitempty"`
	Criterion   AdScheduleCriterion `xml:"criterion"`
	BidModifier *float64            `xml:"bidModifier,omitempty"`
	Errors      []error             `xml:"-"`
}

func (s *CampaignCriterionService) GetScheduleCriterions(selector Selector) (campaignCriterions []CampaignScheduleCriterion, totalCount int64, err error) {
	selector.XMLName = xml.Name{baseUrl, "serviceSelector"}
	respBody, err := s.Auth.request(
		campaignCriterionServiceUrl,
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
	if err != nil {
		return campaignCriterions, totalCount, err
	}
	getResp := struct {
		Size               int64                       `xml:"rval>totalNumEntries"`
		CampaignCriterions []CampaignScheduleCriterion `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return campaignCriterions, totalCount, err
	}
	return getResp.CampaignCriterions, getResp.Size, err
}

func (s *CampaignCriterionService) SetSchedule(campaignid int64, start, end int) error {
	//先拿取以前的  先删除再添加
	campaignCriterions, _, err := s.Get(
		Selector{
			Fields: []string{
				"CampaignId",
				"CriteriaType",
				"DayOfWeek",
				"StartHour",
				"EndHour",
			},
			Predicates: []Predicate{{"CampaignId", "EQUALS", []string{fmt.Sprintf("%v", campaignid)}},
				{"CriteriaType", "EQUALS", []string{"AD_SCHEDULE"}},
			},
		},
	)
	if err != nil {
		return errors.WithMessage(err, "GetScheduleCriterions ")
	}

	//删除
	if len(campaignCriterions) > 0 {
		_, err = s.Mutate(
			CampaignCriterionOperations{
				"REMOVE": campaignCriterions,
			},
		)
		if err != nil {
			return errors.WithMessage(err, "Remove ScheduleCriterions ")
		}

	}
	//添加
	_, err = s.Mutate(
		CampaignCriterionOperations{
			"ADD": GenScheduleCriterions(campaignid, start, end),
		},
	)
	return err
}

func GenScheduleCriterions(campaignid int64, start, end int) CampaignCriterions {
	if start > end {
		return CampaignCriterions{
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "MONDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", 24), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "TUESDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", 24), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "WEDNESDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", 24), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "THURSDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", 24), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "FRIDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", 24), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "SATURDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", 24), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "SUNDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", 24), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "MONDAY", StartHour: fmt.Sprintf("%v", 0), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "TUESDAY", StartHour: fmt.Sprintf("%v", 0), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "WEDNESDAY", StartHour: fmt.Sprintf("%v", 0), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "THURSDAY", StartHour: fmt.Sprintf("%v", 0), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "FRIDAY", StartHour: fmt.Sprintf("%v", 0), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "SATURDAY", StartHour: fmt.Sprintf("%v", 0), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "SUNDAY", StartHour: fmt.Sprintf("%v", 0), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
		}
	} else {
		return CampaignCriterions{
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "MONDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "TUESDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "WEDNESDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "THURSDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "FRIDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "SATURDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
			CampaignCriterion{CampaignId: campaignid, Criterion: AdScheduleCriterion{DayOfWeek: "SUNDAY", StartHour: fmt.Sprintf("%v", start), StartMinute: "ZERO", EndHour: fmt.Sprintf("%v", end), EndMinute: "ZERO"}},
		}
	}

}
