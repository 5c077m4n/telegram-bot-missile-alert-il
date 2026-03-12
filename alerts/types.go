package alerts

var AlertTypeNames = map[int]string{
	1:  "missiles",
	2:  "radiologicalEvent",
	3:  "earthQuake",
	4:  "tsunami",
	5:  "hostileAircraftIntrusion",
	6:  "hazardousMaterials",
	7:  "terroristInfiltration",
	8:  "missilesDrill",
	9:  "earthQuakeDrill",
	10: "radiologicalEventDrill",
	11: "tsunamiDrill",
	12: "hostileAircraftIntrusionDrill",
	13: "hazardousMaterialsDrill",
	14: "terroristInfiltrationDrill",
	20: "newsFlash",
	99: "unknown",
}

type WarningAlert struct {
	ID          int64    `json:"id,string"`
	Category    int      `json:"cat,string"`
	Title       string   `json:"title"`
	Cities      []string `json:"data"`
	Description string   `json:"desc"`
}

func (wa *WarningAlert) CategoryName() string {
	if name, found := AlertTypeNames[wa.Category]; found {
		return name
	}
	return "genric"
}
