package capabilities

import (
	"testing"
)

func TestSpec_Validate(t *testing.T) {
	tests := []struct {
		name    string
		spec   Spec
		wantErr bool
	}{
		{
			name: "valid spec",
			spec: Spec{
				Name:        "test-capability",
				Language:    "python",
				Description: "A test capability",
				MainFile:    "main.py",
				TestFile:    "test.py",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			spec: Spec{
				Language:    "python",
				Description: "A test capability",
				MainFile:    "main.py",
				TestFile:    "test.py",
			},
			wantErr: true,
		},
		{
			name: "empty name",
			spec: Spec{
				Name:        "",
				Language:    "python",
				Description: "A test capability",
				MainFile:    "main.py",
				TestFile:    "test.py",
			},
			wantErr: true,
		},
		{
			name: "missing language",
			spec: Spec{
				Name:        "test-capability",
				Language:    "",
				Description: "A test capability",
				MainFile:    "main.py",
				TestFile:    "test.py",
			},
			wantErr: true,
		},
		{
			name: "invalid language",
			spec: Spec{
				Name:        "test-capability",
				Language:    "javascript",
				Description: "A test capability",
				MainFile:    "main.py",
				TestFile:    "test.py",
			},
			wantErr: true,
		},
		{
			name: "missing mainFile",
			spec: Spec{
				Name:        "test-capability",
				Language:    "python",
				Description: "A test capability",
				MainFile:    "",
				TestFile:    "test.py",
			},
			wantErr: true,
		},
		{
			name: "missing testFile",
			spec: Spec{
				Name:        "test-capability",
				Language:    "python",
				Description: "A test capability",
				MainFile:    "main.py",
				TestFile:    "",
			},
			wantErr: true,
		},
		{
			name: "missing description",
			spec: Spec{
				Name:        "test-capability",
				Language:    "python",
				Description: "",
				MainFile:    "main.py",
				TestFile:    "test.py",
			},
			wantErr: true,
		},
		{
			name: "valid bash",
			spec: Spec{
				Name:        "bash-script",
				Language:    "bash",
				Description: "A bash script",
				MainFile:    "script.sh",
				TestFile:    "test.sh",
			},
			wantErr: false,
		},
		{
			name: "valid go",
			spec: Spec{
				Name:        "go-function",
				Language:    "go",
				Description: "A go function",
				MainFile:    "main.go",
				TestFile:    "main_test.go",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spec.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}