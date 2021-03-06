# Testing stuff

1. [Local testing snippets](#local-testing-snippets)
2. [Running the test suite](#running-the-test-suite)
   1. [S3 Localstack](#s3-localstack)
   2. [Postgres](#postgres)
   3. [Azure Blob Storage](#azure-blob-storage)

## Local testing snippets

Below is a list of curl commands to test the sever.

The `--debug` flag disables claims validation on the JWTs, so you can just reuse these

```shell
# Search Cache
curl -XGET \
     --header 'Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJtM1VTZURvQ1ZtYzdOLXp2YmFpMTlEQ1VEbyJ9.eyJuYW1laWQiOiJkZGRkZGRkZC1kZGRkLWRkZGQtZGRkZC1kZGRkZGRkZGRkZGQiLCJzY3AiOiJBY3Rpb25zLkdlbmVyaWNSZWFkOjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCBBY3Rpb25zLlVwbG9hZEFydGlmYWN0czowMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAvMTpCdWlsZC9CdWlsZC8zMyBMb2NhdGlvblNlcnZpY2UuQ29ubmVjdCBSZWFkQW5kVXBkYXRlQnVpbGRCeVVyaTowMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAvMTpCdWlsZC9CdWlsZC8zMyIsIklkZW50aXR5VHlwZUNsYWltIjoiU3lzdGVtOlNlcnZpY2VJZGVudGl0eSIsImh0dHA6Ly9zY2hlbWFzLnhtbHNvYXAub3JnL3dzLzIwMDUvMDUvaWRlbnRpdHkvY2xhaW1zL3NpZCI6IkRERERERERELUREREQtRERERC1ERERELURERERERERERERERCIsImh0dHA6Ly9zY2hlbWFzLm1pY3Jvc29mdC5jb20vd3MvMjAwOC8wNi9pZGVudGl0eS9jbGFpbXMvcHJpbWFyeXNpZCI6ImRkZGRkZGRkLWRkZGQtZGRkZC1kZGRkLWRkZGRkZGRkZGRkZCIsImF1aSI6IjQ0OGVmMWY2LTMxNzAtNDEwNC04OTBiLTIwMmM1YmIzZmU3NSIsInNpZCI6IjUxNmQ5ODAzLTM4MzctNGEyZC05ZGUwLWI5NTdkMDhmYWU5OSIsImFjIjoiW3tcIlNjb3BlXCI6XCJyZWZzL2hlYWRzL21hc3RlclwiLFwiUGVybWlzc2lvblwiOjN9XSIsIm9yY2hpZCI6ImRiMTZlMzZmLTc1NmEtNGVmYi1hMTZjLTUwYjE4ZTRhMjdiNi50ZXN0Ll9fZGVmYXVsdCIsImlzcyI6InZzdG9rZW4uYWN0aW9ucy5naXRodWJ1c2VyY29udGVudC5jb20iLCJhdWQiOiJ2c3Rva2VuLmFjdGlvbnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tfHZzbzpiNTIyYjg4Yi02MzFlLTQ1MTEtOGZiNi02YzI1OWM1YjI3NzIiLCJuYmYiOjE2MzU4OTM1MzQsImV4cCI6MTYzNTkxNjMzNH0.aJlSr8IW25Xihe3YTL5bAXHSVq1ZcbYgtx22YbSbywnntKaPP0FdzX0c4Be6XR83Or7PGFDj8tusnD4yE2D_BNHkOotLgXkkce569QBv2gjkgACD6vdALjP7eufC1AUiZip-p4NYp_j4W-giCuJtg2x_eJSVmsknwVhTffQeJN58T-sS1eIuZNLhx-gMfMmcJSU3N69BVGtKv6bcrgiCBwfLqPyroHZK_dyfOZQEPxH8Qqob3ImjHmJKyJfIhz8SAf4bjSNSPTSMBAp4Fe7_ca79ikPWVTEyTBcQOvG_zrgR26X9m-lQT_dibNV62Ir4-aY2xk52wKU93pUjBZSSaQ' \
     --header 'Content-Type: application/json' \
     'http://localhost:8080/repokey/_apis/artifactcache/cache?keys=some-cache-key%2Csome-cache-&version=some-hex-1'

# Utilise restore-keys
curl -XGET \
     --header 'Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJtM1VTZURvQ1ZtYzdOLXp2YmFpMTlEQ1VEbyJ9.eyJuYW1laWQiOiJkZGRkZGRkZC1kZGRkLWRkZGQtZGRkZC1kZGRkZGRkZGRkZGQiLCJzY3AiOiJBY3Rpb25zLkdlbmVyaWNSZWFkOjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCBBY3Rpb25zLlVwbG9hZEFydGlmYWN0czowMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAvMTpCdWlsZC9CdWlsZC8zMyBMb2NhdGlvblNlcnZpY2UuQ29ubmVjdCBSZWFkQW5kVXBkYXRlQnVpbGRCeVVyaTowMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAvMTpCdWlsZC9CdWlsZC8zMyIsIklkZW50aXR5VHlwZUNsYWltIjoiU3lzdGVtOlNlcnZpY2VJZGVudGl0eSIsImh0dHA6Ly9zY2hlbWFzLnhtbHNvYXAub3JnL3dzLzIwMDUvMDUvaWRlbnRpdHkvY2xhaW1zL3NpZCI6IkRERERERERELUREREQtRERERC1ERERELURERERERERERERERCIsImh0dHA6Ly9zY2hlbWFzLm1pY3Jvc29mdC5jb20vd3MvMjAwOC8wNi9pZGVudGl0eS9jbGFpbXMvcHJpbWFyeXNpZCI6ImRkZGRkZGRkLWRkZGQtZGRkZC1kZGRkLWRkZGRkZGRkZGRkZCIsImF1aSI6IjQ0OGVmMWY2LTMxNzAtNDEwNC04OTBiLTIwMmM1YmIzZmU3NSIsInNpZCI6IjUxNmQ5ODAzLTM4MzctNGEyZC05ZGUwLWI5NTdkMDhmYWU5OSIsImFjIjoiW3tcIlNjb3BlXCI6XCJyZWZzL2hlYWRzL21hc3RlclwiLFwiUGVybWlzc2lvblwiOjN9XSIsIm9yY2hpZCI6ImRiMTZlMzZmLTc1NmEtNGVmYi1hMTZjLTUwYjE4ZTRhMjdiNi50ZXN0Ll9fZGVmYXVsdCIsImlzcyI6InZzdG9rZW4uYWN0aW9ucy5naXRodWJ1c2VyY29udGVudC5jb20iLCJhdWQiOiJ2c3Rva2VuLmFjdGlvbnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tfHZzbzpiNTIyYjg4Yi02MzFlLTQ1MTEtOGZiNi02YzI1OWM1YjI3NzIiLCJuYmYiOjE2MzU4OTM1MzQsImV4cCI6MTYzNTkxNjMzNH0.aJlSr8IW25Xihe3YTL5bAXHSVq1ZcbYgtx22YbSbywnntKaPP0FdzX0c4Be6XR83Or7PGFDj8tusnD4yE2D_BNHkOotLgXkkce569QBv2gjkgACD6vdALjP7eufC1AUiZip-p4NYp_j4W-giCuJtg2x_eJSVmsknwVhTffQeJN58T-sS1eIuZNLhx-gMfMmcJSU3N69BVGtKv6bcrgiCBwfLqPyroHZK_dyfOZQEPxH8Qqob3ImjHmJKyJfIhz8SAf4bjSNSPTSMBAp4Fe7_ca79ikPWVTEyTBcQOvG_zrgR26X9m-lQT_dibNV62Ir4-aY2xk52wKU93pUjBZSSaQ' \
     --header 'Content-Type: application/json' \
     'http://localhost:8080/repokey/_apis/artifactcache/cache?keys=some-cache-key2%2Csome-cache-&version=some-hex-1'


# -----
# Create Cache
curl -XPOST \
     --header 'Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJtM1VTZURvQ1ZtYzdOLXp2YmFpMTlEQ1VEbyJ9.eyJuYW1laWQiOiJkZGRkZGRkZC1kZGRkLWRkZGQtZGRkZC1kZGRkZGRkZGRkZGQiLCJzY3AiOiJBY3Rpb25zLkdlbmVyaWNSZWFkOjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCBBY3Rpb25zLlVwbG9hZEFydGlmYWN0czowMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAvMTpCdWlsZC9CdWlsZC8zMyBMb2NhdGlvblNlcnZpY2UuQ29ubmVjdCBSZWFkQW5kVXBkYXRlQnVpbGRCeVVyaTowMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAvMTpCdWlsZC9CdWlsZC8zMyIsIklkZW50aXR5VHlwZUNsYWltIjoiU3lzdGVtOlNlcnZpY2VJZGVudGl0eSIsImh0dHA6Ly9zY2hlbWFzLnhtbHNvYXAub3JnL3dzLzIwMDUvMDUvaWRlbnRpdHkvY2xhaW1zL3NpZCI6IkRERERERERELUREREQtRERERC1ERERELURERERERERERERERCIsImh0dHA6Ly9zY2hlbWFzLm1pY3Jvc29mdC5jb20vd3MvMjAwOC8wNi9pZGVudGl0eS9jbGFpbXMvcHJpbWFyeXNpZCI6ImRkZGRkZGRkLWRkZGQtZGRkZC1kZGRkLWRkZGRkZGRkZGRkZCIsImF1aSI6IjQ0OGVmMWY2LTMxNzAtNDEwNC04OTBiLTIwMmM1YmIzZmU3NSIsInNpZCI6IjUxNmQ5ODAzLTM4MzctNGEyZC05ZGUwLWI5NTdkMDhmYWU5OSIsImFjIjoiW3tcIlNjb3BlXCI6XCJyZWZzL2hlYWRzL21hc3RlclwiLFwiUGVybWlzc2lvblwiOjN9XSIsIm9yY2hpZCI6ImRiMTZlMzZmLTc1NmEtNGVmYi1hMTZjLTUwYjE4ZTRhMjdiNi50ZXN0Ll9fZGVmYXVsdCIsImlzcyI6InZzdG9rZW4uYWN0aW9ucy5naXRodWJ1c2VyY29udGVudC5jb20iLCJhdWQiOiJ2c3Rva2VuLmFjdGlvbnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tfHZzbzpiNTIyYjg4Yi02MzFlLTQ1MTEtOGZiNi02YzI1OWM1YjI3NzIiLCJuYmYiOjE2MzU4OTM1MzQsImV4cCI6MTYzNTkxNjMzNH0.aJlSr8IW25Xihe3YTL5bAXHSVq1ZcbYgtx22YbSbywnntKaPP0FdzX0c4Be6XR83Or7PGFDj8tusnD4yE2D_BNHkOotLgXkkce569QBv2gjkgACD6vdALjP7eufC1AUiZip-p4NYp_j4W-giCuJtg2x_eJSVmsknwVhTffQeJN58T-sS1eIuZNLhx-gMfMmcJSU3N69BVGtKv6bcrgiCBwfLqPyroHZK_dyfOZQEPxH8Qqob3ImjHmJKyJfIhz8SAf4bjSNSPTSMBAp4Fe7_ca79ikPWVTEyTBcQOvG_zrgR26X9m-lQT_dibNV62Ir4-aY2xk52wKU93pUjBZSSaQ' \
     --header 'Content-Type: application/json' \
     --data '{"key":"some-cache-key", "version":"some-hex-1"}' \
     'http://localhost:8080/repokey/_apis/artifactcache/cache'

# Should respond with {"cacheId": some_positive_integer}

# -----
# Finish Cache
curl -XPOST \
     --header 'Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJtM1VTZURvQ1ZtYzdOLXp2YmFpMTlEQ1VEbyJ9.eyJuYW1laWQiOiJkZGRkZGRkZC1kZGRkLWRkZGQtZGRkZC1kZGRkZGRkZGRkZGQiLCJzY3AiOiJBY3Rpb25zLkdlbmVyaWNSZWFkOjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCBBY3Rpb25zLlVwbG9hZEFydGlmYWN0czowMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAvMTpCdWlsZC9CdWlsZC8zMyBMb2NhdGlvblNlcnZpY2UuQ29ubmVjdCBSZWFkQW5kVXBkYXRlQnVpbGRCeVVyaTowMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAvMTpCdWlsZC9CdWlsZC8zMyIsIklkZW50aXR5VHlwZUNsYWltIjoiU3lzdGVtOlNlcnZpY2VJZGVudGl0eSIsImh0dHA6Ly9zY2hlbWFzLnhtbHNvYXAub3JnL3dzLzIwMDUvMDUvaWRlbnRpdHkvY2xhaW1zL3NpZCI6IkRERERERERELUREREQtRERERC1ERERELURERERERERERERERCIsImh0dHA6Ly9zY2hlbWFzLm1pY3Jvc29mdC5jb20vd3MvMjAwOC8wNi9pZGVudGl0eS9jbGFpbXMvcHJpbWFyeXNpZCI6ImRkZGRkZGRkLWRkZGQtZGRkZC1kZGRkLWRkZGRkZGRkZGRkZCIsImF1aSI6IjQ0OGVmMWY2LTMxNzAtNDEwNC04OTBiLTIwMmM1YmIzZmU3NSIsInNpZCI6IjUxNmQ5ODAzLTM4MzctNGEyZC05ZGUwLWI5NTdkMDhmYWU5OSIsImFjIjoiW3tcIlNjb3BlXCI6XCJyZWZzL2hlYWRzL21hc3RlclwiLFwiUGVybWlzc2lvblwiOjN9XSIsIm9yY2hpZCI6ImRiMTZlMzZmLTc1NmEtNGVmYi1hMTZjLTUwYjE4ZTRhMjdiNi50ZXN0Ll9fZGVmYXVsdCIsImlzcyI6InZzdG9rZW4uYWN0aW9ucy5naXRodWJ1c2VyY29udGVudC5jb20iLCJhdWQiOiJ2c3Rva2VuLmFjdGlvbnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tfHZzbzpiNTIyYjg4Yi02MzFlLTQ1MTEtOGZiNi02YzI1OWM1YjI3NzIiLCJuYmYiOjE2MzU4OTM1MzQsImV4cCI6MTYzNTkxNjMzNH0.aJlSr8IW25Xihe3YTL5bAXHSVq1ZcbYgtx22YbSbywnntKaPP0FdzX0c4Be6XR83Or7PGFDj8tusnD4yE2D_BNHkOotLgXkkce569QBv2gjkgACD6vdALjP7eufC1AUiZip-p4NYp_j4W-giCuJtg2x_eJSVmsknwVhTffQeJN58T-sS1eIuZNLhx-gMfMmcJSU3N69BVGtKv6bcrgiCBwfLqPyroHZK_dyfOZQEPxH8Qqob3ImjHmJKyJfIhz8SAf4bjSNSPTSMBAp4Fe7_ca79ikPWVTEyTBcQOvG_zrgR26X9m-lQT_dibNV62Ir4-aY2xk52wKU93pUjBZSSaQ' \
     --header 'Content-Type: application/json' \
     --data '{"size":5}' \
     'http://localhost:8080/repokey/_apis/artifactcache/caches/1'
# Should respond with nothing

```

## Running the test suite

go test ./... works, though some env vars and supporting services are needed for various tests.

| Environment Variable | Example Value                           | Description                                                                                                            |
| `DB_POSTGRES`        | `postgres://user:pass@host:port/dbname` | Enables testing the PostgreSQL database backend                                                                        |
| `STORAGE_S3`         | `http://localhost:4566`                 | Enables testing the S3 storage backend, this is the URL to the localstack container                                    |
| `STORAGE_AZUREBLOB`  | `ConnectionString;Container=blah`       | Enables testing the Azure Blob Storage backend, this needs to be a connection string for a live blob storage container |

To generate the Mocks for the storage backends, run `make generate_mocks`

TODO add commands for filtering the tests

### S3 Localstack

This runs [localstack](https://github.com/localstack/localstack) which should act as a fake S3:
```shell
docker run --rm -it -e 'SERVICES=s3' -p 4566:4566 docker.io/localstack/localstack:latest
```

### Postgres

Running a postgres container for the postgres tests
```shell
docker run --rm -it -e 'POSTGRES_DB=actionscache' -e 'POSTGRES_PASSWORD=postgres' -p 5432:5432 docker.io/library/postgres:14
export DB_POSTGRES='postgres://postgres:postgres@localhost:5432/actionscache?sslmode=disable'
```

### Azure Blob Storage

Azure Blob Storage needs a live blob storage account sadly.