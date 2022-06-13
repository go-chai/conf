package conf

import (
	"io"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/cockroachdb/errors"
	"gopkg.in/yaml.v3"
)

type DecoderFunc func(cfg any, r io.Reader) error

func getDecoder(copts *confOptions, path string) (DecoderFunc, error) {
	ext := filepath.Ext(path)
	dec, ok := copts.decoders[ext]
	if !ok {
		return nil, errors.Errorf("no decoder for %s", ext)
	}
	return dec, nil
}

var DefaultDecoders = map[string]DecoderFunc{
	".yaml": YAMLDecoder,
	".yml":  YAMLDecoder,
	".json": JSONDecoder,
	".toml": TOMLDecoder,
}

var YAMLDecoder = func(cfg any, r io.Reader) error {
	err := yaml.NewDecoder(r).Decode(cfg)
	return errors.Wrap(err, "failed to decode yaml")
}
var JSONDecoder = func(cfg any, r io.Reader) error {
	err := yaml.NewDecoder(r).Decode(cfg)
	return errors.Wrap(err, "failed to decode json")
}
var TOMLDecoder = func(cfg any, r io.Reader) error {
	_, err := toml.DecodeReader(r, cfg)
	return errors.Wrap(err, "failed to decode toml")
}
