package slzx

import (
	"github.com/dot-xiaoyuan/srun-demo/common/slzx"
	"github.com/dot-xiaoyuan/srun-demo/helper"
	"strings"
)

func verifySignByPublicKey(sign, plainText string) bool {
	respData := strings.ReplaceAll(plainText, " ", "+")
	signContent := respData + "&key=" + slzx.Md5Key
	return sign == strings.ToLower(helper.Md5(signContent))
}
