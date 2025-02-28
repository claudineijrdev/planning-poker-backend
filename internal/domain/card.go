package domain

type Result struct {
	Average      float64     `json:"average"`
	Distribution map[int]int `json:"distribution"`
}

type Card struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Votes       []int  `json:"votes"`
	Result      Result `json:"result"`
	Closed      bool   `json:"closed"`
}

type Vote struct {
	Score int `json:"score"`
}
