package main

import (
	"fmt"
	error2 "github.com/breathbath/go-learning/error"
	"github.com/breathbath/go_utils/utils/env"
	"github.com/breathbath/go_utils/utils/io"
	"sync"
)

var persons = []Person{
	{
		name:    "Andrey",
		ageRate: 220,
	},
	{
		name:    "Roman",
		ageRate: 230,
	},
	{
		name:    "Ahmed",
		ageRate: 240,
	},
	{
		name:    "",
		ageRate: 250,
	},
}

var dbIsDown = false

type Person struct {
	name    string
	ageRate int
}

func main() {
	foundPerson, err := findPerson("Paul")

	if err != nil {
		switch e := err.(type) {
		case error2.WarningError:
			io.OutputWarning("", e.Error())
		case error2.CriticalError:
			io.OutputError(e, "", "")
			if env.ReadEnv("ENV", "") == "production" {
				sendResponseToUser("Something bad has happened to me")
				return
			}
			panic(e)
		default:
			io.OutputError(e, "", "")
			panic(e)
		}
	}

	if foundPerson == nil {
		io.OutputInfo("", "Noone found\n")
	} else {
		io.OutputInfo("", "Hello '%s', age rate: %d\n", foundPerson.name, foundPerson.ageRate)
	}

	//
	//foundPerson = findPerson(persons, "Roman")
	//if foundPerson == nil {
	//	fmt.Printf("Noone found\n")
	//} else {
	//	fmt.Printf("Hello '%s', age: %d\n", foundPerson.name, foundPerson.age)
	//}
	//
	//foundPerson = findPerson(persons, "Ahmed")
	//if foundPerson == nil {
	//	fmt.Printf("Noone found\n")
	//} else {
	//	fmt.Printf("Hello '%s', age: %d\n", foundPerson.name, foundPerson.age)
	//}
	//
	//foundPerson = findPerson(persons, "Paul")
	//if foundPerson == nil {
	//	fmt.Printf("Noone found\n")
	//} else {
	//	fmt.Printf("Hello '%s', age: %d\n", foundPerson.name, foundPerson.age)
	//}

	wg := sync.WaitGroup{}
	factor := 0
	for i := 0; i < 10; i++ {
		y := i
		wg.Add(1)
		go func() {
			addWithFactor(y, y + 1, &factor)
			wg.Done()
		}()
	}
	wg.Wait()
}

func sendResponseToUser(output string) {
	fmt.Println(output)
}

func findPerson(name string) (person *Person, err error) {
	io.OutputInfo("", "Trying to find person by name '%s'\n", name)
	var dbPersons []Person

	dbPersons, err = getAllPersonsFromDb()
	if err != nil {
		return
	}

	for _, curPerson := range dbPersons {
		if name == curPerson.name {
			return &curPerson, nil
		}
	}

	return nil, nil
}

func getAllPersonsFromDb() (dbPersons []Person, err error) {
	io.OutputInfo("", "Getting users from db")
	if len(persons) == 0 {
		err = error2.WarningError{
			error2.NewErrorWrapper(fmt.Errorf("Database is empty, it has [%d] entries", len(persons))),
		}
		return
	}

	if dbIsDown {
		err = error2.CriticalError{
			error2.NewErrorWrapper(fmt.Errorf("Database is down")),
		}
		return
	}

	io.OutputInfo("", "Got %d users from db\n", len(persons))
	dbPersons = persons
	return
}

func addWithFactor(left, right int, factorPtr *int) int {
	internalFactor := *factorPtr

	if internalFactor == 0 {
		internalFactor = left+right
		*factorPtr = internalFactor
	}

	res := (left + right) * internalFactor
	io.OutputInfo("", "I had input %d and %d and factor %d and result is %d\n", left, right, internalFactor, res)
	return res
}
