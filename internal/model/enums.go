package model

// Status geral de serviço, pagamento ou vaga
type Status string

const (
	StatusPendente   	Status = "PENDENTE"
	StatusEmAndamento 	Status = "EM_ANDAMENTO"
	StatusConcluido   	Status = "CONCLUIDO"
	StatusCancelado   	Status = "CANCELADO"
	StatusDisponivel  	Status = "DISPONIVEL"
	StatusOcupada  	  	Status = "OCUPADA"
	StatusAceito	  	Status = "ACEITO"
	StatusRejeitado	  	Status = "REJEITADO"
)

// Tipo de movimento financeiro
type TipoMovimento string

const (
	TipoCredito TipoMovimento = "CREDITO"
	TipoDebito  TipoMovimento = "DEBITO"
)

// Método de pagamento
type MetodoPagamento string

const (
	MetodoMPesa  MetodoPagamento = "M_PESA"
	MetodoCarteira MetodoPagamento = "CARTEIRA"
	MetodoOutro   MetodoPagamento = "OUTRO"
)

// Função (Role) do usuário no sistema
type Role string

const (
	RoleCliente   Role = "CLIENTE"
	RolePrestador Role = "PRESTADOR"
	RoleAdmin     Role = "ADMIN"
)
