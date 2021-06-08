-- dryland parcels data
SELECT i.parcel_id, a.crop_int crop1, crop1_cov, b.crop_int crop2, crop2_cov, c.crop_int crop3, crop3_cov, d.crop_int crop4, crop4_cov,
       st_area(i.geom)/43560 area, 'np' nrd, st_x(st_transform(st_centroid(i.geom), 4326)) pointx,
       st_y(st_transform(st_centroid(i.geom), 4326)) pointy, sum(st_area(st_intersection(m.geom, i.geom))/43560) s_area,
       m.soil_code, m.coeff_zone
FROM np.t2014_dry i inner join public.model_cells m on st_intersects(i.geom, m.geom)
                    LEFT join public.crops a on crop1 = a.crop_name
                    LEFT join public.crops b on crop2 = b.crop_name
                    LEFT join public.crops c on crop3 = c.crop_name
                    LEFT join public.crops d on crop4 = d.crop_name
GROUP BY i.parcel_id, a.crop_int, parcel_id, crop1_cov, b.crop_int, crop2_cov, c.crop_int, crop3_cov, d.crop_int, crop4_cov,
         st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)), st_y(st_transform(st_centroid(i.geom), 4326)),
         m.soil_code, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, st_area(i.geom)/43560,
         st_x(st_transform(st_centroid(i.geom), 4326)), st_y(st_transform(st_centroid(i.geom), 4326)), nrd, m.soil_code, m.coeff_zone
UNION ALL
SELECT i.parcel_id, a.crop_int crop1, crop1_cov, b.crop_int crop2, crop2_cov, c.crop_int crop3, crop3_cov, d.crop_int crop4, crop4_cov,
       st_area(i.geom)/43560 area, 'sp' nrd, st_x(st_transform(st_centroid(i.geom), 4326)) pointx,
       st_y(st_transform(st_centroid(i.geom), 4326)) pointy, sum(st_area(st_intersection(m.geom, i.geom))/43560) s_area,
       m.soil_code, m.coeff_zone
FROM sp.t2014_dry i inner join public.model_cells m on st_intersects(i.geom, m.geom)
                    LEFT join public.crops a on crop1 = a.crop_name
                    LEFT join public.crops b on crop2 = b.crop_name
                    LEFT join public.crops c on crop3 = c.crop_name
                    LEFT join public.crops d on crop4 = d.crop_name
GROUP BY i.parcel_id, a.crop_int, parcel_id, crop1_cov, b.crop_int, crop2_cov, c.crop_int, crop3_cov, d.crop_int, crop4_cov,
         st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)), st_y(st_transform(st_centroid(i.geom), 4326)),
         m.soil_code, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, st_area(i.geom)/43560,
         st_x(st_transform(st_centroid(i.geom), 4326)), st_y(st_transform(st_centroid(i.geom), 4326)), nrd, m.soil_code, m.coeff_zone;

-- irrigated parcel data
SELECT parcel_id, a.crop_int crop1, crop1_cov, b.crop_int crop2, crop2_cov, c.crop_int crop3, crop3_cov, d.crop_int crop4, crop4_cov, sw, gw,
       irrig_type, sw_fac, cert_num::varchar, model_id, sw_id, st_area(i.geom)/43560 area, 'np' nrd,
       st_x(st_transform(st_centroid(i.geom), 4326)) pointx, st_y(st_transform(st_centroid(i.geom), 4326)) pointy,
       sum(st_area(st_intersection(m.geom, i.geom))/43560) s_area, m.soil_code, m.coeff_zone
FROM np.t2014_irr i inner join public.model_cells m on st_intersects(i.geom, m.geom)
    LEFT join public.crops a on crop1 = a.crop_name
    LEFT join public.crops b on crop2 = b.crop_name
    LEFT join public.crops c on crop3 = c.crop_name
    LEFT join public.crops d on crop4 = d.crop_name
GROUP BY parcel_id, a.crop_int, parcel_id, crop1_cov, b.crop_int, crop2_cov, c.crop_int, crop3_cov, d.crop_int, crop4_cov, sw, gw, irrig_type, sw_fac, cert_num::varchar, model_id, st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)), st_y(st_transform(st_centroid(i.geom), 4326)), m.soil_code, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, sw, gw, irrig_type,
    sw_fac, cert_num::varchar, model_id, st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)),
    st_y(st_transform(st_centroid(i.geom), 4326)), nrd, m.soil_code, m.coeff_zone
UNION ALL
SELECT parcel_id, a.crop_int crop1, crop1_cov, b.crop_int crop2, crop2_cov, c.crop_int crop3, crop3_cov, d.crop_int crop4, crop4_cov, sw, gw,
       irr_type as irrig_type, sw_fac, i.id as cert_num, null as model_id, sw_id, st_area(i.geom)/43560 area, 'sp' nrd,
       st_x(st_transform(st_centroid(i.geom), 4326)) pointx, st_y(st_transform(st_centroid(i.geom), 4326)) pointy,
       sum(st_area(st_intersection(m.geom, i.geom))/43560) s_area, m.soil_code, m.coeff_zone
FROM sp.t2014_irr i inner join public.model_cells m on st_intersects(i.geom, m.geom)
    LEFT join public.crops a on crop1 = a.crop_name
    LEFT join public.crops b on crop2 = b.crop_name
    LEFT join public.crops c on crop3 = c.crop_name
    LEFT join public.crops d on crop4 = d.crop_name
GROUP BY parcel_id, a.crop_int, parcel_id, crop1_cov, b.crop_int, crop2_cov, c.crop_int, crop3_cov, d.crop_int, crop4_cov, sw, gw, irrig_type,
    sw_fac, i.id, model_id, st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)),
    st_y(st_transform(st_centroid(i.geom), 4326)), nrd, m.soil_code, m.coeff_zone;