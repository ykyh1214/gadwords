package gadwords

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func testAssetsService(t *testing.T) (service *AssetsService) {
	return &AssetsService{Auth: testAuthSetup(t)}
}

func TestAssets(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	// load image into []byte

	ms := testAssetsService(t)

	var pageSize int64 = 500
	var offset int64 = 0
	paging := Paging{
		Offset: offset,
		Limit:  pageSize,
	}
	for {
		assets, totalCount, err := ms.Get(
			Selector{
				Fields: []string{
					"AssetId",
					"AssetName",
					"AssetSubtype",
					//"MimeType",
					"AssetText",
					"ImageFullSizeUrl",
					"YouTubeVideoId",
					//"AssetText",
				},
				/* 	Predicates: []Predicate{
					{"Type", "IN", []string{"IMAGE", "VIDEO"}},
				}, */
				Paging: &paging,
			},
		)
		if err != nil {
			fmt.Printf("Error occured finding medias")
		}
		log.Println(err, totalCount, assets)
		// Increment values to request the next page.
		offset += pageSize
		paging.Offset = offset
		if totalCount < offset {
			fmt.Printf("\tFound %d entries.", totalCount)
			break
		}
	}

}

func TestAssetsPut(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	// load image into []byte

	body, err := ioutil.ReadFile("cut_send.png")
	log.Println(err)
	ms := testAssetsService(t)
	assets, err := ms.Mutate(
		AssetOperations{
			"ADD": {
				NewImageAsset("", body),
				//NewImageAsset("yyk3", body),
			},
		},
	)
	log.Println(assets, err)
}
