package main

import (
	"context"
	"fmt"
	"os"

	"github.com/beeper/libgmessages/client"
)

func main() {
	client := client.New()

	if !client.Authenticated() {
		ch, err := client.Pair(context.Background())
		if err != nil {
			fmt.Println("Failed to created client")
			os.Exit(1)
		}

		err = client.Connect()
		if err != nil {
			fmt.Println("Failed to connect:", err)
			os.Exit(1)
		}

		fmt.Println("ch", ch)

		fmt.Println("item", <-ch)
	} else {
		err := client.Connect()
		if err != nil {
			fmt.Println("Failed to connect:", err)
			os.Exit(1)
		}
	}
}
