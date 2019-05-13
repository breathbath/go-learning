package main

import (
	"fmt"
	error2 "github.com/breathbath/go-learning/error"
	"github.com/breathbath/go_utils/utils/env"
	"github.com/breathbath/go_utils/utils/io"
	"sync"
	"cli/cli"
)

func main() {
	err = cli.Run()
	//foundPerson, err := findPerson("Paul")

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
