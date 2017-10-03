package combatClient

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

type CombatClient struct {
	serverURL             string
	sessionID             string
	sessionBeginTimestamp time.Time
	lastSTDOutMessage     string
	SessionTimeout        time.Duration
}

func (t *CombatClient) getServerUrlFromCLI() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("Server URL is required")
	}
	return os.Args[1], nil
}

func NewCombatClient() (*CombatClient, error) {
	var result CombatClient
	var err error

	result.serverURL, err = result.getServerUrlFromCLI()
	if err != nil {
		return &result, err
	}
	result.lastSTDOutMessage = ""
	return &result, nil
}

func (t *CombatClient) packTests() (string, error) {
	fmt.Println("Packing tests")
	tmpFile, err := ioutil.TempFile("", "combatSession")
	if err != nil {
		panic(err)
	}
	tmpFile.Close()
	zipit("./../..", tmpFile.Name())
	return tmpFile.Name(), nil
}

func (t *CombatClient) cleanupTests() error {
	fmt.Println("Cleanup tests")
	return nil
}

func (t *CombatClient) getParams() string {
	params := ""
	for curArgIndex, curArg := range os.Args {
		if curArgIndex > 1 {
			params += curArg + " "
		}
	}
	return params
}

func (t *CombatClient) createSessionOnServer(archiveFileName string) string {
	fmt.Print("Uploading session")
	sessionName := ""

	var err error
	for true { // endless cycle for try upload tests.
		sessionName, err = postSession(archiveFileName, t.getParams(), t.serverURL+"/createSession")
		if err != nil {
			fmt.Print(`.`)
			time.Sleep(5 * time.Second)
		} else {
			fmt.Println()
			break
		}
	}
	return sessionName
}

func (t *CombatClient) CreateNewSession(timeoutMinutes int) (string, error) {
	t.sessionBeginTimestamp = time.Now()
	t.SessionTimeout = time.Minute * time.Duration(timeoutMinutes)
	err := t.cleanupTests()
	if err != nil {
		fmt.Println("Cannot cleanup tests")
		return "", err
	}

	testsArchiveFileName, err := t.packTests()
	if err != nil {
		fmt.Println("Cannot pack tests to zip archive")
		return "", err
	}
	//os.Exit(0)
	sessionName := t.createSessionOnServer(testsArchiveFileName)
	//fmt.Println("Session: " + sessionName)
	combatServerURL, err := t.getServerUrlFromCLI()
	if err != nil {
		fmt.Println("Cannot parse server name as parameter")
		return "", err
	}
	fmt.Println("Session status: " + combatServerURL + "/sessions/" + sessionName)
	t.sessionID = sessionName
	return sessionName, nil
}

func (t *CombatClient) GetSessionResult(sessionID string) int {
	countOfErrors := 1
	for {
		sessionStatusJSON, err := t.getSessionStatusJSON(sessionID)
		if err != nil {
			fmt.Println(err.Error())
		}
		var finished bool
		finished, countOfErrors, err = t.printSessionStatusByJSON(sessionStatusJSON)
		if err == nil {
			if finished {
				break
			}
		}

		if time.Since(t.sessionBeginTimestamp) > t.SessionTimeout {
			fmt.Println(``)
			fmt.Println(`Timeout was reached, but session is still not finished. Check workers and start new session.`)
			os.Exit(1)
		}
		time.Sleep(5 * time.Second)
	}

	// cut microseconds
	timeLongStr := time.Since(t.sessionBeginTimestamp).String()
	r := regexp.MustCompile(`\.\d*s$`)
	timeShortStr := r.ReplaceAllString(timeLongStr, "s")

	fmt.Println("Time of testing: " + timeShortStr)
	return countOfErrors
}
