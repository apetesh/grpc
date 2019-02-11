package server

import (
	"fmt"
	"sync"
)

type Contact struct {
	id          int32
	name        string
	phoneNumber string
}

type Phonebook struct {
	contacts map[int32]*Contact
	nextID   int32
	mapLock  sync.RWMutex
	idLock   sync.RWMutex
}

func NewPhonebook() *Phonebook {
	return &Phonebook{contacts: make(map[int32]*Contact)}
}

func (pb *Phonebook) ListContacts() []*Contact {
	pb.mapLock.RLock()
	defer pb.mapLock.RUnlock()
	contacts := make([]*Contact, 0)
	fmt.Printf("number of contacts: %d", len(pb.contacts))
	for _, contact := range pb.contacts {
		contacts = append(contacts, contact)
	}
	return contacts
}

func (pb *Phonebook) GetContact(id int32) *Contact {
	pb.mapLock.RLock()
	defer pb.mapLock.RUnlock()
	if contact, exists := pb.contacts[id]; exists {
		return contact
	}
	return nil
}

func (pb *Phonebook) AddContact(contact *Contact) int32 {
	pb.mapLock.Lock()
	pb.idLock.Lock()
	defer pb.mapLock.Unlock()
	defer pb.idLock.Unlock()
	contact.id = pb.nextID
	pb.contacts[contact.id] = contact
	pb.nextID = contact.id + 1
	return contact.id
}

func (pb *Phonebook) DeleteContact(id int32) error {
	if _, exists := pb.contacts[id]; exists {
		pb.mapLock.Lock()
		defer pb.mapLock.Unlock()
		delete(pb.contacts, id)
		return nil
	}
	return fmt.Errorf("no contact with id %d was found", id)

}
