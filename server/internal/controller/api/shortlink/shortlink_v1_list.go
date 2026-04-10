package shortlink

import (
	"context"
	v1 "hotgo/api/api/shortlink"
	"hotgo/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	return service.SysShortLink().ListPub(ctx, req)
}
