package hardfork

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Hardforks struct {
	BlockHardforks   map[string]BlockHardforks
	SessionHardforks map[string]SessionHardforks
}

func (h *Hardforks) GetValueHardforkAtBlock(name string, blockNum int64) (interface{}, bool) {
	length := len(h.BlockHardforks[name])
	if length == 0 {
		return nil, false
	}

	lower := sort.Search(length, func(i int) bool {
		return h.BlockHardforks[name][i].BlockNum >= blockNum
	})

	return h.BlockHardforks[name][lower-1].Value, true
}

func (h *Hardforks) GetValueHardforkAtSession(name string, sessionNum int64) (interface{}, bool) {
	length := len(h.SessionHardforks[name])
	if length == 0 {
		return nil, false
	}

	lower := sort.Search(length, func(i int) bool {
		return h.SessionHardforks[name][i].SessionNum >= sessionNum
	})

	return h.SessionHardforks[name][lower-1].Value, true
}

type BlockHardfork struct {
	BlockNum int64
	Value    interface{}
}

type BlockHardforks []BlockHardfork

func (h BlockHardforks) Len() int {
	return len(h)
}

func (h BlockHardforks) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h BlockHardforks) Less(i, j int) bool {
	return h[i].BlockNum < h[j].BlockNum
}

type SessionHardfork struct {
	SessionNum int64
	Value      interface{}
}

type SessionHardforks []SessionHardfork

func (s SessionHardforks) Len() int {
	return len(s)
}

func (s SessionHardforks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SessionHardforks) Less(i, j int) bool {
	return s[i].SessionNum < s[j].SessionNum
}

func NewHardforks(configPath string) *Hardforks {
	hardforks := &Hardforks{
		BlockHardforks:   make(map[string]BlockHardforks),
		SessionHardforks: make(map[string]SessionHardforks),
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		e := fmt.Sprintf("failed to read hardfork config file: %v", err)
		panic(e)
	}

	rawConfigs := []map[string]interface{}{}
	if err := yaml.Unmarshal(content, &rawConfigs); err != nil {
		panic(err)
	}

	for index, config := range rawConfigs {
		block, hasNumber := config["blocknum"]
		session, hasSession := config["session"]
		if !hasNumber && !hasSession {
			panic(fmt.Sprintf("Hardfork config %d must have blocknum or sessionnum", index))
		}

		delete(config, "blocknum")
		delete(config, "desc")
		delete(config, "session")

		viperCfg := viper.New()
		viperCfg.SetConfigType("yaml")

		out, err := yaml.Marshal(config)
		if err != nil {
			panic(err)
		}
		viperCfg.ReadConfig(bytes.NewBuffer(out))

		for _, key := range viperCfg.AllKeys() {
			value := viperCfg.Get(key)

			if hasNumber {
				hardforks.BlockHardforks[key] = append(hardforks.BlockHardforks[key], BlockHardfork{BlockNum: tryConvertToInt64(block), Value: value})
			}

			if hasSession {
				hardforks.SessionHardforks[key] = append(hardforks.SessionHardforks[key], SessionHardfork{SessionNum: tryConvertToInt64(session), Value: value})
			}
		}
	}

	return hardforks
}

func tryConvertToInt64(value interface{}) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	case string:
		ret, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			panic(fmt.Sprintf("Cannot parse value %s. %v", value, err))
		}
		return ret
	default:
		panic("Unexpect value type")
	}
}
