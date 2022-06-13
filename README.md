### Conf

Conf is a front-end for https://github.com/jessevdk/go-flags that supports config files.

### Examples
```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/go-chai/conf"
	"gopkg.in/yaml.v3"
)

func main() {
	cfg, err := conf.Load[Config](conf.ConfigFlag("conf", "examples/conf/config.yaml"))
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

```

```bash

go run ./examples/conf/ --help

Usage:
  conf [OPTIONS]

Application Options:
      --i=      int
      --id=     int with a default (default: 1)
      --f=      float64
      --fd=     float64 with a default (default: -3.14)
  -3, --n       numeric flag
      --str=    string
      --strd=   string with a default (default: abc)
      --strnot= string not unquoted
      --t=      time
      --td=     time with a default (default: 1m)
      --m=      map
      --md=     map with a default (default: a:1)
      --s=      slice
      --sd=     slice with a default (default: 1, 2)

Config:
      --conf= config file paths (default: examples/conf/config.yaml)

Help Options:
  -h, --help    Show this help message

```


```bash

go run ./examples/conf/

int: 3
intDefault: 13
float64: 2.712
float64Default: 1.1234
numericFlag: false
string: asdf
stringDefault: defg
stringNotUnquoted: ""
time: 13s
timeDefault: 11m0s
map:
    val1: 3
    val2: 4
mapDefault:
    a: 1
    val21: 21
    val22: 22
slice:
    - 1
    - 2
    - 3
sliceDefault:
    - 4
    - 5
    - 6
```