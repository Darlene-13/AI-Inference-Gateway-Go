package core

//Package core holds the gateways's own types: what a chat request and response looks like inside the gateway and how errors are classified

type Message struct {
	Role    string //"system", "user", "assistant"
	Content string
}

type ChatRequest struct {
	Model       string    // The alias the client asked for ie: fast
	Messages    []Message // The conversation so far
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
