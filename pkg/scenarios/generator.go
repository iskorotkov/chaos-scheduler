package scenarios

type Config struct {
	Stages int
	Seed   int64
}

type Generator interface {
	Generate(templates []TemplatedAction, config Config) (Scenario, error)
}
