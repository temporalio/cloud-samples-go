# Temporal Cloud Operations Workflows

Workflows that can be used to manage resources on teamporl cloud.

## Workflows

| WorkflowType                                 | Description                                                                                                                                       |
| -------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cloud-operations-workflows.get-user`        | Get an existing user.                                                                                                                             |
| `cloud-operations-workflows.get-users`       | List all users.                                                                                                                                   |
| `cloud-operations-workflows.create-user`     | Create a new user.                                                                                                                                |
| `cloud-operations-workflows.update-user`     | Update an existing user.                                                                                                                          |
| `cloud-operations-workflows.delete-user`     | Delete an existing user.                                                                                                                          |
| `cloud-operations-workflows.reconcile-user`  | Reconcile a user. Creates the user if one does not exist, otherwise updates the existing one.                                                     |
| `cloud-operations-workflows.reconcile-users` | Reconcile set of users. Creates the users that do not exist, updates the existing ones. Optionally can delete the users that are not unaccounted. |

Refer the [demo](demo) to learn how to build a temporal worker that can execute the cloud operations workflows.
