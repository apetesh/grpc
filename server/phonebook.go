package main

import (
	"fmt"
	"sync"
)

type Contact struct {
	name        string
	phoneNumber string
}

type Phonebook struct {
	contacts map[int32]*Contact
	mapLock  sync.RWMutex
	idLock   sync.RWMutex
	nextID   int32
}

func NewPhonebook() *Phonebook {
	return &Phonebook{contacts: make(map[int32]*Contact)}
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
	newContactID := pb.nextID
	pb.contacts[newContactID] = contact
	pb.nextID = newContactID + 1
	return newContactID
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
