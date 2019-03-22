package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sort"
	"syscall"

	"net/textproto"

	"strings"

	"github.com/jordan-wright/email"
	"golang.org/x/crypto/ssh/terminal"
)

//Card contains the inputted data
type Card struct {
	Student string `json:"Student"`
	Port1   string `json:"Port1"`
	Port2   string `json:"Port2"`
}

//Cards struct for array of cards
type Cards struct {
	Card []Card `json:"Card"`
}

func open(w http.ResponseWriter, r *http.Request) {

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'cards' which we defined above

	// o := Card{Student: "st18u138u", Port1: "8080", Port2: "3000"}

	// spew.Dump(cards)

	fmt.Println("Done printing entries")

}

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

func send(p1 string, p2 string, student string) {

	e := &email.Email{
		To:      []string{"s3646416@student.rmit.edu.au"},
		From:    "Inci Keleher <inci.keleher@gmail.com>",
		Subject: "Test",
		Text:    []byte("ports: " + p1 + " and " + p2 + "assigned to " + student),
		Headers: textproto.MIMEHeader{},
	}

	username, password, err := credentials()
	if err != nil {
		panic(err)
	}

	fmt.Println("\nsending mail ...")
	if !strings.HasSuffix(username, "@gmail.com") {
		username += "@gmail.com"
	}
	err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", username, password, "smtp.gmail.com"))
	if err != nil {
		panic(err)
	}

}

func addcard(w http.ResponseWriter, r *http.Request) {
	//create portlist file to store data

	//create a new entry for the student
	card := new(Card)
	card.Student = r.FormValue("student")
	card.Port1 = r.FormValue("port1")
	card.Port2 = r.FormValue("port2")

	fmt.Println("port1: " + card.Port1)
	fmt.Println("port2: " + card.Port2)

	//booleans for port 1 and 2 for email
	var port1bool = false
	var port2bool = false

	// card := new(Card)
	// card.Student = "John"
	// card.Port1 = "7070"
	// card.Port2 = "3000"

	//open file for writing
	fw, err := os.OpenFile("ports.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	//open file for scanner
	file, err := os.Open("ports.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	//sort for faster search
	sort.Strings(lines)

	//search for port 1
	i := sort.SearchStrings(lines, card.Port1)

	//add port 1 if not in database
	if i < len(lines) && lines[i] != card.Port1 {
		fmt.Println(i < len(lines) && lines[i] == card.Port1)

		//write to file
		if _, err := fw.Write([]byte("\n" + card.Port1)); err != nil {
			log.Fatal(err)
		}
		println("Entry added\n")
		port1bool = true

	} else {
		fmt.Fprint(w, "ERROR: PORT 1 TAKEN, PLEASE TRY AGAIN")
	}

	//search for port 2
	i = sort.SearchStrings(lines, card.Port2)

	//add port 2 if not in database
	if i < len(lines) && lines[i] != card.Port2 {
		fmt.Println(i < len(lines) && lines[i] == card.Port2)

		//write to file
		if _, err := fw.Write([]byte("\n" + card.Port2)); err != nil {
			log.Fatal(err)
		}
		println("Entry added\n")
		port2bool = true

	} else {
		fmt.Fprint(w, "ERROR: PORT 2 TAKEN, PLEASE TRY AGAIN")
	}

	if (port1bool == true) && (port2bool == true) {
		send(card.Port1, card.Port2, card.Student)
	}

}

func main() {
	fmt.Println("Welcome to PortSelector")

	//print array for check
	// for _, line := range lines {

	// 	if line != card.Port1 {

	// 		if _, err := fw.Write([]byte("\n" + card.Port1)); err != nil {
	// 			log.Fatal(err)
	// 		}

	// 		println("Entry added\n")
	// 		return
	// 	}

	// 	if line == card.Port1 {
	// 		fmt.Println(line + " is taken")
	// 		return
	// 	}

	// }

	http.HandleFunc("/addcard", addcard)
	http.HandleFunc("/open", open)
	http.ListenAndServe(":8080", nil)

}
