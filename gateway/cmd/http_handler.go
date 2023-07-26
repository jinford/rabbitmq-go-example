package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jinford/rabbitmq-go-example/calculator/pkg/message"
	"github.com/labstack/echo/v4"
)

type PostAddRequestBody struct {
	A int
	B int
}

type PostAddResponseBody struct {
	ID string
}

func PostAddHandler(p *MessageProducer) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("[REST API REQ] POST /add")

		ctx := c.Request().Context()

		reqBody := new(PostAddRequestBody)
		if err := json.NewDecoder(c.Request().Body).Decode(reqBody); err != nil {
			err = fmt.Errorf("json.Decode: %w", err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		fmt.Println("A B:", reqBody.A, reqBody.B)

		id := uuid.NewString()

		log.Println("[PRODUCE MESSAGE] caluculator.add", id)
		if err := InvokeCalcuratorAdd(ctx, p, id, reqBody.A, reqBody.B); err != nil {
			err = fmt.Errorf("InvokeCalcuratorAdd: %w", err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		resBody := &PostAddResponseBody{
			ID: id,
		}

		log.Println("[REST API RES] POST /add")
		return c.JSON(http.StatusCreated, resBody)
	}
}

func InvokeCalcuratorAdd(ctx context.Context, p *MessageProducer, id string, a int, b int) error {
	reqBody := &message.RequestBody{
		ID: id,
		A:  a,
		B:  b,
	}

	msgBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	if err := p.Produce(ctx, message.CommandAdd, msgBody); err != nil {
		return fmt.Errorf("p.Produce: %w", err)
	}

	return nil
}

type GetAddResponseBody struct {
	Result int
}

func GetAddHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("[REST API REQ] GET /add/:id")

		id := c.Param("id")

		result, exists := addResultMem[id]
		if !exists {
			status := http.StatusNotFound
			return c.String(status, http.StatusText(status))
		}

		resBody := &GetAddResponseBody{
			Result: result,
		}

		log.Println("[REST API RES] GET /add/:id")
		return c.JSON(http.StatusOK, resBody)
	}
}
