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

type Producao interface {
	AdicionaFila(ctx context.Context, obj map[string]string) error
}

type producaoClient struct {
	httpClient *http.Client
	url        string
}

func NewProducao() Producao {
	url := os.Getenv("PRODUCAO_URL")
	if url == "" {
		log.Fatal("PRODUCAO_URL environment variable not set")
	}
	return &producaoClient{
		httpClient: http.DefaultClient,
		url:        url,
	}
}
func (c *producaoClient) AdicionaFila(ctx context.Context, obj map[string]string) error {
	jsonBody, err := json.Marshal(obj)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/internal/producao", c.url), bytes.NewBuffer(jsonBody))
	if err != nil {
		return errorx.InternalError.New(fmt.Sprintf("não foi possível inicializar producao client %s", err.Error()))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errorx.InternalError.New(fmt.Sprintf("producao server retornou error %s", err.Error()))
	}

	if resp == nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return errorx.InternalError.New(fmt.Sprintf("producao server retornou status %s", resp.Status))
	}
	defer resp.Body.Close()

	return nil
}
