package database

import (
	"BTgoAWSlambda/types"
	"context"
	"testing"
)

// TestMockUserStore tests the mock implementation
func TestMockUserStore(t *testing.T) {
	mockStore := NewMockUserStore()
	ctx := context.Background()

	// Test user
	user := types.User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}

	// Test DoesUserExist - should return false initially
	exists, err := mockStore.DoesUserExist(ctx, "testuser")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if exists {
		t.Errorf("User should not exist initially")
	}

	// Test InsertUser
	err = mockStore.InsertUser(ctx, user)
	if err != nil {
		t.Errorf("Failed to insert user: %v", err)
	}

	// Test DoesUserExist - should return true now
	exists, err = mockStore.DoesUserExist(ctx, "testuser")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !exists {
		t.Errorf("User should exist after insertion")
	}

	// Test GetUser
	retrievedUser, err := mockStore.GetUser(ctx, "testuser")
	if err != nil {
		t.Errorf("Failed to get user: %v", err)
	}
	if retrievedUser.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrievedUser.Username)
	}
	if retrievedUser.PasswordHash != user.PasswordHash {
		t.Errorf("Expected password hash %s, got %s", user.PasswordHash, retrievedUser.PasswordHash)
	}

	// Test GetUser for non-existent user
	_, err = mockStore.GetUser(ctx, "nonexistent")
	if err == nil {
		t.Errorf("Expected error for non-existent user")
	}
}

// TestDynamoDBClient tests the actual DynamoDB client (requires AWS credentials)
func TestDynamoDBClient(t *testing.T) {
	// Skip if running in CI or without AWS credentials
	if testing.Short() {
		t.Skip("Skipping DynamoDB integration test")
	}

	client := NewDynamoDB()
	ctx := context.Background()

	// Test user
	user := types.User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}

	// Clean up before test
	client.DoesUserExist(ctx, "testuser")

	// Test DoesUserExist - should return false initially
	exists, err := client.DoesUserExist(ctx, "testuser")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if exists {
		t.Errorf("User should not exist initially")
	}

	// Test InsertUser
	err = client.InsertUser(ctx, user)
	if err != nil {
		t.Errorf("Failed to insert user: %v", err)
	}

	// Test DoesUserExist - should return true now
	exists, err = client.DoesUserExist(ctx, "testuser")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !exists {
		t.Errorf("User should exist after insertion")
	}

	// Test GetUser
	retrievedUser, err := client.GetUser(ctx, "testuser")
	if err != nil {
		t.Errorf("Failed to get user: %v", err)
	}
	if retrievedUser.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrievedUser.Username)
	}

	// Clean up after test
	// Note: In a real test, you might want to delete the test user
}
