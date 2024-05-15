package handlers

import (
	"errors"
	"fiap-tech-challenge-pagamentos/internal/core/commons"
	"fiap-tech-challenge-pagamentos/internal/core/domain"
	"fiap-tech-challenge-pagamentos/internal/core/usecase"
	"fmt"
	"github.com/go-openapi/swag"
	"github.com/joomcode/errorx"
	"github.com/labstack/echo/v4"
	serverErr "github.com/rhuandantas/fiap-tech-challenge-commons/pkg/errors"
	"github.com/rhuandantas/fiap-tech-challenge-commons/pkg/middlewares/auth"
	"github.com/rhuandantas/fiap-tech-challenge-commons/pkg/util"
	"net/http"
	"strings"
)

type Pagamento struct {
	pesquisaPorPedidoID usecase.PesquisaPagamento
	realizaCheckoutUC   usecase.RealizarCheckout
	validator           util.Validator
	tokenJwt            auth.Token
}

func NewPagamento(pesquisaPorPedidoID usecase.PesquisaPagamento, validator util.Validator,
	tokenJwt auth.Token,
	realizaCheckoutUC usecase.RealizarCheckout) *Pagamento {
	return &Pagamento{
		pesquisaPorPedidoID: pesquisaPorPedidoID,
		validator:           validator,
		tokenJwt:            tokenJwt,
		realizaCheckoutUC:   realizaCheckoutUC,
	}
}

func (h *Pagamento) RegistraRotasPagamento(server *echo.Echo) {
	server.GET("/pagamento/:pedido_id", h.pesquisaPorPedidoId)
	server.POST("/pagamento/checkout/:pedidoId", h.checkout)
}

// pesquisaPorPedidoId godoc
// @Summary pega um pagamento por pedido id
// @Tags Pagamento
// @Accept */*
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param        pedido_id   path      string  true  "id do pedido"
// @Success 200 {object} domain.Pagamento
// @Router /pagamento/{pedido_id} [get]
func (h *Pagamento) pesquisaPorPedidoId(ctx echo.Context) error {
	pedidoID := ctx.Param("pedido_id")

	cliente, err := h.pesquisaPorPedidoID.PesquisaPorPedidoID(ctx.Request().Context(), pedidoID)
	if err != nil {
		return serverErr.HandleError(ctx, errorx.Cast(err))
	}
	return ctx.JSON(http.StatusOK, cliente)
}

func validatePagamento(d *domain.Pagamento) error {
	if !swag.ContainsStrings(domain.TiposValidos, strings.ToLower(d.Tipo)) {
		return errors.New(fmt.Sprintf("%s tipo de pagamento invalido", d.Tipo))
	}

	if d.Valor <= 0.0 {
		return errors.New(fmt.Sprint("valor invalido"))
	}

	if d.PedidoId == "" {
		return errors.New(fmt.Sprint("pedido id invalido"))
	}

	return nil
}

// checkout godoc
// @Summary checkout do pedido
// @Tags Pedido
// @Accept json
// @Success 200 {object} commons.MessageResponse
// @Param        pedidoId   path      integer  true  "id do pedido a ser feito o checkout"
// @Param        id   body      domain.Pagamento  true  "status permitido: aprovado | recusado"
// @Produce json
// @Router /pagamento/checkout/{pedidoId} [post]
func (h *Pagamento) checkout(ctx echo.Context) error {
	var (
		pagamento domain.Pagamento
		err       error
	)

	id := ctx.Param("pedidoId")

	if err = ctx.Bind(&pagamento); err != nil {
		return serverErr.HandleError(ctx, serverErr.BadRequest.New(err.Error()))
	}

	pagamento.PedidoId = id
	err = validatePagamento(&pagamento)
	if err != nil {
		return serverErr.HandleError(ctx, serverErr.BadRequest.New(err.Error()))
	}

	if !checkoutStatusIsValid(pagamento.Status) {
		return serverErr.HandleError(ctx, serverErr.BadRequest.New(fmt.Sprintf("%s não é um status válido para checkout", pagamento.Status)))
	}

	err = h.realizaCheckoutUC.Checkout(ctx.Request().Context(), &pagamento)
	if err != nil {
		return serverErr.HandleError(ctx, errorx.Cast(err))
	}

	return ctx.JSON(http.StatusOK, commons.MessageResponse{Message: "checkout realizado com sucesso"})
}

func checkoutStatusIsValid(status string) bool {
	return status == domain.StatusAprovado ||
		status == domain.StatusRecusado
}
