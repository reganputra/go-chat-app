package env

import "os"

var Env map[string]string

func GetEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func SetupEnvFile() {

}
