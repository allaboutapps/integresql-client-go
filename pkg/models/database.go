package models

type Database struct {
	TemplateHash string         `json:"templateHash"`
	Config       DatabaseConfig `json:"config"`
}
