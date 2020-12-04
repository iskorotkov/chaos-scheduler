package generators

type Params struct {
	Stages int
	Seed   int64
}

type Generator interface {
	Generate(params Params) (Scenario, error)
}
