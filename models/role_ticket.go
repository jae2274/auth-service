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

// RoleTicket is an object representing the database table.
type RoleTicket struct {
	RoleTicketID int       `boil:"role_ticket_id" json:"role_ticket_id" toml:"role_ticket_id" yaml:"role_ticket_id"`
	UUID         string    `boil:"uuid" json:"uuid" toml:"uuid" yaml:"uuid"`
	UsedBy       null.Int  `boil:"used_by" json:"used_by,omitempty" toml:"used_by" yaml:"used_by,omitempty"`
	CreateDate   time.Time `boil:"create_date" json:"create_date" toml:"create_date" yaml:"create_date"`

	R *roleTicketR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L roleTicketL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var RoleTicketColumns = struct {
	RoleTicketID string
	UUID         string
	UsedBy       string
	CreateDate   string
}{
	RoleTicketID: "role_ticket_id",
	UUID:         "uuid",
	UsedBy:       "used_by",
	CreateDate:   "create_date",
}

var RoleTicketTableColumns = struct {
	RoleTicketID string
	UUID         string
	UsedBy       string
	CreateDate   string
}{
	RoleTicketID: "role_ticket.role_ticket_id",
	UUID:         "role_ticket.uuid",
	UsedBy:       "role_ticket.used_by",
	CreateDate:   "role_ticket.create_date",
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

var RoleTicketWhere = struct {
	RoleTicketID whereHelperint
	UUID         whereHelperstring
	UsedBy       whereHelpernull_Int
	CreateDate   whereHelpertime_Time
}{
	RoleTicketID: whereHelperint{field: "`role_ticket`.`role_ticket_id`"},
	UUID:         whereHelperstring{field: "`role_ticket`.`uuid`"},
	UsedBy:       whereHelpernull_Int{field: "`role_ticket`.`used_by`"},
	CreateDate:   whereHelpertime_Time{field: "`role_ticket`.`create_date`"},
}

// RoleTicketRels is where relationship names are stored.
var RoleTicketRels = struct {
	RoleTicketRole string
}{
	RoleTicketRole: "RoleTicketRole",
}

// roleTicketR is where relationships are stored.
type roleTicketR struct {
	RoleTicketRole *RoleTicketRole `boil:"RoleTicketRole" json:"RoleTicketRole" toml:"RoleTicketRole" yaml:"RoleTicketRole"`
}

// NewStruct creates a new relationship struct
func (*roleTicketR) NewStruct() *roleTicketR {
	return &roleTicketR{}
}

func (r *roleTicketR) GetRoleTicketRole() *RoleTicketRole {
	if r == nil {
		return nil
	}
	return r.RoleTicketRole
}

// roleTicketL is where Load methods for each relationship are stored.
type roleTicketL struct{}

var (
	roleTicketAllColumns            = []string{"role_ticket_id", "uuid", "used_by", "create_date"}
	roleTicketColumnsWithoutDefault = []string{"uuid", "used_by"}
	roleTicketColumnsWithDefault    = []string{"role_ticket_id", "create_date"}
	roleTicketPrimaryKeyColumns     = []string{"role_ticket_id"}
	roleTicketGeneratedColumns      = []string{}
)

type (
	// RoleTicketSlice is an alias for a slice of pointers to RoleTicket.
	// This should almost always be used instead of []RoleTicket.
	RoleTicketSlice []*RoleTicket
	// RoleTicketHook is the signature for custom RoleTicket hook methods
	RoleTicketHook func(context.Context, boil.ContextExecutor, *RoleTicket) error

	roleTicketQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	roleTicketType                 = reflect.TypeOf(&RoleTicket{})
	roleTicketMapping              = queries.MakeStructMapping(roleTicketType)
	roleTicketPrimaryKeyMapping, _ = queries.BindMapping(roleTicketType, roleTicketMapping, roleTicketPrimaryKeyColumns)
	roleTicketInsertCacheMut       sync.RWMutex
	roleTicketInsertCache          = make(map[string]insertCache)
	roleTicketUpdateCacheMut       sync.RWMutex
	roleTicketUpdateCache          = make(map[string]updateCache)
	roleTicketUpsertCacheMut       sync.RWMutex
	roleTicketUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var roleTicketAfterSelectMu sync.Mutex
var roleTicketAfterSelectHooks []RoleTicketHook

var roleTicketBeforeInsertMu sync.Mutex
var roleTicketBeforeInsertHooks []RoleTicketHook
var roleTicketAfterInsertMu sync.Mutex
var roleTicketAfterInsertHooks []RoleTicketHook

var roleTicketBeforeUpdateMu sync.Mutex
var roleTicketBeforeUpdateHooks []RoleTicketHook
var roleTicketAfterUpdateMu sync.Mutex
var roleTicketAfterUpdateHooks []RoleTicketHook

var roleTicketBeforeDeleteMu sync.Mutex
var roleTicketBeforeDeleteHooks []RoleTicketHook
var roleTicketAfterDeleteMu sync.Mutex
var roleTicketAfterDeleteHooks []RoleTicketHook

var roleTicketBeforeUpsertMu sync.Mutex
var roleTicketBeforeUpsertHooks []RoleTicketHook
var roleTicketAfterUpsertMu sync.Mutex
var roleTicketAfterUpsertHooks []RoleTicketHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *RoleTicket) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range roleTicketAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *RoleTicket) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range roleTicketBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *RoleTicket) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range roleTicketAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *RoleTicket) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range roleTicketBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *RoleTicket) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range roleTicketAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *RoleTicket) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range roleTicketBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *RoleTicket) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range roleTicketAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *RoleTicket) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range roleTicketBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *RoleTicket) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range roleTicketAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddRoleTicketHook registers your hook function for all future operations.
func AddRoleTicketHook(hookPoint boil.HookPoint, roleTicketHook RoleTicketHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		roleTicketAfterSelectMu.Lock()
		roleTicketAfterSelectHooks = append(roleTicketAfterSelectHooks, roleTicketHook)
		roleTicketAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		roleTicketBeforeInsertMu.Lock()
		roleTicketBeforeInsertHooks = append(roleTicketBeforeInsertHooks, roleTicketHook)
		roleTicketBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		roleTicketAfterInsertMu.Lock()
		roleTicketAfterInsertHooks = append(roleTicketAfterInsertHooks, roleTicketHook)
		roleTicketAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		roleTicketBeforeUpdateMu.Lock()
		roleTicketBeforeUpdateHooks = append(roleTicketBeforeUpdateHooks, roleTicketHook)
		roleTicketBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		roleTicketAfterUpdateMu.Lock()
		roleTicketAfterUpdateHooks = append(roleTicketAfterUpdateHooks, roleTicketHook)
		roleTicketAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		roleTicketBeforeDeleteMu.Lock()
		roleTicketBeforeDeleteHooks = append(roleTicketBeforeDeleteHooks, roleTicketHook)
		roleTicketBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		roleTicketAfterDeleteMu.Lock()
		roleTicketAfterDeleteHooks = append(roleTicketAfterDeleteHooks, roleTicketHook)
		roleTicketAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		roleTicketBeforeUpsertMu.Lock()
		roleTicketBeforeUpsertHooks = append(roleTicketBeforeUpsertHooks, roleTicketHook)
		roleTicketBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		roleTicketAfterUpsertMu.Lock()
		roleTicketAfterUpsertHooks = append(roleTicketAfterUpsertHooks, roleTicketHook)
		roleTicketAfterUpsertMu.Unlock()
	}
}

// One returns a single roleTicket record from the query.
func (q roleTicketQuery) One(ctx context.Context, exec boil.ContextExecutor) (*RoleTicket, error) {
	o := &RoleTicket{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for role_ticket")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all RoleTicket records from the query.
func (q roleTicketQuery) All(ctx context.Context, exec boil.ContextExecutor) (RoleTicketSlice, error) {
	var o []*RoleTicket

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to RoleTicket slice")
	}

	if len(roleTicketAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all RoleTicket records in the query.
func (q roleTicketQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count role_ticket rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q roleTicketQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if role_ticket exists")
	}

	return count > 0, nil
}

// RoleTicketRole pointed to by the foreign key.
func (o *RoleTicket) RoleTicketRole(mods ...qm.QueryMod) roleTicketRoleQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`role_ticket_id` = ?", o.RoleTicketID),
	}

	queryMods = append(queryMods, mods...)

	return RoleTicketRoles(queryMods...)
}

// LoadRoleTicketRole allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-1 relationship.
func (roleTicketL) LoadRoleTicketRole(ctx context.Context, e boil.ContextExecutor, singular bool, maybeRoleTicket interface{}, mods queries.Applicator) error {
	var slice []*RoleTicket
	var object *RoleTicket

	if singular {
		var ok bool
		object, ok = maybeRoleTicket.(*RoleTicket)
		if !ok {
			object = new(RoleTicket)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeRoleTicket)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeRoleTicket))
			}
		}
	} else {
		s, ok := maybeRoleTicket.(*[]*RoleTicket)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeRoleTicket)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeRoleTicket))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &roleTicketR{}
		}
		args[object.RoleTicketID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &roleTicketR{}
			}

			args[obj.RoleTicketID] = struct{}{}
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
		qm.From(`role_ticket_role`),
		qm.WhereIn(`role_ticket_role.role_ticket_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load RoleTicketRole")
	}

	var resultSlice []*RoleTicketRole
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice RoleTicketRole")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for role_ticket_role")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for role_ticket_role")
	}

	if len(roleTicketRoleAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.RoleTicketRole = foreign
		if foreign.R == nil {
			foreign.R = &roleTicketRoleR{}
		}
		foreign.R.RoleTicket = object
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.RoleTicketID == foreign.RoleTicketID {
				local.R.RoleTicketRole = foreign
				if foreign.R == nil {
					foreign.R = &roleTicketRoleR{}
				}
				foreign.R.RoleTicket = local
				break
			}
		}
	}

	return nil
}

// SetRoleTicketRole of the roleTicket to the related item.
// Sets o.R.RoleTicketRole to related.
// Adds o to related.R.RoleTicket.
func (o *RoleTicket) SetRoleTicketRole(ctx context.Context, exec boil.ContextExecutor, insert bool, related *RoleTicketRole) error {
	var err error

	if insert {
		related.RoleTicketID = o.RoleTicketID

		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	} else {
		updateQuery := fmt.Sprintf(
			"UPDATE `role_ticket_role` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, []string{"role_ticket_id"}),
			strmangle.WhereClause("`", "`", 0, roleTicketRolePrimaryKeyColumns),
		)
		values := []interface{}{o.RoleTicketID, related.RoleTicketID}

		if boil.IsDebug(ctx) {
			writer := boil.DebugWriterFrom(ctx)
			fmt.Fprintln(writer, updateQuery)
			fmt.Fprintln(writer, values)
		}
		if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
			return errors.Wrap(err, "failed to update foreign table")
		}

		related.RoleTicketID = o.RoleTicketID
	}

	if o.R == nil {
		o.R = &roleTicketR{
			RoleTicketRole: related,
		}
	} else {
		o.R.RoleTicketRole = related
	}

	if related.R == nil {
		related.R = &roleTicketRoleR{
			RoleTicket: o,
		}
	} else {
		related.R.RoleTicket = o
	}
	return nil
}

// RoleTickets retrieves all the records using an executor.
func RoleTickets(mods ...qm.QueryMod) roleTicketQuery {
	mods = append(mods, qm.From("`role_ticket`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`role_ticket`.*"})
	}

	return roleTicketQuery{q}
}

// FindRoleTicket retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindRoleTicket(ctx context.Context, exec boil.ContextExecutor, roleTicketID int, selectCols ...string) (*RoleTicket, error) {
	roleTicketObj := &RoleTicket{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `role_ticket` where `role_ticket_id`=?", sel,
	)

	q := queries.Raw(query, roleTicketID)

	err := q.Bind(ctx, exec, roleTicketObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from role_ticket")
	}

	if err = roleTicketObj.doAfterSelectHooks(ctx, exec); err != nil {
		return roleTicketObj, err
	}

	return roleTicketObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *RoleTicket) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no role_ticket provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(roleTicketColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	roleTicketInsertCacheMut.RLock()
	cache, cached := roleTicketInsertCache[key]
	roleTicketInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			roleTicketAllColumns,
			roleTicketColumnsWithDefault,
			roleTicketColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(roleTicketType, roleTicketMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(roleTicketType, roleTicketMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `role_ticket` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `role_ticket` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `role_ticket` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, roleTicketPrimaryKeyColumns))
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
		return errors.Wrap(err, "models: unable to insert into role_ticket")
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

	o.RoleTicketID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == roleTicketMapping["role_ticket_id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.RoleTicketID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for role_ticket")
	}

CacheNoHooks:
	if !cached {
		roleTicketInsertCacheMut.Lock()
		roleTicketInsertCache[key] = cache
		roleTicketInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the RoleTicket.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *RoleTicket) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	roleTicketUpdateCacheMut.RLock()
	cache, cached := roleTicketUpdateCache[key]
	roleTicketUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			roleTicketAllColumns,
			roleTicketPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update role_ticket, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `role_ticket` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, roleTicketPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(roleTicketType, roleTicketMapping, append(wl, roleTicketPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update role_ticket row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for role_ticket")
	}

	if !cached {
		roleTicketUpdateCacheMut.Lock()
		roleTicketUpdateCache[key] = cache
		roleTicketUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q roleTicketQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for role_ticket")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for role_ticket")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o RoleTicketSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), roleTicketPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `role_ticket` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, roleTicketPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in roleTicket slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all roleTicket")
	}
	return rowsAff, nil
}

var mySQLRoleTicketUniqueColumns = []string{
	"role_ticket_id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *RoleTicket) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no role_ticket provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(roleTicketColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLRoleTicketUniqueColumns, o)

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

	roleTicketUpsertCacheMut.RLock()
	cache, cached := roleTicketUpsertCache[key]
	roleTicketUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			roleTicketAllColumns,
			roleTicketColumnsWithDefault,
			roleTicketColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			roleTicketAllColumns,
			roleTicketPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("models: unable to upsert role_ticket, could not build update column list")
		}

		ret := strmangle.SetComplement(roleTicketAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`role_ticket`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `role_ticket` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(roleTicketType, roleTicketMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(roleTicketType, roleTicketMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for role_ticket")
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

	o.RoleTicketID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == roleTicketMapping["role_ticket_id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(roleTicketType, roleTicketMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for role_ticket")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for role_ticket")
	}

CacheNoHooks:
	if !cached {
		roleTicketUpsertCacheMut.Lock()
		roleTicketUpsertCache[key] = cache
		roleTicketUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single RoleTicket record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *RoleTicket) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no RoleTicket provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), roleTicketPrimaryKeyMapping)
	sql := "DELETE FROM `role_ticket` WHERE `role_ticket_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from role_ticket")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for role_ticket")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q roleTicketQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no roleTicketQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from role_ticket")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for role_ticket")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o RoleTicketSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(roleTicketBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), roleTicketPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `role_ticket` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, roleTicketPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from roleTicket slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for role_ticket")
	}

	if len(roleTicketAfterDeleteHooks) != 0 {
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
func (o *RoleTicket) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindRoleTicket(ctx, exec, o.RoleTicketID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *RoleTicketSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := RoleTicketSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), roleTicketPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `role_ticket`.* FROM `role_ticket` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, roleTicketPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in RoleTicketSlice")
	}

	*o = slice

	return nil
}

// RoleTicketExists checks if the RoleTicket row exists.
func RoleTicketExists(ctx context.Context, exec boil.ContextExecutor, roleTicketID int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `role_ticket` where `role_ticket_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, roleTicketID)
	}
	row := exec.QueryRowContext(ctx, sql, roleTicketID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if role_ticket exists")
	}

	return exists, nil
}

// Exists checks if the RoleTicket row exists.
func (o *RoleTicket) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return RoleTicketExists(ctx, exec, o.RoleTicketID)
}
