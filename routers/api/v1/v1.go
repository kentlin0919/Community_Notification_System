package v1

import (
	v1 "Community_Notification_System/app/controller/v1"
	"Community_Notification_System/app/controller/v1/auth"
	"Community_Notification_System/middlewares"

	"github.com/gin-gonic/gin"
)

func V1PublicRoutes(rg *gin.RouterGroup) {
	// 公開路由
	rg.POST("/login", v1.User().UserLogin)
	rg.POST("/register", v1.User().UserRegister)
	rg.GET("/platform/getlist", v1.Platform().Platform_GetList)
	rg.GET("/system/config", v1.Platform().SystemConfig)

	// 認證相關 (Public)
	authCtrl := auth.NewAuthController()
	rg.POST("/auth/refresh", authCtrl.RefreshToken)
	rg.POST("/auth/forgot-password", authCtrl.ForgotPassword)
	rg.POST("/auth/verify-otp", authCtrl.VerifyOTP)
	rg.POST("/auth/reset-password", authCtrl.ResetPassword)
}

func V1PrivateRoutes(rg *gin.RouterGroup) {
	// ── 一般登入使用者即可（無社區隔離限制） ──
	rg.POST("/deleteUser", v1.User().UserDelete)
	rg.POST("/community/register", v1.CommunityManager().CommunityManager_Register)
	rg.GET("/home", v1.Home().GetDashboard)

	// 認證相關 (Private)
	authCtrl := auth.NewAuthController()
	rg.POST("/auth/switch-community", authCtrl.SwitchCommunity)
	rg.POST("/auth/logout", authCtrl.Logout)

	// ── 超級管理員專屬 (Super Admin Only, Level 1) ──
	superAdmin := rg.Group("")
	superAdmin.Use(middlewares.MinimalPermissionMiddleware(1))
	{
		superAdmin.GET("/community/getlist", v1.CommunityManager().CommunityManager_GetList)
		superAdmin.PATCH("/community/register/:id/approve", v1.CommunityManager().CommunityManager_Approve)
		superAdmin.PATCH("/community/register/:id/reject", v1.CommunityManager().CommunityManager_Reject)

		// ── 與前端對齊的 Super Admin 專屬路由 ──
		superAdmin.GET("/super-admin/communities", v1.CommunityManager().CommunityManager_GetList)
		superAdmin.GET("/super-admin/community-applications", v1.CommunityManager().CommunityManager_GetApplicationList)
		superAdmin.PATCH("/super-admin/community-applications/:id/approve", v1.CommunityManager().CommunityManager_Approve)
		superAdmin.PATCH("/super-admin/community-applications/:id/reject", v1.CommunityManager().CommunityManager_Reject)
	}

	// ── 社區住戶級（Resident+, Level 8，有社區資料隔離） ──
	residentScope := rg.Group("")
	residentScope.Use(middlewares.CommunityContextMiddleware())
	{
		residentScope.GET("/messages", v1.Message().GetMessageList)
		residentScope.PATCH("/messages/:id/read", v1.Message().MarkMessageRead)
		residentScope.PATCH("/messages/read-all", v1.Message().MarkAllMessagesRead)
		residentScope.GET("/permissions/profile", v1.Permission().GetCommunityPermissionProfiles)
		residentScope.GET("/facilities", v1.Facility().GetFacilityList)
		residentScope.GET("/facilities/:id", v1.Facility().GetFacilityDetail)
		residentScope.GET("/reservations", v1.Reservation().GetReservationList)
		residentScope.GET("/reservations/:id", v1.Reservation().GetReservationDetail)
		residentScope.POST("/facilities/:facility_id/reservations", v1.Reservation().CreateReservation)
		residentScope.PATCH("/reservations/:id/cancel", v1.Reservation().CancelReservation)
		residentScope.POST("/reservations/:id/reschedule", v1.Reservation().Reschedule)
		residentScope.GET("/parcels", v1.Parcel().GetParcelList)
		residentScope.PUT("/parcels/:id/pickup", v1.Parcel().PickupParcel)
	}

	// ── 社區管理人員級（Staff+, Level 7，有社區資料隔離） ──
	staffScope := rg.Group("")
	staffScope.Use(middlewares.CommunityContextMiddleware(), middlewares.MinimalPermissionMiddleware(7))
	{
		staffScope.PUT("/admin/permissions/profile", v1.Permission().UpdateCommunityPermissionProfile)
		staffScope.POST("/sendmessage", v1.Message().SendMessage)
		staffScope.POST("/messages/send", v1.Message().SendMessage)
		staffScope.POST("/parcels", v1.Parcel().CreateParcel)
	}

	// ── 高級管理人員級（Admin, Level 2，有社區資料隔離） ──
	adminScope := rg.Group("")
	adminScope.Use(middlewares.CommunityContextMiddleware(), middlewares.MinimalPermissionMiddleware(2))
	{
		adminScope.POST("/admin/facilities", v1.Facility().CreateFacility)
		adminScope.PUT("/admin/facilities/:id", v1.Facility().UpdateFacility)
		adminScope.DELETE("/admin/facilities/:id", v1.Facility().DeleteFacility)
		adminScope.PATCH("/admin/reservations/:id/delete", v1.Reservation().AdminDeleteReservation)
		adminScope.PATCH("/admin/reschedule/:id/approve", v1.Reservation().AdminApproveReschedule)
	}
}
