package controllers

import (
	"AAT_Api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// BasePath /api/v1

// GetFormativeAssessments
// @Summary Get all formative assessments
// @Description Get all formative assessments
// @Tags FormativeAssessments
// @Produce json
// @Param tags query string false "Tags to filter assessments"
// @Success 200 {array} models.FormativeAssessmentResponseDTO
// @Router /formative-assessments [get]
func GetFormativeAssessments(c *gin.Context) {
	var assessments []models.FormativeAssessment
	// resolve tags query from context
	tags := c.QueryArray("tags")
	// Get all assessments with multiple optional tags
	if len(tags) > 0 {
		models.DB.Where("tags.name IN ?", tags).Joins("JOIN assessment_tags ON assessment_tags.assessment_id = formative_assessments.id").Joins("JOIN tags ON tags.id = assessment_tags.tag_id").Find(&assessments)
	} else {
		models.DB.Find(&assessments)
	}
	// Convert assessments to DTO
	var dtos []models.FormativeAssessmentResponseDTO
	for _, assessment := range assessments {
		var dto models.FormativeAssessmentResponseDTO
		dto.ID = assessment.ID
		dto.Title = assessment.Title
		dto.Difficulty = assessment.Difficulty
		dto.Points = assessment.Points
		dto.Topic = assessment.Topic
		dto.Randomisation = assessment.Randomisation
		// Get MultiChoiceQuestions
		var multiChoiceQuestions []models.MultiChoiceQuestion
		models.DB.Where("formative_assessment_id = ?", assessment.ID).Find(&multiChoiceQuestions)
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
			dto.MultiChoiceQuestions = append(dto.MultiChoiceQuestions, multiChoiceQuestionDto)
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

// GetFormativeAssessment
// @Summary Get a formative assessment
// @Description Get a formative assessment
// @Tags FormativeAssessments
// @Produce json
// @Param id path uint true "Formative assessment ID"
// @Success 200 {object} models.FormativeAssessmentResponseDTO
// @Router /formative-assessments/{id} [get]
func GetFormativeAssessment(c *gin.Context) {
	var assessment models.FormativeAssessment
	id := c.Param("id")
	err := models.DB.Where("id = ?", id).First(&assessment).Error
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Formative assessment not found"})
		return
	}
	// Convert assessment to DTO
	var dto models.FormativeAssessmentResponseDTO
	dto.ID = assessment.ID
	dto.Title = assessment.Title
	dto.Difficulty = assessment.Difficulty
	dto.Points = assessment.Points
	dto.Topic = assessment.Topic
	dto.Randomisation = assessment.Randomisation
	// Get MultiChoiceQuestions
	var multiChoiceQuestions []models.MultiChoiceQuestion
	models.DB.Where("formative_assessment_id = ?", assessment.ID).Find(&multiChoiceQuestions)
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
		dto.MultiChoiceQuestions = append(dto.MultiChoiceQuestions, multiChoiceQuestionDto)
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

// CreateFormativeAssessment
// @Summary Create a formative assessment
// @Description Create a formative assessment
// @Tags FormativeAssessments
// @Produce json
// @Param formativeAssessment body models.FormativeAssessmentResponseDTO true "Formative assessment"
// @Success 201 {object} models.FormativeAssessmentResponseDTO
// @Router /formative-assessments [post]
func CreateFormativeAssessment(c *gin.Context) {
	var dto models.FormativeAssessmentRequestDTO
	err := c.BindJSON(&dto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error when bind json"})
		return
	}
	var assessment models.FormativeAssessment
	assessment.Title = dto.Title
	assessment.Difficulty = dto.Difficulty
	assessment.Points = dto.Points
	assessment.Topic = dto.Topic
	assessment.Randomisation = dto.Randomisation

	// Create formative assessment
	models.DB.Create(&assessment)

	// Create tags
	for _, tagName := range dto.Tags {
		// check if tag exists, if not create it
		var tag models.Tag
		models.DB.Where("name = ?", tagName).First(&tag)
		if tag.ID == 0 {
			tag.Name = tagName
			models.DB.Create(&tag)
		}
		err := models.DB.Model(&assessment).Association("Tags").Append(&tag)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when create tags"})
			return
		}
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
	var responseDTO models.FormativeAssessmentResponseDTO
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

	// Get ShortAnswerQuestions from assessment entity
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
	c.IndentedJSON(http.StatusCreated, responseDTO)
}

// UpdateFormativeAssessment
// @Summary Update a formative assessment
// @Description Update a formative assessment
// @Tags FormativeAssessments
// @Produce json
// @Param id path int true "Formative assessment ID"
// @Param formativeAssessment body models.FormativeAssessmentResponseDTO true "Formative assessment"
// @Success 200 {object} models.FormativeAssessmentResponseDTO
// @Router /formative-assessments/{id} [put]
func UpdateFormativeAssessment(c *gin.Context) {
	var assessment models.FormativeAssessment
	var dto models.FormativeAssessmentRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get assessment
	models.DB.Where("id = ?", c.Param("id")).First(&assessment)
	assessment.Title = dto.Title
	assessment.Difficulty = dto.Difficulty
	assessment.Points = dto.Points
	assessment.Topic = dto.Topic
	assessment.Randomisation = dto.Randomisation
	// Update assessment
	models.DB.Save(&assessment)
	// Delete tags
	err := models.DB.Model(&assessment).Association("Tags").Clear()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when delete tags"})
		return
	}
	// Create tags
	for _, tagName := range dto.Tags {
		// if tag not exist, create new tag
		var tag models.Tag
		models.DB.Where("name = ?", tagName).First(&tag)
		if tag.ID == 0 {
			tag.Name = tagName
			models.DB.Create(&tag)
		}
		err := models.DB.Model(&assessment).Association("Tags").Append(&tag)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when create tags"})
			return
		}
	}
	// Delete multi choice questions
	err = models.DB.Model(&assessment).Association("MultiChoiceQuestions").Clear()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Error when delete multi choice questions"})
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
	var responseDTO models.FormativeAssessmentResponseDTO
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

// DeleteFormativeAssessment godoc
// @Summary Delete formative assessment
// @Description Delete formative assessment
// @Tags FormativeAssessment
// @Accept  json
// @Produce  json
// @Param id path int true "Formative assessment id"
// @Success 200 {object} models.FormativeAssessmentResponseDTO
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Router /formative-assessments/{id} [delete]
func DeleteFormativeAssessment(c *gin.Context) {
	var assessment models.FormativeAssessment

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
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
