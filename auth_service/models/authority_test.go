// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testAuthorities(t *testing.T) {
	t.Parallel()

	query := Authorities()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testAuthoritiesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Authorities().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthoritiesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Authorities().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Authorities().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthoritiesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := AuthoritySlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Authorities().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthoritiesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := AuthorityExists(ctx, tx, o.AuthorityID)
	if err != nil {
		t.Errorf("Unable to check if Authority exists: %s", err)
	}
	if !e {
		t.Errorf("Expected AuthorityExists to return true, but got false.")
	}
}

func testAuthoritiesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	authorityFound, err := FindAuthority(ctx, tx, o.AuthorityID)
	if err != nil {
		t.Error(err)
	}

	if authorityFound == nil {
		t.Error("want a record, got nil")
	}
}

func testAuthoritiesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Authorities().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testAuthoritiesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Authorities().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testAuthoritiesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorityOne := &Authority{}
	authorityTwo := &Authority{}
	if err = randomize.Struct(seed, authorityOne, authorityDBTypes, false, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}
	if err = randomize.Struct(seed, authorityTwo, authorityDBTypes, false, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = authorityOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = authorityTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Authorities().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testAuthoritiesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	authorityOne := &Authority{}
	authorityTwo := &Authority{}
	if err = randomize.Struct(seed, authorityOne, authorityDBTypes, false, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}
	if err = randomize.Struct(seed, authorityTwo, authorityDBTypes, false, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = authorityOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = authorityTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Authorities().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func authorityBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *Authority) error {
	*o = Authority{}
	return nil
}

func authorityAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *Authority) error {
	*o = Authority{}
	return nil
}

func authorityAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *Authority) error {
	*o = Authority{}
	return nil
}

func authorityBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Authority) error {
	*o = Authority{}
	return nil
}

func authorityAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Authority) error {
	*o = Authority{}
	return nil
}

func authorityBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Authority) error {
	*o = Authority{}
	return nil
}

func authorityAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Authority) error {
	*o = Authority{}
	return nil
}

func authorityBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Authority) error {
	*o = Authority{}
	return nil
}

func authorityAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Authority) error {
	*o = Authority{}
	return nil
}

func testAuthoritiesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &Authority{}
	o := &Authority{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, authorityDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Authority object: %s", err)
	}

	AddAuthorityHook(boil.BeforeInsertHook, authorityBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	authorityBeforeInsertHooks = []AuthorityHook{}

	AddAuthorityHook(boil.AfterInsertHook, authorityAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	authorityAfterInsertHooks = []AuthorityHook{}

	AddAuthorityHook(boil.AfterSelectHook, authorityAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	authorityAfterSelectHooks = []AuthorityHook{}

	AddAuthorityHook(boil.BeforeUpdateHook, authorityBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	authorityBeforeUpdateHooks = []AuthorityHook{}

	AddAuthorityHook(boil.AfterUpdateHook, authorityAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	authorityAfterUpdateHooks = []AuthorityHook{}

	AddAuthorityHook(boil.BeforeDeleteHook, authorityBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	authorityBeforeDeleteHooks = []AuthorityHook{}

	AddAuthorityHook(boil.AfterDeleteHook, authorityAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	authorityAfterDeleteHooks = []AuthorityHook{}

	AddAuthorityHook(boil.BeforeUpsertHook, authorityBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	authorityBeforeUpsertHooks = []AuthorityHook{}

	AddAuthorityHook(boil.AfterUpsertHook, authorityAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	authorityAfterUpsertHooks = []AuthorityHook{}
}

func testAuthoritiesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Authorities().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthoritiesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(authorityColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Authorities().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthorityToManyTicketAuthorities(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Authority
	var b, c TicketAuthority

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, ticketAuthorityDBTypes, false, ticketAuthorityColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, ticketAuthorityDBTypes, false, ticketAuthorityColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.AuthorityID = a.AuthorityID
	c.AuthorityID = a.AuthorityID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.TicketAuthorities().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.AuthorityID == b.AuthorityID {
			bFound = true
		}
		if v.AuthorityID == c.AuthorityID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := AuthoritySlice{&a}
	if err = a.L.LoadTicketAuthorities(ctx, tx, false, (*[]*Authority)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TicketAuthorities); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.TicketAuthorities = nil
	if err = a.L.LoadTicketAuthorities(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TicketAuthorities); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testAuthorityToManyUserAuthorities(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Authority
	var b, c UserAuthority

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, userAuthorityDBTypes, false, userAuthorityColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userAuthorityDBTypes, false, userAuthorityColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.AuthorityID = a.AuthorityID
	c.AuthorityID = a.AuthorityID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.UserAuthorities().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.AuthorityID == b.AuthorityID {
			bFound = true
		}
		if v.AuthorityID == c.AuthorityID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := AuthoritySlice{&a}
	if err = a.L.LoadUserAuthorities(ctx, tx, false, (*[]*Authority)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.UserAuthorities); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.UserAuthorities = nil
	if err = a.L.LoadUserAuthorities(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.UserAuthorities); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testAuthorityToManyAddOpTicketAuthorities(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Authority
	var b, c, d, e TicketAuthority

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorityDBTypes, false, strmangle.SetComplement(authorityPrimaryKeyColumns, authorityColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*TicketAuthority{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, ticketAuthorityDBTypes, false, strmangle.SetComplement(ticketAuthorityPrimaryKeyColumns, ticketAuthorityColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*TicketAuthority{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddTicketAuthorities(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.AuthorityID != first.AuthorityID {
			t.Error("foreign key was wrong value", a.AuthorityID, first.AuthorityID)
		}
		if a.AuthorityID != second.AuthorityID {
			t.Error("foreign key was wrong value", a.AuthorityID, second.AuthorityID)
		}

		if first.R.Authority != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Authority != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.TicketAuthorities[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.TicketAuthorities[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.TicketAuthorities().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testAuthorityToManyAddOpUserAuthorities(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Authority
	var b, c, d, e UserAuthority

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorityDBTypes, false, strmangle.SetComplement(authorityPrimaryKeyColumns, authorityColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*UserAuthority{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, userAuthorityDBTypes, false, strmangle.SetComplement(userAuthorityPrimaryKeyColumns, userAuthorityColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*UserAuthority{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddUserAuthorities(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.AuthorityID != first.AuthorityID {
			t.Error("foreign key was wrong value", a.AuthorityID, first.AuthorityID)
		}
		if a.AuthorityID != second.AuthorityID {
			t.Error("foreign key was wrong value", a.AuthorityID, second.AuthorityID)
		}

		if first.R.Authority != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Authority != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.UserAuthorities[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.UserAuthorities[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.UserAuthorities().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testAuthoritiesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testAuthoritiesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := AuthoritySlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testAuthoritiesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Authorities().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	authorityDBTypes = map[string]string{`AuthorityID`: `int`, `AuthorityCode`: `varchar`, `AuthorityName`: `varchar`, `Summary`: `varchar`}
	_                = bytes.MinRead
)

func testAuthoritiesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(authorityPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(authorityAllColumns) == len(authorityPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Authorities().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testAuthoritiesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(authorityAllColumns) == len(authorityPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Authority{}
	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Authorities().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, authorityDBTypes, true, authorityPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(authorityAllColumns, authorityPrimaryKeyColumns) {
		fields = authorityAllColumns
	} else {
		fields = strmangle.SetComplement(
			authorityAllColumns,
			authorityPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := AuthoritySlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testAuthoritiesUpsert(t *testing.T) {
	t.Parallel()

	if len(authorityAllColumns) == len(authorityPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}
	if len(mySQLAuthorityUniqueColumns) == 0 {
		t.Skip("Skipping table with no unique columns to conflict on")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Authority{}
	if err = randomize.Struct(seed, &o, authorityDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Authority: %s", err)
	}

	count, err := Authorities().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, authorityDBTypes, false, authorityPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Authority struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Authority: %s", err)
	}

	count, err = Authorities().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
