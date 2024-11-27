package shortener

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput = errors.New("invalid input for encoding")
	charset         = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base            = int64(len(charset))
)

type Base62Encoder struct{}

func NewBase62Encoder() *Base62Encoder {
	return &Base62Encoder{}
}

func (b *Base62Encoder) Encode(num int64) string {
	if num == 0 {
		return string(charset[0])
	}

	encoded := make([]byte, 0)
	for num > 0 {
		encoded = append(encoded, charset[num%base])
		num = num / base
	}

	// Reverse the slice
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}

	return string(encoded)
}

func (b *Base62Encoder) Decode(encoded string) (int64, error) {
	var num int64
	for i := 0; i < len(encoded); i++ {
		char := encoded[i]
		pos := int64(-1)

		for j, c := range charset {
			if char == byte(c) {
				pos = int64(j)
				break
			}
		}

		if pos == -1 {
			return 0, fmt.Errorf("%w: invalid character %c", ErrInvalidInput, char)
		}

		num = num*base + pos
	}
	return num, nil
}
