package e2e

import "testing"

func TestApplication(t *testing.T) {
	app := NewApplication(getRandomPort())
	app.Stop()
}
