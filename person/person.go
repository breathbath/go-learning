package person

import (
	"fmt"

	error2 "github.com/breathbath/go-learning/error"
	"github.com/breathbath/go_utils/utils/io"
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

func FindPerson(name string) (person *Person, err error) {
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
