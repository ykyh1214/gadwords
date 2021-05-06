package gadwords

import "encoding/xml"

type ConversionTrackingSettings struct {
	EffectiveConversionTrackingId      int64 `xml:effectiveConversionTrackingId`
	UsesCrossAccountConversionTracking bool  `xml:usesCrossAccountConversionTracking`
}
type ConversionTracker struct {
	Id       int64  `xml:"id,omitempty"`
	Name     string `xml:"name,omitempty"`
	Category string `xml:"category,omitempty"`
	Status   string `xml:"status,omitempty"`
}
type ConversionTrackerService struct {
	Auth
}

func NewConversionTrackerService(auth *Auth) *ConversionTrackerService {
	return &ConversionTrackerService{Auth: *auth}
}

func (s *ConversionTrackerService) Get(selector Selector) (ConversionTrackers []ConversionTracker, totalCount int64, err error) {
	// The default namespace, "", will break in 1.5 with the addition of
	// custom namespace support.  Hence, we have to ensure that the baseUrl is
	// set again as the proper namespace for the service/serviceSelector element
	selector.XMLName = xml.Name{baseUrl, "serviceSelector"}

	respBody, err := s.Auth.request(
		conversionTrackerServiceUrl,
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
		return ConversionTrackers, totalCount, err
	}
	getResp := struct {
		Size               int64               `xml:"rval>totalNumEntries"`
		ConversionTrackers []ConversionTracker `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return ConversionTrackers, totalCount, err
	}
	return getResp.ConversionTrackers, getResp.Size, err

}
