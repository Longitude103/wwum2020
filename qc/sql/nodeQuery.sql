select st_asgeojson(q) geojson, node
from (select st_transform(st_centroid(geom), 4326), node from model_cells where cell_type = $1) q;