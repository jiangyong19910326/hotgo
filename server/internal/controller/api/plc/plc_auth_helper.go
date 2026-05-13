package plc

import (
	"context"

	"hotgo/internal/library/contexts"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
)

func currentFrontUserId(ctx context.Context) (int64, error) {
	user := contexts.GetUser(ctx)
	if user == nil || user.Id <= 0 {
		return 0, gerror.New("未登录")
	}
	return user.Id, nil
}

func ensureDeviceAccess(ctx context.Context, deviceId int) error {
	userId, err := currentFrontUserId(ctx)
	if err != nil {
		return err
	}
	ok, err := service.FrontUser().CanAccessDevice(ctx, userId, deviceId)
	if err != nil {
		return err
	}
	if !ok {
		return gerror.New("无权访问该设备")
	}
	return nil
}

func ensurePointAccess(ctx context.Context, pointId int) error {
	userId, err := currentFrontUserId(ctx)
	if err != nil {
		return err
	}
	ok, err := service.FrontUser().CanAccessPoint(ctx, userId, pointId)
	if err != nil {
		return err
	}
	if !ok {
		return gerror.New("无权访问该点位")
	}
	return nil
}
