package communitydb

type CommunityInfo struct {
	Community_id   uint64 `json:"community_id" gorm:"primarykey"`
	PostalCode     int    `json:"postal_code"`
	Municipality   string `json:"municipality"`
	District       string `json:"district"`
	RoadName       string `json:"road_name"`
	LaneNumber     int    `json:"lane_number"`
	AlleyNumber    int    `json:"alley_number"`
	Community_name string `json:"community_name"`
	Address        string `json:"address"`
}
