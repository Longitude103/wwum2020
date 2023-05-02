SELECT cell_node, rw, clm, dt, rslt
from (SELECT cell_node, rw, clm, dt, sum(result) rslt
      FROM wel_results
          LEFT JOIN cellrc on cell_node = node
      WHERE file_type = $1 group by cell_node, dt) as wrc
where rslt > 0;