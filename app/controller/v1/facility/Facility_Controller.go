package facility

import (
	facilityModel "Community_Notification_System/app/models/facility"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/facility"
	facilitydb "Community_Notification_System/database/Facility_DB"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FacilityController struct{}

func NewFacilityController() *FacilityController {
	return &FacilityController{}
}

// CreateFacility 新增預約設施
// @Summary 新增預約設施
// @Description 社區管理員新增公用設施與預設預約條件
// @Tags Facility Management
// @Accept json
// @Produce json
// @Param body body facilityModel.CreateFacilityRequest true "設施資料"
// @Success 201 {object} map[string]interface{} "設施新增成功"
// @Failure 400 {object} model.ErrorRequest "請求參數錯誤"
// @Failure 401 {object} model.ErrorRequest "無法取得登入資訊"
// @Failure 403 {object} model.ErrorRequest "權限不足"
// @Failure 409 {object} model.ErrorRequest "該社區已存在相同名稱的設施"
// @Security BearerAuth
// @Router /api/v1/admin/facilities [post]
func (c *FacilityController) CreateFacility(ctx *gin.Context) {
	communityID, exists := ctx.Get("community_id")
	userID, userExists := ctx.Get("user_id")

	if !exists || !userExists {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorRequest(http.StatusUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}

	var req facilityModel.CreateFacilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "Invalid input: "+err.Error()))
		return
	}

	cID := communityID.(uint64)

	// API 防呆：只能替自己的社區建設施 (這個防線在 CommunityContextMiddleware 也有一道)
	if req.CommunityID != cID {
		ctx.JSON(http.StatusForbidden, model.NewErrorRequest(http.StatusForbidden, "你只能管理所屬社區之設施"))
		return
	}

	status := "draft"
	if req.Status != "" {
		status = req.Status
	}

	newFacility := &facilitydb.FacilityInfo{
		CommunityID:  cID,
		Name:         req.Name,
		FacilityType: req.FacilityType,
		Location:     req.Location,
		Description:  req.Description,
		Status:       status,
		CoverImage:   req.CoverImage,
		CreatedBy:    userID.(string),
	}

	// 呼叫 Repository 新增
	res := repository.CreateFacilityWithRuleTransactionRepository(newFacility)
	if res.Statue.Error != nil {
		// 判斷是否為重複命名
		if res.Statue.Error.Error() == "該社區已存在相同名稱的設施" {
			ctx.JSON(http.StatusConflict, model.NewErrorRequest(http.StatusConflict, res.Statue.Error.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "建立設施失敗"))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "設施新增成功",
		"data": gin.H{
			"id":            res.Result.ID,
			"community_id":  res.Result.CommunityID,
			"name":          res.Result.Name,
			"facility_type": res.Result.FacilityType,
			"status":        res.Result.Status,
		},
	})
}
