package repository_migrate

import (
	"database/sql"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Mormehtar/go_experiments/db_tests/utils"
)

var migrations goose.Migrations

func TestMigrations(t *testing.T) {
	migrations, err = goose.CollectMigrations(migrationsDir, 0, goose.MaxVersion)
	if err != nil {
		panic(err)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Stair test")
}

const migrationsDir = "../migrations"
const dbUrl = "database=postgres"

var dbManager *utils.CleanDB
var testDB *sql.DB

var err error

var _ = BeforeSuite(func() {
	dbManager = utils.NewCleanDB(dbUrl)
	if err := dbManager.Init(); err != nil {
		panic(err)
	}
	testDB = dbManager.GetTestDB()
	utils.ConfigureGoose()
})

var _ = AfterSuite(func() {
	if err := dbManager.Close(); err != nil {
		panic(err)
	}
})

var _ = Describe("Stair", func() {
	When("all the way", Ordered, func() {
		for i, migration := range migrations {
			inner := i

			When(fmt.Sprintf("Migrate %s", migration.String()), func() {
				It("should migrate forward", func() {
					Expect(goose.UpByOne(testDB, migrationsDir)).To(BeNil())
				})
				It("should migrate backward", func() {
					if inner == 0 {
						Expect(goose.DownTo(testDB, migrationsDir, 0))
					} else {
						Expect(goose.DownTo(testDB, migrationsDir, migrations[inner-1].Version)).To(BeNil())
					}
				})
				It("should migrate forward again", func() {
					Expect(goose.UpByOne(testDB, migrationsDir)).To(BeNil())
				})
			})
		}
	})
})
