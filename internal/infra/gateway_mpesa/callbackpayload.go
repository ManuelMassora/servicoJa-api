package gatewaympesa

type MpesaCallbackPayload struct {
	ResponseCode         string `json:"output_ResponseCode"`
	ResponseDesc         string `json:"output_ResponseDesc"`
	TransactionID        string `json:"output_TransactionID"`
	ConversationID       string `json:"output_ConversationID"`
	ThirdPartyReference  string `json:"output_ThirdPartyReference"`
}