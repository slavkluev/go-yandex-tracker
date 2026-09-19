package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleWorkflowsService_Get() {
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
}
