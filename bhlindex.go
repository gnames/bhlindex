package bhlindex

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

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
	emptyEnvs := make([]string, 0, 5)
	envVars := [5]string{"POSTGRES_HOST", "POSTGRES_USER", "POSTGRES_PASSWORD",
		"POSTGRES_DB", "BHL_DIR"}
	for i, v := range envVars {
		val, ok := os.LookupEnv(v)
		if ok {
			envVars[i] = val
		} else {
			emptyEnvs = append(emptyEnvs, v)
		}
	}
	if len(emptyEnvs) > 0 {
		envs := strings.Join(emptyEnvs, ", ")
		panic(errors.New(fmt.Sprintf("Environment variables %s are not defined",
			envs)))
	}
	return Env{DbHost: envVars[0], DbUser: envVars[1], DbPass: envVars[2],
		Db: envVars[3], BHLDir: envVars[4]}
}

// UUID4 returns random (version 4) UUID as a string
func UUID4() string {
	return uuid.NewV4().String()
}

// DbInit returns a connection to bhlindex database
func DbInit() (*sql.DB, error) {
	var db *sql.DB
	var err error
	env := EnvVars()
	params := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		env.DbUser, env.DbPass, env.DbHost, env.Db)
	db, err = sql.Open("postgres", params)
	return db, err
}
