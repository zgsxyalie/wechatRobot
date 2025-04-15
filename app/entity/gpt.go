package entity

type GptResp struct {
	Id       string        `json:"id"`
	Provider string        `json:"provider"`
	Model    string        `json:"model"`
	Object   string        `json:"object"`
	Created  int           `json:"created"`
	Choices  []RespChoices `json:"choices"`
	Usage    RespUsage     `json:"usage"`
}

type RespChoices struct {
	Logrobs            interface{}        `json:"logprobs"`
	FinishReason       string             `json:"finish_reason"`
	NativeFinishReason string             `json:"native_finish_reason"`
	Index              int                `json:"index"`
	Message            RespChoicesMessage `json:"message"`
}

type RespChoicesMessage struct {
	Role    string      `json:"role"`
	Content string      `json:"content"`
	Refusal interface{} `json:"refusal"`
}

type RespUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
