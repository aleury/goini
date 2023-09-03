package goini

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Section struct {
	Name          string         `json:"name"`
	KeyValuePairs []KeyValuePair `json:"keyValuePairs"`
}

type File struct {
	Name     string    `json:"name"`
	Sections []Section `json:"sections"`
}
