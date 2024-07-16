package appctx

import "context"

type CUSTOM_CONTEXT string

var (
	USERNAME CUSTOM_CONTEXT = "USERNAME"
)

func GetUsername(ctx context.Context) *string {
	name := ctx.Value(USERNAME)
	if name == nil {
		return nil
	}
	return name.(*string)
}
