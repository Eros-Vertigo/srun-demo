package main

import (
	"github.com/dot-xiaoyuan/srun-demo/floger"
	"github.com/dot-xiaoyuan/srun-demo/slzx"
)

func main() {
	pay := slzx.Pay{Model: slzx.PayModel{
		Username: "yuantong",
		//OutTradeNo:  "0620240626162228",
		OutTradeNo:  slzx.GenerateChannelSerialNumber(),
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
		floger.Errorf("Failed to unified order: %v", err)
	}
	floger.Debug5("res", res)
	//time.Sleep(5 * time.Second)
	// search
	//res, err := pay.OrderQuery()
	//if err != nil {
	//	floger.Errorf("Failed to unified order: %v", err)
	//}
	//floger.Debug5("res", res)
}
