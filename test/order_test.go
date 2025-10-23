package dao

import (
	"context"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/WoWBytePaladin/go-mall/dal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/plugin/soft_delete"
	"regexp"
	"testing"
	"time"
)

func TestOrderDao_GetUserOrders(t *testing.T) {
	orderDel := soft_delete.DeletedAt(0)
	now := time.Now()
	emptyPayTime := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	orders := []*model.Order{
		{1, "12345675555", "", 1, 1, 100, 100, 0, 0, emptyPayTime, orderDel, now, now},
		{2, "12345675556", "", 1, 1, 100, 100, 0, 0, emptyPayTime, orderDel, now, now},
	}
	od := dao2.NewOrderDao(context.TODO())
	var userId int64 = 1
	offset := 10
	limit := 50
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders`")).WithArgs(userId, orderDel, limit, offset).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "order_no", "pay_trans_id", "pay_type", "user_id", "bill_money", "pay_money",
				"pay_state", "order_status", "paid_at", "is_del", "created_at", "updated_at"}).
				AddRow(
					orders[0].ID, orders[0].OrderNo, orders[0].PayTransId, orders[0].PayType, orders[0].UserId, orders[0].BillMoney, orders[0].PayMoney,
					orders[0].PayState, orders[0].OrderStatus, orders[0].PaidAt, orders[0].IsDel, orders[0].CreatedAt, orders[0].UpdatedAt,
				).AddRow(
				orders[1].ID, orders[1].OrderNo, orders[1].PayTransId, orders[1].PayType, orders[1].UserId, orders[1].BillMoney, orders[1].PayMoney,
				orders[1].PayState, orders[1].OrderStatus, orders[1].PaidAt, orders[1].IsDel, orders[1].CreatedAt, orders[1].UpdatedAt,
			),
		)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `orders`")).WithArgs(userId, orderDel).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2))
	gotOrders, totalRow, err := od.GetUserOrders(userId, offset, limit)
	assert.Nil(t, err)
	assert.Equal(t, orders, gotOrders)
	assert.Equal(t, totalRow, int64(2))
}

func TestOrderDao_UpdateOrderStatus(t *testing.T) {
	orderNewStatus := 1
	var orderId int64 = 1
	orderDel := 0
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `orders` SET")).
		WithArgs(orderNewStatus, AnyTime{}, orderId, orderDel).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	od := dao2.NewOrderDao(context.TODO())
	err := od.UpdateOrderStatus(orderId, orderNewStatus)
	assert.Nil(t, err)
}

// 定义一个AnyTime 类型，实现 sqlmock.Argument接口
// 参考自：https://qiita.com/isao_e_dev/items/c9da34c6d1f99a112207
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	// Match 方法中：判断字段值只要是time.Time 类型，就能验证通过
	_, ok := v.(time.Time)
	return ok
}
