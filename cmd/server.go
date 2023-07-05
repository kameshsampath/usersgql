package main

import (
	"context"
	"log"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kameshsampath/usersgql/ent"
	"github.com/kameshsampath/usersgql/resolvers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	if err != nil {
		log.Fatalf("error opening DB connection,%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Schema.Create(ctx, schema.WithGlobalUniqueID(true)); err != nil {
		log.Fatalf("error creating/migrating schema,%v", err)
	}

	if err := loadData(ctx, client); err != nil {
		log.Printf("Error loading sample data:%v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", playgroundHandler())
	e.POST("/query", graphqlHandler(client))

	if err := e.Start(":8081"); err != nil {
		log.Fatalf("error starting GraphQL server,%v", err)
	}
}

func playgroundHandler() echo.HandlerFunc {
	h := playground.Handler("Users", "/query")
	return echo.WrapHandler(h)
}

func graphqlHandler(client *ent.Client) echo.HandlerFunc {
	server := handler.NewDefaultServer(resolvers.NewSchema(client))
	return echo.WrapHandler(server)
}

func loadData(ctx context.Context, client *ent.Client) error {
	// fruits data as name:season format
	// this data will be loaded as sample users data
	fruitsData := []string{
		"Tom:M",
		"Jerry:M",
		"Minnie:F",
		"Mickey:M",
	}
	fc := make([]*ent.UserCreate, len(fruitsData))
	for i, fd := range fruitsData {
		d := strings.Split(fd, ":")
		fc[i] = client.User.Create().SetName(d[0]).SetGender(d[1])
	}
	_, err := client.User.CreateBulk(fc...).Save(ctx)

	if err != nil {
		return err
	}

	log.Println("Sample data successfully loaded")

	return nil
}
