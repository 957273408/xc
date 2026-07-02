package config

type Proxy struct {
	UpstreamBaseURL   string `mapstructure:"upstream-base-url" json:"upstream-base-url" yaml:"upstream-base-url"`
	SignKey           string `mapstructure:"sign-key" json:"sign-key" yaml:"sign-key"`
	SignAlgorithm     string `mapstructure:"sign-algorithm" json:"sign-algorithm" yaml:"sign-algorithm"`
	UpstreamQueryMode string `mapstructure:"upstream-query-mode" json:"upstream-query-mode" yaml:"upstream-query-mode"`
	RequestTimeout    int    `mapstructure:"request-timeout" json:"request-timeout" yaml:"request-timeout"`
}
