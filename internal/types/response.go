package types

type MainResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Data        any    `json:"data"`
}
