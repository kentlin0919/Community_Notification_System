package parcel

import (
	parcelModel "Community_Notification_System/app/models/parcel"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/parcel"
	parcel_db "Community_Notification_System/database/Parcel_DB"
	"Community_Notification_System/pkg/firebase"
	"context"
	"fmt"
	"log"
	"net/http"

	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
)

// CreateParcel 管理員代收登錄包裹
// @Summary 代收登錄包裹
// @Description 管理員登錄代收包裹資訊，系統自動發送 FCM 推播通知給該住戶所有成員
// @Tags Parcel Management
// @Accept json
// @Produce json
// @Param body body parcelModel.CreateParcelRequest true "包裹資料"
// @Success 201 {object} map[string]interface{} "包裹登錄成功"
// @Failure 400 {object} model.Response400Error "請求參數錯誤"
// @Failure 401 {object} model.Response401Error "無法取得登入資訊"
// @Failure 409 {object} model.ErrorRequest "重複單號"
// @Failure 500 {object} model.Response500Error "伺服器錯誤"
// @Security BearerAuth
// @Router /api/v1/parcels [post]
func (c *ParcelController) CreateParcel(ctx *gin.Context) {
	communityID, exists := ctx.Get("community_id")
	userID, userExists := ctx.Get("user_id")

	if !exists || !userExists {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorRequest(http.StatusUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}

	var req parcelModel.CreateParcelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "無效的輸入資料: "+err.Error()))
		return
	}

	cID := communityID.(uint64)

	newParcel := &parcel_db.ParcelInfo{
		CommunityID:    cID,
		HomeID:         req.HomeID,
		CourierCompany: req.CourierCompany,
		TrackingNumber: req.TrackingNumber,
		Remark:         req.Remark,
		Status:         1, // 待領取
		CreatedBy:      userID.(string),
	}

	// 呼叫 Repository 建立紀錄
	res := repository.CreateParcelRepository(newParcel)
	if res.Statue.Error != nil {
		if res.Statue.Error.Error() == "該社區已存在相同單號的待領包裹" {
			ctx.JSON(http.StatusConflict, model.NewErrorRequest(http.StatusConflict, res.Statue.Error.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "建立包裹紀錄失敗"))
		return
	}

	// 查詢該住戶的 FCM Token 並推播通知
	tokens := repository.GetFcmTokensByHomeID(req.HomeID)
	if len(tokens) > 0 && firebase.FcmClient != nil {
		for _, token := range tokens {
			msg := &messaging.Message{
				Notification: &messaging.Notification{
					Title: "📦 包裹到府通知",
					Body:  fmt.Sprintf("您有一件來自「%s」的包裹已由管理室代收，請至管理室領取。", req.CourierCompany),
				},
				Token: token,
			}
			response, err := firebase.FcmClient.Send(context.Background(), msg)
			if err != nil {
				log.Printf("FCM 推播失敗 (token: %s): %v\n", token, err)
			} else {
				log.Printf("FCM 推播成功: %s\n", response)
			}
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "包裹登錄成功並已發送通知",
		"data": gin.H{
			"id":              res.Result.ID,
			"community_id":    res.Result.CommunityID,
			"home_id":         res.Result.HomeID,
			"courier_company": res.Result.CourierCompany,
			"tracking_number": res.Result.TrackingNumber,
			"status":          res.Result.Status,
			"received_at":     res.Result.ReceivedAt,
		},
	})
}
