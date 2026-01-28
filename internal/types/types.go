package types

type Text struct {
	Chunks []string `json:"chunks"`
}

type Tool struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type MockResponse struct {
	Agent string `json:"agent,omitempty"`
	Text  Text   `json:"text"`
	Tool  Tool   `json:"tool"`
}
