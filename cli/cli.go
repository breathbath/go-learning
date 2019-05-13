package cli

import (
	"fmt"

	"github.com/breathbath/go-learning/person/person"
)

func Run() error {
	person.FindPerson("Andrey")
	return
}

func sendResponseToUser(output string) {
	fmt.Println(output)
}
