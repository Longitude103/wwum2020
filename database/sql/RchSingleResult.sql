SELECT cell_node, cell_size, rw, clm, dt, rslt
from (SELECT cell_node, cell_size, rw, clm, dt, sum(result) rslt
      FROM results LEFT JOIN cellrc on cell_node = node
      WHERE file_type = $1 group by cell_node, cell_size, dt)
where rslt > 0;