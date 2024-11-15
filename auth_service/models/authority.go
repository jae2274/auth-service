// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Authority is an object representing the database table.
type Authority struct {
	AuthorityID   int    `boil:"authority_id" json:"authority_id" toml:"authority_id" yaml:"authority_id"`
	AuthorityCode string `boil:"authority_code" json:"authority_code" toml:"authority_code" yaml:"authority_code"`
	AuthorityName string `boil:"authority_name" json:"authority_name" toml:"authority_name" yaml:"authority_name"`
	Summary       string `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`

	R *authorityR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L authorityL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var AuthorityColumns = struct {
	AuthorityID   string
	AuthorityCode string
	AuthorityName string
	Summary       string
}{
	AuthorityID:   "authority_id",
	AuthorityCode: "authority_code",
	AuthorityName: "authority_name",
	Summary:       "summary",
}

var AuthorityTableColumns = struct {
	AuthorityID   string
	AuthorityCode string
	AuthorityName string
	Summary       string
}{
	AuthorityID:   "authority.authority_id",
	AuthorityCode: "authority.authority_code",
	AuthorityName: "authority.authority_name",
	Summary:       "authority.summary",
}

// Generated where

var AuthorityWhere = struct {
	AuthorityID   whereHelperint
	AuthorityCode whereHelperstring
	AuthorityName whereHelperstring
	Summary       whereHelperstring
}{
	AuthorityID:   whereHelperint{field: "`authority`.`authority_id`"},
	AuthorityCode: whereHelperstring{field: "`authority`.`authority_code`"},
	AuthorityName: whereHelperstring{field: "`authority`.`authority_name`"},
	Summary:       whereHelperstring{field: "`authority`.`summary`"},
}

// AuthorityRels is where relationship names are stored.
var AuthorityRels = struct {
	TicketAuthorities string
	UserAuthorities   string
}{
	TicketAuthorities: "TicketAuthorities",
	UserAuthorities:   "UserAuthorities",
}

// authorityR is where relationships are stored.
type authorityR struct {
	TicketAuthorities TicketAuthoritySlice `boil:"TicketAuthorities" json:"TicketAuthorities" toml:"TicketAuthorities" yaml:"TicketAuthorities"`
	UserAuthorities   UserAuthoritySlice   `boil:"UserAuthorities" json:"UserAuthorities" toml:"UserAuthorities" yaml:"UserAuthorities"`
}

// NewStruct creates a new relationship struct
func (*authorityR) NewStruct() *authorityR {
	return &authorityR{}
}

func (r *authorityR) GetTicketAuthorities() TicketAuthoritySlice {
	if r == nil {
		return nil
	}
	return r.TicketAuthorities
}

func (r *authorityR) GetUserAuthorities() UserAuthoritySlice {
	if r == nil {
		return nil
	}
	return r.UserAuthorities
}

// authorityL is where Load methods for each relationship are stored.
type authorityL struct{}

var (
	authorityAllColumns            = []string{"authority_id", "authority_code", "authority_name", "summary"}
	authorityColumnsWithoutDefault = []string{"authority_code", "authority_name", "summary"}
	authorityColumnsWithDefault    = []string{"authority_id"}
	authorityPrimaryKeyColumns     = []string{"authority_id"}
	authorityGeneratedColumns      = []string{}
)

type (
	// AuthoritySlice is an alias for a slice of pointers to Authority.
	// This should almost always be used instead of []Authority.
	AuthoritySlice []*Authority
	// AuthorityHook is the signature for custom Authority hook methods
	AuthorityHook func(context.Context, boil.ContextExecutor, *Authority) error

	authorityQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	authorityType                 = reflect.TypeOf(&Authority{})
	authorityMapping              = queries.MakeStructMapping(authorityType)
	authorityPrimaryKeyMapping, _ = queries.BindMapping(authorityType, authorityMapping, authorityPrimaryKeyColumns)
	authorityInsertCacheMut       sync.RWMutex
	authorityInsertCache          = make(map[string]insertCache)
	authorityUpdateCacheMut       sync.RWMutex
	authorityUpdateCache          = make(map[string]updateCache)
	authorityUpsertCacheMut       sync.RWMutex
	authorityUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var authorityAfterSelectMu sync.Mutex
var authorityAfterSelectHooks []AuthorityHook

var authorityBeforeInsertMu sync.Mutex
var authorityBeforeInsertHooks []AuthorityHook
var authorityAfterInsertMu sync.Mutex
var authorityAfterInsertHooks []AuthorityHook

var authorityBeforeUpdateMu sync.Mutex
var authorityBeforeUpdateHooks []AuthorityHook
var authorityAfterUpdateMu sync.Mutex
var authorityAfterUpdateHooks []AuthorityHook

var authorityBeforeDeleteMu sync.Mutex
var authorityBeforeDeleteHooks []AuthorityHook
var authorityAfterDeleteMu sync.Mutex
var authorityAfterDeleteHooks []AuthorityHook

var authorityBeforeUpsertMu sync.Mutex
var authorityBeforeUpsertHooks []AuthorityHook
var authorityAfterUpsertMu sync.Mutex
var authorityAfterUpsertHooks []AuthorityHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Authority) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range authorityAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Authority) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range authorityBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Authority) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range authorityAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Authority) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range authorityBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Authority) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range authorityAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Authority) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range authorityBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Authority) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range authorityAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Authority) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range authorityBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Authority) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range authorityAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddAuthorityHook registers your hook function for all future operations.
func AddAuthorityHook(hookPoint boil.HookPoint, authorityHook AuthorityHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		authorityAfterSelectMu.Lock()
		authorityAfterSelectHooks = append(authorityAfterSelectHooks, authorityHook)
		authorityAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		authorityBeforeInsertMu.Lock()
		authorityBeforeInsertHooks = append(authorityBeforeInsertHooks, authorityHook)
		authorityBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		authorityAfterInsertMu.Lock()
		authorityAfterInsertHooks = append(authorityAfterInsertHooks, authorityHook)
		authorityAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		authorityBeforeUpdateMu.Lock()
		authorityBeforeUpdateHooks = append(authorityBeforeUpdateHooks, authorityHook)
		authorityBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		authorityAfterUpdateMu.Lock()
		authorityAfterUpdateHooks = append(authorityAfterUpdateHooks, authorityHook)
		authorityAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		authorityBeforeDeleteMu.Lock()
		authorityBeforeDeleteHooks = append(authorityBeforeDeleteHooks, authorityHook)
		authorityBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		authorityAfterDeleteMu.Lock()
		authorityAfterDeleteHooks = append(authorityAfterDeleteHooks, authorityHook)
		authorityAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		authorityBeforeUpsertMu.Lock()
		authorityBeforeUpsertHooks = append(authorityBeforeUpsertHooks, authorityHook)
		authorityBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		authorityAfterUpsertMu.Lock()
		authorityAfterUpsertHooks = append(authorityAfterUpsertHooks, authorityHook)
		authorityAfterUpsertMu.Unlock()
	}
}

// One returns a single authority record from the query.
func (q authorityQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Authority, error) {
	o := &Authority{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for authority")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Authority records from the query.
func (q authorityQuery) All(ctx context.Context, exec boil.ContextExecutor) (AuthoritySlice, error) {
	var o []*Authority

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Authority slice")
	}

	if len(authorityAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Authority records in the query.
func (q authorityQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count authority rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q authorityQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if authority exists")
	}

	return count > 0, nil
}

// TicketAuthorities retrieves all the ticket_authority's TicketAuthorities with an executor.
func (o *Authority) TicketAuthorities(mods ...qm.QueryMod) ticketAuthorityQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("`ticket_authority`.`authority_id`=?", o.AuthorityID),
	)

	return TicketAuthorities(queryMods...)
}

// UserAuthorities retrieves all the user_authority's UserAuthorities with an executor.
func (o *Authority) UserAuthorities(mods ...qm.QueryMod) userAuthorityQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("`user_authority`.`authority_id`=?", o.AuthorityID),
	)

	return UserAuthorities(queryMods...)
}

// LoadTicketAuthorities allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (authorityL) LoadTicketAuthorities(ctx context.Context, e boil.ContextExecutor, singular bool, maybeAuthority interface{}, mods queries.Applicator) error {
	var slice []*Authority
	var object *Authority

	if singular {
		var ok bool
		object, ok = maybeAuthority.(*Authority)
		if !ok {
			object = new(Authority)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAuthority)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAuthority))
			}
		}
	} else {
		s, ok := maybeAuthority.(*[]*Authority)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAuthority)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAuthority))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &authorityR{}
		}
		args[object.AuthorityID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &authorityR{}
			}
			args[obj.AuthorityID] = struct{}{}
		}
	}

	if len(args) == 0 {
		return nil
	}

	argsSlice := make([]interface{}, len(args))
	i := 0
	for arg := range args {
		argsSlice[i] = arg
		i++
	}

	query := NewQuery(
		qm.From(`ticket_authority`),
		qm.WhereIn(`ticket_authority.authority_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load ticket_authority")
	}

	var resultSlice []*TicketAuthority
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice ticket_authority")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on ticket_authority")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for ticket_authority")
	}

	if len(ticketAuthorityAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.TicketAuthorities = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &ticketAuthorityR{}
			}
			foreign.R.Authority = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.AuthorityID == foreign.AuthorityID {
				local.R.TicketAuthorities = append(local.R.TicketAuthorities, foreign)
				if foreign.R == nil {
					foreign.R = &ticketAuthorityR{}
				}
				foreign.R.Authority = local
				break
			}
		}
	}

	return nil
}

// LoadUserAuthorities allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (authorityL) LoadUserAuthorities(ctx context.Context, e boil.ContextExecutor, singular bool, maybeAuthority interface{}, mods queries.Applicator) error {
	var slice []*Authority
	var object *Authority

	if singular {
		var ok bool
		object, ok = maybeAuthority.(*Authority)
		if !ok {
			object = new(Authority)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAuthority)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAuthority))
			}
		}
	} else {
		s, ok := maybeAuthority.(*[]*Authority)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAuthority)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAuthority))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &authorityR{}
		}
		args[object.AuthorityID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &authorityR{}
			}
			args[obj.AuthorityID] = struct{}{}
		}
	}

	if len(args) == 0 {
		return nil
	}

	argsSlice := make([]interface{}, len(args))
	i := 0
	for arg := range args {
		argsSlice[i] = arg
		i++
	}

	query := NewQuery(
		qm.From(`user_authority`),
		qm.WhereIn(`user_authority.authority_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load user_authority")
	}

	var resultSlice []*UserAuthority
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice user_authority")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on user_authority")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for user_authority")
	}

	if len(userAuthorityAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.UserAuthorities = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &userAuthorityR{}
			}
			foreign.R.Authority = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.AuthorityID == foreign.AuthorityID {
				local.R.UserAuthorities = append(local.R.UserAuthorities, foreign)
				if foreign.R == nil {
					foreign.R = &userAuthorityR{}
				}
				foreign.R.Authority = local
				break
			}
		}
	}

	return nil
}

// AddTicketAuthorities adds the given related objects to the existing relationships
// of the authority, optionally inserting them as new records.
// Appends related to o.R.TicketAuthorities.
// Sets related.R.Authority appropriately.
func (o *Authority) AddTicketAuthorities(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*TicketAuthority) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.AuthorityID = o.AuthorityID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE `ticket_authority` SET %s WHERE %s",
				strmangle.SetParamNames("`", "`", 0, []string{"authority_id"}),
				strmangle.WhereClause("`", "`", 0, ticketAuthorityPrimaryKeyColumns),
			)
			values := []interface{}{o.AuthorityID, rel.TicketID, rel.AuthorityID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.AuthorityID = o.AuthorityID
		}
	}

	if o.R == nil {
		o.R = &authorityR{
			TicketAuthorities: related,
		}
	} else {
		o.R.TicketAuthorities = append(o.R.TicketAuthorities, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &ticketAuthorityR{
				Authority: o,
			}
		} else {
			rel.R.Authority = o
		}
	}
	return nil
}

// AddUserAuthorities adds the given related objects to the existing relationships
// of the authority, optionally inserting them as new records.
// Appends related to o.R.UserAuthorities.
// Sets related.R.Authority appropriately.
func (o *Authority) AddUserAuthorities(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*UserAuthority) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.AuthorityID = o.AuthorityID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE `user_authority` SET %s WHERE %s",
				strmangle.SetParamNames("`", "`", 0, []string{"authority_id"}),
				strmangle.WhereClause("`", "`", 0, userAuthorityPrimaryKeyColumns),
			)
			values := []interface{}{o.AuthorityID, rel.UserID, rel.AuthorityID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.AuthorityID = o.AuthorityID
		}
	}

	if o.R == nil {
		o.R = &authorityR{
			UserAuthorities: related,
		}
	} else {
		o.R.UserAuthorities = append(o.R.UserAuthorities, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &userAuthorityR{
				Authority: o,
			}
		} else {
			rel.R.Authority = o
		}
	}
	return nil
}

// Authorities retrieves all the records using an executor.
func Authorities(mods ...qm.QueryMod) authorityQuery {
	mods = append(mods, qm.From("`authority`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`authority`.*"})
	}

	return authorityQuery{q}
}

// FindAuthority retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAuthority(ctx context.Context, exec boil.ContextExecutor, authorityID int, selectCols ...string) (*Authority, error) {
	authorityObj := &Authority{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `authority` where `authority_id`=?", sel,
	)

	q := queries.Raw(query, authorityID)

	err := q.Bind(ctx, exec, authorityObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from authority")
	}

	if err = authorityObj.doAfterSelectHooks(ctx, exec); err != nil {
		return authorityObj, err
	}

	return authorityObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Authority) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no authority provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(authorityColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	authorityInsertCacheMut.RLock()
	cache, cached := authorityInsertCache[key]
	authorityInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			authorityAllColumns,
			authorityColumnsWithDefault,
			authorityColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(authorityType, authorityMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(authorityType, authorityMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `authority` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `authority` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `authority` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, authorityPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into authority")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.AuthorityID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == authorityMapping["authority_id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.AuthorityID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for authority")
	}

CacheNoHooks:
	if !cached {
		authorityInsertCacheMut.Lock()
		authorityInsertCache[key] = cache
		authorityInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Authority.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Authority) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	authorityUpdateCacheMut.RLock()
	cache, cached := authorityUpdateCache[key]
	authorityUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			authorityAllColumns,
			authorityPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update authority, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `authority` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, authorityPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(authorityType, authorityMapping, append(wl, authorityPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update authority row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for authority")
	}

	if !cached {
		authorityUpdateCacheMut.Lock()
		authorityUpdateCache[key] = cache
		authorityUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q authorityQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for authority")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for authority")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AuthoritySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authorityPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `authority` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, authorityPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in authority slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all authority")
	}
	return rowsAff, nil
}

var mySQLAuthorityUniqueColumns = []string{
	"authority_id",
	"authority_code",
	"authority_name",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Authority) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no authority provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(authorityColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLAuthorityUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	authorityUpsertCacheMut.RLock()
	cache, cached := authorityUpsertCache[key]
	authorityUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			authorityAllColumns,
			authorityColumnsWithDefault,
			authorityColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			authorityAllColumns,
			authorityPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("models: unable to upsert authority, could not build update column list")
		}

		ret := strmangle.SetComplement(authorityAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`authority`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `authority` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(authorityType, authorityMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(authorityType, authorityMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for authority")
	}

	var lastID int64
	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.AuthorityID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == authorityMapping["authority_id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(authorityType, authorityMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for authority")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for authority")
	}

CacheNoHooks:
	if !cached {
		authorityUpsertCacheMut.Lock()
		authorityUpsertCache[key] = cache
		authorityUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Authority record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Authority) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Authority provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), authorityPrimaryKeyMapping)
	sql := "DELETE FROM `authority` WHERE `authority_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from authority")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for authority")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q authorityQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no authorityQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from authority")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for authority")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AuthoritySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(authorityBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authorityPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `authority` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, authorityPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from authority slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for authority")
	}

	if len(authorityAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Authority) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindAuthority(ctx, exec, o.AuthorityID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AuthoritySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := AuthoritySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authorityPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `authority`.* FROM `authority` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, authorityPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in AuthoritySlice")
	}

	*o = slice

	return nil
}

// AuthorityExists checks if the Authority row exists.
func AuthorityExists(ctx context.Context, exec boil.ContextExecutor, authorityID int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `authority` where `authority_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, authorityID)
	}
	row := exec.QueryRowContext(ctx, sql, authorityID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if authority exists")
	}

	return exists, nil
}

// Exists checks if the Authority row exists.
func (o *Authority) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return AuthorityExists(ctx, exec, o.AuthorityID)
}
