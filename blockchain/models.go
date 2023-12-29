package blockchain

// Transaction contains the required attributes to submit or evaluate a transaction
type Transaction struct {
	Channel   string   `json:"channel" validate:"required"`
	Chaincode string   `json:"chaincode" validate:"required"`
	Function  string   `json:"function" validate:"required"`
	Arguments []string `json:"arguments"`
}

// TransactionResponse contains the result of a submit or evaluate transaction operation
type TransactionResponse struct {
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}
