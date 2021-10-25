package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func parseResponseBodyJSON(body io.ReadCloser, v interface{}) error {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("unable to read response body: %v", err)
	}

	if jsonErr := json.Unmarshal(byt, v); jsonErr != nil {
		return fmt.Errorf("unable to parse response body to json: %s\n", jsonErr)
	}

	return nil
}

func jsonPrint(v interface{}) error {
	byt, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(byt))
	return nil
}
