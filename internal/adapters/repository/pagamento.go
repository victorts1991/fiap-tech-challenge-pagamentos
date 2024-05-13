package repository

import (
	"context"
	"fiap-tech-challenge-pagamentos/internal/core/domain"
	"fmt"
	"github.com/joomcode/errorx"
	db "github.com/rhuandantas/fiap-tech-challenge-commons/pkg/db/mysql"
	_errors "github.com/rhuandantas/fiap-tech-challenge-commons/pkg/errors"
	"log"
	"xorm.io/xorm"
)

const tableNamePagamento string = "pagamento"

type pagamento struct {
	session *xorm.Session
}

type PagamentoRepo interface {
	Insere(ctx context.Context, pagamento *domain.Pagamento) error
	AtualizaStatus(ctx context.Context, status string, pedidoId string) error
	PesquisaPorPedidoID(ctx context.Context, pedidoId string) (*domain.Pagamento, error)
}

func NewPagamentoRepo(connector db.DBConnector) PagamentoRepo {
	session := connector.GetORM().Table(tableNamePagamento)
	err := connector.SyncTables(new(domain.Pagamento))
	if err != nil {
		log.Fatal(err.Error())
	}

	return &pagamento{
		session: session,
	}
}

func (f *pagamento) Insere(ctx context.Context, pagamento *domain.Pagamento) error {
	_, err := f.session.Context(ctx).Insert(pagamento)
	if err != nil {
		if _errors.IsDuplicatedEntryError(err) {
			return errorx.InternalError.New(fmt.Sprintf("pagamento já existe para pedido %s", pagamento.PedidoId))
		}
		return err
	}

	return nil
}

func (f *pagamento) AtualizaStatus(ctx context.Context, status string, pedidoId string) error {
	_, err := f.session.Context(ctx).Where("pedido_id = ?", pedidoId).Update(&domain.Pagamento{Status: status})
	if err != nil {
		return errorx.InternalError.New(err.Error())
	}

	return nil
}

func (f *pagamento) PesquisaPorPedidoID(ctx context.Context, pedidoId string) (*domain.Pagamento, error) {
	dto := &domain.Pagamento{PedidoId: pedidoId}
	has, err := f.session.Context(ctx).Get(dto)
	if err != nil {
		return nil, err
	}

	if has {
		return dto, nil
	}

	return nil, nil
}
