package types

type JSONPostRequestData struct {
	URL string `json:"url"`
}
type JSONPostResponseData struct {
	Result string `json:"result"`
}

type BatchJSONPostRequestData struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchJSONPostResponseData struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
