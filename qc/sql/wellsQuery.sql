select st_asgeojson(q) geojson, wellid, nrd
from (select st_transform(geom, 4326), wellid, 'np' nrd
      from np.npnrd_wells
      union
      select st_transform(geom, 4326), wellid, 'sp' nrd from sp.spnrd_wells) q;