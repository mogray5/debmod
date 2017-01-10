package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var appBase = flag.String("source", "/tmp/zzz", "base folder containing source")
var buildBase = flag.String("build", "/tmp/bbb", "base folder to use for builds")
var pkgName = flag.String("pkg", "mmod-zzz", "desired package name")
var installFolder = flag.String("deploy", "/usr/share/games/minetest/mods", "folder to install mod")
var distName = flag.String("dist", "wheezy", "Debian distribution")
var maintName = flag.String("maintainer", "myname <myemail>", "Maintainer name and email")
var pkgCompat = flag.String("compat", "8", "Package compatibility")
var dhVersion = flag.String("debhelper", ">= 8.0.0", "Version of Debhelper to put into control file")
var repoDir = flag.String("repo", "/var/opt/mmrepo", "Path to APT repository")
var buildMode = flag.String("buildmode", "mods", "Build mode can be mods, games, or meta")

var esDbHost = flag.String("dbhost", "localhost", "Host location of PostgresSQL database.")
var esDbPort = flag.String("dbport", "5432", "Port used by PostgresSQL database.")
var esDbUser = flag.String("dbuser", "myuser", "User to connect to mmrepo database.")
var esDbPwd = flag.String("dbpwd", "xxxx", "Password to connect to mmrepo database.")
var esDbName = flag.String("dbname", "mmrepodb", "Name of mmrepo database.")

var perm os.FileMode = 0776

//file_id, plugin_id, file_nm, checksum, last_changed, new_checksum
type PluginFileRow struct {
	FileId      int64
	PluginId    int64
	FileNm      string
	RelPath     string
	Checksum    string
	LastChanged time.Time
	NewCheckSum string
}

//plugin_id, plugin_nm, vcs_url, vcs_clone_folder, vcs_clone_cmd,
//dest_folder, description, author, forum_link, pkg_nm, pkg_version
type PluginRow struct {
	PluginId       int64
	PluginNm       string
	VcsUrl         string
	VcsCloneFolder string
	VcsCloneCmd    string
	DestFolder     string
	Description    string
	Author         string
	ForumLink      string
	PkgNm          string
	PkgVersion     string
}

func main() {

	flag.Parse()

	db := Repodb{}
	err := db.Open()

	if err != nil {
		log.Fatal("Database connection error:  ", err)
	}

	var pList *sql.Rows
	var listErr error

	// Get plugins from database
	if *pkgName == "mmod-zzz" || *pkgName == "mgame-zzz" || *pkgName == "mmeta-zzz" {
		if *buildMode == "mods" {
			pList, listErr = db.GetPluginList()
		} else if *buildMode == "games" {
			pList, listErr = db.GetGameList()
		} else {
            pList, listErr = db.GetMetaList()
        }
	} else {
		pList, listErr = db.GetPluginByPkgName(*pkgName)
	}

	if listErr != nil {
		log.Fatal("Error loading plugins: ", listErr)
	}

	defer pList.Close()

	pRow := PluginRow{}

	for pList.Next() {

		if err := pList.Scan(&pRow.PluginId, &pRow.PluginNm, &pRow.VcsUrl,
			&pRow.VcsCloneFolder, &pRow.VcsCloneCmd, &pRow.DestFolder,
			&pRow.Description, &pRow.Author, &pRow.ForumLink,
			&pRow.PkgNm, &pRow.PkgVersion); err != nil {

			log.Fatal("Error fetching plugin: ", err)

		}

		if pRow.VcsCloneCmd != "NA" && *buildMode != "meta" {
			pullErr := pullNew(&pRow)
			if pullErr != nil {
				log.Fatal("Error pulling new source: ", pullErr)
			}
			newVersion, versionErr := isNewVersion(&pRow, &db)
			if newVersion && versionErr == nil {
               
               // Build game or mod
               debBuild := ModbuildDebian{Modbuild{"", "",  &pRow, "", &db}}
  
               debBuild.Init()
               
               _, err = debBuild.Debianize()
        
               if err != nil {
                   log.Println("Error calling debianize: ", err)
                   return
               }
        
                err = debBuild.Build()
               
                if err != nil {
					log.Println("Error building plugin: ", err)
				}
			}
        } else if *buildMode == "meta" {
            // Meta package found, use equives tool to create package
            metaBuild := ModbuildMeta{Modbuild{"", "",  &pRow, "", &db}}
            
            metaBuild.Init()
            
            _, err = metaBuild.Debianize()
            
             if err != nil {
                   log.Println("Error calling debianize: ", err)
                   return
               }
            
             err = metaBuild.Build()
            
            if err != nil {
                log.Println("Error building plugin: ", err)
            }
        } else {
			log.Println("Plugin source download not support skipping: ", pRow.PluginNm)
		}
	}

	defer db.Close()

}

func pullNew(plugin *PluginRow) error {

	// Try pull a new version from vcs
	var cmd *exec.Cmd
	log.Println("Pulling source with: ", plugin.VcsCloneCmd+" "+plugin.VcsUrl)

	pathArray := strings.Split(plugin.VcsCloneFolder, "/")

	cmd = exec.Command("rm", "-rf", pathArray[0])
	cmd.Dir = *appBase
	cmd.Run()

	fParts := strings.Split(plugin.VcsCloneCmd, " ")
	if len(fParts) == 1 {
		cmd = exec.Command(plugin.VcsCloneCmd, plugin.VcsUrl)
	} else if len(fParts) == 2 {
		cmd = exec.Command(fParts[0], fParts[1], plugin.VcsUrl)

	}
	cmd.Dir = *appBase
    return cmd.Run()
    
}

func isNewVersion(plugin *PluginRow, db *Repodb) (bool, error) {

	var delCount int
	var addCount int
	var changeCount int

	db.ClearPluginTempFiles(plugin.PluginId)
	db.SyncPluginFileChecksums(plugin.PluginId)
	walkFolder(plugin.VcsCloneFolder, plugin, db)

	dList, dErr := db.GetDeletedFiles(plugin.PluginId)

	if dErr != nil {
		log.Println("Error checking for file deletions: ", dErr)
		return false, dErr
	}

	defer dList.Close()

	pFile := PluginFileRow{}

	for dList.Next() {

		if err := dList.Scan(&pFile.FileId, &pFile.PluginId, &pFile.FileNm, &pFile.RelPath,
			&pFile.Checksum); err != nil {

			log.Println("Error fetching deleted file info: ", err)
			return false, err

		}

		_, dErr = db.DeletePluginFile(&pFile)
		if dErr != nil {
			log.Println("Error deleting plugin file: ", dErr)
			return false, dErr
		}
		delCount++
	}

	aList, aErr := db.GetNewFiles(plugin.PluginId)

	if aErr != nil {
		log.Println("Error checking for file additions: ", aErr)
		return false, aErr
	}

	defer aList.Close()

	aFile := PluginFileRow{}

	for aList.Next() {

		if err := aList.Scan(&aFile.FileId, &aFile.PluginId, &aFile.FileNm, &aFile.RelPath,
			&aFile.Checksum, &aFile.LastChanged, &aFile.NewCheckSum); err != nil {

			log.Println("Error fetching added file info: ", err)
			return false, err

		}

		_, aErr = db.InsertPluginFile(&aFile)
		if aErr != nil {
			log.Println("Error adding plugin file: ", aErr)
			return false, aErr
		}
		addCount++
	}

	cList, cErr := db.GetChangedFiles(plugin.PluginId)

	if cErr != nil {
		log.Println("Error checking for file changes: ", cErr)
		return false, cErr
	}

	defer cList.Close()

	cFile := PluginFileRow{}

	for cList.Next() {

		if err := cList.Scan(&cFile.FileId, &cFile.PluginId, &cFile.FileNm, &cFile.RelPath,
			&cFile.Checksum, &cFile.LastChanged, &cFile.NewCheckSum); err != nil {

			log.Fatal("Error fetching changed file info: ", err)

		}

		log.Println("File has changed: ", cFile.RelPath)

		_, cErr = db.UpdatePluginFileByNm(&cFile)
		if cErr != nil {
			log.Println("Error changing plugin file: ", cErr)
		}
		changeCount++
	}

	return addCount > 0 || delCount > 0 || changeCount > 0, nil

}

func walkFolder(path string, plugin *PluginRow, db *Repodb) {

	ss := string(os.PathSeparator)

	var h hash.Hash = md5.New()
	fmt.Printf("searching folder...%s\n", path)
	contents, _ := ioutil.ReadDir(*appBase + ss + path)
	for i := 0; i < len(contents); i++ {

		if contents[i].IsDir() {
			if contents[i].Name() != ".git" &&
				contents[i].Name() != ".svn" &&
				contents[i].Name() != ".bzr" {
				walkFolder(path+ss+contents[i].Name(), plugin, db)
			}
		} else {

			if contents[i].Name() != ".gitignore" {

				newPath := *appBase + ss + path + ss + contents[i].Name()
				log.Println("Reading: ", contents[i].Name())
				fileContent, err := ioutil.ReadFile(newPath)
				if err != nil {
					fmt.Printf("can't read file %s \n", newPath)
				} else {
					h.Reset()
					h.Write([]byte(fileContent))
					fileChecksum := h.Sum(nil)
					tempFile := PluginFileRow{}
					tempFile.PluginId = plugin.PluginId
					tempFile.FileNm = contents[i].Name()
					tempFile.RelPath = path + ss + contents[i].Name()
					tempFile.Checksum = hex.EncodeToString(fileChecksum)
					t := time.Now()
					tempFile.LastChanged = t
					_, err = db.InsertPluginTempFile(&tempFile)
					if err != nil {
						log.Println("Error inserting temp file: ", err)
					}
				}
			}

		}

	}

}

func validFile(fPath string) bool{
    
    file, err := os.Open( fPath ) 
    if err != nil {
        log.Println("Error opening file: ", err)
    }
    
    defer file.Close()
    
    fi, err := file.Stat()
    if err != nil {
        log.Println("Error statting file: ", err)
    }
    log.Println( "File size is ", fi.Size() )
    
    return fi.Size() > 0
}
