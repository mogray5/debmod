
commit;

INSERT INTO plugin_depends(
            plugin_id, depends_id)
select s.plugin_id, d.plugin_id
from plugin s cross join plugin d
where s.plugin_nm = 'Sweet Foods'
	and d.plugin_nm = 'Food';

delete from plugin_depends where plugin_id = (select plugin_id from plugin where plugin_nm = 'moretrees');

delete from plugin_files where plugin_id = (select plugin_id from plugin where plugin_nm = 'moretrees');

