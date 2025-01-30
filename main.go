package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"RescueHub/database"
	"RescueHub/controllers"
	"RescueHub/middlewares"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	"github.com/gin-contrib/cors"
	"time"
	_ "RescueHub/docs"
	"os"

	_ "github.com/lib/pq"
)

var (
	DB *sql.DB
	err error
)

// @title Rescue Hub API
// @version 1.0
// @description API	untuk mengelola pelaporan bencana, manajemen pengungsi, koordinasi relawan, dan distribusi logistik. Platform ini juga menyediakan rute evakuasi yang aman serta laporan kebutuhan darurat dari daerah terdampak bencana.
// @description Author: Rigel Ramadhani W. - Sanbercode Bootcamp Golang Batch 63 - FINAL PROJECT
// @contact.name Rigel Ramadhani W.
// @contact.url https://github.com/rigelra15

// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @security BearerAuth

func main() {
	err = godotenv.Load("config/.env")
	if err != nil {
		panic("Error loading .env file")
	}

	if os.Getenv("JWT_SECRET") == "" {
		panic("JWT_SECRET tidak ditemukan di environment variables")
	}

	// psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	os.Getenv("PGHOST"),
	// 	os.Getenv("PGPORT"),
	// 	os.Getenv("PGUSER"),
	// 	os.Getenv("PGPASSWORD"),
	// 	os.Getenv("PGDATABASE"),
	// )

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost",
		"5432",
		"postgres",
		"postgres",
		"rescue_hub",
	)

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	database.DBMigrate(DB)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		userRoutes := api.Group("/users")
		{
			userRoutes.POST("/", controllers.CreateUser)
			userRoutes.POST("/login", controllers.Login)
			userRoutes.POST("/verify-otp", controllers.VerifyOTP)
			
			protectedUserRoutes := userRoutes.Group("/", middlewares.JWTAuthMiddleware())
			{
				protectedUserRoutes.PUT("/enable-2fa", middlewares.RequireSelfFor2FA(), controllers.Enable2FA)
				protectedUserRoutes.PUT("/info/:id", middlewares.RequireSameUserOrRole("Anda hanya bisa mengedit info akun Anda sendiri", "admin"), controllers.UpdateUserInfoWithoutEmail)
				protectedUserRoutes.GET("/", middlewares.RequireRoles("Akses ditolak, hanya admin dan donor yang dapat melihat daftar pengguna","admin"), controllers.GetAllUsers)
				protectedUserRoutes.GET("/:id", middlewares.RequireSameUserOrRole("Anda hanya bisa melihat info akun Anda sendiri", "admin"), controllers.GetUserByID)
				protectedUserRoutes.PUT("/:id", middlewares.RequireSameUserOrRole("Anda hanya bisa mengedit akun Anda sendiri", "admin"), controllers.UpdateUser)
				protectedUserRoutes.DELETE("/:id", middlewares.RequireSameUserOrRole("Anda hanya bisa menghapus akun Anda sendiri", "admin"), controllers.DeleteUser)
			}
		}

		disasterRoutes := api.Group("/disasters")
		{
			disasterRoutes.POST("/", controllers.CreateDisaster)
			disasterRoutes.GET("/", controllers.GetAllDisasters)
			disasterRoutes.GET("/:id", controllers.GetDisasterByID)
			disasterRoutes.PUT("/:id", controllers.UpdateDisaster)
			disasterRoutes.DELETE("/:id", controllers.DeleteDisaster)
		}

		shelterRoutes := api.Group("/shelters")
		{
			shelterRoutes.POST("/", controllers.CreateShelter)
			shelterRoutes.GET("/", controllers.GetAllShelters)
			shelterRoutes.GET("/:id", controllers.GetShelterByID)
			shelterRoutes.PUT("/:id", controllers.UpdateShelter)
			shelterRoutes.DELETE("/:id", controllers.DeleteShelter)
		}
		refugeeRoutes := api.Group("/refugees")
		{
			refugeeRoutes.POST("/", controllers.CreateRefugee)
			refugeeRoutes.GET("/", controllers.GetAllRefugees)
			refugeeRoutes.GET("/:id", controllers.GetRefugeeByID)
			refugeeRoutes.PUT("/:id", controllers.UpdateRefugee)
			refugeeRoutes.DELETE("/:id", controllers.DeleteRefugee)
		}
		logisticRoutes := api.Group("/logistics")
		{
			logisticRoutes.POST("/", controllers.CreateLogistic)
			logisticRoutes.GET("/", controllers.GetAllLogistics)
			logisticRoutes.GET("/:id", controllers.GetLogisticByID)
			logisticRoutes.PUT("/:id", controllers.UpdateLogistic)
			logisticRoutes.DELETE("/:id", controllers.DeleteLogistic)
		}
		distributionLogRoutes := api.Group("/distribution_logs")
		{
			distributionLogRoutes.POST("/", controllers.CreateDistributionLog)
			distributionLogRoutes.GET("/", controllers.GetAllDistributionLogs)
			distributionLogRoutes.GET("/:id", controllers.GetDistributionLogByID)
			distributionLogRoutes.PUT("/:id", controllers.UpdateDistributionLog)
			distributionLogRoutes.DELETE("/:id", controllers.DeleteDistributionLog)
		}
		evacuationRouteRoutes := api.Group("/evacuation_routes")
		{
			evacuationRouteRoutes.POST("/", controllers.CreateEvacuationRoute)
			evacuationRouteRoutes.GET("/", controllers.GetAllEvacuationRoutes)
			evacuationRouteRoutes.GET("/:id", controllers.GetEvacuationRouteByID)
			evacuationRouteRoutes.PUT("/:id", controllers.UpdateEvacuationRoute)
			evacuationRouteRoutes.DELETE("/:id", controllers.DeleteEvacuationRoute)
		}
		emergencyReportRoutes := api.Group("/emergency_reports")
		{
			emergencyReportRoutes.POST("/", controllers.CreateEmergencyReport)
			emergencyReportRoutes.GET("/", controllers.GetAllEmergencyReports)
			emergencyReportRoutes.GET("/:id", controllers.GetEmergencyReportByID)
			emergencyReportRoutes.PUT("/:id", controllers.UpdateEmergencyReport)
			emergencyReportRoutes.DELETE("/:id", controllers.DeleteEmergencyReport)
		}
		donationRoutes := api.Group("/donations")
		{
			donationRoutes.POST("/", controllers.CreateDonation)
			donationRoutes.GET("/", controllers.GetAllDonations)
			donationRoutes.GET("/:id", controllers.GetDonationByID)
			donationRoutes.PUT("/:id", controllers.UpdateDonation)
			donationRoutes.DELETE("/:id", controllers.DeleteDonation)
		}
		volunteerRoutes := api.Group("/volunteers")
		{
			volunteerRoutes.POST("/", controllers.CreateVolunteer)
			volunteerRoutes.GET("/", controllers.GetAllVolunteers)
			volunteerRoutes.GET("/:id", controllers.GetVolunteerByID)
			volunteerRoutes.PUT("/:id", controllers.UpdateVolunteer)
			volunteerRoutes.DELETE("/:id", controllers.DeleteVolunteer)
		}
	}

	router.Run(":" + os.Getenv("PORT"))
}