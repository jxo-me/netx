package handler

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
)

type NullStructRes struct{}

func Write(w io.Writer, data interface{}, format string) error {
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(data)
		return nil
	case "yaml":
		fallthrough
	default:
		enc := yaml.NewEncoder(w)
		defer enc.Close()
		enc.SetIndent(2)

		return enc.Encode(data)
	}
}
