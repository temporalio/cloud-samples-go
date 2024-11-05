# Worker

A worker that invokes the temporal cloud apis to perform various operations.

Temporal Cloud APIs are located at Repo: [temporalio/cloud-api](https://github.com/temporalio/api-cloud)

## Running the worker

### Step 1: Generate an apikey
Generate an apikey by either visiting the [Cloud UI](https://cloud.temporal.io/settings/api-keys) or using [tcld](https://github.com/temporalio/tcld#creating-an-api-key). For more information on api keys refer the [api keys documentation](https://docs.temporal.io/cloud/api-keys).

### Step 2: Start the worker 

To start worker that connects to a locally running temporal instance run:
```
TEMPORAL_CLOUD_API_KEY=<apikey> go run ./cmd/worker
```

Or start the worker that connects to a cloud namespace run:
```
TEMPORAL_CLOUD_NAMESPACE=<namespace.accountId> TEMPORAL_CLOUD_NAMESPACE_TLS_CERT=</path/to/cert.pem> TEMPORAL_CLOUD_NAMESPACE_TLS_KEY=</path/to/cert.key> TEMPORAL_CLOUD_API_KEY=<apikey> go run ./cmd/worker
```


Or to start the worker that connects to a cloud namespace using an api key, run:
```
TEMPORAL_CLOUD_NAMESPACE=<namespace.accountId> TEMPORAL_CLOUD_API_KEY=<apikey> TEMPORAL_CLOUD_NAMESPACE_API_KEY=<namespace_apikey> go run ./cmd/worker
```
Parameters:
- `<apikey>` is the api key that the worker will use to invoke the cloud ops apis.
- `<namespace.accountId>` is the Temporal Cloud namespace that the worker should connect to. For e.g. `prod.a2dd6`.
- `<namespace_apikey>` is the apikey to use to connect to the Temporal Cloud namespace.
- `cert.pem`, `cert.key` are the certificate-key pair to use when connecting to the Temporal Cloud namespace using MTLS auth. For more information on how to use mtls in Temporal Cloud refer to the [certificates documentation](https://docs.temporal.io/cloud/certificates).

### Step 3: Run workflows
Run a workflow using `tctl` or `temporal` cli. 

For example to invoke `get-users` workflow for a worker connected to a local temporal instance, run:
```
tctl wf start --tq demo --wt tmprlcloud-wf.get-users -i '{}'
```

## Workflows Supported

### User Workflows
- `tmprlcloud-wf.get-user`: Get a user by id
- `tmprlcloud-wf.get-users`: List users by pages
- `tmprlcloud-wf.get-all-users`: List all users
- `tmprlcloud-wf.create-user`: Create a user
- `tmprlcloud-wf.update-user`: Update a user
- `tmprlcloud-wf.delete-user`: Delete a user
- `tmprlcloud-wf.reconcile-user`: Reconcile a user
- `tmprlcloud-wf.reconcile-users`: Reconcile a list of users

### Region Workflows
- `tmprlcloud-wf.get-region`: Get a region by id
- `tmprlcloud-wf.get-all-regions`: List all regions

### Async Operation Workflows
- `tmprlcloud-wf.get-async-operation`: Get the status of an async operation by id

### Namespace Workflows
- `tmprlcloud-wf.get-namespace`: Get a namespace
- `tmprlcloud-wf.get-namespaces`: List namespaces by pages
- `tmprlcloud-wf.get-all-namespaces`: List all namespaces
- `tmprlcloud-wf.create-namespace`: Create a namespace
- `tmprlcloud-wf.update-namespace`: Update a namespace
- `tmprlcloud-wf.delete-namespace`: Delete a namespace
- `tmprlcloud-wf.reconcile-namespace`: Reconcile a namespace
- `tmprlcloud-wf.reconcile-namespaces`: Reconcile a list of namespaces

