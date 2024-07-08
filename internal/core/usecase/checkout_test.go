package usecase

import (
	"context"
	"errors"
	"fiap-tech-challenge-pagamentos/internal/core/domain"
	mock_client "fiap-tech-challenge-pagamentos/test/mock/client"
	mock_repo "fiap-tech-challenge-pagamentos/test/mock/repository"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

var _ = Describe("realiza checkout use case testes", func() {
	var (
		ctx            = context.Background()
		repo           *mock_repo.MockPagamentoRepo
		pedidoClient   *mock_client.MockPedido
		producaoClient *mock_client.MockProducao
		checkoutUC     RealizarCheckout
	)

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		repo = mock_repo.NewMockPagamentoRepo(mockCtrl)
		pedidoClient = mock_client.NewMockPedido(mockCtrl)
		producaoClient = mock_client.NewMockProducao(mockCtrl)
		checkoutUC = NewRealizaCheckout(repo, pedidoClient, producaoClient)
	})

	Context("realiza checkout", func() {
		objID := primitive.NewObjectID()
		pagamentoReturn := &domain.Pagamento{
			Id:       1,
			Status:   "aprovado",
			PedidoId: objID.Hex(),
			Tipo:     domain.TipoPix,
			Valor:    10.0,
		}
		It("faz checkout com sucesso", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(nil, nil)
			pedidoClient.EXPECT().AtualizaStatus(gomock.Any(), "pagamento_aprovado", objID.Hex()).Return(nil)
			producaoClient.EXPECT().AdicionaFila(gomock.Any(), gomock.Any()).Return(nil)
			repo.EXPECT().Insere(ctx, pagamentoReturn).Return(nil)
			err := checkoutUC.Checkout(ctx, pagamentoReturn)

			gomega.Expect(err).To(gomega.BeNil())
		})
		It("falha ao pesquisar por pedido", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(nil, errors.New("mock error"))
			err := checkoutUC.Checkout(ctx, pagamentoReturn)

			gomega.Expect(err.Error()).To(gomega.Equal("mock error"))
		})
		It("falha com pedido ja existente", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(pagamentoReturn, nil)
			err := checkoutUC.Checkout(ctx, pagamentoReturn)

			gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
		})
		It("atualiza pedido falha", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(nil, nil)
			pedidoClient.EXPECT().AtualizaStatus(ctx, "pagamento_aprovado", objID.Hex()).Return(errors.New("mock error"))
			err := checkoutUC.Checkout(ctx, pagamentoReturn)

			gomega.Expect(err.Error()).To(gomega.Equal("mock error"))
		})
		It("insere pagamento falha", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(nil, nil)
			pedidoClient.EXPECT().AtualizaStatus(ctx, "pagamento_aprovado", objID.Hex()).Return(nil)
			producaoClient.EXPECT().AdicionaFila(ctx, gomock.Any()).Return(nil)
			repo.EXPECT().Insere(ctx, pagamentoReturn).Return(errors.New("mock error"))
			err := checkoutUC.Checkout(ctx, pagamentoReturn)

			gomega.Expect(err.Error()).To(gomega.Equal("mock error"))
		})
		It("falha ao adicionar na fila", func() {
			repo.EXPECT().PesquisaPorPedidoID(ctx, objID.Hex()).Return(nil, nil)
			pedidoClient.EXPECT().AtualizaStatus(ctx, "pagamento_aprovado", objID.Hex()).Return(nil)
			producaoClient.EXPECT().AdicionaFila(ctx, gomock.Any()).Return(errors.New("mock error"))
			err := checkoutUC.Checkout(ctx, pagamentoReturn)

			gomega.Expect(err.Error()).To(gomega.Equal("mock error"))
		})
	})
})
