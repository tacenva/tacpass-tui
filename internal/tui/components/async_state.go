package components

import (
	"github.com/tacenva/tacpass-tui/internal/styles"
	"github.com/tacenva/tacpass-tui/internal/tui/state"
)

func AsyncStateView(
	asyncState *state.Async,
	loadingText string,
) string {
	if asyncState.Loading {
		return loadingText
	}

	if asyncState.Error != nil {
		return styles.Error.Render(
			"[ERROR] " + asyncState.Error.Error(),
		)
	}

	return ""
}
