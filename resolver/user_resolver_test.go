package resolver

import (
	"context"
	"testing"

	"github.com/beinan/graphql-server/service"
	graphql "github.com/graph-gophers/graphql-go"
)

type ID = string
type MockedUserService struct {
	t *testing.T
}

func (mo *MockedUserService) GetByID(ctx context.Context, id ID) (*User, error) {
	mo.t.Logf("Mocked user service get user by id: %v", id)
	user := &User{
		Id:   "123",
		Name: "Mocked Return",
	}
	return user, nil
}

func (mo *MockedUserService) GetByIDs(ctx context.Context, ids []ID) ([]User, error) {
	mo.t.Logf("Mocked user service get user by id: %v", ids)

	users := []User{
		User{
			Id:   "123",
			Name: "Mocked Return",
		},
	}
	return users, nil
}

func (mo *MockedUserService) Create(ctx context.Context, user User) error {
	return nil
}

func TestGetUserResolver(t *testing.T) {
	t.Logf("Testing get user resolver.")
	resolver := Resolver{
		services: &service.Services{
			UserService: &MockedUserService{t},
		},
	}

	input := &struct { //anonymous struct for input
		Id graphql.ID
	}{graphql.ID("123")}
	userResolver, err := resolver.GetUser(
		context.Background(), input)
	if userResolver.user.Name != "Mocked Return" {
		t.Errorf("Get user resolver returns incorrect result %v", userResolver.user.Name)
	}
	if err != nil {
		t.Errorf("Get user resolver return unexpected error: %v", err)
	}
}
