package main

import (
	"fmt"
	"strings"

	sops "github.com/getsops/sops/v3/decrypt"
	"github.com/pkg/errors"
)

func decryptSecretData(data map[string][]byte) (map[string][]byte, error) {
	var inputStrings []string
	for key, value := range data {
		inputStrings = append(inputStrings, fmt.Sprintf("%s=%s", key, string(value)))
	}

	input := strings.Join(inputStrings, "\n")
	output, err := sops.Data([]byte(input), "yaml")
	if err != nil {
		return nil, errors.Wrap(err, "error decrypting secret data")
	}

	newData := make(map[string][]byte)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		key := parts[0]
		value := parts[1]
		newData[key] = []byte(value)
	}

	return newData, nil
}
