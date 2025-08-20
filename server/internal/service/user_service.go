package service

import (
	"context"
	"server/internal/config"
	"server/internal/models"
	"server/internal/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type UserService struct {
	repo    UserRepository
	timeout time.Duration
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo:    repo,
		timeout: config.DefaultTimeout,
	}
}

func (s *UserService) CreateUser(c context.Context, req *models.CreateUserReq) (*models.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	res := &models.CreateUserRes{
		ID:       strconv.FormatInt(r.ID, 10),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *UserService) Login(c context.Context, req *models.LoginUserReq) (*models.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &models.LoginUserRes{}, err
	}

	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		return &models.LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.CustomClaims{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    strconv.FormatInt(user.ID, 10),
		},
	})

	ss, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return &models.LoginUserRes{}, err
	}

	res := &models.LoginUserRes{
		AccessToken: ss,
		ID:          strconv.FormatInt(user.ID, 10),
		Username:    user.Username,
	}

	return res, nil
}
