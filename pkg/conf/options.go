package conf

// RsshubHost RSSHub host used by the app.
// The default https://rsshub.app actually cannot work with Telegram.
// You should deploy RSSHub on your own.
// If you need authentication, see the option as a URL and add them on querystring.
// We can process it properly.
var RsshubHost = env("RSSHUB_HOST", strPtr("https://rsshub.app"))

// TgChanName Telegram channel slug to be forwarded.
// If your channel link is https://t.me/myl7s, use `myl7s` here.
var TgChanName = env("TG_CHAN_NAME", nil)

// PostFilterOut Message containing the slug will be filtered.
// See also https://docs.rsshub.app/parameter.html#nei-rong-guo-lu .
var PostFilterOut = env("POST_FILTER_OUT", strPtr("#notwfwd"))

// PollInterval Polling time interval to request RSSHub.
// The default cache age of RSSHub when caching is enabled is 300s.
// So if you want smaller update interval for quicker syncing,
// you need to not only decrease the variable but also disable/reconfigure RSSHub caching.
// Unit: s.
var PollInterval = envInt("POLL_INTERVAL", intPtr(360))

// PollRange The time range of each polling request to RSSHub.
// The option is directly used as `filter_time` param for RSSHub.
// Unit: s.
var PollRange = envInt("POLL_RANGE", intPtr(360))

// DBPath SQLite database file path
var DBPath = env("DB_PATH", strPtr("/db/db.sqlite"))

// TwConsumerKey The following 4 is Twitter secrets. See `README.md` to know how to get them.
var TwConsumerKey = env("TW_CONSUMER_KEY", nil)
var TwConsumerSecret = env("TW_CONSUMER_SECRET", nil)
var TwTokenKey = env("TW_TOKEN_KEY", nil)
var TwTokenSecret = env("TW_TOKEN_SECRET", nil)

// TwTextSplitBackRate When searching at the end of body text for a possible split position, the range proportion.
// The result of dividing the int with 100 should be in 0-1.
var TwTextSplitBackRate = envInt("TW_TEXT_SPLIT_BACK_RATE", intPtr(90))

// TwTextSplitBackDisableRate When set, disable the searching of TwTextSplitBackRate
var TwTextSplitBackDisableRate = env("TW_TEXT_SPLIT_BACK_DISABLE_RATE", strPtr(""))

// TwTextSplitBackLen Like TwTextSplitBackRate, but in int length counted in UTF-16 encoding.
// The default 20 chars should be large enough to cross an English word and reach a space separating point.
var TwTextSplitBackLen = envInt("TW_TEXT_SPLIT_BACK_LEN", intPtr(20))

// TwTextSplitBackDisableLen Like TwTextSplitBackDisableRate but for TwTextSplitBackLen.
// When all of TwTextSplitBackDisable* is set, the text splitting will fall back to just choose the end point,
// which means no attempt to split at a maybe better position will be taken.
var TwTextSplitBackDisableLen = env("TW_TEXT_SPLIT_BACK_DISABLE_RATE", strPtr(""))

// TmpDir Temporary dir path.
// Make sure it exists and is accessible by the app.
// Multiply directories may be created as tmp dir for different purposes in it.
var TmpDir = env("TMP_DIR", strPtr("/tmp"))
