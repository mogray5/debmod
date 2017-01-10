package main

import (
	"database/sql"
	"database/sql/driver"

	_ "github.com/lib/pq"
)

type Repodb struct {
	Db *sql.DB
}

func (c *Repodb) Open() error {

	var err error

	c.Db, err = sql.Open("postgres", "user="+*esDbUser+
		" password="+*esDbPwd+
		" dbname="+*esDbName+
		" host="+*esDbHost+
		" port="+*esDbPort+
		" sslmode=disable")

	return err

}

func (c *Repodb) Close() {
	c.Db.Close()
}

//GetPluginList returns a list of plugins with package names that start with
//mmod
func (c Repodb) GetPluginList() (*sql.Rows, error) {

	qry := `SELECT plugin_id, plugin_nm, vcs_url, vcs_clone_folder, vcs_clone_cmd, 
       dest_folder, description, author, forum_link, pkg_nm, pkg_version
         FROM plugin where pkg_nm like 'mmod%'`

	return c.Db.Query(qry)
}

//GetGameList returns a list of plugins with package names that start with
//mgame
func (c Repodb) GetGameList() (*sql.Rows, error) {

	qry := `SELECT plugin_id, plugin_nm, vcs_url, vcs_clone_folder, vcs_clone_cmd, 
       dest_folder, description, author, forum_link, pkg_nm, pkg_version
         FROM plugin where pkg_nm like 'mgame%'`

	return c.Db.Query(qry)
}

//GetMetaList returns a list of plugins with package names that start with
//mmeta.  These represent collections of packages.
func (c Repodb) GetMetaList() (*sql.Rows, error) {

	qry := `SELECT plugin_id, plugin_nm, vcs_url, vcs_clone_folder, vcs_clone_cmd, 
       dest_folder, description, author, forum_link, pkg_nm, pkg_version
         FROM plugin where pkg_nm like 'mmeta%'`

	return c.Db.Query(qry)
}

//GetPlugin returns a single plugin mached by the passed pluginId
func (c Repodb) GetPlugin(pluginId int) (*sql.Rows, error) {

	qry := `SELECT plugin_id, plugin_nm, vcs_url, vcs_clone_folder, vcs_clone_cmd, 
       dest_folder, description, author, forum_link, pkg_nm, pkg_version
         FROM plugin where plugin_id = $1`
	return c.Db.Query(qry)
}

//GetPluginByPkgName returns a single plugin mached by the passed package
//name
func (c Repodb) GetPluginByPkgName(pkg string) (*sql.Rows, error) {

	qry := `SELECT plugin_id, plugin_nm, vcs_url, vcs_clone_folder, vcs_clone_cmd, 
       dest_folder, description, author, forum_link, pkg_nm, pkg_version
         FROM plugin where pkg_nm = $1`
	return c.Db.Query(qry, pkg)
}

//GetPluginDepends returns all plugin dependencies for the passed plugin id.
//These will end up as package dependencies.
func (c Repodb) GetPluginDepends(pluginId int64) (*sql.Rows, error) {

	qry := `select d.plugin_id, d.plugin_nm, d.vcs_url, d.vcs_clone_folder, d.vcs_clone_cmd, 
	       d.dest_folder, d.description, d.author, d.forum_link, d.pkg_nm, d.pkg_version
	       from plugin s inner join plugin_depends pd
	         on s.plugin_id = pd.plugin_id
	        inner join plugin d
	         on pd.depends_id = d.plugin_id
		where s.plugin_id = $1`
	return c.Db.Query(qry, pluginId)
}

//GetPluginRecs returns a list of recommended plugins for the passed plugin id.
func (c Repodb) GetPluginRecs(pluginId int) (*sql.Rows, error) {

	qry := `select pd.plugin_id, pd.plugin_nm, pd.vcs_url, pd.vcs_clone_folder, pd.vcs_clone_cmd, 
	       pd.dest_folder, pd.description, pd.author, pd.forum_link, pd.pkg_nm, pd.pkg_version
	       from plugin s inner join plugin_recs pd
	       	on s.plugin_id = pd.plugin_id
		where s.plugin_id = $1`
	return c.Db.Query(qry, pluginId)
}

func (c Repodb) GetPluginConflicts(pluginId int64) (*sql.Rows, error) {

	qry := `select c.plugin_id, c.plugin_nm, c.vcs_url, c.vcs_clone_folder, c.vcs_clone_cmd, 
	       c.dest_folder, c.description, c.author, c.forum_link, c.pkg_nm, c.pkg_version
	       from plugin s inner join plugin_conflicts pd
	       	 on s.plugin_id = pd.plugin_id
		inner join plugin c
			on pd.conflicts_id = c.plugin_id
		where s.plugin_id = $1`
	return c.Db.Query(qry, pluginId)
}

func (c Repodb) GetPluginFiles(pluginId int) (*sql.Rows, error) {

	qry := `SELECT file_id, plugin_id, file_nm, rel_path, checksum, 
	last_changed, new_checksum
	  FROM plugin_files
	  where plugin_id = $1`
	return c.Db.Query(qry, pluginId)
}

func (c Repodb) GetDeletedFiles(pluginId int64) (*sql.Rows, error) {

	qry := `select 0 as file_id, pf.plugin_id, pf.file_nm, pf.rel_path,
		pf.checksum 
		from plugin_files pf
		left outer join temp_files tf
        	on pf.plugin_id = tf.plugin_id
	                and pf.rel_path = tf.rel_path
			where tf.plugin_id is null 
			and pf.plugin_id = $1`

	return c.Db.Query(qry, pluginId)

}

func (c Repodb) GetNewFiles(pluginId int64) (*sql.Rows, error) {

	qry := `select 0 as file_id, tf.plugin_id, tf.file_nm, tf.rel_path,
		tf.checksum, tf.last_changed, tf.checksum as new_checksum
		from temp_files tf
		left outer join plugin_files pf
			on tf.plugin_id = pf.plugin_id
			and tf.rel_path = pf.rel_path
		where pf.plugin_id is null
		and tf.plugin_id = $1`

	return c.Db.Query(qry, pluginId)

}

func (c Repodb) GetChangedFiles(pluginId int64) (*sql.Rows, error) {

	qry := `select 0 as file_id, pf.plugin_id, pf.file_nm, pf.rel_path,
		pf.checksum, tf.last_changed, tf.checksum as new_checksum
		from plugin_files pf
		inner join temp_files tf
			on pf.plugin_id = tf.plugin_id
			and pf.rel_path = tf.rel_path
		where pf.checksum <> tf.checksum
		and pf.plugin_id = $1`

	return c.Db.Query(qry, pluginId)

}

func (c Repodb) InsertPluginFile(row *PluginFileRow) (driver.Result, error) {

	qry := `INSERT INTO plugin_files(
	            plugin_id, file_nm, rel_path, checksum, last_changed, new_checksum)
		        VALUES ($1, $2, $3, $4, $5, $6)`

	return c.Db.Exec(qry,
		row.PluginId,
		row.FileNm,
		row.RelPath,
		row.Checksum,
		row.LastChanged,
		row.NewCheckSum)

}

func (c Repodb) ClearPluginTempFiles(pluginId int64) (driver.Result, error) {

	qry := `delete from temp_files where plugin_id = $1`
	return c.Db.Exec(qry, pluginId)

}

func (c Repodb) InsertPluginTempFile(row *PluginFileRow) (driver.Result, error) {

	qry := `INSERT INTO temp_files(
	            plugin_id, file_nm, rel_path, checksum, last_changed, new_checksum)
		        VALUES ($1, $2, $3, $4, $5, $6)`

	return c.Db.Exec(qry,
		row.PluginId,
		row.FileNm,
		row.RelPath,
		row.Checksum,
		row.LastChanged,
		row.NewCheckSum)

}

func (c Repodb) SyncPluginFileChecksums(pluginId int64) (driver.Result, error) {

	qry := `UPDATE plugin_files
		SET checksum = new_checksum
		WHERE plugin_id=$1`
	return c.Db.Exec(qry, pluginId)
}

func (c Repodb) UpdatePluginFileById(row *PluginFileRow) (driver.Result, error) {

	qry := `UPDATE  plugin_files
	            SET file_nm=$1, rel_path=$2, checksum=$3, last_changed=$4, new_checksum=$5
		WHERE file_id=$6`

	return c.Db.Exec(qry,
		row.FileNm,
		row.RelPath,
		row.Checksum,
		row.LastChanged,
		row.NewCheckSum,
		row.FileId)

}

func (c Repodb) UpdatePluginFileByNm(row *PluginFileRow) (driver.Result, error) {

	qry := `UPDATE  plugin_files
	            SET checksum=$1, last_changed=$2, new_checksum=$3
		WHERE rel_path=$4
		and plugin_id=$5`

	return c.Db.Exec(qry,
		row.Checksum,
		row.LastChanged,
		row.NewCheckSum,
		row.RelPath,
		row.PluginId)

}

func (c Repodb) DeletePluginFile(row *PluginFileRow) (driver.Result, error) {

	qry := `DELETE FROM plugin_files
		WHERE plugin_id = $1
		  AND rel_path = $2`

	return c.Db.Exec(qry,
		row.PluginId,
		row.RelPath)
}
