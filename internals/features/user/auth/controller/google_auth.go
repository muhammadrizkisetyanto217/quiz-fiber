package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"quiz-fiber/internals/configs"
	"quiz-fiber/internals/features/user/user/models"
)

// GoogleConfig holds the configuration for Google OAuth
type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
}

// GoogleUserInfo holds user info from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// GoogleAuthController handles Google authentication
type GoogleAuthController struct {
	DB           *gorm.DB
	GoogleConfig GoogleConfig
}

// NewGoogleAuthController creates a new GoogleAuthController
func NewGoogleAuthController(db *gorm.DB) *GoogleAuthController {
	return &GoogleAuthController{
		DB: db,
		GoogleConfig: GoogleConfig{
			ClientID:     configs.GetEnv("GOOGLE_CLIENT_ID"),
			ClientSecret: configs.GetEnv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  configs.GetEnv("GOOGLE_REDIRECT_URL"),
			Scopes:       []string{"profile", "email"},
			AuthURL:      "https://accounts.google.com/o/oauth2/auth",
			TokenURL:     "https://oauth2.googleapis.com/token",
			UserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
		},
	}
}

// GoogleLogin initiates the Google OAuth flow
func (gc *GoogleAuthController) GoogleLogin(c *fiber.Ctx) error {
	log.Println("[INFO] Starting Google Login process")

	// Generate random state
	state := generateRandomString(32)

	// Store state in session
	c.Cookie(&fiber.Cookie{
		Name:     "google_oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
	})

	// Build the authorization URL
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		gc.GoogleConfig.AuthURL,
		gc.GoogleConfig.ClientID,
		url.QueryEscape(gc.GoogleConfig.RedirectURL),
		url.QueryEscape("email profile"),
		state,
	)

	log.Printf("[INFO] Redirecting to Google auth URL: %s", authURL)
	return c.Redirect(authURL)
}

// GoogleCallback handles the callback from Google
func (gc *GoogleAuthController) GoogleCallback(c *fiber.Ctx) error {
	log.Println("[INFO] Handling Google callback")

	// 1. Get and validate code and state
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		log.Println("[ERROR] No code received from Google!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No code received from Google",
		})
	}
	
	// Verify state
	storedState := c.Cookies("google_oauth_state")
	if storedState == "" || storedState != state {
		log.Println("[ERROR] Invalid state parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid state parameter",
		})
	}

	// Clear the state cookie
	c.Cookie(&fiber.Cookie{
		Name:     "google_oauth_state",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	// 2. Exchange code for token
	tokenData := url.Values{}
	tokenData.Set("code", code)
	tokenData.Set("client_id", gc.GoogleConfig.ClientID)
	tokenData.Set("client_secret", gc.GoogleConfig.ClientSecret)
	tokenData.Set("redirect_uri", gc.GoogleConfig.RedirectURL)
	tokenData.Set("grant_type", "authorization_code")

	tokenReq, _ := http.NewRequest("POST", gc.GoogleConfig.TokenURL, 
		io.NopCloser(strings.NewReader(tokenData.Encode())))
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	client := &http.Client{}
	resp, err := client.Do(tokenReq)
	if err != nil {
		log.Printf("[ERROR] Failed to exchange code for token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to exchange code for token",
		})
	}
	defer resp.Body.Close()

	// 3. Parse token response
	tokenBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read token response: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read token response",
		})
	}
	
	log.Printf("[DEBUG] Google token response: %s", string(tokenBody))

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}

	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		log.Printf("[ERROR] Failed to parse token response: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse token response",
		})
	}

	if tokenResp.AccessToken == "" {
		log.Printf("[ERROR] No access token received from Google")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No access token received from Google",
		})
	}

	// 4. Get user info
	userReq, _ := http.NewRequest("GET", gc.GoogleConfig.UserInfoURL, nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	userResp, err := client.Do(userReq)
	if err != nil {
		log.Printf("[ERROR] Failed to get user info: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user info",
		})
	}
	defer userResp.Body.Close()

	// 5. Parse user info
	userData, err := io.ReadAll(userResp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read user info response: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read user info response",
		})
	}

	log.Printf("[DEBUG] Raw Google User Info Response: %s", string(userData))

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(userData, &userInfo); err != nil {
		log.Printf("[ERROR] Failed to parse user info: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user info",
		})
	}

	// Validasi data user dari Google
	if userInfo.ID == "" || userInfo.Email == "" {
		log.Printf("[ERROR] Invalid user info from Google: ID=%s, Email=%s", userInfo.ID, userInfo.Email)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user info from Google",
		})
	}

	log.Printf("[INFO] User info from Google: ID=%s, Name=%s, Email=%s", userInfo.ID, userInfo.Name, userInfo.Email)

	// 6. Check if user exists and create/update as needed
	// Kita akan gunakan transaction untuk memastikan integritas data
	var user models.UserModel
	tx := gc.DB.Begin()
	
	if tx.Error != nil {
		log.Printf("[ERROR] Failed to begin transaction: %v", tx.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Coba cari user berdasarkan GoogleID
	result := tx.Where("google_id = ?", userInfo.ID).First(&user)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		tx.Rollback()
		log.Printf("[ERROR] Database error when searching by Google ID: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	if result.RowsAffected == 0 {
		// User tidak ditemukan berdasarkan Google ID, coba cari berdasarkan email
		result = tx.Where("email = ?", userInfo.Email).First(&user)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tx.Rollback()
			log.Printf("[ERROR] Database error when searching by email: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		if result.RowsAffected == 0 {
			// User tidak ditemukan, buat user baru
			log.Printf("[INFO] Creating new user for Google ID: %s, Email: %s", userInfo.ID, userInfo.Email)

			// Generate random password
			randomPassword := generateRandomString(12)
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
			if err != nil {
				tx.Rollback()
				log.Printf("[ERROR] Failed to hash password: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to secure password",
				})
			}

			// Setup user data
			googleID := userInfo.ID
			userName := userInfo.Name
			if userName == "" {
				userName = userInfo.Email[:strings.Index(userInfo.Email, "@")]
			}

			newUser := models.UserModel{
				UserName:         userName,
				Email:            userInfo.Email,
				Password:         string(hashedPassword),
				GoogleID:         &googleID,
				Role:             "user",
				SecurityQuestion: "Account created with Google Auth",
				SecurityAnswer:   "google_auth_user",
				OriginalName:     &userName,
			}

			// Debug info
			userJSON, _ := json.Marshal(newUser)
			log.Printf("[DEBUG] New user data: %s", string(userJSON))

			// Buat user baru
			if err := tx.Create(&newUser).Error; err != nil {
				tx.Rollback()
				log.Printf("[ERROR] Failed to create user: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create user account",
					"details": err.Error(),
				})
			}

			// Pastikan user sudah di-commit ke database
			if err := tx.Commit().Error; err != nil {
				log.Printf("[ERROR] Failed to commit transaction: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to save user account",
				})
			}

			// Ambil user yang baru dibuat untuk memastikan data lengkap
			if err := gc.DB.First(&user, newUser.ID).Error; err != nil {
				log.Printf("[ERROR] Failed to retrieve newly created user: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to retrieve user account",
				})
			}

			log.Printf("[SUCCESS] New user created: ID=%d, Email=%s", user.ID, user.Email)
		} else {
			// User ditemukan berdasarkan email, update Google ID
			log.Printf("[INFO] Updating existing user with Google ID: %s", userInfo.ID)
			googleID := userInfo.ID
			user.GoogleID = &googleID
			
			if err := tx.Save(&user).Error; err != nil {
				tx.Rollback()
				log.Printf("[ERROR] Failed to update user with Google ID: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to update user account",
				})
			}
			
			if err := tx.Commit().Error; err != nil {
				log.Printf("[ERROR] Failed to commit transaction: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to save user account",
				})
			}
			
			log.Printf("[SUCCESS] Updated user with Google ID: ID=%d, Email=%s", user.ID, user.Email)
		}
	} else {
		// User ditemukan berdasarkan Google ID
		log.Printf("[INFO] User found by Google ID: %s", userInfo.ID)
		
		// Update user info jika diperlukan
		if user.Email != userInfo.Email {
			user.Email = userInfo.Email
			
			if err := tx.Save(&user).Error; err != nil {
				tx.Rollback()
				log.Printf("[ERROR] Failed to update user email: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to update user account",
				})
			}
			
			if err := tx.Commit().Error; err != nil {
				log.Printf("[ERROR] Failed to commit transaction: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to save user account",
				})
			}
			
			log.Printf("[SUCCESS] Updated user email: ID=%d, Email=%s", user.ID, user.Email)
		} else {
			// Tidak ada perubahan, commit transaction
			if err := tx.Commit().Error; err != nil {
				log.Printf("[ERROR] Failed to commit transaction: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Database error",
				})
			}
		}
	}

	// Generate JWT token
	expirationTime := time.Now().Add(time.Hour * 96) // 4 days
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": expirationTime.Unix(),
	})

	tokenString, err := token.SignedString([]byte(configs.GetEnv("JWT_SECRET")))
	if err != nil {
		log.Printf("[ERROR] Failed to generate token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Return token and user info
	log.Printf("[SUCCESS] Google login successful for user: ID=%d, Email=%s", user.ID, user.Email)
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

// generateRandomString generates a random string of the specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // To ensure uniqueness
	}
	return string(result)
}
