package utils

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/spf13/cast"
)

func SimpleScientificBigIntParse(val string) (*big.Int, error) {
	// Only allowed values in *simple* parsing are [0-9], 'e/E',
	// optional '-/+' in the start, and a optional '+' after 'e/E'.
	// No '-' allowed after 'e' since we are parsing Ints not Floats.
	badFormatErr := fmt.Errorf("error parsing BigInt value '%s' : bad format", val)
	if matched, err := regexp.MatchString(`^[-+]?[0-9]+([e][+]?[0-9]+)?$`, val); err != nil {
		return nil, fmt.Errorf("error parsing BigInt value '%s' : %s", val, err)
	} else if !matched {
		return nil, badFormatErr
	}
	indexOfE := strings.Index(val, "e")
	var firstPart, secondPart string
	if indexOfE == -1 { // no 'e' present
		firstPart = val
		secondPart = "0"
	} else {
		firstPart = val[0:indexOfE]
		secondPart = val[indexOfE+1:] // regular exp above ensures there are chars after 'e'
	}
	base, success1 := big.NewInt(0).SetString(cast.ToString(firstPart), 0)
	exp10, success2 := big.NewInt(0).SetString(cast.ToString(secondPart), 0)
	if !success1 || !success2 {
		return nil, badFormatErr
	}
	// return firstPart * 10^secondPart
	return base.Mul(base, big.NewInt(0).Exp(big.NewInt(10), exp10, nil)), nil
}
