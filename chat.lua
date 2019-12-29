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
