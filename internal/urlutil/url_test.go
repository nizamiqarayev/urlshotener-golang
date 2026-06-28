package urlutil

import "testing"

func TestNormalizeHTTPURL(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		{
			name: "full https URL",
			in:   "https://www.google.com",
			want: "https://www.google.com",
			ok:   true,
		},
		{
			name: "bare domain gets https",
			in:   "www.google.com",
			want: "https://www.google.com",
			ok:   true,
		},
		{
			name: "trims whitespace and quotes",
			in:   ` "https://www.google.com" `,
			want: "https://www.google.com",
			ok:   true,
		},
		{
			name: "rejects unsupported scheme",
			in:   "ftp://example.com",
			ok:   false,
		},
		{
			name: "rejects spaces",
			in:   "https://exa mple.com",
			ok:   false,
		},
		{
			name: "rejects missing dot domain",
			in:   "hello",
			ok:   false,
		},
		{
			name: "allows localhost",
			in:   "http://localhost:3000/test",
			want: "http://localhost:3000/test",
			ok:   true,
		},
		{
			name: "preserves path and query",
			in:   "https://shop.com/products/123?color=black",
			want: "https://shop.com/products/123?color=black",
			ok:   true,
		},
		{
			name: "preserves fragment",
			in:   "https://docs.com/guide#install",
			want: "https://docs.com/guide#install",
			ok:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := NormalizeHTTPURL(tt.in)
			if ok != tt.ok {
				t.Fatalf("expected ok=%v, got %v", tt.ok, ok)
			}

			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
