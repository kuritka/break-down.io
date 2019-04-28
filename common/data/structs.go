package data

type Calendar struct {
	Name     string `json:"name" bson:"_id"`
	Schedule Schedule
}

type Schedule struct {
	Mon Mon
	Tue Tue
	Wed Wed
	Thr Thr
	Fri Fri
}

type Event struct {
	From    string
	To      string
	Note    string
	Subject string
	Week    string
}

type Mon struct {
	Events []Event
}

type Tue struct {
	Events []Event
}

type Wed struct {
	Events []Event
}

type Thr struct {
	Events []Event
}

type Fri struct {
	Events []Event
}
