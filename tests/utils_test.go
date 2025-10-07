package tests

import (
	"fmt"
	"testing"
	"statusframe/backend/utils"
)

func TestMapStatusCode(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{200, "up"},
		{404, "not found"},
		{500, "server error"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.code), func(t *testing.T) {
			got := utils.MapStatusCode(tt.code)
			if got != tt.expected {
				t.Errorf("mapStatusCode(%d) = %q; want %q", tt.code, got, tt.expected)
			}
		})
	}
}
