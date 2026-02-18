package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNew verifies the creation of the Plugin instance.
func TestNew(t *testing.T) {
	tempDir := t.TempDir()
	contentFile := "example.json"

	// Create a dummy JSON content file for testing
	err := os.WriteFile(filepath.Join(tempDir, contentFile), []byte(`{}`), 0644)
	if err != nil {
		t.Fatalf("failed to create test content file: %v", err)
	}

	plugin, err := New(contentFile, tempDir, nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if plugin == nil {
		t.Fatalf("plugin instance is nil")
	}

	if plugin.Options.TemplateDir != tempDir {
		t.Errorf("expected TemplateDir to be %s, got %s", tempDir, plugin.Options.TemplateDir)
	}

	if plugin.Options.TemplateExt != "json" {
		t.Errorf("expected TemplateExt to be json, got %s", plugin.Options.TemplateExt)
	}
}

// TestLoadMessages verifies the loading of content files.
func TestLoadMessages(t *testing.T) {
	tempDir := t.TempDir()
	contentFile := "example.json"

	// Create a dummy JSON content file for testing
	err := os.WriteFile(filepath.Join(tempDir, contentFile), []byte(`{}`), 0644)
	if err != nil {
		t.Fatalf("failed to create test content file: %v", err)
	}

	plugin, _ := New(contentFile, tempDir, nil)
	if err := plugin.loadMessages(contentFile); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestPlugin_GetMessage(t *testing.T) {
	type args struct {
		code string
		lng  string
	}

	path, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get path: %v", err)
	}

	p, err := New("example", path, nil)
	if err != nil {
		t.Fatalf("failed to init plugin: %v", err)
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get kk content",
			args: args{
				code: "greetings",
				lng:  "kk",
			},
			want: "Cәлем, әлем!",
		},
		{
			name: "get ru content",
			args: args{
				code: "greetings",
				lng:  "ru",
			},
			want: "Привет, мир!",
		},

		{
			name: "get en content",
			args: args{
				code: "greetings",
				lng:  "en",
			},
			want: "Hello, World!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := p.GetMessage(tt.args.code, tt.args.lng); got != tt.want {
				t.Errorf("GetMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
