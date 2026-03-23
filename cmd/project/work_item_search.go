package project

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"lark_cli/internal/cmddeps"
	"lark_cli/internal/openapi"
	"lark_cli/internal/query"

	"github.com/spf13/cobra"
)

func NewWorkItemSearchCmd(deps cmddeps.Deps) *cobra.Command {
	var projectKey string
	var workItemTypeKey string
	var persons []string
	var statuses []string
	var createdFrom string
	var createdTo string
	var updatedFrom string
	var updatedTo string
	var fields []string
	var fieldsOut []string
	var pageSize int64
	var pageNum int64
	var me bool
	var searchGroupJSON string
	var rawOnly bool
	var dryRun bool
	var jsonOut bool
	var tuiMode bool

	cmd := &cobra.Command{
		Use:          "search",
		Short:        "Search work items with structured filters or raw search_group JSON",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if deps.PluginTokenProvider == nil {
				return fmt.Errorf("not logged in: run 'lark login' first")
			}
			if projectKey == "" {
				return fmt.Errorf("--project-key is required")
			}
			if workItemTypeKey == "" {
				return fmt.Errorf("--work-item-type-key is required")
			}
			if tuiMode {
				return fmt.Errorf("--tui is not implemented for work-item search yet")
			}

			client := openapi.NewClient(deps.Config.BaseURL, &http.Client{Timeout: deps.Config.HTTPTimeout}, deps.PluginTokenProvider)
			return runWorkItemSearch(ctx, deps.Stdout, client, workItemSearchInput{
				currentUserKey:   deps.Config.UserKey,
				projectKey:       projectKey,
				workItemTypeKey:  workItemTypeKey,
				persons:          persons,
				statuses:         statuses,
				createdFrom:      createdFrom,
				createdTo:        createdTo,
				updatedFrom:      updatedFrom,
				updatedTo:        updatedTo,
				fields:           fields,
				fieldsOut:        fieldsOut,
				pageSize:         pageSize,
				pageNum:          pageNum,
				me:               me,
				searchGroupJSON:  searchGroupJSON,
				rawOnly:          rawOnly,
				dryRun:           dryRun,
				jsonOut:          jsonOut,
			})
		},
	}

	cmd.Flags().StringVarP(&projectKey, "project-key", "k", "", "Project key (required)")
	cmd.Flags().StringVarP(&workItemTypeKey, "work-item-type-key", "w", "", "Work item type key (required)")
	cmd.Flags().BoolVar(&me, "me", false, "Filter work items related to current user")
	cmd.Flags().StringArrayVar(&persons, "person", nil, "Filter by person user_key (repeatable)")
	cmd.Flags().StringArrayVar(&statuses, "status", nil, "Filter by work item status value (repeatable)")
	cmd.Flags().StringVar(&createdFrom, "created-from", "", "Created time start (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().StringVar(&createdTo, "created-to", "", "Created time end (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().StringVar(&updatedFrom, "updated-from", "", "Updated time start (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().StringVar(&updatedTo, "updated-to", "", "Updated time end (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().StringArrayVar(&fields, "field", nil, "Custom field filter in key=value format (repeatable)")
	cmd.Flags().StringSliceVar(&fieldsOut, "fields", nil, "Returned field keys (comma-separated or repeatable)")
	cmd.Flags().StringVar(&searchGroupJSON, "search-group-json", "", "Raw search_group JSON")
	cmd.Flags().BoolVar(&rawOnly, "raw-only", false, "Use only --search-group-json and ignore compiled filters")
	cmd.Flags().BoolVar(&dryRun, "dry-run-query", false, "Print request payload and do not execute")
	cmd.Flags().Int64Var(&pageSize, "page-size", 20, "Page size")
	cmd.Flags().Int64Var(&pageNum, "page-num", 1, "Page number starting from 1")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Print raw JSON output")
	cmd.Flags().BoolVarP(&tuiMode, "tui", "t", false, "Launch interactive TUI (not implemented yet)")

	return cmd
}

type workItemSearchInput struct {
	currentUserKey  string
	projectKey      string
	workItemTypeKey string
	persons         []string
	statuses        []string
	createdFrom     string
	createdTo       string
	updatedFrom     string
	updatedTo       string
	fields          []string
	fieldsOut       []string
	pageSize        int64
	pageNum         int64
	me              bool
	searchGroupJSON string
	rawOnly         bool
	dryRun          bool
	jsonOut         bool
}

func validateFieldsSelection(fields []string) error {
	if len(fields) == 0 {
		return nil
	}
	hasInclude := false
	hasExclude := false
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f == "" {
			return fmt.Errorf("--fields cannot contain empty field key")
		}
		if strings.HasPrefix(f, "-") {
			hasExclude = true
			continue
		}
		hasInclude = true
	}
	if hasInclude && hasExclude {
		return fmt.Errorf("--fields cannot mix include and exclude modes")
	}
	return nil
}
func runWorkItemSearch(ctx context.Context, out io.Writer, client *openapi.Client, in workItemSearchInput) error {
	built, err := query.BuildSearchGroup(query.BuildInput{
		CurrentUserKey:    in.currentUserKey,
		Me:                in.me,
		Persons:           in.persons,
		Statuses:          in.statuses,
		CreatedFrom:       in.createdFrom,
		CreatedTo:         in.createdTo,
		UpdatedFrom:       in.updatedFrom,
		UpdatedTo:         in.updatedTo,
		Fields:            in.fields,
		RawSearchGroupJSON: in.searchGroupJSON,
		RawOnly:           in.rawOnly,
	})
	if err != nil {
		return err
	}

	if in.pageSize <= 0 || in.pageSize > 50 {
		return fmt.Errorf("--page-size must be between 1 and 50")
	}
	if in.pageNum <= 0 {
		return fmt.Errorf("--page-num must be greater than 0")
	}
	if err := validateFieldsSelection(in.fieldsOut); err != nil {
		return err
	}

	payload := map[string]any{
		"search_group": built.SearchGroup,
	}
	payload["page_size"] = in.pageSize
	payload["page_num"] = in.pageNum
	if len(in.fieldsOut) > 0 {
		payload["fields"] = in.fieldsOut
	}

	if in.dryRun {
		if len(built.Warnings) > 0 {
			for _, w := range built.Warnings {
				fmt.Fprintf(os.Stderr, "warning: %s\n", w)
			}
		}
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(payload)
	}

	resp, err := client.SearchWorkItems(ctx, in.projectKey, in.workItemTypeKey, payload)
	if err != nil {
		return fmt.Errorf("failed to search work items: %w", err)
	}

	if in.jsonOut {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(resp)
	}

	if len(resp.Data) == 0 {
		fmt.Fprintln(out, "No work items found.")
		return nil
	}

	fmt.Fprintf(out, "Project: %s  Type: %s\n", in.projectKey, in.workItemTypeKey)
	for i, row := range resp.Data {
		id := row["id"]
		name := row["name"]
		if name == nil || fmt.Sprint(name) == "" {
			name = "(unnamed)"
		}
		fmt.Fprintf(out, "%d. [%v] %v\n", i+1, id, name)
	}
	return nil
}
