package activities

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	v1 "github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/auth/v1"
	"github.com/temporalio/cloud-operations-workflows/protogen/temporal/api/cloud/user/v1"
	"testing"
)

func Test(t *testing.T) {
	req := GetUserSpecsFromFileResponse{
		Specs: []*user.UserSpec{
			{
				Email: "bobadmin@example.com",
				Access: &user.UserAccess{
					AccountAccess: &v1.AccountAccess{
						Role: "admin",
					},
					NamespaceAccesses: map[string]*v1.NamespaceAccess{
						"demo-cloud-ops.temporal-dev": {Permission: "admin"},
					},
				},
			},
		},
	}

	bytes, err := json.Marshal(req)
	assert.NoError(t, err)
	fmt.Println(string(bytes))

}
