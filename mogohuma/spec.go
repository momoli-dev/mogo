package mogohuma

import (
	"fmt"
	"os"
)

// CreateJSONSpecFile writes the OpenAPI spec to path as JSON, replacing any existing file.
func (r *Registry) CreateJSONSpecFile(path string) error {
	jsonData, err := r.HumaAPI.OpenAPI().MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to OpenAPI.MarshalJSON: %w", err)
	}

	const fileMode os.FileMode = 0644

	err = os.WriteFile(path, jsonData, fileMode)
	if err != nil {
		return fmt.Errorf("failed to os.WriteFile: %w", err)
	}

	return nil
}

// CreateYAMLSpecFile writes the OpenAPI spec to path as YAML, replacing any existing file.
func (r *Registry) CreateYAMLSpecFile(path string) error {
	yamlData, err := r.HumaAPI.OpenAPI().YAML()
	if err != nil {
		return fmt.Errorf("failed to OpenAPI.YAML: %w", err)
	}

	const fileMode os.FileMode = 0644

	err = os.WriteFile(path, yamlData, fileMode)
	if err != nil {
		return fmt.Errorf("failed to os.WriteFile: %w", err)
	}

	return nil
}
