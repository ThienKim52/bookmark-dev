package user

import "context"

func (s *service) EditCurrentUser(ctx context.Context, userID, displayName, email string) error {
	updates := map[string]interface{}{}
	if displayName != "" {
		updates["display_name"] = displayName
	}
	if email != "" {
		updates["email"] = email
	}
	if len(updates) == 0 {
		return nil
	}
	return s.repo.UpdateUser(ctx, userID, updates)
}