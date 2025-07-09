package github

import (
	"context"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/go-viper/mapstructure/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/shurcooL/githubv4"
)

// ListProjects lists projects for a given user or organization.
func ListProjects(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_projects",
			mcp.WithDescription(t("TOOL_LIST_PROJECTS_DESCRIPTION", "List Projects for a user or organization")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_LIST_PROJECTS_USER_TITLE", "List projects"),
				ReadOnlyHint: ToBoolPtr(true),
			}),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Owner login (user or organization)"),
			),
			mcp.WithString("owner_type",
				mcp.Description("Owner type"),
				mcp.Enum("user", "organization"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner, err := RequiredParam[string](req, "owner")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ownerType, err := OptionalParam[string](req, "owner_type")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			if ownerType == "" {
				ownerType = "organization"
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if ownerType == "user" {
				var q struct {
					User struct {
						Projects struct {
							Nodes []struct {
								ID     githubv4.ID
								Title  githubv4.String
								Number githubv4.Int
							}
						} `graphql:"projectsV2(first: 100)"`
					} `graphql:"user(login: $login)"`
				}
				if err := client.Query(ctx, &q, map[string]any{
					"login": githubv4.String(owner),
				}); err != nil {
					return mcp.NewToolResultError(err.Error()), nil
				}
				return MarshalledTextResult(q), nil
			}

			var q struct {
				Organization struct {
					Projects struct {
						Nodes []struct {
							ID     githubv4.ID
							Title  githubv4.String
							Number githubv4.Int
						}
					} `graphql:"projectsV2(first: 100)"`
				} `graphql:"organization(login: $login)"`
			}
			if err := client.Query(ctx, &q, map[string]any{
				"login": githubv4.String(owner),
			}); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(q), nil
		}
}

// GetProjectFields lists fields for a project.
func GetProjectFields(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_project_fields",
			mcp.WithDescription(t("TOOL_GET_PROJECT_FIELDS_DESCRIPTION", "Get fields for a project")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_GET_PROJECT_FIELDS_USER_TITLE", "Get project fields"),
				ReadOnlyHint: ToBoolPtr(true),
			}),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Owner login"),
			),
			mcp.WithString("owner_type",
				mcp.Description("Owner type"),
				mcp.Enum("user", "organization"),
			),
			mcp.WithNumber("number",
				mcp.Required(),
				mcp.Description("Project number"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner, err := RequiredParam[string](req, "owner")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			number, err := RequiredInt(req, "number")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ownerType, err := OptionalParam[string](req, "owner_type")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			if ownerType == "" {
				ownerType = "organization"
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if ownerType == "user" {
				var q struct {
					User struct {
						Project struct {
							Fields struct {
								Nodes []struct {
									ID       githubv4.ID
									Name     githubv4.String
									DataType githubv4.String
								}
							} `graphql:"fields(first: 100)"`
						} `graphql:"projectV2(number: $number)"`
					} `graphql:"user(login: $login)"`
				}
				if err := client.Query(ctx, &q, map[string]any{
					"login":  githubv4.String(owner),
					"number": githubv4.Int(number), // #nosec G115 safe narrowing
				}); err != nil {
					return mcp.NewToolResultError(err.Error()), nil
				}
				return MarshalledTextResult(q), nil
			}

			var q struct {
				Organization struct {
					Project struct {
						Fields struct {
							Nodes []struct {
								ID       githubv4.ID
								Name     githubv4.String
								DataType githubv4.String
							}
						} `graphql:"fields(first: 100)"`
					} `graphql:"projectV2(number: $number)"`
				} `graphql:"organization(login: $login)"`
			}
			if err := client.Query(ctx, &q, map[string]any{
				"login":  githubv4.String(owner),
				"number": githubv4.Int(number), // #nosec G115 safe narrowing
			}); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(q), nil
		}
}

// GetProjectItems lists items for a project.
func GetProjectItems(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_project_items",
			mcp.WithDescription(t("TOOL_GET_PROJECT_ITEMS_DESCRIPTION", "Get items for a project")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_GET_PROJECT_ITEMS_USER_TITLE", "Get project items"),
				ReadOnlyHint: ToBoolPtr(true),
			}),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Owner login"),
			),
			mcp.WithString("owner_type",
				mcp.Description("Owner type"),
				mcp.Enum("user", "organization"),
			),
			mcp.WithNumber("number",
				mcp.Required(),
				mcp.Description("Project number"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner, err := RequiredParam[string](req, "owner")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			number, err := RequiredInt(req, "number")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ownerType, err := OptionalParam[string](req, "owner_type")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			if ownerType == "" {
				ownerType = "organization"
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if ownerType == "user" {
				var q struct {
					User struct {
						Project struct {
							Items struct {
								Nodes []struct {
									ID githubv4.ID
								}
							} `graphql:"items(first: 100)"`
						} `graphql:"projectV2(number: $number)"`
					} `graphql:"user(login: $login)"`
				}
				if err := client.Query(ctx, &q, map[string]any{
					"login":  githubv4.String(owner),
					"number": githubv4.Int(number), // #nosec G115 safe narrowing
				}); err != nil {
					return mcp.NewToolResultError(err.Error()), nil
				}
				return MarshalledTextResult(q), nil
			}

			var q struct {
				Organization struct {
					Project struct {
						Items struct {
							Nodes []struct {
								ID githubv4.ID
							}
						} `graphql:"items(first: 100)"`
					} `graphql:"projectV2(number: $number)"`
				} `graphql:"organization(login: $login)"`
			}
			if err := client.Query(ctx, &q, map[string]any{
				"login":  githubv4.String(owner),
				"number": githubv4.Int(number), // #nosec G115 safe narrowing
			}); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(q), nil
		}
}

// CreateProjectIssue creates an issue in a repository and returns its ID.
func CreateProjectIssue(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("create_project_issue",
			mcp.WithDescription(t("TOOL_CREATE_PROJECT_ISSUE_DESCRIPTION", "Create a new issue")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_CREATE_PROJECT_ISSUE_USER_TITLE", "Create issue"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description(DescriptionRepositoryOwner),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description(DescriptionRepositoryName),
			),
			mcp.WithString("title",
				mcp.Required(),
				mcp.Description("Issue title"),
			),
			mcp.WithString("body",
				mcp.Description("Issue body"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var params struct {
				Owner string
				Repo  string
				Title string
				Body  string
			}
			if err := mapstructure.Decode(req.Params.Arguments, &params); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// Fetch repository ID
			var repoQ struct {
				Repository struct {
					ID githubv4.ID
				} `graphql:"repository(owner: $owner, name: $name)"`
			}
			if err := client.Query(ctx, &repoQ, map[string]any{
				"owner": githubv4.String(params.Owner),
				"name":  githubv4.String(params.Repo),
			}); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			input := githubv4.CreateIssueInput{
				RepositoryID: repoQ.Repository.ID,
				Title:        githubv4.String(params.Title),
			}
			if params.Body != "" {
				input.Body = githubv4.NewString(githubv4.String(params.Body))
			}

			var mut struct {
				CreateIssue struct {
					Issue struct {
						ID githubv4.ID
					}
				} `graphql:"createIssue(input: $input)"`
			}
			if err := client.Mutate(ctx, &mut, input, nil); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(mut), nil
		}
}

// AddIssueToProject adds an existing issue to a project.
func AddIssueToProject(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("add_issue_to_project",
			mcp.WithDescription(t("TOOL_ADD_ISSUE_TO_PROJECT_DESCRIPTION", "Add an issue to a project")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_ADD_ISSUE_TO_PROJECT_USER_TITLE", "Add issue to project"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("Project ID"),
			),
			mcp.WithString("issue_id",
				mcp.Required(),
				mcp.Description("Issue node ID"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectID, err := RequiredParam[string](req, "project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			issueID, err := RequiredParam[string](req, "issue_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			var mut struct {
				AddProjectV2ItemByID struct {
					Item struct {
						ID githubv4.ID
					}
				} `graphql:"addProjectV2ItemById(input: $input)"`
			}
			input := githubv4.AddProjectV2ItemByIdInput{
				ProjectID: githubv4.ID(projectID),
				ContentID: githubv4.ID(issueID),
			}
			if err := client.Mutate(ctx, &mut, input, nil); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(mut), nil
		}
}

// UpdateProjectItemField updates a field value on a project item.
func UpdateProjectItemField(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("update_project_item_field",
			mcp.WithDescription(t("TOOL_UPDATE_PROJECT_ITEM_FIELD_DESCRIPTION", "Update a project item field")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_UPDATE_PROJECT_ITEM_FIELD_USER_TITLE", "Update project item field"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("Project ID"),
			),
			mcp.WithString("item_id",
				mcp.Required(),
				mcp.Description("Item ID"),
			),
			mcp.WithString("field_id",
				mcp.Required(),
				mcp.Description("Field ID"),
			),
			mcp.WithString("text_value",
				mcp.Description("Text value"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectID, err := RequiredParam[string](req, "project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			itemID, err := RequiredParam[string](req, "item_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			fieldID, err := RequiredParam[string](req, "field_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			textValue, err := OptionalParam[string](req, "text_value")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			val := githubv4.ProjectV2FieldValue{}
			if textValue != "" {
				val.Text = githubv4.NewString(githubv4.String(textValue))
			}

			input := githubv4.UpdateProjectV2ItemFieldValueInput{
				ProjectID: githubv4.ID(projectID),
				ItemID:    githubv4.ID(itemID),
				FieldID:   githubv4.ID(fieldID),
				Value:     val,
			}
			var mut struct {
				UpdateProjectV2ItemFieldValue struct {
					Typename githubv4.String `graphql:"__typename"`
				} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
			}
			if err := client.Mutate(ctx, &mut, input, nil); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(mut), nil
		}
}

// CreateDraftIssue creates a draft issue directly in a project.
func CreateDraftIssue(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("create_draft_issue",
			mcp.WithDescription(t("TOOL_CREATE_DRAFT_ISSUE_DESCRIPTION", "Create a draft issue in a project")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_CREATE_DRAFT_ISSUE_USER_TITLE", "Create draft issue"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("Project ID"),
			),
			mcp.WithString("title",
				mcp.Required(),
				mcp.Description("Issue title"),
			),
			mcp.WithString("body",
				mcp.Description("Issue body"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectID, err := RequiredParam[string](req, "project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			title, err := RequiredParam[string](req, "title")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body, err := OptionalParam[string](req, "body")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			input := githubv4.AddProjectV2DraftIssueInput{
				ProjectID: githubv4.ID(projectID),
				Title:     githubv4.String(title),
			}
			if body != "" {
				input.Body = githubv4.NewString(githubv4.String(body))
			}

			var mut struct {
				AddProjectV2DraftIssue struct {
					Item struct {
						ID githubv4.ID
					}
				} `graphql:"addProjectV2DraftIssue(input: $input)"`
			}
			if err := client.Mutate(ctx, &mut, input, nil); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(mut), nil
		}
}

// DeleteProjectItem removes an item from a project.
func DeleteProjectItem(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("delete_project_item",
			mcp.WithDescription(t("TOOL_DELETE_PROJECT_ITEM_DESCRIPTION", "Delete a project item")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_DELETE_PROJECT_ITEM_USER_TITLE", "Delete project item"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("Project ID"),
			),
			mcp.WithString("item_id",
				mcp.Required(),
				mcp.Description("Item ID"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectID, err := RequiredParam[string](req, "project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			itemID, err := RequiredParam[string](req, "item_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			input := githubv4.DeleteProjectV2ItemInput{
				ProjectID: githubv4.ID(projectID),
				ItemID:    githubv4.ID(itemID),
			}
			var mut struct {
				DeleteProjectV2Item struct {
					Typename githubv4.String `graphql:"__typename"`
				} `graphql:"deleteProjectV2Item(input: $input)"`
			}
			if err := client.Mutate(ctx, &mut, input, nil); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(mut), nil
		}
}

// CreateProject creates a new Project V2 board.
func CreateProject(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("create_project",
			mcp.WithDescription(t("TOOL_CREATE_PROJECT_DESCRIPTION", "Create a new Project V2 board")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_CREATE_PROJECT_USER_TITLE", "Create project"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Owner login (user or organization)"),
			),
			mcp.WithString("owner_type",
				mcp.Description("Owner type"),
				mcp.Enum("user", "organization"),
			),
			mcp.WithString("title",
				mcp.Required(),
				mcp.Description("Project title"),
			),
			mcp.WithBoolean("public",
				mcp.Description("Whether the project should be public. Defaults to private (false)"),
			),
			mcp.WithString("short_description",
				mcp.Description("Short description for the project"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner, err := RequiredParam[string](req, "owner")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			title, err := RequiredParam[string](req, "title")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ownerType, err := OptionalParam[string](req, "owner_type")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			if ownerType == "" {
				ownerType = "organization"
			}
			// retrieve optional params but ignore unsupported ones to avoid unused variable lint
			_, _ = OptionalParam[bool](req, "public")
			_, _ = OptionalParam[string](req, "short_description")

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// Resolve owner node ID
			var ownerQuery struct {
				User struct {
					ID githubv4.ID
				} `graphql:"user(login: $login)"`
				Organization struct {
					ID githubv4.ID
				} `graphql:"organization(login: $login)"`
			}
			if err := client.Query(ctx, &ownerQuery, map[string]any{
				"login": githubv4.String(owner),
			}); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			var ownerID githubv4.ID
			if ownerType == "user" {
				ownerID = ownerQuery.User.ID
			} else {
				ownerID = ownerQuery.Organization.ID
			}
			if ownerID == "" {
				return mcp.NewToolResultError("failed to resolve owner ID"), nil
			}

			input := githubv4.CreateProjectV2Input{
				OwnerID: ownerID,
				Title:   githubv4.String(title),
			}
			// createProjectV2 does not currently accept shortDescription/public; set only title/owner

			var mut struct {
				CreateProjectV2 struct {
					Project struct {
						ID    githubv4.ID
						URL   githubv4.URI
						Title githubv4.String
					} `graphql:"projectV2"`
				} `graphql:"createProjectV2(input: $input)"`
			}
			if err := client.Mutate(ctx, &mut, input, nil); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(mut), nil
		}
}

// UpdateProject updates mutable attributes of a Project V2 board.
func UpdateProject(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("update_project",
			mcp.WithDescription(t("TOOL_UPDATE_PROJECT_DESCRIPTION", "Update an existing Project V2 board")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_UPDATE_PROJECT_USER_TITLE", "Update project"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("Project ID"),
			),
			mcp.WithString("title",
				mcp.Description("New title"),
			),
			mcp.WithString("short_description",
				mcp.Description("New short description"),
			),
			mcp.WithBoolean("public",
				mcp.Description("Set project visibility to public (true) or private (false)"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectID, err := RequiredParam[string](req, "project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			title, err := OptionalParam[string](req, "title")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			shortDesc, err := OptionalParam[string](req, "short_description")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			isPublic, err := OptionalParam[bool](req, "public")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			input := githubv4.UpdateProjectV2Input{
				ProjectID: githubv4.ID(projectID),
			}
			if title != "" {
				input.Title = githubv4.NewString(githubv4.String(title))
			}
			if shortDesc != "" {
				input.ShortDescription = githubv4.NewString(githubv4.String(shortDesc))
			}
			if _, ok := req.GetArguments()["public"]; ok {
				input.Public = githubv4.NewBoolean(githubv4.Boolean(isPublic))
			}

			var mut struct {
				UpdateProjectV2 struct {
					Project struct {
						ID    githubv4.ID
						URL   githubv4.URI
						Title githubv4.String
					} `graphql:"projectV2"`
				} `graphql:"updateProjectV2(input: $input)"`
			}
			if err := client.Mutate(ctx, &mut, input, nil); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(mut), nil
		}
}

// DeleteProject deletes a Project V2 board.
func DeleteProject(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("delete_project",
			mcp.WithDescription(t("TOOL_DELETE_PROJECT_DESCRIPTION", "Delete a Project V2 board")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_DELETE_PROJECT_USER_TITLE", "Delete project"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("Project ID"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectID, err := RequiredParam[string](req, "project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			input := githubv4.DeleteProjectV2Input{
				ProjectID: githubv4.ID(projectID),
			}
			var mut struct {
				DeleteProjectV2 struct {
					Typename githubv4.String `graphql:"__typename"`
				} `graphql:"deleteProjectV2(input: $input)"`
			}
			if err := client.Mutate(ctx, &mut, input, nil); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(mut), nil
		}
}

// UpdateProjectItem archives or unarchives a project item.
func UpdateProjectItem(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("update_project_item",
			mcp.WithDescription(t("TOOL_UPDATE_PROJECT_ITEM_DESCRIPTION", "Archive / unarchive a project item")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_UPDATE_PROJECT_ITEM_USER_TITLE", "Update project item"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("Project ID"),
			),
			mcp.WithString("item_id",
				mcp.Required(),
				mcp.Description("Item ID"),
			),
			mcp.WithBoolean("archived",
				mcp.Description("Whether the item should be archived (true) or unarchived (false)"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectID, err := RequiredParam[string](req, "project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			itemID, err := RequiredParam[string](req, "item_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			archived, err := OptionalParam[bool](req, "archived")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			var archivedPtr *githubv4.Boolean
			if _, ok := req.GetArguments()["archived"]; ok {
				b := githubv4.Boolean(archived)
				archivedPtr = &b
			}
			input := updateProjectV2ItemInput{
				ProjectID: githubv4.ID(projectID),
				ItemID:    githubv4.ID(itemID),
				Archived:  archivedPtr,
			}

			var mut struct {
				UpdateProjectV2Item struct {
					Item struct {
						ID githubv4.ID
					}
				} `graphql:"updateProjectV2Item(input: $input)"`
			}
			if err := client.Mutate(ctx, &mut, input, nil); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(mut), nil
		}
}

// UpdateProjectItemPosition reorders an item within a project.
func UpdateProjectItemPosition(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("update_project_item_position",
			mcp.WithDescription(t("TOOL_UPDATE_PROJECT_ITEM_POSITION_DESCRIPTION", "Move a project item to a new position")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_UPDATE_PROJECT_ITEM_POSITION_USER_TITLE", "Move project item"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("Project ID"),
			),
			mcp.WithString("item_id",
				mcp.Required(),
				mcp.Description("Item ID to move"),
			),
			mcp.WithString("previous_item_id",
				mcp.Description("Item ID that should come directly before the moved item (optional)"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectID, err := RequiredParam[string](req, "project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			itemID, err := RequiredParam[string](req, "item_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			previousID, err := OptionalParam[string](req, "previous_item_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			var prevPtr *githubv4.ID
			if previousID != "" {
				idVal := githubv4.ID(previousID)
				prevPtr = &idVal
			}
			input := updateProjectV2ItemPositionInput{
				ProjectID:      githubv4.ID(projectID),
				ItemID:         githubv4.ID(itemID),
				PreviousItemID: prevPtr,
			}

			var mut struct {
				UpdateProjectV2ItemPosition struct {
					Item struct {
						ID githubv4.ID
					}
				} `graphql:"updateProjectV2ItemPosition(input: $input)"`
			}
			if err := client.Mutate(ctx, &mut, input, nil); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(mut), nil
		}
}

// ConvertProjectItemToIssue converts a draft-issue item into a real repository issue.
func ConvertProjectItemToIssue(getClient GetGQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("convert_project_item_to_issue",
			mcp.WithDescription(t("TOOL_CONVERT_PROJECT_ITEM_TO_ISSUE_DESCRIPTION", "Convert a draft item to a repository issue")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_CONVERT_PROJECT_ITEM_TO_ISSUE_USER_TITLE", "Convert item to issue"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("Project ID"),
			),
			mcp.WithString("item_id",
				mcp.Required(),
				mcp.Description("Item ID to convert"),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectID, err := RequiredParam[string](req, "project_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			itemID, err := RequiredParam[string](req, "item_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			input := convertProjectV2ItemToIssueInput{
				ProjectID: githubv4.ID(projectID),
				ItemID:    githubv4.ID(itemID),
			}
			var mut struct {
				ConvertProjectV2ItemToIssue struct {
					Issue struct {
						ID    githubv4.ID
						URL   githubv4.URI
						Title githubv4.String
					}
				} `graphql:"convertProjectV2ItemToIssue(input: $input)"`
			}
			if err := client.Mutate(ctx, &mut, input, nil); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return MarshalledTextResult(mut), nil
		}
}

// Custom input structs for GraphQL mutations that are not yet present in the githubv4 package.
// They embed githubv4.Input so the client recognises them as GraphQL input objects.

type updateProjectV2ItemInput struct {
	ProjectID githubv4.ID       `json:"projectId"`
	ItemID    githubv4.ID       `json:"itemId"`
	Archived  *githubv4.Boolean `json:"archived,omitempty"`
	githubv4.Input
}

type updateProjectV2ItemPositionInput struct {
	ProjectID      githubv4.ID  `json:"projectId"`
	ItemID         githubv4.ID  `json:"itemId"`
	PreviousItemID *githubv4.ID `json:"previousItemId,omitempty"`
	githubv4.Input
}

type convertProjectV2ItemToIssueInput struct {
	ProjectID githubv4.ID `json:"projectId"`
	ItemID    githubv4.ID `json:"itemId"`
	githubv4.Input
}
