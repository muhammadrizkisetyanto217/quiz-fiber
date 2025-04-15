package controller

import (
	"errors"
	"log"

	// "os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"

	"quiz-fiber/internals/configs"
	modelAuth "quiz-fiber/internals/features/user/auth/models"
	modelUser "quiz-fiber/internals/features/user/user/models"
)

type AuthController struct {
	DB *gorm.DB
}


func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var input modelUser.UserModel
	if err := c.BodyParser(&input); err != nil {
		log.Printf("[ERROR] Failed to parse request body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}
	if err := input.Validate(); err != nil {
		log.Printf("[ERROR] Validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ERROR] Failed to hash password: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to secure password"})
	}
	input.Password = string(passwordHash)
	if err := ac.DB.Create(&input).Error; err != nil {
		log.Printf("[ERROR] Failed to save user to database: %v", err)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.Status(400).JSON(fiber.Map{"error": "Email already registered"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to register user"})
	}
	log.Printf("[SUCCESS] User registered: ID=%d, Email=%s", input.ID, input.Email)
	return c.Status(201).JSON(fiber.Map{"message": "User registered successfully"})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var input struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		log.Printf("[ERROR] Failed to parse request body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	var user modelUser.UserModel
	if err := ac.DB.Where("email = ? OR user_name = ?", input.Identifier, input.Identifier).First(&user).Error; err != nil {
		log.Printf("[ERROR] User not found: Identifier=%s", input.Identifier)
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email, username, or password"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Printf("[ERROR] Password incorrect for user: %s", user.Email)
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email, username, or password"})
	}
	expirationTime := time.Now().Add(time.Hour * 96)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": expirationTime.Unix(),
	})
	tokenString, err := token.SignedString([]byte(configs.JWTSecret))
	if err != nil {
		log.Printf("[ERROR] Failed to generate token for user: %s", user.Email)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}
	user.Password = ""
	log.Printf("[SUCCESS] User logged in: ID=%d, Email=%s", user.ID, user.Email)
	return c.JSON(fiber.Map{
		"token": tokenString,
		"user": fiber.Map{
			"id":            user.ID,
			"user_name":     user.UserName,
			"email":         user.Email,
			"google_id":     user.GoogleID,
			"role":          user.Role,
			"donation_name": user.DonationName,
			"original_name": user.OriginalName,
			"created_at":    user.CreatedAt,
			"updated_at":    user.UpdatedAt,
		},
	})
}

// 🔥 LOGOUT USER
func (ac *AuthController) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - No token provided"})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - Invalid token format"})
	}
	tokenString := tokenParts[1]

	// Cek apakah token sudah ada di blacklist
	var existingToken modelAuth.TokenBlacklist
	if err := ac.DB.Where("token = ?", tokenString).First(&existingToken).Error; err == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token already blacklisted"})
	}

	// Tambahkan token ke blacklist
	blacklistToken := modelAuth.TokenBlacklist{
		Token:     tokenString,
		ExpiredAt: time.Now().Add(96 * time.Hour), // Sesuai waktu expired token
	}

	if err := ac.DB.Create(&blacklistToken).Error; err != nil {
		log.Printf("[ERROR] Failed to blacklist token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to logout"})
	}

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

// 🔥 CHANGE PASSWORD (Menggunakan c.Locals dan Transaksi)
func (ac *AuthController) ChangePassword(c *fiber.Ctx) error {
	// 🆔 Ambil User ID dari middleware (sudah divalidasi di AuthMiddleware)
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - Invalid token"})
	}

	// 📌 Parsing request body
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// 📌 Validasi input kosong
	if input.OldPassword == "" || input.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Both old and new passwords are required"})
	}

	// 🚨 Cek apakah password baru sama dengan yang lama
	if input.OldPassword == input.NewPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "New password must be different from old password"})
	}

	// 🔍 Cari user di database
	var user modelUser.UserModel
	if err := ac.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// 🔑 Cek apakah password lama cocok
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Old password is incorrect"})
	}

	// 🔒 Hash password baru
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash new password"})
	}

	// 🔥 Update password menggunakan transaksi
	tx := ac.DB.Begin()
	if err := tx.Model(&user).Update("password", string(newHashedPassword)).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update password"})
	}
	tx.Commit()

	// 🎉 Beri response sukses
	return c.JSON(fiber.Map{"message": "Password changed successfully"})
}

// 🔥 CHECK SECURITY ANSWER
func (ac *AuthController) CheckSecurityAnswer(c *fiber.Ctx) error {
	var input struct {
		Email  string `json:"email"`
		Answer string `json:"security_answer"`
	}

	// 📌 Parsing JSON input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// 📌 Cek user berdasarkan email
	var user modelUser.UserModel
	if err := ac.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// 📌 Bandingkan security answer secara langsung
	if strings.TrimSpace(input.Answer) != strings.TrimSpace(user.SecurityAnswer) {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect security answer"})
	}

	// 📌 Response berhasil validasi
	return c.JSON(fiber.Map{
		"message": "Security answer correct",
		"email":   user.Email,
	})
}

// 🔥 RESET PASSWORD
func (ac *AuthController) ResetPassword(c *fiber.Ctx) error {
	var input struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}

	// 📌 Parsing JSON input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// 📌 Cek user berdasarkan email kembali untuk memastikan
	var user modelUser.UserModel
	if err := ac.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// 📌 Hashing password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash new password"})
	}

	// 📌 Update password di database
	if err := ac.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update password"})
	}

	// 📌 Response sukses reset password
	return c.JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}

// 🔥 Middleware untuk proteksi route
func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		log.Println("[DEBUG] Authorization Header:", authHeader)
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - No token provided"})
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - Invalid token format"})
		}
		tokenString := tokenParts[1]
		var existingToken modelAuth.TokenBlacklist
		err := db.Where("token = ?", tokenString).First(&existingToken).Error
		if err == nil {
			log.Println("[WARNING] Token ditemukan di blacklist, akses ditolak.")
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - Token is blacklisted"})
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("[ERROR] Database error saat cek token blacklist:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
		}
		secretKey := configs.JWTSecret
		if secretKey == "" {
			log.Println("[ERROR] JWT_SECRET tidak ditemukan di environment")
			return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error - Missing JWT Secret"})
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			log.Println("[ERROR] Token tidak valid:", err)
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - Invalid token"})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("[ERROR] Token claims tidak valid")
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - Invalid token claims"})
		}
		exp, exists := claims["exp"].(float64)
		if !exists {
			log.Println("[ERROR] Token tidak memiliki exp")
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - Token has no expiration"})
		}
		log.Println("[DEBUG] Token Claims:", claims)

		idStr, exists := claims["id"].(string)
		if !exists {
			log.Println("[ERROR] User ID not found in token claims")
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - No user ID in token"})
		}

		userID, err := uuid.Parse(idStr)
		if err != nil {
			log.Println("[ERROR] Failed to parse UUID from token:", err)
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - Invalid user ID format"})
		}

		c.Locals("user_id", userID)
		log.Println("[SUCCESS] User ID stored in context:", userID)

		expTime := time.Unix(int64(exp), 0)
		log.Printf("[INFO] Token Expiration Time: %v", expTime)
		if time.Now().Unix() > int64(exp) {
			log.Println("[ERROR] Token sudah expired")
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - Token expired"})
		}
		log.Println("[SUCCESS] Token valid, lanjutkan request")
		return c.Next()
	}
}
