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

// Agreement is an object representing the database table.
type Agreement struct {
	AgreementID   int    `boil:"agreement_id" json:"agreement_id" toml:"agreement_id" yaml:"agreement_id"`
	AgreementCode string `boil:"agreement_code" json:"agreement_code" toml:"agreement_code" yaml:"agreement_code"`
	Summary       string `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
	IsRequired    int8   `boil:"is_required" json:"is_required" toml:"is_required" yaml:"is_required"`
	Priority      int    `boil:"priority" json:"priority" toml:"priority" yaml:"priority"`

	R *agreementR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L agreementL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var AgreementColumns = struct {
	AgreementID   string
	AgreementCode string
	Summary       string
	IsRequired    string
	Priority      string
}{
	AgreementID:   "agreement_id",
	AgreementCode: "agreement_code",
	Summary:       "summary",
	IsRequired:    "is_required",
	Priority:      "priority",
}

var AgreementTableColumns = struct {
	AgreementID   string
	AgreementCode string
	Summary       string
	IsRequired    string
	Priority      string
}{
	AgreementID:   "agreement.agreement_id",
	AgreementCode: "agreement.agreement_code",
	Summary:       "agreement.summary",
	IsRequired:    "agreement.is_required",
	Priority:      "agreement.priority",
}

// Generated where

type whereHelperint struct{ field string }

func (w whereHelperint) EQ(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint) NEQ(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint) LT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint) LTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint) GT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint) GTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod   { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod   { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod   { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) LIKE(x string) qm.QueryMod  { return qm.Where(w.field+" LIKE ?", x) }
func (w whereHelperstring) NLIKE(x string) qm.QueryMod { return qm.Where(w.field+" NOT LIKE ?", x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperint8 struct{ field string }

func (w whereHelperint8) EQ(x int8) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint8) NEQ(x int8) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint8) LT(x int8) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint8) LTE(x int8) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint8) GT(x int8) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint8) GTE(x int8) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint8) IN(slice []int8) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint8) NIN(slice []int8) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var AgreementWhere = struct {
	AgreementID   whereHelperint
	AgreementCode whereHelperstring
	Summary       whereHelperstring
	IsRequired    whereHelperint8
	Priority      whereHelperint
}{
	AgreementID:   whereHelperint{field: "`agreement`.`agreement_id`"},
	AgreementCode: whereHelperstring{field: "`agreement`.`agreement_code`"},
	Summary:       whereHelperstring{field: "`agreement`.`summary`"},
	IsRequired:    whereHelperint8{field: "`agreement`.`is_required`"},
	Priority:      whereHelperint{field: "`agreement`.`priority`"},
}

// AgreementRels is where relationship names are stored.
var AgreementRels = struct {
	UserAgreements string
}{
	UserAgreements: "UserAgreements",
}

// agreementR is where relationships are stored.
type agreementR struct {
	UserAgreements UserAgreementSlice `boil:"UserAgreements" json:"UserAgreements" toml:"UserAgreements" yaml:"UserAgreements"`
}

// NewStruct creates a new relationship struct
func (*agreementR) NewStruct() *agreementR {
	return &agreementR{}
}

func (r *agreementR) GetUserAgreements() UserAgreementSlice {
	if r == nil {
		return nil
	}
	return r.UserAgreements
}

// agreementL is where Load methods for each relationship are stored.
type agreementL struct{}

var (
	agreementAllColumns            = []string{"agreement_id", "agreement_code", "summary", "is_required", "priority"}
	agreementColumnsWithoutDefault = []string{"agreement_id", "agreement_code", "summary", "is_required"}
	agreementColumnsWithDefault    = []string{"priority"}
	agreementPrimaryKeyColumns     = []string{"agreement_id"}
	agreementGeneratedColumns      = []string{}
)

type (
	// AgreementSlice is an alias for a slice of pointers to Agreement.
	// This should almost always be used instead of []Agreement.
	AgreementSlice []*Agreement
	// AgreementHook is the signature for custom Agreement hook methods
	AgreementHook func(context.Context, boil.ContextExecutor, *Agreement) error

	agreementQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	agreementType                 = reflect.TypeOf(&Agreement{})
	agreementMapping              = queries.MakeStructMapping(agreementType)
	agreementPrimaryKeyMapping, _ = queries.BindMapping(agreementType, agreementMapping, agreementPrimaryKeyColumns)
	agreementInsertCacheMut       sync.RWMutex
	agreementInsertCache          = make(map[string]insertCache)
	agreementUpdateCacheMut       sync.RWMutex
	agreementUpdateCache          = make(map[string]updateCache)
	agreementUpsertCacheMut       sync.RWMutex
	agreementUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var agreementAfterSelectMu sync.Mutex
var agreementAfterSelectHooks []AgreementHook

var agreementBeforeInsertMu sync.Mutex
var agreementBeforeInsertHooks []AgreementHook
var agreementAfterInsertMu sync.Mutex
var agreementAfterInsertHooks []AgreementHook

var agreementBeforeUpdateMu sync.Mutex
var agreementBeforeUpdateHooks []AgreementHook
var agreementAfterUpdateMu sync.Mutex
var agreementAfterUpdateHooks []AgreementHook

var agreementBeforeDeleteMu sync.Mutex
var agreementBeforeDeleteHooks []AgreementHook
var agreementAfterDeleteMu sync.Mutex
var agreementAfterDeleteHooks []AgreementHook

var agreementBeforeUpsertMu sync.Mutex
var agreementBeforeUpsertHooks []AgreementHook
var agreementAfterUpsertMu sync.Mutex
var agreementAfterUpsertHooks []AgreementHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Agreement) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range agreementAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Agreement) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range agreementBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Agreement) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range agreementAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Agreement) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range agreementBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Agreement) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range agreementAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Agreement) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range agreementBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Agreement) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range agreementAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Agreement) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range agreementBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Agreement) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range agreementAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddAgreementHook registers your hook function for all future operations.
func AddAgreementHook(hookPoint boil.HookPoint, agreementHook AgreementHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		agreementAfterSelectMu.Lock()
		agreementAfterSelectHooks = append(agreementAfterSelectHooks, agreementHook)
		agreementAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		agreementBeforeInsertMu.Lock()
		agreementBeforeInsertHooks = append(agreementBeforeInsertHooks, agreementHook)
		agreementBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		agreementAfterInsertMu.Lock()
		agreementAfterInsertHooks = append(agreementAfterInsertHooks, agreementHook)
		agreementAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		agreementBeforeUpdateMu.Lock()
		agreementBeforeUpdateHooks = append(agreementBeforeUpdateHooks, agreementHook)
		agreementBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		agreementAfterUpdateMu.Lock()
		agreementAfterUpdateHooks = append(agreementAfterUpdateHooks, agreementHook)
		agreementAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		agreementBeforeDeleteMu.Lock()
		agreementBeforeDeleteHooks = append(agreementBeforeDeleteHooks, agreementHook)
		agreementBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		agreementAfterDeleteMu.Lock()
		agreementAfterDeleteHooks = append(agreementAfterDeleteHooks, agreementHook)
		agreementAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		agreementBeforeUpsertMu.Lock()
		agreementBeforeUpsertHooks = append(agreementBeforeUpsertHooks, agreementHook)
		agreementBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		agreementAfterUpsertMu.Lock()
		agreementAfterUpsertHooks = append(agreementAfterUpsertHooks, agreementHook)
		agreementAfterUpsertMu.Unlock()
	}
}

// One returns a single agreement record from the query.
func (q agreementQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Agreement, error) {
	o := &Agreement{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for agreement")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Agreement records from the query.
func (q agreementQuery) All(ctx context.Context, exec boil.ContextExecutor) (AgreementSlice, error) {
	var o []*Agreement

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Agreement slice")
	}

	if len(agreementAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Agreement records in the query.
func (q agreementQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count agreement rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q agreementQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if agreement exists")
	}

	return count > 0, nil
}

// UserAgreements retrieves all the user_agreement's UserAgreements with an executor.
func (o *Agreement) UserAgreements(mods ...qm.QueryMod) userAgreementQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("`user_agreement`.`agreement_id`=?", o.AgreementID),
	)

	return UserAgreements(queryMods...)
}

// LoadUserAgreements allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (agreementL) LoadUserAgreements(ctx context.Context, e boil.ContextExecutor, singular bool, maybeAgreement interface{}, mods queries.Applicator) error {
	var slice []*Agreement
	var object *Agreement

	if singular {
		var ok bool
		object, ok = maybeAgreement.(*Agreement)
		if !ok {
			object = new(Agreement)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAgreement)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAgreement))
			}
		}
	} else {
		s, ok := maybeAgreement.(*[]*Agreement)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAgreement)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAgreement))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &agreementR{}
		}
		args[object.AgreementID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &agreementR{}
			}
			args[obj.AgreementID] = struct{}{}
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
		qm.From(`user_agreement`),
		qm.WhereIn(`user_agreement.agreement_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load user_agreement")
	}

	var resultSlice []*UserAgreement
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice user_agreement")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on user_agreement")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for user_agreement")
	}

	if len(userAgreementAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.UserAgreements = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &userAgreementR{}
			}
			foreign.R.Agreement = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.AgreementID == foreign.AgreementID {
				local.R.UserAgreements = append(local.R.UserAgreements, foreign)
				if foreign.R == nil {
					foreign.R = &userAgreementR{}
				}
				foreign.R.Agreement = local
				break
			}
		}
	}

	return nil
}

// AddUserAgreements adds the given related objects to the existing relationships
// of the agreement, optionally inserting them as new records.
// Appends related to o.R.UserAgreements.
// Sets related.R.Agreement appropriately.
func (o *Agreement) AddUserAgreements(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*UserAgreement) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.AgreementID = o.AgreementID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE `user_agreement` SET %s WHERE %s",
				strmangle.SetParamNames("`", "`", 0, []string{"agreement_id"}),
				strmangle.WhereClause("`", "`", 0, userAgreementPrimaryKeyColumns),
			)
			values := []interface{}{o.AgreementID, rel.UserID, rel.AgreementID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.AgreementID = o.AgreementID
		}
	}

	if o.R == nil {
		o.R = &agreementR{
			UserAgreements: related,
		}
	} else {
		o.R.UserAgreements = append(o.R.UserAgreements, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &userAgreementR{
				Agreement: o,
			}
		} else {
			rel.R.Agreement = o
		}
	}
	return nil
}

// Agreements retrieves all the records using an executor.
func Agreements(mods ...qm.QueryMod) agreementQuery {
	mods = append(mods, qm.From("`agreement`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`agreement`.*"})
	}

	return agreementQuery{q}
}

// FindAgreement retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAgreement(ctx context.Context, exec boil.ContextExecutor, agreementID int, selectCols ...string) (*Agreement, error) {
	agreementObj := &Agreement{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `agreement` where `agreement_id`=?", sel,
	)

	q := queries.Raw(query, agreementID)

	err := q.Bind(ctx, exec, agreementObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from agreement")
	}

	if err = agreementObj.doAfterSelectHooks(ctx, exec); err != nil {
		return agreementObj, err
	}

	return agreementObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Agreement) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no agreement provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(agreementColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	agreementInsertCacheMut.RLock()
	cache, cached := agreementInsertCache[key]
	agreementInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			agreementAllColumns,
			agreementColumnsWithDefault,
			agreementColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(agreementType, agreementMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(agreementType, agreementMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `agreement` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `agreement` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `agreement` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, agreementPrimaryKeyColumns))
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
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into agreement")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.AgreementID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for agreement")
	}

CacheNoHooks:
	if !cached {
		agreementInsertCacheMut.Lock()
		agreementInsertCache[key] = cache
		agreementInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Agreement.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Agreement) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	agreementUpdateCacheMut.RLock()
	cache, cached := agreementUpdateCache[key]
	agreementUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			agreementAllColumns,
			agreementPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update agreement, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `agreement` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, agreementPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(agreementType, agreementMapping, append(wl, agreementPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update agreement row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for agreement")
	}

	if !cached {
		agreementUpdateCacheMut.Lock()
		agreementUpdateCache[key] = cache
		agreementUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q agreementQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for agreement")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for agreement")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AgreementSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), agreementPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `agreement` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, agreementPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in agreement slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all agreement")
	}
	return rowsAff, nil
}

var mySQLAgreementUniqueColumns = []string{
	"agreement_id",
	"agreement_code",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Agreement) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no agreement provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(agreementColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLAgreementUniqueColumns, o)

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

	agreementUpsertCacheMut.RLock()
	cache, cached := agreementUpsertCache[key]
	agreementUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			agreementAllColumns,
			agreementColumnsWithDefault,
			agreementColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			agreementAllColumns,
			agreementPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("models: unable to upsert agreement, could not build update column list")
		}

		ret := strmangle.SetComplement(agreementAllColumns, strmangle.SetIntersect(insert, update))

		cache.query = buildUpsertQueryMySQL(dialect, "`agreement`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `agreement` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(agreementType, agreementMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(agreementType, agreementMapping, ret)
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
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for agreement")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(agreementType, agreementMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for agreement")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for agreement")
	}

CacheNoHooks:
	if !cached {
		agreementUpsertCacheMut.Lock()
		agreementUpsertCache[key] = cache
		agreementUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Agreement record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Agreement) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Agreement provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), agreementPrimaryKeyMapping)
	sql := "DELETE FROM `agreement` WHERE `agreement_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from agreement")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for agreement")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q agreementQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no agreementQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from agreement")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for agreement")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AgreementSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(agreementBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), agreementPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `agreement` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, agreementPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from agreement slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for agreement")
	}

	if len(agreementAfterDeleteHooks) != 0 {
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
func (o *Agreement) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindAgreement(ctx, exec, o.AgreementID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AgreementSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := AgreementSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), agreementPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `agreement`.* FROM `agreement` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, agreementPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in AgreementSlice")
	}

	*o = slice

	return nil
}

// AgreementExists checks if the Agreement row exists.
func AgreementExists(ctx context.Context, exec boil.ContextExecutor, agreementID int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `agreement` where `agreement_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, agreementID)
	}
	row := exec.QueryRowContext(ctx, sql, agreementID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if agreement exists")
	}

	return exists, nil
}

// Exists checks if the Agreement row exists.
func (o *Agreement) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return AgreementExists(ctx, exec, o.AgreementID)
}
