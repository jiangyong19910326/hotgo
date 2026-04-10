package shortlink

import (
	"context"
	v1 "hotgo/api/api/shortlink"
	"hotgo/internal/service"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	return service.SysShortLink().Create(ctx, req)
}
