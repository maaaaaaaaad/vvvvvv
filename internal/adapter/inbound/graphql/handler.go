package graphql

import (
	"jello-mark-backend/internal/domain"
	"jello-mark-backend/internal/ports/inbound"
	"net/http"
	"time"

	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type userDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

func NewGraphQLHandler(
	users inbound.UserService,
	playgroundEnabled bool,
) http.Handler {
	s, err := schema(users)
	if err != nil {
		panic(err)
	}

	return handler.New(
		&handler.Config{
			Schema:     &s,
			Pretty:     true,
			GraphiQL:   playgroundEnabled,
			Playground: playgroundEnabled,
		},
	)
}

func schema(users inbound.UserService) (
	gql.Schema,
	error,
) {
	userType := gql.NewObject(
		gql.ObjectConfig{
			Name: "User",
			Fields: gql.Fields{
				"id":        &gql.Field{Type: gql.NewNonNull(gql.String)},
				"name":      &gql.Field{Type: gql.NewNonNull(gql.String)},
				"createdAt": &gql.Field{Type: gql.NewNonNull(gql.String)},
			},
		},
	)
	query := gql.NewObject(
		gql.ObjectConfig{
			Name: "Query",
			Fields: gql.Fields{
				"health": &gql.Field{
					Type: gql.NewNonNull(gql.String),
					Resolve: func(p gql.ResolveParams) (
						interface{},
						error,
					) {
						return "ok", nil
					},
				},
				"users": &gql.Field{
					Type: gql.NewList(userType),
					Args: gql.FieldConfigArgument{
						"limit": &gql.ArgumentConfig{Type: gql.Int},
					},
					Resolve: func(p gql.ResolveParams) (
						interface{},
						error,
					) {
						limit := 100
						if v, ok := p.Args["limit"]; ok {
							if i, ok2 := v.(int); ok2 {
								limit = i
							}
						}
						ctx := p.Context
						list, err := users.List(
							ctx,
							limit,
						)
						if err != nil {
							return nil, err
						}
						res := make(
							[]userDTO,
							0,
							len(list),
						)
						for _, u := range list {
							user := domain.User(u)
							res = append(
								res,
								userDTO{
									ID:        string(user.ID),
									Name:      user.Name,
									CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
								},
							)
						}
						return res, nil
					},
				},
			},
		},
	)
	mutation := gql.NewObject(
		gql.ObjectConfig{
			Name: "Mutation",
			Fields: gql.Fields{
				"createUser": &gql.Field{
					Type: userType,
					Args: gql.FieldConfigArgument{
						"name": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.String)},
					},
					Resolve: func(p gql.ResolveParams) (
						interface{},
						error,
					) {
						name := ""
						if v, ok := p.Args["name"].(string); ok {
							name = v
						}
						u, err := users.Create(
							p.Context,
							name,
						)
						if err != nil {
							return nil, err
						}
						user := domain.User(u)
						return userDTO{
							ID:        string(user.ID),
							Name:      user.Name,
							CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
						}, nil
					},
				},
			},
		},
	)
	return gql.NewSchema(
		gql.SchemaConfig{
			Query:    query,
			Mutation: mutation,
		},
	)
}
