package usecase

import (
	"context"
	"fiap-tech-challenge-pagamentos/client"
	"fiap-tech-challenge-pagamentos/internal/adapters/repository"
	"fiap-tech-challenge-pagamentos/internal/core/commons"
	"fiap-tech-challenge-pagamentos/internal/core/domain"
	"fmt"
	"github.com/joomcode/errorx"
)

type RealizarCheckout interface {
	Checkout(ctx context.Context, pagamento *domain.Pagamento) error
}

type realizaCheckout struct {
	pagamentoRepo  repository.PagamentoRepo
	pedidoClient   client.Pedido
	producaoClient client.Producao
}

func (uc *realizaCheckout) Checkout(ctx context.Context, pagamento *domain.Pagamento) error {
	existe, err := uc.pagamentoRepo.PesquisaPorPedidoID(ctx, pagamento.PedidoId)
	if err != nil {
		return err
	}

	if existe != nil {
		return errorx.IllegalState.New(fmt.Sprintf("pagamento para pedido %s encontrado com status %s", pagamento.PedidoId, pagamento.Status))
	}

	err = uc.atualizaPedido(ctx, pagamento.PedidoId, pagamento.Status)
	if err != nil {
		return err
	}

	err = uc.producaoClient.AdicionaFila(ctx, map[string]string{
		"pedido_id": pagamento.PedidoId,
		"status":    "recebido",
	})
	if err != nil {
		return err
	}

	err = uc.pagamentoRepo.Insere(ctx, pagamento)
	if err != nil {
		return err
	}

	return nil
}

func (uc *realizaCheckout) atualizaPedido(ctx context.Context, pedidoId string, status string) error {
	var pedidoStatus string
	switch status {
	case domain.StatusAprovado:
		pedidoStatus = commons.StatusPagamentoAprovado
	case domain.StatusRecusado:
		pedidoStatus = commons.StatusPagamentoRecusado
	}

	err := uc.pedidoClient.AtualizaStatus(ctx, pedidoStatus, pedidoId)
	if err != nil {
		return err
	}
	return nil
}

func NewRealizaCheckout(pagamentoRepo repository.PagamentoRepo, pedidoClient client.Pedido, producaoClient client.Producao) RealizarCheckout {
	return &realizaCheckout{
		pagamentoRepo:  pagamentoRepo,
		pedidoClient:   pedidoClient,
		producaoClient: producaoClient,
	}
}
