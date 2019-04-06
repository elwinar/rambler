package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)


var reg *regexp.Regexp

func init() {
	reg2, err := regexp.Compile(`\$\{(.+?)\}`)
	if err != nil {
		panic("couldn't compile reg ex")
	}
	reg = reg2
}

func findEntries(statement string) []string {
	return reg.FindAllString(statement, -1)
}


func findEnvVal(toReplace string) (string, error) {
	key := strings.TrimSuffix(toReplace, "}")
	key = strings.TrimPrefix(key, "${")
	value := os.Getenv(strings.ToUpper(key))

	if value == "" {
		return "", fmt.Errorf("could not find env value for %s", toReplace)
	}

	return value, nil
}

func replace(statement string) (string, error) {
	r := findEntries(statement)
	if len(r) == 0 {
		return statement, nil
	}

	var replacements []string
	for _, tr := range r {
		val, err := findEnvVal(tr)
		if err != nil {
			return "", err
		}

		replacements = append(replacements, val)
	}

	for i, val := range replacements {
		key := r[i]
		statement = strings.Replace(statement, key, val, 1)
	}
	return statement, nil
}