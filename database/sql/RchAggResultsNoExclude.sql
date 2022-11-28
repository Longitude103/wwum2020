SELECT cell_node, cell_size, rw, clm, dt, rslt
from (SELECT cell_node, cell_size, rw, clm, dt, sum(result) rslt
      FROM results LEFT JOIN cellrc on cell_node = node
      group by cell_node, cell_size, dt) as rc
where rslt > 0;