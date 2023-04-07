package lib

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type JSONParser struct {
}

func (p *JSONParser) Parse(data []byte) (schemas []Schema, err error) {
	err = json.Unmarshal(data, &schemas)
	return
}

type YAMLParser struct {
}

func (p *YAMLParser) Parse(data []byte) (schemas []Schema, err error) {
	err = yaml.Unmarshal(data, &schemas)
	return
}
