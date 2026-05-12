package communityManager

import (
	communityModel "Community_Notification_System/app/models/community"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/community"
	communitydb "Community_Notification_System/database/Community_DB"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CommunityManager_Register 社區送出申請
// @Summary 社區送出申請
// @Description 社區申請人填寫基本資料送出申請，狀態預設為 pending。
// @Tags CommunityManager
// @Accept json
// @Produce json
// @Param application body communityModel.RegisterApplicationRequest true "社區申請資料"
// @Success 200 {object} communityModel.RegisterApplicationResponse "社區申請已送出"
// @Failure 400 {object} model.Response400Error "請求參數錯誤"
// @Failure 500 {object} model.Response500Error "系統錯誤"
// @Router /api/v1/community/register [post]
func (c *CommunityManagerController) CommunityManager_Register(ctx *gin.Context) {
	var req communityModel.RegisterApplicationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorModel := model.NewErrorRequest(http.StatusBadRequest, "Invalid input")
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	if utf8.RuneCountInString(req.AdminPassword) < 6 {
		errorModel := model.NewErrorRequest(http.StatusBadRequest, "管理員密碼長度至少需 6 碼")
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		errorModel := model.NewErrorRequest(http.StatusInternalServerError, "密碼處理失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	application := &communitydb.CommunityRegisterApplication{
		Status:            "pending",
		PostalCode:        req.PostalCode,
		Municipality:      req.Municipality,
		District:          req.District,
		RoadName:          req.RoadName,
		LaneNumber:        req.LaneNumber,
		AlleyNumber:       req.AlleyNumber,
		CommunityName:     req.CommunityName,
		Address:           req.Address,
		AdminName:         req.AdminName,
		AdminEmail:        req.AdminEmail,
		AdminPasswordHash: string(hashedPassword),
		AdminPhone:        req.AdminPhone,
		ApplicantName:     req.ApplicantName,
		ApplicantEmail:    req.ApplicantEmail,
		ApplicantPhone:    req.ApplicantPhone,
		Remark:            req.Remark,
	}

	result := repository.CreateApplicationRepository(application)
	if result.Statue.Error != nil {
		errorModel := model.NewErrorRequest(http.StatusInternalServerError, "建立社區申請單失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "社區申請已送出",
		"data": gin.H{
			"application_id": application.ID,
			"status":         application.Status,
		},
	})
}

