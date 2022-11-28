select well_id, file_type, result
from wel_results
where strftime('%%Y', dt) = '$1' and strftime('%%m', dt) = '$2' and file_type < 209;