package dto

type CreateTicketRequest struct {
	TicketAuthorities []*UserAuthorityReq `json:"ticketAuthorities"`
}

type Ticket struct {
	TicketId          string             `json:"ticketId"`
	IsUsed            bool               `json:"isUsed"`
	TicketAuthorities []*TicketAuthority `json:"ticketAuthorities"`
	CreateUnixMilli   int64              `json:"createUnixMilli"`
}

type TicketAuthority struct {
	AuthorityId      int    `json:"-"`
	AuthorityCode    string `json:"authorityCode"`
	AuthorityName    string `json:"authorityName"`
	Summary          string `json:"summary"`
	ExpiryDurationMS *int64 `json:"expiryDurationMS"`
}
