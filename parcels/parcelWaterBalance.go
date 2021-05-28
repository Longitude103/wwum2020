package parcels

// waterBalance method takes all the parcel information (SW delivery and GW Pumping) and creates a water balance to
// determine the amount of Runoff and Deep Percolation that occurs off of each parcel and sets those values within the
// parcel struct.
func (p *Parcel) waterBalance() error {
	// TODO: Get SW AND GW Delivery for parcel

	// TODO: Compare delivery in each month to monthly NIR

	// TODO: Look at how TFG used CropSim DP and RO values?
	// My initial assumption is that CropSim RO and DP are ideal conditions under sprinkler.
	// If we deviate from those conditions then we have to adjust the RO DP by adding water or subtracting water

	// Questions that need to be looked into:
	// 1. Does CropSim account for the inefficient portion of application efficiency in the RO DP?
	// 2. If under irrigation, can we just reduce by the amount of under irrigation or should it be under irrigation times eff?

	// TODO: Return RO and DP based upon delivery.
	// TODO: If Delivery * EFF < NIR then zero additional DP, RO. If CropSim shows DP or RO, then reduce by proportion of under irrigation?
	// TODO: If Delivery * EFF > NIR then if FLOOD: DP = 75% of Over App and RO = 25%, add to CropSim numbers
	// TODO: If Delivery * EFF > NIR then if SPRINKLER: DP = 95% of Over App and RO = 5% ass to CropSim numbers
	// TODO: Over Application has secondary consumption from Coeffcrops.

	return nil
}
