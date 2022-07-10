package process

func (that *Process) sysProcAttrSetPGid(s *syscall.SysProcAttr) {
	s.Setpgid = true
	s.Pdeathsig = syscall.SIGKILL
}
