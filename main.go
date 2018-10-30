package main

import (
	"fmt"
	"os"

	"github.com/graph-uk/combat-client/client"
)

func main() {

	client, err := combatClient.NewCombatClient()
	if err != nil {
		fmt.Println("Cannot init combat client")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	sessionID, err := client.CreateNewSession()
	if err != nil {
		fmt.Println("Cannot create combat session")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(client.GetSessionResult(sessionID))
}
