package repos

import (
	"fmt"
	"os"
	"time"

	"github.com/dzoxploit/go-grpc-crud/entities"
	"github.com/dzoxploit/go-grpc-crud/protobuf/golang_protobuf_contact"
	"google.golang.org/protobuf/proto"
)

const STORAGE_FILE = "./contacts-storage.pb"

type ContactRepo struct {
	contacts []entities.Contact
}

func NewContactRepo() *ContactRepo {
	var br = ContactRepo{make([]entities.Contact, 0)}
	br.loadFromFileStorage()
	return &br
}

func (b *ContactRepo) Create(partial entities.Contact) entities.Contact {
	newItem := entities.Contact{
		    ID:         partial.ID,
			Name:       partial.Name,
			Gender:     partial.Gender,
			Phone:      partial.Phone,
			Email:      partial.Email,
			Created_at: time.Now(),
			Updated_at: time.Now(),
	}
	b.contacts = append(b.contacts, newItem)
	b.saveToFileStorage()
	return newItem
}

func (b *ContactRepo) GetList() []entities.Contact {
	return b.contacts
}

func (p *ContactRepo) GetOne(id string) (entities.Contact, error) {
	for _, it := range p.contacts {
		if it.ID == id {
			return it, nil
		}
	}
	return entities.Contact{}, fmt.Errorf("key '%d' not found", id)
}

func (b *ContactRepo) Update(id string, amended entities.Contact) (entities.Contact, error) {
	for i, it := range b.contacts {
		if it.ID == id {
			amended.ID = id
			b.contacts = append(b.contacts[:i], b.contacts[i+1:]...)
			b.contacts = append(b.contacts, amended)
			b.saveToFileStorage()
			return amended, nil
		}
	}
	return entities.Contact{}, fmt.Errorf("key '%d' not found", amended.ID)
}

func (b *ContactRepo) DeleteOne(id string) (bool, error) {
	for i, it := range b.contacts {
		if it.ID == id {
			b.contacts = append(b.contacts[:i], b.contacts[i+1:]...)
			b.saveToFileStorage()
			return true, nil
		}
	}
	return false, fmt.Errorf("key '%d' not found", id)
}

func (b *ContactRepo) saveToFileStorage() error {

	contactsMessage := &golang_protobuf_contact.ProtoContactRepo{
		Contacts: []*golang_protobuf_contact.ProtoContactRepo_ProtoContact{},
	}

	for _, b := range b.contacts {
		contactsMessage.Contacts = append(contactsMessage.Contacts,
			&golang_protobuf_contact.ProtoContactRepo_ProtoContact{
				ID: string(b.ID), Name: b.Name, Gender: b.Gender, Phone: b.Phone, Email: b.Email,
			})

	}

	data, err := proto.Marshal(contactsMessage)
	if err != nil {
		return fmt.Errorf("cannot marshal to binary: %w", err)
	}

	err = os.WriteFile(STORAGE_FILE, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file: %w", err)
	}

	return nil
}

func (b *ContactRepo) loadFromFileStorage() error {

	_, err := os.Stat(STORAGE_FILE)
	if err != nil {
		fmt.Println("storage file is not found, starting with empty storage")
		return nil
	}

	data, err := os.ReadFile(STORAGE_FILE)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file: %w", err)
	}

	var contactsMessage golang_protobuf_contact.ProtoContactRepo
	err = proto.Unmarshal(data, &contactsMessage)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary data to protobuf: %w", err)
	}

	for _, brand := range contactsMessage.Contacts {
		b.contacts = append(b.contacts,
			entities.Contact{
				ID: string(brand.ID), Name: brand.Name, Gender: brand.Gender, Phone: brand.Phone, Email: brand.Email,
			})

	}
	return nil
}

func ToProtoContract(b entities.Contact) *golang_protobuf_contact.ProtoContactRepo_ProtoContact {
	return &golang_protobuf_contact.ProtoContactRepo_ProtoContact{
		ID: string(b.ID), Name: b.Name, Gender: b.Gender, Phone: b.Phone, Email: b.Email,
	}
}

func ToContact(b *golang_protobuf_contact.ProtoContactRepo_ProtoContact) entities.Contact {
	return entities.Contact{
		ID: string(b.ID), Name: (*b).Name, Gender: (*b).Gender, Phone: (*b).Phone, Email: (*b).Email,
	}
}