package core

import (
	"encoding/json"
	"encoding/xml"

	"gopkg.in/yaml.v2"
)

var (
	JsonEncoding       = &jsonEncoding{}
	JsonPrettyEncoding = &jsonEncoding{Indent: "  "}
	XmlEncoding        = &xmlEncoding{}
	XmlPrettyEncoding  = &xmlEncoding{Indent: "  "}
	YamlEncoding       = &yamlEncoding{}
)

type Encoding interface {
	Marshal(v interface{}) ([]byte, error)

	Unmarshal(data []byte, v interface{}) error
}

type jsonEncoding struct {
	Prefix, Indent string
}

func (e *jsonEncoding) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (e *jsonEncoding) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type xmlEncoding struct {
	Prefix, Indent string
}

func (e *xmlEncoding) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (e *xmlEncoding) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

type yamlEncoding struct {
}

func (e *yamlEncoding) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (e *yamlEncoding) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
