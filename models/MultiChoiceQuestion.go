package models

import "gorm.io/gorm"

type MultiChoiceQuestion struct {
	gorm.Model
	Title       string
	Description string
	Credit      int
	Feedback    string
	Difficulty  string
	Answers     []MultiChoiceAnswer `gorm:"foreignKey:QuestionID"`
	Tags        []Tag               `gorm:"many2many:multi_choice_question_tags;"`
}

type MultiChoiceAnswer struct {
	gorm.Model
	Content    string
	QuestionID uint
	IsCorrect  bool
}

type MultiChoiceQuestionDTO struct {
	ID          uint                   `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Credit      int                    `json:"credit"`
	Feedback    string                 `json:"feedback"`
	Difficulty  string                 `json:"difficulty"`
	Answers     []MultiChoiceAnswerDTO `json:"answers"`
	Tags        []string               `json:"tags"`
}

type MultiChoiceAnswerDTO struct {
	Content   string `json:"content"`
	IsCorrect bool   `json:"is_correct"`
}
