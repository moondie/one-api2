package common

import (
	"encoding/json"
	"strings"
	"time"
)

type ModelType struct {
	Ratio []float64
	Type  int
}

var ModelTypes map[string]ModelType

// ModelRatio
// https://platform.openai.com/docs/models/model-endpoint-compatibility
// https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Blfmc9dlf
// https://openai.com/pricing
// TODO: when a new api is enabled, check the pricing here
// 1 === $0.002 / 1K tokens
// 1 === ￥0.014 / 1k tokens
var ModelRatio map[string][]float64

func init() {
	ModelTypes = map[string]ModelType{
		// 	$0.03 / 1K tokens	$0.06 / 1K tokens
		"gpt-4":      {[]float64{15, 30}, ChannelTypeOpenAI},
		"gpt-4-0314": {[]float64{15, 30}, ChannelTypeOpenAI},
		"gpt-4-0613": {[]float64{15, 30}, ChannelTypeOpenAI},
		// 	$0.06 / 1K tokens	$0.12 / 1K tokens
		"gpt-4-32k":      {[]float64{30, 60}, ChannelTypeOpenAI},
		"gpt-4-32k-0314": {[]float64{30, 60}, ChannelTypeOpenAI},
		"gpt-4-32k-0613": {[]float64{30, 60}, ChannelTypeOpenAI},
		// 	$0.01 / 1K tokens	$0.03 / 1K tokens
		"gpt-4-preview":        {[]float64{5, 15}, ChannelTypeOpenAI},
		"gpt-4-1106-preview":   {[]float64{5, 15}, ChannelTypeOpenAI},
		"gpt-4-0125-preview":   {[]float64{5, 15}, ChannelTypeOpenAI},
		"gpt-4-turbo-preview":  {[]float64{5, 15}, ChannelTypeOpenAI},
		"gpt-4-vision-preview": {[]float64{5, 15}, ChannelTypeOpenAI},
		// 	$0.0005 / 1K tokens	$0.0015 / 1K tokens
		"gpt-3.5-turbo":      {[]float64{0.25, 0.75}, ChannelTypeOpenAI},
		"gpt-3.5-turbo-0125": {[]float64{0.25, 0.75}, ChannelTypeOpenAI},
		// 	$0.0015 / 1K tokens	$0.002 / 1K tokens
		"gpt-3.5-turbo-0301":     {[]float64{0.75, 1}, ChannelTypeOpenAI},
		"gpt-3.5-turbo-0613":     {[]float64{0.75, 1}, ChannelTypeOpenAI},
		"gpt-3.5-turbo-instruct": {[]float64{0.75, 1}, ChannelTypeOpenAI},
		// 	$0.003 / 1K tokens	$0.004 / 1K tokens
		"gpt-3.5-turbo-16k":      {[]float64{1.5, 2}, ChannelTypeOpenAI},
		"gpt-3.5-turbo-16k-0613": {[]float64{1.5, 2}, ChannelTypeOpenAI},
		// 	$0.001 / 1K tokens	$0.002 / 1K tokens
		"gpt-3.5-turbo-1106": {[]float64{0.5, 1}, ChannelTypeOpenAI},
		// 	$0.0020 / 1K tokens
		"davinci-002": {[]float64{1, 1}, ChannelTypeOpenAI},
		// 	$0.0004 / 1K tokens
		"babbage-002":           {[]float64{0.2, 0.2}, ChannelTypeOpenAI},
		"text-ada-001":          {[]float64{0.2, 0.2}, ChannelTypeOpenAI},
		"text-babbage-001":      {[]float64{0.25, 0.25}, ChannelTypeOpenAI},
		"text-curie-001":        {[]float64{1, 1}, ChannelTypeOpenAI},
		"text-davinci-002":      {[]float64{10, 10}, ChannelTypeOpenAI},
		"text-davinci-003":      {[]float64{10, 10}, ChannelTypeOpenAI},
		"text-davinci-edit-001": {[]float64{10, 10}, ChannelTypeOpenAI},
		"code-davinci-edit-001": {[]float64{10, 10}, ChannelTypeOpenAI},
		// $0.006 / minute -> $0.006 / 150 words -> $0.006 / 200 tokens -> $0.03 / 1k tokens
		"whisper-1": {[]float64{15, 15}, ChannelTypeOpenAI},
		// $0.015 / 1K characters
		"tts-1":      {[]float64{7.5, 7.5}, ChannelTypeOpenAI},
		"tts-1-1106": {[]float64{7.5, 7.5}, ChannelTypeOpenAI},
		// $0.030 / 1K characters
		"tts-1-hd":               {[]float64{15, 15}, ChannelTypeOpenAI},
		"tts-1-hd-1106":          {[]float64{15, 15}, ChannelTypeOpenAI},
		"davinci":                {[]float64{10, 10}, ChannelTypeOpenAI},
		"curie":                  {[]float64{10, 10}, ChannelTypeOpenAI},
		"babbage":                {[]float64{10, 10}, ChannelTypeOpenAI},
		"ada":                    {[]float64{10, 10}, ChannelTypeOpenAI},
		"text-embedding-ada-002": {[]float64{0.05, 0.05}, ChannelTypeOpenAI},
		// 	$0.00002 / 1K tokens
		"text-embedding-3-small": {[]float64{0.01, 0.01}, ChannelTypeOpenAI},
		// 	$0.00013 / 1K tokens
		"text-embedding-3-large":  {[]float64{0.065, 0.065}, ChannelTypeOpenAI},
		"text-search-ada-doc-001": {[]float64{10, 10}, ChannelTypeOpenAI},
		"text-moderation-stable":  {[]float64{0.1, 0.1}, ChannelTypeOpenAI},
		"text-moderation-latest":  {[]float64{0.1, 0.1}, ChannelTypeOpenAI},
		// $0.016 - $0.020 / image
		"dall-e-2": {[]float64{8, 8}, ChannelTypeOpenAI},
		// $0.040 - $0.120 / image
		"dall-e-3": {[]float64{20, 20}, ChannelTypeOpenAI},

		// $0.80/million tokens $2.40/million tokens
		"claude-instant-1": {[]float64{0.4, 1.2}, ChannelTypeAnthropic},
		// $8.00/million tokens $24.00/million tokens
		"claude-2":   {[]float64{4, 12}, ChannelTypeAnthropic},
		"claude-2.0": {[]float64{4, 12}, ChannelTypeAnthropic},
		"claude-2.1": {[]float64{4, 12}, ChannelTypeAnthropic},

		// ￥0.012 / 1k tokens ￥0.012 / 1k tokens
		"ERNIE-Bot": {[]float64{0.8572, 0.8572}, ChannelTypeBaidu},
		// 0.024元/千tokens 0.048元/千tokens
		"ERNIE-Bot-8k": {[]float64{1.7143, 3.4286}, ChannelTypeBaidu},
		// ￥0.008 / 1k tokens ￥0.008 / 1k tokens
		"ERNIE-Bot-turbo": {[]float64{0.5715, 0.5715}, ChannelTypeBaidu},
		// ￥0.12 / 1k tokens ￥0.12 / 1k tokens
		"ERNIE-Bot-4": {[]float64{8.572, 8.572}, ChannelTypeBaidu},
		// ￥0.002 / 1k tokens
		"Embedding-V1": {[]float64{0.1429, 0.1429}, ChannelTypeBaidu},

		"PaLM-2":            {[]float64{1, 1}, ChannelTypePaLM},
		"gemini-pro":        {[]float64{1, 1}, ChannelTypeGemini},
		"gemini-pro-vision": {[]float64{1, 1}, ChannelTypeGemini},

		// ￥0.005 / 1k tokens
		"chatglm_turbo": {[]float64{0.3572, 0.3572}, ChannelTypeZhipu},
		"chatglm_std":   {[]float64{0.3572, 0.3572}, ChannelTypeZhipu},
		"glm-3-turbo":   {[]float64{0.3572, 0.3572}, ChannelTypeZhipu},
		// ￥0.01 / 1k tokens
		"chatglm_pro": {[]float64{0.7143, 0.7143}, ChannelTypeZhipu},
		// ￥0.002 / 1k tokens
		"chatglm_lite": {[]float64{0.1429, 0.1429}, ChannelTypeZhipu},
		// ￥0.1 / 1k tokens
		"glm-4":  {[]float64{7.143, 7.143}, ChannelTypeZhipu},
		"glm-4v": {[]float64{7.143, 7.143}, ChannelTypeZhipu},
		// ￥0.0005 / 1k tokens
		"embedding-2": {[]float64{0.0357, 0.0357}, ChannelTypeZhipu},
		// ￥0.25 / 1张图片
		"cogview-3": {[]float64{17.8571, 17.8571}, ChannelTypeZhipu},

		// ￥0.008 / 1k tokens
		"qwen-turbo": {[]float64{0.5715, 0.5715}, ChannelTypeAli},
		// ￥0.02 / 1k tokens
		"qwen-plus":            {[]float64{1.4286, 1.4286}, ChannelTypeAli},
		"qwen-max":             {[]float64{1.4286, 1.4286}, ChannelTypeAli},
		"qwen-max-longcontext": {[]float64{1.4286, 1.4286}, ChannelTypeAli},
		"qwen-vl-plus":         {[]float64{0.5715, 0.5715}, ChannelTypeAli},
		"qwen-vl-max":          {[]float64{0.5715, 0.5715}, ChannelTypeAli},
		// ￥0.0007 / 1k tokens
		"text-embedding-v1": {[]float64{0.05, 0.05}, ChannelTypeAli},

		// ￥0.018 / 1k tokens
		"SparkDesk":      {[]float64{1.2858, 1.2858}, ChannelTypeXunfei},
		"SparkDesk-v1.1": {[]float64{1.2858, 1.2858}, ChannelTypeXunfei},
		"SparkDesk-v2.1": {[]float64{1.2858, 1.2858}, ChannelTypeXunfei},
		"SparkDesk-v3.1": {[]float64{1.2858, 1.2858}, ChannelTypeXunfei},
		"SparkDesk-v3.5": {[]float64{1.2858, 1.2858}, ChannelTypeXunfei},

		// ¥0.012 / 1k tokens
		"360GPT_S2_V9": {[]float64{0.8572, 0.8572}, ChannelType360},
		// ¥0.001 / 1k tokens
		"embedding-bert-512-v1":     {[]float64{0.0715, 0.0715}, ChannelType360},
		"embedding_s1_v1":           {[]float64{0.0715, 0.0715}, ChannelType360},
		"semantic_similarity_s1_v1": {[]float64{0.0715, 0.0715}, ChannelType360},

		// ¥0.1 / 1k tokens  // https://cloud.tencent.com/document/product/1729/97731#e0e6be58-60c8-469f-bdeb-6c264ce3b4d0
		"hunyuan": {[]float64{7.143, 7.143}, ChannelTypeTencent},

		"Baichuan2-Turbo":         {[]float64{0.5715, 0.5715}, ChannelTypeBaichuan}, // ¥0.008 / 1k tokens
		"Baichuan2-Turbo-192k":    {[]float64{1.143, 1.143}, ChannelTypeBaichuan},   // ¥0.016 / 1k tokens
		"Baichuan2-53B":           {[]float64{1.4286, 1.4286}, ChannelTypeBaichuan}, // ¥0.02 / 1k tokens
		"Baichuan-Text-Embedding": {[]float64{0.0357, 0.0357}, ChannelTypeBaichuan}, // ¥0.0005 / 1k tokens

		"abab5.5s-chat": {[]float64{0.3572, 0.3572}, ChannelTypeMiniMax},   // ¥0.005 / 1k tokens
		"abab5.5-chat":  {[]float64{1.0714, 1.0714}, ChannelTypeMiniMax},   // ¥0.015 / 1k tokens
		"abab6-chat":    {[]float64{14.2857, 14.2857}, ChannelTypeMiniMax}, // ¥0.2 / 1k tokens
		"embo-01":       {[]float64{0.0357, 0.0357}, ChannelTypeMiniMax},   // ¥0.0005 / 1k tokens

		"deepseek-coder": {[]float64{0.75, 0.75}, ChannelTypeDeepseek}, // 暂定 $0.0015 / 1K tokens
		"deepseek-chat":  {[]float64{0.75, 0.75}, ChannelTypeDeepseek}, // 暂定 $0.0015 / 1K tokens

		"moonshot-v1-8k":   {[]float64{0.8572, 0.8572}, ChannelTypeMoonshot}, // ¥0.012 / 1K tokens
		"moonshot-v1-32k":  {[]float64{1.7143, 1.7143}, ChannelTypeMoonshot}, // ¥0.024 / 1K tokens
		"moonshot-v1-128k": {[]float64{4.2857, 4.2857}, ChannelTypeMoonshot}, // ¥0.06 / 1K tokens
	}

	ModelRatio = make(map[string][]float64)
	for name, modelType := range ModelTypes {
		ModelRatio[name] = modelType.Ratio
	}
}

var DalleSizeRatios = map[string]map[string]float64{
	"dall-e-2": {
		"256x256":   1,
		"512x512":   1.125,
		"1024x1024": 1.25,
	},
	"dall-e-3": {
		"1024x1024": 1,
		"1024x1792": 2,
		"1792x1024": 2,
	},
}

var DalleGenerationImageAmounts = map[string][2]int{
	"dall-e-2": {1, 10},
	"dall-e-3": {1, 1}, // OpenAI allows n=1 currently.
}

var DalleImagePromptLengthLimitations = map[string]int{
	"dall-e-2": 1000,
	"dall-e-3": 4000,
}

func ModelRatio2JSONString() string {
	jsonBytes, err := json.Marshal(ModelRatio)
	if err != nil {
		SysError("error marshalling model ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateModelRatioByJSONString(jsonStr string) error {
	ModelRatio = make(map[string][]float64)
	return json.Unmarshal([]byte(jsonStr), &ModelRatio)
}

func MergeModelRatioByJSONString(jsonStr string) (newJsonStr string, err error) {
	isNew := false
	inputModelRatio := make(map[string][]float64)
	err = json.Unmarshal([]byte(jsonStr), &inputModelRatio)
	if err != nil {
		inputModelRatioOld := make(map[string]float64)
		err = json.Unmarshal([]byte(jsonStr), &inputModelRatioOld)
		if err != nil {
			return
		}

		inputModelRatio = UpdateModeRatioFormat(inputModelRatioOld)
		isNew = true
	}

	// 与现有的ModelRatio进行比较，如果有新增的模型，需要添加
	for key, value := range ModelRatio {
		if _, ok := inputModelRatio[key]; !ok {
			isNew = true
			inputModelRatio[key] = value
		}
	}

	if !isNew {
		return
	}

	var jsonBytes []byte
	jsonBytes, err = json.Marshal(inputModelRatio)
	if err != nil {
		SysError("error marshalling model ratio: " + err.Error())
	}
	newJsonStr = string(jsonBytes)
	return
}

func UpdateModeRatioFormat(modelRatioOld map[string]float64) map[string][]float64 {
	modelRatioNew := make(map[string][]float64)
	for key, value := range modelRatioOld {
		completionRatio := GetCompletionRatio(key) * value
		modelRatioNew[key] = []float64{value, completionRatio}
	}
	return modelRatioNew
}

func GetModelRatio(name string) []float64 {
	if strings.HasPrefix(name, "qwen-") && strings.HasSuffix(name, "-internet") {
		name = strings.TrimSuffix(name, "-internet")
	}
	ratio, ok := ModelRatio[name]
	if !ok {
		SysError("model ratio not found: " + name)
		return []float64{30, 30}
	}
	return ratio
}

func GetCompletionRatio(name string) float64 {
	if strings.HasPrefix(name, "gpt-3.5") {
		if strings.HasSuffix(name, "1106") {
			return 2
		}
		if name == "gpt-3.5-turbo" || name == "gpt-3.5-turbo-16k" {
			// TODO: clear this after 2023-12-11
			now := time.Now()
			// https://platform.openai.com/docs/models/continuous-model-upgrades
			// if after 2023-12-11, use 2
			if now.After(time.Date(2023, 12, 11, 0, 0, 0, 0, time.UTC)) {
				return 2
			}
		}
		return 1.333333
	}
	if strings.HasPrefix(name, "gpt-4") {
		if strings.HasSuffix(name, "preview") {
			return 3
		}
		return 2
	}
	if strings.HasPrefix(name, "claude-instant-1") {
		return 3.38
	}
	if strings.HasPrefix(name, "claude-2") {
		return 2.965517
	}
	return 1
}
