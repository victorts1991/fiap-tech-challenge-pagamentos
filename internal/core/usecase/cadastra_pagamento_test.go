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

var _ = Describe("cadastra pagamento use case testes", func() {
	var (
		ctx                 = context.Background()
		repo                *mock_repo.MockPagamentoRepo
		cadastraPagamentoUC CadastrarPagamento
	)

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		repo = mock_repo.NewMockPagamentoRepo(mockCtrl)
		cadastraPagamentoUC = NewCadastraPagamento(repo)
	})

	Context("cadastra pagamento", func() {
		objID := primitive.NewObjectID()
		pagamentoReturn := &domain.Pagamento{
			Id:       1,
			Status:   "pendente",
			PedidoId: objID.Hex(),
			Tipo:     domain.TipoPix,
			Valor:    10.0,
		}
		It("cadastra pagamento com sucesso", func() {
			repo.EXPECT().Insere(ctx, pagamentoReturn).Return(nil)
			err := cadastraPagamentoUC.Cadastra(ctx, pagamentoReturn)

			gomega.Expect(err).To(gomega.BeNil())
		})
		It("insere pagamento falha", func() {
			repo.EXPECT().Insere(ctx, pagamentoReturn).Return(errors.New("mock error"))
			err := cadastraPagamentoUC.Cadastra(ctx, pagamentoReturn)

			gomega.Expect(err.Error()).To(gomega.Equal("mock error"))
		})
	})
})
