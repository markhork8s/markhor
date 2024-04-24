package pkg

// Version of this package (Markhor)
const VERSION string = "1.0.0"

// In a MarkhorSecret, this key contains the Markhor parameters
const MARKHORPARAMS_MANIFEST_KEY string = "markhorParams"

// Path where to look for the configuration file by default if none is provided
// in the CLI args
const DEFAULT_CONFIG_PATH string = "/etc/markhor/config.yaml"

// In a MarkhorSecret, the name of the custom label added to Secrets managed by Markhor
const MSPARAMS_MANAGED_ANNOTATION_KEY string = "managedAnnotation"
