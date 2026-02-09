
-- emerge all chunks around the two "protected" chunks in (0,0,0) and (1,0,0)
mtt.register("emerge_area", {
    func = function(callback)
        minetest.emerge_area({x=-100,y=-100,z=-100}, {x=100,y=100,z=100}, function(_, _, calls_remaining)
            if calls_remaining == 0 then
                callback()
            end
        end)
    end,
    timeout = 600
})

mtt.register("setup", function(callback)

    -- chunk 0,0,0
    local pos1 = {x=0, y=0, z=0}
    local pos2 = pos1
    areas:add("singleplayer", "my area", pos1, pos2)
    areas:save()

    -- chunk 1,0,0
    minetest.set_node({x=80,y=0,z=0}, {name="bones:bones"})

    callback()
end)