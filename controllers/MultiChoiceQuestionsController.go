package controllers

import (
	"AAT_Api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// GetMultiChoiceQuestions
// @Summary Get all multi choice questions
// @Description Get all multi choice questions
// @Tags MultiChoiceQuestions
// @Produce json
// @Param tags query string false "Tags to filter questions"
// @Success 200 {array} models.MultiChoiceQuestionDTO
// @Router /multi-choice [get]
func GetMultiChoiceQuestions(c *gin.Context) {
	var questions []models.MultiChoiceQuestion

	// resolve tags query from context
	tags := c.QueryArray("tags")

	// Get all questions with multiple optional tags
	if len(tags) > 0 {
		models.DB.Where("tags.name IN ?", tags).Joins("JOIN question_tags ON question_tags.question_id = multi_choice_questions.id").Joins("JOIN tags ON tags.id = question_tags.tag_id").Find(&questions)
	} else {
		models.DB.Find(&questions)
	}
	// Convert questions to DTO
	var dtos []models.MultiChoiceQuestionDTO
	for _, question := range questions {
		var dto models.MultiChoiceQuestionDTO
		dto.ID = question.ID
		dto.Title = question.Title
		dto.Description = question.Description
		dto.Credit = question.Credit
		dto.Feedback = question.Feedback
		dto.Difficulty = question.Difficulty
		// Get answers
		var answers []models.MultiChoiceAnswer
		models.DB.Where("question_id = ?", question.ID).Find(&answers)
		for _, answer := range answers {
			var answerDto models.MultiChoiceAnswerDTO
			answerDto.Content = answer.Content
			answerDto.IsCorrect = answer.IsCorrect
			dto.Answers = append(dto.Answers, answerDto)
		}
		// Get tags
		var tags []models.Tag
		err := models.DB.Model(&question).Association("Tags").Find(&tags)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
			return
		}
		for _, tag := range tags {
			dto.Tags = append(dto.Tags, tag.Name)
		}
		dtos = append(dtos, dto)
	}
	// Return DTOs as json
	c.IndentedJSON(http.StatusOK, dtos)
}

// GetMultiChoiceQuestion
// @Summary Get a multi choice question
// @Description Get a multi choice question
// @Tags MultiChoiceQuestions
// @Produce json
// @Param id path string true "ID of the question"
// @Success 200 {object} models.MultiChoiceQuestionDTO
// @Router /multi-choice/{id} [get]
func GetMultiChoiceQuestion(c *gin.Context) {
	var question models.MultiChoiceQuestion
	models.DB.First(&question, c.Param("id"))

	var answers []models.MultiChoiceAnswer
	err := models.DB.Where("question_id = ?", question.ID).Find(&answers)
	if err != nil && err.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch answers"})
		return
	}

	var tags []models.Tag
	err1 := models.DB.Model(&question).Association("Tags").Find(&tags)
	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
		return
	}

	var dto models.MultiChoiceQuestionDTO
	dto.ID = question.ID
	dto.Title = question.Title
	dto.Description = question.Description
	dto.Credit = question.Credit
	dto.Feedback = question.Feedback
	dto.Difficulty = question.Difficulty

	for _, answer := range answers {
		var answerDTO models.MultiChoiceAnswerDTO
		answerDTO.Content = answer.Content
		answerDTO.IsCorrect = answer.IsCorrect

		dto.Answers = append(dto.Answers, answerDTO)
	}

	for _, tag := range tags {
		dto.Tags = append(dto.Tags, tag.Name)
	}

	c.IndentedJSON(http.StatusOK, dto)
}

// CreateMultiChoiceQuestion
// @Summary Create a multi choice question
// @Description Create a multi choice question
// @Tags MultiChoiceQuestions
// @Accept json
// @Produce json
// @Param question body models.MultiChoiceQuestionDTO true "Question to create"
// @Success 200 {object} models.MultiChoiceQuestionDTO
// @Failure 400 {object} string "Invalid request body"
// @Router /multi-choice [post]
func CreateMultiChoiceQuestion(c *gin.Context) {
	var dto models.MultiChoiceQuestionDTO
	err := c.BindJSON(&dto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var question models.MultiChoiceQuestion
	question.Title = dto.Title
	question.Description = dto.Description
	question.Credit = dto.Credit
	question.Feedback = dto.Feedback
	question.Difficulty = dto.Difficulty

	models.DB.Create(&question)

	for _, answerDTO := range dto.Answers {
		var answer models.MultiChoiceAnswer
		answer.Content = answerDTO.Content
		answer.IsCorrect = answerDTO.IsCorrect
		answer.QuestionID = question.ID
		models.DB.Create(&answer)
		// append answer to question
		err := models.DB.Model(&question).Association("Answers").Append(&answer)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when append answer to question"})
			return
		}
	}

	for _, tagName := range dto.Tags {
		// check if tag exists, if not create it
		var tag models.Tag
		models.DB.Where("name = ?", tagName).First(&tag)
		if tag.ID == 0 {
			tag.Name = tagName
			models.DB.Create(&tag)
		}
		// append tag to question
		err := models.DB.Model(&question).Association("Tags").Append(&tag)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when append tag to question"})
			return
		}
	}

	// convert question with its tags and answers to dto
	// the tags and answers are already loaded into question object,
	// so we can just convert it to dto
	var responseDTO models.MultiChoiceQuestionDTO
	responseDTO.ID = question.ID
	responseDTO.Title = question.Title
	responseDTO.Description = question.Description
	responseDTO.Credit = question.Credit
	responseDTO.Feedback = question.Feedback
	responseDTO.Difficulty = question.Difficulty
	// convert question.tags to DTOs and append to responseDTO
	for _, tag := range question.Tags {
		responseDTO.Tags = append(responseDTO.Tags, tag.Name)
	}
	// convert question.answers to DTOs and append to responseDTO
	for _, answer := range question.Answers {
		var answerDTO models.MultiChoiceAnswerDTO
		answerDTO.Content = answer.Content
		answerDTO.IsCorrect = answer.IsCorrect
		responseDTO.Answers = append(responseDTO.Answers, answerDTO)
	}
	c.IndentedJSON(http.StatusCreated, responseDTO)
}

// DeleteMultiChoiceQuestion
// @Summary Delete a multi choice question
// @Description Delete a multi choice question
// @Tags MultiChoiceQuestions
// @Produce json
// @Param id path string true "ID of the question"
// @Success 200 {object} string "Question deleted"
// @Router /multi-choice/{id} [delete]
func DeleteMultiChoiceQuestion(c *gin.Context) {
	var question models.MultiChoiceQuestion
	models.DB.First(&question, c.Param("id"))
	models.DB.Delete(&question)
	// Also delete all answers
	models.DB.Where("question_id = ?", question.ID).Delete(&models.MultiChoiceAnswer{})
	c.IndentedJSON(http.StatusNoContent, gin.H{"message": "Question deleted"})
}

// UpdateMultiChoiceQuestion
// @Summary Update a multi choice question
// @Description Update a multi choice question
// @Tags MultiChoiceQuestions
// @Accept json
// @Produce json
// @Param id path string true "ID of the question"
// @Param question body models.MultiChoiceQuestionDTO true "Question to update"
// @Success 200 {object} models.MultiChoiceQuestionDTO
// @Failure 400 {object} string "Invalid request body"
// @Router /multi-choice/{id} [put]
func UpdateMultiChoiceQuestion(c *gin.Context) {
	// Convert request body to DTO
	var dto models.MultiChoiceQuestionDTO
	err := c.BindJSON(&dto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Find question by ID
	var question models.MultiChoiceQuestion
	models.DB.First(&question, c.Param("id"))
	if question.ID == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	// Update question
	question.Title = dto.Title
	question.Description = dto.Description
	question.Credit = dto.Credit
	question.Feedback = dto.Feedback
	question.Difficulty = dto.Difficulty
	models.DB.Save(&question)
	// Delete all answers
	models.DB.Where("question_id = ?", question.ID).Delete(&models.MultiChoiceAnswer{})
	// Create new answers
	for _, answerDTO := range dto.Answers {
		var answer models.MultiChoiceAnswer
		answer.Content = answerDTO.Content
		answer.IsCorrect = answerDTO.IsCorrect
		answer.QuestionID = question.ID
		models.DB.Create(&answer)
		// append answer to question
		err := models.DB.Model(&question).Association("Answers").Append(&answer)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when append answer to question"})
			return
		}
	}
	// Delete all tags
	err = models.DB.Model(&question).Association("Tags").Clear()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when delete tags from question"})
		return
	}
	// Create new tags
	for _, tagName := range dto.Tags {
		// check if tag exists, if not create it
		var tag models.Tag
		models.DB.Where("name = ?", tagName).First(&tag)
		if tag.ID == 0 {
			tag.Name = tagName
			models.DB.Create(&tag)
		}
		// append tag to question
		err := models.DB.Model(&question).Association("Tags").Append(&tag)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when append tag to question"})
			return
		}
	}
	// convert question with its tags and answers to dto
	// the tags and answers are already loaded into question object,
	// so we can just convert it to dto
	var responseDTO models.MultiChoiceQuestionDTO
	responseDTO.ID = question.ID
	responseDTO.Title = question.Title
	responseDTO.Description = question.Description
	responseDTO.Credit = question.Credit
	responseDTO.Feedback = question.Feedback
	responseDTO.Difficulty = question.Difficulty
	// convert question.tags to DTOs and append to responseDTO
	for _, tag := range question.Tags {
		responseDTO.Tags = append(responseDTO.Tags, tag.Name)
	}
	// convert question.answers to DTOs and append to responseDTO
	for _, answer := range question.Answers {
		var answerDTO models.MultiChoiceAnswerDTO
		answerDTO.Content = answer.Content
		answerDTO.IsCorrect = answer.IsCorrect
		responseDTO.Answers = append(responseDTO.Answers, answerDTO)
	}
	c.IndentedJSON(http.StatusOK, responseDTO)
}
