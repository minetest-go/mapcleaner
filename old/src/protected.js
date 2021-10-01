
const protected_nodes = [
	// default protector
	"protector:protect",
	"protector:protect2",
	"protector:protect_hidden",

	// travelnet
	"travelnet:travelnet",
	"travelnet:travelnet_red",
	"travelnet:travelnet_blue",
	"travelnet:travelnet_green",
	"travelnet:travelnet_black",
	"travelnet:travelnet_white",
	"travelnet:elevator",

	// xp protector
	"xp_redo:protector",

	// priv protector
	"priv_protector:protector",

	// default
	"default:chest_locked",
	"bones:bones",
	"default:mese_post_light",

	// various
	"moreblocks:slab_desert_cobble",
	"technic:slab_marble",

	// advtrains nodes
	"advtrains:dtrack_st",
	"advtrains:dtrack_st_45",
	"advtrains:dtrack_cr_60"
];

module.exports = function(mapblock){
	for (let i=0; i<protected_nodes.length; i++){
		if (mapblock.node_id_mapping[protected_nodes[i]] >= 0){
			return true;
		}
	}

	return false;
};
