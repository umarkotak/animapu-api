package contract

type (
	UserActivityData struct {
		Users []UserActivity `json:"users"`
	}

	UserActivity struct {
		VisitorID string    `json:"visitor_id"`
		Email     string    `json:"email"`
		Histories []History `json:"histories"`
	}
)
