package authors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/nu7hatch/gouuid"
)

// INITIALIZATION

var PathToAuthors = getPathToAuthors()

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
	Id         string `json:"id"`
	Age        int    `json:"age"`
	Name       NameType
	Company    string `json:"company"`
	Email      string `json:"email"`
	Registered int64  `json:"registered"`
}

type NameType struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

// TYPE HELPERS

func (n *NameType) isValid() bool {
	var isValid = false
	if n.First != "" && n.Last != "" {
		isValid = true
	}
	return isValid
}

func (a *Author) ToString() string {
	bytes, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err.Error())
		return string("")
	}

	return string(bytes)
}

func (a *Author) isValid() bool {
	var isValid = false
	if a.Age != 0 && a.Name.isValid() && a.Company != "" && a.Email != "" {
		isValid = true
	}
	return isValid
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

// API

func (a *Author) Save() (*Author, error) {
	if !a.isValid() {
		return nil, errors.New("Not a valid author. Can`t save.")
	}

	blogAuthors, err := Get()
	if err != nil {
		return nil, err
	}

	authorArrIndex := getAuthorArrIndex(a.Id, blogAuthors)

	if authorArrIndex == -1 || a.Id == "" {
		u, uuidErr := uuid.NewV4()
		if uuidErr != nil {
			return nil, uuidErr
		}
		a.Id = u.String()
		a.Registered = time.Now().UnixNano() / 1000000
		*blogAuthors = append(*blogAuthors, *a)
	} else {
		(*blogAuthors)[authorArrIndex] = *a
	}

	if saveErr := saveAuthors(blogAuthors); saveErr != nil {
		return nil, saveErr
	}

	return a, nil
}

func (a *Author) Delete() error {
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

func GetById(Id string) (*Author, error) {
	blogAuthors, err := Get()
	if err != nil {
		return nil, err
	}
	authorArrIndex := getAuthorArrIndex(Id, blogAuthors)

	if authorArrIndex == -1 {
		return nil, errors.New("No author found. Can`t delete.")
	}

	return &(*blogAuthors)[authorArrIndex], nil
}

func Get() (*[]Author, error) {
	raw, readFileErr := ioutil.ReadFile(PathToAuthors)
	if readFileErr != nil {
		return nil, readFileErr
	}

	var blogAuthors []Author
	if parseErr := json.Unmarshal(raw, &blogAuthors); parseErr != nil {
		return nil, parseErr
	}

	return &blogAuthors, nil
}
