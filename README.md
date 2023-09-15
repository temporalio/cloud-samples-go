# Temporal Cloud API Samples

This repository contains several sample Workflow applications that demonstrate the various capabilities of the Temporal Cloud APIs.

Workflows that can be used to manage resources on teamporl cloud.

* Temporal Cloud API repo: [temporalio/cloud-api](https://github.com/temporalio/api-cloud)

## Workflows

| WorkflowType                               | Description                                                                                                                                       |
| ------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| `tmprlcloud-wf.get-user`        | Get an existing user.                                                                                                                             |
| `tmprlcloud-wf.get-users`       | List all users.                                                                                                                                   |
| `tmprlcloud-wf.create-user`     | Create a new user.                                                                                                                                |
| `tmprlcloud-wf.update-user`     | Update an existing user.                                                                                                                          |
| `tmprlcloud-wf.delete-user`     | Delete an existing user.                                                                                                                          |
| `tmprlcloud-wf.reconcile-user`  | Reconcile a user. Creates the user if one does not exist, otherwise updates the existing one.                                                     |
| `tmprlcloud-wf.reconcile-users` | Reconcile set of users. Creates the users that do not exist, updates the existing ones. Optionally can delete the users that are not unaccounted. |
