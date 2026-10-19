//go:generate go run github.com/vektra/mockery/v2@latest --name=Repository --output=./mocks --outpkg=mocks_userrepo
package user
import (
	"context"
	"github.com/ThienKim52/golang-dev/internal/app/model"
)
type Repository interface{
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error 
}