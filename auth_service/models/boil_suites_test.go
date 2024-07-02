// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Agreements", testAgreements)
	t.Run("AuthServers", testAuthServers)
	t.Run("Authorities", testAuthorities)
	t.Run("Tickets", testTickets)
	t.Run("TicketAuthorities", testTicketAuthorities)
	t.Run("TicketSubs", testTicketSubs)
	t.Run("Users", testUsers)
	t.Run("UserAgreements", testUserAgreements)
	t.Run("UserAuthorities", testUserAuthorities)
}

func TestDelete(t *testing.T) {
	t.Run("Agreements", testAgreementsDelete)
	t.Run("AuthServers", testAuthServersDelete)
	t.Run("Authorities", testAuthoritiesDelete)
	t.Run("Tickets", testTicketsDelete)
	t.Run("TicketAuthorities", testTicketAuthoritiesDelete)
	t.Run("TicketSubs", testTicketSubsDelete)
	t.Run("Users", testUsersDelete)
	t.Run("UserAgreements", testUserAgreementsDelete)
	t.Run("UserAuthorities", testUserAuthoritiesDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Agreements", testAgreementsQueryDeleteAll)
	t.Run("AuthServers", testAuthServersQueryDeleteAll)
	t.Run("Authorities", testAuthoritiesQueryDeleteAll)
	t.Run("Tickets", testTicketsQueryDeleteAll)
	t.Run("TicketAuthorities", testTicketAuthoritiesQueryDeleteAll)
	t.Run("TicketSubs", testTicketSubsQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
	t.Run("UserAgreements", testUserAgreementsQueryDeleteAll)
	t.Run("UserAuthorities", testUserAuthoritiesQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Agreements", testAgreementsSliceDeleteAll)
	t.Run("AuthServers", testAuthServersSliceDeleteAll)
	t.Run("Authorities", testAuthoritiesSliceDeleteAll)
	t.Run("Tickets", testTicketsSliceDeleteAll)
	t.Run("TicketAuthorities", testTicketAuthoritiesSliceDeleteAll)
	t.Run("TicketSubs", testTicketSubsSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
	t.Run("UserAgreements", testUserAgreementsSliceDeleteAll)
	t.Run("UserAuthorities", testUserAuthoritiesSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Agreements", testAgreementsExists)
	t.Run("AuthServers", testAuthServersExists)
	t.Run("Authorities", testAuthoritiesExists)
	t.Run("Tickets", testTicketsExists)
	t.Run("TicketAuthorities", testTicketAuthoritiesExists)
	t.Run("TicketSubs", testTicketSubsExists)
	t.Run("Users", testUsersExists)
	t.Run("UserAgreements", testUserAgreementsExists)
	t.Run("UserAuthorities", testUserAuthoritiesExists)
}

func TestFind(t *testing.T) {
	t.Run("Agreements", testAgreementsFind)
	t.Run("AuthServers", testAuthServersFind)
	t.Run("Authorities", testAuthoritiesFind)
	t.Run("Tickets", testTicketsFind)
	t.Run("TicketAuthorities", testTicketAuthoritiesFind)
	t.Run("TicketSubs", testTicketSubsFind)
	t.Run("Users", testUsersFind)
	t.Run("UserAgreements", testUserAgreementsFind)
	t.Run("UserAuthorities", testUserAuthoritiesFind)
}

func TestBind(t *testing.T) {
	t.Run("Agreements", testAgreementsBind)
	t.Run("AuthServers", testAuthServersBind)
	t.Run("Authorities", testAuthoritiesBind)
	t.Run("Tickets", testTicketsBind)
	t.Run("TicketAuthorities", testTicketAuthoritiesBind)
	t.Run("TicketSubs", testTicketSubsBind)
	t.Run("Users", testUsersBind)
	t.Run("UserAgreements", testUserAgreementsBind)
	t.Run("UserAuthorities", testUserAuthoritiesBind)
}

func TestOne(t *testing.T) {
	t.Run("Agreements", testAgreementsOne)
	t.Run("AuthServers", testAuthServersOne)
	t.Run("Authorities", testAuthoritiesOne)
	t.Run("Tickets", testTicketsOne)
	t.Run("TicketAuthorities", testTicketAuthoritiesOne)
	t.Run("TicketSubs", testTicketSubsOne)
	t.Run("Users", testUsersOne)
	t.Run("UserAgreements", testUserAgreementsOne)
	t.Run("UserAuthorities", testUserAuthoritiesOne)
}

func TestAll(t *testing.T) {
	t.Run("Agreements", testAgreementsAll)
	t.Run("AuthServers", testAuthServersAll)
	t.Run("Authorities", testAuthoritiesAll)
	t.Run("Tickets", testTicketsAll)
	t.Run("TicketAuthorities", testTicketAuthoritiesAll)
	t.Run("TicketSubs", testTicketSubsAll)
	t.Run("Users", testUsersAll)
	t.Run("UserAgreements", testUserAgreementsAll)
	t.Run("UserAuthorities", testUserAuthoritiesAll)
}

func TestCount(t *testing.T) {
	t.Run("Agreements", testAgreementsCount)
	t.Run("AuthServers", testAuthServersCount)
	t.Run("Authorities", testAuthoritiesCount)
	t.Run("Tickets", testTicketsCount)
	t.Run("TicketAuthorities", testTicketAuthoritiesCount)
	t.Run("TicketSubs", testTicketSubsCount)
	t.Run("Users", testUsersCount)
	t.Run("UserAgreements", testUserAgreementsCount)
	t.Run("UserAuthorities", testUserAuthoritiesCount)
}

func TestHooks(t *testing.T) {
	t.Run("Agreements", testAgreementsHooks)
	t.Run("AuthServers", testAuthServersHooks)
	t.Run("Authorities", testAuthoritiesHooks)
	t.Run("Tickets", testTicketsHooks)
	t.Run("TicketAuthorities", testTicketAuthoritiesHooks)
	t.Run("TicketSubs", testTicketSubsHooks)
	t.Run("Users", testUsersHooks)
	t.Run("UserAgreements", testUserAgreementsHooks)
	t.Run("UserAuthorities", testUserAuthoritiesHooks)
}

func TestInsert(t *testing.T) {
	t.Run("Agreements", testAgreementsInsert)
	t.Run("Agreements", testAgreementsInsertWhitelist)
	t.Run("AuthServers", testAuthServersInsert)
	t.Run("AuthServers", testAuthServersInsertWhitelist)
	t.Run("Authorities", testAuthoritiesInsert)
	t.Run("Authorities", testAuthoritiesInsertWhitelist)
	t.Run("Tickets", testTicketsInsert)
	t.Run("Tickets", testTicketsInsertWhitelist)
	t.Run("TicketAuthorities", testTicketAuthoritiesInsert)
	t.Run("TicketAuthorities", testTicketAuthoritiesInsertWhitelist)
	t.Run("TicketSubs", testTicketSubsInsert)
	t.Run("TicketSubs", testTicketSubsInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
	t.Run("UserAgreements", testUserAgreementsInsert)
	t.Run("UserAgreements", testUserAgreementsInsertWhitelist)
	t.Run("UserAuthorities", testUserAuthoritiesInsert)
	t.Run("UserAuthorities", testUserAuthoritiesInsertWhitelist)
}

func TestReload(t *testing.T) {
	t.Run("Agreements", testAgreementsReload)
	t.Run("AuthServers", testAuthServersReload)
	t.Run("Authorities", testAuthoritiesReload)
	t.Run("Tickets", testTicketsReload)
	t.Run("TicketAuthorities", testTicketAuthoritiesReload)
	t.Run("TicketSubs", testTicketSubsReload)
	t.Run("Users", testUsersReload)
	t.Run("UserAgreements", testUserAgreementsReload)
	t.Run("UserAuthorities", testUserAuthoritiesReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Agreements", testAgreementsReloadAll)
	t.Run("AuthServers", testAuthServersReloadAll)
	t.Run("Authorities", testAuthoritiesReloadAll)
	t.Run("Tickets", testTicketsReloadAll)
	t.Run("TicketAuthorities", testTicketAuthoritiesReloadAll)
	t.Run("TicketSubs", testTicketSubsReloadAll)
	t.Run("Users", testUsersReloadAll)
	t.Run("UserAgreements", testUserAgreementsReloadAll)
	t.Run("UserAuthorities", testUserAuthoritiesReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Agreements", testAgreementsSelect)
	t.Run("AuthServers", testAuthServersSelect)
	t.Run("Authorities", testAuthoritiesSelect)
	t.Run("Tickets", testTicketsSelect)
	t.Run("TicketAuthorities", testTicketAuthoritiesSelect)
	t.Run("TicketSubs", testTicketSubsSelect)
	t.Run("Users", testUsersSelect)
	t.Run("UserAgreements", testUserAgreementsSelect)
	t.Run("UserAuthorities", testUserAuthoritiesSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Agreements", testAgreementsUpdate)
	t.Run("AuthServers", testAuthServersUpdate)
	t.Run("Authorities", testAuthoritiesUpdate)
	t.Run("Tickets", testTicketsUpdate)
	t.Run("TicketAuthorities", testTicketAuthoritiesUpdate)
	t.Run("TicketSubs", testTicketSubsUpdate)
	t.Run("Users", testUsersUpdate)
	t.Run("UserAgreements", testUserAgreementsUpdate)
	t.Run("UserAuthorities", testUserAuthoritiesUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Agreements", testAgreementsSliceUpdateAll)
	t.Run("AuthServers", testAuthServersSliceUpdateAll)
	t.Run("Authorities", testAuthoritiesSliceUpdateAll)
	t.Run("Tickets", testTicketsSliceUpdateAll)
	t.Run("TicketAuthorities", testTicketAuthoritiesSliceUpdateAll)
	t.Run("TicketSubs", testTicketSubsSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
	t.Run("UserAgreements", testUserAgreementsSliceUpdateAll)
	t.Run("UserAuthorities", testUserAuthoritiesSliceUpdateAll)
}
