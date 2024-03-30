package constants

type ActionMessages struct {
	LoadedEnvFile   string
	ShutdownError   string
	ShutdownSuccess string
}
type EnvNames struct {
	EnvSource string
	Port      string
}

type EnvSources struct {
	Env  string
	File string
}

type LocalsKeys struct {
	RequestReceivedTimestamp string
	UserId                   string
}

type ResponseInfo struct {
	InternalServerError string
	Ok                  string
}
