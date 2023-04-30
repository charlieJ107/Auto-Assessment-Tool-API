package main

import (
	"AAT_Api/controllers"
	"AAT_Api/docs"
	"AAT_Api/models"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

// @title           Auto Assessment Tool backend RESTful API
// @version         1.0
// @description     This is the backend RESTful API for the Auto Assessment Tool.

// @contact.name   Charlie
// @contact.email  charlie_j107+aat-backend-swagger@outlook.com

// @license.name  MPL-2.0
// @license.url   https://www.mozilla.org/en-US/MPL/2.0/

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	fmt.Println("Hello, Auto Assessment Tool by Group 5")
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	models.ConnectDatabase()
	v1 := router.Group("/api/v1")
	v1.GET("/multi-choice", controllers.GetMultiChoiceQuestions)
	v1.GET("/multi-choice/:id", controllers.GetMultiChoiceQuestion)
	v1.POST("/multi-choice", controllers.CreateMultiChoiceQuestion)
	v1.DELETE("/multi-choice/:id", controllers.DeleteMultiChoiceQuestion)
	v1.PUT("/multi-choice/:id", controllers.UpdateMultiChoiceQuestion)
	v1.GET("/short-answer", controllers.GetShortAnswerQuestions)
	v1.GET("/short-answer/:id", controllers.GetShortAnswerQuestion)
	v1.POST("/short-answer", controllers.CreateShortAnswerQuestion)
	v1.DELETE("/short-answer/:id", controllers.DeleteShortAnswerQuestion)
	v1.PUT("/short-answer/:id", controllers.UpdateShortAnswerQuestion)
	v1.GET("/formative-assessment", controllers.GetFormativeAssessments)
	v1.GET("/formative-assessment/:id", controllers.GetFormativeAssessment)
	v1.POST("/formative-assessment", controllers.CreateFormativeAssessment)
	v1.DELETE("/formative-assessment/:id", controllers.DeleteFormativeAssessment)
	v1.PUT("/formative-assessment/:id", controllers.UpdateFormativeAssessment)
	v1.GET("/summative-assessment", controllers.GetSummativeAssessments)
	v1.GET("/summative-assessment/:id", controllers.GetSummativeAssessment)
	v1.POST("/summative-assessment", controllers.CreateSummativeAssessment)
	v1.DELETE("/summative-assessment/:id", controllers.DeleteSummativeAssessment)
	v1.PUT("/summative-assessment/:id", controllers.UpdateSummativeAssessment)

	port, ok := os.LookupEnv("AAT_PORT")
	if !ok {
		port = "8080"
	}
	err := router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		return
	}
}
