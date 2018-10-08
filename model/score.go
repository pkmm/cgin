package model

type Score struct {
	Id        int64   `json:"id"`
	StudentId int64   `json:"-" gorm:"-"`
	Xn        string  `json:"xn"`
	Xq        uint8   `json:"xq"`
	Kcmc      string  `json:"kcmc"`
	Type      string  `json:"type"`
	Xf        float32 `json:"xf"`
	Jd        float32 `json:"jd"`
	Cj        string  `json:"cj"`
	Bkcj      string  `json:"bkcj"`
	Cxcj      string  `json:"cxcj"`
	CreatedAt string  `json:"created_at" gorm:"default:null"`
	UpdatedAt string  `json:"updated_at" gorm:"default:null"`
}

func init() {
	//db.AutoMigrate(&Score{})
}

func (Score) TableName() string {
	return "score"
}

func GetScore(studentId int64) []Score {
	scores := make([]Score, 0)
	if studentId > 0 {
		db.Where("stu_id = ?", studentId).Find(&scores)
	}
	return scores
}
