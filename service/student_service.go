package service

import (
	"cgin/model"
)

type studentService struct {
}

var StudentService = &studentService{}

func (s *studentService) GetStudentNeedSyncScore(
	offset, limit int) (students []*model.Student, err error) {
	if err := db.Model(&model.Student{}).
		Where("is_sync = 1").
		Offset(offset).
		Limit(limit).
		Order("id DESC").
		Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

func (s *studentService) UpdateStudentSyncStatus(studentId uint64, syncStatus bool) error {
	if err := db.Model(&model.Student{}).
		Where("id = ?", studentId).
		UpdateColumn("is_sync", syncStatus).Error; err != nil {
		return err
	}
	return nil
}

func (s *studentService) RestSyncStatus() error {
	if err := db.Model(&model.Student{}).
		Update("is_sync", true).Error; err != nil {
		return err
	}
	return nil
}

func (s *studentService) UpdateStudentName(studentId uint64, name string) error {
	if err := db.Model(&model.Student{}).
		Where("id = ?", studentId).
		Update("name", name).Error; err != nil {
		return err
	}
	return nil
}

func (s *studentService) GetScores(studentId uint64) (scores []*model.Score, err error) {
	if err = db.Model(&model.Score{}).
		Where("student_id = ?", studentId).
		Find(&scores).Error; err != nil {
		return nil, err
	}
	return scores, nil
}

func (s *studentService) GetScoreCount(studentId uint64) (count uint64) {
	if err := db.Model(&model.Score{}).
		Where("student_id = ?", studentId).
		Count(&count).Error; err != nil {
		return 0
	}
	return count
}

func (s *studentService) UpdateSyncDetail(syncDetail *model.SyncDetail) *model.SyncDetail {
	if err := db.Where(&model.SyncDetail{StudentId: syncDetail.StudentId}).
		Assign(model.SyncDetail{
			StudentNumber: syncDetail.StudentNumber,
			CostTime:      syncDetail.CostTime,
			Info:          syncDetail.Info,
			Count:         syncDetail.Count,
		}).
		FirstOrCreate(&syncDetail).Error; err != nil {
		return nil
	}
	return syncDetail
}
