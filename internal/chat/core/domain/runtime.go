package domain

// RuntimeSelection represents the requested LLM runtime selection.
type RuntimeSelection struct {
	Provider string `json:"provider" example:"openai"`
	Model    string `json:"model"    example:"gpt-5.4-mini"`
}
