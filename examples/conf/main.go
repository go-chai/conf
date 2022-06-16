package main

import (
	"log"
	"os"
	"time"

	"github.com/go-chai/conf"
	"gopkg.in/yaml.v3"
)

func main() {
	cfg, err := conf.Load[Config](
		conf.ConfigFlag("conf"),
		conf.OptionalPaths("testdata/config.yaml"),
	)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	yaml.NewEncoder(os.Stdout).Encode(cfg)
}

type Config struct {
	Int        int `long:"i" yaml:"int" description:"int"`
	IntDefault int `long:"id" default:"1" yaml:"intDefault" description:"int with a default"`

	Float64        float64 `long:"f" yaml:"float64" description:"float64"`
	Float64Default float64 `long:"fd" default:"-3.14"  yaml:"float64Default" description:"float64 with a default"`

	NumericFlag bool `short:"3" long:"n" yaml:"numericFlag" description:"numeric flag"`

	String            string `long:"str" yaml:"string" description:"string"`
	StringDefault     string `long:"strd" default:"abc" yaml:"stringDefault" description:"string with a default"`
	StringNotUnquoted string `long:"strnot" unquote:"false" yaml:"stringNotUnquoted" description:"string not unquoted"`

	Time        time.Duration `long:"t" yaml:"time" description:"time"`
	TimeDefault time.Duration `long:"td" default:"1m" yaml:"timeDefault" description:"time with a default"`

	Map        map[string]int `long:"m" yaml:"map" description:"map"`
	MapDefault map[string]int `long:"md" default:"a:1" yaml:"mapDefault" description:"map with a default"`

	Slice        []int `long:"s" yaml:"slice" description:"slice"`
	SliceDefault []int `long:"sd" default:"1" default:"2" yaml:"sliceDefault" description:"slice with a default"`
}
