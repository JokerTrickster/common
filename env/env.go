package env

/*
	환경 변수 초기화 및 관리
*/

import (
	"fmt"
	"os"
)

type EnvStruct struct {
	Port               string
	Env                string
	IsLocal            bool
	GoogleClientID     string
	GoogleClientSecret string
}

// Env : Global environment variables
var Env EnvStruct

// InitEnv initializes environment variables
func InitEnv() error {
	envVarNames := []string{"PORT", "ENV", "IS_LOCAL"}
	envs, err := getOSLookupEnv(envVarNames)
	if err != nil {
		return err
	}
	Env = EnvStruct{
		Port:    envs["PORT"],
		Env:     envs["ENV"],
		IsLocal: isLocalEnv(envs["IS_LOCAL"]),
	}
	return nil
}

// isLocalEnv checks if the environment is local
func isLocalEnv(isLocal string) bool {
	return isLocal == "true"
}

// getOSLookupEnv fetches environment variables by their names
func getOSLookupEnv(envVarNames []string) (map[string]string, error) {
	result := map[string]string{}
	for _, envVarName := range envVarNames {
		value, ok := os.LookupEnv(envVarName)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve environment variable: %s", envVarName)
		}
		result[envVarName] = value
	}
	return result, nil
}

func PtrTrue() *bool {
	b := true
	return &b
}

func PtrFalse() *bool {
	b := false
	return &b
}
