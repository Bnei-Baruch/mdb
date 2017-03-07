package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Operations", testOperations)
	t.Run("OperationTypes", testOperationTypes)
	t.Run("Users", testUsers)
	t.Run("StringTranslations", testStringTranslations)
	t.Run("ContentUnitsPersons", testContentUnitsPersons)
	t.Run("ContentTypes", testContentTypes)
	t.Run("Persons", testPersons)
	t.Run("Tags", testTags)
	t.Run("ContentUnits", testContentUnits)
	t.Run("CollectionsContentUnits", testCollectionsContentUnits)
	t.Run("Collections", testCollections)
	t.Run("ContentRoles", testContentRoles)
	t.Run("Files", testFiles)
}

func TestDelete(t *testing.T) {
	t.Run("Operations", testOperationsDelete)
	t.Run("OperationTypes", testOperationTypesDelete)
	t.Run("Users", testUsersDelete)
	t.Run("StringTranslations", testStringTranslationsDelete)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsDelete)
	t.Run("ContentTypes", testContentTypesDelete)
	t.Run("Persons", testPersonsDelete)
	t.Run("Tags", testTagsDelete)
	t.Run("ContentUnits", testContentUnitsDelete)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsDelete)
	t.Run("Collections", testCollectionsDelete)
	t.Run("ContentRoles", testContentRolesDelete)
	t.Run("Files", testFilesDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Operations", testOperationsQueryDeleteAll)
	t.Run("OperationTypes", testOperationTypesQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
	t.Run("StringTranslations", testStringTranslationsQueryDeleteAll)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsQueryDeleteAll)
	t.Run("ContentTypes", testContentTypesQueryDeleteAll)
	t.Run("Persons", testPersonsQueryDeleteAll)
	t.Run("Tags", testTagsQueryDeleteAll)
	t.Run("ContentUnits", testContentUnitsQueryDeleteAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsQueryDeleteAll)
	t.Run("Collections", testCollectionsQueryDeleteAll)
	t.Run("ContentRoles", testContentRolesQueryDeleteAll)
	t.Run("Files", testFilesQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Operations", testOperationsSliceDeleteAll)
	t.Run("OperationTypes", testOperationTypesSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
	t.Run("StringTranslations", testStringTranslationsSliceDeleteAll)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsSliceDeleteAll)
	t.Run("ContentTypes", testContentTypesSliceDeleteAll)
	t.Run("Persons", testPersonsSliceDeleteAll)
	t.Run("Tags", testTagsSliceDeleteAll)
	t.Run("ContentUnits", testContentUnitsSliceDeleteAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsSliceDeleteAll)
	t.Run("Collections", testCollectionsSliceDeleteAll)
	t.Run("ContentRoles", testContentRolesSliceDeleteAll)
	t.Run("Files", testFilesSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Operations", testOperationsExists)
	t.Run("OperationTypes", testOperationTypesExists)
	t.Run("Users", testUsersExists)
	t.Run("StringTranslations", testStringTranslationsExists)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsExists)
	t.Run("ContentTypes", testContentTypesExists)
	t.Run("Persons", testPersonsExists)
	t.Run("Tags", testTagsExists)
	t.Run("ContentUnits", testContentUnitsExists)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsExists)
	t.Run("Collections", testCollectionsExists)
	t.Run("ContentRoles", testContentRolesExists)
	t.Run("Files", testFilesExists)
}

func TestFind(t *testing.T) {
	t.Run("Operations", testOperationsFind)
	t.Run("OperationTypes", testOperationTypesFind)
	t.Run("Users", testUsersFind)
	t.Run("StringTranslations", testStringTranslationsFind)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsFind)
	t.Run("ContentTypes", testContentTypesFind)
	t.Run("Persons", testPersonsFind)
	t.Run("Tags", testTagsFind)
	t.Run("ContentUnits", testContentUnitsFind)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsFind)
	t.Run("Collections", testCollectionsFind)
	t.Run("ContentRoles", testContentRolesFind)
	t.Run("Files", testFilesFind)
}

func TestBind(t *testing.T) {
	t.Run("Operations", testOperationsBind)
	t.Run("OperationTypes", testOperationTypesBind)
	t.Run("Users", testUsersBind)
	t.Run("StringTranslations", testStringTranslationsBind)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsBind)
	t.Run("ContentTypes", testContentTypesBind)
	t.Run("Persons", testPersonsBind)
	t.Run("Tags", testTagsBind)
	t.Run("ContentUnits", testContentUnitsBind)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsBind)
	t.Run("Collections", testCollectionsBind)
	t.Run("ContentRoles", testContentRolesBind)
	t.Run("Files", testFilesBind)
}

func TestOne(t *testing.T) {
	t.Run("Operations", testOperationsOne)
	t.Run("OperationTypes", testOperationTypesOne)
	t.Run("Users", testUsersOne)
	t.Run("StringTranslations", testStringTranslationsOne)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsOne)
	t.Run("ContentTypes", testContentTypesOne)
	t.Run("Persons", testPersonsOne)
	t.Run("Tags", testTagsOne)
	t.Run("ContentUnits", testContentUnitsOne)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsOne)
	t.Run("Collections", testCollectionsOne)
	t.Run("ContentRoles", testContentRolesOne)
	t.Run("Files", testFilesOne)
}

func TestAll(t *testing.T) {
	t.Run("Operations", testOperationsAll)
	t.Run("OperationTypes", testOperationTypesAll)
	t.Run("Users", testUsersAll)
	t.Run("StringTranslations", testStringTranslationsAll)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsAll)
	t.Run("ContentTypes", testContentTypesAll)
	t.Run("Persons", testPersonsAll)
	t.Run("Tags", testTagsAll)
	t.Run("ContentUnits", testContentUnitsAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsAll)
	t.Run("Collections", testCollectionsAll)
	t.Run("ContentRoles", testContentRolesAll)
	t.Run("Files", testFilesAll)
}

func TestCount(t *testing.T) {
	t.Run("Operations", testOperationsCount)
	t.Run("OperationTypes", testOperationTypesCount)
	t.Run("Users", testUsersCount)
	t.Run("StringTranslations", testStringTranslationsCount)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsCount)
	t.Run("ContentTypes", testContentTypesCount)
	t.Run("Persons", testPersonsCount)
	t.Run("Tags", testTagsCount)
	t.Run("ContentUnits", testContentUnitsCount)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsCount)
	t.Run("Collections", testCollectionsCount)
	t.Run("ContentRoles", testContentRolesCount)
	t.Run("Files", testFilesCount)
}

func TestInsert(t *testing.T) {
	t.Run("Operations", testOperationsInsert)
	t.Run("Operations", testOperationsInsertWhitelist)
	t.Run("OperationTypes", testOperationTypesInsert)
	t.Run("OperationTypes", testOperationTypesInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
	t.Run("StringTranslations", testStringTranslationsInsert)
	t.Run("StringTranslations", testStringTranslationsInsertWhitelist)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsInsert)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsInsertWhitelist)
	t.Run("ContentTypes", testContentTypesInsert)
	t.Run("ContentTypes", testContentTypesInsertWhitelist)
	t.Run("Persons", testPersonsInsert)
	t.Run("Persons", testPersonsInsertWhitelist)
	t.Run("Tags", testTagsInsert)
	t.Run("Tags", testTagsInsertWhitelist)
	t.Run("ContentUnits", testContentUnitsInsert)
	t.Run("ContentUnits", testContentUnitsInsertWhitelist)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsInsert)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsInsertWhitelist)
	t.Run("Collections", testCollectionsInsert)
	t.Run("Collections", testCollectionsInsertWhitelist)
	t.Run("ContentRoles", testContentRolesInsert)
	t.Run("ContentRoles", testContentRolesInsertWhitelist)
	t.Run("Files", testFilesInsert)
	t.Run("Files", testFilesInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("OperationToOperationTypeUsingType", testOperationToOneOperationTypeUsingType)
	t.Run("OperationToUserUsingUser", testOperationToOneUserUsingUser)
	t.Run("ContentUnitsPersonToContentUnitUsingContentUnit", testContentUnitsPersonToOneContentUnitUsingContentUnit)
	t.Run("ContentUnitsPersonToPersonUsingPerson", testContentUnitsPersonToOnePersonUsingPerson)
	t.Run("ContentUnitsPersonToContentRoleUsingRole", testContentUnitsPersonToOneContentRoleUsingRole)
	t.Run("PersonToStringTranslationUsingDescription", testPersonToOneStringTranslationUsingDescription)
	t.Run("PersonToStringTranslationUsingName", testPersonToOneStringTranslationUsingName)
	t.Run("TagToStringTranslationUsingLabel", testTagToOneStringTranslationUsingLabel)
	t.Run("TagToTagUsingParent", testTagToOneTagUsingParent)
	t.Run("ContentUnitToStringTranslationUsingDescription", testContentUnitToOneStringTranslationUsingDescription)
	t.Run("ContentUnitToStringTranslationUsingName", testContentUnitToOneStringTranslationUsingName)
	t.Run("ContentUnitToContentTypeUsingType", testContentUnitToOneContentTypeUsingType)
	t.Run("CollectionsContentUnitToCollectionUsingCollection", testCollectionsContentUnitToOneCollectionUsingCollection)
	t.Run("CollectionsContentUnitToContentUnitUsingContentUnit", testCollectionsContentUnitToOneContentUnitUsingContentUnit)
	t.Run("CollectionToStringTranslationUsingDescription", testCollectionToOneStringTranslationUsingDescription)
	t.Run("CollectionToStringTranslationUsingName", testCollectionToOneStringTranslationUsingName)
	t.Run("CollectionToContentTypeUsingType", testCollectionToOneContentTypeUsingType)
	t.Run("ContentRoleToStringTranslationUsingDescription", testContentRoleToOneStringTranslationUsingDescription)
	t.Run("ContentRoleToStringTranslationUsingName", testContentRoleToOneStringTranslationUsingName)
	t.Run("FileToContentUnitUsingContentUnit", testFileToOneContentUnitUsingContentUnit)
	t.Run("FileToFileUsingParent", testFileToOneFileUsingParent)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("OperationToFiles", testOperationToManyFiles)
	t.Run("OperationTypeToTypeOperations", testOperationTypeToManyTypeOperations)
	t.Run("UserToOperations", testUserToManyOperations)
	t.Run("StringTranslationToDescriptionPersons", testStringTranslationToManyDescriptionPersons)
	t.Run("StringTranslationToNamePersons", testStringTranslationToManyNamePersons)
	t.Run("StringTranslationToLabelTags", testStringTranslationToManyLabelTags)
	t.Run("StringTranslationToDescriptionContentUnits", testStringTranslationToManyDescriptionContentUnits)
	t.Run("StringTranslationToNameContentUnits", testStringTranslationToManyNameContentUnits)
	t.Run("StringTranslationToDescriptionCollections", testStringTranslationToManyDescriptionCollections)
	t.Run("StringTranslationToNameCollections", testStringTranslationToManyNameCollections)
	t.Run("StringTranslationToDescriptionContentRoles", testStringTranslationToManyDescriptionContentRoles)
	t.Run("StringTranslationToNameContentRoles", testStringTranslationToManyNameContentRoles)
	t.Run("ContentTypeToTypeContentUnits", testContentTypeToManyTypeContentUnits)
	t.Run("ContentTypeToTypeCollections", testContentTypeToManyTypeCollections)
	t.Run("PersonToContentUnitsPersons", testPersonToManyContentUnitsPersons)
	t.Run("TagToParentTags", testTagToManyParentTags)
	t.Run("ContentUnitToContentUnitsPersons", testContentUnitToManyContentUnitsPersons)
	t.Run("ContentUnitToCollectionsContentUnits", testContentUnitToManyCollectionsContentUnits)
	t.Run("ContentUnitToFiles", testContentUnitToManyFiles)
	t.Run("CollectionToCollectionsContentUnits", testCollectionToManyCollectionsContentUnits)
	t.Run("ContentRoleToRoleContentUnitsPersons", testContentRoleToManyRoleContentUnitsPersons)
	t.Run("FileToOperations", testFileToManyOperations)
	t.Run("FileToParentFiles", testFileToManyParentFiles)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("OperationToOperationTypeUsingType", testOperationToOneSetOpOperationTypeUsingType)
	t.Run("OperationToUserUsingUser", testOperationToOneSetOpUserUsingUser)
	t.Run("ContentUnitsPersonToContentUnitUsingContentUnit", testContentUnitsPersonToOneSetOpContentUnitUsingContentUnit)
	t.Run("ContentUnitsPersonToPersonUsingPerson", testContentUnitsPersonToOneSetOpPersonUsingPerson)
	t.Run("ContentUnitsPersonToContentRoleUsingRole", testContentUnitsPersonToOneSetOpContentRoleUsingRole)
	t.Run("PersonToStringTranslationUsingDescription", testPersonToOneSetOpStringTranslationUsingDescription)
	t.Run("PersonToStringTranslationUsingName", testPersonToOneSetOpStringTranslationUsingName)
	t.Run("TagToStringTranslationUsingLabel", testTagToOneSetOpStringTranslationUsingLabel)
	t.Run("TagToTagUsingParent", testTagToOneSetOpTagUsingParent)
	t.Run("ContentUnitToStringTranslationUsingDescription", testContentUnitToOneSetOpStringTranslationUsingDescription)
	t.Run("ContentUnitToStringTranslationUsingName", testContentUnitToOneSetOpStringTranslationUsingName)
	t.Run("ContentUnitToContentTypeUsingType", testContentUnitToOneSetOpContentTypeUsingType)
	t.Run("CollectionsContentUnitToCollectionUsingCollection", testCollectionsContentUnitToOneSetOpCollectionUsingCollection)
	t.Run("CollectionsContentUnitToContentUnitUsingContentUnit", testCollectionsContentUnitToOneSetOpContentUnitUsingContentUnit)
	t.Run("CollectionToStringTranslationUsingDescription", testCollectionToOneSetOpStringTranslationUsingDescription)
	t.Run("CollectionToStringTranslationUsingName", testCollectionToOneSetOpStringTranslationUsingName)
	t.Run("CollectionToContentTypeUsingType", testCollectionToOneSetOpContentTypeUsingType)
	t.Run("ContentRoleToStringTranslationUsingDescription", testContentRoleToOneSetOpStringTranslationUsingDescription)
	t.Run("ContentRoleToStringTranslationUsingName", testContentRoleToOneSetOpStringTranslationUsingName)
	t.Run("FileToContentUnitUsingContentUnit", testFileToOneSetOpContentUnitUsingContentUnit)
	t.Run("FileToFileUsingParent", testFileToOneSetOpFileUsingParent)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("OperationToUserUsingUser", testOperationToOneRemoveOpUserUsingUser)
	t.Run("PersonToStringTranslationUsingDescription", testPersonToOneRemoveOpStringTranslationUsingDescription)
	t.Run("TagToTagUsingParent", testTagToOneRemoveOpTagUsingParent)
	t.Run("ContentUnitToStringTranslationUsingDescription", testContentUnitToOneRemoveOpStringTranslationUsingDescription)
	t.Run("CollectionToStringTranslationUsingDescription", testCollectionToOneRemoveOpStringTranslationUsingDescription)
	t.Run("ContentRoleToStringTranslationUsingDescription", testContentRoleToOneRemoveOpStringTranslationUsingDescription)
	t.Run("FileToContentUnitUsingContentUnit", testFileToOneRemoveOpContentUnitUsingContentUnit)
	t.Run("FileToFileUsingParent", testFileToOneRemoveOpFileUsingParent)
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("OperationToFiles", testOperationToManyAddOpFiles)
	t.Run("OperationTypeToTypeOperations", testOperationTypeToManyAddOpTypeOperations)
	t.Run("UserToOperations", testUserToManyAddOpOperations)
	t.Run("StringTranslationToDescriptionPersons", testStringTranslationToManyAddOpDescriptionPersons)
	t.Run("StringTranslationToNamePersons", testStringTranslationToManyAddOpNamePersons)
	t.Run("StringTranslationToLabelTags", testStringTranslationToManyAddOpLabelTags)
	t.Run("StringTranslationToDescriptionContentUnits", testStringTranslationToManyAddOpDescriptionContentUnits)
	t.Run("StringTranslationToNameContentUnits", testStringTranslationToManyAddOpNameContentUnits)
	t.Run("StringTranslationToDescriptionCollections", testStringTranslationToManyAddOpDescriptionCollections)
	t.Run("StringTranslationToNameCollections", testStringTranslationToManyAddOpNameCollections)
	t.Run("StringTranslationToDescriptionContentRoles", testStringTranslationToManyAddOpDescriptionContentRoles)
	t.Run("StringTranslationToNameContentRoles", testStringTranslationToManyAddOpNameContentRoles)
	t.Run("ContentTypeToTypeContentUnits", testContentTypeToManyAddOpTypeContentUnits)
	t.Run("ContentTypeToTypeCollections", testContentTypeToManyAddOpTypeCollections)
	t.Run("PersonToContentUnitsPersons", testPersonToManyAddOpContentUnitsPersons)
	t.Run("TagToParentTags", testTagToManyAddOpParentTags)
	t.Run("ContentUnitToContentUnitsPersons", testContentUnitToManyAddOpContentUnitsPersons)
	t.Run("ContentUnitToCollectionsContentUnits", testContentUnitToManyAddOpCollectionsContentUnits)
	t.Run("ContentUnitToFiles", testContentUnitToManyAddOpFiles)
	t.Run("CollectionToCollectionsContentUnits", testCollectionToManyAddOpCollectionsContentUnits)
	t.Run("ContentRoleToRoleContentUnitsPersons", testContentRoleToManyAddOpRoleContentUnitsPersons)
	t.Run("FileToOperations", testFileToManyAddOpOperations)
	t.Run("FileToParentFiles", testFileToManyAddOpParentFiles)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("UserToOperations", testUserToManySetOpOperations)
	t.Run("StringTranslationToDescriptionPersons", testStringTranslationToManySetOpDescriptionPersons)
	t.Run("StringTranslationToDescriptionContentUnits", testStringTranslationToManySetOpDescriptionContentUnits)
	t.Run("StringTranslationToDescriptionCollections", testStringTranslationToManySetOpDescriptionCollections)
	t.Run("StringTranslationToDescriptionContentRoles", testStringTranslationToManySetOpDescriptionContentRoles)
	t.Run("TagToParentTags", testTagToManySetOpParentTags)
	t.Run("ContentUnitToFiles", testContentUnitToManySetOpFiles)
	t.Run("FileToParentFiles", testFileToManySetOpParentFiles)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("UserToOperations", testUserToManyRemoveOpOperations)
	t.Run("StringTranslationToDescriptionPersons", testStringTranslationToManyRemoveOpDescriptionPersons)
	t.Run("StringTranslationToDescriptionContentUnits", testStringTranslationToManyRemoveOpDescriptionContentUnits)
	t.Run("StringTranslationToDescriptionCollections", testStringTranslationToManyRemoveOpDescriptionCollections)
	t.Run("StringTranslationToDescriptionContentRoles", testStringTranslationToManyRemoveOpDescriptionContentRoles)
	t.Run("TagToParentTags", testTagToManyRemoveOpParentTags)
	t.Run("ContentUnitToFiles", testContentUnitToManyRemoveOpFiles)
	t.Run("FileToParentFiles", testFileToManyRemoveOpParentFiles)
}

func TestReload(t *testing.T) {
	t.Run("Operations", testOperationsReload)
	t.Run("OperationTypes", testOperationTypesReload)
	t.Run("Users", testUsersReload)
	t.Run("StringTranslations", testStringTranslationsReload)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsReload)
	t.Run("ContentTypes", testContentTypesReload)
	t.Run("Persons", testPersonsReload)
	t.Run("Tags", testTagsReload)
	t.Run("ContentUnits", testContentUnitsReload)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsReload)
	t.Run("Collections", testCollectionsReload)
	t.Run("ContentRoles", testContentRolesReload)
	t.Run("Files", testFilesReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Operations", testOperationsReloadAll)
	t.Run("OperationTypes", testOperationTypesReloadAll)
	t.Run("Users", testUsersReloadAll)
	t.Run("StringTranslations", testStringTranslationsReloadAll)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsReloadAll)
	t.Run("ContentTypes", testContentTypesReloadAll)
	t.Run("Persons", testPersonsReloadAll)
	t.Run("Tags", testTagsReloadAll)
	t.Run("ContentUnits", testContentUnitsReloadAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsReloadAll)
	t.Run("Collections", testCollectionsReloadAll)
	t.Run("ContentRoles", testContentRolesReloadAll)
	t.Run("Files", testFilesReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Operations", testOperationsSelect)
	t.Run("OperationTypes", testOperationTypesSelect)
	t.Run("Users", testUsersSelect)
	t.Run("StringTranslations", testStringTranslationsSelect)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsSelect)
	t.Run("ContentTypes", testContentTypesSelect)
	t.Run("Persons", testPersonsSelect)
	t.Run("Tags", testTagsSelect)
	t.Run("ContentUnits", testContentUnitsSelect)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsSelect)
	t.Run("Collections", testCollectionsSelect)
	t.Run("ContentRoles", testContentRolesSelect)
	t.Run("Files", testFilesSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Operations", testOperationsUpdate)
	t.Run("OperationTypes", testOperationTypesUpdate)
	t.Run("Users", testUsersUpdate)
	t.Run("StringTranslations", testStringTranslationsUpdate)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsUpdate)
	t.Run("ContentTypes", testContentTypesUpdate)
	t.Run("Persons", testPersonsUpdate)
	t.Run("Tags", testTagsUpdate)
	t.Run("ContentUnits", testContentUnitsUpdate)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsUpdate)
	t.Run("Collections", testCollectionsUpdate)
	t.Run("ContentRoles", testContentRolesUpdate)
	t.Run("Files", testFilesUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Operations", testOperationsSliceUpdateAll)
	t.Run("OperationTypes", testOperationTypesSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
	t.Run("StringTranslations", testStringTranslationsSliceUpdateAll)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsSliceUpdateAll)
	t.Run("ContentTypes", testContentTypesSliceUpdateAll)
	t.Run("Persons", testPersonsSliceUpdateAll)
	t.Run("Tags", testTagsSliceUpdateAll)
	t.Run("ContentUnits", testContentUnitsSliceUpdateAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsSliceUpdateAll)
	t.Run("Collections", testCollectionsSliceUpdateAll)
	t.Run("ContentRoles", testContentRolesSliceUpdateAll)
	t.Run("Files", testFilesSliceUpdateAll)
}

func TestUpsert(t *testing.T) {
	t.Run("Operations", testOperationsUpsert)
	t.Run("OperationTypes", testOperationTypesUpsert)
	t.Run("Users", testUsersUpsert)
	t.Run("StringTranslations", testStringTranslationsUpsert)
	t.Run("ContentUnitsPersons", testContentUnitsPersonsUpsert)
	t.Run("ContentTypes", testContentTypesUpsert)
	t.Run("Persons", testPersonsUpsert)
	t.Run("Tags", testTagsUpsert)
	t.Run("ContentUnits", testContentUnitsUpsert)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsUpsert)
	t.Run("Collections", testCollectionsUpsert)
	t.Run("ContentRoles", testContentRolesUpsert)
	t.Run("Files", testFilesUpsert)
}
