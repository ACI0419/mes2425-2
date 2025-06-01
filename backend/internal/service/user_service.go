package service

import (
	"errors"
	"mes-system/internal/models"
	"mes-system/pkg/jwt"
	"mes-system/pkg/utils"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	db        *gorm.DB
	jwtConfig *jwt.JWTConfig
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB, jwtConfig *jwt.JWTConfig) *UserService {
	return &UserService{
		db:        db,
		jwtConfig: jwtConfig,
	}
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token    string      `json:"token"`
	UserInfo models.User `json:"user_info"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	RealName string `json:"real_name" binding:"required"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Login 用户登录
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	var user models.User
	err := s.db.Where("username = ? AND status = ?", req.Username, 1).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(s.jwtConfig, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 清除密码字段
	user.Password = ""

	return &LoginResponse{
		Token:    token,
		UserInfo: user,
	}, nil
}

// Register 用户注册
func (s *UserService) Register(req *RegisterRequest) (*models.User, error) {
	// 检查用户名是否已存在
	var count int64
	s.db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	s.db.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 设置默认角色
	if req.Role == "" {
		req.Role = "user"
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		RealName: req.RealName,
		Phone:    req.Phone,
		Role:     req.Role,
		Status:   1,
	}

	err = s.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	// 清除密码字段
	user.Password = ""
	return &user, nil
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := s.db.Where("id = ? AND status = ?", userID, 1).First(&user).Error
	if err != nil {
		return nil, err
	}

	// 清除密码字段
	user.Password = ""
	return &user, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	var user models.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return err
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	return s.db.Model(&user).Update("password", hashedPassword).Error
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(page, pageSize int, keyword string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	// 清除密码字段
	for i := range users {
		users[i].Password = ""
	}

	return users, total, nil
}