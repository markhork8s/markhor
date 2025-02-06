package pkg

// Version of this package (Markhor)
const VERSION string = "1.0.5"

// In a MarkhorSecret, this key contains the Markhor parameters
const MARKHORPARAMS_MANIFEST_KEY string = "markhorParams"

// Path where to look for the configuration file by default if none is provided
// in the CLI args
const DEFAULT_CONFIG_PATH string = "/etc/markhor/config.yaml"

// In a MarkhorSecret, the name of the custom label added to Secrets managed by Markhor
const MSPARAMS_MANAGED_LABEL_KEY string = "managedLabel"

const SLOG_EVENT_ID_KEY string = "eventId"

// Timeouts for the http(s) servers that Markhor creates (healthcheck and admission controller)
const SERVER_READ_TIMEOUT_SECONDS = 5
const SERVER_WRITE_TIMEOUT_SECONDS = 10
