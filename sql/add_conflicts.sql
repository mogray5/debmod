select * from plugin order by plugin_nm;
select * from plugin_conflicts;
select * from plugin where dest_folder in ('mobs', 'snow');

commit;

INSERT INTO plugin_conflicts(
            plugin_id, conflicts_id)
select s.plugin_id, c.plugin_id
from plugin s cross join plugin c
where s.plugin_nm = 'Ethereal NG'
	and c.plugin_nm = 'Ethereal';

INSERT INTO plugin_conflicts(
            plugin_id, conflicts_id)
select s.plugin_id, c.plugin_id
from plugin s cross join plugin c
where s.plugin_nm = 'Ethereal'
	and c.plugin_nm = 'Ethereal NG';

INSERT INTO plugin_conflicts(
            plugin_id, conflicts_id)
select s.plugin_id, c.plugin_id
from plugin s cross join plugin c
where s.plugin_nm = 'Fishing! - Mossmanikin''s version'
	and c.plugin_nm = 'Fishing - Minetestforfun version';
