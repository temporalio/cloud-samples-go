# Worker

A worker that invokes the temporal cloud apis to perform various operations.

Temporal Cloud APIs are located at Repo: [temporalio/cloud-api](https://github.com/temporalio/api-cloud)

## Running the worker

### Step 1: Generate an apikey
Generate an apikey by either visiting the [Cloud UI](https://cloud.temporal.io/settings/api-keys) or using [tcld](https://github.com/temporalio/tcld#creating-an-api-key).

### Step 2: Start local worker 
Start local worker using the APIKey
```
TEMPORAL_CLOUD_API_KEY=<apikey> go run ./cmd/worker
```

### Step 2: Start cloud worker
Start cloud worker using the Cloud Namespace, TLS Cert, TLS Key and APIKey
```
TEMPORAL_CLOUD_NAMESPACE=<namespace.accountId> TEMPORAL_CLOUD_TLS_CERT=</path/to/ca.pem> TEMPORAL_CLOUD_TLS_KEY=</path/to/ca.key> TEMPORAL_CLOUD_API_KEY=<apikey> go run ./cmd/worker
```

### Step 3: Run a workflow
Run a workflow using `tctl` for example to invoke `get-users` workflow run:
```
tctl wf start --tq demo --wt tmprlcloud-wf.get-users -i '{}'
```

