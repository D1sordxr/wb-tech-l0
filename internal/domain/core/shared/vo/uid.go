package vo

import (
	"crypto/rand"
	"fmt"
	"math/big"

	sharedErrs "github.com/D1sordxr/wb-tech-l0/internal/domain/core/shared/errors"
)

const (
	lengthUID        = 20
	suffixForUID     = "test"
	randomPartLength = lengthUID - len(suffixForUID)
)

func GenerateUID() (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

	randomPart := make([]byte, randomPartLength)

	for i := range randomPart {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("generate random number: %w", err)
		}
		randomPart[i] = chars[num.Int64()]
	}

	return string(randomPart) + "test", nil
}

func ValidateUID(value string) error {
	if len(value) != lengthUID {
		return sharedErrs.ErrOrderUIDInvalidLength
	}

	if value[randomPartLength:] != suffixForUID {
		return sharedErrs.ErrOrderUIDInvalidSuffix
	}

	for i := 0; i < randomPartLength; i++ {
		c := value[i]
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return sharedErrs.ErrOrderUIDInvalidChars
		}
	}

	return nil
}
