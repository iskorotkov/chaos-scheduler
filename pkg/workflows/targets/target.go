package targets

type Target struct {
	Pod           string            `json:"pod"`
	Deployment    string            `json:"deployment"`
	Node          string            `json:"node"`
	MainContainer string            `json:"mainContainer"`
	Containers    []string          `json:"containers"`
	AppLabel      string            `json:"appLabel"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
}
