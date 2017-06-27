package authors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"sync"

	"github.com/go-playground/validator"
	"github.com/nu7hatch/gouuid"
)

// INITIALIZATION

var PathToAuthors = getPathToAuthors()
var FileMutex = &sync.RWMutex{}
var validate = validator.New()

func getPathToAuthors() string {
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

func (a *Author) ToString() string {
	bytes, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err.Error())
		return string("")
	}

	return string(bytes)
}

func saveAuthors(a *[]Author) error {
	bytes, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(PathToAuthors, bytes, 0600)
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
	FileMutex.Lock()
	defer FileMutex.Unlock()

	blogAuthors, err := Get()
	if err != nil {
		return err
	}

	authorArrIndex := getAuthorArrIndex(a.Id, blogAuthors)

	if authorArrIndex == -1 {
		return errors.New("Author not found. Can`t save")
	}

	authorToUpdate := &(*blogAuthors)[authorArrIndex]
	softAssign(a, authorToUpdate)
	*a = *authorToUpdate

	if err := validate.Struct(authorToUpdate); err != nil {
		return err
	}
	return saveAuthors(blogAuthors)
}

func (a *Author) Delete() error {
	FileMutex.Lock()
	defer FileMutex.Unlock()

	blogAuthors, err := Get()
	if err != nil {
		return err
	}
	authorArrIndex := getAuthorArrIndex(a.Id, blogAuthors)

	if authorArrIndex == -1 {
		return errors.New("No author found. Can`t delete.")
	}

	*blogAuthors = append((*blogAuthors)[:authorArrIndex], (*blogAuthors)[authorArrIndex+1:]...)
	return saveAuthors(blogAuthors)
}

func (a *Author) Get() error {
	blogAuthors, err := Get()
	if err != nil {
		return err
	}
	authorArrIndex := getAuthorArrIndex(a.Id, blogAuthors)

	if authorArrIndex == -1 {
		return errors.New("No author found. Can`t delete.")
	}

	*a = (*blogAuthors)[authorArrIndex]
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

	FileMutex.Lock()
	defer FileMutex.Unlock()

	blogAuthors, err := Get()
	if err != nil {
		return err
	}

	*blogAuthors = append(*blogAuthors, *a)
	return saveAuthors(blogAuthors)
}

func Get() (*[]Author, error) {
	raw, readFileErr := ioutil.ReadFile(PathToAuthors)
	if readFileErr != nil {
		return nil, errors.New("Can`t read authors file")
	}

	var blogAuthors []Author
	if parseErr := json.Unmarshal(raw, &blogAuthors); parseErr != nil {
		return nil, parseErr
	}

	return &blogAuthors, nil
}
