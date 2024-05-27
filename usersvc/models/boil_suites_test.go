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
	t.Run("Tickets", testTickets)
	t.Run("TicketRoles", testTicketRoles)
	t.Run("Users", testUsers)
	t.Run("UserAgreements", testUserAgreements)
	t.Run("UserRoles", testUserRoles)
}

func TestDelete(t *testing.T) {
	t.Run("Agreements", testAgreementsDelete)
	t.Run("AuthServers", testAuthServersDelete)
	t.Run("Tickets", testTicketsDelete)
	t.Run("TicketRoles", testTicketRolesDelete)
	t.Run("Users", testUsersDelete)
	t.Run("UserAgreements", testUserAgreementsDelete)
	t.Run("UserRoles", testUserRolesDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Agreements", testAgreementsQueryDeleteAll)
	t.Run("AuthServers", testAuthServersQueryDeleteAll)
	t.Run("Tickets", testTicketsQueryDeleteAll)
	t.Run("TicketRoles", testTicketRolesQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
	t.Run("UserAgreements", testUserAgreementsQueryDeleteAll)
	t.Run("UserRoles", testUserRolesQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Agreements", testAgreementsSliceDeleteAll)
	t.Run("AuthServers", testAuthServersSliceDeleteAll)
	t.Run("Tickets", testTicketsSliceDeleteAll)
	t.Run("TicketRoles", testTicketRolesSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
	t.Run("UserAgreements", testUserAgreementsSliceDeleteAll)
	t.Run("UserRoles", testUserRolesSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Agreements", testAgreementsExists)
	t.Run("AuthServers", testAuthServersExists)
	t.Run("Tickets", testTicketsExists)
	t.Run("TicketRoles", testTicketRolesExists)
	t.Run("Users", testUsersExists)
	t.Run("UserAgreements", testUserAgreementsExists)
	t.Run("UserRoles", testUserRolesExists)
}

func TestFind(t *testing.T) {
	t.Run("Agreements", testAgreementsFind)
	t.Run("AuthServers", testAuthServersFind)
	t.Run("Tickets", testTicketsFind)
	t.Run("TicketRoles", testTicketRolesFind)
	t.Run("Users", testUsersFind)
	t.Run("UserAgreements", testUserAgreementsFind)
	t.Run("UserRoles", testUserRolesFind)
}

func TestBind(t *testing.T) {
	t.Run("Agreements", testAgreementsBind)
	t.Run("AuthServers", testAuthServersBind)
	t.Run("Tickets", testTicketsBind)
	t.Run("TicketRoles", testTicketRolesBind)
	t.Run("Users", testUsersBind)
	t.Run("UserAgreements", testUserAgreementsBind)
	t.Run("UserRoles", testUserRolesBind)
}

func TestOne(t *testing.T) {
	t.Run("Agreements", testAgreementsOne)
	t.Run("AuthServers", testAuthServersOne)
	t.Run("Tickets", testTicketsOne)
	t.Run("TicketRoles", testTicketRolesOne)
	t.Run("Users", testUsersOne)
	t.Run("UserAgreements", testUserAgreementsOne)
	t.Run("UserRoles", testUserRolesOne)
}

func TestAll(t *testing.T) {
	t.Run("Agreements", testAgreementsAll)
	t.Run("AuthServers", testAuthServersAll)
	t.Run("Tickets", testTicketsAll)
	t.Run("TicketRoles", testTicketRolesAll)
	t.Run("Users", testUsersAll)
	t.Run("UserAgreements", testUserAgreementsAll)
	t.Run("UserRoles", testUserRolesAll)
}

func TestCount(t *testing.T) {
	t.Run("Agreements", testAgreementsCount)
	t.Run("AuthServers", testAuthServersCount)
	t.Run("Tickets", testTicketsCount)
	t.Run("TicketRoles", testTicketRolesCount)
	t.Run("Users", testUsersCount)
	t.Run("UserAgreements", testUserAgreementsCount)
	t.Run("UserRoles", testUserRolesCount)
}

func TestHooks(t *testing.T) {
	t.Run("Agreements", testAgreementsHooks)
	t.Run("AuthServers", testAuthServersHooks)
	t.Run("Tickets", testTicketsHooks)
	t.Run("TicketRoles", testTicketRolesHooks)
	t.Run("Users", testUsersHooks)
	t.Run("UserAgreements", testUserAgreementsHooks)
	t.Run("UserRoles", testUserRolesHooks)
}

func TestInsert(t *testing.T) {
	t.Run("Agreements", testAgreementsInsert)
	t.Run("Agreements", testAgreementsInsertWhitelist)
	t.Run("AuthServers", testAuthServersInsert)
	t.Run("AuthServers", testAuthServersInsertWhitelist)
	t.Run("Tickets", testTicketsInsert)
	t.Run("Tickets", testTicketsInsertWhitelist)
	t.Run("TicketRoles", testTicketRolesInsert)
	t.Run("TicketRoles", testTicketRolesInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
	t.Run("UserAgreements", testUserAgreementsInsert)
	t.Run("UserAgreements", testUserAgreementsInsertWhitelist)
	t.Run("UserRoles", testUserRolesInsert)
	t.Run("UserRoles", testUserRolesInsertWhitelist)
}

func TestReload(t *testing.T) {
	t.Run("Agreements", testAgreementsReload)
	t.Run("AuthServers", testAuthServersReload)
	t.Run("Tickets", testTicketsReload)
	t.Run("TicketRoles", testTicketRolesReload)
	t.Run("Users", testUsersReload)
	t.Run("UserAgreements", testUserAgreementsReload)
	t.Run("UserRoles", testUserRolesReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Agreements", testAgreementsReloadAll)
	t.Run("AuthServers", testAuthServersReloadAll)
	t.Run("Tickets", testTicketsReloadAll)
	t.Run("TicketRoles", testTicketRolesReloadAll)
	t.Run("Users", testUsersReloadAll)
	t.Run("UserAgreements", testUserAgreementsReloadAll)
	t.Run("UserRoles", testUserRolesReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Agreements", testAgreementsSelect)
	t.Run("AuthServers", testAuthServersSelect)
	t.Run("Tickets", testTicketsSelect)
	t.Run("TicketRoles", testTicketRolesSelect)
	t.Run("Users", testUsersSelect)
	t.Run("UserAgreements", testUserAgreementsSelect)
	t.Run("UserRoles", testUserRolesSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Agreements", testAgreementsUpdate)
	t.Run("AuthServers", testAuthServersUpdate)
	t.Run("Tickets", testTicketsUpdate)
	t.Run("TicketRoles", testTicketRolesUpdate)
	t.Run("Users", testUsersUpdate)
	t.Run("UserAgreements", testUserAgreementsUpdate)
	t.Run("UserRoles", testUserRolesUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Agreements", testAgreementsSliceUpdateAll)
	t.Run("AuthServers", testAuthServersSliceUpdateAll)
	t.Run("Tickets", testTicketsSliceUpdateAll)
	t.Run("TicketRoles", testTicketRolesSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
	t.Run("UserAgreements", testUserAgreementsSliceUpdateAll)
	t.Run("UserRoles", testUserRolesSliceUpdateAll)
}
