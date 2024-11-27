package shortener

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/RajNykDhulapkar/gotiny/pkg/interfaces"
)

type Shortener struct {
	encoder        interfaces.Base62EncoderPort
	rangeAllocator interfaces.RangeAllocatorPort // Interface from your range allocator module
}

type Config struct {
	ServiceID string
	RangeSize int64
}

func NewShortener(rangeAllocator interfaces.RangeAllocatorPort, config *Config) *Shortener {
	return &Shortener{
		encoder:        NewBase62Encoder(),
		rangeAllocator: rangeAllocator,
	}
}

func (s *Shortener) GenerateShortLink(ctx context.Context, originalURL string, userID string) (string, error) {
	// Get next available ID from range allocator
	id, err := s.rangeAllocator.GetNextID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get next ID: %w", err)
	}

	// Calculate checksum from original URL & userID
	h := crc32.NewIEEE()
	h.Write([]byte(originalURL))
	h.Write([]byte(userID))
	checksum := int64(h.Sum32()) % 62 // Keep within base62 single char

	// Combine: id + encode(checksum)
	return s.encoder.Encode(id) + s.encoder.Encode(checksum), nil
}

// Base62 test function to demonstrate encoding/decoding
func (s *Shortener) TestBase62() {
	testCases := []int64{0, 1, 61, 62, 1000, 999999}

	for _, num := range testCases {
		encoded := s.encoder.Encode(num)
		decoded, err := s.encoder.Decode(encoded)

		fmt.Printf("Number: %d\n", num)
		fmt.Printf("Encoded: %s\n", encoded)
		if err != nil {
			fmt.Printf("Decode error: %v\n", err)
		} else {
			fmt.Printf("Decoded: %d\n", decoded)
		}
		fmt.Println("---")
	}
}
