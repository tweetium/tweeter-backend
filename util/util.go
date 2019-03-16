package util

import (
	"log"
	"os"
	"strconv"
)

// MustGetEnv fatals if the env variable provided is not found
func MustGetEnv(env string) string {
	envValue := os.Getenv(env)
	if envValue == "" {
		log.Fatalf("$%s must be set", env)
	}
	return envValue
}

// MustGetEnvUInt32 gets the env and parses to uint32
func MustGetEnvUInt32(env string) uint32 {
	envValue := MustGetEnv(env)
	val, err := strconv.ParseUint(envValue, 10, 32)
	if err != nil {
		log.Fatalf("$%s is not set to a valid uint32, got: %s, err: %s", env, envValue, err)
	}
	return uint32(val)
}
