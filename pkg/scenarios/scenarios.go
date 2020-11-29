package scenarios

type ActionTemplate struct {
	Name     string
	Template map[string]interface{}
}

type PlannedAction struct {
	Name     string
	Template map[string]interface{}
}

type Stage []PlannedAction

type Scenario []Stage

type Config struct {
	Stages int
}

type Generator interface {
	Generate(actions []ActionTemplate, config Config) (Scenario, error)
}
