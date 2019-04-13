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
	envValue := MustGetEnv(env)
	val, err := strconv.ParseUint(envValue, 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{"env": env, "value": envValue}).Fatal("MustGetEnvUInt32 did not find uint32 value")
	}
	return uint32(val)
}
