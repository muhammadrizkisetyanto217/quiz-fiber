package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"quiz-fiber/internals/features/users/user/models"

	"gorm.io/gorm"
)

// * Kita membuat sebuah struct bernama UserController, yang memiliki satu property bernama DB. (Property adalah variabel yang terdapat dalam sebuah struct).
// & Property DB ini adalah pointer ke objek database (gorm.DB), yang akan digunakan untuk mengakses database.
type UserController struct {
	DB *gorm.DB
}

//^ Bayangkan UserController ini seperti seorang kasir toko.
// 1. Agar bisa bekerja, kasir butuh akses ke database toko (misalnya, daftar barang dan harga).
// 2. Dalam kode ini, DB adalah akses ke database yang diberikan ke kasir (UserController).
// 3. Tanpa DB, kasir tidak bisa mencari barang, menambahkan transaksi, dll.

// *  Fungsi NewUserController adalah "constructor"
// Constructor ini digunakan untuk membuat objek UserController dengan database yang bisa disesuaikan.
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// 1. Saat Anda mempekerjakan kasir baru (UserController), Anda harus memberi mereka akses ke database toko (DB).
// 2. NewUserController(db) adalah cara memberi kasir akses ke database saat mereka mulai bekerja.

// GET all users
func (uc *UserController) GetUsers(c *fiber.Ctx) error {
	var users []models.UserModel
	if err := uc.DB.Find(&users).Error; err != nil {
		log.Println("[ERROR] Failed to fetch users:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}

	log.Printf("[SUCCESS] Retrieved %d users\n", len(users))
	return c.JSON(fiber.Map{
		"message": "Users fetched successfully",
		"total":   len(users),
		"data":    users,
	})
}

// GET user by ID
func (uc *UserController) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("[ERROR] Invalid UUID format:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	id := userID
	var user models.UserModel

	if err := uc.DB.First(&user, id).Error; err != nil {
		log.Println("[ERROR] User not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"message": "User fetched successfully",
		"data":    user,
	})
}

// POST create new user(s)
func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	var singleUser models.UserModel
	var multipleUsers []models.UserModel

	// Coba parse sebagai array terlebih dahulu
	if err := c.BodyParser(&multipleUsers); err == nil && len(multipleUsers) > 0 {
		if err := uc.DB.Create(&multipleUsers).Error; err != nil {
			log.Println("[ERROR] Failed to create multiple users:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create multiple users"})
		}
		return c.Status(201).JSON(fiber.Map{
			"message": "Users created successfully",
			"data":    multipleUsers,
		})
	}

	// Jika gagal diparse sebagai array, parse sebagai satu user
	if err := c.BodyParser(&singleUser); err != nil {
		log.Println("[ERROR] Invalid input format:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input format"})
	}

	if err := uc.DB.Create(&singleUser).Error; err != nil {
		log.Println("[ERROR] Failed to create user:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    singleUser,
	})
}

// PUT update user by ID
func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.UserModel

	if err := uc.DB.First(&user, id).Error; err != nil {
		log.Println("[ERROR] User not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if err := c.BodyParser(&user); err != nil {
		log.Println("[ERROR] Invalid input for update:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := uc.DB.Save(&user).Error; err != nil {
		log.Println("[ERROR] Failed to update user:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	log.Printf("[SUCCESS] User updated: ID=%d\n", user.ID)
	return c.JSON(fiber.Map{
		"message": "User updated successfully",
		"data":    user,
	})
}

// DELETE user by ID
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := uc.DB.Delete(&models.UserModel{}, id).Error; err != nil {
		log.Println("[ERROR] Failed to delete user:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	log.Printf("[SUCCESS] User with ID %s deleted\n", id)
	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
