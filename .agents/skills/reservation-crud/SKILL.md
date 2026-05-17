# Community Reservation CRUD Skill

Use this skill when implementing or extending the Facility Reservation functionality in the Community Notification System.

## 1. Principles
- **Clean Architecture**: Strictly separate Controller, Repository, and Database layers.
- **Role-Based Visibility**: Automatically filter data based on user roles (Resident vs Admin).
- **Integrity**: Use Soft Delete for all master data (Facilities) to preserve audit trails.
- **Safety**: Always perform overlap checks and rule validation before committing a reservation.

## 2. Implementation Patterns

### Soft Delete Implementation
In `database/Facility_DB/Facility_Schema.go`:
```go
type FacilityInfo struct {
    gorm.Model // Includes ID, CreatedAt, UpdatedAt, DeletedAt
    // ... fields
}
```
Or manually adding `DeletedAt`:
```go
DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
```

### Overlap Check Query (GORM)
```go
var count int64
db.Model(&facility_db.FacilityReservation{}).
    Where("facility_id = ? AND reservation_date = ? AND status NOT IN ?", facilityID, date, []string{"cancelled", "rejected"}).
    Where("(start_time < ? AND end_time > ?)", newEndTime, newStartTime).
    Count(&count)
```

### Role-Based List Filter
In Controller:
```go
// Get User context
userID, _ := ctx.Get("user_id")
communityID, _ := ctx.Get("community_id")
permissionLevel, _ := ctx.Get("permission_level") // Assuming middleware sets this

// Apply filtering logic
if permissionLevel.(int) < 2 { // Resident
    // filter by userID
} else {
    // filter by communityID only
}
```

## 3. Checklist
- [ ] Update `FacilityInfo` schema to include `gorm.DeletedAt`.
- [ ] Implement `GetFacilityList` with role-based filtering (`active` status for residents).
- [ ] Implement `UpdateFacility` and `DeleteFacility` (soft delete).
- [ ] Implement `GetReservationList` with role-based filtering (own records for residents).
- [ ] Add overlap check logic in `CreateReservationRepository`.
- [ ] Ensure all new routes are registered in `routers/api/v1/v1.go` under `PrivateRoutes`.
- [ ] Generate Swagger documentation using `swag init`.
