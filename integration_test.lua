
minetest.log("warning", "[TEST] integration-test enabled!")

local function after_emerge()
	local data = minetest.write_json({ success = true }, true);
	local file = io.open(minetest.get_worldpath().."/integration_test.json", "w" );
	if file then
		file:write(data)
		file:close()
	end

	minetest.log("warning", "[TEST] integration tests done!")
	minetest.request_shutdown("success")
end

local function emerge()
	minetest.emerge_area({
		x = 0,
		y = 0,
		z = 0,
	}, {
		x = 120,
		y = 120,
		z = 120
	}, function(_, _, calls_remaining)
		if calls_remaining == 0 then
			after_emerge()
		end
	end)
end

minetest.register_on_mods_loaded(function()
	minetest.after(1, emerge)
end)
