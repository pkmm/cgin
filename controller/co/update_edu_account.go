package co

type EduAccount struct {
	AuthCredit
	StudentNumber string `json:"student_number" example:"1923"`
	Password      string `json:"password" example:"34"`
}
