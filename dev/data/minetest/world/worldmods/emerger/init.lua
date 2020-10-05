
local function execute_mapgen(callback)
	local pos1 = { x=-50, y=-10, z=-50 }
	local pos2 = { x=50, y=50, z=50 }
	minetest.emerge_area(pos1, pos2, callback)
end

local function execute_test(callback)
	execute_mapgen(function(blockpos, action, calls_remaining)
		print("Emerged: " .. minetest.pos_to_string(blockpos))
		if calls_remaining == 0 then
      callback()
		end
	end)
end

minetest.register_on_mods_loaded(function()
	minetest.after(1, function()
		execute_test(function()
			-- place bones
			minetest.set_node({ x=0, y=0, z=0 }, {name="bones:bones"})
		end)
	end)
end)
