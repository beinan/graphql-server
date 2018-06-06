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
  input AuthInput {
    id: ID!
    password: String!
  }
    
  type Query {
    getUser(id: ID!): User  
  }
  type Mutation {
    signUp(input: AuthInput!): User
    signIn(input: AuthInput!): String!
  }
  schema {
		query: Query
		mutation: Mutation
	}
`
