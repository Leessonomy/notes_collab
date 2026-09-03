package usecase

import (
	"fmt"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/utils"
	"strings"
)

func validateInput(input any) error {
	if err := utils.ValidateStruct(input); err != nil {
		return fmt.Errorf("%w: %s", domain.ErrValidation, err)
	}
	return nil
}

func uniqeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
