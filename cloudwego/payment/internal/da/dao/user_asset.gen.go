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

func newUserAsset(db *gorm.DB, opts ...gen.DOOption) userAsset {
	_userAsset := userAsset{}

	_userAsset.userAssetDo.UseDB(db, opts...)
	_userAsset.userAssetDo.UseModel(&model.UserAsset{})

	tableName := _userAsset.userAssetDo.TableName()
	_userAsset.ALL = field.NewAsterisk(tableName)
	_userAsset.UserID = field.NewInt64(tableName, "userId")
	_userAsset.AssetCn = field.NewInt64(tableName, "assetCn")
	_userAsset.AssetType = field.NewInt32(tableName, "assetType")
	_userAsset.Version = field.NewInt64(tableName, "version")
	_userAsset.CreatedAt = field.NewTime(tableName, "createdAt")
	_userAsset.UpdatedAt = field.NewTime(tableName, "updatedAt")

	_userAsset.fillFieldMap()

	return _userAsset
}

type userAsset struct {
	userAssetDo userAssetDo

	ALL       field.Asterisk
	UserID    field.Int64
	AssetCn   field.Int64
	AssetType field.Int32
	Version   field.Int64
	CreatedAt field.Time
	UpdatedAt field.Time

	fieldMap map[string]field.Expr
}

func (u userAsset) Table(newTableName string) *userAsset {
	u.userAssetDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userAsset) As(alias string) *userAsset {
	u.userAssetDo.DO = *(u.userAssetDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userAsset) updateTableName(table string) *userAsset {
	u.ALL = field.NewAsterisk(table)
	u.UserID = field.NewInt64(table, "userId")
	u.AssetCn = field.NewInt64(table, "assetCn")
	u.AssetType = field.NewInt32(table, "assetType")
	u.Version = field.NewInt64(table, "version")
	u.CreatedAt = field.NewTime(table, "createdAt")
	u.UpdatedAt = field.NewTime(table, "updatedAt")

	u.fillFieldMap()

	return u
}

func (u *userAsset) WithContext(ctx context.Context) *userAssetDo {
	return u.userAssetDo.WithContext(ctx)
}

func (u userAsset) TableName() string { return u.userAssetDo.TableName() }

func (u userAsset) Alias() string { return u.userAssetDo.Alias() }

func (u *userAsset) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userAsset) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 6)
	u.fieldMap["userId"] = u.UserID
	u.fieldMap["assetCn"] = u.AssetCn
	u.fieldMap["assetType"] = u.AssetType
	u.fieldMap["version"] = u.Version
	u.fieldMap["createdAt"] = u.CreatedAt
	u.fieldMap["updatedAt"] = u.UpdatedAt
}

func (u userAsset) clone(db *gorm.DB) userAsset {
	u.userAssetDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userAsset) replaceDB(db *gorm.DB) userAsset {
	u.userAssetDo.ReplaceDB(db)
	return u
}

type userAssetDo struct{ gen.DO }

func (u userAssetDo) Debug() *userAssetDo {
	return u.withDO(u.DO.Debug())
}

func (u userAssetDo) WithContext(ctx context.Context) *userAssetDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userAssetDo) ReadDB() *userAssetDo {
	return u.Clauses(dbresolver.Read)
}

func (u userAssetDo) WriteDB() *userAssetDo {
	return u.Clauses(dbresolver.Write)
}

func (u userAssetDo) Session(config *gorm.Session) *userAssetDo {
	return u.withDO(u.DO.Session(config))
}

func (u userAssetDo) Clauses(conds ...clause.Expression) *userAssetDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userAssetDo) Returning(value interface{}, columns ...string) *userAssetDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userAssetDo) Not(conds ...gen.Condition) *userAssetDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userAssetDo) Or(conds ...gen.Condition) *userAssetDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userAssetDo) Select(conds ...field.Expr) *userAssetDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userAssetDo) Where(conds ...gen.Condition) *userAssetDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userAssetDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *userAssetDo {
	return u.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (u userAssetDo) Order(conds ...field.Expr) *userAssetDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userAssetDo) Distinct(cols ...field.Expr) *userAssetDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userAssetDo) Omit(cols ...field.Expr) *userAssetDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userAssetDo) Join(table schema.Tabler, on ...field.Expr) *userAssetDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userAssetDo) LeftJoin(table schema.Tabler, on ...field.Expr) *userAssetDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userAssetDo) RightJoin(table schema.Tabler, on ...field.Expr) *userAssetDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userAssetDo) Group(cols ...field.Expr) *userAssetDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userAssetDo) Having(conds ...gen.Condition) *userAssetDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userAssetDo) Limit(limit int) *userAssetDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userAssetDo) Offset(offset int) *userAssetDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userAssetDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *userAssetDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userAssetDo) Unscoped() *userAssetDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userAssetDo) Create(values ...*model.UserAsset) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userAssetDo) CreateInBatches(values []*model.UserAsset, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userAssetDo) Save(values ...*model.UserAsset) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userAssetDo) First() (*model.UserAsset, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAsset), nil
	}
}

func (u userAssetDo) Take() (*model.UserAsset, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAsset), nil
	}
}

func (u userAssetDo) Last() (*model.UserAsset, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAsset), nil
	}
}

func (u userAssetDo) Find() ([]*model.UserAsset, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserAsset), err
}

func (u userAssetDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserAsset, err error) {
	buf := make([]*model.UserAsset, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userAssetDo) FindInBatches(result *[]*model.UserAsset, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userAssetDo) Attrs(attrs ...field.AssignExpr) *userAssetDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userAssetDo) Assign(attrs ...field.AssignExpr) *userAssetDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userAssetDo) Joins(fields ...field.RelationField) *userAssetDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userAssetDo) Preload(fields ...field.RelationField) *userAssetDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userAssetDo) FirstOrInit() (*model.UserAsset, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAsset), nil
	}
}

func (u userAssetDo) FirstOrCreate() (*model.UserAsset, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAsset), nil
	}
}

func (u userAssetDo) FindByPage(offset int, limit int) (result []*model.UserAsset, count int64, err error) {
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

func (u userAssetDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userAssetDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userAssetDo) Delete(models ...*model.UserAsset) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userAssetDo) withDO(do gen.Dao) *userAssetDo {
	u.DO = *do.(*gen.DO)
	return u
}
