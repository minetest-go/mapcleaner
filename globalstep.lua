local has_monitoring_mod = minetest.get_modpath("monitoring")
local storage = mapcleaner.storage

local generated_count_metric, protected_count_metric, delete_count_metric, cpu_time_metric
local delete_count_total_metric, visited_count_total_metric, visited_count_metric
local chunk_x_metric, chunk_y_metric, chunk_z_metric

if has_monitoring_mod then
	generated_count_metric = monitoring.counter("mapcleaner_generated", "Visited generated chunks")
	protected_count_metric = monitoring.counter("mapcleaner_protected", "Visited protected chunks")
	delete_count_metric = monitoring.counter("mapcleaner_deleted", "Visited deleted chunks")
	delete_count_total_metric = monitoring.gauge("mapcleaner_deleted_total", "Total visited deleted chunks")
	visited_count_metric = monitoring.counter("mapcleaner_visited", "Total visited chunks")
	visited_count_total_metric = monitoring.gauge("mapcleaner_visited_total", "Total visited chunks")
	cpu_time_metric = monitoring.counter("mapcleaner_time_usage", "time usage in microseconds")
	chunk_x_metric = monitoring.gauge("mapcleaner_chunk_x", "Current x chunk")
	chunk_y_metric = monitoring.gauge("mapcleaner_chunk_y", "Current y chunk")
	chunk_z_metric = monitoring.gauge("mapcleaner_chunk_z", "Current z chunk")
end

local timer = 0
minetest.register_globalstep(function(dtime)
	timer = timer + dtime
	if timer < mapcleaner.step_interval then
		return
	end

	timer = 0
	if mapcleaner.get_max_lag() > mapcleaner.max_lag then
		return
	end

	local chunk_x = tonumber(storage:get("chunk_x") or "-388")
	local chunk_y = tonumber(storage:get("chunk_y") or "-388")
	local chunk_z = tonumber(storage:get("chunk_z") or "-388")

	local generated_count = storage:get_int("generated_count")
	local protected_count = storage:get_int("protected_count")
	local delete_count = storage:get_int("delete_count")
	local visited_count = storage:get_int("visited_count")
	local start = minetest.get_us_time()

	local function increment_pos()
		if chunk_y >= 400 then
			-- reset pos
			chunk_x = -400
			chunk_y = -400
			chunk_z = -400
		elseif chunk_z > 400 then
			chunk_z = -400
			chunk_y = chunk_y + 1
		elseif chunk_x > 400 then
			chunk_x = -400
			chunk_z = chunk_z + 1
		else
			chunk_x = chunk_x + 1
		end
	end

	while (minetest.get_us_time() - start) < mapcleaner.max_time_usage do
		visited_count = visited_count + 1

		local chunk_pos = {
			x = chunk_x,
			y = chunk_y,
			z = chunk_z
		}

		local generated = mapcleaner.is_generated(chunk_pos)
		if generated then
			generated_count = generated_count + 1
			if has_monitoring_mod then
				generated_count_metric.inc(1)
			end
			local protected = mapcleaner.is_chunk_or_neighbours_protected(chunk_pos)
			if protected then
				protected_count = protected_count + 1
				if has_monitoring_mod then
					protected_count_metric.inc(1)
				end
			else
				delete_count = delete_count + 1
				mapcleaner.delete_chunk(chunk_pos)
				if has_monitoring_mod then
					delete_count_metric.inc(1)
				end
			end
		end

		if has_monitoring_mod then
			delete_count_total_metric.set(delete_count)
			visited_count_metric.inc(1)
			visited_count_total_metric.set(visited_count)
			chunk_x_metric.set(chunk_x)
			chunk_y_metric.set(chunk_y)
			chunk_z_metric.set(chunk_z)
		end

		increment_pos()
	end

	local millis = minetest.get_us_time() - start

	if has_monitoring_mod then
		cpu_time_metric.inc(millis)
	end

	storage:set_int("generated_count", generated_count)
	storage:set_int("protected_count", protected_count)
	storage:set_int("delete_count", delete_count)
	storage:set_int("visited_count", visited_count)

	storage:set_string("chunk_x", chunk_x)
	storage:set_string("chunk_y", chunk_y)
	storage:set_string("chunk_z", chunk_z)

end)
