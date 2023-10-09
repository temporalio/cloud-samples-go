# Worker

A worker that invokes the temporal cloud apis to perform various operations.

Temporal Cloud APIs are located at Repo: [temporalio/cloud-api](https://github.com/temporalio/api-cloud)

## Supported Workflows

| WorkflowType                    | Description                                                                                                                                       |
| ------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| `tmprlcloud-wf.get-user`        | Get an existing user.                                                                                                                             |
| `tmprlcloud-wf.get-users`       | List all users.                                                                                                                                   |
| `tmprlcloud-wf.create-user`     | Create a new user.                                                                                                                                |
| `tmprlcloud-wf.update-user`     | Update an existing user.                                                                                                                          |
| `tmprlcloud-wf.delete-user`     | Delete an existing user.                                                                                                                          |
| `tmprlcloud-wf.reconcile-user`  | Reconcile a user. Creates the user if one does not exist, otherwise updates the existing one.                                                     |
| `tmprlcloud-wf.reconcile-users` | Reconcile set of users. Creates the users that do not exist, updates the existing ones. Optionally can delete the users that are unaccounted. |

## Running the worker

### Step 1: Generate an apikey
Generate an apikey by either visiting the [Cloud UI](https://cloud.temporal.io/settings/api-keys) or using [tcld](https://github.com/temporalio/tcld#creating-an-api-key).

### Step 2: Start the worker 
Start the worker using the APIKey
```
TEMPORAL_CLOUD_API_KEY=<apikey> go run ./cmd/worker
```

### Step 3: Run a workflow
Run a workflow using `tctl` for example to invoke `get-users` workflow run:
```
tctl wf start --tq demo --wt tmprlcloud-wf.get-users -i '{}'
```

