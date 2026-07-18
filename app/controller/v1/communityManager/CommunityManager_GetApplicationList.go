package communityManager

import (
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/community"
	utilsErr "Community_Notification_System/utils/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CommunityManager_GetApplicationList 取得社區申請列表
// @Summary 取得社區申請列表
// @Description 僅允許 Super admin 操作，可依狀態篩選。
// @Tags CommunityManager
// @Produce json
// @Param status query string false "Status filter (pending, approved, rejected)"
// @Success 200 {object} map[string]interface{} "成功返回申請列表"
// @Failure 401 {object} model.Response401Error "未授權"
// @Failure 500 {object} model.Response500Error "伺服器錯誤"
// @Security BearerAuth
// @Router /api/v1/super-admin/community-applications [get]
func (c *CommunityManagerController) CommunityManager_GetApplicationList(ctx *gin.Context) {
	status := ctx.Query("status")
	res := repository.GetApplicationsRepository(status)
	if res.Statue.Error != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "無法取得申請列表: "+res.Statue.Error.Error()))
		return
	}

	type ApplicationResponse struct {
		ID            uint64    `json:"id"`
		CommunityName string    `json:"community_name"`
		ApplicantName string    `json:"applicant_name"`
		AdminEmail    string    `json:"admin_email"`
		Address       string    `json:"address"`
		BuilderName   string    `json:"builder_name"`
		Status        string    `json:"status"`
		CreatedAt     time.Time `json:"created_at"`
		Phone         string    `json:"phone"`
		RejectReason  string    `json:"reject_reason"`
	}

	list := make([]ApplicationResponse, 0, len(res.Result))
	for _, app := range res.Result {
		phone := app.AdminPhone
		if phone == "" {
			phone = app.ApplicantPhone
		}
		list = append(list, ApplicationResponse{
			ID:            app.ID,
			CommunityName: app.CommunityName,
			ApplicantName: app.ApplicantName,
			AdminEmail:    app.AdminEmail,
			Address:       app.Address,
			BuilderName:   "", // 後端無此欄位，給空
			Status:        app.Status,
			CreatedAt:     app.CreatedAt,
			Phone:         phone,
			RejectReason:  app.RejectReason,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}
