package platform_exercise

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	os.Setenv("postgresURL", "postgres://root:postgres@localhost:5432/postgres?sslmode=disable")
}

func shutdown() {}
