package service

import "cgin/errno"

type baseService struct {
}

func (b *baseService) CheckError(err error) {
	if err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
	}
}
