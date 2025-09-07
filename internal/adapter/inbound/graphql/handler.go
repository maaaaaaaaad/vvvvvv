package graphql

import (
	"encoding/json"
	"jello-mark-backend/internal/application"
	"net/http"
	"strconv"
	"time"

	gql "github.com/graphql-go/graphql"
)

type Handler struct {
	Users application.UserService
}

type gqlReq struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type userDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

func schema(users application.UserService) (
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
							res = append(
								res,
								userDTO{
									ID:        string(u.ID),
									Name:      u.Name,
									CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
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
						return userDTO{
							ID:        string(u.ID),
							Name:      u.Name,
							CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
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

func (h Handler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
	dec := json.NewDecoder(r.Body)
	var req gqlReq
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request"))
		return
	}
	s, err := schema(h.Users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("schema error"))
		return
	}
	res := gql.Do(
		gql.Params{
			Schema:         s,
			RequestString:  req.Query,
			VariableValues: req.Variables,
			OperationName:  req.OperationName,
			Context:        r.Context(),
		},
	)
	bs, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("encode error"))
		return
	}
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.Header().Set(
		"Content-Length",
		strconv.Itoa(len(bs)),
	)
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}
