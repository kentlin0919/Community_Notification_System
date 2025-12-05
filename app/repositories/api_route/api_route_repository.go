package api_route

import (
	"Community_Notification_System/database"
	"Community_Notification_System/database/ApiRoute_DB"
)

type ApiRouteRepository struct{}

func NewApiRouteRepository() *ApiRouteRepository {
	return &ApiRouteRepository{}
}

// GetRequiredPermission fetches the required permission ID for a given API route.
func (r *ApiRouteRepository) GetRequiredPermission(path string, method string) (uint, error) {
	var route ApiRoute_DB.ApiRoute
	err := database.DB.Where("path = ? AND method = ?", path, method).First(&route).Error
	if err != nil {
		return 0, err // Return 0 and the error if route not found or other db error
	}
	return uint(route.RequiredPermissionID), nil
}
