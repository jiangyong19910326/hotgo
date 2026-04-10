package shortlink

import (
	"context"
	v1 "hotgo/api/api/shortlink"
	"hotgo/internal/service"
)

func (c *ControllerV1) Stats(ctx context.Context, req *v1.StatsReq) (res *v1.StatsRes, err error) {
	return service.SysShortLink().Stats(ctx, req.Code)
}
