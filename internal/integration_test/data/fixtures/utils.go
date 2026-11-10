package fixtures

import (
	"time"

	"github.com/ThienKim52/golang-dev/internal/app/model"
)

var TestTime = time.Date(2026, 10, 4, 21, 3, 10, 0, time.UTC) 

func GetTestBase(id string) model.Base {
	return model.Base{
		ID: id,
		CreatedAt: TestTime,
		UpdatedAt: TestTime,
	}
}