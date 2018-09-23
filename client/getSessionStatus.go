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
	Finished                  bool
	TotalCasesCount           int
	FinishedCasesCount        int
	CasesExploringFailMessage string
	FailReports               []string
}

func (t *CombatClient) getSessionStatusJSON(sessionID string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("sessionID", sessionID)
	bodyWriter.Close()

	resp, err := http.Post(t.serverURL+"/api/v1/sessions", "application/json", bodyBuf)
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
	msg := ""
	var sessionStatus SessionStatus
	err := json.Unmarshal([]byte(sessionStatusJSON), &sessionStatus)
	if err != nil {
		msg += ("Cannot parse session status JSON\r\n")
		msg += (err.Error() + "\r\n")
		msg += (sessionStatusJSON)
		if t.lastSTDOutMessage != msg {
			t.lastSTDOutMessage = msg
			fmt.Println()
			fmt.Print(msg)
		} else {
			fmt.Print(`.`)
		}
		return false, 1, err
	}
	if !sessionStatus.Finished {
		if sessionStatus.TotalCasesCount == 0 {
			msg += ("Cases exploring")
		} else {
			msg += ("Testing (" + strconv.Itoa(sessionStatus.FinishedCasesCount) + "/" + strconv.Itoa(sessionStatus.TotalCasesCount) + ")")

			if len(sessionStatus.FailReports) != 0 {
				msg += (" " + strconv.Itoa(len(sessionStatus.FailReports)) + " errors:")
				msg += "\r\n"
			}

			for _, curFail := range sessionStatus.FailReports {
				msg += ("    " + curFail + "\r\n")
			}
		}
	} else { // if session finished
		if sessionStatus.CasesExploringFailMessage == "" {
			if len(sessionStatus.FailReports) == 0 { // if no errors
				msg += ("Finished success\r\n")
			} else { // if errors found
				msg += "Finished with " + strconv.Itoa(len(sessionStatus.FailReports)) + " errors:\r\n"
				for _, curFail := range sessionStatus.FailReports {
					msg += ("    " + curFail + "\r\n")
				}
				msg += ("More info at: " + t.serverURL + "/sessions/" + t.sessionID + "\r\n")
			}
		} else { // if session finished on failed cases exploring
			msg += ("Cases exploring failed. Combat says: \r\n" + sessionStatus.CasesExploringFailMessage)
			fmt.Println()
			fmt.Print(msg)
			return sessionStatus.Finished, 1, nil
		}
	}
	if t.lastSTDOutMessage != msg {
		t.lastSTDOutMessage = msg
		fmt.Println()
		fmt.Print(msg)
	} else {
		fmt.Print(`.`)
	}
	return sessionStatus.Finished, len(sessionStatus.FailReports), nil
}
