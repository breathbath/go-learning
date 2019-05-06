package main

import (
	"fmt"
	error2 "github.com/breathbath/go-learning/error"
	"github.com/breathbath/go_utils/utils/env"
	"github.com/breathbath/go_utils/utils/io"
)

var persons = []Person{
	{
		name: "Andrey",
		ageRate: 220,
	},
	{
		name: "Roman",
		ageRate: 230,
	},
	{
		name: "Ahmed",
		ageRate: 240,
	},
	{
		name: "",
		ageRate: 250,
	},
}

var dbIsDown = false

type Person struct {
	name string
	ageRate int
}

func main() {
	foundPerson, err := findPerson("Paul")

	switch e := err.(type) {
	case error2.WarningError:
		io.OutputWarning( "", e.Error())
	case error2.CriticalError:
		io.OutputError(e, "", "")
		panic(e)
	default:
		io.OutputError(e, "", "")
		panic(e)
	}
	if err != nil {
		io.OutputError(err, "", "")
		if env.ReadEnv("ENV", "") == "production" {
			sendResponseToUser("Something bad has happened to me")
			return
		}

		return
	}

	if foundPerson == nil {
		io.OutputInfo("", "Noone found\n")
	} else {
		io.OutputInfo("","Hello '%s', age rate: %d\n", foundPerson.name, foundPerson.ageRate)
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

	left := 1
	right := 2
	factor := 0
	result := addWithFactor(left, right, &factor)

	for i:=0; i<10; i++ {
		go addWithFactor(i*1, i*2, &factor)
	}
	io.OutputInfo("", "I had input %d and %d and factor %d and result is %d\n", left, right, factor, result)
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

	return nil,nil
}

func getAllPersonsFromDb() (dbPersons []Person, err error) {
	io.OutputInfo("","Getting users from db")
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

	io.OutputInfo("","Got %d users from db\n", len(persons))
	dbPersons = persons
	return
}

func addWithFactor(left, right int, factorPtr *int) int {
	internalFactor := *factorPtr

	if internalFactor == 0 {
		internalFactor = 1
		*factorPtr = internalFactor
	}

	return (left + right) * internalFactor
}

