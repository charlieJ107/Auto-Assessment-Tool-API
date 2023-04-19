package models

import "gorm.io/gorm"

type FormativeAssessment struct {
	gorm.Model
	Title                string
	Difficulty           string
	Points               int
	Topic                string
	MultiChoiceQuestions []MultiChoiceQuestion `gorm:"many2many:formative_assessment_multi_choice_questions;"`
	ShortAnswerQuestions []ShortAnswerQuestion `gorm:"many2many:formative_assessment_short_answer_questions;"`
	Randomisation        bool
	Tags                 []Tag `gorm:"many2many:formative_assessment_tags;"`
}

type FormativeAssessmentResponseDTO struct {
	ID                   uint                     `json:"id"`
	Title                string                   `json:"title"`
	Difficulty           string                   `json:"difficulty"`
	Points               int                      `json:"points"`
	Topic                string                   `json:"topic"`
	MultiChoiceQuestions []MultiChoiceQuestionDTO `json:"multi_choice_questions"`
	ShortAnswerQuestions []ShortAnswerQuestionDTO `json:"short_answer_questions"`
	Randomisation        bool                     `json:"randomisation"`
	Tags                 []string                 `json:"tags"`
}

type FormativeAssessmentRequestDTO struct {
	ID                     uint     `json:"id"`
	Title                  string   `json:"title"`
	Difficulty             string   `json:"difficulty"`
	Points                 int      `json:"points"`
	Topic                  string   `json:"topic"`
	MultiChoiceQuestionIds []uint   `json:"multi_choice_question_ids"`
	ShortAnswerQuestionIds []uint   `json:"short_answer_question_ids"`
	Randomisation          bool     `json:"randomisation"`
	Tags                   []string `json:"tags"`
}
