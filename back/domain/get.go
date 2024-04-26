package domain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"pastebin/store"
)

func ListBins(svc store.Store) func(context.Context) error {
	return func(ctx context.Context) error {
		bins, err := svc.GetAllBins(ctx)
		if err != nil {
			return errors.Wrap(err, "couldnt get all bins")
		}

		PrintBins(bins...)

		return nil
	}
}

// func GetStats(svc store.Store) func(context.Context) error {
// 	return func(ctx context.Context) error {
// 		bins, err := svc.GetStats(ctx)
// 		if err != nil {
// 			return errors.Wrap(err, "couldnt get all bins")
// 		}

		// à faire : coder un PrintStats(stats...)
// 		// PrintStats(stats...)

// 		return nil
// 	}
// }

func GetBinByAlias(svc store.Store) func(context.Context, string) error {
	return func(ctx context.Context, binID string) error {
		bin, err := svc.GetBinByAlias(ctx, binID)
		if err != nil {
			return errors.Wrapf(err, "couldnt get bin with %s", binID)
		}

		PrintBins(*bin)

		return nil
	}
}

func PrintBins(bins ...store.Bin) {
	numberOfBins := len(bins)

	if numberOfBins == 0 {
		fmt.Println("no bins to print")
		return
	}

	fmt.Printf("printing %d bins", numberOfBins)
	for _, bin := range bins {
		fmt.Println(jsonify(bin))
	}
}

func jsonify(data store.Bin) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
