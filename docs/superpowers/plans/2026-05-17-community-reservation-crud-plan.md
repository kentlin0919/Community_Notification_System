# Community Reservation Facility CRUD Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement CRUD operations for Community Facilities and Reservations, including role-based visibility, soft deletes, and booking overlap prevention.

**Architecture:** We are updating the existing Clean Architecture layers. Database schemas (GORM) will be updated for soft deletes. Repositories will handle logic like overlap queries and soft deletions. Controllers will apply JWT-based role filtering and route requests.

**Tech Stack:** Go, Gin, GORM, PostgreSQL, Testify (for testing)

---

## File Structure Map
- Modify: `database/Facility_DB/Facility_Schema.go` (Add DeletedAt)
- Modify: `database/Facility_DB/Reservation_Schema.go` (Add DeletedAt)
- Modify: `app/models/facility/Facility_Model.go` (Add List/Update structs)
- Modify: `app/models/reservation/Reservation_Model.go` (Add List structs)
- Modify: `app/repositories/facility/Facility_Repository.go` (Add GetList, GetDetail, Update, Delete)
- Modify: `app/repositories/reservation/Reservation_Repository.go` (Add overlap check, GetList, Delete)
- Modify: `app/controller/v1/facility/Facility_Controller.go` (Add HTTP handlers)
- Modify: `app/controller/v1/reservation/Reservation_Controller.go` (Add HTTP handlers)
- Modify: `routers/api/v1/v1.go` (Add new routes)
- Create: `app/repositories/facility/Facility_Repository_test.go`
- Create: `app/repositories/reservation/Reservation_Repository_test.go`

---

### Task 1: Update Database Schemas for Soft Delete

**Files:**
- Modify: `database/Facility_DB/Facility_Schema.go`
- Modify: `database/Facility_DB/Reservation_Schema.go`

- [ ] **Step 1: Add DeletedAt to FacilityInfo**

```go
package facility_db

import (
	"time"
	"gorm.io/gorm"
)

// FacilityInfo 可供住戶預約的公共設施主檔
type FacilityInfo struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	CommunityID  uint64         `gorm:"index;not null;uniqueIndex:idx_community_facility_name" json:"community_id"`
	Name         string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_community_facility_name" json:"name"`
	FacilityType string         `gorm:"type:varchar(50);not null" json:"facility_type"` 
	Location     string         `gorm:"type:varchar(255)" json:"location"`
	Description  string         `gorm:"type:text" json:"description"`
	Status       string         `gorm:"type:varchar(20);not null;default:'draft'" json:"status"` 
	CoverImage   string         `gorm:"type:varchar(255)" json:"cover_image"`
	CreatedBy    string         `gorm:"type:varchar(100)" json:"created_by"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index;uniqueIndex:idx_community_facility_name" json:"-"`
}
```
*Note: We included DeletedAt in the unique index so you can create a facility with the same name after deleting the old one.*

- [ ] **Step 2: Add DeletedAt to FacilityReservation**

```go
package facility_db

import (
	"time"
	"gorm.io/gorm"
)

// FacilityReservation 住戶預約主單
type FacilityReservation struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	FacilityID      uint64         `gorm:"index;not null" json:"facility_id"`
	UserID          string         `gorm:"index;type:varchar(100);not null" json:"user_id"`
	CommunityID     uint64         `gorm:"index;not null" json:"community_id"`
	HomeID          uint64         `gorm:"index" json:"home_id"`
	ReservationDate string         `gorm:"type:date;not null" json:"reservation_date"` 
	StartTime       string         `gorm:"type:varchar(5);not null" json:"start_time"`   
	EndTime         string         `gorm:"type:varchar(5);not null" json:"end_time"`     
	PeopleCount     int            `gorm:"not null" json:"people_count"`
	Status          string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` 
	Remark          string         `gorm:"type:text" json:"remark"`
	CancelReason    string         `gorm:"type:text" json:"cancel_reason"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
```
*(Leave RescheduleRequest as is for now as it's just a request log).*

- [ ] **Step 3: Commit**

```bash
git add database/Facility_DB/Facility_Schema.go database/Facility_DB/Reservation_Schema.go
git commit -m "feat(db): add soft delete support to facility and reservation schemas"
```

---

### Task 2: Implement Facility Repository CRUD (List, Update, Delete)

**Files:**
- Create: `app/repositories/facility/Facility_Repository_test.go`
- Modify: `app/repositories/facility/Facility_Repository.go`

- [ ] **Step 1: Write Repository Interface/Methods**

Add to `app/repositories/facility/Facility_Repository.go`:

```go
package facility

// ... existing code ...

func GetFacilityListRepository(communityID uint64, onlyActive bool) repositoryModels.RepositoryModel[[]facilitydb.FacilityInfo] {
	var result repositoryModels.RepositoryModel[[]facilitydb.FacilityInfo]
	var facilities []facilitydb.FacilityInfo

	query := database.DB.Where("community_id = ?", communityID)
	if onlyActive {
		query = query.Where("status = ?", "active")
	}

	err := query.Find(&facilities).Error
	if err != nil {
		result.Statue.Error = err
	} else {
		result.Result = facilities
	}
	return result
}

func GetFacilityDetailRepository(id uint64, communityID uint64) repositoryModels.RepositoryModel[facilitydb.FacilityInfo] {
	var result repositoryModels.RepositoryModel[facilitydb.FacilityInfo]
	var facility facilitydb.FacilityInfo

	err := database.DB.Where("id = ? AND community_id = ?", id, communityID).First(&facility).Error
	if err != nil {
		result.Statue.Error = err
	} else {
		result.Result = facility
	}
	return result
}

func UpdateFacilityRepository(facility *facilitydb.FacilityInfo) repositoryModels.RepositoryModel[facilitydb.FacilityInfo] {
	var result repositoryModels.RepositoryModel[facilitydb.FacilityInfo]
	
	// Check name conflict excluding current ID
	var exist facilitydb.FacilityInfo
	if err := database.DB.Where("community_id = ? AND name = ? AND id != ?", facility.CommunityID, facility.Name, facility.ID).First(&exist).Error; err == nil {
		result.Statue.Error = errors.New("該社區已存在相同名稱的設施")
		return result
	}

	err := database.DB.Save(facility).Error
	if err != nil {
		result.Statue.Error = err
	} else {
		result.Result = *facility
	}
	return result
}

func DeleteFacilityRepository(id uint64, communityID uint64) error {
	return database.DB.Where("id = ? AND community_id = ?", id, communityID).Delete(&facilitydb.FacilityInfo{}).Error
}
```

- [ ] **Step 2: Commit**

```bash
git add app/repositories/facility/Facility_Repository.go
git commit -m "feat(repo): implement facility get list, update, and delete"
```

---

### Task 3: Implement Facility Controller API

**Files:**
- Modify: `app/models/facility/Facility_Model.go`
- Modify: `app/controller/v1/facility/Facility_Controller.go`
- Modify: `routers/api/v1/v1.go`

- [ ] **Step 1: Add Update Request Model**

Add to `app/models/facility/Facility_Model.go`:

```go
type UpdateFacilityRequest struct {
	Name         string `json:"name" binding:"required,max=100" example:"健身房"`
	FacilityType string `json:"facility_type" binding:"required,max=50" example:"運動設施"`
	Location     string `json:"location" example:"B1 棟"`
	Description  string `json:"description" example:"提供跑步機、啞鈴等健身設備。"`
	Status       string `json:"status" binding:"omitempty,oneof=draft active inactive maintenance" example:"active"`
	CoverImage   string `json:"cover_image" example:"https://example.com/facility.jpg"`
}
```

- [ ] **Step 2: Add Controller Methods**

Add to `app/controller/v1/facility/Facility_Controller.go`:

```go
import "strconv"

// GetFacilityList 取得設施列表
// ... Swagger annotations ...
func (c *FacilityController) GetFacilityList(ctx *gin.Context) {
	communityID, _ := ctx.Get("community_id")
	permissionLevel, _ := ctx.Get("permission_level")

	// If permission is resident (<2), only show active
	onlyActive := false
	if level, ok := permissionLevel.(int); ok && level < 2 {
		onlyActive = true
	} else if levelStr, ok := permissionLevel.(string); ok {
        levelInt, _ := strconv.Atoi(levelStr)
        if levelInt < 2 {
            onlyActive = true
        }
    }

	res := repository.GetFacilityListRepository(communityID.(uint64), onlyActive)
	if res.Statue.Error != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "獲取列表失敗"))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res.Result})
}

// GetFacilityDetail 取得單一設施
func (c *FacilityController) GetFacilityDetail(ctx *gin.Context) {
	communityID, _ := ctx.Get("community_id")
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	res := repository.GetFacilityDetailRepository(id, communityID.(uint64))
	if res.Statue.Error != nil {
		ctx.JSON(http.StatusNotFound, model.NewErrorRequest(http.StatusNotFound, "找不到該設施"))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res.Result})
}

// UpdateFacility 更新設施
func (c *FacilityController) UpdateFacility(ctx *gin.Context) {
	communityID, _ := ctx.Get("community_id")
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	var req facilityModel.UpdateFacilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "Invalid input"))
		return
	}

	// Fetch existing
	res := repository.GetFacilityDetailRepository(id, communityID.(uint64))
	if res.Statue.Error != nil {
		ctx.JSON(http.StatusNotFound, model.NewErrorRequest(http.StatusNotFound, "找不到該設施"))
		return
	}

	facility := res.Result
	facility.Name = req.Name
	facility.FacilityType = req.FacilityType
	facility.Location = req.Location
	facility.Description = req.Description
	if req.Status != "" {
		facility.Status = req.Status
	}
	facility.CoverImage = req.CoverImage

	updateRes := repository.UpdateFacilityRepository(&facility)
	if updateRes.Statue.Error != nil {
		ctx.JSON(http.StatusConflict, model.NewErrorRequest(http.StatusConflict, updateRes.Statue.Error.Error()))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": updateRes.Result})
}

// DeleteFacility 刪除設施
func (c *FacilityController) DeleteFacility(ctx *gin.Context) {
	communityID, _ := ctx.Get("community_id")
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	err := repository.DeleteFacilityRepository(id, communityID.(uint64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "刪除失敗"))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "刪除成功"})
}
```

- [ ] **Step 3: Register Routes**

In `routers/api/v1/v1.go`, under `V1PrivateRoutes`:

```go
	// 設施管理
	rg.GET("/facilities", v1.Facility().GetFacilityList)
	rg.GET("/facilities/:id", v1.Facility().GetFacilityDetail)
	rg.POST("/admin/facilities", v1.Facility().CreateFacility)
	rg.PUT("/admin/facilities/:id", v1.Facility().UpdateFacility)
	rg.DELETE("/admin/facilities/:id", v1.Facility().DeleteFacility)
```

- [ ] **Step 4: Commit**

```bash
git add app/models/facility/Facility_Model.go app/controller/v1/facility/Facility_Controller.go routers/api/v1/v1.go
git commit -m "feat(api): implement facility crud endpoints"
```

---

### Task 4: Add Overlap Check to Reservation Creation

**Files:**
- Modify: `app/repositories/reservation/Reservation_Repository.go`

- [ ] **Step 1: Modify CreateReservationRepository**

In `app/repositories/reservation/Reservation_Repository.go`, update `CreateReservationRepository` to add the overlap check before creating:

```go
import "errors"

// In CreateReservationRepository:
func CreateReservationRepository(reservation *facilitydb.FacilityReservation, rule *facilitydb.FacilityRule) repositoryModels.RepositoryModel[facilitydb.FacilityReservation] {
	var result repositoryModels.RepositoryModel[facilitydb.FacilityReservation]

	// 1. Check Overlap
	var count int64
	err := database.DB.Model(&facilitydb.FacilityReservation{}).
		Where("facility_id = ? AND reservation_date = ? AND status NOT IN ?", reservation.FacilityID, reservation.ReservationDate, []string{"cancelled", "rejected"}).
		Where("(start_time < ? AND end_time > ?)", reservation.EndTime, reservation.StartTime).
		Count(&count).Error
		
	if err != nil {
		result.Statue.Error = errors.New("檢查時段衝突失敗")
		return result
	}
	if count > 0 {
		result.Statue.Error = errors.New("該時段已被預約")
		return result
	}

	// 2. Existing Creation Logic
	if err := database.DB.Create(reservation).Error; err != nil {
		result.Statue.Error = err
	} else {
		result.Result = *reservation
		result.Statue.RowsAffected = 1
	}

	return result
}
```

- [ ] **Step 2: Commit**

```bash
git add app/repositories/reservation/Reservation_Repository.go
git commit -m "feat(repo): add time overlap prevention for facility reservations"
```

---

### Task 5: Implement Reservation List and Admin Delete

**Files:**
- Modify: `app/repositories/reservation/Reservation_Repository.go`
- Modify: `app/controller/v1/reservation/Reservation_Controller.go`
- Modify: `routers/api/v1/v1.go`

- [ ] **Step 1: Add Repo Methods**

In `app/repositories/reservation/Reservation_Repository.go`:

```go
func GetReservationListRepository(communityID uint64, userID string, filterByUserID bool) repositoryModels.RepositoryModel[[]facilitydb.FacilityReservation] {
	var result repositoryModels.RepositoryModel[[]facilitydb.FacilityReservation]
	var reservations []facilitydb.FacilityReservation

	query := database.DB.Where("community_id = ?", communityID)
	if filterByUserID {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Order("reservation_date DESC, start_time DESC").Find(&reservations).Error
	if err != nil {
		result.Statue.Error = err
	} else {
		result.Result = reservations
	}
	return result
}

func DeleteReservationRepository(id uint64, communityID uint64) error {
	return database.DB.Where("id = ? AND community_id = ?", id, communityID).Delete(&facilitydb.FacilityReservation{}).Error
}
```

- [ ] **Step 2: Add Controller Methods**

In `app/controller/v1/reservation/Reservation_Controller.go`:

```go
// GetReservationList 取得預約列表
func (c *ReservationController) GetReservationList(ctx *gin.Context) {
	communityID, _ := ctx.Get("community_id")
	userID, _ := ctx.Get("user_id")
	permissionLevel, _ := ctx.Get("permission_level")

	filterByUserID := false
	if level, ok := permissionLevel.(int); ok && level < 2 {
		filterByUserID = true
	} else if levelStr, ok := permissionLevel.(string); ok {
        levelInt, _ := strconv.Atoi(levelStr)
        if levelInt < 2 {
            filterByUserID = true
        }
    }

	res := repository.GetReservationListRepository(communityID.(uint64), userID.(string), filterByUserID)
	if res.Statue.Error != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "獲取列表失敗"))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res.Result})
}

// AdminDeleteReservation 管理員刪除預約
func (c *ReservationController) AdminDeleteReservation(ctx *gin.Context) {
	communityID, _ := ctx.Get("community_id")
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	err := repository.DeleteReservationRepository(id, communityID.(uint64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "刪除失敗"))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "刪除成功"})
}
```

- [ ] **Step 3: Register Routes**

In `routers/api/v1/v1.go`, under `V1PrivateRoutes`:

```go
	// 設施預約相關
	rg.GET("/reservations", v1.Reservation().GetReservationList)
	rg.DELETE("/admin/reservations/:id", v1.Reservation().AdminDeleteReservation)
```

- [ ] **Step 4: Commit**

```bash
git add app/repositories/reservation/Reservation_Repository.go app/controller/v1/reservation/Reservation_Controller.go routers/api/v1/v1.go
git commit -m "feat(api): add reservation list and admin delete endpoints"
```
