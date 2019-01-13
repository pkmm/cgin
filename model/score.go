package model

type Score struct {
	ID     uint64  `json:"id" gorm:"primary_key"`
	UserId uint64  `json:"user_id" gorm:"index"`
	Xn     string  `json:"xn" gorm:"size:20;not null"`
	Xq     uint8   `json:"xq"`
	Kcmc   string  `json:"kcmc" gorm:"size:50;not null"`
	Type   string  `json:"type" gorm:"size:200;not null"`
	Xf     float64 `json:"xf"`
	Jd     float64 `json:"jd"`
	Cj     string  `json:"cj" gorm:"size:200;not null"`
	Bkcj   string  `json:"bkcj" gorm:"size:200"`
	Cxcj   string  `json:"cxcj" gorm:"size:200"`

	Model
}

//func GetScore(studentId int64) []Score {
//	scores := make([]Score, 0)
//	if studentId > 0 {
//		db.Where("student_id = ?", studentId).Order("id desc").Find(&scores)
//	}
//	return scores
//}
//
//func UpdateOrCreateScore(score *Score) *Score {
//	db.Where(Score{StudentId: score.StudentId, Xn: score.Xn, Xq: score.Xq, Kcmc: score.Kcmc}).
//		Assign(score).
//		FirstOrCreate(&score)
//	return score
//}
