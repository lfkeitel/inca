package grabber

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/dragonrider23/infrastructure-config-archive/interfaces"
)

func juniperExecute(host host, filename string, conf interfaces.Config) error {
	cmd := exec.Command("scripts/juniper-"+host.proto+"-config-grab.exp", host.address, conf.RemotePassword, conf.RemoteUsername, "tmp/"+filename)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		appLogger.Error(err.Error())
	}
	stdOutLogger.Info(out.String())

	err = mv("tmp/"+filename, conf.FullConfDir+"/"+filename)
	if err != nil {
		appLogger.Error(err.Error())
	}
	os.Remove("tmp/" + filename)
	return nil
}

func mv(oldpath, newpath string) error {
	oldFile, err := ioutil.ReadFile(oldpath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(newpath, oldFile, 0777)
	if err != nil {
		return err
	}
	return nil
}
