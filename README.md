# Github Actions Cache Server for self-hosted runners

# How to run:

TODO, something something run container with volume possibly with DB


### Why:
I ran into an issue where I am running self-hosted runners on-premise to access some local resources.
I wanted some caching, which was fine, but then I ended up with a rather large cache and uploading that from where the runners are took ages hence this project :)

### What doesn't work
* Sharing caches on forks
* Probably something else I've not come across

### Roadmap: 

^^ sounds better than a todo list

* ~~Add S3 Backend~~
* Dockerfile, either scratch or distroless
* Make test harness to test end to end
* More documentation
* PostgreSQL and possibly Azure Blob storage/MySQL backends
* Add Helm Chart / Kustomize manifests
* Better instructions on how to use


# How it works:

## Architecture

TODO

TODO
GitHub will remove any cache entries that have not been accessed in over 7 days. There is no limit on the number of caches you can store, but the total size of all caches in a repository is limited to 5 GB. If you exceed this limit, GitHub will save your cache but will begin evicting caches until the total size is less than 5 GB.

## GitHub Actions Cache API

Below are my findings on how the `actions/cache` plugin works. It may not be 100% correct but seems good enough.

All requests have a `{repokey}` parameter in the URL, this is some opaque string that seems to represent the current repository.
A JWT token is also passed on every request, the only thing of value would be the scope aka ref of the repo.

### `GET {repokey}/_apis/artifactcache/cache?keys={keys}}&version={version}`

Params:
* keys - This is a comma separated list of cache paths. The first item is the cache-`key`, and any `restore-keys` are added after in order.
* version - This is some hash computed from file paths, its used as a version.

So this endpoint looks up any existing caches and returns a 204 when nothing is found, and a 200 with JSON when a cache entry is found.

The logic as I understand it goes like this:
* Get scopes from JWT
* Look for `keys[0]` (primary cache key) and `version` cache entry, if exists, return that.
* For scope in JWT scopes:
  * For `key` in `keys`
    * if key is a prefix of entry, return
* Return 204

Which comes from this paragraph:
```
key:
  npm-feature-d5ea0750
restore-keys: |
  npm-feature-
  npm-
  
For example, if a pull request contains a feature branch (the current scope) and targets the default branch (main), the action searches for key and restore-keys in the following order:

1. Key npm-feature-d5ea0750 in the feature branch scope
2. Key npm-feature- in the feature branch scope
3. Key npm- in the feature branch scope
4. Key npm-feature-d5ea0750 in the main branch scope
5. Key npm-feature- in the main branch scope
6. Key npm- in the main branch scope
```

Below is a response for 200:
```json
{
  "scope":"refs/heads/master",
  "cacheKey":"test11",
  "cacheVersion":"e5172428cbbbc7a2b72d8804c1481a209c7a49fd065144bd594e1c70b03637cf",
  "creationTime":"2021-11-02T23:02:58.89Z",
  "archiveLocation":"SOME_AZURE_BLOB_STORAGE_URL/49dfe502313cec119820a04a5ea900c9?sv=2019-07-07&sr=b&sig=zFdPWpYE1M8fnrZNSq1mn1DbCCVkzBUJIWH7d7Fvk5A%3D&se=2021-11-03T00%3A07%3A22Z&sp=r&rscl=x-e2eid-38579201-ed7f4991-b4af6075-1300d4db"
}
```

### `POST {repokey}/_apis/artifactcache/caches` - StartCache

Body:
```json
{
  "key": "{key}",
  "version": "{version}"
}
```

Params:
* key - The cache key (not the restore keys)
* version - The hash version of the filepaths (as seen above)

This endpoint seems to "initiate" a cache upload job, potentially does some prep work like set up storage in the backend for cache upload.
The cacheId seems to be an arbitary integer that means something in the backend but is used later on for future requests.

Response 201:
```json
{
  "cacheId": {idint}
}
```

### `PATCH {repokey}/_apis/artifactcache/caches/{idint}`

Params:
* idint - The cache id returned from the POST to _apis/artifactcache/cache

Chunked upload of binary data. The data seems to be zstandard compressed directories

Response 204

### `POST {repokey}/_apis/artifactcache/caches/{idint}`

Params:
* idint - The cache id returned from the POST to _apis/artifactcache/cache

Body:
```json
{
  "size": {sizeint}
}
```

This endpoint seems to submit the total number of bytes uploaded, as a santiy check, and seems to "finalize" the cache upload operation.

Response 204


## GitHub Actions JWT

Example JWT Header:
```json
{
  "typ": "JWT",
  "alg": "RS256",
  "x5t": "2m3USeDoCVmc7N-zvbai19DCUDo"
}
```

Example JWT Body
```json
{
  "nameid": "dddddddd-dddd-dddd-dddd-dddddddddddd",
  "scp": "Actions.GenericRead:00000000-0000-0000-0000-000000000000 Actions.UploadArtifacts:00000000-0000-0000-0000-000000000000/1:Build/Build/22 LocationService.Connect ReadAndUpdateBuildByUri:00000000-0000-0000-0000-000000000000/1:Build/Build/22",
  "IdentityTypeClaim": "System:ServiceIdentity",
  "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/sid": "DDDDDDDD-DDDD-DDDD-DDDD-DDDDDDDDDDDD",
  "http://schemas.microsoft.com/ws/2008/06/identity/claims/primarysid": "dddddddd-dddd-dddd-dddd-dddddddddddd",
  "aui": "a3abad71-f654-40cf-96d1-6f40995aa9ff",
  "sid": "d50e4de7-d9ba-4486-87b9-a713b61c7424",
  "ac": "[{\"Scope\":\"refs/heads/master\",\"Permission\":3}]",
  "orchid": "f6144899-8f53-494b-b543-8fe3232fafd3.test.__default",
  "iss": "vstoken.actions.githubusercontent.com",
  "aud": "vstoken.actions.githubusercontent.com|vso:b522b88b-631e-4511-8fb6-6c259c5b2772",
  "nbf": 1635891632,
  "exp": 1635914432
}
```

Well-Known URL:
```
https://token.actions.githubusercontent.com/.well-known/openid-configuration
```