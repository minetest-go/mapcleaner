mapcleaner = {
	storage = minetest.get_mod_storage(),

	-- step interval
	step_interval = 1.0,

	max_lag = 1.5,

	-- time usage per step
	max_time_usage = 50000
}

local MP = minetest.get_modpath("mapcleaner")

dofile(MP .. "/functions.lua")
dofile(MP .. "/protection.lua")
dofile(MP .. "/generated.lua")
dofile(MP .. "/max_lag.lua")
dofile(MP .. "/delete.lua")
dofile(MP .. "/chat.lua")
dofile(MP .. "/globalstep.lua")
