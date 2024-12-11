package main

import (
  "strings"
  "errors"
  "log"
  "reflect"
  "fmt"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/mitchellh/mapstructure"
  "github.com/buildkite/agent-stack-k8s/v2/internal/controller/config"
	"github.com/spf13/viper"
)

func main() {
  configFile := "config.yaml"
	v := viper.NewWithOptions(
		viper.KeyDelimiter("::"),
		viper.EnvKeyReplacer(strings.NewReplacer("-", "_")),
	)
	v.SetConfigFile(configFile)

	if err := v.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			log.Fatalf("failed to read config: %w", err)
		}
	}
	cfg := &config.Config{}
	decodeHook := viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		stringToResourceQuantity,
		config.StringToInterposer,
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	))
	if err := v.UnmarshalExact(cfg, useJSONTagForDecoder, decodeHook); err != nil {
		fmt.Errorf("failed to parse config: %w", err)
	}

  fmt.Printf("%+v", cfg)
}

var resourceQuantityType = reflect.TypeOf(resource.Quantity{})

func stringToResourceQuantity(f, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}
	if t != resourceQuantityType {
		return data, nil
	}
	return resource.ParseQuantity(data.(string))
}

func useJSONTagForDecoder(c *mapstructure.DecoderConfig) {
	c.TagName = "json"
}

