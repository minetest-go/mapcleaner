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

minetest.register_chatcommand("mapcleaner_max_time", {
	description = "sets the max time usage in microseconds",
	privs = { server = true },
	func = function(name, params)
		local value = tonumber(params)
		if value then
			mapcleaner.max_time_usage = value
			storage:set_string("max_time_usage", value)
			return true, "New value: " .. value

		else
			return true, "Value: " .. mapcleaner.max_time_usage
		end
	end
})

minetest.register_chatcommand("mapcleaner_step_interval", {
	description = "sets the step_interval in seconds",
	privs = { server = true },
	func = function(name, params)
		local value = tonumber(params)
		if value then
			mapcleaner.step_interval = value
			storage:set_string("step_interval", value)
			return true, "New value: " .. value

		else
			return true, "Value: " .. mapcleaner.step_interval
		end
	end
})

minetest.register_chatcommand("mapcleaner_run", {
	description = "sets or gets the run state",
	privs = { server = true },
	func = function(name, params)
		if params == "true" then
			mapcleaner.run = true
			storage:set_string("run", "1")

		elseif params == "false" then
			mapcleaner.run = false
			storage:set_string("run", "0")

		end

		if mapcleaner.run then
			return true, "mapcleaner running!"
		else
			return true, "mapcleaner stopped"
		end
	end
})
