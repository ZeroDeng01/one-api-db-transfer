package main

// 新版通道类别
const (
	ChannelTypeUnknown        = 0
	ChannelTypeOpenAI         = 1
	ChannelTypeAPI2D          = 2
	ChannelTypeAzure          = 3
	ChannelTypeCloseAI        = 4
	ChannelTypeOpenAISB       = 5
	ChannelTypeOpenAIMax      = 6
	ChannelTypeOhMyGPT        = 7
	ChannelTypeCustom         = 8
	ChannelTypeAILS           = 9
	ChannelTypeAIProxy        = 10
	ChannelTypePaLM           = 11
	ChannelTypeAPI2GPT        = 12
	ChannelTypeAIGC2D         = 13
	ChannelTypeAnthropic      = 14
	ChannelTypeBaidu          = 15
	ChannelTypeZhipu          = 16
	ChannelTypeAli            = 17
	ChannelTypeXunfei         = 18
	ChannelType360            = 19
	ChannelTypeOpenRouter     = 20
	ChannelTypeAIProxyLibrary = 21
	ChannelTypeFastGPT        = 22
	ChannelTypeTencent        = 23
	ChannelTypeAzureSpeech    = 24
	ChannelTypeGemini         = 25
	ChannelTypeBaichuan       = 26
	ChannelTypeMiniMax        = 27
	ChannelTypeDeepseek       = 28
	ChannelTypeMoonshot       = 29
	ChannelTypeMistral        = 30
	ChannelTypeGroq           = 31
	ChannelTypeBedrock        = 32
	ChannelTypeLingyi         = 33
	ChannelTypeMidjourney     = 34
	ChannelTypeCloudflareAI   = 35
	ChannelTypeCohere         = 36
	ChannelTypeStabilityAI    = 37
	ChannelTypeCoze           = 38
	ChannelTypeOllama         = 39
	ChannelTypeHunyuan        = 40
)

// 旧通道类别
const (
	Unknown = iota
	OpenAI
	API2D
	Azure
	CloseAI
	OpenAISB
	OpenAIMax
	OhMyGPT
	Custom
	Ails
	AIProxy
	PaLM
	API2GPT
	AIGC2D
	Anthropic
	Baidu
	Zhipu
	Ali
	Xunfei
	AI360
	OpenRouter
	AIProxyLibrary
	FastGPT
	Tencent
	Gemini
	Moonshot
	Baichuan
	Minimax
	Mistral
	Groq
	Ollama
	LingYiWanWu
	StepFun
	AwsClaude
	Coze
	Cohere
	DeepSeek
	Cloudflare
	DeepL
	TogetherAI
	Doubao
	Dummy
)

// 旧通道类别映射新通道类别假定映射
var channelOldToNew = map[int]int{
	Unknown:        ChannelTypeUnknown,
	OpenAI:         ChannelTypeOpenAI,
	API2D:          ChannelTypeAPI2D,
	Azure:          ChannelTypeAzure,
	CloseAI:        ChannelTypeCloseAI,
	OpenAISB:       ChannelTypeOpenAISB,
	OpenAIMax:      ChannelTypeOpenAIMax,
	OhMyGPT:        ChannelTypeOhMyGPT,
	Custom:         ChannelTypeCustom,
	Ails:           ChannelTypeAILS,
	AIProxy:        ChannelTypeAIProxy,
	PaLM:           ChannelTypePaLM,
	API2GPT:        ChannelTypeAPI2GPT,
	AIGC2D:         ChannelTypeAIGC2D,
	Anthropic:      ChannelTypeAnthropic,
	Baidu:          ChannelTypeBaidu,
	Zhipu:          ChannelTypeZhipu,
	Ali:            ChannelTypeAli,
	Xunfei:         ChannelTypeXunfei,
	AI360:          ChannelType360,
	OpenRouter:     ChannelTypeOpenRouter,
	AIProxyLibrary: ChannelTypeAIProxyLibrary,
	FastGPT:        ChannelTypeFastGPT,
	Tencent:        ChannelTypeTencent,
	Gemini:         ChannelTypeGemini,
	Moonshot:       ChannelTypeMoonshot,
	Baichuan:       ChannelTypeBaichuan,
	Minimax:        ChannelTypeMiniMax,
	Mistral:        ChannelTypeMistral,
	Groq:           ChannelTypeGroq,
	Ollama:         ChannelTypeOllama,
	LingYiWanWu:    ChannelTypeLingyi,
	StepFun:        ChannelTypeMidjourney,
	AwsClaude:      ChannelTypeCloudflareAI,
	Coze:           ChannelTypeCoze,
	Cohere:         ChannelTypeCohere,
	DeepSeek:       ChannelTypeDeepseek,
	Cloudflare:     ChannelTypeCloudflareAI,
	DeepL:          ChannelTypeStabilityAI,
	TogetherAI:     ChannelTypeCoze,
	Doubao:         ChannelTypeOllama,
	Dummy:          ChannelTypeHunyuan,
}
