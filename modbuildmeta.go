package main

import (
    "io/ioutil"
    "log"
    "os"
	"os/exec"
    "time"
)

type ModbuildMeta struct {
	Modbuild
}


func (c *ModbuildMeta) Build () error {
    
    err := c.BuildInit()
    
    if err != nil {
        log.Println("Error calling build init: ", err)
        return err
    }
    
    cmd := exec.Command("cp", *appBase+c.Ss+c.Plugin.PkgNm+c.Ss+"ns-control", c.BuildFolder+"/")
    cmd.Dir = *appBase
    err = cmd.Run()

    if err != nil {
        log.Println("Error copying ns-control file from app base folder: ", err)
        return err
    }
    
    log.Println("Executing: sh " + *buildBase + c.Ss + "dobuildmeta.sh " + c.Plugin.PkgNm +
        "_" + c.BuildVersion + c.Ss +
        " " + *repoDir + " " + *distName)
    cmd = exec.Command("sh", *buildBase+c.Ss+"dobuildmeta.sh", c.Plugin.PkgNm+"_"+c.BuildVersion+c.Ss, *repoDir, *distName)
    cmd.Dir = *buildBase
    err = cmd.Run()

    if err != nil {
        log.Println("Error running equivs-build/reprepro: ", err)
        return err
    }

    return err

}

func (c *ModbuildMeta) Debianize() (string, error) {
    
   	log.Println("Starting meta debian config")

	pkgPath := *appBase + c.Ss + c.Plugin.PkgNm

	os.RemoveAll(pkgPath)
    var err error

    os.Mkdir(pkgPath, perm)

	log.Println("Creating ns-control file in ", pkgPath)

	if err != nil {
		log.Println("error saving ns-control file: error ", err)
        return c.BuildVersion, err
	}

	t := time.Now()
	c.BuildVersion = "0~" + t.Format("20060102150405")

	dpList, dpErr := c.Db.GetPluginDepends(c.Plugin.PluginId)
	defer dpList.Close()

	// Control
	control := "Section: games\n" +
		"Priority: optional\n" +
		"Maintainer: " + *maintName + "\n" +
		"Standards-Version: 3.9.4\n\n" +
		"Package: " + c.Plugin.PkgNm + "\n" +
        "Version: " + c.BuildVersion + "\n" +
		"Depends: "

    firstDepends := true

	if dpErr == nil {
		dpFile := PluginRow{}

		for dpList.Next() {

			if err := dpList.Scan(&dpFile.PluginId, &dpFile.PluginNm, &dpFile.VcsUrl,
				&dpFile.VcsCloneFolder, &dpFile.VcsCloneCmd, &dpFile.DestFolder,
				&dpFile.Description, &dpFile.Author, &dpFile.ForumLink,
				&dpFile.PkgNm, &dpFile.PkgVersion); err != nil {

				log.Fatal("Error fetching plugin dependencies: ", err)
			}

            if firstDepends {
                firstDepends = false;
			} else {
				control += ", "
			}

			control += dpFile.PkgNm
		}

	} else {
		log.Println("Error pulling plugin dependencies: ", dpErr)
        return c.BuildVersion, err
	}

	cfList, cfErr := c.Db.GetPluginConflicts(c.Plugin.PluginId)
	defer cfList.Close()
	firstConflict := true

	if cfErr == nil {
		cfFile := PluginRow{}

		for cfList.Next() {

			if err := cfList.Scan(&cfFile.PluginId, &cfFile.PluginNm, &cfFile.VcsUrl,
				&cfFile.VcsCloneFolder, &cfFile.VcsCloneCmd, &cfFile.DestFolder,
				&cfFile.Description, &cfFile.Author, &cfFile.ForumLink,
				&cfFile.PkgNm, &cfFile.PkgVersion); err != nil {

				log.Fatal("Error fetching plugin conflicts: ", err)

			}

			if firstConflict {
				control += "\nConflicts: " + cfFile.PkgNm
			} else {
				control += ", " + cfFile.PkgNm
			}
		}
	}

	control += "\n"
	control += "Description: Minetest Meta Plugin Package - " + c.Plugin.PluginNm + "\n" +
		" " + c.Plugin.Description + "\n"

	err = ioutil.WriteFile(pkgPath+c.Ss+"ns-control", []byte(control), perm)

	if err != nil {
		log.Println("error saving ns-control file: error ", err)
        return c.BuildVersion, err
	}

	return c.BuildVersion, err 
}