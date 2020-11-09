package payment

type CardData struct {
	Number string `json:"number"`
	Holder string `json:"holder"`
	CVC    int    `json:"cvc"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
}
