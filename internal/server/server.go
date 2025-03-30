package server

import (
	graph2 "TestOzon/internal/handler/graph"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func graghQlServer(handlers *graph2.Resolver) (*handler.Server, error) {
	srv := handler.New(graph2.NewExecutableSchema(graph2.Config{Resolvers: handlers}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv, nil
}

func (s *Server) StartServer(ctx context.Context, port string, handlers *graph2.Resolver) error {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := s.httpServer.Shutdown(ctx); err != nil {
			slog.Error(fmt.Sprintf("server is not shutting down! Reason: %v", err.Error()))
		}

		close(idleConnsClosed)
	}()

	addr := ":" + port

	s.httpServer = &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	srvQL, err := graghQlServer(handlers)
	if err != nil {
		slog.Error(fmt.Sprintf("falied to init graphQl server handleers. Reason: %s", err.Error()))
		return err
	}

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srvQL)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			slog.Error(fmt.Sprintf("server is not running. Reason: %s", err.Error()))
		}
	}()

	<-idleConnsClosed
	return nil
}
