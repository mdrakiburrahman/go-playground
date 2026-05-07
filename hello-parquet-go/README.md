# Delta/Parquet

### Delta

```bash
# Check arrow schema conversion
go run from_arrow_to_delta/delta_schema_converter.go

# APPEND generation
go run to_delta_append/delta_append_metadata.go
```

Run Delta Bulk Loader:

```bash
cd ${GIT_ROOT}/docker/delta-bulk-loader
az acr login -n arcdataanalyticsacr
docker compose down -v --remove-orphans && docker compose up -d --force-recreate
docker logs delta-bulk-loader -f
```

Run OTEL simulator:

```bash
# multi-tenant kusto
go run to_delta_adls_streaming/azure_blob_delta_streaming.go "mdrrahmansandbox" "onelake" "tenant-1" "multi-tenant" "multi-tenant" "mdrrahmansandbox"
go run to_delta_adls_streaming/azure_blob_delta_streaming.go "mdrrahmansandbox" "onelake" "tenant-2" "multi-tenant" "multi-tenant" "mdrrahmansandbox"

# single-tenant kusto
go run to_delta_adls_streaming/azure_blob_delta_streaming.go "mdrrahmansandbox" "onelake" "tenant-1" "single-tenant" "tenant-1" "mdrrahmansandbox"
go run to_delta_adls_streaming/azure_blob_delta_streaming.go "mdrrahmansandbox" "onelake" "tenant-2" "single-tenant" "tenant-2" "mdrrahmansandbox"
```

First, run `export PATH='/home/boor/.duckdb/cli/latest':$PATH`.

Query via DuckDB `duckdb`:

```sql
SET azure_transport_option_type = 'curl';
CREATE SECRET (
    TYPE AZURE,
    PROVIDER CREDENTIAL_CHAIN,
    CHAIN 'cli',
    ACCOUNT_NAME 'mdrrahmansandbox'
);
SELECT * FROM delta_scan('abfss://onelake@mdrrahmansandbox.dfs.core.windows.net/warehouse/tenant-1');
```