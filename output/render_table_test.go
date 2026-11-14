package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderTable(t *testing.T) {
	for _, rows := range [][][]string{
		nil,
		{{"✓", "überwachung", "0 / 3"}, {"✖", "東京", "2 / 3"}},
	} {
		var out bytes.Buffer
		if err := RenderTable(&out, []string{"Status", "Name", "Recent failures"}, rows); err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(strings.TrimSpace(out.String()), "\n")
		if !strings.Contains(strings.ToUpper(lines[0]), "STATUS") || !strings.Contains(strings.ToUpper(lines[0]), "NAME") {
			t.Fatalf("missing header: %q", out.String())
		}
		if len(rows) > 0 && (len(lines) != len(rows)+2 || !strings.Contains(lines[1], "---")) {
			t.Fatalf("expected header separator and %d rows: %q", len(rows), out.String())
		}
		for _, row := range rows {
			for _, cell := range row {
				if !strings.Contains(out.String(), cell) {
					t.Errorf("missing cell %q: %q", cell, out.String())
				}
			}
		}
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "|") || strings.HasSuffix(strings.TrimSpace(line), "|") {
				t.Fatalf("unexpected outer border: %q", line)
			}
		}
	}
}
