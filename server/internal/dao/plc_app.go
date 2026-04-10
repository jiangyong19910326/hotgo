package dao

import (
	"hotgo/internal/dao/internal"
)

type plcAppDao struct {
	*internal.PlcAppDao
}

var (
	PlcApp = plcAppDao{internal.NewPlcAppDao()}
)
