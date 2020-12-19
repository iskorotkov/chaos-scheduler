package targets

type Target struct {
	Pod           string
	Deployment    string
	Node          string
	MainContainer string
	Containers    []string
	AppLabel      string
	Labels        map[string]string
	Annotations   map[string]string
}
