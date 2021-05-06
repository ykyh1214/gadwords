package gadwords

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

// A campaignService holds the connection information for the
// campaign service.
type GAService struct {
	Auth
}

// NewCampaignService creates a new campaignService
func NewGAService(auth *Auth) *GAService {
	return &GAService{Auth: *auth}
}

type GAContext struct {
	ViewId     string
	StartDay   string
	EndDay     string
	Dimensions []string
	Metrics    []string
	Actid      string
}

//测试 //以mcc为粒度
func (s *GAService) Query(ctx GAContext) ([]byte, error) {
	p := url.Values{}
	p.Set("ids", fmt.Sprintf("ga:%v", ctx.ViewId))
	p.Set("start-date", ctx.StartDay)
	p.Set("end-date", ctx.EndDay)
	p.Set("dimensions", strings.Join(ctx.Dimensions, ","))
	p.Set("metrics", strings.Join(ctx.Metrics, ","))
	p.Set("filters", fmt.Sprintf("ga:adwordsCustomerID==%v;ga:adCost>0", ctx.Actid))
	//p.Set("access_token", token)
	req, err := http.NewRequest("", "https://www.googleapis.com/analytics/v3/data/ga?"+p.Encode(), strings.NewReader(""))
	res, err := s.Auth.Client.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	if gjson.GetBytes(data, "error").Exists() {
		return nil, errors.New(string(data))
	}
	return data, nil
}
