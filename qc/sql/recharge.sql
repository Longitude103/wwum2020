select st_asgeojson(q) geojson, area_ac, node from (select st_transform(geom, 4326), node, st_area(geom)/43560 area_ac from model_cells where cell_type = %d) q;

select cell_node node, sum(result) rslt from results where strftime('%%Y', dt) = '%d' and strftime('%%m', dt) = '%s' group by cell_node, strftime('%%m', dt);