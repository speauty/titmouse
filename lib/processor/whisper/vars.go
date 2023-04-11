package whisper

var (
	LanguagesSupported = []string{
		"英文", "中文", "德语", "西班牙语", "俄罗斯语", "韩语", "法语", "日语", "葡萄牙语",
	}
	languagesMap = map[string]string{
		"英文": "en", "中文": "zh", "德语": "de", "西班牙语": "es", "俄罗斯语": "ru",
		"韩语": "ko", "法语": "fr", "日语": "ja", "葡萄牙语": "pt",
	}
)
