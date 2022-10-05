package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
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

type ResponseGetBookById struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    *SingleBookData `json:"data"`
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

type SingleBookData struct {
	Book Book `json:"book"`
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
	query := req.URL.Query()

	if query.Get("reading") != "" {
		if query.Get("reading") == "1" {
			for _, item := range books {
				if item.Reading {
					mappedBook = append(mappedBook, MappedBook{Id: item.Id, Name: item.Name, Publisher: item.Publisher})
				}
			}
		} else if query.Get("reading") == "0" {
			for _, item := range books {
				if !item.Reading {
					mappedBook = append(mappedBook, MappedBook{Id: item.Id, Name: item.Name, Publisher: item.Publisher})
				}
			}
		}
	} else if query.Get("finished") != "" {
		if query.Get("finished") == "1" {
			for _, item := range books {
				if item.Finished {
					mappedBook = append(mappedBook, MappedBook{Id: item.Id, Name: item.Name, Publisher: item.Publisher})
				}
			}
		} else if query.Get("finished") == "0" {
			for _, item := range books {
				if !item.Finished {
					mappedBook = append(mappedBook, MappedBook{Id: item.Id, Name: item.Name, Publisher: item.Publisher})
				}
			}
		}

	} else if query.Get("name") != "" {
		for _, item := range books {
			if strings.Contains(strings.ToLower(item.Name), strings.ToLower(query.Get("name"))) {
				mappedBook = append(mappedBook, MappedBook{Id: item.Id, Name: item.Name, Publisher: item.Publisher})
			}
		}
	} else {
		for _, item := range books {
			mappedBook = append(mappedBook, MappedBook{Id: item.Id, Name: item.Name, Publisher: item.Publisher})
		}
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
		return

	} else if book.ReadPage > book.PageCount {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ClientError{Status: "fail", Message: "Gagal menambahkan buku. readPage tidak boleh lebih besar dari pageCount"})
		return
	} else if !(book.Name == "" && book.ReadPage > book.PageCount) {
		books = append(books, book)
		res.WriteHeader(http.StatusCreated)
		for _, item := range books {
			if item.Id == book.Id {
				json.NewEncoder(res).Encode(ResponseAddBook{Status: "success", Message: "Buku berhasil ditambahkan", BookId: &BookId{book.Id}})
				break
			}
		}
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ClientError{Status: "fail", Message: "gagal menambahkan buku"})
		return
	}
}

func getBookById(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=utf-8")
	params := mux.Vars(req)

	var book Book
	for _, item := range books {
		if item.Id == params["id"] {
			book = item
			break
		}
	}
	if book.Id != "" {
		json.NewEncoder(res).Encode(ResponseGetBookById{Status: "success", Message: "successfully get book", Data: &SingleBookData{Book: book}})
		return
	} else {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(ResponseWithOutData{Status: "fail", Message: "Buku tidak ditemukan"})
		return

	}
}

func updateBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=utf-8")
	params := mux.Vars(req)

	var book Book
	_ = json.NewDecoder(req.Body).Decode(&book)

	if book.Name == "" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ClientError{Status: "fail", Message: "Gagal memperbarui buku. Mohon isi nama buku"})
		return

	} else if book.ReadPage > book.PageCount {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ClientError{Status: "fail", Message: "Gagal memperbarui buku. readPage tidak boleh lebih besar dari pageCount"})
		return
	} else if !(book.Name == "" && book.ReadPage > book.PageCount) {
		for index, item := range books {
			if item.Id == params["id"] {
				books = append(books[:index], books[index+1:]...)

				_ = json.NewDecoder(req.Body).Decode(&book)
				book.Id = params["id"]
				book.InsertedAt = item.InsertedAt
				book.UpdatedAt = time.Now().Local().String()
				books = append(books, book)
			}
		}

		if book.Id != "" {
			json.NewEncoder(res).Encode(ResponseGetBookById{Status: "success", Message: "Buku berhasil diperbarui", Data: &SingleBookData{book}})
			return
		} else {
			res.WriteHeader(http.StatusNotFound)
			json.NewEncoder(res).Encode(ResponseWithOutData{Status: "fail", Message: "Gagal memperbarui buku. Id tidak ditemukan"})
			return
		}
	}
}

func deleteBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=utf-8")
	params := mux.Vars(req)

	var isDelete bool = false

	for index, item := range books {
		if item.Id == params["id"] {
			books = append(books[:index], books[index+1:]...)
			isDelete = true
			break
		}
	}

	if isDelete {
		json.NewEncoder(res).Encode(ResponseWithOutData{Status: "success", Message: "Buku berhasil dihapus"})
		return
	} else {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(ResponseWithOutData{Status: "fail", Message: "Buku gagal dihapus. Id tidak ditemukan"})
		return
	}

}

func main() {
	routes := mux.NewRouter()
	routes.HandleFunc("/books", getAllBooks).Methods("GET")
	routes.HandleFunc("/books", addBook).Methods("POST")
	routes.HandleFunc("/books/{id}", getBookById).Methods("GET")
	routes.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	routes.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")

	log.Fatal(http.ListenAndServe(":5000", routes))

}
