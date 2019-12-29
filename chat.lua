minetest.register_chatcommand("mapcleaner_check", {
	description = "checks the current chunk",
	privs = { server = true },
	func = function(name)
		local player = minetest.get_player_by_name(name)
		if not player then
			return true, "no such player: " .. name
		end

		local pos = player:get_pos()
		local mapblock = mapcleaner.get_mapblock_from_pos(pos)
		local chunk = mapcleaner.get_chunkpos_from_pos(pos)
		local protected = mapcleaner.is_mapblock_protected(mapblock)
		local protected_str = "false"
		if protected then
			protected_str = "true"
		end

		return true, "Mapblock: " .. minetest.pos_to_string(mapblock) ..
			" Chunk: " .. minetest.pos_to_string(chunk) ..
			" Protected: " .. protected_str

	end
})


minetest.register_chatcommand("mapcleaner_iter", {
	privs = { server = true },
	func = function(name)
		local generated_count = 0
		local protected_count = 0
		local delete_count = 0
		local count = 0
		local start = minetest.get_us_time()

		for x=-2, 2 do
			for y=-2, 2 do
				for z=-2, 2 do
					count = count + 1
					local chunk_pos = {x=x, y=y, z=z}
					local generated = mapcleaner.is_generated(chunk_pos)
					if generated then
						generated_count = generated_count + 1
						local protected = mapcleaner.is_chunk_protected(chunk_pos)
						if protected then
							protected_count = protected_count + 1
						else
							delete_count = delete_count + 1
							mapcleaner.delete_chunk(chunk_pos)
						end
					end
				end
			end
		end

		local millis = minetest.get_us_time() - start

		return true, "Generated: " .. generated_count ..
			" Protected: " .. protected_count ..
			" Deleted: " .. delete_count ..
			" Count: " .. count ..
			" Millis: " .. millis
	end
})
