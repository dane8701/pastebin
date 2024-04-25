package domain

import (
	"context"

	"pastebin/store"
)

func CreateBin(svc store.Store) func(context.Context, string, string) error {
	return func(ctx context.Context, contain string, alias string) error {
		bin, err := svc.CreateBin(ctx, store.Bin{
			Alias: alias,
			Contain: contain,
		})
		if err != nil {
			return err
		}

		PrintBins(*bin)

		return nil
	}
}
