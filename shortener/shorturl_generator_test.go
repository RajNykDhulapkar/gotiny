package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortLink(t *testing.T) {
	encoder := NewBase58Encoder()
	shortener := NewShortener(encoder)

	initialLink := "https://www.example.com"
	userId := "user123"

	shortLink, err := shortener.GenerateShortLink(initialLink, userId)

	assert.NoError(t, err)
	assert.NotEmpty(t, shortLink)
	assert.Len(t, shortLink, 8)
}

func TestGenerateShortLinkDifferentUsers(t *testing.T) {
	encoder := NewBase58Encoder()
	shortener := NewShortener(encoder)

	initialLink := "https://www.example.com"
	userId1 := "user123"
	userId2 := "user456"

	shortLink1, err := shortener.GenerateShortLink(initialLink, userId1)
	assert.NoError(t, err)
	assert.NotEmpty(t, shortLink1)

	shortLink2, err := shortener.GenerateShortLink(initialLink, userId2)
	assert.NoError(t, err)
	assert.NotEmpty(t, shortLink2)

	assert.NotEqual(t, shortLink1, shortLink2)
}

func TestGenerateShortLinkSameUserDifferentLinks(t *testing.T) {
	encoder := NewBase58Encoder()
	shortener := NewShortener(encoder)

	initialLink1 := "https://www.example.com"
	initialLink2 := "https://www.example.org"
	userId := "user123"

	shortLink1, err := shortener.GenerateShortLink(initialLink1, userId)
	assert.NoError(t, err)
	assert.NotEmpty(t, shortLink1)

	shortLink2, err := shortener.GenerateShortLink(initialLink2, userId)
	assert.NoError(t, err)
	assert.NotEmpty(t, shortLink2)

	assert.NotEqual(t, shortLink1, shortLink2)
}
