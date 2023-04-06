package controllers

import (
	"AAT_Api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

//@BasePath /api/v1

// GetShortAnswerQuestions
// @Summary Get all short answer questions
// @Description Get all short answer questions
// @Tags ShortAnswerQuestions
// @Produce json
// @Param tags query string false "Tags to filter questions"
// @Success 200 {array} models.ShortAnswerQuestionDTO
// @Router /short-answer [get]
func GetShortAnswerQuestions(c *gin.Context) {
	var questions []models.ShortAnswerQuestion
	// resolve tags query from context
	tags := c.QueryArray("tags")
	// Get all questions with multiple optional tags
	if len(tags) > 0 {
		models.DB.Where("tags.name IN ?", tags).Joins("JOIN question_tags ON question_tags.question_id = short_answer_questions.id").Joins("JOIN tags ON tags.id = question_tags.tag_id").Find(&questions)
	} else {
		models.DB.Find(&questions)
	}
	// Convert questions to DTO
	var dtos []models.ShortAnswerQuestionDTO
	for _, question := range questions {
		var dto models.ShortAnswerQuestionDTO
		dto.ID = question.ID
		dto.Title = question.Title
		dto.Description = question.Description
		dto.Credit = question.Credit
		dto.CorrectAnswer = question.CorrectAnswer
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

// GetShortAnswerQuestion
// @Summary Get a short answer question
// @Description Get a short answer question
// @Tags ShortAnswerQuestions
// @Produce json
// @Param id path int true "Question ID"
// @Success 200 {object} models.ShortAnswerQuestionDTO
// @Router /short-answer/{id} [get]
func GetShortAnswerQuestion(c *gin.Context) {
	var question models.ShortAnswerQuestion
	id := c.Param("id")
	err := models.DB.Where("id = ?", id).First(&question).Error
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	// Convert question to DTO
	var dto models.ShortAnswerQuestionDTO
	dto.ID = question.ID
	dto.Title = question.Title
	dto.Description = question.Description
	dto.Credit = question.Credit
	dto.CorrectAnswer = question.CorrectAnswer
	// Get tags
	var tags []models.Tag
	err = models.DB.Model(&question).Association("Tags").Find(&tags)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
		return
	}
	for _, tag := range tags {
		dto.Tags = append(dto.Tags, tag.Name)
	}
	// Return DTO as json
	c.IndentedJSON(http.StatusOK, dto)
}

// CreateShortAnswerQuestion
// @Summary Create a short answer question
// @Description Create a short answer question
// @Tags ShortAnswerQuestions
// @Produce json
// @Param question body models.ShortAnswerQuestionDTO true "Question to create"
// @Success 200 {object} models.ShortAnswerQuestionDTO
// @Router /short-answer [post]
func CreateShortAnswerQuestion(c *gin.Context) {
	var dto models.ShortAnswerQuestionDTO
	err := c.BindJSON(&dto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when binding json"})
		return
	}
	// Create question
	var question models.ShortAnswerQuestion
	question.Title = dto.Title
	question.Description = dto.Description
	question.Credit = dto.Credit
	question.CorrectAnswer = dto.CorrectAnswer
	err = models.DB.Create(&question).Error
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when creating question"})
		return
	}
	// Create tags
	for _, tagName := range dto.Tags {
		var tag models.Tag
		err = models.DB.Where("name = ?", tagName).First(&tag).Error
		if err != nil {
			tag.Name = tagName
			err = models.DB.Create(&tag).Error
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when creating tag"})
				return
			}
		}
		err = models.DB.Model(&question).Association("Tags").Append(&tag)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when appending tag to question"})
			return
		}
	}
	// Convert question to DTO
	dto.ID = question.ID
	// Return DTO as json
	c.IndentedJSON(http.StatusCreated, dto)
}

// UpdateShortAnswerQuestion
// @Summary Update a short answer question
// @Description Update a short answer question
// @Tags ShortAnswerQuestions
// @Produce json
// @Param id path int true "Question ID"
// @Param question body models.ShortAnswerQuestionDTO true "Question to update"
// @Success 200 {object} models.ShortAnswerQuestionDTO
// @Router /short-answer/{id} [put]
func UpdateShortAnswerQuestion(c *gin.Context) {
	var dto models.ShortAnswerQuestionDTO
	err := c.BindJSON(&dto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when binding json"})
		return
	}
	// Find question by ID
	var question models.ShortAnswerQuestion
	models.DB.First(&question, c.Param("id"))
	if question.ID == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	// Update question
	question.Title = dto.Title
	question.Description = dto.Description
	question.Credit = dto.Credit
	question.CorrectAnswer = dto.CorrectAnswer
	err = models.DB.Save(&question).Error
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when updating question"})
		return
	}
	// Delete old tags
	err = models.DB.Model(&question).Association("Tags").Clear()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when clearing tags"})
		return
	}
	// Create new tags
	for _, tagName := range dto.Tags {
		var tag models.Tag
		err = models.DB.Where("name = ?", tagName).First(&tag).Error
		if err != nil {
			tag.Name = tagName
			err = models.DB.Create(&tag).Error
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when creating tag"})
				return
			}
		}
		err = models.DB.Model(&question).Association("Tags").Append(&tag)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when appending tag to question"})
			return
		}
	}
	// convert question to DTO
	dto.ID = question.ID
	// Return DTO as json
	c.IndentedJSON(http.StatusOK, dto)
}

// DeleteShortAnswerQuestion
// @Summary Delete a short answer question
// @Description Delete a short answer question
// @Tags ShortAnswerQuestions
// @Produce json
// @Param id path int true "Question ID"
// @Success 200 {object} models.ShortAnswerQuestionDTO
// @Router /short-answer/{id} [delete]
func DeleteShortAnswerQuestion(c *gin.Context) {
	// Get question
	var question models.ShortAnswerQuestion
	id := c.Param("id")
	err := models.DB.Where("id = ?", id).First(&question).Error
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	// Delete question
	err = models.DB.Delete(&question).Error
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when deleting question"})
		return
	}
	// Return DTO as json
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}
