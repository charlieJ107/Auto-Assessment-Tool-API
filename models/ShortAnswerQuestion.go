package models

import "gorm.io/gorm"

type ShortAnswerQuestion struct {
	gorm.Model
	Title         string
	Description   string
	Credit        int
	CorrectAnswer string
	Tags          []Tag `gorm:"many2many:short_answer_question_tags;"`
}

type ShortAnswerQuestionDTO struct {
	ID            uint     `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Credit        int      `json:"credit"`
	CorrectAnswer string   `json:"correct_answer"`
	Tags          []string `json:"tags"`
}
