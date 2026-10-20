package core

type Message struct {
	Role    string
	Content string
}

type ChatRequest struct {
	Model       string
	Messages    []Message
	MaxTokens   *int
	Temperature *float64
}

// Usage Tracks token usage
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type ChatResponse struct {
	Provider     string
	Model        string
	Content      string
	FinishReason string
	Usage        Usage
}
