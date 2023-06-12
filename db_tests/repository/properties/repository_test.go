package properties

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Mormehtar/go_experiments/db_tests/repository/interfaces"
	names2 "github.com/Mormehtar/go_experiments/db_tests/repository/names"
	"github.com/Mormehtar/go_experiments/db_tests/utils"
)

func TestProperties(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Properties test")
}

const migrationsDir = "../../migrations"

var dbUrl = "database=postgres"
var err error
var transactional *utils.Transactional

var _ = BeforeSuite(func() {
	transactional = utils.NewTransactional(dbUrl)
	if err := transactional.Init(migrationsDir); err != nil {
		panic(err)
	}
})

var _ = AfterSuite(func() {
	if err := transactional.Close(); err != nil {
		panic(err)
	}
})

var _ = When("Test properties", func() {
	var properties *Properties
	var names *names2.Names
	var testDB *sql.DB

	BeforeEach(OncePerOrdered, func() {
		testDB, err = transactional.GetTestDB()
		if err != nil {
			panic(err)
		}
		properties = New(testDB)
		names = names2.New(testDB)
	})

	AfterEach(OncePerOrdered, func() {
		if err := transactional.ReturnTestDB(); err != nil {
			panic(err)
		}
	})

	When("Create property", func() {
		var property *interfaces.Property

		When("Create property correctly", Ordered, func() {
			var name *interfaces.Name
			BeforeAll(func() {
				name, err = names.Create("test")
				if err != nil {
					panic(err)
				}
				property, err = properties.Create(name.Id, "key", "value")
			})

			It("Should not return error", func() { Expect(err).To(BeNil()) })
			It("Should return property", func() { Expect(property).NotTo(BeNil()) })
			It("Should return property with correct key", func() { Expect(property.Key).To(Equal("key")) })
			It("Should return property with correct value", func() { Expect(property.Value).To(Equal("value")) })
			It("Should return property with correct name", func() { Expect(property.NameId).To(Equal(name.Id)) })
			It("Should return property with created_at", func() { Expect(property.CreatedAt).NotTo(BeNil()) })
			It("Should return property with updated_at", func() { Expect(property.UpdatedAt).NotTo(BeNil()) })
			It("Should create property with equal created_at and updated_at", func() {
				Expect(property.CreatedAt).To(Equal(property.UpdatedAt))
			})
		})

		When("Name does not exist", Ordered, func() {
			BeforeAll(func() {
				property, err = properties.Create(100500, "key", "value")
			})
			It("should return error", func() { Expect(err).To(Equal(interfaces.NameIsNotFound)) })
			It("should not return property", func() { Expect(property).To(BeNil()) })
		})
	})

	When("Get property", func() {
		var name *interfaces.Name
		var property *interfaces.Property

		BeforeEach(OncePerOrdered, func() {
			name, err = names.Create("test")
			if err != nil {
				panic(err)
			}
			property, err = properties.Create(name.Id, "key", "value")
		})

		When("Get property correctly", Ordered, func() {
			var getProperty []*interfaces.Property
			BeforeAll(func() {
				getProperty, err = properties.Get(name.Id, property.Key)
			})

			It("Should not return error", func() { Expect(err).To(BeNil()) })
			It("Should return property", func() { Expect(getProperty).NotTo(BeNil()) })
			It("should return one property", func() { Expect(getProperty).To(HaveLen(1)) })
			It("Should return correct property", func() { Expect(getProperty[0]).To(Equal(property)) })
		})
	})

	When("Update property", func() {
		var name *interfaces.Name
		var property *interfaces.Property

		BeforeEach(OncePerOrdered, func() {
			name, err = names.Create("test")
			if err != nil {
				panic(err)
			}
			property, err = properties.Create(name.Id, "key", "value")
		})

		When("Update property correctly", Ordered, func() {
			var updated *interfaces.Property
			BeforeAll(func() {
				updated, err = properties.Update(property.Id, "new_key", "new_value")
			})

			It("Should not return error", func() { Expect(err).To(BeNil()) })
			It("Should update property key", func() { Expect(updated.Key).To(Equal("new_key")) })
			It("Should update property value", func() { Expect(updated.Value).To(Equal("new_value")) })
			// Fails if now() is used in update query
			It("Should update property updated_at", func() {
				Expect(updated.UpdatedAt).To(BeTemporally(">", updated.CreatedAt))
			})
		})

		When("Property does not exist", Ordered, func() {
			var updated *interfaces.Property
			BeforeAll(func() {
				updated, err = properties.Update(property.Id+100500, "new_key", "new_value")
			})

			It("Should return error", func() { Expect(err).To(Equal(interfaces.PropertyIsNotFound)) })
			It("Should not return property", func() { Expect(updated).To(BeNil()) })
		})
	})

	When("Delete property", func() {
		When("Name exists", Ordered, func() {
			var name *interfaces.Name

			BeforeAll(func() {
				name, _ = names.Create("test")
				property, _ := properties.Create(name.Id, "key", "value")
				err = properties.Delete(property.Id)
			})
			It("should not raise error", func() { Expect(err).To(BeNil()) })
			It("should delete name", func() {
				result, err := properties.Get(name.Id, "key")
				Expect(err).To(BeNil())
				Expect(result).To(HaveLen(0))
			})
		})

		When("Name not exists", Ordered, func() {
			BeforeAll(func() {
				err = properties.Delete(100500)
			})
			It("should raise error", func() { Expect(err).To(MatchError(interfaces.PropertyIsNotFound)) })
		})
	})

	// ginkgo -procs=10 ./repository/properties
	// Should raise due to conflict on unique index in name
	// It's much faster than with copies of DB.
	When("try overload parallel tests", func() {
		for i := 0; i < 1000; i++ {
			When("Create property correctly", Ordered, func() {
				var property *interfaces.Property
				var name *interfaces.Name
				BeforeAll(func() {
					name, err = names.Create("test")
					if err != nil {
						panic(err)
					}
					property, err = properties.Create(name.Id, "key", "value")
				})

				It("Should not return error", func() { Expect(err).To(BeNil()) })
				It("Should return property", func() { Expect(property).NotTo(BeNil()) })
			})
		}
	})
})
