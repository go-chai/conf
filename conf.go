package conf

import (
	stderr "errors"
	"io"
	"io/fs"
	"os"
	"reflect"

	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
	"github.com/jessevdk/go-flags"
)

var validate = validator.New()

func Load[T any](opts ...ConfOption) (*T, error) {
	cfg := new(T)
	cfgDefaults := new(T)
	var err error

	copts := &confOptions{
		paths:        nil,
		args:         os.Args[1:],
		delimiter:    "-",
		noValidation: false,
		decoders:     DefaultDecoders,
		flagOpts:     flags.Default,
	}

	for _, opt := range opts {
		opt.apply(copts)
	}

	// Step 1:
	// 	load the defaults
	// 	obtain the config file paths
	// 	handle the Help message
	paths, err := mergeDefaults(cfgDefaults, copts)
	if err != nil {
		flagsErr := new(flags.Error)
		if errors.As(err, &flagsErr) {
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		return nil, errors.Wrap(err, "failed to parse command line args")
	}

	// Step 2:
	// 	override with values from the config files
	_, err = mergeConfigFiles(copts, cfg, append(copts.paths, paths...)...)
	if err != nil {
		return nil, err
	}

	// Step 3:
	// 	create a parser that does not add default values
	// 	override with values from flags + env variables
	err = mergeWithoutDefaults(cfg, copts)
	if err != nil {
		return nil, err
	}

	// Step 4: Merge defaults
	// 	override the empty values with defaults
	//  use a custom transformer to avoid map Options containing the default values if they are set in a config file
	err = mergo.Merge(cfg, cfgDefaults, mergo.WithTransformers(mapDefaultsTransformer{}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to merge defaults")
	}

	if !copts.noValidation {
		err = validate.Struct(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to validate config")
		}
	}

	return cfg, nil
}

type mapDefaultsTransformer struct {
}

func (mapDefaultsTransformer) Transformer(t reflect.Type) func(dst, src reflect.Value) error {
	// use a custom transformer only for maps
	if t.Kind() != reflect.Map {
		return nil
	}
	return func(dst, src reflect.Value) error {
		// only merge if dst is empty
		if dst.IsNil() || len(dst.MapKeys()) == 0 && dst.CanSet() {
			dst.Set(src)
		}
		return nil
	}
}

type fileConfig struct {
	ConfigFilePaths []string `long:"conf" description:"config file paths" default:"config.yaml"`
}

func parseFlags(cfg any, defaults bool, copts *confOptions) ([]configPath, error) {
	cfgF := &fileConfig{}
	p := flags.NewParser(cfg, copts.flagOpts)
	p.NamespaceDelimiter = copts.delimiter

	if copts.configFlagOption != nil {
		g, err := p.AddGroup("Config", "", cfgF)
		if err != nil {
			return nil, errors.Wrap(err, "failed to add config group")
		}
		err = mergo.Merge(g.Options()[0], copts.configFlagOption, mergo.WithOverride)
		if err != nil {
			return nil, errors.Wrap(err, "failed to merge config flag option")
		}
	}

	if !defaults {
		eachOption(p.Command, func(c *flags.Command, g *flags.Group, o *flags.Option) {
			o.Default = []string{}
		})
	}

	_, err := p.ParseArgs(copts.args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse command line args")
	}

	if copts.configFlagOption != nil {
		paths := make([]configPath, len(cfgF.ConfigFilePaths))

		for i, path := range cfgF.ConfigFilePaths {
			paths[i] = configPath{
				path: path,
			}
		}
		return paths, nil
	}
	return nil, nil
}

func mergeDefaults(cfg any, copts *confOptions) ([]configPath, error) {
	return parseFlags(cfg, true, copts)
}

func mergeConfigFiles(copts *confOptions, cfg any, paths ...configPath) (loadedPaths []string, err error) {
	loadedPaths = make([]string, 0)
	for _, path := range paths {
		ok, err := mergeConfigFile(path.optional, copts, cfg, path.path)
		if err != nil {
			return loadedPaths, errors.Wrapf(err, "failed to merge config file %s", path.path)
		}
		if ok {
			loadedPaths = append(loadedPaths, path.path)
		}
	}
	return loadedPaths, nil
}

func mergeConfigFile(optional bool, copts *confOptions, cfg any, path string) (ok bool, err error) {
	f, err := os.Open(path)
	if err != nil {
		if optional && stderr.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, errors.Wrapf(err, "failed to open required config file %s", path)
	}
	defer f.Close()

	dec, err := getDecoder(copts, path)
	if err != nil {
		return false, err
	}
	if err := dec(cfg, f); err != nil && !stderr.Is(err, io.EOF) {
		return false, err
	}
	return true, nil
}

func mergeWithoutDefaults(cfg any, copts *confOptions) error {
	_, err := parseFlags(cfg, false, copts)
	return err
}

func eachOption(c *flags.Command, f func(*flags.Command, *flags.Group, *flags.Option)) {
	eachCommand(c, func(c *flags.Command) {
		eachGroup(c.Group, func(g *flags.Group) {
			for _, option := range g.Options() {
				f(c, g, option)
			}
		})
	}, true)
}

func eachCommand(c *flags.Command, f func(*flags.Command), recurse bool) {
	f(c)

	for _, cc := range c.Commands() {
		if recurse {
			eachCommand(cc, f, true)
		} else {
			f(cc)
		}
	}
}

func eachGroup(g *flags.Group, f func(*flags.Group)) {
	f(g)

	for _, gg := range g.Groups() {
		eachGroup(gg, f)
	}
}
