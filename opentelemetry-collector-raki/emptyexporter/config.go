package emptyexporter

type Config struct {
	ShouldLog bool `mapstructure:"should_log"`
}

func (c *Config) Validate() error {
	return nil
}
