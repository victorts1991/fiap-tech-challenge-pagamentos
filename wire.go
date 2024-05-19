// wire.go
//go:build wireinject

package main

import (
	"fiap-tech-challenge-pagamentos/client"
	"fiap-tech-challenge-pagamentos/internal/adapters/http"
	"fiap-tech-challenge-pagamentos/internal/adapters/http/handlers"
	"fiap-tech-challenge-pagamentos/internal/adapters/repository"
	"fiap-tech-challenge-pagamentos/internal/core/usecase"
	db "github.com/rhuandantas/fiap-tech-challenge-commons/pkg/db/mysql"
	"github.com/rhuandantas/fiap-tech-challenge-commons/pkg/middlewares/auth"
	"github.com/rhuandantas/fiap-tech-challenge-commons/pkg/util"

	"github.com/google/wire"
)

func InitializeWebServer() (*http.Server, error) {
	wire.Build(db.NewMySQLConnector,
		util.NewCustomValidator,
		auth.NewJwtToken,
		repository.NewPagamentoRepo,
		usecase.NewRealizaCheckout,
		usecase.NewPesquisaPagamento,
		handlers.NewHealthCheck,
		handlers.NewPagamento,
		client.NewPedido,
		http.NewAPIServer)
	return &http.Server{}, nil
}
