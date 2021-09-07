# Notes for Model Processes

## Fallow
The fallow in irrigated parcels, which was new in 2016, when the NRD started reporting the data needed to be handled
so that WSPP can run it. Fallow is crop 15 in CropSim but is part of a dryland only. When Fallow is listed in the irrigated
dataset, WSPP has to have an increase in ET from dryland ET to operate appropriately. To fix this, irrigated fallow is 
changed to irrigated pasture to ensure that we have an ET increase during the WSPP operations.

## Pumping Outside NRDs
The pumping amounts that are from outside SP or NP NRDs is included in ext_pumping and is recorded in ft^3 per month.
We need to use that value in the db to create the new pumping.
