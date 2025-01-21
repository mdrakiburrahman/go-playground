package emptyexporter

import "fmt"

type EncodingType int

const (
	OTLPCSV EncodingType = iota
)

func (e EncodingType) String() string {
	return [...]string{"otlp_csv"}[e]
}

func ParseEncodingType(s string) (EncodingType, error) {
	switch s {
	case "otlp_csv":
		return OTLPCSV, nil
	default:
		return 0, fmt.Errorf("invalid encoding type: %s", s)
	}
}

type Config struct {
	ShouldLog bool   `mapstructure:"should_log"`
	Encoding  string `mapstructure:"encoding"`
	encoding  EncodingType
}

func (c *Config) Validate() error {
	encoding, err := ParseEncodingType(c.Encoding)
	if err != nil {
		return err
	}
	c.encoding = encoding
	return nil
}
