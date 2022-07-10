//go:build !linux
// +build !linux

package process

/*  Note:  This is a *nix only implementation.  */

// ReapZombie 回收僵尸进程
func ReapZombie() {
}
