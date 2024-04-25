package domain

import (
	"context"
	"time"

	"pastebin/store"
)

func UpdateBinByID(svc store.Store) func(context.Context, string, string, string) error {
	return func(ctx context.Context, id string, contain string, alias string) error {
		bin, err := svc.UpdateBin(ctx, store.Bin{
			ID: id,
			Alias: alias,
			Contain: contain,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		PrintBins(*bin)

		return nil
	}
}
