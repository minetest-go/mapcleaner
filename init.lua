local storage = minetest.get_mod_storage()

mapcleaner = {
	storage = storage,

	-- step interval
	step_interval = tonumber(storage:get("step_interval") or "1.0"),

	-- above that lag the process is stopped
	max_lag = tonumber(storage:get("max_lag") or "1.5"),

	-- time usage per step
	max_time_usage = tonumber(storage:get("max_time_usage") or "50000"),

	-- run state
	run = (storage:get("run") or "0") == 1
}

local MP = minetest.get_modpath("mapcleaner")

dofile(MP .. "/functions.lua")
dofile(MP .. "/protection.lua")
dofile(MP .. "/presence.lua")
dofile(MP .. "/generated.lua")
dofile(MP .. "/max_lag.lua")
dofile(MP .. "/delete.lua")
dofile(MP .. "/chat.lua")
dofile(MP .. "/globalstep.lua")
