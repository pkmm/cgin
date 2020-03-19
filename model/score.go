package model

import "cgin/conf"

type Score struct {
	Id        uint64  `json:"id" gorm:"primary_key;auto_increment;"`
	StudentId uint64  `json:"student_id" gorm:"index"`
	Xn        string  `json:"xn" gorm:"size:20;not null"`
	Xq        uint8   `json:"xq"`
	Kcmc      string  `json:"kcmc" gorm:"size:50;not null"`
	Type      string  `json:"type" gorm:"size:200;not null"`
	Xf        float64 `json:"xf"`
	Jd        float64 `json:"jd"`
	Cj        string  `json:"cj" gorm:"size:200;not null"`
	Bkcj      string  `json:"bkcj" gorm:"size:200"`
	Cxcj      string  `json:"cxcj" gorm:"size:200"`
	// todo 提取课程的唯一编码，主键
	Model
}

func BatchCreateScores(scores []*Score) {
	if len(scores) == 0 {
		return
	}
	// TODO: 使用协程 或者 是线程池来做
	for _, score := range scores {
		_, _ = score.UpdateOrCreate()
	}
}

func (s *Score) UpdateOrCreate() (err error, _s *Score) {
	err = conf.DB.Where("xn = ? and xq = ? and kcmc = ?", s.Xn, s.Xq, s.Kcmc).
		Assign(*s).
		FirstOrCreate(&s).Error
	return err, s
}

// model method.
func GetScoresByStudentId(studentId uint64) (err error, s []*Score) {
	err = conf.DB.Find(&s, Score{StudentId: studentId}).Error
	return
}