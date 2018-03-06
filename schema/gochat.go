package schema

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
