package names

import (
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Mormehtar/go_experiments/db_tests/repository/interfaces"
	"github.com/Mormehtar/go_experiments/db_tests/utils"
)

func TestNames(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Names test")
}

const migrationsDir = "../../migrations"

var dbUrl = "database=postgres"
var err error
var templateDB *utils.Template

var _ = BeforeSuite(func() {
	templateDB = utils.NewTemplate(dbUrl)
	err := templateDB.Init(migrationsDir)
	if err != nil {
		panic(err)
	}
})

var _ = AfterSuite(func() {
	err := templateDB.Close()
	if err != nil {
		panic(err)
	}
})

var _ = When("Test names repository", func() {
	var names *Names
	var name *interfaces.Name
	var testDB *utils.TestDB

	BeforeEach(OncePerOrdered, func() {
		testDB, err = templateDB.GetTestDB()
		if err != nil {
			panic(err)
		}
		names = New(testDB.GetDb())
	})

	AfterEach(OncePerOrdered, func() {
		err = templateDB.DropTestDB(testDB)
		if err != nil {
			panic(err)
		}
	})

	When("Create name", func() {
		When("create name", Ordered, func() {
			BeforeAll(func() { name, err = names.Create("test") })
			It("should not raise error", func() { Expect(err).To(BeNil()) })
			It("should return result", func() { Expect(name).ToNot(BeNil()) })
			It("should return name with id", func() { Expect(name.Id).ToNot(BeNil()) })
			It("should return name with name", func() { Expect(name.Name).To(Equal("test")) })
			It("should return name with created_at", func() { Expect(name.CreatedAt).ToNot(BeNil()) })
			It("should return name with updated_at", func() { Expect(name.UpdatedAt).ToNot(BeNil()) })
			It("should return name with created_at equal to updated_at", func() {
				Expect(name.CreatedAt).To(Equal(name.UpdatedAt))
			})
		})

		When("create with used name", Ordered, func() {
			BeforeAll(func() {
				_, err = names.Create("test")
				name, err = names.Create("test")
			})
			It("should raise error", func() { Expect(err).To(MatchError(interfaces.NameIsUsedAlready)) })
			It("should not return result", func() { Expect(name).To(BeNil()) })
		})
	})

	When("Update name", func() {
		BeforeEach(OncePerOrdered, func() {
			name, _ = names.Create("test")
		})

		When("Name is not used", Ordered, func() {
			BeforeAll(func() {
				name, err = names.Update(name.Id, "test2")
			})
			It("should not raise error", func() { Expect(err).To(BeNil()) })
			It("should return result", func() { Expect(name).ToNot(BeNil()) })
			It("should update name", func() { Expect(name.Name).To(Equal("test2")) })
			It("should change DateUpdated", func() {
				Expect(name.UpdatedAt).To(BeTemporally(">", name.CreatedAt))
			})
		})

		When("Name is used", Ordered, func() {
			BeforeAll(func() {
				_, _ = names.Create("test2")
				name, err = names.Update(name.Id, "test2")
			})
			It("should raise error", func() { Expect(err).To(MatchError(interfaces.NameIsUsedAlready)) })
			It("should not return result", func() { Expect(name).To(BeNil()) })
		})
	})

	When("Get name", func() {
		When("Name exists", Ordered, func() {
			var existing *interfaces.Name
			BeforeAll(func() {
				existing, _ = names.Create("test")
				name, err = names.Get("test")
			})

			It("should not raise error", func() { Expect(err).To(BeNil()) })
			It("should return result", func() { Expect(name).ToNot(BeNil()) })
			It("should return existing name", func() { Expect(name).To(Equal(existing)) })
		})

		When("Name does not exist", Ordered, func() {
			BeforeAll(func() { name, err = names.Get("test") })
			It("should raise error", func() { Expect(err).To(MatchError(interfaces.NameIsNotFound)) })
			It("should not return result", func() { Expect(name).To(BeNil()) })
		})
	})

	When("Delete name", func() {
		When("Name exists", Ordered, func() {
			BeforeAll(func() {
				name, _ = names.Create("test")
				err = names.Delete(name.Id)
			})
			It("should not raise error", func() { Expect(err).To(BeNil()) })
			It("should delete name", func() {
				_, err = names.Get("test")
				Expect(err).To(MatchError(interfaces.NameIsNotFound))
			})
		})

		When("Name not exists", Ordered, func() {
			BeforeAll(func() {
				name, _ = names.Create("test")
				err = names.Delete(name.Id + 100500)
			})
			It("should raise error", func() { Expect(err).To(MatchError(interfaces.NameIsNotFound)) })
		})
	})

	// ginkgo -procs=10 ./repository/names
	// Should raise due to conflict on unique index in name
	// It's quite slow.
	When("try overload parallel tests", func() {
		for i := 0; i < 1000; i++ {
			When("Create name correctly", Ordered, func() {
				BeforeAll(func() { name, err = names.Create("test") })

				It("Should not return error", func() { Expect(err).To(BeNil()) })
				It("Should return property", func() { Expect(name).NotTo(BeNil()) })
			})
		}
	})
})
