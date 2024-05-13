package domain

import (
	"time"
)

const (
	StatusAprovado = "aprovado"
	StatusRecusado = "recusado"
	StatusPendente = "pendente"
	TipoDinheiro   = "dinheiro"
	TipoCredito    = "credito"
	TipoDebito     = "debito"
	TipoPix        = "pix"
)

var TiposValidos = []string{TipoDinheiro, TipoPix, TipoCredito, TipoDebito}

type Pagamento struct {
	Id        int64     `json:"id" xorm:"pk autoincr 'pagamento_id'"`
	PedidoId  string    `json:"pedido_id" validate:"required" xorm:"index unique"`
	Status    string    `json:"status"`
	Tipo      string    `json:"tipo" validate:"required, oneof=credito,debito,pix,dinheiro"`
	Valor     float32   `json:"valor"`
	CreatedAt time.Time `xorm:"created"`
	Update    time.Time `xorm:"updated"`
}
