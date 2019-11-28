---
name: "File"
headline: "Local file-system storage"
---
### File storage

import: `github.com/reddec/storages/filestorage`

* `New`, `NewDefault`

Puts each data to separate file. File name generates from hash function (by default SHA256) applied to key. To prevent
generates too much files in one directory, each filename is chopped to 4 slices by 4 characters.

* `NewFlat`

Key is equal to file name. Sub-directories (`/` in key name) are not allowed.

Namespace are share key space with regular values.