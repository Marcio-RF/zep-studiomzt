package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/getzep/zep/pkg/models"
	"github.com/getzep/zep/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

// TODO: Complete the tests

func TestCreateCollectionRoute(t *testing.T) {
	collectionName := testutils.GenerateRandomString(10)

	autoEmbeded := false

	// Create a collection
	collection := &models.CreateDocumentCollectionRequest{
		Name:        collectionName,
		Description: "Test collection",
		Metadata: map[string]interface{}{
			"key": "value",
		},
		EmbeddingDimensions: 128,
		IsAutoEmbedded:      &autoEmbeded,
	}

	// Convert collection to JSON
	collectionJSON, err := json.Marshal(collection)
	assert.NoError(t, err)

	// Create a request
	req, err := http.NewRequest(
		"POST",
		testServer.URL+"/api/v1/collection/"+collectionName,
		bytes.NewBuffer(collectionJSON),
	)
	assert.NoError(t, err)

	// Create a client and do the request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get the newly created collection
	rc, err := appState.DocumentStore.GetCollection(testCtx, collectionName)
	assert.NoError(t, err)

	assert.NotEmpty(t, rc.UUID)
	assert.Equal(t, rc.Name, strings.ToLower(collectionName))
	assert.Equal(t, rc.Metadata["key"], "value")
}

func TestUpdateCollectionHandler(t *testing.T) {
	collectionName := testutils.GenerateRandomString(10)

	autoEmbeded := false
	collectionCreateRequest := models.DocumentCollection{
		Name:        collectionName,
		Description: "Test collection",
		Metadata: map[string]interface{}{
			"key": "value",
		},
		EmbeddingDimensions: 128,
		IsAutoEmbedded:      autoEmbeded,
	}

	err := appState.DocumentStore.CreateCollection(testCtx, collectionCreateRequest)
	assert.NoError(t, err)

	// Update a collection
	collection := &models.UpdateDocumentCollectionRequest{
		Description: "Updated Test collection",
		Metadata: map[string]interface{}{
			"key": "updated value",
		},
	}

	// Convert collection to JSON
	collectionJSON, err := json.Marshal(collection)
	assert.NoError(t, err)

	// Create a request
	req, err := http.NewRequest(
		"PATCH",
		testServer.URL+"/api/v1/collection/"+collectionName,
		bytes.NewBuffer(collectionJSON),
	)
	assert.NoError(t, err)

	// Create a client and do the request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get the updated collection
	rc, err := appState.DocumentStore.GetCollection(testCtx, collectionName)
	assert.NoError(t, err)

	assert.Equal(t, rc.Description, "Updated Test collection")
	assert.Equal(t, rc.Metadata["key"], "updated value")
}

func TestDeleteCollectionHandler(t *testing.T) {
	collectionName := testutils.GenerateRandomString(10)
	// Create a collection
	cr := models.DocumentCollection{
		Name:        collectionName,
		Description: "Test collection",
		Metadata: map[string]interface{}{
			"key": "value",
		},
		EmbeddingDimensions: 128,
		IsAutoEmbedded:      false,
	}

	err := appState.DocumentStore.CreateCollection(testCtx, cr)
	assert.NoError(t, err)

	// Delete the collection
	req, err := http.NewRequest(
		"DELETE",
		testServer.URL+"/api/v1/collection/"+collectionName,
		nil,
	)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Try to get the deleted collection
	_, err = appState.DocumentStore.GetCollection(testCtx, collectionName)
	assert.ErrorAs(t, err, &models.ErrNotFound)
}

func TestGetCollectionHandler(t *testing.T) {
	collectionName := testutils.GenerateRandomString(10)
	// Create a collection
	cr := models.DocumentCollection{
		Name:        collectionName,
		Description: "Test collection",
		Metadata: map[string]interface{}{
			"key": "value",
		},
		EmbeddingDimensions: 128,
		IsAutoEmbedded:      false,
	}

	err := appState.DocumentStore.CreateCollection(testCtx, cr)
	assert.NoError(t, err)

	// Get the collection
	req, err := http.NewRequest(
		"GET",
		testServer.URL+"/api/v1/collection/"+collectionName,
		nil,
	)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get the collection
	rc, err := appState.DocumentStore.GetCollection(testCtx, collectionName)
	assert.NoError(t, err)

	assert.Equal(t, rc.Description, "Test collection")
	assert.Equal(t, rc.Metadata["key"], "value")
}

func TestCreateDocumentsHandler(t *testing.T) {
	collectionName := testutils.GenerateRandomString(10)
	// Create a collection
	cr := models.DocumentCollection{
		Name:        collectionName,
		Description: "Test collection",
		Metadata: map[string]interface{}{
			"key": "value",
		},
		EmbeddingDimensions: 128,
		IsAutoEmbedded:      true,
	}

	err := appState.DocumentStore.CreateCollection(testCtx, cr)
	assert.NoError(t, err)

	// Create documents
	docs := []models.CreateDocumentRequest{
		{
			DocumentID: "doc1",
			Content:    "This is a test document",
			Metadata: map[string]interface{}{
				"key": "value",
			},
		},
		{
			DocumentID: "doc2",
			Content:    "This is another test document",
			Metadata: map[string]interface{}{
				"key": "value",
			},
		},
	}

	j, err := json.Marshal(docs)
	assert.NoError(t, err)

	req, err := http.NewRequest(
		"POST",
		testServer.URL+"/api/v1/collection/"+collectionName+"/document",
		bytes.NewBuffer(j),
	)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get the documents
	for _, doc := range docs {
		rd, err := appState.DocumentStore.GetDocuments(
			testCtx,
			collectionName,
			nil,
			[]string{doc.DocumentID},
		)
		assert.NoError(t, err)

		assert.Equal(t, rd[0].Content, doc.Content)
		assert.Equal(t, rd[0].Metadata["key"], doc.Metadata["key"])
	}
}

// TestCreateDocumentsHandler with request body size greater than appState.Config.Server.MaxRequestSize
func TestCreateDocumentsHandler_MaxRequestBodySize(t *testing.T) {
	collectionName := testutils.GenerateRandomString(10)
	// Create a collection
	cr := models.DocumentCollection{
		Name:        collectionName,
		Description: "Test collection",
		Metadata: map[string]interface{}{
			"key": "value",
		},
		EmbeddingDimensions: 128,
		IsAutoEmbedded:      true,
	}

	err := appState.DocumentStore.CreateCollection(testCtx, cr)
	assert.NoError(t, err)

	// Create a large document
	largeDoc := strings.Repeat("a", int(appState.Config.Server.MaxRequestSize+1))

	// Create a document request with the large document
	docReq := []models.CreateDocumentRequest{
		{
			DocumentID: "largeDoc",
			Content:    largeDoc,
			Metadata: map[string]interface{}{
				"key": "value",
			},
		},
	}

	// Marshal the document request into JSON
	j, err := json.Marshal(docReq)
	assert.NoError(t, err)

	// Create a new HTTP request
	req, err := http.NewRequest(
		"POST",
		testServer.URL+"/api/v1/collection/"+collectionName+"/document",
		bytes.NewBuffer(j),
	)
	assert.NoError(t, err)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	// Check the status code
	assert.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)
	assert.Equal(t, "413 Request Entity Too Large", resp.Status)
}