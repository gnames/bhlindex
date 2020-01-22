package bhlindex

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// Env is a collection of environment variables.
type Env struct {
	DbHost      string
	DbUser      string
	Db          string
	BHLDir      string
	PrefSources []int
}

// Check handles error checking, and panicks if error is not nil.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// EnvVars imports all environment variables relevant for the data conversion.
func EnvVars() Env {
	emptyEnvs := make([]string, 0, 4)
	envVars := [5]string{"POSTGRES_HOST", "POSTGRES_USER", "POSTGRES_DB",
		"BHL_DIR", "PREF_SOURCES"}
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
		panic(fmt.Errorf("environment variables %s are not defined", envs))
	}
	var sources []int
	for _, v := range strings.Split(envVars[4], ",") {
		v = strings.Trim(v, " ")
		source, err := strconv.Atoi(v)
		if err != nil {
			log.Println(err)
			log.Fatal(`
PREF_SOURCES env variable should be a comma-separated list of ints:
Example: "1,2,3,4"`)
		}
		sources = append(sources, source)
	}
	return Env{DbHost: envVars[0], DbUser: envVars[1], Db: envVars[2],
		BHLDir: envVars[3], PrefSources: sources}
}

// UUID4 returns random (version 4) UUID as a string.
func UUID4() string {
	return uuid.NewV4().String()
}

// DbInit returns a connection to bhlindex database.
func DbInit() (*sql.DB, error) {
	var db *sql.DB
	var err error
	env := EnvVars()
	params := fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable",
		env.DbUser, env.DbHost, env.Db)
	db, err = sql.Open("postgres", params)
	return db, err
}
