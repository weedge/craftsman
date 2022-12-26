// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/weedge/craftsman/cloudwego/payment/internal/da/model"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
)

func init() {
	InitializeDB()
	err := db.AutoMigrate(&model.UserAsset{})
	if err != nil {
		fmt.Printf("Error: AutoMigrate(&model.UserAsset{}) fail: %s", err)
	}
}

func Test_userAssetQuery(t *testing.T) {
	userAsset := newUserAsset(db)
	userAsset = *userAsset.As(userAsset.TableName())
	_do := userAsset.WithContext(context.Background()).Debug()

	primaryKey := field.NewString(userAsset.TableName(), clause.PrimaryKey)
	_, err := _do.Unscoped().Where(primaryKey.IsNotNull()).Delete()
	if err != nil {
		t.Error("clean table <user_asset> fail:", err)
		return
	}

	_, ok := userAsset.GetFieldByName("")
	if ok {
		t.Error("GetFieldByName(\"\") from userAsset success")
	}

	err = _do.Create(&model.UserAsset{})
	if err != nil {
		t.Error("create item in table <user_asset> fail:", err)
	}

	err = _do.Save(&model.UserAsset{})
	if err != nil {
		t.Error("create item in table <user_asset> fail:", err)
	}

	err = _do.CreateInBatches([]*model.UserAsset{{}, {}}, 10)
	if err != nil {
		t.Error("create item in table <user_asset> fail:", err)
	}

	_, err = _do.Select(userAsset.ALL).Take()
	if err != nil {
		t.Error("Take() on table <user_asset> fail:", err)
	}

	_, err = _do.First()
	if err != nil {
		t.Error("First() on table <user_asset> fail:", err)
	}

	_, err = _do.Last()
	if err != nil {
		t.Error("First() on table <user_asset> fail:", err)
	}

	_, err = _do.Where(primaryKey.IsNotNull()).FindInBatch(10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatch() on table <user_asset> fail:", err)
	}

	err = _do.Where(primaryKey.IsNotNull()).FindInBatches(&[]*model.UserAsset{}, 10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatches() on table <user_asset> fail:", err)
	}

	_, err = _do.Select(userAsset.ALL).Where(primaryKey.IsNotNull()).Order(primaryKey.Desc()).Find()
	if err != nil {
		t.Error("Find() on table <user_asset> fail:", err)
	}

	_, err = _do.Distinct(primaryKey).Take()
	if err != nil {
		t.Error("select Distinct() on table <user_asset> fail:", err)
	}

	_, err = _do.Select(userAsset.ALL).Omit(primaryKey).Take()
	if err != nil {
		t.Error("Omit() on table <user_asset> fail:", err)
	}

	_, err = _do.Group(primaryKey).Find()
	if err != nil {
		t.Error("Group() on table <user_asset> fail:", err)
	}

	_, err = _do.Scopes(func(dao gen.Dao) gen.Dao { return dao.Where(primaryKey.IsNotNull()) }).Find()
	if err != nil {
		t.Error("Scopes() on table <user_asset> fail:", err)
	}

	_, _, err = _do.FindByPage(0, 1)
	if err != nil {
		t.Error("FindByPage() on table <user_asset> fail:", err)
	}

	_, err = _do.ScanByPage(&model.UserAsset{}, 0, 1)
	if err != nil {
		t.Error("ScanByPage() on table <user_asset> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrInit()
	if err != nil {
		t.Error("FirstOrInit() on table <user_asset> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrCreate()
	if err != nil {
		t.Error("FirstOrCreate() on table <user_asset> fail:", err)
	}

	var _a _another
	var _aPK = field.NewString(_a.TableName(), clause.PrimaryKey)

	err = _do.Join(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("Join() on table <user_asset> fail:", err)
	}

	err = _do.LeftJoin(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("LeftJoin() on table <user_asset> fail:", err)
	}

	_, err = _do.Not().Or().Clauses().Take()
	if err != nil {
		t.Error("Not/Or/Clauses on table <user_asset> fail:", err)
	}
}
