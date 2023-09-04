# Temporal Cloud Operations Workflows

Workflows that can be used to invoke the temporal cloud apis.

Currently supported workflows

| WorkflowType                               | Cloud API                  |                          |
| ------------------------------------------ | -------------------------- | ------------------------ |
| `cloud-operations-workflows.get-user`      | `CloudService/GetUser`     | Get an existing user     |
| `cloud-operations-workflows.get-users`     | `CloudService/GetUsers`    | List all users           |
| `cloud-operations-workflows.create-user`   | `CloudService/CreateUser`  | Create a new user        |
| `cloud-operations-workflows.update-user`   | `CloudService/UpdateUser`  | Update an existing user  |
| `cloud-operations-workflows.delete-user`   | `CloudService/DeleteUser`  | Delete an existing user  |


Refer the (demo)[demo/README.md] to learn how to build a temporal worker that can execute the cloud operations workflows.
