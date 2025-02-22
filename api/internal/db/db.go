package db

import (
	"fmt"
	"log"
	"os"
	"pulse/internal/db/models"
	"time"

	charm_log "github.com/charmbracelet/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC",
		"localhost",
		"mariusihring",
		"password",
		"pulse",
		"5432",
	)

	// Configure GORM logger
	logLevel := logger.Info

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Log queries slower than 1 second
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true, // Enable prepared statement cache
	})
	if err != nil {
		charm_log.Fatalf("failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		charm_log.Fatalf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

	db.Exec("DEALLOCATE ALL")
	// Run migrations
	if err := autoMigrate(db); err != nil {
		charm_log.Fatalf("failed to run migrations: %v", err)
	}
	if err := seedRoles(db); err != nil {
		charm_log.Fatalf("failed to run seeding roles: %v", err)
	}
	charm_log.Info("Database connection established & Migrations run")

	return db
}

func autoMigrate(db *gorm.DB) error {
	db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Chain{},
		&models.TransactionCategory{},
		&models.Wallet{},
		&models.Subwallet{},
		&models.Token{},
		&models.SubwalletToken{},
		&models.HistoricalPrice{},
		&models.Transaction{},
		&models.Snapshot{},
		&models.TokenSnapshot{},
		&models.Alert{},           // depends on User and Token
		&models.PortfolioMetric{}, // depends on User
	)
	// Step 5: Create or replace the function to handle new user signups
	// props: {"title": "Create Trigger for Users Table", "runQuery": "false", "isChart": "false"}

	// Step 1: Create or replace the function to handle new user signups
	db.Exec(`
    CREATE OR REPLACE FUNCTION public.handle_new_user()
    RETURNS trigger 
    LANGUAGE plpgsql
    AS $$
    DECLARE
        default_role_id uuid;
    BEGIN
        RAISE NOTICE 'New user created with ID: %', NEW.id;

        -- Get the default USER role ID
        SELECT id INTO default_role_id FROM public.roles WHERE name = 'USER';
        RAISE NOTICE 'Default role ID: %', default_role_id;

        -- Check if default_role_id is NULL
        IF default_role_id IS NULL THEN
            RAISE EXCEPTION 'Default role ID not found for USER role';
        END IF;

        -- Insert into public.users
        INSERT INTO public.users (
            id, 
            email,
            created_at,
            updated_at
        ) VALUES (
            NEW.id,
            NEW.email,
            NEW.created_at,
            NEW.updated_at
        );

        -- Insert default role
        INSERT INTO public.user_roles (user_id, role_id)
        VALUES (NEW.id, default_role_id);

        RETURN NEW;
    END;
    $$;
`)

	// Step 2: Drop the trigger if it exists
	db.Exec(`
    DROP TRIGGER IF EXISTS on_user_created ON auth.users;
`)

	// Step 3: Create the trigger for new user signups
	db.Exec(`
    CREATE TRIGGER on_user_created
        AFTER INSERT ON auth.users
        FOR EACH ROW EXECUTE FUNCTION public.handle_new_user();
`)

	db.Exec(`
GRANT SELECT ON public.roles TO supabase_auth_admin;
`)

	db.Exec(`
GRANT INSERT, UPDATE, DELETE ON public.roles TO supabase_auth_admin;
`)
	return nil
}

func seedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{Name: "USER"},
		{Name: "ADMIN"},
	}

	for _, role := range roles {
		if err := db.Where("name = ?", role.Name).FirstOrCreate(&role).Error; err != nil {
			return err
		}
	}

	return nil
}
