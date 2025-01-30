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
				protectedUserRoutes.GET("/", middlewares.RequireRoles("Akses ditolak, hanya admin dan donor yang dapat melihat daftar pengguna","admin", "donor"), controllers.GetAllUsers)
				protectedUserRoutes.GET("/:id", middlewares.RequireSameUserOrRole("Anda hanya bisa melihat info akun Anda sendiri", "admin"), controllers.GetUserByID)
				protectedUserRoutes.PUT("/:id", middlewares.RequireSameUserOrRole("Anda hanya bisa mengedit akun Anda sendiri", "admin"), controllers.UpdateUser)
				protectedUserRoutes.DELETE("/:id", middlewares.RequireSameUserOrRole("Anda hanya bisa menghapus akun Anda sendiri", "admin"), controllers.DeleteUser)
				protectedUserRoutes.PUT("/enable-2fa", middlewares.RequireSelfFor2FA(), controllers.Enable2FA)
				protectedUserRoutes.PUT("/:id/edit-info", middlewares.RequireSameUserOrRole("Anda hanya bisa mengedit info akun Anda sendiri", "admin"), controllers.UpdateUserInfoWithoutEmail)
				protectedUserRoutes.GET("/:id/emergency-reports", middlewares.RequireSameUserOrRole("Anda hanya bisa melihat laporan darurat Anda sendiri", "user"), controllers.GetEmergencyReportsByUserID)
			}
		}

		disasterRoutes := api.Group("/disasters")
		{
			disasterRoutes.GET("/", controllers.GetAllDisasters)
			disasterRoutes.GET("/:id", controllers.GetDisasterByID)
			disasterRoutes.GET("/:id/emergency-reports", controllers.GetEmergencyReportsByDisasterID)
			disasterRoutes.GET("/:id/evacuation-routes", controllers.GetEvacuationRoutesByDisasterID)

			protectedDisasterRoutes := disasterRoutes.Group("/", middlewares.JWTAuthMiddleware())
			{
				protectedDisasterRoutes.POST("/", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa melaporkan bencana",
					"admin",
				), controllers.CreateDisaster)

				protectedDisasterRoutes.PUT("/:id", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa mengedit laporan bencana",
					"admin",
				), controllers.UpdateDisaster)

				protectedDisasterRoutes.DELETE("/:id", middlewares.RequireRoles(
					"Akses ditolak, hanya admin yang bisa menghapus laporan bencana",
					"admin",
				), controllers.DeleteDisaster)

				protectedDisasterRoutes.GET("/:id/shelters", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa melihat daftar shelter",
					"admin",
				), controllers.GetSheltersByDisasterID)

				protectedDisasterRoutes.GET("/:id/volunteers", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa melihat daftar relawan",
					"admin",
				), controllers.GetVolunteersByDisasterID)

				protectedDisasterRoutes.GET("/:id/logistics", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa melihat daftar logistik",
					"admin",
				), controllers.GetLogisticsByDisasterID)		
			}
		}

		shelterRoutes := api.Group("/shelters")
		{
			shelterRoutes.GET("/", controllers.GetAllShelters)
			shelterRoutes.GET("/:id", controllers.GetShelterByID)
			shelterRoutes.GET("/:id/refugees", controllers.GetRefugeesByShelterID)
			shelterRoutes.GET("/:id/logistics", controllers.GetLogisticsByShelterID)

			protectedShelterRoutes := shelterRoutes.Group("/", middlewares.JWTAuthMiddleware())
			{
				protectedShelterRoutes.POST("/", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa menambahkan shelter",
					"admin",
				), controllers.CreateShelter)

				protectedShelterRoutes.PUT("/:id", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa mengedit shelter",
					"admin",
				), controllers.UpdateShelter)

				protectedShelterRoutes.DELETE("/:id", middlewares.RequireRoles(
					"Akses ditolak, hanya admin yang bisa menghapus shelter",
					"admin",
				), controllers.DeleteShelter)
			}
		}

		refugeeRoutes := api.Group("/refugees")
		{
			refugeeRoutes.GET("/", controllers.GetAllRefugees)
			refugeeRoutes.GET("/:id", controllers.GetRefugeeByID)
			refugeeRoutes.GET("/:id/distribution-logs", controllers.GetDistributionLogsByRefugeeID)

			protectedRefugeeRoutes := refugeeRoutes.Group("/", middlewares.JWTAuthMiddleware())
			{
				protectedRefugeeRoutes.POST("/", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa mencatat pengungsi",
					"admin",
				), controllers.CreateRefugee)

				protectedRefugeeRoutes.PUT("/:id", middlewares.RequireRoles(
					"Akses ditolak, hanya admin yang bisa mengedit data pengungsi",
					"admin",
				), controllers.UpdateRefugee)

				protectedRefugeeRoutes.DELETE("/:id", middlewares.RequireRoles(
					"Akses ditolak, hanya admin yang bisa menghapus data pengungsi",
					"admin",
				), controllers.DeleteRefugee)
			}
		}

		logisticRoutes := api.Group("/logistics")
		{
			logisticRoutes.GET("/", controllers.GetAllLogistics)
			logisticRoutes.GET("/:id", controllers.GetLogisticByID)

			protectedLogisticRoutes := logisticRoutes.Group("/", middlewares.JWTAuthMiddleware())
			{
				protectedLogisticRoutes.POST("/", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa mencatat logistik",
					"admin",
				), controllers.CreateLogistic)

				protectedLogisticRoutes.PUT("/:id", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa mengedit data logistik",
					"admin",
				), controllers.UpdateLogistic)

				protectedLogisticRoutes.DELETE("/:id", middlewares.RequireRoles(
					"Akses ditolak, hanya admin yang bisa menghapus data logistik",
					"admin",
				), controllers.DeleteLogistic)
			}
		}

		distributionLogRoutes := api.Group("/distribution_logs")
		{
			distributionLogRoutes.GET("/", controllers.GetAllDistributionLogs)
			distributionLogRoutes.GET("/:id", controllers.GetDistributionLogByID)

			protectedDistributionLogRoutes := distributionLogRoutes.Group("/", middlewares.JWTAuthMiddleware())
			{
				protectedDistributionLogRoutes.POST("/", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa mencatat distribusi bantuan",
					"admin",
				), controllers.CreateDistributionLog)

				protectedDistributionLogRoutes.PUT("/:id", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa mengedit distribusi bantuan",
					"admin",
				), controllers.UpdateDistributionLog)

				protectedDistributionLogRoutes.DELETE("/:id", middlewares.RequireRoles(
					"Akses ditolak, hanya admin yang bisa menghapus distribusi bantuan",
					"admin",
				), controllers.DeleteDistributionLog)
			}
		}

		evacuationRouteRoutes := api.Group("/evacuation_routes")
		{
			evacuationRouteRoutes.GET("/", controllers.GetAllEvacuationRoutes)
			evacuationRouteRoutes.GET("/:id", controllers.GetEvacuationRouteByID)

			protectedEvacuationRouteRoutes := evacuationRouteRoutes.Group("/", middlewares.JWTAuthMiddleware())
			{
				protectedEvacuationRouteRoutes.POST("/", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa mencatat jalur evakuasi",
					"admin",
				), controllers.CreateEvacuationRoute)

				protectedEvacuationRouteRoutes.PUT("/:id", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa mengedit jalur evakuasi",
					"admin",
				), controllers.UpdateEvacuationRoute)

				protectedEvacuationRouteRoutes.DELETE("/:id", middlewares.RequireRoles(
					"Akses ditolak, hanya admin yang bisa menghapus jalur evakuasi",
					"admin",
				), controllers.DeleteEvacuationRoute)
			}
		}


		emergencyReportRoutes := api.Group("/emergency_reports")
		{
			protectedEmergencyReportRoutes := emergencyReportRoutes.Group("/", middlewares.JWTAuthMiddleware())
			{
				protectedEmergencyReportRoutes.POST("/", controllers.CreateEmergencyReport)

				protectedEmergencyReportRoutes.GET("/", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa melihat semua laporan darurat",
					"admin",
				), controllers.GetAllEmergencyReports)

				protectedEmergencyReportRoutes.GET("/:id", middlewares.RequireSameUserOrRole(
					"Anda hanya bisa melihat laporan darurat Anda sendiri",
					"admin",
				), controllers.GetEmergencyReportByID)

				protectedEmergencyReportRoutes.PUT("/:id", middlewares.RequireSameUserOrRole(
					"Anda hanya bisa mengedit laporan darurat Anda sendiri",
					"admin",
				), controllers.UpdateEmergencyReport)

				protectedEmergencyReportRoutes.DELETE("/:id", middlewares.RequireSameUserOrRole(
					"Anda hanya bisa menghapus laporan darurat Anda sendiri",
					"admin",
				), controllers.DeleteEmergencyReport)
			}
		}


		donationRoutes := api.Group("/donations")
		{
			protectedDonationRoutes := donationRoutes.Group("/", middlewares.JWTAuthMiddleware())
			{
				protectedDonationRoutes.POST("/", controllers.CreateDonation)

				protectedDonationRoutes.GET("/", middlewares.RequireRoles(
					"Akses ditolak, hanya admin dan donatur yang bisa melihat semua donasi",
					"admin", "donor",
				), controllers.GetAllDonations)

				protectedDonationRoutes.GET("/:id", middlewares.RequireSameUserOrRole(
					"Anda hanya bisa melihat donasi Anda sendiri",
					"admin",
				), controllers.GetDonationByID)

				protectedDonationRoutes.PUT("/:id", middlewares.RequireSameUserOrRole(
					"Anda hanya bisa mengedit donasi Anda sendiri",
					"admin",
				), controllers.UpdateDonation)

				protectedDonationRoutes.DELETE("/:id", middlewares.RequireSameUserOrRole(
					"Anda hanya bisa menghapus donasi Anda sendiri",
					"admin",
				), controllers.DeleteDonation)
			}
		}

		volunteerRoutes := api.Group("/volunteers")
		{
			protectedVolunteerRoutes := volunteerRoutes.Group("/", middlewares.JWTAuthMiddleware())
			{
				protectedVolunteerRoutes.POST("/", controllers.CreateVolunteer)

				protectedVolunteerRoutes.GET("/", middlewares.RequireVolunteerOrRole(
					"Akses ditolak, hanya admin dan relawan yang bisa melihat semua relawan",
					"admin",
				), controllers.GetAllVolunteers)

				protectedVolunteerRoutes.GET("/:id", middlewares.RequireSameUserOrRole(
					"Anda hanya bisa melihat informasi relawan Anda sendiri",
					"admin",
				), controllers.GetVolunteerByID)

				protectedVolunteerRoutes.PUT("/:id", middlewares.RequireSameUserOrRole(
					"Anda hanya bisa mengedit informasi relawan Anda sendiri",
					"admin",
				), controllers.UpdateVolunteer)

				protectedVolunteerRoutes.DELETE("/:id", middlewares.RequireSameUserOrRole(
					"Anda hanya bisa menghapus status relawan Anda sendiri",
					"admin",
				), controllers.DeleteVolunteer)
			}
		}
	}

	router.Run(":" + os.Getenv("PORT"))
}