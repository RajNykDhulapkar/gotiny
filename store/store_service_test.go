package store

import (
	"context"
	"testing"

	"github.com/RajNykDhulapkar/gotiny/interfaces"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

var testStoreService interfaces.StorageServiceInterface

func init() {
	mockRedis, err := miniredis.Run()
	if err != nil {
		panic("could not start miniredis")
	}

	testStoreService, err = InitializeStore(mockRedis.Addr(), "", 0)
	if err != nil {
		panic("could not initialize store")
	}
}

func TestInsertionAndRetrieval(t *testing.T) {
	initialLink := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	userUUId := "e0dba740-fc4b-4977-872c-d360239e6b1a"
	shortURL := "Jsz4k57oAX"

	err := testStoreService.SaveUrlMapping(context.Background(), shortURL, initialLink, userUUId)
	assert.Nil(t, err)

	retrievedUrl, err := testStoreService.RetrieveOriginalUrl(context.Background(), shortURL)
	assert.Nil(t, err)
	assert.Equal(t, initialLink, retrievedUrl)
}
