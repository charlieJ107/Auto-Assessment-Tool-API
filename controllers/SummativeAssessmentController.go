package controllers

import (
	"AAT_Api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Base path: /api/v1

// GetSummativeAssessments
// @Summary Get all summative assessments
// @Description Get all summative assessments
// @Tags SummativeAssessment
// @Accept  json
// @Produce  json
// @Param tags query string false "tags"
// @Success 200 {object} []string
// @Router /summativeAssessments [get]
func GetSummativeAssessments(c *gin.Context) {
	var assessments []models.SummativeAssessment
	tags := c.Query("tags")
	if len(tags) > 0 {
		models.DB.
			Model(&models.SummativeAssessment{}).
			Preload("MultipleChoiceQuestions").Preload("ShortAnswerQuestions").
			Where("tags.name IN ?", tags).
			Joins("JOIN assessment_tags ON assessment_tags.assessment_id = summative_assessments.id").
			Joins("JOIN tags ON tags.id = assessment_tags.tag_id").
			Find(&assessments)
	} else {
		models.DB.
			Model(&models.SummativeAssessment{}).
			Preload("MultipleChoiceQuestions").Preload("ShortAnswerQuestions").
			Find(&assessments)
	}
	// Convert assessments to DTO
	var dtos []models.SummativeAssessmentResponseDTO
	for _, assessment := range assessments {
		var dto models.SummativeAssessmentResponseDTO
		dto.ID = assessment.ID
		dto.Title = assessment.Title
		dto.Difficulty = assessment.Difficulty
		dto.Points = assessment.Points
		dto.Topic = assessment.Topic
		dto.Randomisation = assessment.Randomisation
		// Get MultiChoiceQuestions
		for _, multiChoiceQuestion := range assessment.MultiChoiceQuestions {
			var multiChoiceQuestionDto models.MultiChoiceQuestionDTO
			multiChoiceQuestionDto.ID = multiChoiceQuestion.ID
			multiChoiceQuestionDto.Title = multiChoiceQuestion.Title
			multiChoiceQuestionDto.Description = multiChoiceQuestion.Description
			multiChoiceQuestionDto.Credit = multiChoiceQuestion.Credit
			// Get answers
			var answers []models.MultiChoiceAnswer
			models.DB.Where("question_id = ?", multiChoiceQuestion.ID).Find(&answers)
			for _, answer := range answers {
				var answerDto models.MultiChoiceAnswerDTO
				answerDto.Content = answer.Content
				answerDto.IsCorrect = answer.IsCorrect
				multiChoiceQuestionDto.Answers = append(multiChoiceQuestionDto.Answers, answerDto)
			}
			// Get tags
			var tags []models.Tag
			err := models.DB.Model(&multiChoiceQuestion).Association("Tags").Find(&tags)
			if err != nil {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
				return
			}
			for _, tag := range tags {
				multiChoiceQuestionDto.Tags = append(multiChoiceQuestionDto.Tags, tag.Name)
			}
			dto.MultiChoiceQuestions = append(dto.MultiChoiceQuestions, multiChoiceQuestionDto)
		}

		// Get ShortAnswerQuestions
		for _, shortAnswerQuestion := range assessment.ShortAnswerQuestions {
			var shortAnswerQuestionDto models.ShortAnswerQuestionDTO
			shortAnswerQuestionDto.ID = shortAnswerQuestion.ID
			shortAnswerQuestionDto.Title = shortAnswerQuestion.Title
			shortAnswerQuestionDto.Description = shortAnswerQuestion.Description
			shortAnswerQuestionDto.Credit = shortAnswerQuestion.Credit
			shortAnswerQuestionDto.CorrectAnswer = shortAnswerQuestion.CorrectAnswer
			// Get tags
			var tags []models.Tag
			err := models.DB.Model(&shortAnswerQuestion).Association("Tags").Find(&tags)
			if err != nil {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
				return
			}
			for _, tag := range tags {
				shortAnswerQuestionDto.Tags = append(shortAnswerQuestionDto.Tags, tag.Name)
			}
		}
		// Get tags
		var tags []models.Tag
		err := models.DB.Model(&assessment).Association("Tags").Find(&tags)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
			return
		}
		for _, tag := range tags {
			dto.Tags = append(dto.Tags, tag.Name)
		}
		dtos = append(dtos, dto)
	}
	c.IndentedJSON(http.StatusOK, dtos)
}

// GetSummativeAssessment
// @Summary Get summative assessment by ID
// @Description Get summative assessment by ID
// @Tags SummativeAssessment
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} string
// @Router /summativeAssessments/{id} [get]
func GetSummativeAssessment(c *gin.Context) {
	var assessment models.SummativeAssessment
	id := c.Param("id")
	err := models.DB.Model(&models.SummativeAssessment{}).
		Preload("MultipleChoiceQuestions").Preload("ShortAnswerQuestions").
		Where("id = ?", id).
		First(&assessment).Error
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Summative assessment not found"})
		return
	}
	// Convert assessment to DTO
	var dto models.SummativeAssessmentResponseDTO
	dto.ID = assessment.ID
	dto.Title = assessment.Title
	dto.Difficulty = assessment.Difficulty
	dto.Points = assessment.Points
	dto.Topic = assessment.Topic
	dto.Randomisation = assessment.Randomisation
	// Get MultiChoiceQuestions
	for _, multiChoiceQuestion := range assessment.MultiChoiceQuestions {
		var multiChoiceQuestionDto models.MultiChoiceQuestionDTO
		multiChoiceQuestionDto.ID = multiChoiceQuestion.ID
		multiChoiceQuestionDto.Title = multiChoiceQuestion.Title
		multiChoiceQuestionDto.Description = multiChoiceQuestion.Description
		multiChoiceQuestionDto.Credit = multiChoiceQuestion.Credit
		// Get answers
		var answers []models.MultiChoiceAnswer
		models.DB.Where("question_id = ?", multiChoiceQuestion.ID).Find(&answers)
		for _, answer := range answers {
			var answerDto models.MultiChoiceAnswerDTO
			answerDto.Content = answer.Content
			answerDto.IsCorrect = answer.IsCorrect
			multiChoiceQuestionDto.Answers = append(multiChoiceQuestionDto.Answers, answerDto)
		}
		// Get tags
		var tags []models.Tag
		err := models.DB.Model(&multiChoiceQuestion).Association("Tags").Find(&tags)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
			return
		}
		for _, tag := range tags {
			multiChoiceQuestionDto.Tags = append(multiChoiceQuestionDto.Tags, tag.Name)
		}
		dto.MultiChoiceQuestions = append(dto.MultiChoiceQuestions, multiChoiceQuestionDto)
	}

	// Get ShortAnswerQuestions
	for _, shortAnswerQuestion := range assessment.ShortAnswerQuestions {
		var shortAnswerQuestionDto models.ShortAnswerQuestionDTO
		shortAnswerQuestionDto.ID = shortAnswerQuestion.ID
		shortAnswerQuestionDto.Title = shortAnswerQuestion.Title
		shortAnswerQuestionDto.Description = shortAnswerQuestion.Description
		shortAnswerQuestionDto.Credit = shortAnswerQuestion.Credit
		shortAnswerQuestionDto.CorrectAnswer = shortAnswerQuestion.CorrectAnswer
		// Get tags
		var tags []models.Tag
		err := models.DB.Model(&shortAnswerQuestion).Association("Tags").Find(&tags)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
			return
		}
		for _, tag := range tags {
			shortAnswerQuestionDto.Tags = append(shortAnswerQuestionDto.Tags, tag.Name)
		}
	}
	// Get tags
	var tags []models.Tag
	err = models.DB.Model(&assessment).Association("Tags").Find(&tags)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
		return
	}
	for _, tag := range tags {
		dto.Tags = append(dto.Tags, tag.Name)
	}
	c.IndentedJSON(http.StatusOK, dto)
}

// CreateSummativeAssessment
// @Summary Create summative assessment
// @Description Create summative assessment
// @Tags SummativeAssessment
// @Accept  json
// @Produce  json
// @Param assessment body CreateSummativeAssessmentRequestDTO true "Assessment"
// @Success 200 {object} string
// @Router /summativeAssessments [post]
func CreateSummativeAssessment(c *gin.Context) {
	var dto models.SummativeAssessmentRequestDTO
	err := c.BindJSON(&dto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when binding request body"})
		return
	}
	// Create assessment
	var assessment models.SummativeAssessment
	assessment.Title = dto.Title
	assessment.Difficulty = dto.Difficulty
	assessment.Points = dto.Points
	assessment.Topic = dto.Topic
	assessment.Randomisation = dto.Randomisation
	models.DB.Create(&assessment)
	// create tags
	for _, tagName := range dto.Tags {
		var tag models.Tag
		models.DB.Where("name = ?", tagName).First(&tag)
		if tag.ID == 0 {
			tag.Name = tagName
			models.DB.Create(&tag)
		}
		err := models.DB.Model(&assessment).Association("Tags").Append(&tag)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when append tags"})
			return
		}
	}
	// Set multi choice questions from database with multi choice question id
	for _, multiChoiceQuestionId := range dto.MultiChoiceQuestionIds {
		var multiChoiceQuestion models.MultiChoiceQuestion
		models.DB.Where("id = ?", multiChoiceQuestionId).First(&multiChoiceQuestion)
		err := models.DB.Model(&assessment).Association("MultiChoiceQuestions").Append(&multiChoiceQuestion)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when append multi choice questions"})
			return
		}
	}
	// Set short answer questions from database with short answer question id
	for _, shortAnswerQuestionId := range dto.ShortAnswerQuestionIds {
		var shortAnswerQuestion models.ShortAnswerQuestion
		models.DB.Where("id = ?", shortAnswerQuestionId).First(&shortAnswerQuestion)
		err := models.DB.Model(&assessment).Association("ShortAnswerQuestions").Append(&shortAnswerQuestion)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when append short answer questions"})
			return
		}
	}
	// Convert assessment to DTO
	var responseDto models.SummativeAssessmentResponseDTO
	responseDto.ID = assessment.ID
	responseDto.Title = assessment.Title
	responseDto.Difficulty = assessment.Difficulty
	responseDto.Points = assessment.Points
	responseDto.Topic = assessment.Topic
	responseDto.Randomisation = assessment.Randomisation

	// Get MultiChoiceQuestions with ids from assessmentDTO multiChoiceQuestionIds
	var multiChoiceQuestions []models.MultiChoiceQuestion
	for _, multiChoiceQuestionId := range dto.MultiChoiceQuestionIds {
		var multiChoiceQuestion models.MultiChoiceQuestion
		models.DB.Where("id = ?", multiChoiceQuestionId).First(&multiChoiceQuestion)
		multiChoiceQuestions = append(multiChoiceQuestions, multiChoiceQuestion)
	}
	for _, multiChoiceQuestion := range multiChoiceQuestions {
		var multiChoiceQuestionDto models.MultiChoiceQuestionDTO
		multiChoiceQuestionDto.ID = multiChoiceQuestion.ID
		multiChoiceQuestionDto.Title = multiChoiceQuestion.Title
		multiChoiceQuestionDto.Description = multiChoiceQuestion.Description
		multiChoiceQuestionDto.Credit = multiChoiceQuestion.Credit
		// Get answers
		var answers []models.MultiChoiceAnswer
		models.DB.Where("question_id = ?", multiChoiceQuestion.ID).Find(&answers)
		for _, answer := range answers {
			var answerDto models.MultiChoiceAnswerDTO
			answerDto.Content = answer.Content
			answerDto.IsCorrect = answer.IsCorrect
			multiChoiceQuestionDto.Answers = append(multiChoiceQuestionDto.Answers, answerDto)
		}
		// Get tags
		var tags []models.Tag
		err := models.DB.Model(&multiChoiceQuestion).Association("Tags").Find(&tags)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
			return
		}
		for _, tag := range tags {
			multiChoiceQuestionDto.Tags = append(multiChoiceQuestionDto.Tags, tag.Name)
		}
		responseDto.MultiChoiceQuestions = append(responseDto.MultiChoiceQuestions, multiChoiceQuestionDto)
	}
	// Get ShortAnswerQuestions with ids from assessmentDTO shortAnswerQuestionIds
	var shortAnswerQuestions []models.ShortAnswerQuestion
	err = models.DB.Model(&assessment).Association("ShortAnswerQuestions").Find(&shortAnswerQuestions)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch short answer questions"})
		return
	}
	for _, shortAnswerQuestion := range shortAnswerQuestions {
		var shortAnswerQuestionDto models.ShortAnswerQuestionDTO
		shortAnswerQuestionDto.ID = shortAnswerQuestion.ID
		shortAnswerQuestionDto.Title = shortAnswerQuestion.Title
		shortAnswerQuestionDto.Description = shortAnswerQuestion.Description
		shortAnswerQuestionDto.Credit = shortAnswerQuestion.Credit
		shortAnswerQuestionDto.CorrectAnswer = shortAnswerQuestion.CorrectAnswer
		// Get tags
		var tags []models.Tag
		err := models.DB.Model(&shortAnswerQuestion).Association("Tags").Find(&tags)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
			return
		}
		for _, tag := range tags {
			shortAnswerQuestionDto.Tags = append(shortAnswerQuestionDto.Tags, tag.Name)
		}
		responseDto.ShortAnswerQuestions = append(responseDto.ShortAnswerQuestions, shortAnswerQuestionDto)
	}
	c.IndentedJSON(http.StatusOK, responseDto)
}

// UpdateSummativeAssessment
// @Summary Update summative assessment
// @Description Update summative assessment
// @Tags SummativeAssessment
// @Accept  json
// @Produce  json
// @Param id path int true "SummativeAssessment ID"
// @Param summativeAssessment body models.SummativeAssessmentResponseDTO true "SummativeAssessment"
// @Success 200 {object} string "ok"
// @Router /summativeAssessments/{id} [put]
func UpdateSummativeAssessment(c *gin.Context) {
	var assessment models.SummativeAssessment
	var dto models.SummativeAssessmentRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get Assessment
	models.DB.Where("id = ?", c.Param("id")).First(&assessment)
	assessment.Title = dto.Title
	assessment.Difficulty = dto.Difficulty
	assessment.Points = dto.Points
	assessment.Topic = dto.Topic
	assessment.Randomisation = dto.Randomisation
	// Update assessment
	models.DB.Save(&assessment)
	// Delete all tags
	err := models.DB.Model(&assessment).Association("Tags").Clear()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when clear tags"})
		return
	}
	// Create tags
	for _, tagName := range dto.Tags {
		// if tag not exist, create tag
		var tag models.Tag
		models.DB.Where("name = ?", tagName).First(&tag)
		if tag.ID == 0 {
			tag.Name = tagName
			models.DB.Create(&tag)
		}
		err := models.DB.Model(&assessment).Association("Tags").Append(&tag)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when append tags"})
			return
		}
	}
	// Delete all multi choice questions
	err = models.DB.Model(&assessment).Association("MultiChoiceQuestions").Clear()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when clear multi choice questions"})
		return
	}
	// Set multi choice questions from database with multi choice question id
	for _, multiChoiceQuestionId := range dto.MultiChoiceQuestionIds {
		var multiChoiceQuestion models.MultiChoiceQuestion
		models.DB.Where("id = ?", multiChoiceQuestionId).First(&multiChoiceQuestion)
		err := models.DB.Model(&assessment).Association("MultiChoiceQuestions").Append(&multiChoiceQuestion)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when create multi choice questions"})
			return
		}
	}
	// Delete short answer questions
	err = models.DB.Model(&assessment).Association("ShortAnswerQuestions").Clear()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when delete short answer questions"})
		return
	}
	// Set short answer questions from database with short answer question id
	for _, shortAnswerQuestionId := range dto.ShortAnswerQuestionIds {
		var shortAnswerQuestion models.ShortAnswerQuestion
		models.DB.Where("id = ?", shortAnswerQuestionId).First(&shortAnswerQuestion)
		err := models.DB.Model(&assessment).Association("ShortAnswerQuestions").Append(&shortAnswerQuestion)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when create short answer questions"})
			return
		}
	}
	// Convert assessment to DTO
	var responseDTO models.SummativeAssessmentResponseDTO
	responseDTO.ID = assessment.ID
	responseDTO.Title = assessment.Title
	responseDTO.Difficulty = assessment.Difficulty
	responseDTO.Points = assessment.Points
	responseDTO.Topic = assessment.Topic
	responseDTO.Randomisation = assessment.Randomisation
	// Get MultiChoiceQuestions with ids from assessmentDTO multiChoiceQuestionIds
	var multiChoiceQuestions []models.MultiChoiceQuestion
	for _, multiChoiceQuestionId := range dto.MultiChoiceQuestionIds {
		var multiChoiceQuestion models.MultiChoiceQuestion
		models.DB.Where("id = ?", multiChoiceQuestionId).First(&multiChoiceQuestion)
		multiChoiceQuestions = append(multiChoiceQuestions, multiChoiceQuestion)
	}
	for _, multiChoiceQuestion := range multiChoiceQuestions {
		var multiChoiceQuestionDto models.MultiChoiceQuestionDTO
		multiChoiceQuestionDto.ID = multiChoiceQuestion.ID
		multiChoiceQuestionDto.Title = multiChoiceQuestion.Title
		multiChoiceQuestionDto.Description = multiChoiceQuestion.Description
		multiChoiceQuestionDto.Credit = multiChoiceQuestion.Credit
		// get multi choice answers
		var multiChoiceAnswers []models.MultiChoiceAnswer
		err := models.DB.Model(&multiChoiceQuestion).Association("Answers").Find(&multiChoiceAnswers)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch multi choice answers"})
			return
		}
		for _, multiChoiceAnswer := range multiChoiceAnswers {
			var multiChoiceAnswerDto models.MultiChoiceAnswerDTO
			multiChoiceAnswerDto.Content = multiChoiceAnswer.Content
			multiChoiceAnswerDto.IsCorrect = multiChoiceAnswer.IsCorrect
			multiChoiceQuestionDto.Answers = append(multiChoiceQuestionDto.Answers, multiChoiceAnswerDto)
		}

		// Get tags
		var tags []models.Tag
		err = models.DB.Model(&multiChoiceQuestion).Association("Tags").Find(&tags)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
			return
		}
		for _, tag := range tags {
			multiChoiceQuestionDto.Tags = append(multiChoiceQuestionDto.Tags, tag.Name)
		}
		responseDTO.MultiChoiceQuestions = append(responseDTO.MultiChoiceQuestions, multiChoiceQuestionDto)
	}
	// Get ShortAnswerQuestions
	var shortAnswerQuestions []models.ShortAnswerQuestion
	models.DB.Where("formative_assessment_id = ?", assessment.ID).Find(&shortAnswerQuestions)
	for _, shortAnswerQuestion := range shortAnswerQuestions {
		var shortAnswerQuestionDto models.ShortAnswerQuestionDTO
		shortAnswerQuestionDto.ID = shortAnswerQuestion.ID
		shortAnswerQuestionDto.Title = shortAnswerQuestion.Title
		shortAnswerQuestionDto.Description = shortAnswerQuestion.Description
		shortAnswerQuestionDto.Credit = shortAnswerQuestion.Credit
		shortAnswerQuestionDto.CorrectAnswer = shortAnswerQuestion.CorrectAnswer
		// Get tags
		var tags []models.Tag
		err := models.DB.Model(&shortAnswerQuestion).Association("Tags").Find(&tags)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
			return
		}
		for _, tag := range tags {
			shortAnswerQuestionDto.Tags = append(shortAnswerQuestionDto.Tags, tag.Name)
		}
		responseDTO.ShortAnswerQuestions = append(responseDTO.ShortAnswerQuestions, shortAnswerQuestionDto)
	}
	// Get tags
	var tags []models.Tag
	err = models.DB.Model(&assessment).Association("Tags").Find(&tags)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when fetch tags"})
		return
	}
	for _, tag := range tags {
		responseDTO.Tags = append(responseDTO.Tags, tag.Name)
	}
	c.IndentedJSON(http.StatusOK, responseDTO)
}

// DeleteSummativeAssessment godoc
// @Summary Delete summative assessment
// @Description Delete summative assessment
// @Tags SummativeAssessments
// @Accept  json
// @Produce  json
// @Param id path int true "Summative assessment id"
// @Success 200
// @Router /summative-assessments/{id} [delete]
func DeleteSummativeAssessment(c *gin.Context) {
	var assessment models.SummativeAssessment
	// Get assessment
	models.DB.Where("id = ?", c.Param("id")).First(&assessment)
	// Delete tags
	err := models.DB.Model(&assessment).Association("Tags").Clear()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when delete tags"})
		return
	}
	// Delete multi choice questions
	err = models.DB.Model(&assessment).Association("MultiChoiceQuestions").Clear()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when delete multi choice questions"})
		return
	}
	// Delete short answer questions
	err = models.DB.Model(&assessment).Association("ShortAnswerQuestions").Clear()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when delete short answer questions"})
		return
	}
	// Delete assessment
	models.DB.Delete(&assessment)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Summative assessment deleted successfully"})
}
