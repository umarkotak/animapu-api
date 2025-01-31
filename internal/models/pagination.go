package models

type (
	Pagination struct {
		Limit    int64 `json:"limit" db:"limit"`
		Page     int64 `json:"page" db:"-"`
		Offset   int64 `json:"offset" db:"offset"`
		Total    int64 `json:"total" db:"-"`
		NextPage bool  `json:"next_page" db:"-"`
	}
)

func (p *Pagination) SetDefault(defaultLimit int64) {
	if p.Limit <= 0 {
		p.Limit = defaultLimit
	}

	if p.Page <= 0 {
		p.Page = 1
	}

	p.Offset = (p.Page - 1) * p.Limit
}

func (p *Pagination) SetIsNextPage() {
	offset := (p.Page - 1) * p.Limit
	if offset+p.Limit < p.Total {
		p.NextPage = true
	} else {
		p.NextPage = false
	}
}
