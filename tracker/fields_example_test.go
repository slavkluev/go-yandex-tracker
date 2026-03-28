package tracker_test

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-yandex-tracker/tracker"
)

func ExampleFieldsService_List() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	fields, _, err := client.Fields.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, field := range fields {
		fmt.Println(*field.ID, *field.Name)
	}
}

func ExampleFieldsService_Create() {
	client := tracker.NewClient(
		tracker.WithOAuthToken("your-oauth-token"),
		tracker.WithOrgID("your-org-id"),
	)

	field, _, err := client.Fields.Create(context.Background(), &tracker.FieldCreateRequest{
		ID: tracker.Ptr("myField"),
		Name: &tracker.FieldName{
			EN: tracker.Ptr("My Field"),
			RU: tracker.Ptr("Мое поле"),
		},
		Category: tracker.Ptr("000000000000000000000002"),
		Type:     tracker.Ptr("ru.yandex.startrek.core.fields.StringFieldType"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*field.ID)
}
