package usecase

import (
	"context"
	"fiap-tech-challenge-pagamentos/internal/adapters/repository"
)

type AtualizaPagamento interface {
	Atualiza(ctx context.Context, status, pedidoID string) error
}

type atualizaPagamentoUC struct {
	pagamentoRepoRepo repository.PagamentoRepo
}

func (uc *atualizaPagamentoUC) Atualiza(ctx context.Context, status, pedidoID string) error {
	err := uc.pagamentoRepoRepo.AtualizaStatus(ctx, status, pedidoID)

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
