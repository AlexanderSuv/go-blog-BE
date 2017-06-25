package main

import (
	"github.com/AlexanderSuv/goblog/server"
)

func main() {
	//myauthors, err := authors.Get()
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(len(*myauthors))
	//
	//newAuthor := &authors.Author{
	//	Name: authors.NameType{
	//		First: "Alex",
	//		Last:  "Suvorov",
	//	},
	//	Age:     28,
	//	Company: "BMW",
	//	Email:   "alex.suv@gmail.com",
	//}
	//
	//newAuthor, _ = newAuthor.Save()
	//fmt.Printf("saved author with ID: %s\n", newAuthor.Id)
	//
	//myauthors, err = authors.Get()
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(len(*myauthors))
	//
	//author, _ := authors.GetById(newAuthor.Id)
	//fmt.Printf("found author by ID: %s\n", author.ToString())
	//
	//author.Delete()
	//fmt.Printf("Author deleted")
	//
	//myauthors, err = authors.Get()
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(len(*myauthors))

	server.Start()
}
