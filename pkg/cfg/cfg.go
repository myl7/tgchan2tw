// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package cfg

type Config struct {
	// RsshubUrl RSSHub url used by the app.
	// The default https://rsshub.app actually cannot work with Telegram.
	// You should deploy RSSHub on your own.
	// If you need authentication, see the option as a URL and add them on querystring.
	// We can process it properly.
	RsshubUrl string

	// TgChanName Telegram channel slug to be forwarded.
	// If your channel link is https://t.me/myl7s, use `myl7s` here.
	TgChanName string

	// PostFilterOut Message containing the slug will be filtered.
	// See also https://docs.rsshub.app/parameter.html#nei-rong-guo-lu .
	// You may use | to separate multiply keywords.
	PostFilterOut string

	// PollInterval Polling time interval to request RSSHub.
	// The default cache age of RSSHub when caching is enabled is 300s.
	// So if you want smaller update interval for quicker syncing,
	// you need to not only decrease the variable but also disable/reconfigure RSSHub caching.
	// Unit: s.
	PollInterval int

	// PollRange The time range of each polling request to RSSHub.
	// The option is directly used as `filter_time` param for RSSHub.
	// PollRange / PollInterval indicates how many times will a message be tried to forward.
	// If your RSSHub is fine enough, set PollRange = 1 * PollInterval, which means only one try.
	// Otherwise, PollRange = 2 * PollInterval or 3 * PollInterval should be better.
	// Unit: s.
	PollRange int

	// DBPath SQLite database file path
	DBPath string

	// TwConsumerKey The following 4 is Twitter secrets. See `README.md` to know how to get them.
	TwConsumerKey    string
	TwConsumerSecret string
	TwTokenKey       string
	TwTokenSecret    string

	// TwTextSplitBackRate When searching at the end of body text for a possible split position, the range proportion.
	// The result of dividing the int with 100 should be in 0-1.
	TwTextSplitBackRate int

	// TwTextSplitBackDisableRate When set, disable the searching of TwTextSplitBackRate
	TwTextSplitBackDisableRate bool

	// TwTextSplitBackLen Like TwTextSplitBackRate, but in int length counted in UTF-16 encoding.
	// The default 20 chars should be large enough to cross an English word and reach a space separating point.
	TwTextSplitBackLen int

	// TwTextSplitBackDisableLen Like TwTextSplitBackDisableRate but for TwTextSplitBackLen.
	// When all of TwTextSplitBackDisable* is set, the text splitting will fall back to just choose the end point,
	// which means no attempt to split at a maybe better position will be taken.
	TwTextSplitBackDisableLen bool

	// TmpDir Temporary dir path.
	// Make sure it exists and is accessible by the app.
	// Multiply directories may be created as tmp dir for different purposes in it.
	TmpDir string
}

func GetConfig() (*Config, error) {
	var c Config
	var err error

	c.RsshubUrl, err = getEnvStrDef("APP_RSSHUB_URL", "https://rsshub.app")
	if err != nil {
		return nil, err
	}

	c.TgChanName, err = getEnvStr("APP_TG_CHAN_NAME")
	if err != nil {
		return nil, err
	}

	c.PostFilterOut, err = getEnvStrDef("APP_POST_FILTER_OUT", "#notwfwd")
	if err != nil {
		return nil, err
	}

	c.PollInterval, err = getEnvIntDef("APP_POLL_INTERVAL", 360)
	if err != nil {
		return nil, err
	}

	c.PollRange, err = getEnvIntDef("APP_POLL_RANGE", 720)
	if err != nil {
		return nil, err
	}

	c.DBPath, err = getEnvStrDef("APP_DB_PATH", "/db/db.sqlite")
	if err != nil {
		return nil, err
	}

	c.TwConsumerKey, err = getEnvStr("APP_TW_CONSUMER_KEY")
	if err != nil {
		return nil, err
	}

	c.TwConsumerSecret, err = getEnvStr("APP_TW_CONSUMER_SECRET")
	if err != nil {
		return nil, err
	}

	c.TwTokenKey, err = getEnvStr("APP_TW_TOKEN_KEY")
	if err != nil {
		return nil, err
	}

	c.TwTokenSecret, err = getEnvStr("APP_TW_TOKEN_SECRET")
	if err != nil {
		return nil, err
	}

	c.TwTextSplitBackRate, err = getEnvIntDef("APP_TW_TEXT_SPLIT_BACK_RATE", 90)
	if err != nil {
		return nil, err
	}

	twTextSplitBackDisableRate, err := getEnvStrDef("APP_TW_TEXT_SPLIT_BACK_DISABLE_RATE", "")
	if err != nil {
		return nil, err
	}
	c.TwTextSplitBackDisableRate = twTextSplitBackDisableRate != ""

	c.TwTextSplitBackLen, err = getEnvIntDef("APP_TW_TEXT_SPLIT_BACK_LEN", 20)
	if err != nil {
		return nil, err
	}

	twTextSplitBackDisableLen, err := getEnvStrDef("TW_TEXT_SPLIT_BACK_DISABLE_RATE", "")
	if err != nil {
		return nil, err
	}
	c.TwTextSplitBackDisableLen = twTextSplitBackDisableLen != ""

	c.TmpDir, err = getEnvStrDef("TMP_DIR", "/tmp")
	if err != nil {
		return nil, err
	}

	return &c, nil
}

var Cfg *Config

func LoadConfig() error {
	if Cfg != nil {
		return nil
	}

	var err error
	Cfg, err = GetConfig()
	if err != nil {
		return err
	}
	return nil
}
