package services

import (
	"e-shop-api/internal/constants"
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/models"
	"e-shop-api/internal/pkg/utils"
	"e-shop-api/internal/repositories"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dtos.RegisterRequest) (dtos.UserResponse, error)
	Login(req dtos.LoginRequest) (dtos.LoginResponse, error)
	RefreshToken(req dtos.RefreshTokenRequest) (dtos.RefreshTokenResponse, error)
	ForgotPassword(req dtos.ForgotPasswordRequest) error
	ResetPassword(req dtos.ResetPasswordRequest) error
}

type authService struct {
	db            *gorm.DB
	userRepo      repositories.UserRepository
	userQueryRepo repositories.UserQueryRepository
	notifService  NotificationService
	rdb           *redis.Client
}

func NewAuthService(
	db *gorm.DB,
	userRepo repositories.UserRepository,
	userQueryRepo repositories.UserQueryRepository,
	notifService NotificationService,
	rdb *redis.Client,
) AuthService {
	return &authService{
		db,
		userRepo,
		userQueryRepo,
		notifService,
		rdb,
	}
}

func (s *authService) Register(req dtos.RegisterRequest) (dtos.UserResponse, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	u, err := s.userQueryRepo.FindByEmail(req.Email)
	if err == nil && u != nil {
		tx.Rollback()
		return dtos.UserResponse{}, utils.BadRequestException("Email already use by another account", err)
	}

	newUser := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := s.userRepo.Create(tx, &newUser); err != nil {
		tx.Rollback()
		return dtos.UserResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dtos.UserResponse{}, err
	}

	return dtos.UserResponse{
		ID:       newUser.ID.String(),
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     newUser.Role,
	}, nil
}

func (s *authService) Login(req dtos.LoginRequest) (dtos.LoginResponse, error) {
	cacheKey := "user:email:" + req.Email

	u, err := utils.GetCache[*models.User](s.rdb, cacheKey)

	if err != nil {
		u, err = s.userQueryRepo.FindByEmail(req.Email)
		if err != nil {
			return dtos.LoginResponse{}, utils.UnprocessableEntityException("User email " + req.Email + " is not registered")
		}

		ttl := utils.GetEnvTime("REDIS_CACHE_TTL", constants.RedisCacheTtl)
		_ = utils.SetCache(s.rdb, cacheKey, u, ttl)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return dtos.LoginResponse{}, utils.UnauthorizedException("Invalid email or password")
	}

	accessToken, err := utils.GenerateAccessToken(u.ID.String(), u.Username, u.Email, u.Picture, u.Role)
	if err != nil {
		return dtos.LoginResponse{}, utils.UnauthorizedException("Failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(u.ID.String())
	if err != nil {
		return dtos.LoginResponse{}, utils.UnauthorizedException("Failed to generate refresh token")
	}

	refreshTTL := utils.GetEnvTime("JWT_REFRESH_TTL", constants.JwtRefreshTtl)
	redisKey := "refresh_token:" + u.ID.String()
	_ = utils.SetCache(s.rdb, redisKey, refreshToken, refreshTTL)

	accessTTL := utils.GetEnvTime("JWT_ACCESS_TTL", constants.JwtAccessTtl)
	intAccessTTL := int64(accessTTL.Seconds())
	expiresIn := intAccessTTL

	return dtos.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    constants.BearerPrefix,
	}, nil
}

func (s *authService) RefreshToken(req dtos.RefreshTokenRequest) (dtos.RefreshTokenResponse, error) {
	claims, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return dtos.RefreshTokenResponse{}, utils.UnauthorizedException("Invalid refresh token")
	}

	redisKey := "refresh_token:" + claims.UserID
	storedToken, err := utils.GetCache[string](s.rdb, redisKey)
	if err != nil || storedToken != req.RefreshToken {
		return dtos.RefreshTokenResponse{}, utils.UnauthorizedException("Refresh token revoked or expired")
	}

	u, err := s.userQueryRepo.FindByID(claims.UserID)
	if err != nil {
		return dtos.RefreshTokenResponse{}, utils.UnauthorizedException("User not found")
	}

	accessToken, err := utils.GenerateAccessToken(u.ID.String(), u.Username, u.Email, u.Picture, u.Role)
	if err != nil {
		return dtos.RefreshTokenResponse{}, utils.UnauthorizedException("Failed to generate access token")
	}

	newRefreshToken, err := utils.GenerateRefreshToken(u.ID.String())
	if err != nil {
		return dtos.RefreshTokenResponse{}, utils.UnauthorizedException("Failed to generate refresh token")
	}

	refreshTTL := utils.GetEnvTime("JWT_REFRESH_TTL", constants.JwtRefreshTtl)
	_ = utils.SetCache(s.rdb, redisKey, newRefreshToken, refreshTTL)

	expiresIn := int64(utils.GetEnvTime("JWT_ACCESS_TTL", constants.JwtAccessTtl).Seconds())

	return dtos.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    constants.BearerPrefix,
	}, nil
}

func (s *authService) ForgotPassword(req dtos.ForgotPasswordRequest) error {
	// Check email is registered
	u, err := s.userQueryRepo.FindByEmail(req.Email)
	if err != nil {
		return utils.UnprocessableEntityException("Email not found")
	}

	// Generate new token for reset password
	token := uuid.New().String()
	cacheKey := "reset_password:" + token

	// Save token to Redis with TTL 5 minutes
	ttl := utils.GetEnvTime("REDIS_CACHE_TTL", constants.RedisCacheTtl)
	err = utils.SetCache(s.rdb, cacheKey, u.Email, ttl)
	if err != nil {
		return err
	}

	// Set email body
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
	emailBody := fmt.Sprintf(`
        <h1>Reset Password Request</h1>
        <p>Halo %s, we received a request to reset your password.</p>
        <p>Click the link below or copy and paste it into your browser (link expires in 5 minutes):</p>
        <a href="%s" style="background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Reset Password</a>
        <p>If you did not make this request, please ignore this email.</p>
    `, u.Username, resetLink)

	// Send reset password email
	s.notifService.QueueSendEmail(
		u.Email,
		"Reset Password",
		emailBody,
	)

	return nil
}

func (s *authService) ResetPassword(req dtos.ResetPasswordRequest) error {
	cacheKey := "reset_password:" + req.Token

	// Get email from Redis based on token
	email, err := utils.GetCache[string](s.rdb, cacheKey)
	if err != nil {
		return utils.UnauthorizedException("Token invalid or expired")
	}

	// Get data user by email
	u, err := s.userQueryRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	// Update password
	if err := s.userRepo.Update(s.db, u); err != nil {
		return err
	}

	// Delete token and email from Redis
	_ = utils.DeleteCache(s.rdb, cacheKey)
	_ = utils.DeleteCache(s.rdb, "user:email:"+email)

	// Revoke all refresh tokens (force logout all devices)
	_ = utils.DeleteCache(s.rdb, "refresh_token:"+u.ID.String())

	return nil
}
