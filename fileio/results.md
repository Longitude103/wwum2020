# CropSim Output Files

CropSim output file has the following format:

Station, year, soil, crop, Tillage, Irrigation Type

| Station | Year | Soil | Crop No | Tillage | Irr |
| ------- | ---- | ---- | ------- | ---- | --- |
| SDN6 | 2009 | 622 | 8 | 1 | 1 |

This is at station Sydney in year 2009 for soil 622 for crop 8 (corn) Tillage = 1, irrigation = Dryland

Then monthly information follows which is:

- ET month (in)
- EFF Precip (in)
- NIR (in)
- DP (in)
- RO (in)
- Precip (in)

## Crop Numbers
CropSim results crop number conversion:

| Crop | CropSim Result |
| ----- | ------------: |
| Corn | 8 |
| SugarBeets | 5 |
| EdibleBeans | 2 |
| Alfalfa | 10 |
| WinterWheat | 7 |
| Potatoes | 4 |
| Milo | 6 |
| Sunflower | 9 |
| SoyBeans | 3 |
| SmallGrain | 1 |
| Fallow | 15 |
| Past | 12 |

## Soils
Soils were greatly simplified, Soils in CropSim are:

| Soil | Class No |
| :--- | -------: |
| Sandy Soil | 412 |
| Table Lands Soils | 622 |
| Valley Soils | 722 |

Soil 412 is mostly in the NE part of the model in the sand hills and areas along the river, there are small amounts just north of the upper most canals. Soil 622 is mostly the table lands and and the northern part of the model. Soil 722 looks like the majority of the farmlands.

Soils: Identified By A 3 Digit Code That Represents The Available Water Holding Capacity (In Quarter Of Inch/Foot), Hydrologic Group (1=a,...4=d), And Depth To Groundwater Indicator (1<6ft, 2>6ft)

## Irrigation Types
Irrigation is split among types, however it appears WWUMM is only using 1 and 3.

| Irrigation Number | Type |
| ----------------- | :--- |
| 1 | Dryland |
| 2 | Fixed Irrig. Dates |
| 3 | Pivot - Sprinkler |
| 4 | Furrow Irrigation |
| 5 | Other |

# Trim Results
The cropsim files appear to be run from the start of data through current. We start the WWUMM at 1950, so the files need cut down.

We can also reduce the file sizes so that there are just the crops we have in the model represented as well as just the soils.

This folder contains a python script to trim the CSResults file down to the model start and end year. This will include to open the file, filter and remove the years that are not within the start and end range and then save it with the -tt.txt extension on the external drive.

Contains 4 variables:
str_yr = start year of model
end_yr = end year of model
raw_filedir = path to raw cropsim files currently: '/Volumes/G-Raid with TB3/WWUMM 2016/rswb/WBP/v3_1/CSResults/Crops_Raw/'
out_filedir = output files to currently: '/Volumes/G-Raid with TB3/WWUMM 2016/rswb/WBP/v3_1/CSResults/Crops'
