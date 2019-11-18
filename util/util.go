package util

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

// MustGetEnv fatals if the env variable provided is not found
func MustGetEnv(env string) string {
	envValue := os.Getenv(env)
	if envValue == "" {
		logrus.WithField("env", env).Fatal("MustGetEnv did not find env")
	}
	return envValue
}

// MustGetEnvUInt32 gets the env and parses to uint32
func MustGetEnvUInt32(env string) uint32 {
	return envParseUInt32(env, MustGetEnv(env))
}

// GetEnvOrDefault gets an environment variable or uses the default value
func GetEnvOrDefault(env, defaultValue string) string {
	envValue := os.Getenv(env)
	if envValue == "" {
		return defaultValue
	}
	return envValue
}

// GetEnvOrDefaultUInt32 gets an environment variable (or default) and parses to uint32
func GetEnvOrDefaultUInt32(env, defaultValue string) uint32 {
	return envParseUInt32(env, GetEnvOrDefault(env, defaultValue))
}

func envParseUInt32(env, envValue string) uint32 {
	val, err := strconv.ParseUint(envValue, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{"env": env, "value": envValue}).Fatal("envParseUInt32 did not find uint32 value")
	}
	return uint32(val)
}
