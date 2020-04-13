local has_areas_mod = minetest.get_modpath("areas")

local cache = {}

-- list of protected nodes
local protected_nodes = {
	-- default protector
	"protector:protect",
	"protector:protect2",

	-- travelnet
	"travelnet:travelnet",

	-- xp protector
	"xp_redo:protector",

	-- priv protector
	"priv_protector:protector",

	-- default
	"default:chest_protected",
	"bones:bones",

	-- various
	"moreblocks:slab_desert_cobble",
	"technic:slab_marble",

	-- advtrains nodes
	"advtrains:dtrack_st",
	"advtrains:dtrack_st_45",
	"advtrains:dtrack_cr_60",
	"group:save_in_at_nodedb" -- just in case...
}

function mapcleaner.is_chunk_or_neighbours_protected(chunk_pos)
	for x=chunk_pos.x-1,chunk_pos.x+1 do
		for y=chunk_pos.y-1,chunk_pos.y+1 do
			for z=chunk_pos.z-1,chunk_pos.z+1 do
				local current_chunk = {x=x, y=y, z=z}
				if mapcleaner.is_chunk_protected(current_chunk) then
					return true
				end
			end
		end
	end

	return false
end

function mapcleaner.is_chunk_protected(chunk_pos)
	local hash = minetest.hash_node_position(chunk_pos)
	if cache[hash] then
		return true
	end

	if not mapcleaner.is_generated(chunk_pos) then
		return false
	end

	local min_mapblock_pos, max_mapblock_pos = mapcleaner.get_mapblocks_from_chunk(chunk_pos)

	for x=min_mapblock_pos.x,max_mapblock_pos.x do
		for y=min_mapblock_pos.y,max_mapblock_pos.y do
			for z=min_mapblock_pos.z,max_mapblock_pos.z do
				if mapcleaner.is_mapblock_protected({x=x, y=y, z=z}) then
					cache[hash] = true
					return true
				end
			end
		end
	end

	return false
end

-- returns true if the mapblock is protected
function mapcleaner.is_mapblock_protected(mapblock_pos)
	local min, max = mapcleaner.get_blocks_from_mapblock(mapblock_pos)

	if has_areas_mod then
		local areas_map = areas:getAreasIntersectingArea(min, max)
		local area_count = 0
		for _ in pairs(areas_map) do
			area_count = area_count + 1
		end

		if area_count > 0 then
			return true
		end
	end

	-- load area
	minetest.get_voxel_manip(min, max)

	local nodes = minetest.find_nodes_in_area(min, max, protected_nodes)

	return nodes and #nodes > 0
end
