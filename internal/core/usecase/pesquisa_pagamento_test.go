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

var _ = Describe("pesquisa pagamento use case testes", func() {
	var (
		ctx                 = context.Background()
		repo                *mock_repo.MockPagamentoRepo
		pesquisaPagamentoUC PesquisaPagamento
	)

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		repo = mock_repo.NewMockPagamentoRepo(mockCtrl)
		pesquisaPagamentoUC = NewPesquisaPagamento(repo)
	})

	Context("pesquisa pagamento", func() {
		objID := primitive.NewObjectID()
		pagamentoReturn := &domain.Pagamento{
			Id:       1,
			Status:   "aprovado",
			PedidoId: objID.Hex(),
			Tipo:     domain.TipoPix,
			Valor:    10.0,
		}
		It("faz pesquisa com sucesso", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(pagamentoReturn, nil)
			pagamento, err := pesquisaPagamentoUC.PesquisaPorPedidoID(ctx, objID.Hex())

			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(pagamento).ToNot(gomega.BeNil())
		})
		It("não encontra pedido", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(nil, nil)
			pagamento, err := pesquisaPagamentoUC.PesquisaPorPedidoID(ctx, objID.Hex())

			gomega.Expect(pagamento).To(gomega.BeNil())
			gomega.Expect(err).ToNot(gomega.BeNil())
		})
		It("falha ao pesquisar pedido", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(nil, errors.New("erro ao pesquisar pedido"))
			pagamento, err := pesquisaPagamentoUC.PesquisaPorPedidoID(ctx, objID.Hex())

			gomega.Expect(err.Error()).To(gomega.Equal("erro ao pesquisar pedido"))
			gomega.Expect(pagamento).To(gomega.BeNil())
		})
	})
})
