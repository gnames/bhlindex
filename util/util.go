package util

import (
	"os"

	uuid "github.com/satori/go.uuid"
)

// Env is a collection of environment variables
type Env struct {
	DbHost string
	DbUser string
	DbPass string
	Db     string
	BHLDir string
}

// Check handles error checking, and panicks if error is not nil
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// EnvVars imports all environment variables relevant for the data conversion.
func EnvVars() Env {
	env := Env{
		DbHost: os.Getenv("POSTGRES_HOST"),
		DbUser: os.Getenv("POSTGRES_USER"),
		DbPass: os.Getenv("POSTGRES_PASSWORD"),
		Db:     os.Getenv("POSTGRES_DB"),
		BHLDir: os.Getenv("BHL_DIR"),
	}
	return env
}

// UUID4 returns random (version 4) UUID as a string
func UUID4() string {
	return uuid.NewV4().String()
}
