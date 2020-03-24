package entity

// Point オブジェクト構造
type Point struct {
	ID      uint   `json:"id"`
	UserID  uint   `json:"userId"`
	Number  int    `json:"number"`
	Comment string `json:"comment"`
}
