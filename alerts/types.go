package alerts

import "sync"

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

type Alert struct {
	ID          int64    `json:"id,string"`
	Category    int      `json:"cat,string"`
	Title       string   `json:"title"`
	Cities      []string `json:"data"`
	Description string   `json:"desc"`
}

func (wa *Alert) CategoryName() string {
	if name, found := AlertTypeNames[wa.Category]; found {
		return name
	}
	return "genric"
}

var (
	lastAlertID  int64
	lastAlertMtx sync.Mutex
)

func (a *Alert) ShouldSend() bool {
	lastAlertMtx.Lock()
	defer lastAlertMtx.Unlock()

	if a.ID < lastAlertID {
		lastAlertID = a.ID
		return true
	}
	return false
}
