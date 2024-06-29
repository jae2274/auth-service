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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Ticket is an object representing the database table.
type Ticket struct {
	TicketID   int       `boil:"ticket_id" json:"ticket_id" toml:"ticket_id" yaml:"ticket_id"`
	UUID       string    `boil:"uuid" json:"uuid" toml:"uuid" yaml:"uuid"`
	TicketName string    `boil:"ticket_name" json:"ticket_name" toml:"ticket_name" yaml:"ticket_name"`
	UsedBy     null.Int  `boil:"used_by" json:"used_by,omitempty" toml:"used_by" yaml:"used_by,omitempty"`
	UsedDate   null.Time `boil:"used_date" json:"used_date,omitempty" toml:"used_date" yaml:"used_date,omitempty"`
	CreateDate time.Time `boil:"create_date" json:"create_date" toml:"create_date" yaml:"create_date"`

	R *ticketR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L ticketL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var TicketColumns = struct {
	TicketID   string
	UUID       string
	TicketName string
	UsedBy     string
	UsedDate   string
	CreateDate string
}{
	TicketID:   "ticket_id",
	UUID:       "uuid",
	TicketName: "ticket_name",
	UsedBy:     "used_by",
	UsedDate:   "used_date",
	CreateDate: "create_date",
}

var TicketTableColumns = struct {
	TicketID   string
	UUID       string
	TicketName string
	UsedBy     string
	UsedDate   string
	CreateDate string
}{
	TicketID:   "ticket.ticket_id",
	UUID:       "ticket.uuid",
	TicketName: "ticket.ticket_name",
	UsedBy:     "ticket.used_by",
	UsedDate:   "ticket.used_date",
	CreateDate: "ticket.create_date",
}

// Generated where

type whereHelpernull_Int struct{ field string }

func (w whereHelpernull_Int) EQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Int) NEQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Int) LT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Int) LTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Int) GT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Int) GTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelpernull_Int) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelpernull_Int) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

func (w whereHelpernull_Int) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Int) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

type whereHelpernull_Time struct{ field string }

func (w whereHelpernull_Time) EQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Time) NEQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Time) LT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Time) LTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Time) GT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Time) GTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_Time) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Time) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var TicketWhere = struct {
	TicketID   whereHelperint
	UUID       whereHelperstring
	TicketName whereHelperstring
	UsedBy     whereHelpernull_Int
	UsedDate   whereHelpernull_Time
	CreateDate whereHelpertime_Time
}{
	TicketID:   whereHelperint{field: "`ticket`.`ticket_id`"},
	UUID:       whereHelperstring{field: "`ticket`.`uuid`"},
	TicketName: whereHelperstring{field: "`ticket`.`ticket_name`"},
	UsedBy:     whereHelpernull_Int{field: "`ticket`.`used_by`"},
	UsedDate:   whereHelpernull_Time{field: "`ticket`.`used_date`"},
	CreateDate: whereHelpertime_Time{field: "`ticket`.`create_date`"},
}

// TicketRels is where relationship names are stored.
var TicketRels = struct {
	TicketAuthorities string
}{
	TicketAuthorities: "TicketAuthorities",
}

// ticketR is where relationships are stored.
type ticketR struct {
	TicketAuthorities TicketAuthoritySlice `boil:"TicketAuthorities" json:"TicketAuthorities" toml:"TicketAuthorities" yaml:"TicketAuthorities"`
}

// NewStruct creates a new relationship struct
func (*ticketR) NewStruct() *ticketR {
	return &ticketR{}
}

func (r *ticketR) GetTicketAuthorities() TicketAuthoritySlice {
	if r == nil {
		return nil
	}
	return r.TicketAuthorities
}

// ticketL is where Load methods for each relationship are stored.
type ticketL struct{}

var (
	ticketAllColumns            = []string{"ticket_id", "uuid", "ticket_name", "used_by", "used_date", "create_date"}
	ticketColumnsWithoutDefault = []string{"uuid", "ticket_name", "used_by", "used_date"}
	ticketColumnsWithDefault    = []string{"ticket_id", "create_date"}
	ticketPrimaryKeyColumns     = []string{"ticket_id"}
	ticketGeneratedColumns      = []string{}
)

type (
	// TicketSlice is an alias for a slice of pointers to Ticket.
	// This should almost always be used instead of []Ticket.
	TicketSlice []*Ticket
	// TicketHook is the signature for custom Ticket hook methods
	TicketHook func(context.Context, boil.ContextExecutor, *Ticket) error

	ticketQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	ticketType                 = reflect.TypeOf(&Ticket{})
	ticketMapping              = queries.MakeStructMapping(ticketType)
	ticketPrimaryKeyMapping, _ = queries.BindMapping(ticketType, ticketMapping, ticketPrimaryKeyColumns)
	ticketInsertCacheMut       sync.RWMutex
	ticketInsertCache          = make(map[string]insertCache)
	ticketUpdateCacheMut       sync.RWMutex
	ticketUpdateCache          = make(map[string]updateCache)
	ticketUpsertCacheMut       sync.RWMutex
	ticketUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var ticketAfterSelectMu sync.Mutex
var ticketAfterSelectHooks []TicketHook

var ticketBeforeInsertMu sync.Mutex
var ticketBeforeInsertHooks []TicketHook
var ticketAfterInsertMu sync.Mutex
var ticketAfterInsertHooks []TicketHook

var ticketBeforeUpdateMu sync.Mutex
var ticketBeforeUpdateHooks []TicketHook
var ticketAfterUpdateMu sync.Mutex
var ticketAfterUpdateHooks []TicketHook

var ticketBeforeDeleteMu sync.Mutex
var ticketBeforeDeleteHooks []TicketHook
var ticketAfterDeleteMu sync.Mutex
var ticketAfterDeleteHooks []TicketHook

var ticketBeforeUpsertMu sync.Mutex
var ticketBeforeUpsertHooks []TicketHook
var ticketAfterUpsertMu sync.Mutex
var ticketAfterUpsertHooks []TicketHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Ticket) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ticketAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Ticket) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ticketBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Ticket) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ticketAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Ticket) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ticketBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Ticket) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ticketAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Ticket) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ticketBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Ticket) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ticketAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Ticket) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ticketBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Ticket) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range ticketAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddTicketHook registers your hook function for all future operations.
func AddTicketHook(hookPoint boil.HookPoint, ticketHook TicketHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		ticketAfterSelectMu.Lock()
		ticketAfterSelectHooks = append(ticketAfterSelectHooks, ticketHook)
		ticketAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		ticketBeforeInsertMu.Lock()
		ticketBeforeInsertHooks = append(ticketBeforeInsertHooks, ticketHook)
		ticketBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		ticketAfterInsertMu.Lock()
		ticketAfterInsertHooks = append(ticketAfterInsertHooks, ticketHook)
		ticketAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		ticketBeforeUpdateMu.Lock()
		ticketBeforeUpdateHooks = append(ticketBeforeUpdateHooks, ticketHook)
		ticketBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		ticketAfterUpdateMu.Lock()
		ticketAfterUpdateHooks = append(ticketAfterUpdateHooks, ticketHook)
		ticketAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		ticketBeforeDeleteMu.Lock()
		ticketBeforeDeleteHooks = append(ticketBeforeDeleteHooks, ticketHook)
		ticketBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		ticketAfterDeleteMu.Lock()
		ticketAfterDeleteHooks = append(ticketAfterDeleteHooks, ticketHook)
		ticketAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		ticketBeforeUpsertMu.Lock()
		ticketBeforeUpsertHooks = append(ticketBeforeUpsertHooks, ticketHook)
		ticketBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		ticketAfterUpsertMu.Lock()
		ticketAfterUpsertHooks = append(ticketAfterUpsertHooks, ticketHook)
		ticketAfterUpsertMu.Unlock()
	}
}

// One returns a single ticket record from the query.
func (q ticketQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Ticket, error) {
	o := &Ticket{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for ticket")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Ticket records from the query.
func (q ticketQuery) All(ctx context.Context, exec boil.ContextExecutor) (TicketSlice, error) {
	var o []*Ticket

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Ticket slice")
	}

	if len(ticketAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Ticket records in the query.
func (q ticketQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count ticket rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q ticketQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if ticket exists")
	}

	return count > 0, nil
}

// TicketAuthorities retrieves all the ticket_authority's TicketAuthorities with an executor.
func (o *Ticket) TicketAuthorities(mods ...qm.QueryMod) ticketAuthorityQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("`ticket_authority`.`ticket_id`=?", o.TicketID),
	)

	return TicketAuthorities(queryMods...)
}

// LoadTicketAuthorities allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (ticketL) LoadTicketAuthorities(ctx context.Context, e boil.ContextExecutor, singular bool, maybeTicket interface{}, mods queries.Applicator) error {
	var slice []*Ticket
	var object *Ticket

	if singular {
		var ok bool
		object, ok = maybeTicket.(*Ticket)
		if !ok {
			object = new(Ticket)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeTicket)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeTicket))
			}
		}
	} else {
		s, ok := maybeTicket.(*[]*Ticket)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeTicket)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeTicket))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &ticketR{}
		}
		args[object.TicketID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &ticketR{}
			}
			args[obj.TicketID] = struct{}{}
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
		qm.WhereIn(`ticket_authority.ticket_id in ?`, argsSlice...),
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
			foreign.R.Ticket = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.TicketID == foreign.TicketID {
				local.R.TicketAuthorities = append(local.R.TicketAuthorities, foreign)
				if foreign.R == nil {
					foreign.R = &ticketAuthorityR{}
				}
				foreign.R.Ticket = local
				break
			}
		}
	}

	return nil
}

// AddTicketAuthorities adds the given related objects to the existing relationships
// of the ticket, optionally inserting them as new records.
// Appends related to o.R.TicketAuthorities.
// Sets related.R.Ticket appropriately.
func (o *Ticket) AddTicketAuthorities(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*TicketAuthority) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.TicketID = o.TicketID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE `ticket_authority` SET %s WHERE %s",
				strmangle.SetParamNames("`", "`", 0, []string{"ticket_id"}),
				strmangle.WhereClause("`", "`", 0, ticketAuthorityPrimaryKeyColumns),
			)
			values := []interface{}{o.TicketID, rel.TicketID, rel.AuthorityID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.TicketID = o.TicketID
		}
	}

	if o.R == nil {
		o.R = &ticketR{
			TicketAuthorities: related,
		}
	} else {
		o.R.TicketAuthorities = append(o.R.TicketAuthorities, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &ticketAuthorityR{
				Ticket: o,
			}
		} else {
			rel.R.Ticket = o
		}
	}
	return nil
}

// Tickets retrieves all the records using an executor.
func Tickets(mods ...qm.QueryMod) ticketQuery {
	mods = append(mods, qm.From("`ticket`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`ticket`.*"})
	}

	return ticketQuery{q}
}

// FindTicket retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindTicket(ctx context.Context, exec boil.ContextExecutor, ticketID int, selectCols ...string) (*Ticket, error) {
	ticketObj := &Ticket{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `ticket` where `ticket_id`=?", sel,
	)

	q := queries.Raw(query, ticketID)

	err := q.Bind(ctx, exec, ticketObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from ticket")
	}

	if err = ticketObj.doAfterSelectHooks(ctx, exec); err != nil {
		return ticketObj, err
	}

	return ticketObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Ticket) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no ticket provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(ticketColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	ticketInsertCacheMut.RLock()
	cache, cached := ticketInsertCache[key]
	ticketInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			ticketAllColumns,
			ticketColumnsWithDefault,
			ticketColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(ticketType, ticketMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(ticketType, ticketMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `ticket` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `ticket` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `ticket` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, ticketPrimaryKeyColumns))
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
		return errors.Wrap(err, "models: unable to insert into ticket")
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

	o.TicketID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == ticketMapping["ticket_id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.TicketID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for ticket")
	}

CacheNoHooks:
	if !cached {
		ticketInsertCacheMut.Lock()
		ticketInsertCache[key] = cache
		ticketInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Ticket.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Ticket) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	ticketUpdateCacheMut.RLock()
	cache, cached := ticketUpdateCache[key]
	ticketUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			ticketAllColumns,
			ticketPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update ticket, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `ticket` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, ticketPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(ticketType, ticketMapping, append(wl, ticketPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update ticket row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for ticket")
	}

	if !cached {
		ticketUpdateCacheMut.Lock()
		ticketUpdateCache[key] = cache
		ticketUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q ticketQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for ticket")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for ticket")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o TicketSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), ticketPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `ticket` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, ticketPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in ticket slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all ticket")
	}
	return rowsAff, nil
}

var mySQLTicketUniqueColumns = []string{
	"ticket_id",
	"uuid",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Ticket) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no ticket provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(ticketColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLTicketUniqueColumns, o)

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

	ticketUpsertCacheMut.RLock()
	cache, cached := ticketUpsertCache[key]
	ticketUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			ticketAllColumns,
			ticketColumnsWithDefault,
			ticketColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			ticketAllColumns,
			ticketPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("models: unable to upsert ticket, could not build update column list")
		}

		ret := strmangle.SetComplement(ticketAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`ticket`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `ticket` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(ticketType, ticketMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(ticketType, ticketMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for ticket")
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

	o.TicketID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == ticketMapping["ticket_id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(ticketType, ticketMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for ticket")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for ticket")
	}

CacheNoHooks:
	if !cached {
		ticketUpsertCacheMut.Lock()
		ticketUpsertCache[key] = cache
		ticketUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Ticket record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Ticket) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Ticket provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), ticketPrimaryKeyMapping)
	sql := "DELETE FROM `ticket` WHERE `ticket_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from ticket")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for ticket")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q ticketQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no ticketQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from ticket")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for ticket")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o TicketSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(ticketBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), ticketPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `ticket` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, ticketPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from ticket slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for ticket")
	}

	if len(ticketAfterDeleteHooks) != 0 {
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
func (o *Ticket) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindTicket(ctx, exec, o.TicketID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *TicketSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := TicketSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), ticketPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `ticket`.* FROM `ticket` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, ticketPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in TicketSlice")
	}

	*o = slice

	return nil
}

// TicketExists checks if the Ticket row exists.
func TicketExists(ctx context.Context, exec boil.ContextExecutor, ticketID int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `ticket` where `ticket_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, ticketID)
	}
	row := exec.QueryRowContext(ctx, sql, ticketID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if ticket exists")
	}

	return exists, nil
}

// Exists checks if the Ticket row exists.
func (o *Ticket) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return TicketExists(ctx, exec, o.TicketID)
}
