package adapters

import "testing"

func TestPublicHostFromAPIURL(t *testing.T) {
	got, err := publicHostFromAPIURL("https://blog.example/api/v1")
	if err != nil {
		t.Fatalf("publicHostFromAPIURL: %v", err)
	}
	if got != "https://blog.example" {
		t.Fatalf("host = %q", got)
	}

	for _, raw := range []string{
		"",
		"file:///tmp/api",
		"https://user:pass@blog.example/api/v1",
		"https://blog.example/api/v1?x=1",
		"https://blog.example/api/v1#fragment",
	} {
		if _, err := publicHostFromAPIURL(raw); err == nil {
			t.Errorf("publicHostFromAPIURL(%q) error = nil", raw)
		}
	}
}
