
local function execute_mapgen(callback)
	local pos1 = { x=-50, y=-10, z=-50 }
	local pos2 = { x=150, y=50, z=150 }
	minetest.emerge_area(pos1, pos2, callback)
end

local function execute_test(callback)
	execute_mapgen(function(blockpos, _, calls_remaining)
		print("Emerged: " .. minetest.pos_to_string(blockpos))
		if calls_remaining == 0 then
      callback()
		end
	end)
end

local chunks = 0
minetest.register_on_generated(function(minp)
	chunks = chunks + 1
end)

minetest.register_on_mods_loaded(function()
	minetest.after(1, function()
		execute_test(function()
			-- place bones
			minetest.set_node({ x=0, y=0, z=0 }, {name="bones:bones"})
			print("Done emerging " .. chunks .. " chunks")
		end)
	end)
end)
