package usecase

import (
	"context"
	"errors"
	"fiap-tech-challenge-pagamentos/internal/core/domain"
	mock_repo "fiap-tech-challenge-pagamentos/test/mock/repository"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

var _ = Describe("atualiza pagamento use case testes", func() {
	var (
		ctx               = context.Background()
		repo              *mock_repo.MockPagamentoRepo
		atualizaPagamento AtualizaPagamento
	)

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		repo = mock_repo.NewMockPagamentoRepo(mockCtrl)
		atualizaPagamento = NewAtualizaPagamento(repo)
	})

	Context("atualiza pagamento", func() {
		objID := primitive.NewObjectID()
		pagamentoReturn := &domain.Pagamento{
			Id:       1,
			Status:   "aprovado",
			PedidoId: objID.Hex(),
			Tipo:     domain.TipoPix,
			Valor:    10.0,
		}
		It("atualiza com sucesso", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(pagamentoReturn, nil)
			repo.EXPECT().AtualizaStatus(ctx, "aprovado", objID.Hex()).Return(nil)
			err := atualizaPagamento.Atualiza(ctx, "aprovado", objID.Hex())

			gomega.Expect(err).To(gomega.BeNil())
		})
		It("falha na pesquisa por pedido", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(nil, errors.New("mock error"))
			err := atualizaPagamento.Atualiza(ctx, "aprovado", objID.Hex())
			gomega.Expect(err.Error()).To(gomega.Equal("mock error"))
		})
		It("falha ao atualizar pagamento", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(pagamentoReturn, nil)
			repo.EXPECT().AtualizaStatus(ctx, "aprovado", objID.Hex()).Return(errors.New("mock error"))
			err := atualizaPagamento.Atualiza(ctx, "aprovado", objID.Hex())

			gomega.Expect(err.Error()).To(gomega.Equal("mock error"))
		})
		It("falha ao atualizar pagamento", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(nil, nil)
			err := atualizaPagamento.Atualiza(ctx, "aprovado", objID.Hex())

			gomega.Expect(err).ToNot(gomega.BeNil())
		})
	})
})
