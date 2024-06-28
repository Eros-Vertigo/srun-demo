package slzx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dot-xiaoyuan/srun-demo/common/slzx"
	"github.com/dot-xiaoyuan/srun-demo/floger"
	"github.com/dot-xiaoyuan/srun-demo/helper"
	"net/http"
	"strings"
	"time"
)

type QueryResponse struct {
	ReturnCode     string   `json:"returncode"`
	ReturnType     string   `json:"returntype"`
	ReturnMsg      string   `json:"returnmsg"`
	transDate      string   `json:"trans_date"`
	createTime     string   `json:"create_time"`
	accountDate    string   `json:"account_date"`
	DiscountAmount float64  `json:"discount_amount"`
	sign           string   `json:"sign"`
	respData       string   `json:"respdata"`
	DealAmount     float64  `json:"deal_amount"`
	OverAmount     float64  `json:"over_amount"`
	SLBusIId       string   `json:"slbusiid"`
	PresetId       string   `json:"preset_id"`
	OrderAmount    float64  `json:"order_amount"`
	Paid           string   `json:"paid"`
	ResponseTime   string   `json:"responsetime"`
	ReceiveAmount  float64  `json:"receive_amount"`
	ResData        RespData `json:"resdata"`
}

type PayModel struct {
	Username    string `json:"phone" form:"phone"`
	OutTradeNo  string `json:"order_id" form:"order_id"`
	Money       string `json:"price" form:"price"`
	ProductId   int    `json:"product_id" form:"product_id"`
	ProductName string `json:"product_name" form:"product_name"`
	PayMethod   string `json:"pay_method" form:"pay_method"`
	BuyTime     int64  `json:"buy_time"`
	Status      string `json:"status"`
	Payment     string `json:"payment"`
	PayType     string `json:"pay_type"`
	Remark      string `json:"remark"`
	Mobile      string `json:"mobile" form:"mobile"`
	PackageId   string `json:"package_id" form:"package_id"`
	SyncUrl     string `json:"sync_url" form:"sync_url"`
	ClientIP    string `json:"client_ip"`
	Email       string `json:"email" form:"email"`
}

func (p *Pay) OrderQuery() (res QueryResponse, err error) {
	// 接口地址
	nativePayRequestUrl := fmt.Sprintf("%s%s", slzx.ServiceUrl, "sltf-outside/inter/pmManageOrderQuery")

	reqData := map[string]interface{}{
		"txtype":             "03",
		"oldproduid":         "PM2000",
		"oldchannelserialno": p.Model.OutTradeNo,
		"oldmerchantorderid": p.Model.OutTradeNo,
	}

	reqDataJson, _ := json.Marshal(reqData)

	floger.Debug5("des 加密前的数据reqDataJson:", string(reqDataJson))

	reqDataDes, _ := EncryptDES(slzx.Key, string(reqDataJson))

	floger.Debug5("des 加密后的数据reqDataDes:", reqDataDes)

	signContent := reqDataDes + "&key=" + slzx.Md5Key

	floger.Debug5("md5 签名前 signContent:", signContent)

	signMsg := strings.ToLower(helper.Md5(signContent))

	floger.Debug5("md5 签名后 signMsg:", signMsg)

	headMap := map[string]interface{}{
		"version":         "V1.0",
		"charset":         "1",
		"platid":          slzx.PlatId,
		"produid":         "PM4000",
		"channelserialno": p.Model.OutTradeNo,
		"channeltime":     time.Now().Format("20060102150405"),
		"reqdata":         reqDataDes,
		"sign":            signMsg,
	}

	sendJson, _ := json.Marshal(headMap)

	floger.Debug5("发送数据 sendJson:", string(sendJson))
	resp, err := doRequest(nativePayRequestUrl, http.MethodPost, sendJson)
	if err != nil {
		floger.Error("Failed to request pmManagerOrderQuery, err:", err)
		return res, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		floger.Error("Failed to json_decode response:", err)
		return res, err
	}
	floger.Debug5("返回结果 resultStr:", res)
	// 解密
	if res.ReturnCode != "0000" {
		return res, errors.New(res.ReturnMsg)
	}
	res.respData, _ = DecryptDES(slzx.Key, res.respData)
	if err = json.Unmarshal([]byte(res.respData), &res.ResData); err != nil {
		floger.Error("Failed to json_decode response:", err)
		return res, err
	}
	return res, nil
}
