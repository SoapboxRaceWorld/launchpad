package launchpad

import (
	"os"
	"fmt"
	"sync"
)

var once sync.Once
var isDebug bool

func EnableDebug() {
	isDebug = true
}

func CheckEnvVar(varName string) bool {
	_, res := os.LookupEnv(varName)

	return res
}

func CheckEnv() error {
	if !CheckEnvVar("DB_HOST") {
		return fmt.Errorf("missing environent variable: DB_HOST")
	}

	if !CheckEnvVar("DB_NAME") {
		return fmt.Errorf("missing environent variable: DB_NAME")
	}

	if !CheckEnvVar("DB_USER") {
		return fmt.Errorf("missing environent variable: DB_USER")
	}

	if !CheckEnvVar("DB_PASS") {
		return fmt.Errorf("missing environent variable: DB_PASS")
	}

	if !CheckEnvVar("GITHUB_AUTH_TOKEN") {
		return fmt.Errorf("missing environent variable: GITHUB_AUTH_TOKEN")
	}

	return nil
}
