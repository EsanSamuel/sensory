package models

type Runtime struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Fn   string `json:"fn"`
}

type LogEntry struct {
	Level     string  `json:"level"`
	Timestamp string  `json:"timestamp"`
	Project   string  `json:"project"`
	Service   string  `json:"service"`
	Message   string  `json:"message"`
	Runtime   Runtime `json:"runtime"`
	ApiKey    string  `json:"api_key"`
}
