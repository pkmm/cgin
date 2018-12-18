package model

import "time"

type Score struct {
	Id        int64     `json:"id"`
	StudentId int64     `json:"-" gorm:"column:student_id"`
	Xn        string    `json:"xn"`
	Xq        uint8     `json:"xq"`
	Kcmc      string    `json:"kcmc"`
	Type      string    `json:"type"`
	Xf        float64   `json:"xf"`
	Jd        float64   `json:"jd"`
	Cj        string    `json:"cj"`
	Bkcj      string    `json:"bkcj"`
	Cxcj      string    `json:"cxcj"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func init() {
	//db.AutoMigrate(&Score{})
}

func (Score) TableName() string {
	return "scores"
}

func GetScore(studentId int64) []Score {
	scores := make([]Score, 0)
	if studentId > 0 {
		db.Where("student_id = ?", studentId).Order("id desc").Find(&scores)
	}
	return scores
}

func UpdateOrCreateScore(score *Score) *Score {
	db.Where(Score{StudentId: score.StudentId, Xn: score.Xn, Xq: score.Xq, Kcmc: score.Kcmc}).
		Assign(score).
		FirstOrCreate(&score)
	return score
}
