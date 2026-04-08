# Mapcleaner — Functional Architecture

## Overview

Mapcleaner is a CLI tool that operates on a Minetest world directory. It scans the map database chunk-by-chunk and either **removes unprotected terrain** to reclaim storage, or **exports protected regions** to a separate SQLite database for backup/migration.

It is designed to be resumable: scan progress is persisted to `mapcleaner.json` after every Z-stride so the process can be safely stopped and restarted.

---

## Coordinate System

Minetest uses a 3-tier coordinate hierarchy:

```
Node        — smallest unit (1×1×1 block)
Mapblock    — 16×16×16 nodes
Chunk       — 5×5×5 mapblocks (80×80×80 nodes)
```

Mapcleaner operates at the **chunk** level but interacts with the database at the **mapblock** level.

```
chunk(x,y,z) → mapblocks [x*5 .. x*5+4] × [y*5 .. y*5+4] × [z*5 .. z*5+4]
```

---

## SQLite Map Database Formats

Minetest has used two different SQLite schemas for `map.sqlite` over its history. Mapcleaner (via the `mtdb` library) handles both, but they have different capabilities.

### Old Format — `pos` column (pre-5.12.0)

Introduced when SQLite support was first added (circa 0.3.x era, implemented by JacobF on a suggestion by celeron55). All mapblock coordinates are encoded into a **single 64-bit integer** primary key:

```sql
CREATE TABLE blocks (
    pos  INT  NOT NULL PRIMARY KEY,
    data BLOB
);
```

**Encoding formula:**

$$pos = (z \times 2^{24}) + (y \times 2^{12}) + x$$

where x, y, z are 12-bit signed integers (range −2048 to +2047 in mapblock coords).

**Decoding:**
```
pos = pos + 0x800800800
x = (pos & 0xFFF) - 0x800
y = ((pos >> 12) & 0xFFF) - 0x800
z = ((pos >> 24) & 0xFFF) - 0x800
```

**Limitation for mapcleaner:** Range queries over encoded `pos` values do not map cleanly to axis-aligned X/Y/Z bounding boxes, so `prune_unprotected_batched` cannot be used — fall back to `prune_unprotected`.

### New Format — `x, y, z` columns (5.12.0+)

Introduced in **Minetest 5.12.0** (released 2024). Coordinates are stored as separate integer columns with a composite primary key:

```sql
CREATE TABLE blocks (
    x    INTEGER,
    y    INTEGER,
    z    INTEGER,
    data BLOB NOT NULL,
    PRIMARY KEY (x, z, y)
);
```

This format enables straightforward range queries like:

```sql
SELECT x, y, z, data FROM blocks
WHERE x >= ? AND x <= ?
  AND y >= ? AND y <= ?
  AND z >= ? AND z <= ?
```

**This is what `prune_unprotected_batched` requires.**

### Format Detection

The `mtdb` library (and mapcleaner's `sqliteHasPosColumn`) detects the format at runtime by inspecting `PRAGMA table_info(blocks)` and checking whether a column named `pos` exists.

```mermaid
flowchart LR
    A[Open map.sqlite] --> B["PRAGMA table_info(blocks)"]
    B --> C{column 'pos' exists?}
    C -- yes --> D["Old format
    pre-5.12.0
    single encoded integer"]
    C -- no --> E["New format
    5.12.0+
    separate x/y/z columns"]
    D --> F["prune_unprotected only
    batched mode unsupported"]
    E --> G["All modes supported
    incl. batched"]
```

### PostgreSQL

PostgreSQL has always used separate `posX`, `posY`, `posZ` columns (no legacy encoding), so all modes including `prune_unprotected_batched` are fully supported.

```sql
-- PostgreSQL schema (all versions)
CREATE TABLE blocks (
    posX  INT NOT NULL,
    posY  INT NOT NULL,
    posZ  INT NOT NULL,
    data  BYTEA,
    PRIMARY KEY (posX, posY, posZ)
);
```

---

## Startup Flow

```mermaid
flowchart TD
    A([Start]) --> B[Parse CLI flags]
    B --> C["Open block DB
    via mtdb.NewBlockDB"]
    C --> D{"areas.dat
    exists?"}
    D -- yes --> E["Parse areas.dat
    areasparser.ParseFile"]
    E --> F["PopulateAreaProtection
    for each area"]
    D -- no --> G[Warn: no areas found]
    F --> H{mode flag}
    G --> H
    H -- prune_unprotected --> I[ProcessRemoveUnprotected]
    H -- prune_unprotected_batched --> J[ProcessRemoveUnprotectedBatched]
    H -- export_protected --> K{export-all flag}
    K -- false --> L[ProcessExportProtected]
    K -- true --> M[ProccessExportAllProtected]
```

---

## Modes

### 1. `prune_unprotected`

Scans every chunk in the configured bounding box. For each chunk, it makes individual DB calls to determine if the chunk is emerged and/or protected. Unprotected emerged chunks are deleted.

**DB access pattern:** up to 133 individual `GetByPos` queries per chunk (8 for emergence check + 125 for protection check).

```mermaid
flowchart TD
    A([ProcessRemoveUnprotected]) --> B["LoadProtectedNodes
    mapcleaner_protect.txt"]
    B --> C["LoadState
    mapcleaner.json"]
    C --> D{ChunkX > ToX?}
    D -- yes --> E["Advance Z
    SaveState"]
    E --> D
    D -- no --> F{ChunkZ > ToZ?}
    F -- yes --> G["Advance Y
    PurgeCaches"]
    G --> F
    F -- no --> H{ChunkY > ToY?}
    H -- yes --> I([Done — SaveState])
    H -- no --> J["IsEmerged
    check 8 corner mapblocks"]
    J -- not emerged --> K[Advance ChunkX]
    K --> D
    J -- emerged --> L["IsProtectedWithNeighbors
    check chunk + 26 neighbors"]
    L -- protected --> M[RetainedChunks++]
    M --> K
    L -- unprotected --> N["RemoveChunk
    delete 125 mapblocks"]
    N --> O[RemovedChunks++]
    O --> K
```

---

### 2. `prune_unprotected_batched` ⚠️ Experimental

Same logic as `prune_unprotected` but replaces per-mapblock DB calls with a **single bulk SQL range query per Y-layer**. All mapblocks for `chunk_y ±1` across the full X/Z scan area are fetched once into an in-memory map, then all chunk checks for that Y-layer are served from it.

**DB access pattern:** 1 bulk query per Y-layer (instead of up to 133 queries per chunk).

```mermaid
flowchart TD
    A([ProcessRemoveUnprotectedBatched]) --> B[LoadProtectedNodes]
    B --> C[LoadState]
    C --> D["InitBatchDB
    open direct SQL connection"]
    D --> E["LoadYLayer chunk_y
    SELECT all mapblocks for Y±1"]
    E --> F{ChunkX > ToX?}
    F -- yes --> G["Advance Z
    SaveState"]
    G --> F
    F -- no --> H{ChunkZ > ToZ?}
    H -- yes --> I["Advance Y
    LoadYLayer new chunk_y"]
    I --> H
    H -- no --> J{ChunkY > ToY?}
    J -- yes --> K([Done — SaveState])
    J -- no --> L["batchIsEmerged
    lookup 8 corners in map"]
    L -- not emerged --> M[Advance ChunkX]
    M --> F
    L -- emerged --> N["batchIsProtectedWithNeighbors
    lookup 27 chunks x 125 blocks in map"]
    N -- protected --> O[RetainedChunks++]
    O --> M
    N -- unprotected --> P["RemoveChunk
    delete 125 mapblocks via DB"]
    P --> Q[RemovedChunks++]
    Q --> M
```

#### Y-Layer Cache Layout

```mermaid
block-beta
  columns 3
  block:A["chunk_y - 1 (neighbor below)"]:3
  block:B["chunk_y (current layer)"]:3
  block:C["chunk_y + 1 (neighbor above)"]:3
  style A fill:#dde,stroke:#99a
  style B fill:#bdf,stroke:#48a
  style C fill:#dde,stroke:#99a
```

Each layer is 5 mapblock-rows tall. The cache covers 15 mapblock-rows in Y, and the full configured X/Z scan range plus ±1 chunk borders.

---

### 3. `export_protected`

Iterates over areas defined in `areas.dat` and copies all chunks within each area's bounding box from the source DB into a new SQLite database at `area-export/`.

```mermaid
flowchart TD
    A([ProcessExportProtected]) --> B{"areas list
    empty?"}
    B -- yes --> C([Error: no areas])
    B -- no --> D["initializeExportDirectory
    create area-export/ + world.mt"]
    D --> E[For each area]
    E --> F[SortPos: normalize p1/p2]
    F --> G[Compute chunk bounding box]
    G --> H[For each chunk x/y/z in box]
    H --> I{Already exported?}
    I -- yes --> H
    I -- no --> J["ExportChunk
    copy 125 mapblocks src→dst"]
    J --> H
    H -- done --> E
    E -- done --> K([Done])
```

---

### 4. `export_protected` with `--export-all`

Parses **every** block in the database using the iterator and exports chunks that contain any protected node (from `mapcleaner_protect.txt`) or are within an areas-protected region. Also exports the ±1 surrounding chunks as buffer.

```mermaid
flowchart TD
    A([ProccessExportAllProtected]) --> B[LoadProtectedNodes]
    B --> C[initializeExportDirectory]
    C --> D["block_repo.Iterator
    from -33000,-33000,-33000"]
    D --> E[For each block b]
    E --> F["IsBlockProtected
    check areas + nodenames"]
    F -- not protected --> E
    F -- protected --> G[GetChunkPosFromMapblock]
    G --> H[For each neighbor chunk ±1]
    H --> I{Already exported?}
    I -- yes --> H
    I -- no --> J[ExportChunk src→dst]
    J --> H
    H -- done --> E
    E -- done --> K([Done])
```

---

## State Machine (Scan Position)

The scan traverses the 3D bounding box in **X → Z → Y** order (X is innermost, Y is outermost). State is persisted after every Z-stride advance.

```mermaid
stateDiagram-v2
    [*] --> ScanningX : LoadState
    ScanningX --> ScanningX : ChunkX++
    ScanningX --> AdvanceZ : ChunkX > ToX
    AdvanceZ --> ScanningX : ChunkX = FromX, ChunkZ++, SaveState
    AdvanceZ --> AdvanceY : ChunkZ > ToZ
    AdvanceY --> ScanningX : ChunkX = FromX, ChunkY++, ChunkZ = FromZ, PurgeCaches / LoadYLayer
    AdvanceY --> Done : ChunkY > ToY
    Done --> [*] : SaveState
```

---

## Protection System

A chunk is considered **protected** if any of the following is true:

1. It falls within a bounding box registered from `areas.dat` (`protected_areas` map)
2. Any of its 125 mapblocks contain a node whose name appears in `mapcleaner_protect.txt` (`protected_nodenames` set)

Before deleting a chunk, its **27 neighbors** (3×3×3 including itself) are also checked — a chunk adjacent to a protected region is preserved as a safety buffer.

```mermaid
flowchart LR
    A[Chunk] --> B{In protected_areas?}
    B -- yes --> P([Protected])
    B -- no --> C{"In emerged_chunks
    cache?"}
    C -- hit --> D{cached result}
    D -- protected --> P
    D -- unprotected --> U([Unprotected])
    C -- miss --> E["Load 125 mapblocks
    from DB"]
    E --> F{"Any block name
    in protected_nodenames?"}
    F -- yes --> G[Cache: protected]
    G --> P
    F -- no --> H[Cache: unprotected]
    H --> U
```

---

## File Reference

| File | Responsibility |
|---|---|
| [main.go](../main.go) | CLI entry point, flag parsing, mode dispatch |
| [state.go](../state.go) | `State` struct, `LoadState` / `SaveState` (mapcleaner.json) |
| [util.go](../util.go) | Coordinate conversion helpers (`GetChunkKey`, `GetMapblockBoundsFromChunk`, etc.) |
| [protected.go](../protected.go) | Protection registry, `IsEmerged`, `IsProtected`, `IsProtectedWithNeighbors`, `IsBlockProtected` |
| [remove.go](../remove.go) | `RemoveChunk` — deletes all 125 mapblocks of a chunk |
| [export.go](../export.go) | `ExportChunk` — copies 125 mapblocks from src to dst repository |
| [process_remove_unprotected.go](../process_remove_unprotected.go) | `ProcessRemoveUnprotected` — standard per-block scan mode |
| [process_remove_unprotected_batched.go](../process_remove_unprotected_batched.go) | `ProcessRemoveUnprotectedBatched` — Y-layer batch scan mode |
| [batch.go](../batch.go) | `InitBatchDB`, `LoadYLayer` — bulk range query infrastructure |
| [process_export_protected.go](../process_export_protected.go) | `ProcessExportProtected`, `ProccessExportAllProtected` |

---

## Function Summary

### `util.go`

| Function | Description |
|---|---|
| `GetChunkKey(x, y, z)` | Returns `"x/y/z"` string key for maps |
| `GetMapblockPosFromNode(x, y, z)` | Converts node coords → mapblock coords (÷16) |
| `GetMapblockBoundsFromChunk(x, y, z)` | Returns the 6 mapblock-coordinate bounds of a chunk |
| `GetChunkPosFromMapblock(x, y, z)` | Converts mapblock coords → chunk coords (÷5) |
| `GetChunkPosFromNode(x, y, z)` | Converts node coords → chunk coords |
| `SortPos(p1, p2)` | Returns (lo, hi) pair with min/max per axis |

### `protected.go`

| Function | Description |
|---|---|
| `LoadProtectedNodes()` | Reads `mapcleaner_protect.txt` into `protected_nodenames` |
| `PopulateAreaProtection(area)` | Marks all chunk keys in an area's bounding box as protected |
| `PurgeCaches()` | Clears `emerged_chunks` and `protected_chunks` LRU caches |
| `IsEmerged(x, y, z)` | Returns true if any of the 8 corner mapblocks of a chunk exist in DB |
| `IsProtected(x, y, z)` | Returns true if the chunk is area-protected or contains a protected node |
| `IsProtectedWithNeighbors(x, y, z)` | Returns true if the chunk or any of its 26 neighbors is protected |
| `IsBlockProtected(b)` | Returns true if a single mapblock is protected (used in export-all mode) |

### `batch.go`

| Function | Description |
|---|---|
| `InitBatchDB()` | Opens a direct SQL connection for range queries; validates SQLite format |
| `sqliteHasPosColumn(db)` | Detects legacy SQLite single-pos column format |
| `LoadYLayer(chunk_y)` | Executes one range query to load all mapblocks for chunk_y ±1 into a map |

### `process_remove_unprotected_batched.go`

| Function | Description |
|---|---|
| `batchBlockKey(x, y, z)` | Same as `GetChunkKey` — mapblock key for the layer map |
| `batchGetBlock(layer, x, y, z)` | Looks up a mapblock in the in-memory layer map |
| `batchIsEmerged(layer, x, y, z)` | Checks 8 corner mapblocks against the layer map (no DB call) |
| `batchIsProtected(layer, x, y, z)` | Checks 125 mapblocks against the layer map for protected nodes |
| `batchIsProtectedWithNeighbors(layer, x, y, z)` | Checks the chunk and 26 neighbors using the layer map |

---

## Performance Comparison

| Mode | DB queries per chunk | Suitable for |
|---|---|---|
| `prune_unprotected` | up to ~133 (`GetByPos` each) | Any backend, low-memory |
| `prune_unprotected_batched` ⚠️ | ~1 (bulk range query per Y-layer) | PostgreSQL, SQLite (new format) |

The batched mode trades memory (holding one Y-layer of blocks in RAM) for a dramatic reduction in database round-trips, making it significantly faster on large maps especially with PostgreSQL.
