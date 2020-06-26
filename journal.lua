
local journal_path = minetest.get_worldpath() .. "/mapgen_journal"

minetest.mkdir(journal_path)

-- write generated chunk coordinates to a journal for later cleanup
minetest.register_on_generated(function(minp, maxp)
	local date = os.date("*t")
	local filename = journal_path .. "/" .. date.year .. "-" .. date.month .. "-" .. date.day

	local file = io.open(filename, "a+")
	if file then
		file:write(minetest.pos_to_string(minp) .. "_" .. minetest.pos_to_string(maxp) .. "\n")
		io.close(file)
	end
end)
