package node

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func commandLine() string{
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n> ")
	userInput, _ := reader.ReadString('\n')

	return userInput
}

func clear() {
	cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()
}


