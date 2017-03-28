package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("ContentTypes", testContentTypes)
	t.Run("ContentUnits", testContentUnits)
	t.Run("Files", testFiles)
	t.Run("OperationTypes", testOperationTypes)
	t.Run("Users", testUsers)
	t.Run("Sources", testSources)
	t.Run("CollectionI18ns", testCollectionI18ns)
	t.Run("CollectionsContentUnits", testCollectionsContentUnits)
	t.Run("Collections", testCollections)
	t.Run("ContentUnitI18ns", testContentUnitI18ns)
	t.Run("Operations", testOperations)
	t.Run("Tags", testTags)
	t.Run("TagI18ns", testTagI18ns)
	t.Run("SourceTypes", testSourceTypes)
	t.Run("SourceI18ns", testSourceI18ns)
	t.Run("AuthorI18ns", testAuthorI18ns)
	t.Run("Authors", testAuthors)
}

func TestDelete(t *testing.T) {
	t.Run("ContentTypes", testContentTypesDelete)
	t.Run("ContentUnits", testContentUnitsDelete)
	t.Run("Files", testFilesDelete)
	t.Run("OperationTypes", testOperationTypesDelete)
	t.Run("Users", testUsersDelete)
	t.Run("Sources", testSourcesDelete)
	t.Run("CollectionI18ns", testCollectionI18nsDelete)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsDelete)
	t.Run("Collections", testCollectionsDelete)
	t.Run("ContentUnitI18ns", testContentUnitI18nsDelete)
	t.Run("Operations", testOperationsDelete)
	t.Run("Tags", testTagsDelete)
	t.Run("TagI18ns", testTagI18nsDelete)
	t.Run("SourceTypes", testSourceTypesDelete)
	t.Run("SourceI18ns", testSourceI18nsDelete)
	t.Run("AuthorI18ns", testAuthorI18nsDelete)
	t.Run("Authors", testAuthorsDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("ContentTypes", testContentTypesQueryDeleteAll)
	t.Run("ContentUnits", testContentUnitsQueryDeleteAll)
	t.Run("Files", testFilesQueryDeleteAll)
	t.Run("OperationTypes", testOperationTypesQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
	t.Run("Sources", testSourcesQueryDeleteAll)
	t.Run("CollectionI18ns", testCollectionI18nsQueryDeleteAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsQueryDeleteAll)
	t.Run("Collections", testCollectionsQueryDeleteAll)
	t.Run("ContentUnitI18ns", testContentUnitI18nsQueryDeleteAll)
	t.Run("Operations", testOperationsQueryDeleteAll)
	t.Run("Tags", testTagsQueryDeleteAll)
	t.Run("TagI18ns", testTagI18nsQueryDeleteAll)
	t.Run("SourceTypes", testSourceTypesQueryDeleteAll)
	t.Run("SourceI18ns", testSourceI18nsQueryDeleteAll)
	t.Run("AuthorI18ns", testAuthorI18nsQueryDeleteAll)
	t.Run("Authors", testAuthorsQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("ContentTypes", testContentTypesSliceDeleteAll)
	t.Run("ContentUnits", testContentUnitsSliceDeleteAll)
	t.Run("Files", testFilesSliceDeleteAll)
	t.Run("OperationTypes", testOperationTypesSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
	t.Run("Sources", testSourcesSliceDeleteAll)
	t.Run("CollectionI18ns", testCollectionI18nsSliceDeleteAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsSliceDeleteAll)
	t.Run("Collections", testCollectionsSliceDeleteAll)
	t.Run("ContentUnitI18ns", testContentUnitI18nsSliceDeleteAll)
	t.Run("Operations", testOperationsSliceDeleteAll)
	t.Run("Tags", testTagsSliceDeleteAll)
	t.Run("TagI18ns", testTagI18nsSliceDeleteAll)
	t.Run("SourceTypes", testSourceTypesSliceDeleteAll)
	t.Run("SourceI18ns", testSourceI18nsSliceDeleteAll)
	t.Run("AuthorI18ns", testAuthorI18nsSliceDeleteAll)
	t.Run("Authors", testAuthorsSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("ContentTypes", testContentTypesExists)
	t.Run("ContentUnits", testContentUnitsExists)
	t.Run("Files", testFilesExists)
	t.Run("OperationTypes", testOperationTypesExists)
	t.Run("Users", testUsersExists)
	t.Run("Sources", testSourcesExists)
	t.Run("CollectionI18ns", testCollectionI18nsExists)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsExists)
	t.Run("Collections", testCollectionsExists)
	t.Run("ContentUnitI18ns", testContentUnitI18nsExists)
	t.Run("Operations", testOperationsExists)
	t.Run("Tags", testTagsExists)
	t.Run("TagI18ns", testTagI18nsExists)
	t.Run("SourceTypes", testSourceTypesExists)
	t.Run("SourceI18ns", testSourceI18nsExists)
	t.Run("AuthorI18ns", testAuthorI18nsExists)
	t.Run("Authors", testAuthorsExists)
}

func TestFind(t *testing.T) {
	t.Run("ContentTypes", testContentTypesFind)
	t.Run("ContentUnits", testContentUnitsFind)
	t.Run("Files", testFilesFind)
	t.Run("OperationTypes", testOperationTypesFind)
	t.Run("Users", testUsersFind)
	t.Run("Sources", testSourcesFind)
	t.Run("CollectionI18ns", testCollectionI18nsFind)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsFind)
	t.Run("Collections", testCollectionsFind)
	t.Run("ContentUnitI18ns", testContentUnitI18nsFind)
	t.Run("Operations", testOperationsFind)
	t.Run("Tags", testTagsFind)
	t.Run("TagI18ns", testTagI18nsFind)
	t.Run("SourceTypes", testSourceTypesFind)
	t.Run("SourceI18ns", testSourceI18nsFind)
	t.Run("AuthorI18ns", testAuthorI18nsFind)
	t.Run("Authors", testAuthorsFind)
}

func TestBind(t *testing.T) {
	t.Run("ContentTypes", testContentTypesBind)
	t.Run("ContentUnits", testContentUnitsBind)
	t.Run("Files", testFilesBind)
	t.Run("OperationTypes", testOperationTypesBind)
	t.Run("Users", testUsersBind)
	t.Run("Sources", testSourcesBind)
	t.Run("CollectionI18ns", testCollectionI18nsBind)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsBind)
	t.Run("Collections", testCollectionsBind)
	t.Run("ContentUnitI18ns", testContentUnitI18nsBind)
	t.Run("Operations", testOperationsBind)
	t.Run("Tags", testTagsBind)
	t.Run("TagI18ns", testTagI18nsBind)
	t.Run("SourceTypes", testSourceTypesBind)
	t.Run("SourceI18ns", testSourceI18nsBind)
	t.Run("AuthorI18ns", testAuthorI18nsBind)
	t.Run("Authors", testAuthorsBind)
}

func TestOne(t *testing.T) {
	t.Run("ContentTypes", testContentTypesOne)
	t.Run("ContentUnits", testContentUnitsOne)
	t.Run("Files", testFilesOne)
	t.Run("OperationTypes", testOperationTypesOne)
	t.Run("Users", testUsersOne)
	t.Run("Sources", testSourcesOne)
	t.Run("CollectionI18ns", testCollectionI18nsOne)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsOne)
	t.Run("Collections", testCollectionsOne)
	t.Run("ContentUnitI18ns", testContentUnitI18nsOne)
	t.Run("Operations", testOperationsOne)
	t.Run("Tags", testTagsOne)
	t.Run("TagI18ns", testTagI18nsOne)
	t.Run("SourceTypes", testSourceTypesOne)
	t.Run("SourceI18ns", testSourceI18nsOne)
	t.Run("AuthorI18ns", testAuthorI18nsOne)
	t.Run("Authors", testAuthorsOne)
}

func TestAll(t *testing.T) {
	t.Run("ContentTypes", testContentTypesAll)
	t.Run("ContentUnits", testContentUnitsAll)
	t.Run("Files", testFilesAll)
	t.Run("OperationTypes", testOperationTypesAll)
	t.Run("Users", testUsersAll)
	t.Run("Sources", testSourcesAll)
	t.Run("CollectionI18ns", testCollectionI18nsAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsAll)
	t.Run("Collections", testCollectionsAll)
	t.Run("ContentUnitI18ns", testContentUnitI18nsAll)
	t.Run("Operations", testOperationsAll)
	t.Run("Tags", testTagsAll)
	t.Run("TagI18ns", testTagI18nsAll)
	t.Run("SourceTypes", testSourceTypesAll)
	t.Run("SourceI18ns", testSourceI18nsAll)
	t.Run("AuthorI18ns", testAuthorI18nsAll)
	t.Run("Authors", testAuthorsAll)
}

func TestCount(t *testing.T) {
	t.Run("ContentTypes", testContentTypesCount)
	t.Run("ContentUnits", testContentUnitsCount)
	t.Run("Files", testFilesCount)
	t.Run("OperationTypes", testOperationTypesCount)
	t.Run("Users", testUsersCount)
	t.Run("Sources", testSourcesCount)
	t.Run("CollectionI18ns", testCollectionI18nsCount)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsCount)
	t.Run("Collections", testCollectionsCount)
	t.Run("ContentUnitI18ns", testContentUnitI18nsCount)
	t.Run("Operations", testOperationsCount)
	t.Run("Tags", testTagsCount)
	t.Run("TagI18ns", testTagI18nsCount)
	t.Run("SourceTypes", testSourceTypesCount)
	t.Run("SourceI18ns", testSourceI18nsCount)
	t.Run("AuthorI18ns", testAuthorI18nsCount)
	t.Run("Authors", testAuthorsCount)
}

func TestInsert(t *testing.T) {
	t.Run("ContentTypes", testContentTypesInsert)
	t.Run("ContentTypes", testContentTypesInsertWhitelist)
	t.Run("ContentUnits", testContentUnitsInsert)
	t.Run("ContentUnits", testContentUnitsInsertWhitelist)
	t.Run("Files", testFilesInsert)
	t.Run("Files", testFilesInsertWhitelist)
	t.Run("OperationTypes", testOperationTypesInsert)
	t.Run("OperationTypes", testOperationTypesInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
	t.Run("Sources", testSourcesInsert)
	t.Run("Sources", testSourcesInsertWhitelist)
	t.Run("CollectionI18ns", testCollectionI18nsInsert)
	t.Run("CollectionI18ns", testCollectionI18nsInsertWhitelist)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsInsert)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsInsertWhitelist)
	t.Run("Collections", testCollectionsInsert)
	t.Run("Collections", testCollectionsInsertWhitelist)
	t.Run("ContentUnitI18ns", testContentUnitI18nsInsert)
	t.Run("ContentUnitI18ns", testContentUnitI18nsInsertWhitelist)
	t.Run("Operations", testOperationsInsert)
	t.Run("Operations", testOperationsInsertWhitelist)
	t.Run("Tags", testTagsInsert)
	t.Run("Tags", testTagsInsertWhitelist)
	t.Run("TagI18ns", testTagI18nsInsert)
	t.Run("TagI18ns", testTagI18nsInsertWhitelist)
	t.Run("SourceTypes", testSourceTypesInsert)
	t.Run("SourceTypes", testSourceTypesInsertWhitelist)
	t.Run("SourceI18ns", testSourceI18nsInsert)
	t.Run("SourceI18ns", testSourceI18nsInsertWhitelist)
	t.Run("AuthorI18ns", testAuthorI18nsInsert)
	t.Run("AuthorI18ns", testAuthorI18nsInsertWhitelist)
	t.Run("Authors", testAuthorsInsert)
	t.Run("Authors", testAuthorsInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("ContentUnitToContentTypeUsingType", testContentUnitToOneContentTypeUsingType)
	t.Run("FileToContentUnitUsingContentUnit", testFileToOneContentUnitUsingContentUnit)
	t.Run("FileToFileUsingParent", testFileToOneFileUsingParent)
	t.Run("SourceToSourceUsingParent", testSourceToOneSourceUsingParent)
	t.Run("SourceToSourceTypeUsingType", testSourceToOneSourceTypeUsingType)
	t.Run("CollectionI18nToCollectionUsingCollection", testCollectionI18nToOneCollectionUsingCollection)
	t.Run("CollectionI18nToUserUsingUser", testCollectionI18nToOneUserUsingUser)
	t.Run("CollectionsContentUnitToCollectionUsingCollection", testCollectionsContentUnitToOneCollectionUsingCollection)
	t.Run("CollectionsContentUnitToContentUnitUsingContentUnit", testCollectionsContentUnitToOneContentUnitUsingContentUnit)
	t.Run("CollectionToContentTypeUsingType", testCollectionToOneContentTypeUsingType)
	t.Run("ContentUnitI18nToContentUnitUsingContentUnit", testContentUnitI18nToOneContentUnitUsingContentUnit)
	t.Run("ContentUnitI18nToUserUsingUser", testContentUnitI18nToOneUserUsingUser)
	t.Run("OperationToOperationTypeUsingType", testOperationToOneOperationTypeUsingType)
	t.Run("OperationToUserUsingUser", testOperationToOneUserUsingUser)
	t.Run("TagToTagUsingParent", testTagToOneTagUsingParent)
	t.Run("TagI18nToTagUsingTag", testTagI18nToOneTagUsingTag)
	t.Run("TagI18nToUserUsingUser", testTagI18nToOneUserUsingUser)
	t.Run("SourceI18nToSourceUsingSource", testSourceI18nToOneSourceUsingSource)
	t.Run("AuthorI18nToAuthorUsingAuthor", testAuthorI18nToOneAuthorUsingAuthor)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("ContentTypeToTypeContentUnits", testContentTypeToManyTypeContentUnits)
	t.Run("ContentTypeToTypeCollections", testContentTypeToManyTypeCollections)
	t.Run("ContentUnitToFiles", testContentUnitToManyFiles)
	t.Run("ContentUnitToSources", testContentUnitToManySources)
	t.Run("ContentUnitToCollectionsContentUnits", testContentUnitToManyCollectionsContentUnits)
	t.Run("ContentUnitToContentUnitI18ns", testContentUnitToManyContentUnitI18ns)
	t.Run("FileToParentFiles", testFileToManyParentFiles)
	t.Run("FileToOperations", testFileToManyOperations)
	t.Run("OperationTypeToTypeOperations", testOperationTypeToManyTypeOperations)
	t.Run("UserToCollectionI18ns", testUserToManyCollectionI18ns)
	t.Run("UserToContentUnitI18ns", testUserToManyContentUnitI18ns)
	t.Run("UserToOperations", testUserToManyOperations)
	t.Run("UserToTagI18ns", testUserToManyTagI18ns)
	t.Run("SourceToParentSources", testSourceToManyParentSources)
	t.Run("SourceToContentUnits", testSourceToManyContentUnits)
	t.Run("SourceToAuthors", testSourceToManyAuthors)
	t.Run("SourceToSourceI18ns", testSourceToManySourceI18ns)
	t.Run("CollectionToCollectionI18ns", testCollectionToManyCollectionI18ns)
	t.Run("CollectionToCollectionsContentUnits", testCollectionToManyCollectionsContentUnits)
	t.Run("OperationToFiles", testOperationToManyFiles)
	t.Run("TagToParentTags", testTagToManyParentTags)
	t.Run("TagToTagI18ns", testTagToManyTagI18ns)
	t.Run("SourceTypeToTypeSources", testSourceTypeToManyTypeSources)
	t.Run("AuthorToSources", testAuthorToManySources)
	t.Run("AuthorToAuthorI18ns", testAuthorToManyAuthorI18ns)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("ContentUnitToContentTypeUsingType", testContentUnitToOneSetOpContentTypeUsingType)
	t.Run("FileToContentUnitUsingContentUnit", testFileToOneSetOpContentUnitUsingContentUnit)
	t.Run("FileToFileUsingParent", testFileToOneSetOpFileUsingParent)
	t.Run("SourceToSourceUsingParent", testSourceToOneSetOpSourceUsingParent)
	t.Run("SourceToSourceTypeUsingType", testSourceToOneSetOpSourceTypeUsingType)
	t.Run("CollectionI18nToCollectionUsingCollection", testCollectionI18nToOneSetOpCollectionUsingCollection)
	t.Run("CollectionI18nToUserUsingUser", testCollectionI18nToOneSetOpUserUsingUser)
	t.Run("CollectionsContentUnitToCollectionUsingCollection", testCollectionsContentUnitToOneSetOpCollectionUsingCollection)
	t.Run("CollectionsContentUnitToContentUnitUsingContentUnit", testCollectionsContentUnitToOneSetOpContentUnitUsingContentUnit)
	t.Run("CollectionToContentTypeUsingType", testCollectionToOneSetOpContentTypeUsingType)
	t.Run("ContentUnitI18nToContentUnitUsingContentUnit", testContentUnitI18nToOneSetOpContentUnitUsingContentUnit)
	t.Run("ContentUnitI18nToUserUsingUser", testContentUnitI18nToOneSetOpUserUsingUser)
	t.Run("OperationToOperationTypeUsingType", testOperationToOneSetOpOperationTypeUsingType)
	t.Run("OperationToUserUsingUser", testOperationToOneSetOpUserUsingUser)
	t.Run("TagToTagUsingParent", testTagToOneSetOpTagUsingParent)
	t.Run("TagI18nToTagUsingTag", testTagI18nToOneSetOpTagUsingTag)
	t.Run("TagI18nToUserUsingUser", testTagI18nToOneSetOpUserUsingUser)
	t.Run("SourceI18nToSourceUsingSource", testSourceI18nToOneSetOpSourceUsingSource)
	t.Run("AuthorI18nToAuthorUsingAuthor", testAuthorI18nToOneSetOpAuthorUsingAuthor)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("FileToContentUnitUsingContentUnit", testFileToOneRemoveOpContentUnitUsingContentUnit)
	t.Run("FileToFileUsingParent", testFileToOneRemoveOpFileUsingParent)
	t.Run("SourceToSourceUsingParent", testSourceToOneRemoveOpSourceUsingParent)
	t.Run("CollectionI18nToUserUsingUser", testCollectionI18nToOneRemoveOpUserUsingUser)
	t.Run("ContentUnitI18nToUserUsingUser", testContentUnitI18nToOneRemoveOpUserUsingUser)
	t.Run("OperationToUserUsingUser", testOperationToOneRemoveOpUserUsingUser)
	t.Run("TagToTagUsingParent", testTagToOneRemoveOpTagUsingParent)
	t.Run("TagI18nToUserUsingUser", testTagI18nToOneRemoveOpUserUsingUser)
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
	t.Run("ContentTypeToTypeContentUnits", testContentTypeToManyAddOpTypeContentUnits)
	t.Run("ContentTypeToTypeCollections", testContentTypeToManyAddOpTypeCollections)
	t.Run("ContentUnitToFiles", testContentUnitToManyAddOpFiles)
	t.Run("ContentUnitToSources", testContentUnitToManyAddOpSources)
	t.Run("ContentUnitToCollectionsContentUnits", testContentUnitToManyAddOpCollectionsContentUnits)
	t.Run("ContentUnitToContentUnitI18ns", testContentUnitToManyAddOpContentUnitI18ns)
	t.Run("FileToParentFiles", testFileToManyAddOpParentFiles)
	t.Run("FileToOperations", testFileToManyAddOpOperations)
	t.Run("OperationTypeToTypeOperations", testOperationTypeToManyAddOpTypeOperations)
	t.Run("UserToCollectionI18ns", testUserToManyAddOpCollectionI18ns)
	t.Run("UserToContentUnitI18ns", testUserToManyAddOpContentUnitI18ns)
	t.Run("UserToOperations", testUserToManyAddOpOperations)
	t.Run("UserToTagI18ns", testUserToManyAddOpTagI18ns)
	t.Run("SourceToParentSources", testSourceToManyAddOpParentSources)
	t.Run("SourceToContentUnits", testSourceToManyAddOpContentUnits)
	t.Run("SourceToAuthors", testSourceToManyAddOpAuthors)
	t.Run("SourceToSourceI18ns", testSourceToManyAddOpSourceI18ns)
	t.Run("CollectionToCollectionI18ns", testCollectionToManyAddOpCollectionI18ns)
	t.Run("CollectionToCollectionsContentUnits", testCollectionToManyAddOpCollectionsContentUnits)
	t.Run("OperationToFiles", testOperationToManyAddOpFiles)
	t.Run("TagToParentTags", testTagToManyAddOpParentTags)
	t.Run("TagToTagI18ns", testTagToManyAddOpTagI18ns)
	t.Run("SourceTypeToTypeSources", testSourceTypeToManyAddOpTypeSources)
	t.Run("AuthorToSources", testAuthorToManyAddOpSources)
	t.Run("AuthorToAuthorI18ns", testAuthorToManyAddOpAuthorI18ns)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("ContentUnitToFiles", testContentUnitToManySetOpFiles)
	t.Run("FileToParentFiles", testFileToManySetOpParentFiles)
	t.Run("UserToCollectionI18ns", testUserToManySetOpCollectionI18ns)
	t.Run("UserToContentUnitI18ns", testUserToManySetOpContentUnitI18ns)
	t.Run("UserToOperations", testUserToManySetOpOperations)
	t.Run("UserToTagI18ns", testUserToManySetOpTagI18ns)
	t.Run("SourceToParentSources", testSourceToManySetOpParentSources)
	t.Run("TagToParentTags", testTagToManySetOpParentTags)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("ContentUnitToFiles", testContentUnitToManyRemoveOpFiles)
	t.Run("FileToParentFiles", testFileToManyRemoveOpParentFiles)
	t.Run("UserToCollectionI18ns", testUserToManyRemoveOpCollectionI18ns)
	t.Run("UserToContentUnitI18ns", testUserToManyRemoveOpContentUnitI18ns)
	t.Run("UserToOperations", testUserToManyRemoveOpOperations)
	t.Run("UserToTagI18ns", testUserToManyRemoveOpTagI18ns)
	t.Run("SourceToParentSources", testSourceToManyRemoveOpParentSources)
	t.Run("TagToParentTags", testTagToManyRemoveOpParentTags)
}

func TestReload(t *testing.T) {
	t.Run("ContentTypes", testContentTypesReload)
	t.Run("ContentUnits", testContentUnitsReload)
	t.Run("Files", testFilesReload)
	t.Run("OperationTypes", testOperationTypesReload)
	t.Run("Users", testUsersReload)
	t.Run("Sources", testSourcesReload)
	t.Run("CollectionI18ns", testCollectionI18nsReload)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsReload)
	t.Run("Collections", testCollectionsReload)
	t.Run("ContentUnitI18ns", testContentUnitI18nsReload)
	t.Run("Operations", testOperationsReload)
	t.Run("Tags", testTagsReload)
	t.Run("TagI18ns", testTagI18nsReload)
	t.Run("SourceTypes", testSourceTypesReload)
	t.Run("SourceI18ns", testSourceI18nsReload)
	t.Run("AuthorI18ns", testAuthorI18nsReload)
	t.Run("Authors", testAuthorsReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("ContentTypes", testContentTypesReloadAll)
	t.Run("ContentUnits", testContentUnitsReloadAll)
	t.Run("Files", testFilesReloadAll)
	t.Run("OperationTypes", testOperationTypesReloadAll)
	t.Run("Users", testUsersReloadAll)
	t.Run("Sources", testSourcesReloadAll)
	t.Run("CollectionI18ns", testCollectionI18nsReloadAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsReloadAll)
	t.Run("Collections", testCollectionsReloadAll)
	t.Run("ContentUnitI18ns", testContentUnitI18nsReloadAll)
	t.Run("Operations", testOperationsReloadAll)
	t.Run("Tags", testTagsReloadAll)
	t.Run("TagI18ns", testTagI18nsReloadAll)
	t.Run("SourceTypes", testSourceTypesReloadAll)
	t.Run("SourceI18ns", testSourceI18nsReloadAll)
	t.Run("AuthorI18ns", testAuthorI18nsReloadAll)
	t.Run("Authors", testAuthorsReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("ContentTypes", testContentTypesSelect)
	t.Run("ContentUnits", testContentUnitsSelect)
	t.Run("Files", testFilesSelect)
	t.Run("OperationTypes", testOperationTypesSelect)
	t.Run("Users", testUsersSelect)
	t.Run("Sources", testSourcesSelect)
	t.Run("CollectionI18ns", testCollectionI18nsSelect)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsSelect)
	t.Run("Collections", testCollectionsSelect)
	t.Run("ContentUnitI18ns", testContentUnitI18nsSelect)
	t.Run("Operations", testOperationsSelect)
	t.Run("Tags", testTagsSelect)
	t.Run("TagI18ns", testTagI18nsSelect)
	t.Run("SourceTypes", testSourceTypesSelect)
	t.Run("SourceI18ns", testSourceI18nsSelect)
	t.Run("AuthorI18ns", testAuthorI18nsSelect)
	t.Run("Authors", testAuthorsSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("ContentTypes", testContentTypesUpdate)
	t.Run("ContentUnits", testContentUnitsUpdate)
	t.Run("Files", testFilesUpdate)
	t.Run("OperationTypes", testOperationTypesUpdate)
	t.Run("Users", testUsersUpdate)
	t.Run("Sources", testSourcesUpdate)
	t.Run("CollectionI18ns", testCollectionI18nsUpdate)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsUpdate)
	t.Run("Collections", testCollectionsUpdate)
	t.Run("ContentUnitI18ns", testContentUnitI18nsUpdate)
	t.Run("Operations", testOperationsUpdate)
	t.Run("Tags", testTagsUpdate)
	t.Run("TagI18ns", testTagI18nsUpdate)
	t.Run("SourceTypes", testSourceTypesUpdate)
	t.Run("SourceI18ns", testSourceI18nsUpdate)
	t.Run("AuthorI18ns", testAuthorI18nsUpdate)
	t.Run("Authors", testAuthorsUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("ContentTypes", testContentTypesSliceUpdateAll)
	t.Run("ContentUnits", testContentUnitsSliceUpdateAll)
	t.Run("Files", testFilesSliceUpdateAll)
	t.Run("OperationTypes", testOperationTypesSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
	t.Run("Sources", testSourcesSliceUpdateAll)
	t.Run("CollectionI18ns", testCollectionI18nsSliceUpdateAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsSliceUpdateAll)
	t.Run("Collections", testCollectionsSliceUpdateAll)
	t.Run("ContentUnitI18ns", testContentUnitI18nsSliceUpdateAll)
	t.Run("Operations", testOperationsSliceUpdateAll)
	t.Run("Tags", testTagsSliceUpdateAll)
	t.Run("TagI18ns", testTagI18nsSliceUpdateAll)
	t.Run("SourceTypes", testSourceTypesSliceUpdateAll)
	t.Run("SourceI18ns", testSourceI18nsSliceUpdateAll)
	t.Run("AuthorI18ns", testAuthorI18nsSliceUpdateAll)
	t.Run("Authors", testAuthorsSliceUpdateAll)
}

func TestUpsert(t *testing.T) {
	t.Run("ContentTypes", testContentTypesUpsert)
	t.Run("ContentUnits", testContentUnitsUpsert)
	t.Run("Files", testFilesUpsert)
	t.Run("OperationTypes", testOperationTypesUpsert)
	t.Run("Users", testUsersUpsert)
	t.Run("Sources", testSourcesUpsert)
	t.Run("CollectionI18ns", testCollectionI18nsUpsert)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsUpsert)
	t.Run("Collections", testCollectionsUpsert)
	t.Run("ContentUnitI18ns", testContentUnitI18nsUpsert)
	t.Run("Operations", testOperationsUpsert)
	t.Run("Tags", testTagsUpsert)
	t.Run("TagI18ns", testTagI18nsUpsert)
	t.Run("SourceTypes", testSourceTypesUpsert)
	t.Run("SourceI18ns", testSourceI18nsUpsert)
	t.Run("AuthorI18ns", testAuthorI18nsUpsert)
	t.Run("Authors", testAuthorsUpsert)
}
