package main

import (
    "os"
    "os/exec"
)

type Modbuild struct {
	BuildFolder string
    BuildVersion string
    Plugin *PluginRow
    Ss string //PathSeparator
    Db *Repodb
}

func (c *Modbuild) Init () {
	c.Ss = string(os.PathSeparator)
}

func (c *Modbuild) BuildInit() error {
    c.BuildFolder = *buildBase + c.Ss + c.Plugin.PkgNm + "_" + c.BuildVersion
    var cmd *exec.Cmd
    cmd = exec.Command("mkdir", c.BuildFolder)
	cmd.Dir = *buildBase
	return cmd.Run()
}
 