---
name: "Mock"
headline: "Mocking storage that do nothing"
---
### NOP

import: `github.com/reddec/storages/memstorage`

No-Operation storage that drops any content and returns not-exists on any request.

Useful for mocking, performance testing or for any other logic that needs discard storage.

## Features

{% include feature_batch_writer.md %}