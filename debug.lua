

--[[
minetest.register_on_generated(function(minp, maxp)
	minetest.chat_send_all("minp: " .. minetest.pos_to_string(minp) ..
		" Chunk: " .. minetest.pos_to_string(get_chunkpos_from_pos(minp)))
	minetest.chat_send_all("maxp: " .. minetest.pos_to_string(maxp) ..
		" Chunk: " .. minetest.pos_to_string(get_chunkpos_from_pos(maxp)))
end)
--]]
