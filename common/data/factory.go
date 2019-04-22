package data


func Get() Calendar {
	return  Calendar{
		Name: "myCalendar",
		Schedule: Schedule {
			Mon {
				Events: []Event{
					{
						From: "11:15",
						To: "12:30",
						Note: "blah blah",
						Subject: "TRX",
						Week: "odd",
					},
					{
						From: "15:30",
						To: "17:20",
						Note: "Something else",
						Subject: "WWW",
						Week: "odd",
					},
				}},
			Tue{
				Events: []Event{
					{
						From:    "11:30",
						To:      "12:30",
						Note:    "blah blah",
						Subject: "TRX",
						Week:    "odd",
					},
				}},
			Wed{},
			Thr{},
			Fri{
				Events: []Event {
					{
						From: "11:30",
						To: "12:30",
						Note: "blah blah",
						Subject: "FitBox",
						Week: "even",
					},
					{
						From: "18:00",
						To: "19:00",
						Note: "blah blah",
						Subject: "FitBox",
						Week: "once",
					},
				},
			},
		},
	}
}
