package base58

import (
	"crypto/sha256"
	"fmt"
	"github.com/itchyny/base58-go"
	"math/big"
	"os"
)

func GenerateShortLink(initialLink string, userId string) string {
	urlHashBytes := sha256Of(fmt.Sprintf("%s%s", initialLink, userId))
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:6]
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}
