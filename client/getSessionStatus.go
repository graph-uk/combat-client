package combatClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type SessionStatus struct {
	ID                  string
	Status              string
	SessionError        string
	CasesCount          int
	CasesProcessedCount int
	CasesFailed         []string
}

func (t *CombatClient) getSessionStatusJSON(sessionID string) (string, error) {
	resp, err := http.Get(t.serverURL + "/api/v1/sessions")
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

	casesFailed := []string{}
	if sessionStatus.CasesFailed != nil {
		casesFailed = sessionStatus.CasesFailed
	}

	msg = fmt.Sprintf("%s - %s\nProcessed %d of %d tests", sessionStatus.ID, sessionStatus.Status, sessionStatus.CasesProcessedCount, sessionStatus.CasesCount)

	if sessionStatus.Status == "Failed" || sessionStatus.Status == "Success" || sessionStatus.Status == "Incomplete" {
		fmt.Println()
		fmt.Print(msg)
		if len(casesFailed) > 0 {
			fmt.Println("Failed cases:")
			for _, caseFailed := range casesFailed {
				fmt.Println(caseFailed)
			}
		}

		return true, len(casesFailed), nil
	}

	if sessionStatus.CasesCount == 0 {
		msg = fmt.Sprintf("%s - %s\nCase exploring", sessionStatus.ID, sessionStatus.Status)
	}

	if msg != t.lastSTDOutMessage {
		t.lastSTDOutMessage = msg
		fmt.Println()
		fmt.Print(msg)
	} else {
		fmt.Print(".")
	}

	return false, len(casesFailed), nil
}
