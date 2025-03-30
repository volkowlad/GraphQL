package server

import (
	graph2 "TestOzon/internal/handler/graph"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"log/slog"
	"net/http"
)

func StartServer(port string, handlers *graph2.Resolver) error {
	srv := handler.New(graph2.NewExecutableSchema(graph2.Config{Resolvers: handlers}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	slog.Info(fmt.Sprintf("connect to localhost:%s for GraphQL playground on port", port))
	return http.ListenAndServe(":"+port, nil)
}
