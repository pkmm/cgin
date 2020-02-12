package service

import "testing"

func TestBaseService_BackupMysqlService(t *testing.T) {
	(&backupService{}).BackupMysqlService("")
}
