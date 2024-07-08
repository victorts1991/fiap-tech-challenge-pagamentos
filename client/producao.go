package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rhuandantas/fiap-tech-challenge-commons/pkg/messaging"
	"log"
	"os"
)

type Producao interface {
	AdicionaFila(ctx context.Context, obj map[string]string) error
}

type producaoClient struct {
	publisher messaging.Client
	queueUrl  string
}

func NewProducao(publisher messaging.SqsClient) Producao {
	url := os.Getenv("PRODUCAO_QUEUE")
	if url == "" {
		log.Fatal("PRODUCAO_QUEUE environment variable not set")
	}
	return &producaoClient{
		publisher: publisher,
		queueUrl:  url,
	}
}
func (c *producaoClient) AdicionaFila(ctx context.Context, obj map[string]string) error {
	jsonBody, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = c.publisher.Publish(ctx, c.queueUrl, jsonBody)
	if err != nil {
		return err
	}

	fmt.Println("message published successfully")

	return nil
}
