package combatClient

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func postSession(filename string, params string, targetUrl string) (string, error) {
	fileContent, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", err
	}

	content := base64.StdEncoding.EncodeToString(fileContent)

	json := fmt.Sprintf("{\"content\": \"%s\", \"arguments\":\"%s\"}", content, params)

	body := bytes.NewBuffer([]byte(json))

	resp, err := http.Post(targetUrl, "application/json", body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("Incorrect request status: " + strconv.Itoa(resp.StatusCode))
	}
	return string(responseBody), nil
}
