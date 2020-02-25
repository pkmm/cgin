package service

import (
	"cgin/errno"
	"cgin/model"
	"cgin/util"
)

type hermannService struct {
}

var HermannService = &hermannService{}

var (
	//RecordNotFoundError = errors.New("未找到相关的任务计划记录")
	//TaskNotBegin        = errors.New("还没有到计划开始的时间")
	//TaskHasDone         = errors.New("计划的任务已经结束了")
	reviewDayAt = []int{0, 1, 3, 7, 14}
)

// 背诵的单元的开始结束区间
type UnitInterval struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// 每日的任务
type TodayTask struct {
	Remember   *UnitInterval  `json:"remember"`
	ReviewList []UnitInterval `json:"review_list"`
	CurrentDay int            `json:"current_day"`
}

// 任务到今天已经过去的天数, 每天的单位数量
type HoursPast struct {
	Days  int `json:"days"`
	Unit  int `json:"unit"`
	Total int `json:"total"`
}

// 计算今天背诵以及复习的任务
// TODO: 设置总共的单元的数量，然后计算复习的单元
func (h *hermannService) GetTodayTask(userId uint64) (*TodayTask, error) {
	record := h.GetHoursPastOfTask(userId)
	if record == nil {
		return nil, errno.RememberRecordNotFound
	}
	if record.Days < 0 {
		return nil, errno.RememberTaskNotBegin
	}

	// 单词复习结束的日期，在这个日期之后的几天中每天复习unit个单元
	reviewEndAt := record.Total/record.Unit + reviewDayAt[len(reviewDayAt)-1]
	// 任务结束的日期
	taskEndAt := reviewEndAt + record.Total/record.Unit + 1
	if taskEndAt < record.Days {
		return nil, errno.RememberTaskHasDone
	}
	var result TodayTask
	result.CurrentDay = record.Days
	result.Remember = nil
	list := make([]UnitInterval, 0)
	result.ReviewList = list

	// 百科上的这一段是3天的空档期（本实现中不在留有空白的日期）
	//if record.Days > reviewEndAt && record.Days <= reviewEndAt + 3 {
	//	return &result, nil
	//}
	//if record.Days > reviewEndAt + 3 {
	//	result.Remember = h.getNthRememberList(record.Days - reviewEndAt - 3, record.Unit)
	//	return &result, nil
	//}

	if record.Days > reviewEndAt {
		result.Remember = h.getNthRememberList(record.Days-reviewEndAt, record.Unit)
		return &result, nil
	}

	result.Remember = h.getNthRememberList(record.Days, record.Unit)
	if result.Remember.Start > record.Total {
		result.Remember = nil
	}
	for _, day := range reviewDayAt {
		if day >= record.Days {
			break
		}
		review := h.getNthRememberList(record.Days-day, record.Unit)
		if review.Start > record.Total {
			continue
		}
		list = append(list, *review)
	}
	result.ReviewList = list
	return &result, nil
}

// 第n天背诵的单词
func (h *hermannService) getNthRememberList(n, unit int) *UnitInterval {
	return &UnitInterval{n*unit - unit + 1, n * unit}
}

func (h *hermannService) GetTaskRecord(userId uint64) *model.HermannMemorial {
	var result *model.HermannMemorial
	if err := db.Model(&model.HermannMemorial{}).
		Where("user_id = ?", userId).
		First(&result).Error; err != nil {
		// TODO: LOG something.
		return nil
	}
	return result
}

func (h *hermannService) GetHoursPastOfTask(userId uint64) *HoursPast {
	var result HoursPast
	if err := db.
		Raw("SELECT DATEDIFF(NOW(), start_at) + 1 AS days, remember_unit AS unit, total_unit as total FROM hermann_memorials WHERE user_id = ?", userId).
		Scan(&result).Error; err != nil {
		return nil
	}
	return &result
}

func (h *hermannService) SaveTask(unit, totalUnit uint, startAt util.JSONTime, userId uint64) error {
	var result model.HermannMemorial
	if err := db.Model(&model.HermannMemorial{}).
		Where("user_id = ?", userId).
		Assign(model.HermannMemorial{RememberUnit: unit, StartAt: startAt, UserId: userId, TotalUnit: totalUnit}).
		FirstOrCreate(&result).Error; err != nil {
		return err
	}
	return nil
}
