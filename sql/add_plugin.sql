--delete  from plugin_files where plugin_id = 68671;
--delete from plugin where plugin_id = 4351;
select * from plugin order by plugin_nm;
select * from plugin where pkg_nm like 'mgame%' order by plugin_nm;
select * from plugin where pkg_nm like 'mmod%' order by plugin_nm;
select plugin_nm from plugin where pkg_nm like 'mmod%' order by plugin_nm;
select plugin_nm,dest_folder, pkg_nm from plugin where pkg_nm like 'mmod%' order by dest_folder;
select * from plugin order by vcs_clone_folder;
select * from plugin where pkg_nm = 'mmod-3darmor';
select * from plugin order by dest_folder;
select * from plugin where forum_link = 'https://forum.minetest.net/viewtopic.php?f=11&t=9429';

 INSERT INTO plugin(
            plugin_nm, vcs_url, vcs_clone_folder, vcs_clone_cmd, 
            dest_folder, description, author, forum_link, pkg_nm, pkg_version)
    VALUES ('LegendofMinetest', --plugin_nm
	'https://github.com/maikerumine/aftermath', --vcs_url 
	'aftermath', --vcs_clone_folder
	'git clone', --vcs_clone_cmd
        'aftermath', --dest_folder 
        'Nothing is as it seems. Can you survive a night?',  --description
        'maikerumine',  --author
        'https://forum.minetest.net/viewtopic.php?f=15&t=13700',  --forump_link
        'mgame-aftermath',  --pkg_nm
        '0~20140921034815-3');
        --pgk_version

        commit;
  
