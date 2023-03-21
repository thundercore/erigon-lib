package hardfork

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/spf13/cast"
)

type HardforkConfig struct {
	name       string
	prettyName string
	desc       string
}

func newHardforkConfig(name, desc string) *HardforkConfig {
	return &HardforkConfig{
		name:       strings.ToLower(name),
		prettyName: name,
		desc:       desc,
	}
}

type Int64HardforkConfig struct {
	HardforkConfig
}

func NewInt64HardforkConfig(name, desc string) *Int64HardforkConfig {
	return &Int64HardforkConfig{
		HardforkConfig: *newHardforkConfig(name, desc),
	}
}

func (c *Int64HardforkConfig) GetValueHardforkAtBlock(hardforks *Hardforks, blockNum int64) int64 {
	if value, found := hardforks.GetValueHardforkAtBlock(c.name, blockNum); found {
		return tryConvertToInt64(value)
	}
	return -1
}

func (c *Int64HardforkConfig) GetValueHardforkAtSession(hardforks *Hardforks, sessionNum int64) int64 {
	if value, found := hardforks.GetValueHardforkAtSession(c.name, sessionNum); found {
		return tryConvertToInt64(value)
	}
	return -1
}

type Float64HardforkConfig struct {
	HardforkConfig
}

func NewFloat64HardforkConfig(name, desc string) *Float64HardforkConfig {
	return &Float64HardforkConfig{
		HardforkConfig: *newHardforkConfig(name, desc),
	}
}

func (c *Float64HardforkConfig) GetValueHardforkAtBlock(hardforks *Hardforks, blockNum int64) float64 {
	if value, found := hardforks.GetValueHardforkAtBlock(c.name, blockNum); found {
		return value.(float64)
	}
	return -1
}

func (c *Float64HardforkConfig) GetValueHardforkAtSession(hardforks *Hardforks, sessionNum int64) float64 {
	if value, found := hardforks.GetValueHardforkAtSession(c.name, sessionNum); found {
		return value.(float64)
	}
	return -1
}

type BoolHardforkConfig struct {
	HardforkConfig
}

func NewBoolHardforkConfig(name, desc string) *BoolHardforkConfig {
	return &BoolHardforkConfig{
		HardforkConfig: *newHardforkConfig(name, desc),
	}
}

func (c *BoolHardforkConfig) GetValueHardforkAtBlock(hardforks *Hardforks, blockNum int64) bool {
	if value, found := hardforks.GetValueHardforkAtBlock(c.name, blockNum); found {
		return value.(bool)
	}
	return false
}

func (c *BoolHardforkConfig) GetValueHardforkAtSession(hardforks *Hardforks, sessionNum int64) bool {
	if value, found := hardforks.GetValueHardforkAtSession(c.name, sessionNum); found {
		return value.(bool)
	}
	return false
}

type StringHardforkConfig struct {
	HardforkConfig
}

func NewStringHardforkConfig(name, desc string) *StringHardforkConfig {
	return &StringHardforkConfig{
		HardforkConfig: *newHardforkConfig(name, desc),
	}
}

func (c *StringHardforkConfig) GetValueHardforkAtBlock(hardforks *Hardforks, blockNum int64) string {
	if value, found := hardforks.GetValueHardforkAtBlock(c.name, blockNum); found {
		return value.(string)
	}
	return ""
}

func (c *StringHardforkConfig) GetValueHardforkAtSession(hardforks *Hardforks, sessionNum int64) string {
	if value, found := hardforks.GetValueHardforkAtSession(c.name, sessionNum); found {
		return value.(string)
	}
	return ""
}

type BigIntHardforkConfig struct {
	HardforkConfig
}

func NewBigIntHardforkConfig(name, desc string) *BigIntHardforkConfig {
	return &BigIntHardforkConfig{
		HardforkConfig: *newHardforkConfig(name, desc),
	}
}

func (c *BigIntHardforkConfig) GetValueHardforkAtBlock(hardforks *Hardforks, blockNum int64) *big.Int {
	if value, found := hardforks.GetValueHardforkAtBlock(c.name, blockNum); found {
		bi, err := SimpleScientificBigIntParse(cast.ToString(value))
		if err != nil {
			panic(err)
		}
		return bi
	}
	bi, _ := SimpleScientificBigIntParse("-1")
	return bi
}

func (c *BigIntHardforkConfig) GetValueHardforkAtSession(hardforks *Hardforks, sessionNum int64) *big.Int {
	if value, found := hardforks.GetValueHardforkAtSession(c.name, sessionNum); found {
		bi, err := SimpleScientificBigIntParse(cast.ToString(value))
		if err != nil {
			panic(err)
		}
		return bi
	}
	bi, _ := SimpleScientificBigIntParse("-1")
	return bi
}

type AddressHardforkConfig struct {
	HardforkConfig
}

func NewAddressHardforkConfig(name, desc string) *AddressHardforkConfig {
	return &AddressHardforkConfig{
		HardforkConfig: *newHardforkConfig(name, desc),
	}
}

func (c *AddressHardforkConfig) GetValueHardforkAtBlock(hardforks *Hardforks, blockNum int64) common.Address {
	if value, found := hardforks.GetValueHardforkAtBlock(c.name, blockNum); found {
		return common.HexToAddress(cast.ToString(value))
	}
	return common.HexToAddress("0x0")
}

func (c *AddressHardforkConfig) GetValueHardforkAtSession(hardforks *Hardforks, sessionNum int64) common.Address {
	if value, found := hardforks.GetValueHardforkAtSession(c.name, sessionNum); found {
		return common.HexToAddress(cast.ToString(value))
	}
	return common.HexToAddress("0x0")
}

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
