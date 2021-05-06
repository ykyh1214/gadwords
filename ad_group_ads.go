package gadwords

import (
	"encoding/xml"
	"fmt"
	"log"
)

//type AppAds []interface{}

type AdGroupAds []interface{}

func (a1 AdGroupAds) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	a := a1[0]
	if _, ok := a.(OptAd); ok {
		e.EncodeElement(a1[0], start)
		return nil
	}
	e.EncodeToken(start)

	switch a.(type) {
	case TextAd:
		ad := a.(TextAd)
		e.EncodeElement(ad.AdGroupId, xml.StartElement{Name: xml.Name{"", "adGroupId"}})
		e.EncodeElement(ad, xml.StartElement{
			xml.Name{"", "ad"},
			[]xml.Attr{
				xml.Attr{xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"}, "TextAd"},
			},
		})
		e.EncodeElement(ad.Status, xml.StartElement{Name: xml.Name{"", "status"}})
		e.EncodeElement(ad.Labels, xml.StartElement{Name: xml.Name{"", "labels"}})
	case ExpandedTextAd:
		ad := a.(ExpandedTextAd)
		e.EncodeElement(ad.AdGroupId, xml.StartElement{Name: xml.Name{"", "adGroupId"}})
		e.EncodeElement(ad, xml.StartElement{
			xml.Name{"", "ad"},
			[]xml.Attr{
				xml.Attr{xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"}, "ExpandedTextAd"},
			},
		})
		e.EncodeElement(ad.Status, xml.StartElement{Name: xml.Name{"", "status"}})
		e.EncodeElement(ad.Labels, xml.StartElement{Name: xml.Name{"", "labels"}})
	case AppAd:
		ad := a.(AppAd)
		e.EncodeElement(ad.AdGroupId, xml.StartElement{Name: xml.Name{"", "adGroupId"}})
		e.EncodeElement(ad, xml.StartElement{
			xml.Name{"", "ad"},
			[]xml.Attr{
				xml.Attr{xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"}, "UniversalAppAd"},
			},
		})
		e.EncodeElement(ad.Labels, xml.StartElement{Name: xml.Name{"", "labels"}})
	case ResponsiveSearchAd:
		ad := a.(ResponsiveSearchAd)
		e.EncodeElement(ad.AdGroupId, xml.StartElement{Name: xml.Name{"", "adGroupId"}})
		e.EncodeElement(ad, xml.StartElement{
			xml.Name{"", "ad"},
			[]xml.Attr{
				xml.Attr{xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"}, "ResponsiveSearchAd"},
			},
		})
		e.EncodeElement(ad.Status, xml.StartElement{Name: xml.Name{"", "status"}})
		e.EncodeElement(ad.Labels, xml.StartElement{Name: xml.Name{"", "labels"}})
		/* 	case Ad:
		ad := a.(Ad)
		e.EncodeElement(ad.AdGroupId, xml.StartElement{Name: xml.Name{"", "adGroupId"}})
		e.EncodeElement(ad, xml.StartElement{
			xml.Name{"", "ad"},
			[]xml.Attr{
				xml.Attr{xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"}, "Ad"},
			},
		})
		e.EncodeElement(ad.Status, xml.StartElement{Name: xml.Name{"", "status"}}) */
	case ImageAd:
		return ERROR_NOT_YET_IMPLEMENTED
	case TemplateAd:
		return ERROR_NOT_YET_IMPLEMENTED
	default:
		return fmt.Errorf("unknown Ad type -> %#v", start)
	}
	e.EncodeToken(start.End())
	return nil
}

func (aga *AdGroupAds) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {

	typeName := xml.Name{Space: "http://www.w3.org/2001/XMLSchema-instance", Local: "type"}
	var adGroupId int64
	var status, approvalStatus string
	var disapprovalReasons []string
	var trademarkDisapproved bool
	var labels []Label
	var experimentData *AdGroupExperimentData
	var trademarks []string
	var baseCampaignId *int64
	var baseAdGroupId *int64
	var ad interface{}
	var headlines, descriptions, images, videos []AssetLink
	log.Println("yyk")
	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		if err != nil {
			return err
		}
		switch start := token.(type) {
		case xml.StartElement:
			tag := start.Name.Local
			switch tag {
			case "adGroupId":
				err := dec.DecodeElement(&adGroupId, &start)
				if err != nil {
					return err
				}
			case "ad":
				typeName, err := findAttr(start.Attr, typeName)
				if err != nil {
					return err
				}
				fmt.Println(typeName)
				switch typeName {
				case "TextAd":
					a := TextAd{AdGroupId: adGroupId}
					err := dec.DecodeElement(&a, &start)
					if err != nil {
						return err
					}
					ad = a
				case "ExpandedTextAd":
					a := ExpandedTextAd{AdGroupId: adGroupId}
					err := dec.DecodeElement(&a, &start)
					if err != nil {
						return err
					}
					ad = a
				case "ImageAd":
					a := ImageAd{AdGroupId: adGroupId}
					err := dec.DecodeElement(&a, &start)
					if err != nil {
						return err
					}
					ad = a
				case "TemplateAd":
					a := TemplateAd{AdGroupId: adGroupId}
					err := dec.DecodeElement(&a, &start)
					if err != nil {
						return err
					}
					ad = a
				case "UniversalAppAd":
					a := AppAd{AdGroupId: adGroupId}
					err := dec.DecodeElement(&a, &start)
					if err != nil {
						return err
					}
					ad = a
					log.Println(a.Headlines)
				case "DynamicSearchAd":
					a := DynamicSearchAd{AdGroupId: adGroupId}
					err := dec.DecodeElement(&a, &start)
					if err != nil {
						return err
					}
					ad = a
				case "ProductAd":
					a := ProductAd{AdGroupId: adGroupId}
					err := dec.DecodeElement(&a, &start)
					if err != nil {
						return err
					}
					ad = a
				case "ResponsiveSearchAd":
					a := ResponsiveSearchAd{AdGroupId: adGroupId}
					err := dec.DecodeElement(&a, &start)
					if err != nil {
						return err
					}
					ad = a
				default:
					return fmt.Errorf("unknown AdGroupCriterion -> %#v", start)
				}
			case "experimentData":
				err := dec.DecodeElement(&experimentData, &start)
				if err != nil {
					return err
				}
			case "status":
				err := dec.DecodeElement(&status, &start)
				if err != nil {
					return err
				}
			case "approvalStatus":
				err := dec.DecodeElement(&approvalStatus, &start)
				if err != nil {
					return err
				}
			case "trademarks":
				err := dec.DecodeElement(&trademarks, &start)
				if err != nil {
					return err
				}
			case "disapprovalReasons":
				err := dec.DecodeElement(&disapprovalReasons, &start)
				if err != nil {
					return err
				}
			case "trademarkDisapproved":
				err := dec.DecodeElement(&trademarkDisapproved, &start)
				if err != nil {
					return err
				}
			case "labels":
				err := dec.DecodeElement(&labels, &start)
				if err != nil {
					return err
				}
			case "baseCampaignId":
				err := dec.DecodeElement(&baseCampaignId, &start)
				if err != nil {
					return err
				}
			case "baseAdGroupId":
				err := dec.DecodeElement(&baseAdGroupId, &start)
				if err != nil {
					return err
				}
			case "headlines":
				err := dec.DecodeElement(&headlines, &start)
				if err != nil {
					return err
				}
				log.Println(headlines)
			case "descriptions":
				err := dec.DecodeElement(&descriptions, &start)
				if err != nil {
					return err
				}
			case "images":
				err := dec.DecodeElement(&images, &start)
				if err != nil {
					return err
				}
			case "videos":
				err := dec.DecodeElement(&videos, &start)
				if err != nil {
					return err
				}
			default:
				//return fmt.Errorf("unknown AdGroupAd field -> %#v", tag)
			}

		}
	}
	switch a := ad.(type) {
	case TextAd:
		a.Status = status
		*aga = append(*aga, a)
	case ExpandedTextAd:
		a.ExperimentData = experimentData
		a.Status = status
		a.Labels = labels
		a.BaseCampaignId = baseCampaignId
		a.BaseAdGroupId = baseAdGroupId
		*aga = append(*aga, a)
	case ImageAd:
		a.Status = status
		*aga = append(*aga, a)
	case TemplateAd:
		a.Status = status
		*aga = append(*aga, a)
	case DynamicSearchAd:
		a.Status = status
		*aga = append(*aga, a)
	case ProductAd:
		a.Status = status
		*aga = append(*aga, a)
	case AppAd:
		a.Labels = labels
		*aga = append(*aga, a)
	case ResponsiveSearchAd:
		a.Status = status
		*aga = append(*aga, a)
	}

	return nil
}
