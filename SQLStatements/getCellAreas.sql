DROP FUNCTION public.npIrr(yr int, cell_type int);

-- NP Irrigated Acres
create or replace FUNCTION public.npIrr(yr int, cell_type int)
    returns TABLE(node bigint, nip_area double precision)
AS
$$
DECLARE
    tableName text;
    qry text;
BEGIN
    tableName = 't' || yr || '_irr';
    qry = format('select node, sum(public.st_area(public.st_intersection(c.geom, ni.geom))/43560) nip_area
    from public.model_cells c inner join np.%s ni on public.st_intersects(c.geom, ni.geom) where c.cell_type = %s group by node;', tableName, cell_type);
    return QUERY EXECUTE qry;
END
$$
    language plpgsql;

select * from public.npIrr(1953, 1);

drop function public.npDry(yr int, cell_type int);

-- NP Dryland Acres
create or replace FUNCTION public.npDry(yr int, cell_type int)
    returns TABLE(node bigint, ndp_area double precision)
AS
$$
DECLARE
    tableName text;
    qry text;
BEGIN
    tableName = 't' || yr || '_dry';
    qry = format('select node, sum(public.st_area(public.st_intersection(c.geom, ni.geom))/43560) ndp_area
    from public.model_cells c inner join np.%s ni on public.st_intersects(c.geom, ni.geom) where c.cell_type = %s group by node;', tableName, cell_type);
    return QUERY EXECUTE qry;
END
$$
    language plpgsql;

select * from public.npDry(1953, 1);

drop function spirr(yr integer, cell_type integer);

-- SP Irrigated Acres
create or replace FUNCTION public.spIrr(yr int, cell_type int)
    returns TABLE(node bigint, sip_area double precision)
AS
$$
DECLARE
    tableName text;
    qry text;
BEGIN
    tableName = 't' || yr || '_irr';
    qry = format('select node, sum(public.st_area(public.st_intersection(c.geom, ni.geom))/43560) sip_area
    from public.model_cells c inner join sp.%s ni on public.st_intersects(c.geom, ni.geom) where c.cell_type = %s group by node;', tableName, cell_type);
    return QUERY EXECUTE qry;
END
$$
    language plpgsql;

select * from public.spIrr(1953, 1);

drop function spDry(yr int, cell_type int);

-- SP Dryland Acres
create or replace FUNCTION public.spDry(yr int, cell_type int)
    returns TABLE(node bigint, sdp_area double precision)
AS
$$
DECLARE
    tableName text;
    qry text;
BEGIN
    tableName = 't' || yr || '_dry';
    qry = format('select node, sum(public.st_area(public.st_intersection(c.geom, ni.geom))/43560) sdp_area
    from public.model_cells c inner join sp.%s ni on public.st_intersects(c.geom, ni.geom) where c.cell_type = %s group by node;', tableName, cell_type);
    return QUERY EXECUTE qry;
END
$$
    language plpgsql;

select * from public.spDry(1953, 1);

drop function getCellAcres(yr int, cell_type int);

create or replace FUNCTION public.getCellAcres(yr int, cell_type int)
    returns TABLE(node bigint, soil_code bigint, coeff_zone bigint, mtg double precision, cell_area double precision, pointx double precision,
                  pointy double precision, nip_area double precision, ndp_area double precision, sip_area double precision, sdp_area double precision)
AS
$$
DECLARE
    qry text;
BEGIN
    qry = format('select m.node, m.soil_code, m.coeff_zone, m.mtg, st_area(geom)/43560 cell_area, st_x(st_transform(st_centroid(geom), 4326)) pointx,
       st_y(st_transform(st_centroid(geom), 4326)) pointy, nip_area, ndp_area, sip_area, sdp_area from model_cells m
           left join (select * from npIrr(%1$s, %2$s)) ni on m.node = ni.node
           left join (select * from npDry(%1$s, %2$s)) nd on m.node = nd.node
           left join (select * from spIrr(%1$s, %2$s)) si on m.node = si.node
           left join (select * from spDry(%1$s, %2$s)) sd on m.node = sd.node where m.nat_veg = true and m.cell_type = %2$s;', yr, cell_type);
--     raise notice '%', qry;
    return QUERY EXECUTE qry;
END
$$
    language plpgsql;

select node, soil_code, coeff_zone, mtg, cell_area, pointx, pointy, nip_area, ndp_area, sip_area, sdp_area from getCellAcres(1953, 1);