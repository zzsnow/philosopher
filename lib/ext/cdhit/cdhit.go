package cdhit

import (
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	ucdhit "github.com/prvst/philosopher/lib/ext/cdhit/unix"
	wcdhit "github.com/prvst/philosopher/lib/ext/cdhit/win"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
)

// CDhit represents the tool configuration
type CDhit struct {
	meta.Data
	OS           string
	Arch         string
	UnixBin      string
	WinBin       string
	DefaultBin   string
	FileName     string
	FastaDB      string
	ClusterFile  string
	ClusterFasta string
}

// New constructor
func New() CDhit {

	var o CDhit
	var m meta.Data
	m.Restore(sys.Meta())

	o.UUID = m.UUID
	o.Distro = m.Distro
	o.Home = m.Home
	o.MetaFile = m.MetaFile
	o.MetaDir = m.MetaDir
	o.DB = m.DB
	o.Temp = m.Temp
	o.TimeStamp = m.TimeStamp
	o.OS = m.OS
	o.Arch = m.Arch

	o.OS = runtime.GOOS
	o.Arch = runtime.GOARCH
	o.UnixBin = m.Temp + string(filepath.Separator) + "cd-hit"
	o.WinBin = m.Temp + string(filepath.Separator) + "cd-hit.exe"
	o.FastaDB = m.Temp + string(filepath.Separator) + o.UUID + ".fasta"
	o.FileName = m.Temp + string(filepath.Separator) + "cdhit"

	return o
}

// Deploy generates binaries on workdir
func (c *CDhit) Deploy() {

	if c.OS == sys.Windows() {

		// deploy cd-hit binary
		wcdhit.Win64(c.WinBin)
		c.DefaultBin = c.WinBin

	} else {

		// deploy cd-hit binary
		ucdhit.Unix64(c.UnixBin)
		c.DefaultBin = c.UnixBin

	}

	return
}

// Run runs the cdhit binary with user's information
func (c *CDhit) Run(level float64) error {

	l := strconv.FormatFloat(level, 'E', -1, 64)
	var err error

	cmd := c.DefaultBin
	args := []string{"-i", c.FastaDB, "-o", c.ClusterFasta, "-c", l}

	run := exec.Command(cmd, args...)
	err = run.Start()
	_ = run.Wait()

	if err != nil {
		return errors.New("Could not run CD-Hit")
	}

	return nil
}