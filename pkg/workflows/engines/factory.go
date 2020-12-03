package engines

type Factory interface {
	Create(target string) Engine
	Type() ExperimentType
}
