# Demo
Cloud Operations Workflows in Action

## Step 1
Generate an apikey by either visiting the [Cloud UI](https://cloud.temporal.io/settings/api-keys) or using [tcld](https://github.com/temporalio/tcld#creating-an-api-key).

## Step 2
Start the demo worker with the APIKey
```
TEMPORAL_CLOUD_API_KEY=<apikey> go run ./demo
```

## Step 3
Run a workflow using `tctl` for example to invoke `get-users` workflow run:
```
tctl wf start --tq demo --wt cloud-operations-workflows.get-users -i '{}'
```

