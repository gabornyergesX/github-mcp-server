package github

import (
	"context"
	"testing"

	"github.com/github/github-mcp-server/internal/githubv4mock"
	"github.com/github/github-mcp-server/internal/toolsnaps"
	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateProject(t *testing.T) {
	// Mock owner lookup + mutation
	mockClient := githubv4mock.NewMockedHTTPClient(
		githubv4mock.NewQueryMatcher(
			struct {
				User         struct{ ID githubv4.ID } `graphql:"user(login: $login)"`
				Organization struct{ ID githubv4.ID } `graphql:"organization(login: $login)"`
			}{},
			map[string]any{"login": githubv4.String("acme")},
			githubv4mock.DataResponse(map[string]any{"user": map[string]any{"id": "U_123"}}),
		),
		githubv4mock.NewMutationMatcher(
			struct {
				CreateProjectV2 struct {
					Project struct{ ID githubv4.ID } `graphql:"projectV2"`
				} `graphql:"createProjectV2(input: $input)"`
			}{},
			githubv4.CreateProjectV2Input{OwnerID: "U_123", Title: githubv4.String("Demo Board")},
			nil,
			githubv4mock.DataResponse(map[string]any{"createProjectV2": map[string]any{"projectV2": map[string]any{"id": "PVT_1"}}}),
		),
	)

	tool, handler := CreateProject(stubGetGQLClientFn(githubv4.NewClient(mockClient)), translations.NullTranslationHelper)
	require.NoError(t, toolsnaps.Test(tool.Name, tool))

	res, err := handler(context.Background(), createMCPRequest(map[string]any{
		"owner":      "acme",
		"owner_type": "user",
		"title":      "Demo Board",
	}))
	require.NoError(t, err)
	assert.NotNil(t, res)
}

func Test_UpdateProject(t *testing.T) {
	titlePtr := githubv4.NewString(githubv4.String("Updated Title"))
	mockClient := githubv4mock.NewMockedHTTPClient(
		githubv4mock.NewMutationMatcher(
			struct {
				UpdateProjectV2 struct {
					Project struct{ ID githubv4.ID } `graphql:"projectV2"`
				} `graphql:"updateProjectV2(input: $input)"`
			}{},
			githubv4.UpdateProjectV2Input{ProjectID: "P_1", Title: titlePtr},
			nil,
			githubv4mock.DataResponse(map[string]any{"updateProjectV2": map[string]any{"projectV2": map[string]any{"id": "P_1"}}}),
		),
	)

	tool, handler := UpdateProject(stubGetGQLClientFn(githubv4.NewClient(mockClient)), translations.NullTranslationHelper)
	require.NoError(t, toolsnaps.Test(tool.Name, tool))

	res, err := handler(context.Background(), createMCPRequest(map[string]any{
		"project_id": "P_1",
		"title":      "Updated Title",
	}))
	require.NoError(t, err)
	assert.NotNil(t, res)
}

func Test_DeleteProject(t *testing.T) {
	mockClient := githubv4mock.NewMockedHTTPClient(
		githubv4mock.NewMutationMatcher(
			struct {
				DeleteProjectV2 struct {
					Typename githubv4.String `graphql:"__typename"`
				} `graphql:"deleteProjectV2(input: $input)"`
			}{},
			githubv4.DeleteProjectV2Input{ProjectID: "P_1"},
			nil,
			githubv4mock.DataResponse(map[string]any{"deleteProjectV2": map[string]any{"__typename": "ProjectV2"}}),
		),
	)

	tool, handler := DeleteProject(stubGetGQLClientFn(githubv4.NewClient(mockClient)), translations.NullTranslationHelper)
	require.NoError(t, toolsnaps.Test(tool.Name, tool))

	res, err := handler(context.Background(), createMCPRequest(map[string]any{"project_id": "P_1"}))
	require.NoError(t, err)
	assert.NotNil(t, res)
}

func Test_UpdateProjectItem(t *testing.T) {
	archivedBool := githubv4.Boolean(true)
	input := updateProjectV2ItemInput{ProjectID: "P_1", ItemID: "I_1", Archived: &archivedBool}

	mockClient := githubv4mock.NewMockedHTTPClient(
		githubv4mock.NewMutationMatcher(
			struct {
				UpdateProjectV2Item struct {
					Item struct{ ID githubv4.ID }
				} `graphql:"updateProjectV2Item(input: $input)"`
			}{},
			input,
			nil,
			githubv4mock.DataResponse(map[string]any{"updateProjectV2Item": map[string]any{"item": map[string]any{"id": "I_1"}}}),
		),
	)

	tool, handler := UpdateProjectItem(stubGetGQLClientFn(githubv4.NewClient(mockClient)), translations.NullTranslationHelper)
	require.NoError(t, toolsnaps.Test(tool.Name, tool))

	res, err := handler(context.Background(), createMCPRequest(map[string]any{
		"project_id": "P_1",
		"item_id":    "I_1",
		"archived":   true,
	}))
	require.NoError(t, err)
	assert.NotNil(t, res)
}

func Test_UpdateProjectItemPosition(t *testing.T) {
	input := updateProjectV2ItemPositionInput{ProjectID: "P_1", ItemID: "I_2"}
	mockClient := githubv4mock.NewMockedHTTPClient(
		githubv4mock.NewMutationMatcher(
			struct {
				UpdateProjectV2ItemPosition struct {
					Item struct{ ID githubv4.ID }
				} `graphql:"updateProjectV2ItemPosition(input: $input)"`
			}{},
			input,
			nil,
			githubv4mock.DataResponse(map[string]any{"updateProjectV2ItemPosition": map[string]any{"item": map[string]any{"id": "I_2"}}}),
		),
	)

	tool, handler := UpdateProjectItemPosition(stubGetGQLClientFn(githubv4.NewClient(mockClient)), translations.NullTranslationHelper)
	require.NoError(t, toolsnaps.Test(tool.Name, tool))

	res, err := handler(context.Background(), createMCPRequest(map[string]any{
		"project_id": "P_1",
		"item_id":    "I_2",
	}))
	require.NoError(t, err)
	assert.NotNil(t, res)
}

func Test_ConvertProjectItemToIssue(t *testing.T) {
	input := convertProjectV2ItemToIssueInput{ProjectID: "P_1", ItemID: "I_3"}
	mockClient := githubv4mock.NewMockedHTTPClient(
		githubv4mock.NewMutationMatcher(
			struct {
				ConvertProjectV2ItemToIssue struct {
					Issue struct{ ID githubv4.ID }
				} `graphql:"convertProjectV2ItemToIssue(input: $input)"`
			}{},
			input,
			nil,
			githubv4mock.DataResponse(map[string]any{"convertProjectV2ItemToIssue": map[string]any{"issue": map[string]any{"id": "ISS_1"}}}),
		),
	)

	tool, handler := ConvertProjectItemToIssue(stubGetGQLClientFn(githubv4.NewClient(mockClient)), translations.NullTranslationHelper)
	require.NoError(t, toolsnaps.Test(tool.Name, tool))

	res, err := handler(context.Background(), createMCPRequest(map[string]any{
		"project_id": "P_1",
		"item_id":    "I_3",
	}))
	require.NoError(t, err)
	assert.NotNil(t, res)
}
