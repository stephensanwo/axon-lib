package types

const (
	DEVELOPMENT string = "development"
	PRODUCTION  string = "production"
)

type Metadata struct {
	Environment string `yaml:"environment"`
	Version     string `yaml:"version"`
}
