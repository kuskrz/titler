package environment

import "os"

var EnvVars map[string]string = map[string]string{
	"TAG":      "title",
	"LOGLEVEL": "",
}

func InitEnv() {
	for k := range EnvVars {
		tmp_val, not_empty := os.LookupEnv(k)
		if not_empty {
			EnvVars[k] = tmp_val
		}
	}
}
