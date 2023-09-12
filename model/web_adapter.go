package model

import "context"

type IWebAdapter interface {
	Run(ctx context.Context) error
}
