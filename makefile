mock:
	mockgen -source=internal/adapters/repository/pagamento.go -package=mock_repo -destination=test/mock/repository/pagamento.go
	mockgen -source=client/pedido.go -package=mock_client -destination=test/mock/client/pedido.go
	mockgen -source=internal/core/usecase/atualiza_pagamento.go -package=mock_usecase -destination=test/mock/usecase/atualiza_pagamento.go
install:
	go install github.com/swaggo/swag/cmd/swag@latest
generate-swagger: install
	swag init