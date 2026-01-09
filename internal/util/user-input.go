package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Returns 1 for yes, 0 for no. Defaults to yes if user does not enter anything
func AskYesNoQuestionDefaultYes(question string) int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	input, err := reader.ReadString('\n')

	if err != nil {
			// default to "y" on read error
			input = "y"
	}
	
	response := strings.TrimSpace(input)
	if response == "" {
			response = "y"
	}

	if strings.ToLower(response) != "y" {
			return 0
	}
	return 1
}

func WaitForResponse(question string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}