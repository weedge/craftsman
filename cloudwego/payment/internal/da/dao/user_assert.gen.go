// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/weedge/craftsman/cloudwego/payment/internal/da/model"
)

func newUserAssert(db *gorm.DB, opts ...gen.DOOption) userAssert {
	_userAssert := userAssert{}

	_userAssert.userAssertDo.UseDB(db, opts...)
	_userAssert.userAssertDo.UseModel(&model.UserAssert{})

	tableName := _userAssert.userAssertDo.TableName()
	_userAssert.ALL = field.NewAsterisk(tableName)
	_userAssert.ID = field.NewInt64(tableName, "id")
	_userAssert.UserID = field.NewInt64(tableName, "userId")
	_userAssert.AssetCn = field.NewInt64(tableName, "assetCn")
	_userAssert.AssetType = field.NewInt32(tableName, "assetType")
	_userAssert.Version = field.NewInt64(tableName, "version")
	_userAssert.CreatedAt = field.NewTime(tableName, "createdAt")
	_userAssert.UpdatedAt = field.NewTime(tableName, "updatedAt")

	_userAssert.fillFieldMap()

	return _userAssert
}

type userAssert struct {
	userAssertDo userAssertDo

	ALL       field.Asterisk
	ID        field.Int64
	UserID    field.Int64
	AssetCn   field.Int64
	AssetType field.Int32
	Version   field.Int64
	CreatedAt field.Time
	UpdatedAt field.Time

	fieldMap map[string]field.Expr
}

func (u userAssert) Table(newTableName string) *userAssert {
	u.userAssertDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userAssert) As(alias string) *userAssert {
	u.userAssertDo.DO = *(u.userAssertDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userAssert) updateTableName(table string) *userAssert {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewInt64(table, "id")
	u.UserID = field.NewInt64(table, "userId")
	u.AssetCn = field.NewInt64(table, "assetCn")
	u.AssetType = field.NewInt32(table, "assetType")
	u.Version = field.NewInt64(table, "version")
	u.CreatedAt = field.NewTime(table, "createdAt")
	u.UpdatedAt = field.NewTime(table, "updatedAt")

	u.fillFieldMap()

	return u
}

func (u *userAssert) WithContext(ctx context.Context) *userAssertDo {
	return u.userAssertDo.WithContext(ctx)
}

func (u userAssert) TableName() string { return u.userAssertDo.TableName() }

func (u userAssert) Alias() string { return u.userAssertDo.Alias() }

func (u *userAssert) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userAssert) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 7)
	u.fieldMap["id"] = u.ID
	u.fieldMap["userId"] = u.UserID
	u.fieldMap["assetCn"] = u.AssetCn
	u.fieldMap["assetType"] = u.AssetType
	u.fieldMap["version"] = u.Version
	u.fieldMap["createdAt"] = u.CreatedAt
	u.fieldMap["updatedAt"] = u.UpdatedAt
}

func (u userAssert) clone(db *gorm.DB) userAssert {
	u.userAssertDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userAssert) replaceDB(db *gorm.DB) userAssert {
	u.userAssertDo.ReplaceDB(db)
	return u
}

type userAssertDo struct{ gen.DO }

func (u userAssertDo) Debug() *userAssertDo {
	return u.withDO(u.DO.Debug())
}

func (u userAssertDo) WithContext(ctx context.Context) *userAssertDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userAssertDo) ReadDB() *userAssertDo {
	return u.Clauses(dbresolver.Read)
}

func (u userAssertDo) WriteDB() *userAssertDo {
	return u.Clauses(dbresolver.Write)
}

func (u userAssertDo) Session(config *gorm.Session) *userAssertDo {
	return u.withDO(u.DO.Session(config))
}

func (u userAssertDo) Clauses(conds ...clause.Expression) *userAssertDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userAssertDo) Returning(value interface{}, columns ...string) *userAssertDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userAssertDo) Not(conds ...gen.Condition) *userAssertDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userAssertDo) Or(conds ...gen.Condition) *userAssertDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userAssertDo) Select(conds ...field.Expr) *userAssertDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userAssertDo) Where(conds ...gen.Condition) *userAssertDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userAssertDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *userAssertDo {
	return u.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (u userAssertDo) Order(conds ...field.Expr) *userAssertDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userAssertDo) Distinct(cols ...field.Expr) *userAssertDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userAssertDo) Omit(cols ...field.Expr) *userAssertDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userAssertDo) Join(table schema.Tabler, on ...field.Expr) *userAssertDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userAssertDo) LeftJoin(table schema.Tabler, on ...field.Expr) *userAssertDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userAssertDo) RightJoin(table schema.Tabler, on ...field.Expr) *userAssertDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userAssertDo) Group(cols ...field.Expr) *userAssertDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userAssertDo) Having(conds ...gen.Condition) *userAssertDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userAssertDo) Limit(limit int) *userAssertDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userAssertDo) Offset(offset int) *userAssertDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userAssertDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *userAssertDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userAssertDo) Unscoped() *userAssertDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userAssertDo) Create(values ...*model.UserAssert) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userAssertDo) CreateInBatches(values []*model.UserAssert, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userAssertDo) Save(values ...*model.UserAssert) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userAssertDo) First() (*model.UserAssert, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAssert), nil
	}
}

func (u userAssertDo) Take() (*model.UserAssert, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAssert), nil
	}
}

func (u userAssertDo) Last() (*model.UserAssert, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAssert), nil
	}
}

func (u userAssertDo) Find() ([]*model.UserAssert, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserAssert), err
}

func (u userAssertDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserAssert, err error) {
	buf := make([]*model.UserAssert, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userAssertDo) FindInBatches(result *[]*model.UserAssert, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userAssertDo) Attrs(attrs ...field.AssignExpr) *userAssertDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userAssertDo) Assign(attrs ...field.AssignExpr) *userAssertDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userAssertDo) Joins(fields ...field.RelationField) *userAssertDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userAssertDo) Preload(fields ...field.RelationField) *userAssertDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userAssertDo) FirstOrInit() (*model.UserAssert, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAssert), nil
	}
}

func (u userAssertDo) FirstOrCreate() (*model.UserAssert, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAssert), nil
	}
}

func (u userAssertDo) FindByPage(offset int, limit int) (result []*model.UserAssert, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userAssertDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userAssertDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userAssertDo) Delete(models ...*model.UserAssert) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userAssertDo) withDO(do gen.Dao) *userAssertDo {
	u.DO = *do.(*gen.DO)
	return u
}