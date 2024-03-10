package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"one-api/common"
	"os"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Redemption struct {
	Id           int    `json:"id"`
	UserId       int    `json:"user_id"`
	Key          string `json:"key" gorm:"type:char(32);uniqueIndex"`
	Status       int    `json:"status" gorm:"default:1"`
	Name         string `json:"name" gorm:"index"`
	Quota        int    `json:"quota" gorm:"default:100"`
	CreatedTime  int64  `json:"created_time" gorm:"bigint"`
	RedeemedTime int64  `json:"redeemed_time" gorm:"bigint"`
	Count        int    `json:"count" gorm:"-:all"` // only for api request
}

type RechargeLog struct {
	Id           int    `json:"id"`
	UserId       int    `json:"user_id"`
	TradeNo      string `json:"trade_no" gorm:"type:char(20);uniqueIndex"`
	Status       int    `json:"status" gorm:"default:1"` //1:未支付，2：已支付
	Name         string `json:"name" gorm:"index"`
	RedeemedTime int64  `json:"redeemed_time" gorm:"bigint"`
	CreatedTime  int64  `json:"created_time" gorm:"bigint"`
	Quota        int    `json:"quota" gorm:"default:100"`
}

var allowedRedemptionslOrderFields = map[string]bool{
	"id":            true,
	"name":          true,
	"status":        true,
	"quota":         true,
	"created_time":  true,
	"redeemed_time": true,
}

func GetRedemptionsList(params *GenericParams) (*DataResult[Redemption], error) {
	var redemptions []*Redemption
	db := DB
	if params.Keyword != "" {
		db = db.Where("id = ? or name LIKE ?", common.String2Int(params.Keyword), params.Keyword+"%")
	}

	return PaginateAndOrder[Redemption](db, &params.PaginationParams, &redemptions, allowedRedemptionslOrderFields)
}

func GetRedemptionById(id int) (*Redemption, error) {
	if id == 0 {
		return nil, errors.New("id 为空！")
	}
	redemption := Redemption{Id: id}
	var err error = nil
	err = DB.First(&redemption, "id = ?", id).Error
	return &redemption, err
}

func Redeem(key string, userId int) (quota int, err error) {
	if key == "" {
		return 0, errors.New("未提供兑换码")
	}
	if userId == 0 {
		return 0, errors.New("无效的 user id")
	}
	redemption := &Redemption{}

	keyCol := "`key`"
	if common.UsingPostgreSQL {
		keyCol = `"key"`
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Set("gorm:query_option", "FOR UPDATE").Where(keyCol+" = ?", key).First(redemption).Error
		if err != nil {
			return errors.New("无效的兑换码")
		}
		if redemption.Status != common.RedemptionCodeStatusEnabled {
			return errors.New("该兑换码已被使用")
		}
		err = tx.Model(&User{}).Where("id = ?", userId).Update("quota", gorm.Expr("quota + ?", redemption.Quota)).Error
		if err != nil {
			return err
		}
		redemption.RedeemedTime = common.GetTimestamp()
		redemption.Status = common.RedemptionCodeStatusUsed
		err = tx.Save(redemption).Error
		return err
	})
	if err != nil {
		return 0, errors.New("兑换失败，" + err.Error())
	}
	RecordLog(userId, LogTypeTopup, fmt.Sprintf("通过兑换码充值 %s", common.LogQuota(redemption.Quota)))
	return redemption.Quota, nil
}

type RechargeResponse struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg,omitempty"`
	TradeNo   string `json:"trade_no"`
	PayUrl    string `json:"payurl,omitempty"`
	QRcode    string `json:"qrcode,omitempty"`
	UrlScheme string `json:"urlscheme,omitempty"`
}

func InitRecharge(quota int, payType string, userId int) (payUrl string, TradeNo string, err error) {
	if quota < 1 {
		RecordLog(userId, LogTypeTopup, "至少充值1$")
		return "", "", errors.New("至少充值1$")

	}
	if quota > 50 {
		RecordLog(userId, LogTypeTopup, "单次最多充值50$")
		return "", "", errors.New("单次最多充值50$")
	}
	if userId == 0 {
		RecordLog(userId, LogTypeTopup, "无效的 user id")
		return "", "", errors.New("无效的 user id")
	}
	if (payType != "wxpay") && (payType != "alipay") {
		RecordLog(userId, LogTypeTopup, "无效的支付方式")
		return "", "", errors.New("无效的支付方式")
	}
	RecordLog(userId, LogTypeTopup, fmt.Sprintf("创建充值请求，请求金额: %d$", quota))
	rand.Seed(time.Now().UnixNano())
	// 生成交易号：年月日时分秒+随机数
	transactionID := fmt.Sprintf("%s%d", time.Now().Format("20060102150405"), rand.Intn(999999-100000)+100000)
	initRechargeLog := RechargeLog{
		UserId:      userId,
		TradeNo:     transactionID,
		Status:      1,
		Name:        "VIP会员",
		Quota:       quota * int(common.QuotaPerUnit),
		CreatedTime: common.GetTimestamp(),
	}
	DB.Create(&initRechargeLog)
	paymentData := map[string]string{
		"pid":          os.Getenv("YI_PAY_PID"),
		"type":         payType,
		"out_trade_no": transactionID,
		"notify_url":   "https://platform.hustgpt.com/api/user/rechargenotify",
		"return_url":   "https://platform.hustgpt.com/panel/topup",
		"name":         "VIP会员",
		"money":        fmt.Sprintf("%.2f", float64(quota)*7.3),
		"clientip":     "192.168.1.1",
		//"param":       "", // 业务扩展参数，没有可以留空
		"sign":      "", // 这一行在签名生成后会被更新
		"sign_type": "MD5",
	}

	// 商户密钥KEY（这里用XXXX代替真实的KEY）
	key := os.Getenv("YI_PAY_KEY")

	// 生成签名
	signature := GenerateSignature(paymentData, key)
	paymentData["sign"] = signature

	data := url.Values{}
	for k, v := range paymentData {
		data.Set(k, v)
	}

	resp, err := http.PostForm("https://yi-pay.com/mapi.php", data)
	if err != nil {
		RecordLog(userId, LogTypeTopup, "支付通道异常1")
		return "", "", errors.New("支付通道异常1")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // 读取响应体
	if err != nil {
		// 处理错误
		RecordLog(userId, LogTypeTopup, "支付通道异常2")
		return "", "", errors.New("支付通道异常2")
	}
	respJson := RechargeResponse{}
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		RecordLog(userId, LogTypeTopup, "支付通道异常3")
		return "", "", errors.New("支付通道异常3")
	}

	if respJson.Code != 1 {
		RecordLog(userId, LogTypeTopup, "支付通道异常4: "+respJson.Msg)
		return "", "", errors.New("支付通道异常4: " + respJson.Msg)
	}
	if respJson.PayUrl != "" {
		return respJson.PayUrl, transactionID, nil
	} else if respJson.QRcode != "" {
		return respJson.QRcode, transactionID, nil
	}
	return "", "", errors.New("返回的是小程序支付链接")
}

func CompeleteRecharge(TradeNo string, userId int) error {
	keyCol := "`trade_no`"
	if common.UsingPostgreSQL {
		keyCol = `"trade_no"`
	}
	rechargeLog := &RechargeLog{}
	err := DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Set("gorm:query_option", "FOR UPDATE").Where(keyCol+" = ?", TradeNo).First(rechargeLog).Error
		if err != nil {
			RecordLog(userId, LogTypeTopup, "未创建订单！")
			return errors.New("未创建订单！")
		}
		if rechargeLog.Status != 1 {
			RecordLog(userId, LogTypeTopup, "订单已支付！")
			return errors.New("订单已支付！")
		}
		err = tx.Model(&User{}).Where("id = ?", userId).Update("quota", gorm.Expr("quota + ?", rechargeLog.Quota)).Error
		if err != nil {
			RecordLog(userId, LogTypeTopup, err.Error())
			return err
		}
		rechargeLog.RedeemedTime = common.GetTimestamp()
		rechargeLog.Status = 2
		err = tx.Save(rechargeLog).Error
		RecordLog(userId, LogTypeTopup, err.Error())
		return err
	})
	if err != nil {
		return err
	}
	RecordLog(userId, LogTypeTopup, fmt.Sprintf("完成充值 %stkens", common.LogQuota(rechargeLog.Quota)))
	return nil
}

//	{err = DB.Transaction(func(tx *gorm.DB) error {
//		err := tx.Set("gorm:query_option", "FOR UPDATE").Where(keyCol+" = ?", key).First(redemption).Error
//		if err != nil {
//			return errors.New("无效的兑换码")
//		}
//		if redemption.Status != common.RedemptionCodeStatusEnabled {
//			return errors.New("该兑换码已被使用")
//		}
//		err = tx.Model(&User{}).Where("id = ?", userId).Update("quota", gorm.Expr("quota + ?", redemption.Quota)).Error
//		if err != nil {
//			return err
//		}
//		redemption.RedeemedTime = common.GetTimestamp()
//		redemption.Status = common.RedemptionCodeStatusUsed
//		err = tx.Save(redemption).Error
//		return err
//	})
//	if err != nil {
//		return 0, errors.New("兑换失败，" + err.Error())
//	}
//	RecordLog(userId, LogTypeTopup, fmt.Sprintf("通过兑换码充值 %s", common.LogQuota(redemption.Quota)))
//	return redemption.Quota, nil
//}

func GenerateSignature(parameters map[string]string, key string) string {
	// 1. 将参数按照键名ASCII升序排序
	keys := make([]string, 0, len(parameters))
	for k := range parameters {
		if k != "sign" && k != "sign_type" {
			keys = append(keys, k)
		}
	}
	slices.Sort(keys)
	// 2. 将排序后的参数拼接成URL键值对的格式，不进行URL编码
	var concatenatedParams string
	for _, k := range keys {
		value := parameters[k]
		if value != "" {
			concatenatedParams += fmt.Sprintf("%s=%s&", k, value)
		}
	}
	concatenatedParams = strings.TrimRight(concatenatedParams, "&") // 删除最后一个"&"
	signStr := concatenatedParams + key                             // 拼接密钥KEY

	// 3. 使用MD5算法加密
	hasher := md5.New()
	hasher.Write([]byte(signStr))
	sign := hex.EncodeToString(hasher.Sum(nil)) // 结果为小写

	return sign
}

func (redemption *Redemption) Insert() error {
	var err error
	err = DB.Create(redemption).Error
	return err
}

func (redemption *Redemption) SelectUpdate() error {
	// This can update zero values
	return DB.Model(redemption).Select("redeemed_time", "status").Updates(redemption).Error
}

// Update Make sure your token's fields is completed, because this will update non-zero values
func (redemption *Redemption) Update() error {
	var err error
	err = DB.Model(redemption).Select("name", "status", "quota", "redeemed_time").Updates(redemption).Error
	return err
}

func (redemption *Redemption) Delete() error {
	var err error
	err = DB.Delete(redemption).Error
	return err
}

func DeleteRedemptionById(id int) (err error) {
	if id == 0 {
		return errors.New("id 为空！")
	}
	redemption := Redemption{Id: id}
	err = DB.Where(redemption).First(&redemption).Error
	if err != nil {
		return err
	}
	return redemption.Delete()
}

type RedemptionStatistics struct {
	Count  int64 `json:"count"`
	Quota  int64 `json:"quota"`
	Status int   `json:"status"`
}

func GetStatisticsRedemption() (redemptionStatistics []*RedemptionStatistics, err error) {
	err = DB.Model(&Redemption{}).Select("status", "count(*) as count", "sum(quota) as quota").Where("status != ?", 2).Group("status").Scan(&redemptionStatistics).Error
	return redemptionStatistics, err
}

type RedemptionStatisticsGroup struct {
	Date      string `json:"date"`
	Quota     int64  `json:"quota"`
	UserCount int64  `json:"user_count"`
}

func GetStatisticsRedemptionByPeriod(startTimestamp, endTimestamp int64) (redemptionStatistics []*RedemptionStatisticsGroup, err error) {
	groupSelect := getTimestampGroupsSelect("redeemed_time", "day", "date")

	err = DB.Raw(`
		SELECT `+groupSelect+`,
		sum(quota) as quota,
		count(distinct user_id) as user_count
		FROM redemptions
		WHERE status=3
		AND redeemed_time BETWEEN ? AND ?
		GROUP BY date
		ORDER BY date
	`, startTimestamp, endTimestamp).Scan(&redemptionStatistics).Error

	return redemptionStatistics, err
}
