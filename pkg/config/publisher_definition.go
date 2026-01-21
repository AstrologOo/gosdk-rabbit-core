package config

func NewPublisherDefinition(name string, cfg ...Config) PublisherDefinition {
	return PublisherDefinition{
		Name:   name,
		Config: cfg,
	}
}

type PublisherDefinition struct {
	Name   string
	Config []Config
}
