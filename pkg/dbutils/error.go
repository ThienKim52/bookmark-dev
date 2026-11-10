package dbutils

import (
	"errors"
	"strings"
)

var errorFilters = []func(err error) (bool, error) {
	filterDuplicationUsername,
	filterDuplicationEmail,
	filterRecordNotFound,
	filterUniqueConstraint,
	filterBookmarkNotFound,
}
func CatchDBError(err error) error {
	if err == nil {
		return nil
	}

	for _, filter := range errorFilters {
		match, filteredErr := filter(err)
		if match {
			return filteredErr
		}
	}
	return err
}

var (
	ErrDuplicationUsername = errors.New("username already exists") 
	ErrDuplicationEmail = errors.New("email already exists")
	ErrRecordNotFound = errors.New("record not found")
	ErrUniqueConstraint = errors.New("unique constraint")
	ErrBookmarkNotFound = errors.New("bookmark not found")
)

func filterDuplicationUsername(err error) (bool, error) {
    if strings.Contains(strings.ToLower(err.Error()), "uni_users_username") {
        return true, ErrDuplicationUsername
    }
    return false, nil
}

func filterDuplicationEmail(err error) (bool, error) {
    if strings.Contains(strings.ToLower(err.Error()), "uni_users_email") {
        return true, ErrDuplicationEmail
    }
    return false, nil
}

func filterRecordNotFound(err error) (bool, error) {
    if strings.Contains(strings.ToLower(err.Error()), "record not found") {
        return true, ErrRecordNotFound
    }
    return false, nil
}

func filterUniqueConstraint(err error) (bool, error) {
	errLower := strings.ToLower(err.Error())
	if strings.Contains(errLower, "unique constraint") || strings.Contains(errLower, "duplicate key") || strings.Contains(errLower, "uni_code") {
		return true, ErrUniqueConstraint
	}
	return false, nil
}

func filterBookmarkNotFound(err error) (bool, error) {
	if strings.Contains(strings.ToLower(err.Error()), "bookmark not found") {
		return true, ErrBookmarkNotFound
	}
	return false, nil
}
