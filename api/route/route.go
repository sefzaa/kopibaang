package route

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"kopibang/api/controller"
	"kopibang/api/middleware"
	"kopibang/bootstrap"
	"kopibang/repository"
	"kopibang/usecase"
	_ "kopibang/docs"
)

// Setup melakukan inisialisasi semua layer dan mendaftarkan route
func Setup(env *bootstrap.Env, timeout time.Duration, app *bootstrap.Application, r *gin.Engine) {
	// 1. Setup Firebase Client
	fcmClient := bootstrap.NewFirebaseMessagingClient(env)

	// 2. Inisialisasi Repositories (Data Layer)
	userRepo := repository.NewUserRepository(app.DB)
	productRepo := repository.NewProductRepository(app.DB)
	txRepo := repository.NewTransactionRepository(app.DB)
	rawMaterialRepo := repository.NewRawMaterialRepository(app.DB)
	voucherRepo := repository.NewVoucherRepository(app.DB)
	dashboardRepo := repository.NewDashboardRepository(app.DB)
	settingRepo := repository.NewSettingRepository(app.DB)
	redisRepo := repository.NewRedisRepository(app.Redis)

	// 3. Inisialisasi Usecases (Business Logic Layer)
	authUc := usecase.NewAuthUsecase(userRepo, redisRepo, env)
	userUc := usecase.NewUserUsecase(userRepo)
	
	// TAMBAHAN: Masukkan app.Minio dan env ke ProductUsecase
	productUc := usecase.NewProductUsecase(productRepo, fcmClient, app.Minio, env)
	
	txUc := usecase.NewTransactionUsecase(txRepo, productRepo, voucherRepo, redisRepo)
	dashboardUc := usecase.NewDashboardUsecase(dashboardRepo)
	rawMaterialUc := usecase.NewRawMaterialUsecase(rawMaterialRepo)
	voucherUc := usecase.NewVoucherUsecase(voucherRepo, fcmClient)
	settingUc := usecase.NewSettingUsecase(settingRepo, redisRepo, fcmClient)

	// 4. Inisialisasi Controllers (Presentation Layer)
	authController := controller.NewAuthController(authUc)
	userController := controller.NewUserController(userUc)
	productController := controller.NewProductController(productUc)
	txController := controller.NewTransactionController(txUc)
	dashboardController := controller.NewDashboardController(dashboardUc)
	rawMaterialController := controller.NewRawMaterialController(rawMaterialUc)
	voucherController := controller.NewVoucherController(voucherUc)
	settingController := controller.NewSettingController(settingUc)
	notificationController := controller.NewNotificationController(fcmClient)

	// 5. Setup Routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		// --- 1. PUBLIC ROUTES (Tanpa Login) ---
		v1.GET("/settings/barista-status", settingController.GetBaristaStatus)

		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.RegisterCustomer)
			auth.POST("/login", authController.LoginCustomer)
			auth.POST("/refresh", authController.RefreshToken)
			auth.POST("/forgot-password", authController.ForgotPassword)
			auth.POST("/verify-otp", authController.VerifyOTP)
			auth.POST("/reset-password", authController.ResetPassword)
		}

		adminAuthPub := v1.Group("/admin/auth")
		{
			adminAuthPub.POST("/login", authController.LoginAdmin)
		}

		// --- 2. PROTECTED ROUTES (Bisa Diakses Customer & Barista) ---
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(""))
		{
			protected.POST("/auth/logout", authController.Logout)
			
			protected.GET("/profile", userController.GetProfile)
			protected.PUT("/profile", userController.UpdateProfile)

			// Menu Catalog
			protected.GET("/menus", productController.GetMenus)
			protected.GET("/menus/:id", productController.GetMenuDetail)

			// Points & Transaction
			protected.POST("/points/redeem-qr", txController.RequestRedeemQR)
			protected.POST("/points/scan-earn", txController.ScanEarnQR)
		}

		// --- 3. ADMIN ROUTES (Hanya Barista) ---
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware("barista"))
		{
			// Dashboard Analytics
			admin.GET("/dashboard", dashboardController.GetAdminDashboard)
			admin.PATCH("/settings/barista-status", settingController.PatchBaristaStatus)

			// TAMBAHAN: Route Upload File ke MinIO
			admin.POST("/upload", productController.UploadImage)

			// CRUD Menu
			admin.POST("/menus", productController.CreateMenu)
			admin.PUT("/menus/:id", productController.UpdateMenu)
			admin.PATCH("/menus/:id/status", productController.ToggleMenuStatus) 
			admin.DELETE("/menus/:id", productController.DeleteMenu)

			// POS / Kasir
			admin.POST("/orders", txController.CreateOrder)

			// CRUD Raw Materials
			admin.POST("/raw-materials", rawMaterialController.AddMaterial)
			admin.PUT("/raw-materials/:id", rawMaterialController.UpdateMaterial)
			admin.DELETE("/raw-materials/:id", rawMaterialController.DeleteMaterial)
			admin.GET("/raw-materials", rawMaterialController.GetAllMaterials)

			// CRUD Vouchers
			admin.POST("/vouchers", voucherController.CreateVoucher)
			admin.PUT("/vouchers/:id", voucherController.UpdateVoucher)
			admin.DELETE("/vouchers/:id", voucherController.DeleteVoucher)
			admin.GET("/vouchers", voucherController.GetAllVouchers)

			// Notifications
			admin.POST("/notifications/push", notificationController.SendCustomPush)
		}
	}
}