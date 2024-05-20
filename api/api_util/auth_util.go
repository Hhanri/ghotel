package api_util

import (
	"context"
	"fmt"

	"github.com/hhanri/ghotel/types"
)

func GetAuth(ctx context.Context) (*types.User, error) {
	user, ok := ctx.Value("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("Unauthorized")
	}
	return user, nil
}
