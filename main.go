package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    *Data  `json:"data"`
}

type ResponseWithOutData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseAddBook struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	BookId  *BookId `json:"data"`
}

type BookId struct {
	BookId string `json:"bookId"`
}

type Book struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Year       uint   `json:"year"`
	Author     string `json:"author"`
	Summary    string `json:"summary"`
	Publisher  string `json:"publisher"`
	PageCount  uint   `json:"pageCount,"`
	ReadPage   uint   `json:"readPage,"`
	Finished   bool   `json:"finished,"`
	Reading    bool   `json:"reading,"`
	InsertedAt string `json:"insertedAt,"`
	UpdatedAt  string `json:"updatedAt,"`
}

type MappedBook struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
}
type Data struct {
	Book []MappedBook `json:"books"`
}

type ClientError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var response Response
var books []Book

func getAllBooks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=utf-8")
	var mappedBook []MappedBook
	for _, item := range books {
		mappedBook = append(mappedBook, MappedBook{Id: item.Id, Name: item.Name, Publisher: item.Publisher})
	}
	response = Response{Status: "success", Message: " successfully get all books", Data: &Data{Book: mappedBook}}
	json.NewEncoder(res).Encode(response)
}

func addBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=utf-8")
	var book Book
	_ = json.NewDecoder(req.Body).Decode(&book)

	book.Id = strconv.Itoa(rand.Intn(1000000))
	book.InsertedAt = time.Now().Local().String()
	book.UpdatedAt = time.Now().Local().String()

	if book.ReadPage == book.PageCount {
		book.Finished = true
	} else {
		book.Finished = false
	}

	if book.Name == "" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ClientError{Status: "fail", Message: "Gagal menambahkan buku. Mohon isi nama buku"})

	} else if book.ReadPage > book.PageCount {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ClientError{Status: "fail", Message: "Gagal menambahkan buku. readPage tidak boleh lebih besar dari pageCount"})
	} else if !(book.Name == "" && book.ReadPage > book.PageCount) {
		books = append(books, book)
		res.WriteHeader(http.StatusCreated)
		for _, item := range books {
			if item.Id == book.Id {
				json.NewEncoder(res).Encode(ResponseAddBook{Status: "success", Message: "Buku berhasil ditambahkan", BookId: &BookId{book.Id}})
			}
		}
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ClientError{Status: "fail", Message: "gagal menambahkan buku"})
	}
}

func main() {
	routes := mux.NewRouter()

	routes.HandleFunc("/books", getAllBooks).Methods("GET")
	routes.HandleFunc("/books", addBook).Methods("POST")
	fmt.Printf("Starting server at port 8000\n")

	log.Fatal(http.ListenAndServe(":8000", routes))

}
