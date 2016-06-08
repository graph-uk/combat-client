package main

import (
	"fmt"
	"os"

	"github.com/graph-uk/combat-client/client"
)

func main() {
	defaultSessionTimeout := 60 //minutes

	client, err := combatClient.NewCombatClient()
	if err != nil {
		fmt.Println("Cannot init combat client")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	sessionID, err := client.CreateNewSession(defaultSessionTimeout)
	if err != nil {
		fmt.Println("Cannot create combat session")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	client.GetSessionResult(sessionID)
}
