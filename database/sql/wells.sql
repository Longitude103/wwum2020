SELECT mi.id, mi.wellname, mi.defaultq, mi.muni_well, mi.indust_well, mi.stop_97, mi.start_97, mc.node FROM mi_wells mi inner join model_cells mc on st_contains(mc.geom, st_translate(mi.geom, 20, 20)) where mc.cell_type = %d;

SELECT well_id, dt, pumping FROM mi_pumping where extract(YEAR from dt) between %d and %d;