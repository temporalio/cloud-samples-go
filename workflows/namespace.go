package workflows

import (
	"errors"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/proto"

	"github.com/temporalio/cloud-samples-go/internal/validator"
	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"go.temporal.io/cloud-sdk/api/namespace/v1"
	"github.com/temporalio/cloud-samples-go/workflows/activities"
)

const (
	namespaceUpdateTimeout = 30 * time.Minute

	// namespace management workflow types
	GetNamespaceWorkflowType                      = workflowPrefix + "get-namespace"
	GetNamespacesWorkflowType                     = workflowPrefix + "get-namespaces"
	GetAllNamespacesWorkflowType                  = workflowPrefix + "get-all-namespaces"
	GetNamespaceWithNameWorkflow                  = workflowPrefix + "get-namespace-with-name"
	GetAllNamespacesWithAccessToNamespaceWorkflow = workflowPrefix + "get-all-namespaces-with-access-to-namespace"
	CreateNamespaceWorkflowType                   = workflowPrefix + "create-namespace"
	UpdateNamespaceWorkflowType                   = workflowPrefix + "update-namespace"
	DeleteNamespaceWorkflowType                   = workflowPrefix + "delete-namespace"
	ReconcileNamespaceWorkflowType                = workflowPrefix + "reconcile-namespace"
	ReconcileNamespacesWorkflowType               = workflowPrefix + "reconcile-namespaces"
)

type (
	ReconcileNamespaceInput struct {
		Spec *namespace.NamespaceSpec `required:"true" json:"spec"`
	}
	ReconcileNamespaceOutput struct {
		Namespace *namespace.Namespace `json:"namespace"`
		Outcome   ReconcileOutcome     `json:"outcome"`
		Error     string               `json:"error"`
	}

	ReconcileNamespacesInput struct {
		Specs             []*namespace.NamespaceSpec `required:"true" json:"specs"`
		DeleteUnaccounted bool                       `json:"delete_unaccounted"`
	}
	ReconcileNamespacesOutput struct {
		Results []*ReconcileNamespaceOutput `json:"results"`
	}

	NamespaceWorkflows interface {
		// Namespace Management Workflows
		GetNamespace(ctx workflow.Context, in *cloudservice.GetNamespaceRequest) (*cloudservice.GetNamespaceResponse, error)
		GetNamespaces(ctx workflow.Context, in *cloudservice.GetNamespacesRequest) (*cloudservice.GetNamespacesResponse, error)
		GetAllNamespaces(ctx workflow.Context) ([]*namespace.Namespace, error)
		GetNamespaceWithName(ctx workflow.Context, name string) (*namespace.Namespace, error)
		CreateNamespace(ctx workflow.Context, in *cloudservice.CreateNamespaceRequest) (*cloudservice.CreateNamespaceResponse, error)
		UpdateNamespace(ctx workflow.Context, in *cloudservice.UpdateNamespaceRequest) (*cloudservice.UpdateNamespaceResponse, error)
		DeleteNamespace(ctx workflow.Context, in *cloudservice.DeleteNamespaceRequest) (*cloudservice.DeleteNamespaceResponse, error)
		ReconcileNamespace(ctx workflow.Context, in *ReconcileNamespaceInput) (*ReconcileNamespaceOutput, error)
		ReconcileNamespaces(ctx workflow.Context, in *ReconcileNamespacesInput) (*ReconcileNamespacesOutput, error)
	}
)

func registerNamespaceWorkflows(w worker.Worker, wf NamespaceWorkflows) {
	for k, v := range map[string]any{
		GetNamespaceWorkflowType:        wf.GetNamespace,
		GetNamespacesWorkflowType:       wf.GetNamespaces,
		GetAllNamespacesWorkflowType:    wf.GetAllNamespaces,
		GetNamespaceWithNameWorkflow:    wf.GetNamespaceWithName,
		CreateNamespaceWorkflowType:     wf.CreateNamespace,
		UpdateNamespaceWorkflowType:     wf.UpdateNamespace,
		DeleteNamespaceWorkflowType:     wf.DeleteNamespace,
		ReconcileNamespaceWorkflowType:  wf.ReconcileNamespace,
		ReconcileNamespacesWorkflowType: wf.ReconcileNamespaces,
	} {
		w.RegisterWorkflowWithOptions(v, workflow.RegisterOptions{Name: k})
	}
}

func (o *ReconcileNamespaceOutput) setError(err error) {
	var applicationErr *temporal.ApplicationError
	if errors.As(err, &applicationErr) {
		o.Error = applicationErr.Error()
		o.Outcome = ReconcileOutcomeError
	}
}

// Get a namespace
func (w *workflows) GetNamespace(ctx workflow.Context, in *cloudservice.GetNamespaceRequest) (*cloudservice.GetNamespaceResponse, error) {
	return activities.GetNamespace(withInfiniteRetryActivityOptions(ctx), in)
}

// Get multiple namespaces
func (w *workflows) GetNamespaces(ctx workflow.Context, in *cloudservice.GetNamespacesRequest) (*cloudservice.GetNamespacesResponse, error) {
	return activities.GetNamespaces(withInfiniteRetryActivityOptions(ctx), in)
}

func (w *workflows) getAllNamespaces(ctx workflow.Context, name string) ([]*namespace.Namespace, error) {
	var (
		namespaces = make([]*namespace.Namespace, 0)
		pageToken  = ""
	)
	for {
		resp, err := w.GetNamespaces(ctx, &cloudservice.GetNamespacesRequest{
			Name:      name,
			PageToken: pageToken,
		})
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, resp.Namespaces...)
		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}
	return namespaces, nil
}

// Get all known namespaces
func (w *workflows) GetAllNamespaces(ctx workflow.Context) ([]*namespace.Namespace, error) {
	return w.getAllNamespaces(ctx, "")
}

// Get the namespace with name
func (w *workflows) GetNamespaceWithName(ctx workflow.Context, name string) (*namespace.Namespace, error) {
	namespaces, err := w.getAllNamespaces(ctx, name)
	if err != nil {
		return nil, err
	}
	if len(namespaces) == 0 {
		return nil, nil
	}
	if len(namespaces) > 1 {
		return nil, fmt.Errorf("multiple namespaces found for name %q", name)
	}
	return namespaces[0], nil
}

// Create a namespace
func (w *workflows) CreateNamespace(ctx workflow.Context, in *cloudservice.CreateNamespaceRequest) (*cloudservice.CreateNamespaceResponse, error) {
	return activities.CreateNamespace(withInfiniteRetryActivityOptions(ctx), in)
}

// Update a namespace
func (w *workflows) UpdateNamespace(ctx workflow.Context, in *cloudservice.UpdateNamespaceRequest) (*cloudservice.UpdateNamespaceResponse, error) {
	return activities.UpdateNamespace(withInfiniteRetryActivityOptions(ctx), in)
}

// Delete a namespace
func (w *workflows) DeleteNamespace(ctx workflow.Context, in *cloudservice.DeleteNamespaceRequest) (*cloudservice.DeleteNamespaceResponse, error) {
	return activities.DeleteNamespace(withInfiniteRetryActivityOptions(ctx), in)
}

func (w *workflows) reconcileNamespace(ctx workflow.Context, spec *namespace.NamespaceSpec, ns *namespace.Namespace) (*ReconcileNamespaceOutput, error) {
	var (
		namespaceID string
		asyncOpID   string
		out         = &ReconcileNamespaceOutput{}
		err         error
	)
	defer func() {
		if err != nil {
			out.setError(err)
		}
		if ns != nil {
			out.Namespace = ns
		} else if spec != nil {
			out.Namespace = &namespace.Namespace{
				Namespace: namespaceID,
				Spec:      spec,
			}
		}
	}()
	if ns == nil {
		var createResp *cloudservice.CreateNamespaceResponse
		// no namespace found, create one
		createResp, err = w.CreateNamespace(ctx, &cloudservice.CreateNamespaceRequest{
			Spec: spec,
		})
		if err != nil {
			return out, err
		}
		namespaceID = createResp.Namespace
		asyncOpID = createResp.AsyncOperation.Id
		out.Outcome = ReconcileOutcomeCreated

	} else if !proto.Equal(ns.Spec, spec) {
		var updateResp *cloudservice.UpdateNamespaceResponse
		// namespace found, and specs don't match,  update it
		updateResp, err = w.UpdateNamespace(ctx, &cloudservice.UpdateNamespaceRequest{
			Namespace:       ns.Namespace,
			Spec:            spec,
			ResourceVersion: ns.ResourceVersion,
		})
		if err != nil {
			return out, err
		}
		namespaceID = ns.Namespace
		asyncOpID = updateResp.AsyncOperation.Id
		out.Outcome = ReconcileOutcomeUpdated

	} else {
		// nothing to change, get the latest namespace and return
		namespaceID = ns.Namespace
		out.Outcome = ReconcileOutcomeUnchanged
		return out, nil
	}

	if asyncOpID != "" {
		// wait for the operation to complete
		_, err = w.WaitForAsyncOperation(ctx, &WaitForAsyncOperationInput{
			AsyncOperationID: asyncOpID,
			Timeout:          namespaceUpdateTimeout,
		})
		if err != nil {
			return out, err
		}
	}
	var getResp *cloudservice.GetNamespaceResponse
	getResp, err = w.GetNamespace(ctx, &cloudservice.GetNamespaceRequest{
		Namespace: namespaceID,
	})
	if err != nil {
		return out, err
	}
	ns = getResp.Namespace
	return out, nil
}

// Reconcile a namespace, create the namespace if one does not exist, or update the namespace if one does exist.
func (w *workflows) ReconcileNamespace(ctx workflow.Context, in *ReconcileNamespaceInput) (*ReconcileNamespaceOutput, error) {
	if err := validator.ValidateStruct(in); err != nil {
		return nil, fmt.Errorf("invalid input: %s", err)
	}
	namespace, err := w.GetNamespaceWithName(ctx, in.Spec.Name)
	if err != nil {
		return nil, err
	}
	out, err := w.reconcileNamespace(ctx, in.Spec, namespace)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Reconcile multiple namespaces, create missing namespaces, update existing namespaces, and optionally delete unaccounted namespaces.
func (w *workflows) ReconcileNamespaces(ctx workflow.Context, in *ReconcileNamespacesInput) (*ReconcileNamespacesOutput, error) {
	if err := validator.ValidateStruct(in); err != nil {
		return nil, fmt.Errorf("invalid input: %s", err)
	}
	namespaces, err := w.GetAllNamespaces(ctx)
	if err != nil {
		return nil, err
	}
	out := &ReconcileNamespacesOutput{}
	for i := range in.Specs {
		var namespace *namespace.Namespace
		for _, ns := range namespaces {
			if ns.Spec.Name == in.Specs[i].Name {
				namespace = ns
				namespaces = append(namespaces[:i], namespaces[i+1:]...) // remove the namespace from the list
				break
			}
		}
		// reconcile the namespace
		o, _ := w.reconcileNamespace(ctx, in.Specs[i], namespace)
		out.Results = append(out.Results, o)
	}
	// whats left in maps is only the unaccounted namespaces
	for _, ns := range namespaces {
		if in.DeleteUnaccounted {
			o := &ReconcileNamespaceOutput{
				Namespace: ns,
				Outcome:   ReconcileOutcomeDeleted,
			}
			_, err := w.DeleteNamespace(ctx, &cloudservice.DeleteNamespaceRequest{
				Namespace:       ns.Namespace,
				ResourceVersion: ns.ResourceVersion,
			})
			if err != nil {
				o.setError(err)
			}
			out.Results = append(out.Results, o)
		}
	}
	return out, nil
}
