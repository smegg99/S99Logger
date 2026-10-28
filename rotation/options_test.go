// rotation/options_test.go

package rotation

import (
	"errors"
	"testing"
)

// New and Configure share one compression default, and an unknown name is refused.
func TestCompressionDefaultsToGzip(t *testing.T) {
	opts, err := normalizeOptions(Options{Directory: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Compression != "gzip" {
		t.Fatalf("Compression = %q, want gzip", opts.Compression)
	}
	if _, err := normalizeOptions(Options{Directory: t.TempDir(), Compression: "zip"}); !errors.Is(err, errCompressionUnknown) {
		t.Fatalf("unknown compression error = %v", err)
	}
}
