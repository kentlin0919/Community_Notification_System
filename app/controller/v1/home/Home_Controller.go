package home

import (
	homeModel "Community_Notification_System/app/models/Home"
	"Community_Notification_System/app/models/model"
	homeRepo "Community_Notification_System/app/repositories/home"
	"net/http"

	"github.com/gin-gonic/gin"

	utilsErr "Community_Notification_System/utils/errors"
)

// HomeController 首頁控制器
type HomeController struct{}

// NewHomeController 建立 HomeController 實例
func NewHomeController() *HomeController {
	return &HomeController{}
}

// GetDashboard 取得首頁儀表板資訊
// @Summary 取得首頁儀表板資訊
// @Description 根據目前登入者的權限，聚合回傳對應的統計數據、預約清單與訊息。
// @Tags Home
// @Produce json
// @Success 200 {object} homeModel.HomeResponse "成功返回儀表板數據"
// @Failure 401 {object} model.Response401Error "未授權"
// @Failure 500 {object} model.Response500Error "伺服器錯誤"
// @Router /api/v1/home [get]
// @Security BearerAuth
func (h *HomeController) GetDashboard(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法取得使用者登入資訊"))
		return
	}

	permissionID, _ := ctx.Get("permission_id")
	permID, ok := permissionID.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法取得權限資訊"))
		return
	}

	communityID, _ := ctx.Get("community_id")
	commID, _ := communityID.(uint64)

	switch {
	case permID <= 2:
		// 超級管理員
		h.getSuperAdminDashboard(ctx)
	case permID <= 7:
		// 社區管理員
		h.getAdminDashboard(ctx, commID)
	default:
		// 住戶
		h.getResidentDashboard(ctx, userID, commID)
	}
}

// getResidentDashboard 住戶首頁
func (h *HomeController) getResidentDashboard(ctx *gin.Context, userID string, communityID uint64) {
	// 計數
	counts := homeModel.ResidentCounts{
		UnreadMessages:       homeRepo.CountUnreadMessages(userID),
		UnclaimedParcels:     homeRepo.CountUnclaimedParcels(communityID),
		UpcomingReservations: homeRepo.CountUpcomingReservations(userID),
	}

	// 最新 5 筆訊息
	rawMessages := homeRepo.GetRecentMessages(userID, 5)
	messages := make([]homeModel.MessageSummary, 0, len(rawMessages))
	for _, m := range rawMessages {
		messages = append(messages, homeModel.MessageSummary{
			ID:         m.ID,
			Title:      m.Title,
			Content:    m.Detail,
			Category:   m.Category,
			IsRead:     m.IsRead,
			CreateTime: m.CreateTime.Format("2006-01-02 15:04"),
		})
	}

	// 最近 3 筆預約
	rawReservations := homeRepo.GetUpcomingReservations(userID, 3)
	reservations := make([]homeModel.ReservationSummary, 0, len(rawReservations))
	for _, r := range rawReservations {
		reservations = append(reservations, homeModel.ReservationSummary{
			ID:              r.ID,
			FacilityName:    r.FacilityName,
			ReservationDate: r.ReservationDate,
			StartTime:       r.StartTime,
			EndTime:         r.EndTime,
			Status:          r.Status,
		})
	}

	ctx.JSON(http.StatusOK, homeModel.HomeResponse{
		Role: "resident",
		Data: homeModel.ResidentDashboard{
			Counts:               counts,
			RecentMessages:       messages,
			UpcomingReservations: reservations,
		},
	})
}

// getAdminDashboard 社區管理員首頁
func (h *HomeController) getAdminDashboard(ctx *gin.Context, communityID uint64) {
	stats := homeModel.AdminStats{
		ResidentCount:     homeRepo.CountCommunityResidents(communityID),
		TodayReservations: homeRepo.CountTodayReservations(communityID),
		UnhandledParcels:  homeRepo.CountUnclaimedParcels(communityID),
	}

	// 今日預約列表 (最多 5 筆)
	rawReservations := homeRepo.GetTodayReservations(communityID, 5)
	reservations := make([]homeModel.ReservationSummary, 0, len(rawReservations))
	for _, r := range rawReservations {
		reservations = append(reservations, homeModel.ReservationSummary{
			ID:              r.ID,
			FacilityName:    r.FacilityName,
			ReservationDate: r.ReservationDate,
			StartTime:       r.StartTime,
			EndTime:         r.EndTime,
			Status:          r.Status,
		})
	}

	// 最新訊息列表 (最多 5 筆)
	rawMessages := homeRepo.GetCommunityRecentMessages(communityID, 5)
	messages := make([]homeModel.MessageSummary, 0, len(rawMessages))
	for _, m := range rawMessages {
		messages = append(messages, homeModel.MessageSummary{
			ID:         m.ID,
			Title:      m.Title,
			Content:    m.Detail,
			Category:   m.Category,
			IsRead:     m.IsRead,
			CreateTime: m.CreateTime.Format("2006-01-02 15:04"),
		})
	}

	ctx.JSON(http.StatusOK, homeModel.HomeResponse{
		Role: "community_admin",
		Data: homeModel.AdminDashboard{
			Stats:             stats,
			TodayReservations: reservations,
			RecentMessages:    messages,
		},
	})
}

// getSuperAdminDashboard 超級管理員首頁
func (h *HomeController) getSuperAdminDashboard(ctx *gin.Context) {
	platformStats := homeModel.PlatformStats{
		CommunityCount:     homeRepo.CountAllCommunities(),
		UserCount:          homeRepo.CountAllUsers(),
		PendingCommunities: homeRepo.CountPendingCommunityApplications(),
	}

	rawApps := homeRepo.GetPendingCommunityApplications(10)
	pendingList := make([]homeModel.PendingCommunitySummary, 0, len(rawApps))
	for _, app := range rawApps {
		pendingList = append(pendingList, homeModel.PendingCommunitySummary{
			ID:            app.ID,
			CommunityName: app.CommunityName,
			ApplicantName: app.ApplicantName,
			Status:        app.Status,
			CreatedAt:     app.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	ctx.JSON(http.StatusOK, homeModel.HomeResponse{
		Role: "super_admin",
		Data: homeModel.SuperAdminDashboard{
			PlatformStats:      platformStats,
			PendingCommunities: pendingList,
		},
	})
}
