package relay

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"one-api/common"
	"one-api/model"
	"one-api/types"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type QuotaInfo struct {
	modelName         string
	promptTokens      int
	preConsumedTokens int
	modelRatio        []float64
	groupRatio        float64
	ratio             float64
	preConsumedQuota  int
	userId            int
	channelId         int
	tokenId           int
	HandelStatus      bool
}

func generateQuotaInfo(c *gin.Context, modelName string, promptTokens int) (*QuotaInfo, *types.OpenAIErrorWithStatusCode) {
	quotaInfo := &QuotaInfo{
		modelName:    modelName,
		promptTokens: promptTokens,
		userId:       c.GetInt("id"),
		channelId:    c.GetInt("channel_id"),
		tokenId:      c.GetInt("token_id"),
		HandelStatus: false,
	}
	quotaInfo.initQuotaInfo(c.GetString("group"))

	errWithCode := quotaInfo.preQuotaConsumption()
	if errWithCode != nil {
		return nil, errWithCode
	}

	return quotaInfo, nil
}

func (q *QuotaInfo) initQuotaInfo(groupName string) {
	modelRatio := common.GetModelRatio(q.modelName)
	groupRatio := common.GetGroupRatio(groupName)
	preConsumedTokens := common.PreConsumedQuota
	ratio := modelRatio[0] * groupRatio
	preConsumedQuota := int(float64(q.promptTokens+preConsumedTokens) * ratio)

	q.preConsumedTokens = preConsumedTokens
	q.modelRatio = modelRatio
	q.groupRatio = groupRatio
	q.ratio = ratio
	q.preConsumedQuota = preConsumedQuota

}

func (q *QuotaInfo) preQuotaConsumption() *types.OpenAIErrorWithStatusCode {
	userQuota, err := model.CacheGetUserQuota(q.userId)
	if err != nil {
		return common.ErrorWrapper(err, "get_user_quota_failed", http.StatusInternalServerError)
	}

	if userQuota < q.preConsumedQuota && q.modelRatio[0] != 0 {
		return common.ErrorWrapper(errors.New("user quota is not enough"), "insufficient_user_quota", http.StatusForbidden)
	}

	token, err := model.GetTokenById(q.tokenId)
	if err != nil {
		return common.ErrorWrapper(err, "get_token_failed", http.StatusInternalServerError)
	}
	if token.RemainQuota < q.preConsumedQuota && q.modelRatio[0] != 0 && !token.UnlimitedQuota {
		return common.ErrorWrapper(errors.New("token is not enough"), "insufficient_token_quota", http.StatusForbidden)
	}

	if q.modelRatio[0] == 0 {
		times, err := model.GetTimesByUserId(q.userId)
		if err != nil {
			return common.ErrorWrapper(err, "get_Times_failed", http.StatusInternalServerError)
		}
		t := time.Unix(times.ChangeTime, 0)

		// 创建中国时区的 Location
		loc, _ := time.LoadLocation("Asia/Shanghai")

		// 将 t 转换为中国时区的时间
		t = t.In(loc)

		// 获取今天的日期
		now := time.Now().In(loc)
		year, month, day := now.Date()

		// 比较时间是否在今天
		tYear, tMonth, tDay := t.Date()
		if !(year == tYear && month == tMonth && day == tDay) {
			times.Times = 30
		} else {
			if times.Times <= 0 {
				return common.ErrorWrapper(err, "free time is used out!", http.StatusForbidden)
			}
		}
		times.ChangeTime = common.GetTimestamp()
		err = times.Update()
		if err != nil {
			return common.ErrorWrapper(err, "update_Times_failed", http.StatusInternalServerError)
		}
	}

	err = model.CacheDecreaseUserQuota(q.userId, q.preConsumedQuota)
	if err != nil {
		return common.ErrorWrapper(err, "decrease_user_quota_failed", http.StatusInternalServerError)
	}

	if userQuota > 100*q.preConsumedQuota {
		// in this case, we do not pre-consume quota
		// because the user has enough quota
		q.preConsumedQuota = 0
		// common.LogInfo(c.Request.Context(), fmt.Sprintf("user %d has enough quota %d, trusted and no need to pre-consume", userId, userQuota))
	}

	if q.preConsumedQuota > 0 {
		err := model.PreConsumeTokenQuota(q.tokenId, q.preConsumedQuota)
		if err != nil {
			return common.ErrorWrapper(err, "pre_consume_token_quota_failed", http.StatusForbidden)
		}
		q.HandelStatus = true
	}

	return nil
}

func (q *QuotaInfo) completedQuotaConsumption(usage *types.Usage, tokenName string, ctx context.Context) error {
	quota := 0
	completionRatio := q.modelRatio[1] * q.groupRatio
	promptTokens := usage.PromptTokens
	completionTokens := usage.CompletionTokens
	quota = int(math.Ceil(((float64(promptTokens) * q.ratio) + (float64(completionTokens) * completionRatio))))
	if q.ratio != 0 && quota <= 0 {
		quota = 1
	}
	totalTokens := promptTokens + completionTokens
	if totalTokens == 0 {
		// in this case, must be some error happened
		// we cannot just return, because we may have to return the pre-consumed quota
		quota = 0
	}
	quotaDelta := quota - q.preConsumedQuota
	err := model.PostConsumeTokenQuota(q.tokenId, quotaDelta)
	if err != nil {
		return errors.New("error consuming token remain quota: " + err.Error())
	}
	err = model.CacheUpdateUserQuota(q.userId)
	if err != nil {
		return errors.New("error consuming token remain quota: " + err.Error())
	}
	if quota >= 0 {
		requestTime := 0
		requestStartTimeValue := ctx.Value("requestStartTime")
		if requestStartTimeValue != nil {
			requestStartTime, ok := requestStartTimeValue.(time.Time)
			if ok {
				requestTime = int(time.Since(requestStartTime).Milliseconds())
			}
		}
		var modelRatioStr string
		if q.modelRatio[0] == q.modelRatio[1] {
			modelRatioStr = fmt.Sprintf("%.2f", q.modelRatio[0])
		} else {
			modelRatioStr = fmt.Sprintf("%.2f (输入)/%.2f (输出)", q.modelRatio[0], q.modelRatio[1])
		}

		logContent := fmt.Sprintf("模型倍率 %s", modelRatioStr)
		model.RecordConsumeLog(ctx, q.userId, q.channelId, promptTokens, completionTokens, q.modelName, tokenName, quota, logContent, requestTime)
		model.UpdateUserUsedQuotaAndRequestCount(q.userId, quota)
		model.UpdateChannelUsedQuota(q.channelId, quota)

		times, err := model.GetTimesByUserId(q.userId)
		if err != nil {
			return errors.New("error get times: " + err.Error())
		}
		if os.Getenv("ENABLE_LIMIT_FREE") == "true" {
			times.Times = times.Times - 1
		}

		err = times.Update()
		if err != nil {
			return errors.New("error update times: " + err.Error())
		}
	}

	return nil
}

func (q *QuotaInfo) undo(c *gin.Context) {
	tokenId := c.GetInt("token_id")
	if q.HandelStatus {
		go func(ctx context.Context) {
			// return pre-consumed quota
			err := model.PostConsumeTokenQuota(tokenId, -q.preConsumedQuota)
			if err != nil {
				common.LogError(ctx, "error return pre-consumed quota: "+err.Error())
			}
		}(c.Request.Context())
	}
}

func (q *QuotaInfo) consume(c *gin.Context, usage *types.Usage) {
	tokenName := c.GetString("token_name")
	// 如果没有报错，则消费配额
	go func(ctx context.Context) {
		err := q.completedQuotaConsumption(usage, tokenName, ctx)
		if err != nil {
			common.LogError(ctx, err.Error())
		}
	}(c.Request.Context())
}
