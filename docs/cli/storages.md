
## CLI access

```go get -v github.com/reddec/storages/cmd/storages```

There is a simple command line wrapper around all currently supported database: `storage`. It provides get, put, del and
list operations over file, leveldb and redis storage.

Important: empty value implies stream.

Usage:

```
Usage:
  storages [OPTIONS] [Command] [key] [Value]

Application Options:
  -t, --db=[file|leveldb|redis|s3] DB mode (default: file) [$DB]
  -s, --stream                     Use STDIN as source of value [$STREAM]
  -0, --null                       Use zero byte as terminator for list instead of new line [$NULL]

File storage params:
      --file.location=             Root dir to store data (default: ./db) [$FILE_LOCATION]

LevelDB storage params:
      --leveldb.location=          Root dir to store data (default: ./db) [$LEVELDB_LOCATION]

Redis storage params:
      --redis.url=                 Redis URL (default: redis://localhost) [$REDIS_URL]
      --redis.namespace=           Hashmap name (default: db) [$REDIS_NAMESPACE]

S3 storage:
      --s3.bucket=                 S3 AWS bucket [$S3_BUCKET]
      --s3.endpoint=               Override AWS endpoint for AWS-capable services [$S3_ENDPOINT]
      --s3.force-path-style        Force the request to use path-style addressing [$S3_FORCE_PATH_STYLE]

Help Options:
  -h, --help                       Show this help message

Arguments:
  Command:                         what to do (put, list, get, del)
  key:                             key name
  Value:                           Value to put if stream flag is not enabled

```
