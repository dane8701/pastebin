package domain

import (
	"context"

	"github.com/pkg/errors"

	"pastebin/store"
)

func DeleteBinByID(svc store.Store) func(context.Context, string) error {
	return func(ctx context.Context, binID string) error {
		bin, err := svc.DeleteBinByID(ctx, binID)
		if err != nil {
			return errors.Wrapf(err, "couldnt delete bin with %s", binID)
		}

		PrintBins(*bin)

		return nil
	}
}
