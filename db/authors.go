package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"sync"

	"github.com/nu7hatch/gouuid"
)

// INITIALIZATION

var pathToAuthors = getpathToAuthors()
var authorsFileMutex = &sync.RWMutex{}

func getpathToAuthors() string {
	path, filenameErr := filepath.Abs("./db/db_data/authors.json")
	if filenameErr != nil {
		fmt.Println(filenameErr.Error())
		os.Exit(1)
	}
	return path
}

// TYPES

type Author struct {
	Id         string   `json:"id" validate:"required"`
	Age        int      `json:"age" validate:"required"`
	Name       NameType `json:"name" validate:"required,dive,required"`
	Company    string   `json:"company" validate:"required"`
	Email      string   `json:"email" validate:"required,email"`
	Registered int64    `json:"registered" validate:"required"`
}

type NameType struct {
	First string `json:"first" validate:"required"`
	Last  string `json:"last" validate:"required"`
}

// TYPE HELPERS

func readAuthors() (*[]Author, error) {
	raw, readFileErr := ioutil.ReadFile(pathToAuthors)
	if readFileErr != nil {
		fmt.Println(readFileErr)
		return nil, errors.New("Can`t read authors file")
	}

	var authors []Author
	if parseErr := json.Unmarshal(raw, &authors); parseErr != nil {
		return nil, parseErr
	}

	return &authors, nil
}

func writeAuthors(a *[]Author) error {
	bytes, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pathToAuthors, bytes, 0600)
}

func getAuthorArrIndex(authorId string, authors *[]Author) int {
	result := -1
	for index, author := range *authors {
		if author.Id == authorId {
			result = index
			break
		}
	}

	return result
}

func softAssign(from *Author, to *Author) {
	if from.Age != 0 {
		to.Age = from.Age
	}
	if from.Name.First != "" {
		to.Name.First = from.Name.First
	}
	if from.Name.Last != "" {
		to.Name.Last = from.Name.Last
	}
	if from.Email != "" {
		to.Email = from.Email
	}
	if from.Company != "" {
		to.Company = from.Company
	}
}

// API

func (a *Author) Update() error {
	authorsFileMutex.Lock()
	defer authorsFileMutex.Unlock()

	authors, err := readAuthors()
	if err != nil {
		return err
	}

	authorArrIndex := getAuthorArrIndex(a.Id, authors)

	if authorArrIndex == -1 {
		return errors.New("Author not found. Can`t save")
	}

	authorToUpdate := &(*authors)[authorArrIndex]
	softAssign(a, authorToUpdate)
	*a = *authorToUpdate

	if err := validate.Struct(authorToUpdate); err != nil {
		return err
	}
	return writeAuthors(authors)
}

func (a *Author) Delete() error {
	authorsFileMutex.Lock()
	defer authorsFileMutex.Unlock()

	authors, err := readAuthors()
	if err != nil {
		return err
	}
	authorArrIndex := getAuthorArrIndex(a.Id, authors)

	if authorArrIndex == -1 {
		return errors.New("No author found. Can`t delete.")
	}

	*authors = append((*authors)[:authorArrIndex], (*authors)[authorArrIndex+1:]...)
	return writeAuthors(authors)
}

func (a *Author) Get() error {
	authorsFileMutex.Lock()
	defer authorsFileMutex.Unlock()
	authors, err := readAuthors()
	if err != nil {
		return err
	}
	authorArrIndex := getAuthorArrIndex(a.Id, authors)

	if authorArrIndex == -1 {
		return errors.New("No author found. Can`t delete.")
	}

	*a = (*authors)[authorArrIndex]
	return nil
}

func NewAuthor(a *Author) error {
	u, uuidErr := uuid.NewV4()
	if uuidErr != nil {
		return uuidErr
	}

	a.Id = u.String()
	a.Registered = time.Now().UnixNano() / 1000000

	if err := validate.Struct(a); err != nil {
		return err
	}

	authorsFileMutex.Lock()
	defer authorsFileMutex.Unlock()

	authors, err := readAuthors()
	if err != nil {
		return err
	}

	*authors = append(*authors, *a)
	return writeAuthors(authors)
}

func GetAuthors() (*[]Author, error) {
	authorsFileMutex.RLock()
	defer authorsFileMutex.RUnlock()
	return readAuthors()
}
