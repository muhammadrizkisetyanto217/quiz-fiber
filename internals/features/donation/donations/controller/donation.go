package controller

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"quiz-fiber/internals/features/donation/donations/model"

	midtrans "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var SnapClient snap.Client

func InitMidtrans() {
	fmt.Println("[MIDTRANS] SnapClient initialized")
	SnapClient.New("SB-Mid-server-tGx-RfyDGj8L5aUi_I17Hqxf", midtrans.Sandbox)
}

type DonationController struct {
	DB *gorm.DB
}

func NewDonationController(db *gorm.DB) *DonationController {
	return &DonationController{DB: db}
}

// ========== BUAT DONASI + SNAP TOKEN ==========
type CreateDonationRequest struct {
	UserID  string `json:"user_id"` // UUID dalam bentuk string
	Amount  int    `json:"amount"`
	Message string `json:"message"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

func (ctrl *DonationController) CreateDonation(c *fiber.Ctx) error {
	var body CreateDonationRequest
	if err := c.BodyParser(&body); err != nil {
		log.Println("[ERROR] Invalid body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 🔐 Validasi & parsing UUID
	userUUID, err := uuid.Parse(body.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id tidak valid"})
	}

	orderID := fmt.Sprintf("DONATION-%d", time.Now().UnixNano())

	donation := model.Donation{
		UserID:  userUUID,
		Amount:  body.Amount,
		Message: body.Message,
		Status:  "pending",
		OrderID: orderID,
	}

	// 💾 Simpan donasi
	if err := ctrl.DB.Create(&donation).Error; err != nil {
		log.Println("[ERROR] Gagal simpan donasi:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan donasi"})
	}

	// 🪙 Generate token Snap Midtrans
	token, err := ctrl.generateSnapToken(donation, body.Name, body.Email)
	if err != nil {
		log.Println("[ERROR] Gagal generate token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat token Midtrans"})
	}

	donation.PaymentToken = token
	if err := ctrl.DB.Save(&donation).Error; err != nil {
		log.Println("[ERROR] Gagal update token:", err)
	}

	return c.JSON(fiber.Map{
		"message":    "Donasi berhasil dibuat",
		"order_id":   donation.OrderID,
		"snap_token": token,
	})
}

// ========== GENERATE SNAP TOKEN ==========
func (ctrl *DonationController) generateSnapToken(d model.Donation, name string, email string) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  d.OrderID,
			GrossAmt: int64(d.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: name,
			Email: email,
		},
	}
	resp, err := SnapClient.CreateTransaction(req)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

// ========== WEBHOOK MIDTRANS ==========
func HandleMidtransNotification(c *fiber.Ctx) error {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		log.Println("[ERROR] Webhook invalid:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid webhook"})
	}

	orderID, ok := body["order_id"].(string)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Missing order_id"})
	}
	transactionStatus := body["transaction_status"].(string)

	db := c.Locals("db").(*gorm.DB)

	var donation model.Donation
	if err := db.Where("order_id = ?", orderID).First(&donation).Error; err != nil {
		log.Println("[ERROR] Order not found:", orderID)
		return c.SendStatus(200) // Tetap OK biar Midtrans gak retry
	}

	switch transactionStatus {
	case "capture", "settlement":
		now := time.Now()
		donation.Status = "paid"
		donation.PaidAt = &now
	case "expire":
		donation.Status = "expired"
	case "cancel":
		donation.Status = "canceled"
	default:
		log.Println("[INFO] Status tidak diproses:", transactionStatus)
	}

	db.Save(&donation)

	log.Println("[SUCCESS] Update status:", donation.OrderID, "→", donation.Status)
	return c.SendStatus(200)
}
