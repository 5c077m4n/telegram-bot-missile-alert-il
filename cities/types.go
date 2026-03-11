package cities

type CityData struct {
	Label           string `json:"label"`
	Value           string `json:"value"`
	ID              int64  `json:"id,string"`
	AreaID          int    `json:"areaid"`
	AreaName        string `json:"areaname"`
	LabelHebrew     string `json:"label_he"`
	SecondsToBunker int    `json:"migun_time"`
}
