package gadwords

import (
	"encoding/base64"
	"encoding/xml"
)

type AssetsService struct {
	Auth
}

func NewAssetsService(auth *Auth) *AssetsService {
	return &AssetsService{Auth: *auth}
}

// Media represents an audio, image or video file.
type Asset struct {
	Type           string        `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	Id             int64         `xml:"assetId,omitempty"`
	AssetName      string        `xml:"assetName,omitempty"`
	ImageData      string        `xml:"imageData,omitempty"`
	AssetText      string        `xml:"assetText,omitempty"`
	YouTubeVideoId string        `xml:"youTubeVideoId,omitempty"`
	FullSizeInfo   *FullSizeInfo `xml:"fullSizeInfo,omitempty"`
}
type FullSizeInfo struct {
	ImageHeight int64  `xml:"imageHeight,omitempty"`
	ImageWidth  int64  `xml:"imageWidth,omitempty"`
	ImageUrl    string `xml:"imageUrl,omitempty"`
}

func NewTextAsset(text string) Asset {
	return Asset{
		Type:      "TextAsset",
		AssetText: text,
	}
}

func NewImageAsset(name string, data []byte) Asset {
	imageData := base64.StdEncoding.EncodeToString(data)
	return Asset{
		Type:      "ImageAsset",
		AssetName: name,
		ImageData: imageData,
	}
}

func NewYouTubeVideoAsset(name string, youTubeVideoId string) Asset {
	return Asset{
		Type:           "YouTubeVideoAsset",
		AssetName:      name,
		YouTubeVideoId: youTubeVideoId,
	}
}

func (s *AssetsService) Get(selector Selector) (assets []Asset, totalCount int64, err error) {
	selector.XMLName = xml.Name{baseUrl, "selector"}
	respBody, err := s.Auth.request(
		assetServiceUrl,
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
		return assets, totalCount, err
	}
	getResp := struct {
		Size   int64   `xml:"rval>totalNumEntries"`
		Assets []Asset `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return assets, totalCount, err
	}
	return getResp.Assets, getResp.Size, err
}

type AssetOperations map[string][]Asset

func (s *AssetsService) Mutate(assetOperations AssetOperations) (assets []Asset, err error) {
	type assetOperation struct {
		Action string  `xml:"operator"`
		Assets []Asset `xml:"operand"`
	}
	operations := []assetOperation{}
	for action, assets := range assetOperations {
		for _, asset := range assets {
			ad := []Asset{asset}
			operations = append(operations,
				assetOperation{
					Action: action,
					Assets: ad,
				},
			)
		}
	}
	mutation := struct {
		XMLName xml.Name
		Ops     []assetOperation `xml:"operations"`
	}{
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "mutate",
		},
		Ops: operations,
	}

	respBody, err := s.Auth.request(assetServiceUrl, "mutate", mutation)
	if err != nil {
		return assets, err
	}
	mutateResp := struct {
		Assets []Asset `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return assets, err
	}
	return mutateResp.Assets, err
}
