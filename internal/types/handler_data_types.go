package types

type JsonPostRequestData struct {
	URL string `json:"url"`
}
type JsonPostResponseData struct {
	Result string `json:"result"`
}

type BatchJsonPostRequestData struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchJsonPostResponseData struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
