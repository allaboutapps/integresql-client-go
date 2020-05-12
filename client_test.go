package integresql

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"

	_ "github.com/lib/pq"
)

func TestClientInitializeTemplate(t *testing.T) {
	ctx := context.Background()

	c, err := DefaultClientFromEnv()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if err := c.ResetAllTracking(ctx); err != nil {
		t.Fatalf("failed to reset all test pool tracking: %v", err)
	}

	hash := "hashinghash1"

	template, err := c.InitializeTemplate(ctx, hash)
	if err != nil {
		t.Fatalf("failed to initialize template: %v", err)
	}

	if len(template.Config.Database) == 0 {
		t.Error("received invalid template database config")
	}
}

func TestClientDiscardTemplate(t *testing.T) {
	ctx := context.Background()

	c, err := DefaultClientFromEnv()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if err := c.ResetAllTracking(ctx); err != nil {
		t.Fatalf("failed to reset all test pool tracking: %v", err)
	}

	hash := "hashinghash2"

	if _, err := c.InitializeTemplate(ctx, hash); err != nil {
		t.Fatalf("failed to initialize template: %v", err)
	}

	if err := c.DiscardTemplate(ctx, hash); err != nil {
		t.Fatalf("failed to discard template: %v", err)
	}

	if _, err := c.InitializeTemplate(ctx, hash); err != nil {
		t.Fatalf("failed to reinitialize template: %v", err)
	}

	if err := c.FinalizeTemplate(ctx, hash); err != nil {
		t.Fatalf("failed to refinalize template: %v", err)
	}
}

func TestClientFinalizeTemplate(t *testing.T) {
	ctx := context.Background()

	c, err := DefaultClientFromEnv()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if err := c.ResetAllTracking(ctx); err != nil {
		t.Fatalf("failed to reset all test pool tracking: %v", err)
	}

	hash := "hashinghash2"

	if _, err := c.InitializeTemplate(ctx, hash); err != nil {
		t.Fatalf("failed to initialize template: %v", err)
	}

	if err := c.FinalizeTemplate(ctx, hash); err != nil {
		t.Fatalf("failed to finalize template: %v", err)
	}
}

func TestClientGetTestDatabase(t *testing.T) {
	ctx := context.Background()

	c, err := DefaultClientFromEnv()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if err := c.ResetAllTracking(ctx); err != nil {
		t.Fatalf("failed to reset all test pool tracking: %v", err)
	}

	hash := "hashinghash3"

	if _, err := c.InitializeTemplate(ctx, hash); err != nil {
		t.Fatalf("failed to initialize template: %v", err)
	}

	if err := c.FinalizeTemplate(ctx, hash); err != nil {
		t.Fatalf("failed to finalize template: %v", err)
	}

	test, err := c.GetTestDatabase(ctx, hash)
	if err != nil {
		t.Fatalf("failed to get test database: %v", err)
	}

	if test.TemplateHash != hash {
		t.Errorf("test database has invalid template hash, got %q, want %q", test.TemplateHash, hash)
	}

	db, err := sql.Open("postgres", test.Config.ConnectionString())
	if err != nil {
		t.Fatalf("failed to open test database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping test database connection: %v", err)
	}

	test2, err := c.GetTestDatabase(ctx, hash)
	if err != nil {
		t.Fatalf("failed to get second test database: %v", err)
	}

	if test2.TemplateHash != hash {
		t.Errorf("test database has invalid second template hash, got %q, want %q", test2.TemplateHash, hash)
	}

	if test2.ID == test.ID {
		t.Error("received same test database a second time without returning")
	}

	db2, err := sql.Open("postgres", test2.Config.ConnectionString())
	if err != nil {
		t.Fatalf("failed to open second test database connection: %v", err)
	}
	defer db2.Close()

	if err := db2.Ping(); err != nil {
		t.Fatalf("failed to ping second test database connection: %v", err)
	}
}

func TestClientReturnTestDatabase(t *testing.T) {
	ctx := context.Background()

	c, err := DefaultClientFromEnv()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if err := c.ResetAllTracking(ctx); err != nil {
		t.Fatalf("failed to reset all test pool tracking: %v", err)
	}

	hash := "hashinghash4"

	if _, err := c.InitializeTemplate(ctx, hash); err != nil {
		t.Fatalf("failed to initialize template: %v", err)
	}

	if err := c.FinalizeTemplate(ctx, hash); err != nil {
		t.Fatalf("failed to finalize template: %v", err)
	}

	test, err := c.GetTestDatabase(ctx, hash)
	if err != nil {
		t.Fatalf("failed to get test database: %v", err)
	}

	if test.TemplateHash != hash {
		t.Errorf("test database has invalid template hash, got %q, want %q", test.TemplateHash, hash)
	}

	if err := c.ReturnTestDatabase(ctx, hash, test.ID); err != nil {
		t.Fatalf("failed to return test database: %v", err)
	}

	test2, err := c.GetTestDatabase(ctx, hash)
	if err != nil {
		t.Fatalf("failed to get second test database: %v", err)
	}

	if test2.TemplateHash != hash {
		t.Errorf("test database has invalid second template hash, got %q, want %q", test2.TemplateHash, hash)
	}

	if test2.ID != test.ID {
		t.Errorf("received invalid test database, want %d, got %d", test.ID, test2.ID)
	}
}

func populateTemplateDB(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `
		CREATE EXTENSION "uuid-ossp";
		CREATE TABLE pilots (
			id uuid NOT NULL DEFAULT uuid_generate_v4(),
			"name" text NOT NULL,
			created_at timestamptz NOT NULL,
			updated_at timestamptz NULL,
			CONSTRAINT pilot_pkey PRIMARY KEY (id)
		);
		CREATE TABLE jets (
			id uuid NOT NULL DEFAULT uuid_generate_v4(),
			pilot_id uuid NOT NULL,
			age int4 NOT NULL,
			"name" text NOT NULL,
			color text NOT NULL,
			created_at timestamptz NOT NULL,
			updated_at timestamptz NULL,
			CONSTRAINT jet_pkey PRIMARY KEY (id)
		);
		ALTER TABLE jets ADD CONSTRAINT jet_pilots_fkey FOREIGN KEY (pilot_id) REFERENCES pilots(id);
	`); err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO pilots (id, "name", created_at, updated_at) VALUES ('744a1a87-5ef7-4309-8814-0f1054751156', 'Mario', '2020-03-23 09:44:00.548', '2020-03-23 09:44:00.548');
		INSERT INTO pilots (id, "name", created_at, updated_at) VALUES ('20d9d155-2e95-49a2-8889-2ae975a8617e', 'Nick', '2020-03-23 09:44:00.548', '2020-03-23 09:44:00.548');
		INSERT INTO jets (id, pilot_id, age, "name", color, created_at, updated_at) VALUES ('67d9d0c7-34e5-48b0-9c7d-c6344995353c', '744a1a87-5ef7-4309-8814-0f1054751156', 26, 'F-14B', 'grey', '2020-03-23 09:44:00.000', '2020-03-23 09:44:00.000');
		INSERT INTO jets (id, pilot_id, age, "name", color, created_at, updated_at) VALUES ('facaf791-21b4-401a-bbac-67079ae4921f', '20d9d155-2e95-49a2-8889-2ae975a8617e', 27, 'F-14B', 'grey/red', '2020-03-23 09:44:00.000', '2020-03-23 09:44:00.000');
	`); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func TestSetupTemplateWithDBClient(t *testing.T) {
	ctx := context.Background()

	c, err := DefaultClientFromEnv()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if err := c.ResetAllTracking(ctx); err != nil {
		t.Fatalf("failed to reset all test pool tracking: %v", err)
	}

	hash := "hashinghash5"

	if err := c.SetupTemplateWithDBClient(ctx, hash, func(db *sql.DB) error {

		// setup code
		if err := populateTemplateDB(ctx, db); err != nil {
			t.Fatalf("failed to populate template db: %v", err)
		}

		return err
	}); err != nil {
		t.Fatalf("Failed to setup template database for hash %q: %v", hash, err)
	}
}

func getTestDB(wg *sync.WaitGroup, errs chan<- error, c *Client, hash string) {
	defer wg.Done()

	_, err := c.GetTestDatabase(context.Background(), hash)
	if err != nil {
		errs <- err
		return
	}

	errs <- nil
}

func TestSetupTemplateWithDBClientFailingSetupCodeAndReinitialize(t *testing.T) {
	ctx := context.Background()

	c, err := DefaultClientFromEnv()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if err := c.ResetAllTracking(ctx); err != nil {
		t.Fatalf("failed to reset all test pool tracking: %v", err)
	}

	hash := "hashinghash5"

	testDBWhileInitializeCount := 5
	testDBPreDiscardCount := 5
	testDBAfterDiscardCount := 5

	allTestDbCount := testDBWhileInitializeCount + testDBPreDiscardCount + testDBAfterDiscardCount

	var errs = make(chan error, allTestDbCount)
	var wg sync.WaitGroup

	if err := c.SetupTemplateWithDBClient(ctx, hash, func(db *sql.DB) error {

		// setup code
		if err := populateTemplateDB(ctx, db); err != nil {
			t.Fatalf("failed to populate template db: %v", err)
		}

		// some other participents are already connecting...
		wg.Add(testDBWhileInitializeCount)

		for i := 0; i < testDBWhileInitializeCount; i++ {
			go getTestDB(&wg, errs, c, hash)
		}

		// but then we throw an error during our test setup!
		err = errors.New("FAILED ERR DURING INITIALIZE")

		return err
	}); err == nil {
		t.Fatalf("we expected this to error!!")
	}

	// some other participents are still want to be part of this...
	wg.Add(testDBPreDiscardCount)

	for i := 0; i < testDBPreDiscardCount; i++ {
		go getTestDB(&wg, errs, c, hash)
	}

	// SIGNAL DISCARD!
	err = c.DiscardTemplate(ctx, hash)

	if err != nil {
		t.Fatalf("failed to discard template database after error during initialize: %v", err)
	}

	// finalize template should now no longer work
	err = c.FinalizeTemplate(ctx, hash)

	if err == nil {
		t.Fatalf("finalize template should not work after a successful discard!: %v", err)
	}

	// haven't learned other participents are still want to be part of this...
	wg.Add(testDBAfterDiscardCount)

	for i := 0; i < testDBAfterDiscardCount; i++ {
		go getTestDB(&wg, errs, c, hash)
	}

	// wait for all the errors with getTestDB to arrive...
	wg.Wait()

	var results = make([]error, 0, allTestDbCount)
	for i := 0; i < allTestDbCount; i++ {
		results = append(results, <-errs)
	}

	close(errs)

	// check all getTestDatabase clients also errored out...
	success := 0
	errored := 0
	for _, err := range results {
		if err == nil {
			success++
		} else {
			// fmt.Println(err)
			errored++
		}
	}

	if errored != allTestDbCount {
		t.Errorf("invalid number of errored retrievals, got %d, want %d", errored, allTestDbCount)
	}

	if success != 0 {
		t.Errorf("invalid number of successful retrievals, got %d, want %d", success, 0)
	}

	// then test a successful reinitialize...
	if err := c.SetupTemplateWithDBClient(ctx, hash, func(db *sql.DB) error {

		errr := populateTemplateDB(ctx, db)

		// setup code
		if errr != nil {
			t.Fatalf("failed to repopulate template db: %v", errr)
		}

		return errr
	}); err != nil {
		t.Fatalf("Failed to resetup template database for hash %q: %v", hash, err)
	}
}
