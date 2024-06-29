package dto

type CreateTicketRequest struct {
	TicketName        string              `json:"ticketName"`
	TicketAuthorities []*UserAuthorityReq `json:"ticketAuthorities"`
}

type Ticket struct {
	TicketId          string             `json:"ticketId"`
	TicketName        string             `json:"ticketName"`
	TicketAuthorities []*TicketAuthority `json:"ticketAuthorities"`
	CreateUnixMilli   int64              `json:"createUnixMilli"`
	UsedUnixMilli     *int64             `json:"usedUnixMilli"`
}

type TicketAuthority struct {
	AuthorityId      int    `json:"-"`
	AuthorityCode    string `json:"authorityCode"`
	AuthorityName    string `json:"authorityName"`
	Summary          string `json:"summary"`
	ExpiryDurationMS *int64 `json:"expiryDurationMS"`
}
