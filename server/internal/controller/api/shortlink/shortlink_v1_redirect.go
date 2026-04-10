package shortlink

import (
	"context"
	v1 "hotgo/api/api/shortlink"
	"hotgo/internal/service"
)

func (c *ControllerV1) Redirect(ctx context.Context, req *v1.RedirectReq) (res *v1.RedirectRes, err error) {
	return service.SysShortLink().Redirect(ctx, req.Code)
}
