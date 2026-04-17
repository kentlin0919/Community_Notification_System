package permission

// UpdatePermissionProfileRequest 讓社區 Admin 更改該社區的自訂職等顯示名稱
type UpdatePermissionProfileRequest struct {
	PermissionID int    `json:"permission_id" binding:"required,min=3,max=7"`
	DisplayName  string `json:"display_name" binding:"required,max=50"`
	Description  string `json:"description"`
}
