package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Exercise the linker flags from the release configuration, rather than a
// separate copy that could stay correct while the release configuration drifts.
func TestReleaseMetadata(t *testing.T) {
	data, err := os.ReadFile(".goreleaser.yml")
	if err != nil {
		t.Fatal(err)
	}
	var config struct {
		Builds []struct {
			Ldflags []string `yaml:"ldflags"`
		} `yaml:"builds"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}
	if len(config.Builds) == 0 {
		t.Fatal("release configuration has no builds")
	}

	metadata := struct{ Tag, ShortCommit, Date string }{
		Tag:         "v0.0.0-metadata-test",
		ShortCommit: "abc1234",
		Date:        "2026-01-02T03:04:05Z",
	}
	for i, build := range config.Builds {
		if len(build.Ldflags) == 0 {
			t.Fatalf("release build %d has no linker flags", i)
		}
		flags, err := template.New("ldflags").Parse(strings.Join(build.Ldflags, " "))
		if err != nil {
			t.Fatal(err)
		}
		var rendered bytes.Buffer
		if err := flags.Execute(&rendered, metadata); err != nil {
			t.Fatal(err)
		}

		binary := filepath.Join(t.TempDir(), "checkson")
		if runtime.GOOS == "windows" {
			binary += ".exe"
		}
		cmd := exec.Command("go", "build", "-mod=readonly", "-ldflags", rendered.String(), "-o", binary, ".")
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("release build %d failed: %v\n%s", i, err, output)
		}
		output, err := exec.Command(binary, "version").CombinedOutput()
		if err != nil {
			t.Fatalf("version command failed: %v\n%s", err, output)
		}
		for _, want := range []string{
			`version:"` + metadata.Tag + `"`,
			`buildTime:"` + metadata.Date + `"`,
			`gitCommit:"` + metadata.ShortCommit + `"`,
		} {
			if !strings.Contains(string(output), want) {
				t.Errorf("release build %d: expected %s in %s", i, want, output)
			}
		}
	}
}
