package conf

import "github.com/jessevdk/go-flags"

type confOptions struct {
	paths            []configPath
	args             []string
	delimiter        string
	noValidation     bool
	decoders         map[string]DecoderFunc
	configFlagOption *flags.Option
	flagOpts         flags.Options
}

type configPath struct {
	path     string
	optional bool
}

type ConfOption interface {
	apply(*confOptions)
}

type ConfOptions []ConfOption

func (opts ConfOptions) apply(o *confOptions) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

type funcConfOption struct {
	f func(*confOptions)
}

func (fo *funcConfOption) apply(o *confOptions) {
	fo.f(o)
}

func newFuncConfOption(f func(*confOptions)) *funcConfOption {
	return &funcConfOption{
		f: f,
	}
}

func Paths(paths ...string) ConfOption {
	return newFuncConfOption(func(o *confOptions) {
		for _, path := range paths {
			o.paths = append(o.paths, configPath{path: path})
		}
	})
}

func OptionalPaths(paths ...string) ConfOption {
	return newFuncConfOption(func(o *confOptions) {
		for _, path := range paths {
			o.paths = append(o.paths, configPath{path: path, optional: true})
		}
	})
}

func ConfigFlag(longName string, paths ...string) ConfOption {
	return newFuncConfOption(func(o *confOptions) {
		if o.configFlagOption == nil {
			o.configFlagOption = &flags.Option{}
		}
		o.configFlagOption.LongName = longName
		o.configFlagOption.Default = paths
	})
}

func ConfigFlagOption(configFlagOption *flags.Option) ConfOption {
	return newFuncConfOption(func(o *confOptions) {
		o.configFlagOption = configFlagOption
	})
}

func Delimiter(delimiter string) ConfOption {
	return newFuncConfOption(func(o *confOptions) {
		o.delimiter = delimiter
	})
}

func Args(args []string) ConfOption {
	return newFuncConfOption(func(o *confOptions) {
		o.args = args
	})
}

func NoValidation() ConfOption {
	return newFuncConfOption(func(o *confOptions) {
		o.noValidation = true
	})
}

func AddDecoder(ext string, dec DecoderFunc) ConfOption {
	return newFuncConfOption(func(o *confOptions) {
		o.decoders[ext] = dec
	})
}

func WithFlagOpts(flagOpts flags.Options) ConfOption {
	return newFuncConfOption(func(o *confOptions) {
		o.flagOpts = flagOpts
	})
}
