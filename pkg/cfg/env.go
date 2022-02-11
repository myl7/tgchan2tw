// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package cfg

import (
	"errors"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

func getEnvStr(k string) (string, error) {
	v := os.Getenv(k)
	if v != "" {
		return v, nil
	}
	return "", errors.New(fmt.Sprintf("env %s is not found", k))
}

func getEnvStrDef(k string, dv string) (string, error) {
	v := os.Getenv(k)
	if v != "" {
		return v, nil
	}
	return dv, nil
}

func getEnvIntDef(k string, dv int) (int, error) {
	v := os.Getenv(k)
	if v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("env %s is not an integer", k))
		}
		return i, nil
	}
	return dv, nil
}
