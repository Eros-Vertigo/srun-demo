package main

import (
	"fmt"
	"github.com/dot-xiaoyuan/srun-demo/floger"
	"github.com/dot-xiaoyuan/srun-demo/slzx"
	"time"
)

func main() {
	pay := slzx.Pay{Model: slzx.PayModel{
		Username:    "yuantong",
		OutTradeNo:  time.Now().Format("0120060102150405"),
		Money:       "1",
		ProductId:   0,
		ProductName: "test",
		PayMethod:   "",
		BuyTime:     0,
		Status:      "",
		Payment:     "",
		PayType:     "",
		Remark:      "",
		Mobile:      "",
		PackageId:   "",
		SyncUrl:     "",
		ClientIP:    "",
		Email:       "",
	}}
	res, err := pay.UnifiedOrder()
	if err != nil {
		fmt.Errorf("Failed to unified order: %v", err)
	}
	floger.Debug5("res", res)
}
