package store

import (
	"context"
	"time"
)

type Bin struct {
	ID        string    `json:"id"`
	Alias     string    `json:"alias"`
	Contain   string    `json:"contain"`
	Clic      int32     `json:"clic"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Statistics struct {
	binNumber int32 `json:"bin_number"`
	clicByBin any   `json:"clic_by_bin"`
}

type Store interface {
	CreateBin(ctx context.Context, task Bin) (*Bin, error)
	GetBinByAlias(ctx context.Context, alias string) (*Bin, error)
	GetAllBins(ctx context.Context) ([]Bin, error)
	GetStats(ctx context.Context) (*Statistics, error)
	UpdateBin(ctx context.Context, task Bin) (*Bin, error)
	DeleteBinByID(ctx context.Context, id string) (*Bin, error)
}