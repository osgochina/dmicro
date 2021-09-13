package process

type ProcessInfo struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Start         int    `json:"start"`
	Stop          int    `json:"stop"`
	Now           int    `json:"now"`
	State         int    `json:"state"`
	Statename     string `json:"statename"`
	Spawnerr      string `json:"spawnerr"`
	Exitstatus    int    `json:"exitstatus"`
	Logfile       string `json:"logfile"`
	StdoutLogfile string `json:"stdout_logfile"`
	StderrLogfile string `json:"stderr_logfile"`
	Pid           int    `json:"pid"`
}
