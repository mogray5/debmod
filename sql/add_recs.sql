select * from plugin order by plugin_nm;
select * from plugin_recs where plugin_id =67730;

commit;

INSERT INTO plugin_recs(
            plugin_id, rec_id)
select s.plugin_id, d.plugin_id
from plugin s cross join plugin d
where s.plugin_nm = 'Fishing! - Mossmanikin''s version'
	and d.plugin_nm in ('Canyon river systems');

INSERT INTO plugin_recs(
            plugin_id, rec_id)
select s.plugin_id, d.plugin_id
from plugin s cross join plugin d
where s.plugin_nm = 'Fishing! - Mossmanikin''s version'
	and d.plugin_nm in ('plantlife');
