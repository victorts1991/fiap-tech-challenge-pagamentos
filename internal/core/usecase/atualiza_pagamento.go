package usecase

import (
	"context"
	"fiap-tech-challenge-pagamentos/internal/adapters/repository"
	"fmt"
	"github.com/rhuandantas/fiap-tech-challenge-commons/pkg/errors"
)

type AtualizaPagamento interface {
	Atualiza(ctx context.Context, status, pedidoID string) error
}

type atualizaPagamentoUC struct {
	pagamentoRepoRepo repository.PagamentoRepo
}

func (uc *atualizaPagamentoUC) Atualiza(ctx context.Context, status, pedidoID string) error {
	pagamento, err := uc.pagamentoRepoRepo.PesquisaPorPedidoID(ctx, pedidoID)
	if err != nil {
		return err
	}

	if pagamento == nil {
		return errors.NotFound.New(fmt.Sprintf("pagamento não encontrado para pedido %s", pedidoID))
	}

	err = uc.pagamentoRepoRepo.AtualizaStatus(ctx, status, pedidoID)
	if err != nil {
		return err
	}

	return nil
}

func NewAtualizaPagamento(pagamentoRepoRepo repository.PagamentoRepo) AtualizaPagamento {
	return &atualizaPagamentoUC{
		pagamentoRepoRepo: pagamentoRepoRepo,
	}
}
