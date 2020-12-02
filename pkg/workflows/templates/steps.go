package templates

type Step struct {
	Name     string `yaml:"name" json:"name"`
	Template string `yaml:"template" json:"template"`
}

type StepsTemplate struct {
	Name  string   `yaml:"name" json:"name"`
	Steps [][]Step `yaml:"steps" json:"steps"`
}

func (s StepsTemplate) Id() string {
	return s.Name
}

func NewStepsTemplate(ids [][]string) StepsTemplate {
	res := StepsTemplate{"entry", make([][]Step, 0)}

	for _, stage := range ids {
		newStage := make([]Step, 0)

		for _, id := range stage {
			newStage = append(newStage, Step{id, id})
		}

		res.Steps = append(res.Steps, newStage)
	}

	return res
}
