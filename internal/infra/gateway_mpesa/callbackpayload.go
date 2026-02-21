package gatewaympesa

type MpesaCallbackPayload struct {
	ResponseCode              string `json:"output_ResponseCode"`
	ResponseDesc              string `json:"output_ResponseDesc"`
	TransactionID             string `json:"output_TransactionID"`
	ConversationID            string `json:"output_ConversationID"`
	ThirdPartyReference       string `json:"output_ThirdPartyReference"`
	ResponseTransactionStatus string `json:"output_ResponseTransactionStatus"`
}

type MpesaQueryResponse struct {
	ResponseCode              string `json:"output_ResponseCode"`
	ResponseDesc              string `json:"output_ResponseDesc"`
	ResponseResult            string `json:"output_ResponseResult"`
	TransactionID             string `json:"output_TransactionID"`
	ConversationID            string `json:"output_ConversationID"`
	ThirdPartyReference       string `json:"output_ThirdPartyReference"`
	TransactionStatus         string `json:"output_TransactionStatus"`
	ResponseTransactionStatus string `json:"output_ResponseTransactionStatus"`
}
