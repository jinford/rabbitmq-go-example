package main

import (
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

const (
	port = 1234
)

type HttpServer struct {
	engine *echo.Echo
}

func NewHttpServer(p *MessageProducer) *HttpServer {
	engine := echo.New()

	engine.POST("/add", PostAddHandler(p))
	engine.GET("/add/:id", GetAddHandler())

	return &HttpServer{
		engine: engine,
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	errStream := make(chan error)
	go func() {
		addr := fmt.Sprintf("0.0.0.0:%d", port)
		if err := s.engine.Start(addr); err != nil {
			errStream <- fmt.Errorf("s.engine.Start: %w", err)
		}
	}()

	select {
	case err := <-errStream:
		return err
	case <-ctx.Done():
		if err := s.engine.Shutdown(context.Background()); err != nil {
			log.Println("s.engine.Shutdown: %w", err)
		}
		return nil
	}
}
