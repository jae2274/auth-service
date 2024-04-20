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

func testAuthServers(t *testing.T) {
	t.Parallel()

	query := AuthServers()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testAuthServersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
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

	count, err := AuthServers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthServersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := AuthServers().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := AuthServers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthServersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := AuthServerSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := AuthServers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthServersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := AuthServerExists(ctx, tx, o.AuthServerID)
	if err != nil {
		t.Errorf("Unable to check if AuthServer exists: %s", err)
	}
	if !e {
		t.Errorf("Expected AuthServerExists to return true, but got false.")
	}
}

func testAuthServersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	authServerFound, err := FindAuthServer(ctx, tx, o.AuthServerID)
	if err != nil {
		t.Error(err)
	}

	if authServerFound == nil {
		t.Error("want a record, got nil")
	}
}

func testAuthServersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = AuthServers().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testAuthServersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := AuthServers().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testAuthServersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authServerOne := &AuthServer{}
	authServerTwo := &AuthServer{}
	if err = randomize.Struct(seed, authServerOne, authServerDBTypes, false, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}
	if err = randomize.Struct(seed, authServerTwo, authServerDBTypes, false, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = authServerOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = authServerTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := AuthServers().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testAuthServersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	authServerOne := &AuthServer{}
	authServerTwo := &AuthServer{}
	if err = randomize.Struct(seed, authServerOne, authServerDBTypes, false, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}
	if err = randomize.Struct(seed, authServerTwo, authServerDBTypes, false, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = authServerOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = authServerTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := AuthServers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func authServerBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *AuthServer) error {
	*o = AuthServer{}
	return nil
}

func authServerAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *AuthServer) error {
	*o = AuthServer{}
	return nil
}

func authServerAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *AuthServer) error {
	*o = AuthServer{}
	return nil
}

func authServerBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *AuthServer) error {
	*o = AuthServer{}
	return nil
}

func authServerAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *AuthServer) error {
	*o = AuthServer{}
	return nil
}

func authServerBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *AuthServer) error {
	*o = AuthServer{}
	return nil
}

func authServerAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *AuthServer) error {
	*o = AuthServer{}
	return nil
}

func authServerBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *AuthServer) error {
	*o = AuthServer{}
	return nil
}

func authServerAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *AuthServer) error {
	*o = AuthServer{}
	return nil
}

func testAuthServersHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &AuthServer{}
	o := &AuthServer{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, authServerDBTypes, false); err != nil {
		t.Errorf("Unable to randomize AuthServer object: %s", err)
	}

	AddAuthServerHook(boil.BeforeInsertHook, authServerBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	authServerBeforeInsertHooks = []AuthServerHook{}

	AddAuthServerHook(boil.AfterInsertHook, authServerAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	authServerAfterInsertHooks = []AuthServerHook{}

	AddAuthServerHook(boil.AfterSelectHook, authServerAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	authServerAfterSelectHooks = []AuthServerHook{}

	AddAuthServerHook(boil.BeforeUpdateHook, authServerBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	authServerBeforeUpdateHooks = []AuthServerHook{}

	AddAuthServerHook(boil.AfterUpdateHook, authServerAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	authServerAfterUpdateHooks = []AuthServerHook{}

	AddAuthServerHook(boil.BeforeDeleteHook, authServerBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	authServerBeforeDeleteHooks = []AuthServerHook{}

	AddAuthServerHook(boil.AfterDeleteHook, authServerAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	authServerAfterDeleteHooks = []AuthServerHook{}

	AddAuthServerHook(boil.BeforeUpsertHook, authServerBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	authServerBeforeUpsertHooks = []AuthServerHook{}

	AddAuthServerHook(boil.AfterUpsertHook, authServerAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	authServerAfterUpsertHooks = []AuthServerHook{}
}

func testAuthServersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := AuthServers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthServersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(authServerColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := AuthServers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthServersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
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

func testAuthServersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := AuthServerSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testAuthServersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := AuthServers().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	authServerDBTypes = map[string]string{`AuthServerID`: `int`, `AuthServerName`: `varchar`}
	_                 = bytes.MinRead
)

func testAuthServersUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(authServerPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(authServerAllColumns) == len(authServerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := AuthServers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testAuthServersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(authServerAllColumns) == len(authServerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &AuthServer{}
	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := AuthServers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, authServerDBTypes, true, authServerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(authServerAllColumns, authServerPrimaryKeyColumns) {
		fields = authServerAllColumns
	} else {
		fields = strmangle.SetComplement(
			authServerAllColumns,
			authServerPrimaryKeyColumns,
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

	slice := AuthServerSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testAuthServersUpsert(t *testing.T) {
	t.Parallel()

	if len(authServerAllColumns) == len(authServerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}
	if len(mySQLAuthServerUniqueColumns) == 0 {
		t.Skip("Skipping table with no unique columns to conflict on")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := AuthServer{}
	if err = randomize.Struct(seed, &o, authServerDBTypes, false); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert AuthServer: %s", err)
	}

	count, err := AuthServers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, authServerDBTypes, false, authServerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize AuthServer struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert AuthServer: %s", err)
	}

	count, err = AuthServers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}