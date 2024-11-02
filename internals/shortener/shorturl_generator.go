package shortener

import (
	"fmt"
	"math/big"

	"github.com/RajNykDhulapkar/gotiny/internals/utils"
	"github.com/RajNykDhulapkar/gotiny/pkg/interfaces"
)

type Shortener struct {
	Base58Encoder interfaces.Base58EncoderInterface
}

func NewShortener(encoder interfaces.Base58EncoderInterface) interfaces.ShortenerInterface {
	return &Shortener{Base58Encoder: encoder}
}

func (s *Shortener) GenerateShortLink(initialLink string, userId string) (string, error) {
	urlHashBytes := utils.Sha256Of(initialLink + userId)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, err := s.Base58Encoder.Encode([]byte(fmt.Sprintf("%d", generatedNumber)))
	if err != nil {
		return "", fmt.Errorf("error encoding the URL: %v", err)
	}
	return finalString[:8], nil
}
