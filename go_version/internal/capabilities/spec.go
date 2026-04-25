package capabilities

import (
	"fmt"
	"strings"
)

type Spec struct {
	Name        string
	Language    string
	Description string
	Inputs      []string
	Outputs     []string
	Constraints []string
	MainFile    string
	TestFile    string
}

func (s *Spec) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name is required")
	}
	if s.Language == "" {
		return fmt.Errorf("language is required")
	}
	validLanguages := []string{"python", "bash", "go"}
	if !contains(validLanguages, s.Language) {
		return fmt.Errorf("language must be one of: %s", strings.Join(validLanguages, ", "))
	}
	if s.MainFile == "" {
		return fmt.Errorf("mainFile is required")
	}
	if s.TestFile == "" {
		return fmt.Errorf("testFile is required")
	}
	if s.Description == "" {
		return fmt.Errorf("description is required")
	}
	return nil
}

func (s *Spec) ToMarkdown() string {
	var sb strings.Builder
	sb.WriteString("# " + s.Name + "\n\n")
	sb.WriteString("## Description\n" + s.Description + "\n\n")
	sb.WriteString("## Language\n" + s.Language + "\n\n")
	sb.WriteString("## Inputs\n")
	for _, input := range s.Inputs {
		sb.WriteString("- " + input + "\n")
	}
	sb.WriteString("\n## Outputs\n")
	for _, output := range s.Outputs {
		sb.WriteString("- " + output + "\n")
	}
	sb.WriteString("\n## Constraints\n")
	for _, constraint := range s.Constraints {
		sb.WriteString("- " + constraint + "\n")
	}
	sb.WriteString("\n## Files\n")
	sb.WriteString("- MainFile: " + s.MainFile + "\n")
	sb.WriteString("- TestFile: " + s.TestFile + "\n")
	return sb.String()
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}