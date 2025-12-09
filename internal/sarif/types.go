package sarif

// ROOT
type Sarif struct {
	Version string `json:"version"`
	Runs    []Run  `json:"runs"`
}

type Run struct {
	Tool    Tool     `json:"tool"`
	Results []Result `json:"results"`
}

type Tool struct {
	Driver Driver `json:"driver"`
}

type Driver struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// One res  =  one finding
type Result struct {
	RuleID     string     `json:"ruleId"`
	Message    Message    `json:"message"`
	Locations  []Location `json:"locations"`
	Properties Properties `json:"properties"`
}

type Message struct {
	Text string `json:"text"`
}

type Location struct {
	PhysicalLocation PhysicalLocation `json:"physicalLocation"`
}

type PhysicalLocation struct {
	ArtifactLocation ArtifactLocation `json:"artifactLocation"`
	Region           Region           `json:"region"`
}

type ArtifactLocation struct {
	URI string `json:"uri"`
}

type Region struct {
	StartLine int `json:"startLine"`
}

type Properties struct {
	Snippet    string  `json:"snippet"`
	Severity   string  `json:"severity"`
	Confidence float64 `json:"confidence"`

	AIVerdict string `json:"aiVerdict"`
	Source    string `json:"source"`
}
