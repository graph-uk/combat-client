package combatClient

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (t *CombatClient) getSessionStatusJSON(sessionID string) (bool, string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("sessionID", sessionID)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(t.serverURL+"/getSessionStatus", contentType, bodyBuf)
	if err != nil {
		return false, err.Error(), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	finishedString := resp.Header.Get("Finished")
	finished := false
	if finishedString == "True" {
		finished = true
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Fail: incorrect request status - " + strconv.Itoa(resp.StatusCode))
		return false, "Incorrect request status: " + strconv.Itoa(resp.StatusCode), errors.New("Incorrect request status: " + strconv.Itoa(resp.StatusCode))
	} else {
		fmt.Println(string(body))
	}

	return finished, string(body), nil
}
