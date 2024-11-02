package shortener

import (
	"github.com/RajNykDhulapkar/gotiny/pkg/interfaces"
	"github.com/itchyny/base58-go"
)

type Base58Encoder struct{}

func NewBase58Encoder() interfaces.Base58EncoderInterface {
	return &Base58Encoder{}
}

func (b *Base58Encoder) Encode(input []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(input)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
