package main

import (
	"log"

	"github.com/galiherlangga/clinic-appointment/configs"
	"github.com/galiherlangga/clinic-appointment/handlers"
	"github.com/galiherlangga/clinic-appointment/models"
	"github.com/galiherlangga/clinic-appointment/repositories"
	"github.com/galiherlangga/clinic-appointment/routes"
	"github.com/galiherlangga/clinic-appointment/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Configuration
	configs.LoadConfig()

	// 2. Connect to Database
	configs.ConnectDB()

	// 2.1 Connect to Redis
	configs.ConnectRedis()

	// 3. Initial Migration (GORM AutoMigrate for convenience)
	err := configs.DB.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Provider{},
		&models.Service{},
		&models.Appointment{},
	)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// 4. Initialize Layers
	userRepo := repositories.NewUserRepository(configs.DB)
	blacklistService := services.NewBlacklistService()
	authService := services.NewAuthService(userRepo, blacklistService)
	authHandler := handlers.NewAuthHandler(authService)

	// 4.1 Initialize Clinic Layers
	cacheService := services.NewCacheService(configs.RedisClient)
	customerRepo := repositories.NewCustomerRepository(configs.DB)
	customerService := services.NewCustomerService(configs.DB, customerRepo, cacheService)
	customerHandler := handlers.NewCustomerHandler(customerService)

	providerRepo := repositories.NewProviderRepository(configs.DB)
	providerService := services.NewProviderService(configs.DB, providerRepo, cacheService)
	providerHandler := handlers.NewProviderHandler(providerService)

	appointmentRepo := repositories.NewAppointmentRepository(configs.DB)
	appointmentService := services.NewAppointmentService(configs.DB, appointmentRepo, cacheService)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)

	// 5. Setup Router
	r := gin.Default()

	// 5.1 Configure Trusted Proxies (Uncomment and set your proxy IPs/ranges)
	// Example for Cloudflare: r.SetTrustedProxies([]string{"173.245.48.0/20", "103.21.244.0/22", ...})
	// r.SetTrustedProxies(nil) // Trust all proxies (use with caution)

	// 6. Setup Routes
	routes.SetupRoutes(r, &routes.RouterContainer{
		AuthHandler:        authHandler,
		CustomerHandler:    customerHandler,
		ProviderHandler:    providerHandler,
		AppointmentHandler: appointmentHandler,
		BlacklistService:   blacklistService,
	})

	// 7. Start Server
	port := configs.AppConfig.AppPort
	if port != "" && port[0] != ':' {
		port = ":" + port
	}
	log.Println("Server starting on", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
