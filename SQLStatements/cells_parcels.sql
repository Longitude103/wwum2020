select m.node, st_area(geom)/43560 cell_area, st_x(st_transform(st_centroid(geom), 4326)) pointx,
       st_y(st_transform(st_centroid(geom), 4326)) pointy, nip_area, ndp_area, sip_area, sdp_area from model_cells m
           left join (select node, sum(st_area(st_intersection(c.geom, ni.geom))/43560) nip_area from public.model_cells c
               inner join np.t2014_irr ni on st_intersects(c.geom, ni.geom) group by node) ni on m.node = ni.node
           left join (select node, sum(st_area(st_intersection(c.geom, nd.geom))/43560) ndp_area from public.model_cells c
               inner join np.t2014_dry nd on st_intersects(c.geom, nd.geom) group by node) nd on m.node = nd.node
           left join (select node, sum(st_area(st_intersection(c.geom, si.geom))/43560) sip_area from public.model_cells c
               inner join sp.t2014_irr si on st_intersects(c.geom, si.geom) group by node) si on m.node = si.node
           left join (select node, sum(st_area(st_intersection(c.geom, sd.geom))/43560) sdp_area from public.model_cells c
               inner join sp.t2014_dry sd on st_intersects(c.geom, sd.geom) group by node) sd on m.node = sd.node;