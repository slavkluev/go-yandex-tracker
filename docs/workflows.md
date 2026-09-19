# Workflows

The Workflows service fetches a Yandex Tracker workflow: its steps (statuses) and the actions (transitions) available from each step. To find a queue's workflows, call `Queues.ListWorkflows`, which maps each workflow ID to its issue types.

## Usage

### Get a Workflow

```go
client := tracker.NewClient(
    tracker.WithOAuthToken("your-oauth-token"),
    tracker.WithOrgID("your-org-id"),
)

workflow, _, err := client.Workflows.Get(context.Background(), "W21")
if err != nil {
    log.Fatal(err)
}

for _, step := range workflow.Steps {
    for _, action := range step.Actions {
        fmt.Println(*step.Status.Key, "->", *action.Target.Key)
    }
}
```

`InitialAction` is the transition that sets a new issue's first status; its `Target` is that status. A terminal step has no actions, so its `Actions` is nil. Tracker omits optional fields that have no value, so `Queue`, `CreatedBy`, `UpdatedBy`, and `Type` can be nil.

## Methods

| Method | Description |
|--------|-------------|
| `Get` | Get a workflow by ID, with its steps and transitions |

## See Also

- [ExampleWorkflowsService_Get](https://pkg.go.dev/github.com/slavkluev/go-yandex-tracker/tracker#example-WorkflowsService.Get)
- [Queues](queues.md) -- `ListWorkflows`
- [Error Handling](errors.md)
