package capabilities

import (
	"fmt"
)

type Spec struct {
	Name        string
	Type        string
	Description string
	Commands    map[string]string
}

func (s *Spec) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name is required")
	}
	if s.Type == "" {
		return fmt.Errorf("type is required")
	}
	if s.Description == "" {
		return fmt.Errorf("description is required")
	}
	return nil
}