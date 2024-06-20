// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("TicketAuthorityToAuthorityUsingAuthority", testTicketAuthorityToOneAuthorityUsingAuthority)
	t.Run("TicketAuthorityToTicketUsingTicket", testTicketAuthorityToOneTicketUsingTicket)
	t.Run("UserAgreementToAgreementUsingAgreement", testUserAgreementToOneAgreementUsingAgreement)
	t.Run("UserAgreementToUserUsingUser", testUserAgreementToOneUserUsingUser)
	t.Run("UserAuthorityToAuthorityUsingAuthority", testUserAuthorityToOneAuthorityUsingAuthority)
	t.Run("UserAuthorityToUserUsingUser", testUserAuthorityToOneUserUsingUser)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("AgreementToUserAgreements", testAgreementToManyUserAgreements)
	t.Run("AuthorityToTicketAuthorities", testAuthorityToManyTicketAuthorities)
	t.Run("AuthorityToUserAuthorities", testAuthorityToManyUserAuthorities)
	t.Run("TicketToTicketAuthorities", testTicketToManyTicketAuthorities)
	t.Run("UserToUserAgreements", testUserToManyUserAgreements)
	t.Run("UserToUserAuthorities", testUserToManyUserAuthorities)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("TicketAuthorityToAuthorityUsingTicketAuthorities", testTicketAuthorityToOneSetOpAuthorityUsingAuthority)
	t.Run("TicketAuthorityToTicketUsingTicketAuthorities", testTicketAuthorityToOneSetOpTicketUsingTicket)
	t.Run("UserAgreementToAgreementUsingUserAgreements", testUserAgreementToOneSetOpAgreementUsingAgreement)
	t.Run("UserAgreementToUserUsingUserAgreements", testUserAgreementToOneSetOpUserUsingUser)
	t.Run("UserAuthorityToAuthorityUsingUserAuthorities", testUserAuthorityToOneSetOpAuthorityUsingAuthority)
	t.Run("UserAuthorityToUserUsingUserAuthorities", testUserAuthorityToOneSetOpUserUsingUser)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("AgreementToUserAgreements", testAgreementToManyAddOpUserAgreements)
	t.Run("AuthorityToTicketAuthorities", testAuthorityToManyAddOpTicketAuthorities)
	t.Run("AuthorityToUserAuthorities", testAuthorityToManyAddOpUserAuthorities)
	t.Run("TicketToTicketAuthorities", testTicketToManyAddOpTicketAuthorities)
	t.Run("UserToUserAgreements", testUserToManyAddOpUserAgreements)
	t.Run("UserToUserAuthorities", testUserToManyAddOpUserAuthorities)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {}