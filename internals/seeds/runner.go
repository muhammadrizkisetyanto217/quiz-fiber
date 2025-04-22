package seeds

import (
	categories "quiz-fiber/internals/seeds/lessons/categories"
	difficulties "quiz-fiber/internals/seeds/lessons/difficulties"
	subcategories "quiz-fiber/internals/seeds/lessons/subcategories"
	themes_or_levels "quiz-fiber/internals/seeds/lessons/themes_or_levels"
	units "quiz-fiber/internals/seeds/lessons/units"

	evaluations "quiz-fiber/internals/seeds/quizzes/evaluations"
	exams "quiz-fiber/internals/seeds/quizzes/exams"
	questions "quiz-fiber/internals/seeds/quizzes/questions"
	quizzes "quiz-fiber/internals/seeds/quizzes/quizzes"
	reading "quiz-fiber/internals/seeds/quizzes/readings"
	section_quizzes "quiz-fiber/internals/seeds/quizzes/section_quizzes"

	level "quiz-fiber/internals/seeds/progress/levels"
	rank "quiz-fiber/internals/seeds/progress/ranks"

	"gorm.io/gorm"
)

func RunAllSeeds(db *gorm.DB) {

	//* Category
	difficulties.SeedDifficultiesFromJSON(db, "internals/seeds/category/difficulty/data_difficulty.json")
	categories.SeedCategoriesFromJSON(db, "internals/seeds/category/category/data_category.json")
	subcategories.SeedSubcategoriesFromJSON(db, "internals/seeds/category/subcategory/data_subcategory.json")
	themes_or_levels.SeedThemesOrLevelsFromJSON(db, "internals/seeds/category/themes_or_levels/data_themes_or_levels.json")
	units.SeedUnitsFromJSON(db, "internals/seeds/category/units/data_units.json")

	//* User

	//* Quizzes
	evaluations.SeedEvaluationsFromJSON(db, "internals/seeds/quizzes/evaluations/data_evaluations.json")
	exams.SeedExamsFromJSON(db, "internals/seeds/quizzes/exams/data_exams.json")
	questions.SeedQuestionsFromJSON(db, "internals/seeds/quizzes/questions/data_questions.json")
	quizzes.SeedQuizzesFromJSON(db, "internals/seeds/quizzes/quizzes/data_quizzes.json")
	reading.SeedReadingsFromJSON(db, "internals/seeds/quizzes/readings/data_readings.json")
	section_quizzes.SeedSectionQuizzesFromJSON(db, "internals/seeds/quizzes/section_quizzes/data_section_quizzes.json")

	//* Progress
	level.SeedLevelRequirementsFromJSON(db, "internals/seeds/progress/levels/data_levels_requirements.json")
	rank.SeedRanksRequirementsFromJSON(db, "internals/seeds/progress/ranks/data_ranks_requirements.json")

}
