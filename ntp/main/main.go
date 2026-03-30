package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sparxfort1ano/wb-level-2/ntp/ntpclient"
)

func main() {
	ntpTime, err := ntpclient.GetCurrentTime()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println("Current time:", ntpTime)
	fmt.Println("Local time:  ", time.Now())
}
