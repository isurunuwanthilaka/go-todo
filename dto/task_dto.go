package dto

type Task struct {
	Id          uint   `json:"id"`
	Description string `json:"description"`
	Title       string `json:"title"`
}
