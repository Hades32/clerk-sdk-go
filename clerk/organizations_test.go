package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationsService_ListAll_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := fmt.Sprintf(`{
		"data": [%s],
		"total_count": 1
	}`, dummyOrganizationJson)

	mux.HandleFunc("/organizations", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	var want *OrganizationsResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Organizations().ListAll(ListAllOrganizationsParams{})
	if len(got.Data) != len(want.Data) {
		t.Errorf("Was expecting %d organizations to be returned, instead got %d", len(want.Data), len(got.Data))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestOrganizationsService_ListAll_happyPathWithParameters(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := fmt.Sprintf(`{
		"data": [%s],
		"total_count": 1
	}`, dummyOrganizationJson)

	mux.HandleFunc("/organizations", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")

		actualQuery := req.URL.Query()
		expectedQuery := url.Values(map[string][]string{
			"limit":  {"5"},
			"offset": {"6"},
		})
		assert.Equal(t, expectedQuery, actualQuery)
		fmt.Fprint(w, expectedResponse)
	})

	var want *OrganizationsResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	limit := 5
	offset := 6
	got, _ := client.Organizations().ListAll(ListAllOrganizationsParams{
		Limit:  &limit,
		Offset: &offset,
	})
	if len(got.Data) != len(want.Data) {
		t.Errorf("Was expecting %d organizations to be returned, instead got %d", len(want.Data), len(got.Data))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestOrganizationsService_ListAll_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	organizations, err := client.Organizations().ListAll(ListAllOrganizationsParams{})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if organizations != nil {
		t.Errorf("Was not expecting any organizations to be returned, instead got %v", organizations)
	}
}

const dummyOrganizationJson = `{
        "object": "organization",
        "id": "org_1mebQggrD3xO5JfuHk7clQ94ysA",
        "name": "test-org",
        "slug": "org_slug",
        "created_at": 1610783813,
        "updated_at": 1610783813,
		"public_metadata": {
			"address": {
				"street": "Pennsylvania Avenue",
				"number": "1600"
			}
		},
		"private_metadata": {
			"app_id": 5
		}
    }`
