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
	UserId    User      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Statistics struct {
	BinNumber int32 		`json:"bin_number"`
	ClicByBin []ClicByBin  `json:"clic_by_bin"`
}

type ClicByBin struct {
	BinID string `json:"bin_id"`
	Clic int32   `json:"clic"`
}

type User struct {
	ID        		string    	`json:"id"`
	Email     		string    	`json:"email"`
	MotDePasse      string    	`json:"mot_de_passe"`
}

type Store interface {
	CreateBin(ctx context.Context, task Bin) (*Bin, error)
	GetBinByAlias(ctx context.Context, alias string) (*Bin, error)
	GetAllBins(ctx context.Context) ([]Bin, error)
	GetStats(ctx context.Context) (*Statistics, error)
	UpdateBin(ctx context.Context, task Bin) (*Bin, error)
	DeleteBinByID(ctx context.Context, id string) (*Bin, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user User) (*User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	DropAllUsers(ctx context.Context) error
}