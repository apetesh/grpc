package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	api "github.com/apetesh/grpc/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	phonebookClient := api.NewPhonebookClient(conn)
	resp, err := phonebookClient.AddContact(context.Background(), &api.AddContactRequest{Name: "asaf petesh", PhoneNumber: "1700707070"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("new contact id: %d\n", resp.Id)

	contact, err := phonebookClient.GetContactByID(context.Background(), &api.GetContactRequest{Id: 0})
	if err != nil {
		panic(err)
	}
	fmt.Printf("new contact: %+v\n", contact)

	addedIds, err := sendContactsFromFile("./contacts", phonebookClient)
	for _, addedID := range addedIds {
		fmt.Println(addedID)
	}
}

func sendContactsFromFile(filePath string, phonebookClient api.PhonebookClient) ([]int32, error) {
	file, err := os.Open("./contacts")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	stream, err := phonebookClient.AddContacts(context.Background())
	for scanner.Scan() {
		if err != nil {
			log.Fatal(err)
		}
		contactLine := scanner.Text()
		splittedContactLine := strings.Split(contactLine, "-")
		stream.Send(&api.AddContactRequest{Name: splittedContactLine[0], PhoneNumber: splittedContactLine[1]})
	}
	if err != nil {
		log.Fatal(err)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	resp, err := stream.CloseAndRecv()
	return resp.Ids, err
}
