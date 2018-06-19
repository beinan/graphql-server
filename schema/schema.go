package schema

var Schema = `
  enum Gender {
    MALE
    FEMALE
  }
  interface IDNode{
    id: ID!
  }
  type Edge{
    node: IDNode!
    cursor: String!
  }
  type Edges {
    edges: [Edge]
    hasMore: Boolean!
  }
  input EdgesInput {
    cursor: String
    pageSize: Int!
    isRev: Boolean #is reversed ordering
  }
  type User implements IDNode{
    id: ID!
    name: String!
    gender: Gender
    friends(pageNum: Int = 0, pageSize: Int = 20): [User!]
    friendEdges(input: EdgesInput): Edges
  }
  input AuthInput {
    loginName: String!
    password: String!
  }
    
  type Query {
    getUser(id: ID!): User  
  }
  type Mutation {
    signUp(input: AuthInput!): User
    signIn(input: AuthInput!): String!
    addFriend(fromId: ID!, toId: ID!): Boolean!
  }
  schema {
		query: Query
		mutation: Mutation
	}
`
