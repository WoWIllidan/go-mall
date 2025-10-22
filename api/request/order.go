package request

type OrderCreate struct {
	CartItemIdList []int64 `json:"cart_item_id_list" binding:"required"`
	UserAddressId  int64   `json:"user_address_id" binding:"required"`
}
