package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleBoardsService_Edit() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	// The version parameter is used for optimistic locking via the If-Match header.
	// Get the current version from the board's Version field first.
	board, _, err := client.Boards.Edit(context.Background(), 1, 5, &tracker.BoardEditRequest{
		Name: tracker.Ptr("Updated Board Name"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*board.Name)
}

func ExampleBoardsService_CreateColumn() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	// Version for optimistic locking -- get from board's Version field.
	column, _, err := client.Boards.CreateColumn(context.Background(), 1, 5, &tracker.ColumnCreateRequest{
		Name: tracker.Ptr("In Progress"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*column.Name)
}
