package main

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/databricks/terraform-provider-databricks/common"
	"github.com/databricks/terraform-provider-databricks/qa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	qa.HTTPFixturesApply(t, []qa.HTTPFixture{
		{
			MatchAny:     true,
			ReuseRequest: true,
			Status:       200,
			Response:     `{}`,
		},
	}, func(ctx context.Context, client *common.DatabricksClient) {
		responseWriter := httptest.NewRecorder()
		azure.PublicCloud.ResourceManagerEndpoint = client.Config.Host
		t.Setenv("MSI_ENDPOINT", client.Config.Host)
		t.Setenv("MSI_SECRET", "secret")
		t.Setenv("ACI_CONTAINER_GROUP", "")
		triggerStart(responseWriter, nil)
		assert.Equal(t, "400 Bad Request", responseWriter.Result().Status)

		responseWriter = httptest.NewRecorder()
		t.Setenv("ACI_CONTAINER_GROUP", "/abc")
		triggerStart(responseWriter, nil)
		assert.Equal(t, "200 OK", responseWriter.Result().Status)

		// test that app properly fails
		t.Setenv("FUNCTIONS_CUSTOMHANDLER_PORT", "abc")
		defer func() {
			err := recover()
			require.NotNil(t, err)
		}()
		main()
	})
}
