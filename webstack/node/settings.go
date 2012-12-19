package main

var (
	PORT string
	CERTIFICATE string
	KEY string
	PASSWORD string
)

func init() {
	getConfig := func(name string, def string) string {
		v, ok := syscall.Getenv(name)
		if ok {
			return v
		}
		return def
	}

	PORT = getConfig("PORT", "9000")
	CERTIFICATE = getConfig("CERTIFICATE", "node.crt")
	KEY = getConfig("KEY", "node.key")
	PASSWORD = getConfig("PASSWORD", "")
}
