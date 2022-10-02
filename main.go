package main

import "fmt"

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    *Data  `json:"data"`
}

type Data struct {
	Book []Book `json: "books"`
}

type Book struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

var response Response
var books []Book

func main() {
	books = append(books, Book{Id: "1", Title: "book1", Publisher: "name"})
	response = Response{Status: "success", Message: " success", Data: &Data{Book{books}}}

	fmt.Println(response)

}
