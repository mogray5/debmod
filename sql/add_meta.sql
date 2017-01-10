INSERT INTO plugin(
            plugin_nm,
            vcs_clone_cmd, 
            description, author, forum_link, pkg_nm, pkg_version)
    VALUES ('Meta Ecosystem', --plugin_nm
        'NA', --vcs_clone_cmd
        'Collection of mods to add forest and sea life.',  --description
        'mogray5',  --author
        'https://forum.minetest.net/viewtopic.php?f=14&t=13051',  --forum_link
        'mmeta-ecosystem',  --pkg_nm
        '0~20140921034815-3');
        --pgk_version

commit; 
        

INSERT INTO plugin_depends(
            plugin_id, depends_id)
select s.plugin_id, d.plugin_id
from plugin s cross join plugin d
where s.plugin_nm = 'Meta Ecosystem'
	and d.plugin_nm = 'Food Chain'
UNION ALL
select s.plugin_id, d.plugin_id
from plugin s cross join plugin d
where s.plugin_nm = 'Meta Ecosystem'
	and d.plugin_nm = 'Farming Plus'
UNION ALL
select s.plugin_id, d.plugin_id
from plugin s cross join plugin d
where s.plugin_nm = 'Meta Ecosystem'
	and d.plugin_nm = 'plantlife'
UNION ALL
select s.plugin_id, d.plugin_id
from plugin s cross join plugin d
where s.plugin_nm = 'Meta Ecosystem'
	and d.plugin_nm = 'Sea'
;
INSERT INTO plugin_depends(
            plugin_id, depends_id)
select s.plugin_id, d.plugin_id
from plugin s cross join plugin d
where s.plugin_nm = 'Meta Ecosystem'
	and d.plugin_nm = 'Ambiance';