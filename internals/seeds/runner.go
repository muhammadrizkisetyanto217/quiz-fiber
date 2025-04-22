package seeds

import (
	category "quiz-fiber/internals/seeds/lessons/categories"
	difficulty "quiz-fiber/internals/seeds/lessons/difficulties"
	subcategory "quiz-fiber/internals/seeds/lessons/subcategories"
	themes_or_levels "quiz-fiber/internals/seeds/lessons/themes_or_levels"
	units "quiz-fiber/internals/seeds/lessons/units"

	"gorm.io/gorm"
)

func RunAllSeeds(db *gorm.DB) {

	//* Category
	difficulty.SeedDifficultiesFromJSON(db, "internals/seeds/category/difficulty/data_difficulty.json")
	category.SeedCategoriesFromJSON(db, "internals/seeds/category/category/data_category.json")
	subcategory.SeedSubcategoriesFromJSON(db, "internals/seeds/category/subcategory/data_subcategory.json")
	themes_or_levels.SeedThemesOrLevelsFromJSON(db, "internals/seeds/category/themes_or_levels/data_themes_or_levels.json")
	units.SeedUnitsFromJSON(db, "internals/seeds/category/units/data_units.json")


	//* User


	//* Quizzes


}
