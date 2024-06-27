package main

import (
	"encoding/json"
	"github.com/dot-xiaoyuan/srun-demo/floger"
	"github.com/dot-xiaoyuan/srun-demo/helper"
	"github.com/dot-xiaoyuan/srun-demo/slzx"
	"strings"
)

func main2() {
	//private_key := "MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQg6lrARQdATyiHlNoEHb59BnbD0c57E29slYheO7eNWJ2gCgYIKoEcz1UBgi2hRANCAATMMCG30/6aILN5bdhk4KxfEoP3VDNrY02Wk2QuZ5AV0wKbkEk5I0MnzyNujUETp4xTpaYO0m9d0fqkx497jRk7"

	plainText := "IqtnuC/RfC/QGWkzIPpRJMdOBHs/bYIVQ4cmmBeX7Qlg6GXyYXd5Zqqa7zHqgpSN4LrAhyRUPxu0Q7kALQ8N8cBkFSfuie4FsVI7RO/kBbyIA/gOdWZEUoYu8zpi+507U7NzmZTmrDFeQzcpwJaT4CzaeK6HfUVMoA0Ve/7iMYp3V1hPm9jb7tI3a/jJ8PWFM2SGZ5Cu3QKGLfxeKdZAYJ9jiSlRicK+hi7zOmL7nTsTSjZFpxn7XoMTR8q/pO4quLpgFWzxwY6vU44CzzOBYT1BMyrv0u3X8k94LXe1wdQoW9aMNtyl6CyotvVhanH2B4Lw8csIzfTEeo5fTx4EDoILWG1vpc0CJunl/X3A1acKKsKUclo+MkR8huNfhCGDBBwIWYcMy4M91i+/4513GijM9jQSli/sW4VT2tK9NQ+7vgiN77mIs91NxgEy6MIGUn0IzQG828D2QugxjEVqw/YKfg9MS/DLAkWqM9JhnTCuHCkdjEZoyBTkoIXqkDG+pmWTimBX25+uI8YDTySLf8X+LVMqWKsoayRFkGfx7P81WWg7b3ylrATS9gl/+UVS+2/xi3CZSZT/NyTNpV94Hhl2fQkZZQdSriPGA08ki38mg7kr6ZrDcfvsy23qQupvhIFQxvQFasnuJyDX1ZzIUZKRMDOIh2x7v0fKL9VSr1azwuChGQSLuW/Tb9pIkmmgZJVuvrnUlfXsWQyxRdwcTJOzCVNPtNM/2U8hqEfRcGF7FCFjyv9fkQM4c1NcJcGMFkOD0MuM9HR9VJV4TafRW3XnOAk+hZpTnC+rn9BqQizGu/3uR0gGVuukLb28b0f6bpAMUwajV6NIC0h7/sYEYxHFaDcPtKjHlCntX7q1vKDQq6M63nbUXrmP0Q1QOLR2"
	signStr := "a22c71f649d3443032c30f18320887ab"

	respdata := strings.ReplaceAll(plainText, " ", "+")
	signContent := respdata + "&key=" + slzx.Md5Key
	floger.Debug5(signContent, strings.ToLower(helper.Md5(signContent)))
	floger.Debug5(strings.ToLower(helper.Md5(signContent)) == signStr)

	text, _ := DecryptDES(slzx.Key, plainText)
	temp := make(map[string]interface{})
	json.Unmarshal([]byte(text), &temp)
	floger.Debug5(temp)
}
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
