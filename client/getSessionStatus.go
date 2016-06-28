package combatClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

type SessionStatus struct {
	Finished           bool
	TotalCasesCount    int
	FinishedCasesCount int
	FailReports        []string
}

func (t *CombatClient) getSessionStatusJSON(sessionID string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("sessionID", sessionID)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(t.serverURL+"/getSessionStatus", contentType, bodyBuf)
	if err != nil {
		return err.Error(), err
	}
	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "Incorrect request status: " + strconv.Itoa(resp.StatusCode), errors.New("Incorrect request status: " + strconv.Itoa(resp.StatusCode))
	}
	return string(body), nil
}

func (t *CombatClient) printSessionStatusByJSON(sessionStatusJSON string) (bool, int, error) {
	var sessionStatus SessionStatus
	err := json.Unmarshal([]byte(sessionStatusJSON), &sessionStatus)
	if err != nil {
		fmt.Println("Cannot parse session status JSON")
		fmt.Println(err.Error())
		fmt.Println(sessionStatusJSON)
		return false, 1, err
	}
	if !sessionStatus.Finished {
		if sessionStatus.TotalCasesCount == 0 {
			fmt.Println("Cases exploring")
		} else {
			fmt.Print("Testing (" + strconv.Itoa(sessionStatus.FinishedCasesCount) + "/" + strconv.Itoa(sessionStatus.TotalCasesCount) + ")")

			if len(sessionStatus.FailReports) != 0 {
				fmt.Print(" " + strconv.Itoa(len(sessionStatus.FailReports)) + " errors:")
			}
			fmt.Println()

			for _, curFail := range sessionStatus.FailReports {
				fmt.Println("    " + curFail)
			}
			if len(sessionStatus.FailReports) != 0 {
				fmt.Println()
			}
		}
	} else { // if session finished
		if len(sessionStatus.FailReports) == 0 { // if no errors
			fmt.Println("Finished success")
		} else { // if errors found
			fmt.Println("Finished with " + strconv.Itoa(len(sessionStatus.FailReports)) + " errors:")
			for _, curFail := range sessionStatus.FailReports {
				fmt.Println("    " + curFail)
			}
			fmt.Println("More info at: "+t.serverURL+"/sessions/"+t.sessionID)
		}
	}
	return sessionStatus.Finished, len(sessionStatus.FailReports), nil
}
