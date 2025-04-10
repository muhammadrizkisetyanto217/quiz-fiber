package services

import (
	"errors"
	"log"

	userCategoryModel "quiz-fiber/internals/features/category/category/model"
	userSubcategoryModel "quiz-fiber/internals/features/category/subcategory/model"
	userThemesOrLevelsModel "quiz-fiber/internals/features/category/themes_or_levels/model"
	userUnitModel "quiz-fiber/internals/features/category/units/model"
	quizzesModel "quiz-fiber/internals/features/quizzes/quizzes/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func UpdateUserSectionIfQuizCompleted(db *gorm.DB, userID uuid.UUID, sectionID int, quizID int) error {
	log.Println("[SERVICE] UpdateUserSectionIfQuizCompleted - userID:", userID, "sectionID:", sectionID, "quizID:", quizID)

	// Hitung semua quiz yang aktif di section tersebut sebagai total
	var allQuizzes []quizzesModel.QuizModel
	if err := db.Where("section_quizzes_id = ? AND deleted_at IS NULL", sectionID).Find(&allQuizzes).Error; err != nil {
		log.Println("[ERROR] Failed to fetch quizzes for section:", err)
		return err
	}

	totalQuizIDs := pq.Int64Array{}
	for _, quiz := range allQuizzes {
		totalQuizIDs = append(totalQuizIDs, int64(quiz.ID))
	}

	var userSection quizzesModel.UserSectionQuizzesModel
	err := db.Where("user_id = ? AND section_quizzes_id = ?", userID, sectionID).First(&userSection).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userSection = quizzesModel.UserSectionQuizzesModel{
			UserID:           userID,
			SectionQuizzesID: sectionID,
			CompleteQuiz:     pq.Int64Array{int64(quizID)},
			TotalQuiz:        totalQuizIDs,
		}
		log.Println("[SERVICE] Creating new UserSectionQuizzesModel")
		return db.Create(&userSection).Error
	}

	// Cek apakah quizID sudah tercatat
	for _, id := range userSection.CompleteQuiz {
		if id == int64(quizID) {
			log.Println("[SERVICE] Quiz ID already recorded, skipping update")
			return nil
		}
	}

	userSection.CompleteQuiz = append(userSection.CompleteQuiz, int64(quizID))
	userSection.TotalQuiz = totalQuizIDs

	log.Println("[SERVICE] Updating existing UserSectionQuizzesModel")
	return db.Save(&userSection).Error
}

func UpdateUserUnitIfSectionCompleted(db *gorm.DB, userID uuid.UUID, unitID uint, sectionID uint) error {
	log.Println("[SERVICE] UpdateUserUnitIfSectionCompleted - userID:", userID, "unitID:", unitID, "sectionID:", sectionID)

	// Ambil semua section dari unit (untuk TotalSectionQuizzes)
	var allSections []quizzesModel.SectionQuizzesModel
	if err := db.Where("unit_id = ? AND deleted_at IS NULL", unitID).Find(&allSections).Error; err != nil {
		log.Println("[ERROR] Failed to fetch sections for unit:", err)
		return err
	}
	totalSectionIDs := pq.Int64Array{}
	for _, section := range allSections {
		totalSectionIDs = append(totalSectionIDs, int64(section.ID))
	}

	// Cek apakah semua quiz dalam sectionID sudah diselesaikan
	var userSection quizzesModel.UserSectionQuizzesModel
	err := db.Where("user_id = ? AND section_quizzes_id = ?", userID, sectionID).First(&userSection).Error
	if err != nil {
		log.Println("[ERROR] UserSectionQuizzesModel not found, skipping unit update")
		return nil // tidak buat user_unit kalau belum ada section progress
	}

	if len(userSection.CompleteQuiz) < len(userSection.TotalQuiz) {
		log.Println("[INFO] Section belum selesai, tidak update ke UserUnitModel")
	}

	// Cek user_unit (buat kalau belum ada)
	var userUnit userUnitModel.UserUnitModel
	err = db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userUnit = userUnitModel.UserUnitModel{
			UserID:                 userID,
			UnitID:                 unitID,
			CompleteSectionQuizzes: pq.Int64Array{},
			TotalSectionQuizzes:    totalSectionIDs,
		}
		log.Println("[SERVICE] Creating new UserUnitModel")
		if len(userSection.CompleteQuiz) == len(userSection.TotalQuiz) {
			userUnit.CompleteSectionQuizzes = pq.Int64Array{int64(sectionID)}
		}
		return db.Create(&userUnit).Error
	}

	// Jika section lengkap, dan belum masuk array → tambahkan
	if len(userSection.CompleteQuiz) == len(userSection.TotalQuiz) {
		for _, id := range userUnit.CompleteSectionQuizzes {
			if id == int64(sectionID) {
				log.Println("[SERVICE] Section ID already recorded in UserUnitModel")
				return nil
			}
		}
		userUnit.CompleteSectionQuizzes = append(userUnit.CompleteSectionQuizzes, int64(sectionID))
	}

	userUnit.TotalSectionQuizzes = totalSectionIDs

	log.Println("[SERVICE] Updating existing UserUnitModel")
	return db.Save(&userUnit).Error
}

// func UpdateUserThemesOrLevelsIfUnitCompleted(db *gorm.DB, userID uuid.UUID, unitID int, themesOrLevelID int) error {
// 	log.Println("[SERVICE] UpdateUserThemesOrLevelsIfUnitCompleted - userID:", userID, "unitID:", unitID, "themesOrLevelID:", themesOrLevelID)

// 	// Ambil semua unit dari themes/level
// 	var allUnits []userUnitModel.UnitModel
// 	if err := db.Where("themes_or_level_id = ? AND deleted_at IS NULL", themesOrLevelID).Find(&allUnits).Error; err != nil {
// 		log.Println("[ERROR] Failed to fetch units for theme:", err)
// 		return err
// 	}
// 	totalUnitIDs := pq.Int64Array{}
// 	for _, unit := range allUnits {
// 		totalUnitIDs = append(totalUnitIDs, int64(unit.ID))
// 	}

// 	// Ambil user_unit untuk unit ini
// 	var userUnit userUnitModel.UserUnitModel
// 	if err := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit).Error; err != nil {
// 		log.Println("[SERVICE] UserUnit belum ada, skip update theme progress")
// 		return nil
// 	}

// 	unitCompleted := len(userUnit.CompleteSectionQuizzes) == len(userUnit.TotalSectionQuizzes)

// 	// Cek apakah sudah ada
// 	var userTheme userThemesOrLevelsModel.UserThemesOrLevelsModel
// 	err := db.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesOrLevelID).First(&userTheme).Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		userTheme = userThemesOrLevelsModel.UserThemesOrLevelsModel{
// 			UserID:           userID,
// 			ThemesOrLevelsID: themesOrLevelID,
// 			CompleteUnit:     pq.Int64Array{},
// 			TotalUnit:        totalUnitIDs,
// 		}

// 		if unitCompleted {
// 			userTheme.CompleteUnit = pq.Int64Array{int64(unitID)}
// 		}

// 		log.Println("[SERVICE] Creating new UserThemesOrLevelsModel")
// 		return db.Create(&userTheme).Error
// 	}

// 	// Update total unit selalu
// 	userTheme.TotalUnit = totalUnitIDs

// 	// Tambahkan unit jika lengkap dan belum tercatat
// 	if unitCompleted {
// 		for _, id := range userTheme.CompleteUnit {
// 			if id == int64(unitID) {
// 				log.Println("[SERVICE] Unit already recorded in CompleteUnit")
// 				return db.Save(&userTheme).Error
// 			}
// 		}
// 		userTheme.CompleteUnit = append(userTheme.CompleteUnit, int64(unitID))
// 	}

// 	log.Println("[SERVICE] Updating existing UserThemesOrLevelsModel")
// 	return db.Save(&userTheme).Error
// }

// ✅ Final: UpdateUserThemesOrLevelsIfUnitCompleted
func UpdateUserThemesOrLevelsIfUnitCompleted(db *gorm.DB, userID uuid.UUID, unitID int, themesOrLevelID int) error {
	log.Println("[SERVICE] UpdateUserThemesOrLevelsIfUnitCompleted - userID:", userID, "unitID:", unitID, "themesOrLevelID:", themesOrLevelID)

	// Ambil semua unit dari themes/level
	var allUnits []userUnitModel.UnitModel
	if err := db.Where("themes_or_level_id = ? AND deleted_at IS NULL", themesOrLevelID).Find(&allUnits).Error; err != nil {
		log.Println("[ERROR] Failed to fetch units for theme:", err)
		return err
	}
	totalUnitIDs := pq.Int64Array{}
	for _, unit := range allUnits {
		totalUnitIDs = append(totalUnitIDs, int64(unit.ID))
	}

	// Ambil user_unit untuk unit ini
	var userUnit userUnitModel.UserUnitModel
	if err := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit).Error; err != nil {
		log.Println("[SERVICE] UserUnit belum ada, skip update theme progress")
		return nil
	}

	// Cek apakah sudah ada
	var userTheme userThemesOrLevelsModel.UserThemesOrLevelsModel
	err := db.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesOrLevelID).First(&userTheme).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userTheme = userThemesOrLevelsModel.UserThemesOrLevelsModel{
			UserID:           userID,
			ThemesOrLevelsID: themesOrLevelID,
			CompleteUnit:     pq.Int64Array{int64(unitID)},
			TotalUnit:        totalUnitIDs,
		}
		log.Println("[SERVICE] Creating new UserThemesOrLevelsModel")
		return db.Create(&userTheme).Error
	}

	// Update total unit selalu
	userTheme.TotalUnit = totalUnitIDs

	// Tambahkan unitID ke CompleteUnit jika belum ada
	alreadyAdded := false
	for _, id := range userTheme.CompleteUnit {
		if id == int64(unitID) {
			alreadyAdded = true
			break
		}
	}
	if !alreadyAdded {
		userTheme.CompleteUnit = append(userTheme.CompleteUnit, int64(unitID))
	}

	log.Println("[SERVICE] Updating existing UserThemesOrLevelsModel")
	return db.Save(&userTheme).Error
}

func UpdateUserSubcategoryIfThemeCompleted(db *gorm.DB, userID uuid.UUID, themesOrLevelsID int, subcategoryID int) error {
	log.Println("[SERVICE] UpdateUserSubcategoryIfThemeCompleted - userID:", userID, "themesOrLevelsID:", themesOrLevelsID, "subcategoryID:", subcategoryID)

	// Ambil semua theme yang ada di subcategory
	var allThemes []userThemesOrLevelsModel.ThemesOrLevelsModel
	if err := db.Where("subcategories_id = ? AND deleted_at IS NULL", subcategoryID).Find(&allThemes).Error; err != nil {
		log.Println("[ERROR] Failed to fetch themes for subcategory:", err)
		return err
	}
	totalThemeIDs := pq.Int64Array{}
	for _, theme := range allThemes {
		totalThemeIDs = append(totalThemeIDs, int64(theme.ID))
	}

	// Ambil user_theme untuk theme ini
	var userTheme userThemesOrLevelsModel.UserThemesOrLevelsModel
	if err := db.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesOrLevelsID).First(&userTheme).Error; err != nil {
		log.Println("[SERVICE] UserThemesOrLevelsModel not found, skip update subcategory")
		return nil
	}

	// Jika belum selesai, jangan update complete
	if len(userTheme.CompleteUnit) < len(userTheme.TotalUnit) {
		log.Println("[SERVICE] Theme belum selesai, skip marking as complete")
	}

	// Cek apakah user_subcategory sudah ada
	var userSub userSubcategoryModel.UserSubcategoryModel
	err := db.Where("user_id = ? AND subcategory_id = ?", userID, subcategoryID).First(&userSub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userSub = userSubcategoryModel.UserSubcategoryModel{
			UserID:                 userID,
			SubcategoryID:          subcategoryID,
			CompleteThemesOrLevels: pq.Int64Array{},
			TotalThemesOrLevels:    totalThemeIDs,
		}

		if len(userTheme.CompleteUnit) == len(userTheme.TotalUnit) {
			userSub.CompleteThemesOrLevels = pq.Int64Array{int64(themesOrLevelsID)}
		}

		log.Println("[SERVICE] Creating new UserSubcategoryModel")
		return db.Create(&userSub).Error
	}

	// Selalu update total
	userSub.TotalThemesOrLevels = totalThemeIDs

	// Tambahkan themes jika selesai dan belum ada
	if len(userTheme.CompleteUnit) == len(userTheme.TotalUnit) {
		for _, id := range userSub.CompleteThemesOrLevels {
			if id == int64(themesOrLevelsID) {
				log.Println("[SERVICE] Theme already recorded in UserSubcategoryModel")
				return db.Save(&userSub).Error
			}
		}
		userSub.CompleteThemesOrLevels = append(userSub.CompleteThemesOrLevels, int64(themesOrLevelsID))
	}

	log.Println("[SERVICE] Updating existing UserSubcategoryModel")
	return db.Save(&userSub).Error
}

func UpdateUserCategoryIfSubcategoryCompleted(db *gorm.DB, userID uuid.UUID, subcategoryID int, categoryID int) error {
	log.Println("[SERVICE] UpdateUserCategoryIfSubcategoryCompleted - userID:", userID, "subcategoryID:", subcategoryID, "categoryID:", categoryID)

	// Ambil semua subcategory dalam category
	var allSubcategories []userSubcategoryModel.SubcategoryModel
	if err := db.Where("categories_id = ? AND deleted_at IS NULL", categoryID).Find(&allSubcategories).Error; err != nil {
		log.Println("[ERROR] Failed to fetch subcategories for category:", err)
		return err
	}
	totalSubIDs := pq.Int64Array{}
	for _, sub := range allSubcategories {
		totalSubIDs = append(totalSubIDs, int64(sub.ID))
	}

	// Ambil user_subcategory untuk cek apakah sudah selesai
	var userSub userSubcategoryModel.UserSubcategoryModel
	if err := db.Where("user_id = ? AND subcategory_id = ?", userID, subcategoryID).First(&userSub).Error; err != nil {
		log.Println("[SERVICE] UserSubcategoryModel not found, skip update category")
		return nil
	}

	if len(userSub.CompleteThemesOrLevels) < len(userSub.TotalThemesOrLevels) {
		log.Println("[SERVICE] Subcategory belum selesai, skip marking category")
	}

	// Cek apakah user_category sudah ada
	var userCat userCategoryModel.UserCategoryModel
	err := db.Where("user_id = ? AND category_id = ?", userID, categoryID).First(&userCat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userCat = userCategoryModel.UserCategoryModel{
			UserID:           userID,
			CategoryID:       categoryID,
			CompleteCategory: pq.Int64Array{},
			TotalCategory:    totalSubIDs,
		}
		if len(userSub.CompleteThemesOrLevels) == len(userSub.TotalThemesOrLevels) {
			userCat.CompleteCategory = pq.Int64Array{int64(subcategoryID)}
		}
		log.Println("[SERVICE] Creating new UserCategoryModel")
		return db.Create(&userCat).Error
	}

	userCat.TotalCategory = totalSubIDs

	if len(userSub.CompleteThemesOrLevels) == len(userSub.TotalThemesOrLevels) {
		for _, id := range userCat.CompleteCategory {
			if id == int64(subcategoryID) {
				log.Println("[SERVICE] Subcategory already recorded in UserCategoryModel")
				return db.Save(&userCat).Error
			}
		}
		userCat.CompleteCategory = append(userCat.CompleteCategory, int64(subcategoryID))
	}

	log.Println("[SERVICE] Updating existing UserCategoryModel")
	return db.Save(&userCat).Error
}
