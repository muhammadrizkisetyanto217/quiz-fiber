DROP INDEX IF EXISTS idx_question_saved_question;
DROP INDEX IF EXISTS idx_question_saved_user;
DROP TABLE IF EXISTS question_saved;


DROP INDEX IF EXISTS idx_question_mistakes_question_id;
DROP INDEX IF EXISTS idx_question_mistakes_user_id;
DROP TABLE IF EXISTS question_mistakes;


DROP INDEX IF EXISTS idx_user_questions_user_id;
DROP INDEX IF EXISTS idx_user_questions_question_id;
DROP INDEX IF EXISTS idx_user_questions_source_type_id;
DROP INDEX IF EXISTS idx_user_questions_source_id;

DROP TABLE IF EXISTS user_questions;
