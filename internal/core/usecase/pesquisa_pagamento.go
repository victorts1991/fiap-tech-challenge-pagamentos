package usecase

import (
	"context"
	"fiap-tech-challenge-pagamentos/internal/adapters/repository"
	"fiap-tech-challenge-pagamentos/internal/core/domain"
	"fmt"
	serverErr "github.com/rhuandantas/fiap-tech-challenge-commons/pkg/errors"
)

type PesquisaPagamento interface {
	PesquisaPorPedidoID(ctx context.Context, pedidoId string) (*domain.Pagamento, error)
}

type pesquisaPagamento struct {
	pagamentoRepo repository.PagamentoRepo
}

func (uc *pesquisaPagamento) PesquisaPorPedidoID(ctx context.Context, pedidoId string) (*domain.Pagamento, error) {
	existe, err := uc.pagamentoRepo.PesquisaPorPedidoID(ctx, pedidoId)
	if err != nil {
		return nil, err
	}

	if existe == nil {
		return nil, serverErr.NotFound.New(fmt.Sprintf("nenhum pagamento encontrado para pedido %s", pedidoId))
	}

	return existe, nil
}

func NewPesquisaPagamento(pagamentoRepo repository.PagamentoRepo) PesquisaPagamento {
	return &pesquisaPagamento{
		pagamentoRepo: pagamentoRepo,
	}
}
