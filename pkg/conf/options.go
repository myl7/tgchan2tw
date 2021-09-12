package conf

// RsshubHost Actually https://rsshub.app cannot work with Telegram.
// Maybe you have to deploy RSSHub on your own.
// If you need authentication, see the option as a URL and add them on querystring.
// We can process it properly.
var RsshubHost = env("RSSHUB_HOST", strPtr("https://rsshub.app"))

var TgChanName = env("TG_CHAN_NAME", nil)

var PostFilterOut = env("POST_FILTER_OUT", strPtr("#notwfwd"))

// PollInterval The default cache age of RSSHub when caching is enabled is 300s.
// So if you want smaller update interval for quick syncing,
// you need to not only decrease the value but also disable/reconfigure RSSHub caching.
// Unit: s.
var PollInterval = envInt("POLL_INTERVAL", intPtr(360))

var DBPath = env("DB_PATH", strPtr("/db/db.sqlite"))

var TwConsumerKey = env("TW_CONSUMER_KEY", nil)
var TwConsumerSecret = env("TW_CONSUMER_SECRET", nil)
var TwTokenKey = env("TW_TOKEN_KEY", nil)
var TwTokenSecret = env("TW_TOKEN_SECRET", nil)
