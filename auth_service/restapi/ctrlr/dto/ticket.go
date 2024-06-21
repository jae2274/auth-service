package dto

type GetTicketInfoResponse struct {
	TicketId          string `json:"ticketId"`
	TicketAuthorities []*TicketAuthorityReq
}

type TicketAuthorityReq struct {
	AuthorityCode    string `json:"authorityCode"`
	AuthorityName    string `json:"authorityName"`
	Summary          string `json:"summary"`
	ExpiryDurationMS *int64 `json:"expiryDurationMS"`
}
