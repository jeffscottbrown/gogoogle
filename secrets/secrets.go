package secrets

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretspb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

var projectId = os.Getenv("PROJECT_ID")

// RetrieveSecret retrieves the secret from the secret manager.
// The secret is identified by the secret name.
// The PROJECT_ID environment variable must be set to the
// project id of the project where the secret is stored.
func RetrieveSecret(secretName string) (string, error) {
	if projectId == "" {
		slog.Error("PROJECT_ID environment variable not set")
		return "", fmt.Errorf("PROJECT_ID environment variable not set")
	}
	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secret manager client: %w", err)
	}
	defer client.Close()

	secretResource := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectId, secretName)

	req := &secretspb.AccessSecretVersionRequest{
		Name: secretResource,
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		slog.Error("Error retrieving client secret", "secretName", secretName, "error", err)
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	return string(result.Payload.Data), nil
}
