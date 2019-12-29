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

local storage = mapcleaner.storage

minetest.register_chatcommand("mapcleaner_status", {
	func = function(name)
		local chunk_x = storage:get_int("chunk_x")
		local chunk_y = storage:get_int("chunk_y")
		local chunk_z = storage:get_int("chunk_z")
		local generated_count = storage:get_int("generated_count")
		local protected_count = storage:get_int("protected_count")
		local delete_count = storage:get_int("delete_count")
		local visited_count = storage:get_int("visited_count")

		return true, "Generated: " .. generated_count ..
			" Protected: " .. protected_count ..
			" Deleted: " .. delete_count ..
			" Count: " .. visited_count ..
			" current chunk: " .. chunk_x .. "/" .. chunk_y .. "/" .. chunk_z
	end
})
