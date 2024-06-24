package dto

type CreateTicketRequest struct {
	TicketAuthorities []*UserAuthorityReq `json:"ticketAuthorities"`
}

type Ticket struct {
	TicketId          string `json:"ticketId"`
	IsUsed            bool   `json:"isUsed"`
	TicketAuthorities []*TicketAuthority
}

type TicketAuthority struct {
	AuthorityId      int       `json:"-"`
	AuthorityCode    string    `json:"authorityCode"`
	AuthorityName    string    `json:"authorityName"`
	Summary          string    `json:"summary"`
	ExpiryDuration   *Duration `json:"expiryDuration"`
	ExpiryDurationMS *int64    `json:"expiryDurationMS"`
}
