package api

type Params struct {
	Status  int    `json:"status"`
	States  []int  `json:"states"`
	UserId  uint64 `json:"user_id"`
	Type    int    `json:"type"`
	Types   []int  `json:"types"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}
