select node, st_x(st_transform(st_centroid(geom), 4326)) pointx, st_y(st_transform(st_centroid(geom), 4326)) pointy,
       soil_code, coeff_zone, zone, mtg from public.model_cells;

SELECT code, st_x(st_transform(st_centroid(geom), 4326)) pointx, st_y(st_transform(st_centroid(geom), 4326)) pointy FROM public.weather_stations;