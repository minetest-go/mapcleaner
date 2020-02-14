
function mapcleaner.is_player_near_chunk(chunk_pos)
  local min_mapblock_pos = mapcleaner.get_mapblocks_from_chunk(chunk_pos)
  local min_pos = mapcleaner.get_blocks_from_mapblock(min_mapblock_pos)

  for _, player in ipairs(minetest.get_connected_players()) do
    local ppos = player:get_pos()
    local distance = vector.distance(min_pos, ppos)

    if distance < (80 * 3) then
      return true
    end
  end

  return false
end
