package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	api "github.com/apetesh/grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	phonebook  *Phonebook
	port       int
	grpcServer *grpc.Server
}

func NewServer(phonebook *Phonebook, port int) *Server {
	return &Server{phonebook, port, grpc.NewServer()}
}

func (s *Server) Stop() {
	log.Printf("gracefully stopping phonebook service")
	s.grpcServer.GracefulStop()
	log.Printf("phonebook service stopped")
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	api.RegisterPhonebookServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)
	return s.grpcServer.Serve(lis)
}

func (s *Server) AddContact(ctx context.Context, contact *api.AddContactRequest) (*api.AddContactResponse, error) {
	newContact := &Contact{name: contact.Name, phoneNumber: contact.PhoneNumber}
	id := s.phonebook.AddContact(newContact)
	return &api.AddContactResponse{Id: id}, nil
}

func (s *Server) GetContactByID(ctx context.Context, id *api.GetContactRequest) (*api.GetContactResponse, error) {
	contact := s.phonebook.GetContact(id.Id)
	if contact == nil {
		return nil, status.Errorf(codes.NotFound, "contact with id %d was not found", id.Id)
	}
	return &api.GetContactResponse{Name: contact.name, PhoneNumber: contact.phoneNumber}, nil
}

func (s *Server) DeleteContact(ctx context.Context, id *api.DeleteContactRequest) (*api.DeleteContactResponse, error) {
	err := s.phonebook.DeleteContact(id.Id)
	if err != nil {
		return nil, fmt.Errorf("contact with id %d was not found", id.Id)
	}
	return &api.DeleteContactResponse{Id: id.Id}, nil
}

func (s *Server) AddContacts(stream api.Phonebook_AddContactsServer) error {
	ids := make([]int32, 0)
	for {
		contact, err := stream.Recv()
		log.Println("Received value")
		if err == io.EOF {
			stream.SendAndClose(&api.AddContactsResponse{Ids: ids})
			return nil
		}
		if err != nil {
			return err
		}
		newContact := &Contact{name: contact.Name, phoneNumber: contact.PhoneNumber}
		id := s.phonebook.AddContact(newContact)
		ids = append(ids, id)
	}
}
