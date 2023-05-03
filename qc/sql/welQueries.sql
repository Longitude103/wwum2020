select cell_node from wel_results group by cell_node;
select cell_node, sum(result) result from wel_results WHERE strftime('%%Y', dt) = '%d' GROUP BY cell_node;
select well_id, file_type, result from wel_results where strftime('%%Y', dt) = '%d' and strftime('%%m', dt) = '%s' and file_type < 209;
select cell_node, sum(result) result from wel_results where strftime('%%Y', dt) = '%d' and strftime('%%m', dt) = '%s' and file_type > 208 GROUP BY cell_node;