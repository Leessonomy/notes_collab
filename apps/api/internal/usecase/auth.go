package usecase

import (
	"context"
	"errors"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/dto"
	"notes-collab-api/internal/utils"
	"time"

	"github.com/google/uuid"
)

type UserRepo interface {
	Create(ctx context.Context, u domain.User) error
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type TokenRepo interface {
	Save(ctx context.Context, t domain.RefreshToken) error
	Refresh(ctx context.Context, oldToken, newToken string, expiresAt time.Time) (string, error)
	Delete(ctx context.Context, token string) error
}

type Session interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken() (string, error)
	AccessExpire() time.Duration
	RefreshExpire() time.Duration
}

type Auth struct {
	userRepo  UserRepo
	tokenRepo TokenRepo
	session   Session
}

func NewAuth(userRepo UserRepo, tokenRepo TokenRepo, session Session) *Auth {
	return &Auth{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		session:   session,
	}
}

func (a *Auth) GetUserSession(ctx context.Context, userID string) (*domain.User, error) {
	return a.userRepo.GetByID(ctx, userID)
}

func (a *Auth) SignUp(ctx context.Context, input dto.SignUpInput) (dto.SessionOutput, error) {
	var output dto.SessionOutput

	input.Email = uniqeEmail(input.Email)

	if err := validateInput(input); err != nil {
		return output, err
	}

	_, err := a.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return output, domain.ErrEmailTaken
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return output, err
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return output, err
	}

	user := domain.User{
		ID:        uuid.NewString(),
		Name:      input.Name,
		Email:     input.Email,
		Password:  hash,
		CreatedAt: time.Now().UTC(),
	}

	if err := a.userRepo.Create(ctx, user); err != nil {
		return output, err
	}

	output, err = a.newSession(ctx, user.ID)
	if err != nil {
		return output, err
	}

	output.User = &user

	return output, nil
}

func (a *Auth) LogIn(ctx context.Context, input dto.LogInInput) (dto.SessionOutput, error) {
	var output dto.SessionOutput

	input.Email = uniqeEmail(input.Email)

	if err := validateInput(input); err != nil {
		return output, err
	}

	user, err := a.userRepo.GetByEmail(ctx, input.Email)
	if err != nil || !utils.CheckPassword(user.Password, input.Password) {
		return output, domain.ErrInvalidCredentials
	}

	output, err = a.newSession(ctx, user.ID)
	if err != nil {
		return output, err
	}

	output.User = user

	return output, nil
}

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (dto.SessionOutput, error) {
	var output dto.SessionOutput

	newRefreshToken, err := a.session.GenerateRefreshToken()
	if err != nil {
		return output, err
	}

	now := time.Now().UTC()
	refreshExpiresAt := now.Add(a.session.RefreshExpire())

	userID, err := a.tokenRepo.Refresh(ctx, refreshToken, newRefreshToken, refreshExpiresAt)
	if err != nil {
		return output, err
	}

	accessToken, err := a.session.GenerateAccessToken(userID)
	if err != nil {
		return output, err
	}

	return dto.SessionOutput{
		AccessToken:      accessToken,
		RefreshToken:     newRefreshToken,
		AccessExpiresAt:  now.Add(a.session.AccessExpire()),
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

func (a *Auth) LogOut(ctx context.Context, refreshToken string) error {
	return a.tokenRepo.Delete(ctx, refreshToken)
}

func (a *Auth) newSession(ctx context.Context, userID string) (dto.SessionOutput, error) {
	var output dto.SessionOutput

	accessToken, err := a.session.GenerateAccessToken(userID)
	if err != nil {
		return output, err
	}

	refreshToken, err := a.session.GenerateRefreshToken()
	if err != nil {
		return output, err
	}

	now := time.Now().UTC()
	accessExpiresAt := now.Add(a.session.AccessExpire())
	refreshExpiresAt := now.Add(a.session.RefreshExpire())

	err = a.tokenRepo.Save(ctx, domain.RefreshToken{
		Token:     refreshToken,
		UserID:    userID,
		ExpiresAt: refreshExpiresAt,
		CreatedAt: now,
	})
	if err != nil {
		return output, err
	}

	return dto.SessionOutput{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}
