package usecase

import (
	"context"
	"fiap-tech-challenge-pagamentos/internal/adapters/repository"
	"fiap-tech-challenge-pagamentos/internal/core/domain"
)

type CadastrarPagamento interface {
	Cadastra(ctx context.Context, pagamento *domain.Pagamento) error
}

type cadastraFila struct {
	pagamentoRepoRepo repository.PagamentoRepo
}

func (uc *cadastraFila) Cadastra(ctx context.Context, pagamento *domain.Pagamento) error {
	pagamento.Status = domain.StatusPendente
	err := uc.pagamentoRepoRepo.Insere(ctx, pagamento)

	if err != nil {
		return err
	}
	return nil
}

func NewCadastraPagamento(pagamentoRepoRepo repository.PagamentoRepo) CadastrarPagamento {
	return &cadastraFila{
		pagamentoRepoRepo: pagamentoRepoRepo,
	}
}
