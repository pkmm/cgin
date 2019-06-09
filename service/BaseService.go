package service

import "cgin/errno"

type BaseService struct {
}

func (b *BaseService) CheckError(err error) {
	if err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
	}
}

