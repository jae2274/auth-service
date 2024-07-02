package dto

type CreateTicketRequest struct {
	UseableCount      int                 `json:"useableCount"`
	TicketName        string              `json:"ticketName"`
	TicketAuthorities []*UserAuthorityReq `json:"ticketAuthorities"`
}

type Ticket struct {
	TicketId          string             `json:"ticketId"`
	TicketName        string             `json:"ticketName"`
	TicketAuthorities []*TicketAuthority `json:"ticketAuthorities"`
	CreateUnixMilli   int64              `json:"createUnixMilli"`
	UseableCount      int                `json:"useableCount"`
	UsedCount         int                `json:"usedCount"`
}

type TicketDetail struct {
	Ticket
	UsedInfos []*UsedInfo `json:"usedInfos"`
	CreatedBy int         `json:"createdBy"`
}

type UsedInfo struct {
	UsedBy        int    `json:"usedBy"`
	UsedUserName  string `json:"usedUserName"`
	UsedUnixMilli int64  `json:"usedUnixMilli"`
}

type TicketAuthority struct {
	AuthorityId      int    `json:"-"`
	AuthorityCode    string `json:"authorityCode"`
	AuthorityName    string `json:"authorityName"`
	Summary          string `json:"summary"`
	ExpiryDurationMS *int64 `json:"expiryDurationMS"`
}
