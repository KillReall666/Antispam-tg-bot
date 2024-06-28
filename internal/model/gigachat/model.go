package gigachat

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatModel struct {
	Model             string    `json:"model"`
	Messages          []Message `json:"messages"`
	Temperature       float64   `json:"temperature"`
	TopP              float64   `json:"top_p"`
	N                 int       `json:"n"`
	Stream            bool      `json:"stream"`
	MaxTokens         int       `json:"max_tokens"`
	RepetitionPenalty float64   `json:"repetition_penalty"`
}

type Choice struct {
	Message      Message `json:"message"`
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
}

type Completion struct {
	Choices []Choice `json:"choices"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
