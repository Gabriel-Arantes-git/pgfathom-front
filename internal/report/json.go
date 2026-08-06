package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/lvcas-dotcom/pgfathom/internal/model"
)

// JSON writes the result as the versioned public contract. It is what the CI
// baseline and third-party tooling consume, so the shape is an API rather than
// a rendering choice.
func JSON(w io.Writer, r *model.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("writing JSON result: %w", err)
	}
	return nil
}
