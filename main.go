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

// @host rescuehub-production.up.railway.app
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

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
	)

	// psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	"localhost",
	// 	"5432",
	// 	"postgres",
	// 	"postgres",
	// 	"rescue_hub",
	// )

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
				protectedUserRoutes.GET("/", middlewares.RequireRoles("Akses ditolak, hanya admin dan donor yang dapat melihat daftar pengguna","admin", "donor"), controllers.GetAllUsers)
				protectedUserRoutes.GET("/:id", middlewares.RequireSameUserOrRole("Anda hanya bisa melihat info akun Anda sendiri", "admin"), controllers.GetUserByID)
				protectedUserRoutes.PUT("/:id", middlewares.RequireSameUserOrRole("Anda hanya bisa mengedit akun Anda sendiri", "admin"), controllers.UpdateUser)
				protectedUserRoutes.DELETE("/:id", middlewares.RequireSameUserOrRole("Anda hanya bisa menghapus akun Anda sendiri", "admin"), controllers.DeleteUser)
				protectedUserRoutes.PUT("/enable-2fa", middlewares.RequireSelfFor2FA(), controllers.Enable2FA)
				protectedUserRoutes.GET("/:id/donations", middlewares.RequireSameUserOrRole("Anda hanya bisa melihat donasi Anda sendiri", "admin"), controllers.GetDonationsByUserID)
				protectedUserRoutes.PUT("/:id/change-role", middlewares.RequireAdminOrSelfForRoleChange(), controllers.ChangeUserRole)
				protectedUserRoutes.PUT("/:id/edit-info", middlewares.RequireSameUserOrRole("Anda hanya bisa mengedit info akun Anda sendiri", "admin"), controllers.UpdateUserInfoWithoutEmail)
				protectedUserRoutes.GET("/:id/emergency-reports", middlewares.RequireSameUserOrRole("Anda hanya bisa melihat laporan darurat Anda sendiri", "admin"), controllers.GetEmergencyReportsByUserID)
			}
		}

		disasterRoutes := api.Group("/disasters", middlewares.JWTAuthMiddleware()) 
		{
			disasterRoutes.GET("/", controllers.GetAllDisasters)
			disasterRoutes.GET("/:id", controllers.GetDisasterByID)
			disasterRoutes.GET("/:id/refugees", controllers.GetRefugeesByDisasterID)
			disasterRoutes.GET("/:id/shelters", controllers.GetSheltersByDisasterID)
			disasterRoutes.GET("/:id/emergency-reports", controllers.GetEmergencyReportsByDisasterID)
			disasterRoutes.GET("/:id/evacuation-routes", controllers.GetEvacuationRoutesByDisasterID)

			disasterRoutes.POST("/", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa melaporkan bencana",
				"admin",
			), controllers.CreateDisaster)

			disasterRoutes.PUT("/:id", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa mengedit laporan bencana",
				"admin",
			), controllers.UpdateDisaster)

			disasterRoutes.DELETE("/:id", middlewares.RequireRoles(
				"Akses ditolak, hanya admin yang bisa menghapus laporan bencana",
				"admin",
			), controllers.DeleteDisaster)

			disasterRoutes.GET("/:id/volunteers", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa melihat daftar relawan",
				"admin",
			), controllers.GetVolunteersByDisasterID)

			disasterRoutes.GET("/:id/logistics", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa melihat daftar logistik",
				"admin",
			), controllers.GetLogisticsByDisasterID)
		}


		shelterRoutes := api.Group("/shelters", middlewares.JWTAuthMiddleware()) 
		{
			shelterRoutes.GET("/", controllers.GetAllShelters)
			shelterRoutes.GET("/:id", controllers.GetShelterByID)
			shelterRoutes.GET("/:id/refugees", controllers.GetRefugeesByShelterID)
			shelterRoutes.GET("/:id/logistics", controllers.GetLogisticsByShelterID)

			shelterRoutes.POST("/", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa menambahkan shelter",
				"admin",
			), controllers.CreateShelter)

			shelterRoutes.PUT("/:id", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa mengedit shelter",
				"admin",
			), controllers.UpdateShelter)

			shelterRoutes.DELETE("/:id", middlewares.RequireRoles(
				"Akses ditolak, hanya admin yang bisa menghapus shelter",
				"admin",
			), controllers.DeleteShelter)
		}

		refugeeRoutes := api.Group("/refugees", middlewares.JWTAuthMiddleware()) 
		{
		refugeeRoutes.GET("/", controllers.GetAllRefugees)
		refugeeRoutes.GET("/:id", controllers.GetRefugeeByID)

		refugeeRoutes.POST("/", middlewares.RequireVolunteerOrRole(
			"Akses ditolak, hanya admin dan relawan yang bisa mencatat pengungsi",
			"admin",
		), controllers.CreateRefugee)

		refugeeRoutes.PUT("/:id", middlewares.RequireRoles(
			"Akses ditolak, hanya admin yang bisa mengedit data pengungsi",
			"admin",
		), controllers.UpdateRefugee)

		refugeeRoutes.DELETE("/:id", middlewares.RequireRoles(
			"Akses ditolak, hanya admin yang bisa menghapus data pengungsi",
			"admin",
		), controllers.DeleteRefugee)
		}


		logisticRoutes := api.Group("/logistics", middlewares.JWTAuthMiddleware()) 
		{
			logisticRoutes.GET("/", controllers.GetAllLogistics)
			logisticRoutes.GET("/:id", controllers.GetLogisticByID)

			logisticRoutes.POST("/", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa mencatat logistik",
				"admin",
			), controllers.CreateLogistic)

			logisticRoutes.PUT("/:id", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa mengedit data logistik",
				"admin",
			), controllers.UpdateLogistic)

			logisticRoutes.DELETE("/:id", middlewares.RequireRoles(
				"Akses ditolak, hanya admin yang bisa menghapus data logistik",
				"admin",
			), controllers.DeleteLogistic)
		}

		distributionLogRoutes := api.Group("/distribution_logs", middlewares.JWTAuthMiddleware()) 
		{
			distributionLogRoutes.GET("/", controllers.GetAllDistributionLogs)
			distributionLogRoutes.GET("/:id", controllers.GetDistributionLogByID)

			distributionLogRoutes.POST("/", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa mencatat distribusi bantuan",
				"admin",
			), controllers.CreateDistributionLog)

			distributionLogRoutes.PUT("/:id", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa mengedit distribusi bantuan",
				"admin",
			), controllers.UpdateDistributionLog)

			distributionLogRoutes.DELETE("/:id", middlewares.RequireRoles(
				"Akses ditolak, hanya admin yang bisa menghapus distribusi bantuan",
				"admin",
			), controllers.DeleteDistributionLog)
		}


		evacuationRouteRoutes := api.Group("/evacuation_routes", middlewares.JWTAuthMiddleware()) 
		{
			evacuationRouteRoutes.GET("/", controllers.GetAllEvacuationRoutes)
			evacuationRouteRoutes.GET("/:id", controllers.GetEvacuationRouteByID)

			evacuationRouteRoutes.POST("/", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa mencatat jalur evakuasi",
				"admin",
			), controllers.CreateEvacuationRoute)

			evacuationRouteRoutes.PUT("/:id", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa mengedit jalur evakuasi",
				"admin",
			), controllers.UpdateEvacuationRoute)

			evacuationRouteRoutes.DELETE("/:id", middlewares.RequireRoles(
				"Akses ditolak, hanya admin yang bisa menghapus jalur evakuasi",
				"admin",
			), controllers.DeleteEvacuationRoute)
		}


		emergencyReportRoutes := api.Group("/emergency_reports", middlewares.JWTAuthMiddleware()) 
		{
			emergencyReportRoutes.POST("/", controllers.CreateEmergencyReport)

			emergencyReportRoutes.GET("/", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa melihat semua laporan darurat",
				"admin",
			), controllers.GetAllEmergencyReports)

			emergencyReportRoutes.GET("/:id", middlewares.RequireSelfForRelatedEntities(
				"Anda hanya bisa melihat laporan darurat Anda sendiri",
				"emergency_reports",
				"user_id",
			), controllers.GetEmergencyReportByID)

			emergencyReportRoutes.PUT("/:id", middlewares.RequireSelfForRelatedEntities(
				"Anda hanya bisa mengedit laporan darurat Anda sendiri",
				"emergency_reports",
				"user_id",
			), controllers.UpdateEmergencyReport)

			emergencyReportRoutes.DELETE("/:id", middlewares.RequireSelfForRelatedEntities(
				"Anda hanya bisa menghapus laporan darurat Anda sendiri",
				"emergency_reports",
				"user_id",
			), controllers.DeleteEmergencyReport)
		}

		donationRoutes := api.Group("/donations", middlewares.JWTAuthMiddleware()) 
		{
			donationRoutes.POST("/", controllers.CreateDonation)

			donationRoutes.GET("/", middlewares.RequireRoles(
				"Akses ditolak, hanya admin dan donatur yang bisa melihat semua donasi",
				"admin", "donor",
			), controllers.GetAllDonations)

			donationRoutes.GET("/:id", middlewares.RequireSelfForRelatedEntities(
				"Anda hanya bisa melihat donasi Anda sendiri",
				"donations",
				"donor_id",
			), controllers.GetDonationByID)

			donationRoutes.PUT("/:id", middlewares.RequireSelfForRelatedEntities(
				"Anda hanya bisa mengedit donasi Anda sendiri",
				"donations",
				"donor_id",
			), controllers.UpdateDonation)

			donationRoutes.DELETE("/:id", middlewares.RequireSelfForRelatedEntities(
				"Anda hanya bisa menghapus donasi Anda sendiri",
				"donations",
				"donor_id",
			), controllers.DeleteDonation)
		}

		volunteerRoutes := api.Group("/volunteers", middlewares.JWTAuthMiddleware()) 
		{
			volunteerRoutes.POST("/", controllers.CreateVolunteer)

			volunteerRoutes.GET("/", middlewares.RequireVolunteerOrRole(
				"Akses ditolak, hanya admin dan relawan yang bisa melihat semua relawan",
				"admin",
			), controllers.GetAllVolunteers)

			volunteerRoutes.GET("/:id", middlewares.RequireSelfForRelatedEntities(
				"Anda hanya bisa melihat informasi relawan Anda sendiri",
				"volunteers",
				"user_id",
			), controllers.GetVolunteerByID)

			volunteerRoutes.PUT("/:id", middlewares.RequireSelfForRelatedEntities(
				"Anda hanya bisa mengedit informasi relawan Anda sendiri",
				"volunteers",
				"user_id",
			), controllers.UpdateVolunteer)

			volunteerRoutes.DELETE("/:id", middlewares.RequireSelfForRelatedEntities(
				"Anda hanya bisa menghapus status relawan Anda sendiri",
				"volunteers",
				"user_id",
			), controllers.DeleteVolunteer)
		}
	}

	router.Run(":" + os.Getenv("PORT"))
}