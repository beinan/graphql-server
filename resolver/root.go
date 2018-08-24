package resolver

import (
	"github.com/beinan/graphql-server/service"
)

//resolver type
type Resolver struct {
	services *service.Services
}

func MkRootResolver(services *service.Services) *Resolver {
	return &Resolver{services}
}
