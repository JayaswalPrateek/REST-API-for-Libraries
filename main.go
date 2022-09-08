package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	switch homePage() {
	case 1:
		handleBooks()
	case 2:
		handleMembers()
	case 3:
		handleBookStatus()
	default:
		unexpected()
	}
}

type book struct {
	ISBN   string `json:"isbn"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Status string `json:"status"`
}

var books = []book{
	{ISBN: "0001", Title: "RANDOM1", Author: "TONY", Status: "issued"},
	{ISBN: "0002", Title: "RANDOM2", Author: "PETER", Status: "available"},
	{ISBN: "0003", Title: "RANDOM3", Author: "SAM", Status: "available"},
	{ISBN: "0004", Title: "RANDOM4", Author: "JOE", Status: "issued"},
}

type member struct {
	UID        string `json:"uid"`
	Name       string `json:"name"`
	DateJoined string `json:"datejoined"`
	ContactNum string `json:"contactnum"`
	Status     string `json:"status"`
	History    string `json:"history"`
}

var members = []member{
	{UID: "0001", Name: "ME", DateJoined: "04082004", ContactNum: "9136125577", Status: "BOOK1", History: "BOOK2,BOOK3,BOOK4"},
	{UID: "0002", Name: "HE", DateJoined: "14082010", ContactNum: "9219922580", Status: "BOOK2", History: "BOOK1,BOOK3,BOOK4"},
	{UID: "0003", Name: "SHE", DateJoined: "24092014", ContactNum: "9919922911", Status: "BOOK3", History: "BOOK2,BOOK1,BOOK4"},
	{UID: "0004", Name: "Teacher", DateJoined: "05102006", ContactNum: "8724944789", Status: "BOOK4", History: "BOOK2,BOOK3,BOOK1"},
}

func clearScr() {
	fmt.Print("\x1bc")
}
func exitChoice() {
	clearScr()
	fmt.Println("See You Later, Keep Reading, Keep Learning :)")
	os.Exit(0)
}
func unexpected() {
	fmt.Print("You weren't supposed to see this message, You have encountered a bug, Please contact us @xx-xxxxxxxxxx")
}
func checkErr(err error, message string) {
	if err != nil {
		fmt.Println(message)
		time.Sleep(10 * time.Second)
	}
}

func homePage() int {
	var choice int
	fmt.Println("Enter Your Choice:")
	fmt.Println("[1] Books")
	fmt.Println("[2] Members")
	fmt.Println("[3] Manage")
	fmt.Println("[0] Exit")
	fmt.Scanln(&choice)
	if choice == 0 {
		exitChoice()
	}
	if !(choice == 1 || choice == 2 || choice == 3) {
		clearScr()
		fmt.Print("Invalid Input, Please Reconsider\n\n")
		homePage()
	}
	clearScr()
	return choice
}

func handleBooks() {
	var choice int
	fmt.Println("Enter Your Choice: ")
	fmt.Println("[1] Get all books in the database")
	fmt.Println("[2] Get books available for issue")
	fmt.Println("[3] Get all issued books")
	fmt.Println("[4] Get book details by ISBN number")
	fmt.Println("[5] Browse books using book title")
	fmt.Println("[6] Browse books using author name")
	fmt.Println("[0] Exit")
	fmt.Scanln(&choice)
	if choice == 0 {
		exitChoice()
	}
	if !(choice == 1 || choice == 2 || choice == 3 || choice == 4 || choice == 5 || choice == 6) {
		clearScr()
		fmt.Print("Invalid Input, Please Reconsider\n\n")
		handleBooks()
	}
	clearScr()
	router := gin.Default()
	switch choice {
	case 1:
		router.GET("/books", getAllBooks)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to getAllBooks()")
	case 2:
		router.GET("/available-books", getAvailableBooks)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to getAvailableBooks()")
	case 3:
		router.GET("/issued-books", getIssuedBooks)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to getIssuedBooks()")
	case 4:
		router.GET("/books/:isbn", handleLookupBooksByISBN)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to handlelookupBooksByISBN()")
	case 5:
		router.GET("/books/:title", handlelookupBooksByTitle)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to handlelookupBooksByTitle()")
	case 6:
		router.GET("/books/:author", handlelookupBooksByAuthor)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to handlelookupBooksByAuthor()")
	default:
		unexpected()
	}
}
func getAllBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}
func sortAvailable() []book {
	var availableBooks []book
	for _, b := range books {
		if b.Status == "available" {
			availableBooks = append(availableBooks, b)
		}
	}
	return availableBooks
}
func getAvailableBooks(c *gin.Context) {
	availableBooks := sortAvailable()
	c.IndentedJSON(http.StatusOK, availableBooks)
}
func sortIssued() []book {
	var issuedBooks []book
	for _, b := range books {
		if b.Status == "issued" {
			issuedBooks = append(issuedBooks, b)
		}
	}
	return issuedBooks
}
func getIssuedBooks(c *gin.Context) {
	issuedBooks := sortIssued()
	c.IndentedJSON(http.StatusOK, issuedBooks)
}
func lookupBooksByISBN(isbn string) (*book, error) {
	for i, b := range books {
		if b.ISBN == isbn {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book Not Found")
}
func handleLookupBooksByISBN(c *gin.Context) {
	isbn := c.Param("isbn")
	book, err := lookupBooksByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book Not Found"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}
func lookupBooksByTitle(title string) (*book, error) {
	for i, b := range books {
		if b.Title == title {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book Not Found")
}
func handlelookupBooksByTitle(c *gin.Context) {
	title := c.Param("title")
	book, err := lookupBooksByTitle(title)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book Not Found"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}
func lookupBooksByAuthor(author string) (*book, error) {
	for i, b := range books {
		if b.Author == author {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book Not Found")
}
func handlelookupBooksByAuthor(c *gin.Context) {
	author := c.Param("author")
	book, err := lookupBooksByAuthor(author)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book Not Found"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func handleMembers() {
	var choice int
	fmt.Println("Enter Your Choice:")
	fmt.Println("[1] Get all Members in the database")
	fmt.Println("[2] Get Member details by ID")
	fmt.Println("[3] Get books currently issued by a member")
	fmt.Println("[4] Get total books issued by a member")
	fmt.Println("[5] Browse Members using Member name")
	fmt.Println("[6] Browse Members using Contact number")
	fmt.Println("[0] Exit")
	fmt.Scanln(&choice)
	if choice == 0 {
		exitChoice()
	}
	if !(choice == 1 || choice == 2 || choice == 3 || choice == 4 || choice == 5 || choice == 6) {
		clearScr()
		fmt.Print("Invalid Input, Please Reconsider\n\n")
		handleMembers()
	}
	clearScr()
	router := gin.Default()
	switch choice {
	case 1:
		router.GET("/members", getAllMembers)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to getAllMembers()")
	case 2:
		router.GET("/members/:uid", handleLookupMembersByID)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to handleLookupMembersByID()")
	case 3:
		router.GET("/members/:name", handleLookupMembersByName)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to handleLookupMembersByName()")
	case 4:
		router.GET("/members/:name", handleLookupMembersByName)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to handleLookupMembersByName()")
	case 5:
		router.GET("/members/:name", handleLookupMembersByName)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to handleLookupMembersByName()")
	case 6:
		router.GET("/members/:contactnum", handleLookupMembersByContactNum)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to handleLookupMembersByContactNum()")
	default:
		unexpected()
	}
}
func getAllMembers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}
func lookupMembersByID(uid string) (*member, error) {
	for i, m := range members {
		if m.UID == uid {
			return &members[i], nil
		}
	}
	return nil, errors.New("Member Not Found")
}
func handleLookupMembersByID(c *gin.Context) {
	uid := c.Param("uid")
	member, err := lookupMembersByID(uid)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Member Not Found"})
		return
	}
	c.IndentedJSON(http.StatusOK, member)
}
func lookupMembersByName(name string) (*member, error) {
	for i, m := range members {
		if m.Name == name {
			return &members[i], nil
		}
	}
	return nil, errors.New("Member Not Found")
}
func handleLookupMembersByName(c *gin.Context) {
	name := c.Param("name")
	member, err := lookupMembersByName(name)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Member Not Found"})
		return
	}
	c.IndentedJSON(http.StatusOK, member)
}
func lookupMembersByContactNum(contactNum string) (*member, error) {
	for i, m := range members {
		if m.ContactNum == contactNum {
			return &members[i], nil
		}
	}
	return nil, errors.New("Member Not Found")
}
func handleLookupMembersByContactNum(c *gin.Context) {
	contactNum := c.Param("contactnum")
	member, err := lookupMembersByContactNum(contactNum)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Member Not Found"})
		return
	}
	c.IndentedJSON(http.StatusOK, member)
}

func handleBookStatus() {
	var choice int
	fmt.Println("Enter Your Choice:")
	fmt.Println("[1] Issue Book")
	fmt.Println("[2] Return Book")
	fmt.Println("[0] Exit")
	fmt.Scanln(&choice)
	if choice == 0 {
		exitChoice()
	}
	if !(choice == 1 || choice == 2) {
		clearScr()
		fmt.Print("Invalid Input, Please Reconsider\n\n")
		handleBookStatus()
	}
	router := gin.Default()
	switch choice {
	case 1:
		router.PATCH("/issue", issueBook)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to issueBook()")
	case 2:
		router.PATCH("/return", returnBook)
		checkErr(router.Run("localhost:8080"), "Error due to GIN Router due to returnBook()")
	default:
		unexpected()
	}
}
func issueBook(c *gin.Context) {
	isbn, ok := c.GetQuery("isbn")
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Missing id query parameter"})
		return
	}
	book, err := lookupBooksByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book Not found"})
		return
	}
	if book.Status == "issued" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book Not available"})
		return
	}
	book.Status = "issued"
	c.IndentedJSON(http.StatusOK, book)
}
func returnBook(c *gin.Context) {
	isbn, ok := c.GetQuery("isbn")
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Missing id query parameter"})
		return
	}
	book, err := lookupBooksByISBN(isbn)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book Not found"})
		return
	}
	if book.Status == "available" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book already available"})
		return
	}
	book.Status = "available"
	c.IndentedJSON(http.StatusOK, book)
}
