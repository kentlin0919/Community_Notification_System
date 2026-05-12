package v1

import (
	v1 "Community_Notification_System/app/controller/v1"

	"github.com/gin-gonic/gin"
)

func V1Routes(rg *gin.RouterGroup) {

	// 處理登入請求
	rg.POST("/login", v1.User().UserLogin)
	/// 處理註冊請求
	rg.POST("/register", v1.User().UserRegister)

	/// 刪除使用者
	rg.POST("/deleteUser", v1.User().UserDelete)

	//處理送通知
	rg.POST("/sendmessage", v1.Message().SendMessage)
	rg.GET("/messages", v1.Message().GetMessageList)
	rg.PATCH("/messages/:id/read", v1.Message().MarkMessageRead)
	rg.PATCH("/messages/read-all", v1.Message().MarkAllMessagesRead)

	// 取得社區列表
	rg.GET("/community/getlist", v1.CommunityManager().CommunityManager_GetList)

	// 新增社區（提交申請）
	rg.POST("/community/register", v1.CommunityManager().CommunityManager_Register)

	// 社區申請單審核（Super admin only）
	rg.PATCH("/community/register/:id/approve", v1.CommunityManager().CommunityManager_Approve)
	rg.PATCH("/community/register/:id/reject", v1.CommunityManager().CommunityManager_Reject)

	// 社區自有角色設定 (權限 2+)
	// 請確保前端請求這些 API 時有帶 JWT，因為我們需要從 JWT 撈取 community_id 來決定改哪一棟的設定
	rg.PUT("/admin/permissions/profile", v1.Permission().UpdateCommunityPermissionProfile)
	rg.GET("/permissions/profile", v1.Permission().GetCommunityPermissionProfiles)

	// 新增預約設施主檔 (管理員 權限2+)
	rg.POST("/admin/facilities", v1.Facility().CreateFacility)

	// 設施預約相關
	rg.POST("/facilities/:facility_id/reservations", v1.Reservation().CreateReservation)
	rg.PATCH("/reservations/:id/cancel", v1.Reservation().CancelReservation)
	rg.POST("/reservations/:id/reschedule", v1.Reservation().Reschedule)
	rg.PATCH("/admin/reschedule/:id/approve", v1.Reservation().AdminApproveReschedule)

	// 取得平台列表
	rg.GET("/platform/getlist", v1.Platform().Platform_GetList)

	// 包裹管理
	rg.POST("/parcels", v1.Parcel().CreateParcel)
	rg.GET("/parcels", v1.Parcel().GetParcelList)
	rg.PUT("/parcels/:id/pickup", v1.Parcel().PickupParcel)

}
