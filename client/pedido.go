package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joomcode/errorx"
	"log"
	"net/http"
	"os"
)

type Pedido interface {
	AtualizaStatus(ctx context.Context, status, id string) error
}

type pedidoCliente struct {
	httpClient *http.Client
	url        string
}

func NewPedido() Pedido {
	url := os.Getenv("PEDIDOS_URL")
	if url == "" {
		log.Fatal("PEDIDOS_URL environment variable not set")
	}
	return &pedidoCliente{
		httpClient: http.DefaultClient,
		url:        url,
	}
}
func (c *pedidoCliente) AtualizaStatus(ctx context.Context, status, id string) error {
	body := map[string]string{
		"status": status,
	}

	jsonBody, err := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("%s/pedido/%s", c.url, id), bytes.NewBuffer(jsonBody))
	if err != nil {
		return errorx.InternalError.New(fmt.Sprintf("não foi possível inicializar pedido client %s", err.Error()))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errorx.InternalError.New(fmt.Sprintf("pedido server retornou error %s", err.Error()))
	}

	if resp == nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return errorx.InternalError.New(fmt.Sprintf("pedido server retornou status %s", resp.Status))
	}
	defer resp.Body.Close()

	return nil
}
