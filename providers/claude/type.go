package claude

import "one-api/types"

type ClaudeError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ClaudeMetadata struct {
	UserId string `json:"user_id"`
}

type ResContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type ClaudeRequest struct {
	Model             string   `json:"model"`
	System            string   `json:"system,omitempty"`
	Messages          string   `json:"messages"`
	MaxTokens 	  int      `json:"max_tokens"`
	StopSequences     []string `json:"stop_sequences,omitempty"`
	Temperature       float64  `json:"temperature,omitempty"`
	TopP              float64  `json:"top_p,omitempty"`
	TopK              int      `json:"top_k,omitempty"`
	//ClaudeMetadata    `json:"metadata,omitempty"`
	Stream bool `json:"stream,omitempty"`
}

type ClaudeResponseError struct {
	Error ClaudeError `json:"error,omitempty"`
}
type ClaudeResponse struct {
	Content	   	[]ResContent 	`json:"content"`
	Id	   	string		`json:"id"`
	Role	   	string		`json:"role"`
	StopReason 	string       	`json:"stop_reason"`
	StopSequence	string		`json:"stop_sequence,omitempty"`
	Model      	string      	`json:"model"`
	Usage      	*types.Usage 	`json:"usage,omitempty"`
	ClaudeResponseError
}
