// https://docs.datadoghq.com/agent/guide/secrets-management/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// secretsPayload is stdin from datadog agent
// sample
// '{ "version": "1.0", "secrets": ["secret1", "secret2"] }'
type secretsPayload struct {
	Secrets []string `json:"secrets"`
	Version int      `json:"version"`
}

type secretResponse struct {
	Secrets []secretValue
}

// The expected payload is a JSON object, where each key is one of the handles requested in the input payload.
// The value for each handle is a JSON object with 2 fields:
//
// value: a string; the actual secret value to be used in the check configurations (can be null in the case of error).
// error: a string; the error message, if needed. If error is anything other than null, the integration configuration that uses this handle is considered erroneous and is dropped.
// learn more: https://docs.datadoghq.com/agent/guide/secrets-management/
type secretValue struct {
	SecretKey string  `json:"-"`
	Value     *string `json:"value"` // can be null in the case of error
	Error     *string `json:"error"` // If error is not happened must be nil.
}

func (r secretResponse) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	for _, secret := range r.Secrets {
		data[secret.SecretKey] = secret
	}

	return json.Marshal(data)
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read from stdin: %s", err)
		os.Exit(1)
	}

	secrets := secretsPayload{}
	json.Unmarshal(data, &secrets)

	res := secretResponse{}

	for _, secretKey := range secrets.Secrets {
		val := secretValue{
			SecretKey: secretKey,
			Value:     nil,
			Error:     nil,
		}

		if value, ok := os.LookupEnv(strings.ToUpper(secretKey)); ok {
			val.Value = &value
		} else {
			errStr := fmt.Sprintf("environment variable [%s] is not set", strings.ToUpper(secretKey))
			val.Error = &errStr
		}
		res.Secrets = append(res.Secrets, val)
	}

	output, err := json.Marshal(res)

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not serialize res: %s", err)
		os.Exit(1)
	}

	fmt.Printf(string(output))
}
