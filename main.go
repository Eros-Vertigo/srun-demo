package main

import (
	"bytes"
	"encoding/json"
	"github.com/dot-xiaoyuan/srun-demo/floger"
	"github.com/dot-xiaoyuan/srun-demo/pay/slzx"
	"io"
	"log"
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
	// unified
	//res, err := pay.UnifiedOrder()
	//if err != nil {
	//	floger.Errorf("Failed to unified order: %v", err)
	//}
	//res = res.(*slzx.PayResponse)

	// search

	res, err := pay.OrderQuery()
	if err != nil {
		floger.Errorf("Failed to search order: %v", err)
	}

	var buffer bytes.Buffer

	err = PrettyEncode(res, &buffer)
	if err != nil {
		log.Fatal(err)
	}
	floger.Debug5(buffer.String())
}

//time.Sleep(5 * time.Second)
// search
//res, err := pay.OrderQuery()
//if err != nil {
//	floger.Errorf("Failed to unified order: %v", err)
//}
//floger.Debug5("res", res)

func PrettyEncode(data interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}
