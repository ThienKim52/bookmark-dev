//go:generate go run github.com/vektra/mockery/v2@latest --name=Hasher --output=./mocks --outpkg=mocks_utils
package utils
import(
	"golang.org/x/crypto/bcrypt"
	"errors"
)

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type hasher struct {}

func NewHasher() Hasher {
	return &hasher{}
}

var ErrHashFailed = errors.New("failed to hash password")
func (h *hasher) Hash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashFailed
	}
	return string(hashBytes), nil
}

func (h *hasher) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}