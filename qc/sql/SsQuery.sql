select cell_node, sum(result) result
from wel_results
where strftime('%%Y', dt) = '$1' and strftime('%%m', dt) = '$2' and file_type > 208
GROUP BY cell_node;