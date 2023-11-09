package main

import (
	"fmt"
	"os"
	"os/exec"
)

func commandLine() string{
	var command string
	fmt.Print("\n> ")
	fmt.Scanf("%s", &command)

	return command
}

func clear() {
	cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()
}


