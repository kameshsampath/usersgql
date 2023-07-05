package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/kameshsampath/usersgql/ent"

	_ "github.com/mattn/go-sqlite3"
)

func Example_AddUser() {
	// Open the Connection to the Database
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	if err != nil {
		log.Fatalf("error opening DB connection,%v", err)
	}
	// Create Schema with migration support
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Schema.Create(ctx, schema.WithGlobalUniqueID(true)); err != nil {
		log.Fatalf("error creating/migrating schema,%v", err)
	}
	//Create a user and verify
	u, err := createUser(ctx, client)
	if err != nil {
		log.Fatalf("error creating user,%v", err)
	}
	fmt.Printf("%d: %s/%s", u.ID, u.Name, u.Gender)
	//Output:
	//1: Tom/M
}

func Example_QueryUser() {
	// Open the Connection to the Database
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	if err != nil {
		log.Fatalf("error opening DB connection,%v", err)
	}
	// Create Schema with migration support
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Schema.Create(ctx, schema.WithGlobalUniqueID(true)); err != nil {
		log.Fatalf("error creating/migrating schema,%v", err)
	}
	_, err = createUser(ctx, client)
	if err != nil {
		log.Fatalf("error creating user,%v", err)
	}
	//Query all users
	users := client.User.Query().AllX(ctx)
	for _, u := range users {
		fmt.Printf("%d: %s", u.ID, u.Name)
	}
	//Output:
	//1: Tom
}

func createUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	return client.User.Create().
		SetName("Tom").
		SetGender("M").
		Save(ctx)
}
