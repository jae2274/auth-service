package dto

type CreateTicketRequest struct {
	TicketAuthorities []*UserAuthorityReq `json:"ticketAuthorities"`
}

type Ticket struct {
	TicketId          string `json:"ticketId"`
	TicketAuthorities []*TicketAuthority
}

type TicketAuthority struct {
	AuthorityCode    string `json:"authorityCode"`
	AuthorityName    string `json:"authorityName"`
	Summary          string `json:"summary"`
	ExpiryDurationMS *int64 `json:"expiryDurationMS"`
}
