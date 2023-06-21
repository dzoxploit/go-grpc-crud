package server

import (
	"context"
	"log"

	"github.com/dzoxploit/go-grpc-crud/entities"
	"github.com/dzoxploit/go-grpc-crud/protobuf/golang_protobuf_contact"
	"github.com/dzoxploit/go-grpc-crud/repos"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CRUDServiceServer struct {
	repo *repos.GenericRepo[entities.Contact]
}

func NewCRUDServiceServer(repo *repos.GenericRepo[entities.Contact]) *CRUDServiceServer {
	server := CRUDServiceServer{
		repo: repo,
	}
	return &server
}

func (c CRUDServiceServer) Create(context.Context, *golang_protobuf_contact.ProtoContactRepo_ProtoContact) (*golang_protobuf_contact.ProtoContactRepo_ProtoContact, error) {
	myMessage := &golang_protobuf_contact.ProtoContactRepo_ProtoContact{}
	myMessage.ID = uuid.New().String()
	return myMessage, nil
}

func (c CRUDServiceServer) GetOne(_ context.Context, id *wrapperspb.StringValue) (*golang_protobuf_contact.ProtoContactRepo_ProtoContact, error) {
	contact, err := (*c.repo).GetOne(string(id.Value))
	if err != nil {
		log.Printf("failed to get Brand: %v", err)
		return &golang_protobuf_contact.ProtoContactRepo_ProtoContact{}, err
	}
	return repos.ToProtoContract(contact), nil
}

func (c CRUDServiceServer) GetList(_ *emptypb.Empty, stream golang_protobuf_contact.CRUD_GetListServer) error {
	for _, contact := range (*c.repo).GetList() {
		if err := stream.Send(repos.ToProtoContract(contact)); err != nil {
			return err
		}
	}
	return nil
}

func (c CRUDServiceServer) Update(_ context.Context, message *golang_protobuf_contact.UpdateRequest) (*golang_protobuf_contact.ProtoContactRepo_ProtoContact, error) {
	contact, err := (*c.repo).Update(string(message.ID.Value), repos.ToContact(message.Contact))
	if err != nil {
		log.Printf("failed to update Brand: %v", err)
		return &golang_protobuf_contact.ProtoContactRepo_ProtoContact{}, err
	}
	return repos.ToProtoContract(contact), nil
}

func (c CRUDServiceServer) Delete(_ context.Context, message *wrapperspb.StringValue) (*wrapperspb.BoolValue, error) {
	success, err := (*c.repo).DeleteOne(string(message.Value))
	if err != nil {
		log.Printf("failed to delete Brand: %v", err)
		return wrapperspb.Bool(false), err
	}
	return wrapperspb.Bool(success), nil
}