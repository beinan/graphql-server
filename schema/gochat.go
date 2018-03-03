package schema

import (
	graphql "github.com/graph-gophers/graphql-go"
)

var Schema = `
  enum Gender {
    MALE
    FEMALE
  }
  type User {
    id: ID!
    name: String!
    gender: Gender
  }
  input SignUpInput {
    id: ID!
    password: String!
  }
  type Query {
    getUser(id: ID!): User  
  }
  type Mutation {
    signUp(input: SignUpInput!): User
  }
  schema {
		query: Query
		mutation: Mutation
	}
`

//resolver type
type Resolver struct{}

func (r *Resolver) GetUser(args *struct {
	Id graphql.ID
}) *userResolver {
	user := &User{
		Id:     args.Id,
		Name:   "aaaa",
		Gender: "MALE",
	}
	return &userResolver{user}
}

func (r *Resolver) SignUp(args *struct {
	Input *signUpInput
}) *userResolver {
	user := &User{
		Id:     args.Input.Id,
		Name:   "" + args.Input.Password,
		Gender: "MALE",
	}
	return &userResolver{user}
}

type User struct {
	Id     graphql.ID
	Name   string
	Gender string
}

type userResolver struct {
	user *User
}

func (r *userResolver) ID() graphql.ID {
	return r.user.Id
}

func (r *userResolver) Name() string {
	return r.user.Name
}

func (r *userResolver) Gender() *string {
	return &r.user.Gender
}

type signUpInput struct {
	Id       graphql.ID
	Password string
}
