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
				Type:        "system",
				Description: "A test capability",
				Commands:    map[string]string{"apt": "apt install test"},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			spec: Spec{
				Type:        "system",
				Description: "A test capability",
				Commands:    map[string]string{"apt": "apt install test"},
			},
			wantErr: true,
		},
		{
			name: "empty name",
			spec: Spec{
				Name:        "",
				Type:        "system",
				Description: "A test capability",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			spec: Spec{
				Name:        "test-capability",
				Type:        "",
				Description: "A test capability",
			},
			wantErr: true,
		},
		{
			name: "missing description",
			spec: Spec{
				Name:        "test-capability",
				Type:        "system",
				Description: "",
			},
			wantErr: true,
		},
		{
			name: "valid with no commands",
			spec: Spec{
				Name:        "test-capability",
				Type:        "system",
				Description: "A test capability",
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