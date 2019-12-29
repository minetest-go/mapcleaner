mapcleaner = {
	storage = minetest.get_mod_storage()
}

local MP = minetest.get_modpath("mapcleaner")

dofile(MP .. "/functions.lua")
dofile(MP .. "/protection.lua")
dofile(MP .. "/generated.lua")
dofile(MP .. "/max_lag.lua")
dofile(MP .. "/delete.lua")
dofile(MP .. "/chat.lua")
dofile(MP .. "/globalstep.lua")
