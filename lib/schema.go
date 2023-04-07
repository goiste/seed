package lib

type Schema struct {
	Name    string   `json:"schema" yaml:"schema"`
	Tables  []Table  `json:"tables" yaml:"tables"`
	EnvVars []EnvVar `json:"env_vars" yaml:"env_vars"`
}

type EnvVar struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

type Table struct {
	Name        string           `json:"name" yaml:"name"`
	PrimaryKey  string           `json:"primary_key,omitempty" yaml:"primary_key,omitempty"`
	ForeignKeys []ForeignKey     `json:"foreign_keys,omitempty" yaml:"foreign_keys,omitempty"`
	Columns     []string         `json:"columns" yaml:"columns"`
	Values      []map[string]any `json:"values" yaml:"values"`
}

type ForeignKey struct {
	RefTable string `json:"ref_table" yaml:"ref_table"`
	Column   string `json:"column" yaml:"column"`
}
