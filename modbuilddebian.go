package main

import (
    "io/ioutil"
    "log"
    "os"
	"os/exec"
    "path/filepath"
    "strings"
    "time"
    "unicode/utf8"
)

type ModbuildDebian struct {
	Modbuild
}

func (c *ModbuildDebian) Build () error {
    
    err := c.BuildInit()
    
    if err != nil {
        log.Println("Error calling build init: ", err)
        return err
    }
    
    cmd := exec.Command("cp", "-R", *appBase+c.Ss+c.Plugin.VcsCloneFolder, c.BuildFolder+"/") 
    cmd.Dir = *appBase
    err = cmd.Run()

    if err != nil {
        log.Println("Error copying from app base folder: ", err)
        return err
    }

    cmd = exec.Command("tar", "-cvzf", c.BuildFolder+c.Ss+c.Plugin.PkgNm+"_"+c.BuildVersion+".orig.tar.gz", c.BuildFolder+c.Ss+filepath.Base(c.Plugin.VcsCloneFolder))
       
    log.Println("Calling tar command : tar -cvzf " + c.BuildFolder + c.Ss + c.Plugin.PkgNm + "_" + c.BuildVersion + ".orig.tar.gz" + c.BuildFolder + c.Ss + filepath.Base(c.Plugin.VcsCloneFolder))
    cmd.Dir = *buildBase
    err = cmd.Run()

    if err != nil {
        log.Println("Error creating tarball: ", err)
        return err
    }

    log.Println("Executing: sh " + *buildBase + c.Ss + "dobuild.sh " + c.Plugin.PkgNm +
        "_" + c.BuildVersion + c.Ss +
        " " + *repoDir + " " + *distName)
    cmd = exec.Command("sh", *buildBase+c.Ss+"dobuild.sh", c.Plugin.PkgNm+"_"+c.BuildVersion+c.Ss+filepath.Base(c.Plugin.VcsCloneFolder), *repoDir, *distName)
    cmd.Dir = *buildBase
    err = cmd.Run()

    if err != nil {
        log.Println("Error running debuild/reprepro: ", err)
        return err
    }
        
    return nil
    
}

func (c *ModbuildDebian) Debianize() (string, error) {
 
 	log.Println("Starting debian config")

	c.Ss = string(os.PathSeparator)

	pkgPath := *appBase + c.Ss + c.Plugin.VcsCloneFolder
	debPath := pkgPath + c.Ss + "debian"
	os.RemoveAll(debPath)
	var sInstall string
	var fileNm string

	log.Println("Searching folder...", pkgPath)
	contents, _ := ioutil.ReadDir(pkgPath)

	for i := 0; i < len(contents); i++ {
		fileNm = contents[i].Name()
		
		if fileNm =="modpack.txt" && !validFile(pkgPath + "/" + fileNm) {
			// modpack.txt was probably empty.  Create one with some text in it.
			err := ioutil.WriteFile(pkgPath + "/modpack.txt", []byte("Added by debmod"), perm)
			
			if err != nil {
				log.Println("error adding text to modpack.txt", err)
			}
		}
		
		
		if fileNm != "debian" && !strings.Contains(fileNm, " ") &&
            !strings.Contains(fileNm, ".git") && 
            !strings.Contains(fileNm, "?") &&
            utf8.ValidString(fileNm) &&
            len(strings.Trim(fileNm, " ")) > 1 &&
            validFile(pkgPath + "/" + fileNm) {
                
			log.Println("...adding item ", fileNm)
			sInstall += fileNm + " " + *installFolder +
				c.Ss + c.Plugin.DestFolder + c.Ss + "\n"
		}
	}

	log.Println("Creating debian folder in ", pkgPath)
	os.Mkdir(debPath, perm)
	err := ioutil.WriteFile(debPath+c.Ss+"install", []byte(sInstall), perm)

	if err != nil {
		log.Println("error saving install file: error ", err)
        return c.BuildVersion, err
	}

	t := time.Now()
	c.BuildVersion = "0~" + t.Format("20060102150405")

	// Change log
	changeLog := c.Plugin.PkgNm + " (" + c.BuildVersion + "-1) " + *distName + "; urgency=low" +
		"\n\n" +
		"  * release\n\n" +
		" -- " + *maintName + "  " +
		t.Format("Mon, 02 Jan 2006 15:04:05 -0000")

	err = ioutil.WriteFile(debPath+c.Ss+"changelog", []byte(changeLog), perm)

	if err != nil {
		log.Println("error saving changelog file: error ", err)
        return c.BuildVersion, err
	}

	// Compat
	err = ioutil.WriteFile(debPath+c.Ss+"compat", []byte(*pkgCompat), perm)

	if err != nil {
		log.Println("error saving compat file: error ", err)
        return c.BuildVersion, err
	}

	dpList, dpErr := c.Db.GetPluginDepends(c.Plugin.PluginId)
	defer dpList.Close()

	// Control
	control := "Source: " + c.Plugin.PkgNm + "\n" +
		"Section: games\n" +
		"Priority: optional\n" +
		"Maintainer: " + *maintName + "\n" +
		"Build-Depends: debhelper (" + *dhVersion + ")\n" +
		"Standards-Version: 3.9.4\n\n" +
		"Package: " + c.Plugin.PkgNm + "\n" +
		"Architecture: all\n" +
		"Depends: ${shlibs:Depends}, ${misc:Depends}"

	if dpErr == nil {
		dpFile := PluginRow{}

		for dpList.Next() {

			if err := dpList.Scan(&dpFile.PluginId, &dpFile.PluginNm, &dpFile.VcsUrl,
				&dpFile.VcsCloneFolder, &dpFile.VcsCloneCmd, &dpFile.DestFolder,
				&dpFile.Description, &dpFile.Author, &dpFile.ForumLink,
				&dpFile.PkgNm, &dpFile.PkgVersion); err != nil {

				log.Fatal("Error fetching plugin dependencies: ", err)
			}

			control += ", " + dpFile.PkgNm
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
                firstConflict = false
			} else {
				control += ", " + cfFile.PkgNm
			}
		}
	}

	control += "\n"
	control += "Description: Minetest mod - " + c.Plugin.PluginNm + "\n" +
		" " + c.Plugin.Description + "\n"

	err = ioutil.WriteFile(debPath+c.Ss+"control", []byte(control), perm)

	if err != nil {
		log.Println("error saving control file: error ", err)
        return c.BuildVersion, err
	}

	// Copyright
	err = ioutil.WriteFile(debPath+c.Ss+"copyright", []byte("<todo>"), perm)

	if err != nil {
		log.Println("error saving copyright file: error ", err)
        return c.BuildVersion, err
	}

	// Rules
	rules := "#!/usr/bin/make -f\n\n" +
		"%:\n" +
		"\tdh $@"

	err = ioutil.WriteFile(debPath+c.Ss+"rules", []byte(rules), perm)

	if err != nil {
		log.Println("error saving rules file: error ", err)
	}
	return c.BuildVersion, err
    
}