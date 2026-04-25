package service

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/model"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/repository"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (dto.UserResponse, error)
	Login(req dto.LoginRequest) (dto.LoginResponse, error)
	ForgotPassword(req dto.ForgotPasswordRequest) error
	ResetPassword(req dto.ResetPasswordRequest) error
}

type authService struct {
	db            *gorm.DB
	userRepo      repository.UserRepository
	userQueryRepo repository.UserQueryRepository
	notifService  NotificationService
	rdb           *redis.Client
}

func NewAuthService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	userQueryRepo repository.UserQueryRepository,
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

func (s *authService) Register(req dto.RegisterRequest) (dto.UserResponse, error) {
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
		return dto.UserResponse{}, util.BadRequestException("Email already use by another account", err)
	}

	newUser := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := s.userRepo.Create(tx, &newUser); err != nil {
		tx.Rollback()
		return dto.UserResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       newUser.ID.String(),
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     newUser.Role,
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	// Define cache key
    cacheKey := "user:email:" + req.Email
    
    // Try to get data from Redis
    u, err := util.GetCache[*model.User](s.rdb, cacheKey)
    
    // If data not found in Redis
    if err != nil {
        // Get data from database
        u, err = s.userQueryRepo.FindByEmail(req.Email)
        if err != nil {
            return dto.LoginResponse{}, util.UnprocessableEntityException("User email " + req.Email + " is not registered")
        }

        // Set data to Redis
		ttl := util.GetEnvTime("REDIS_CACHE_TTL", "5m")
        _ = util.SetCache(s.rdb, cacheKey, u, ttl)
    }

    // Validate password
    err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
    if err != nil {
        return dto.LoginResponse{}, util.UnauthorizedException("Invalid email or password")
    }

    // Generate token
    token, err := util.GenerateToken(u.ID.String(), u.Username, u.Email, u.Picture, u.Role)
    if err != nil {
        return dto.LoginResponse{}, util.UnauthorizedException("Token is invalid or expired")
    }

    return dto.LoginResponse{Token: token}, nil
}

func (s *authService) ForgotPassword(req dto.ForgotPasswordRequest) error {
    // Check email is registered
    u, err := s.userQueryRepo.FindByEmail(req.Email)
    if err != nil {
        return util.UnprocessableEntityException("Email not found")
    }

    // Generate new token for reset password
    token := uuid.New().String()
    cacheKey := "reset_password:" + token

    // Save token to Redis with TTL 5 minutes
	ttl := util.GetEnvTime("REDIS_CACHE_TTL", "5m")
    err = util.SetCache(s.rdb, cacheKey, u.Email, ttl)
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

func (s *authService) ResetPassword(req dto.ResetPasswordRequest) error {
    cacheKey := "reset_password:" + req.Token

    // Get email from Redis based on token
    email, err := util.GetCache[string](s.rdb, cacheKey)
    if err != nil {
        return util.UnauthorizedException("Token invalid or expired")
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
    _ = util.DeleteCache(s.rdb, cacheKey)
    _ = util.DeleteCache(s.rdb, "user:email:"+email)

    return nil
}
