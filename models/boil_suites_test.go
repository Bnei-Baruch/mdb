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
	t.Run("ContentTypes", testContentTypes)
	t.Run("Files", testFiles)
	t.Run("Collections", testCollections)
	t.Run("CollectionsContentUnits", testCollectionsContentUnits)
	t.Run("CollectionI18ns", testCollectionI18ns)
	t.Run("ContentUnitI18ns", testContentUnitI18ns)
	t.Run("Tags", testTags)
	t.Run("Users", testUsers)
	t.Run("TagsI18ns", testTagsI18ns)
	t.Run("ContentUnits", testContentUnits)
	t.Run("AuthorI18ns", testAuthorI18ns)
	t.Run("SourceTypes", testSourceTypes)
	t.Run("Authors", testAuthors)
	t.Run("Sources", testSources)
	t.Run("SourceI18ns", testSourceI18ns)
}

func TestDelete(t *testing.T) {
	t.Run("Operations", testOperationsDelete)
	t.Run("OperationTypes", testOperationTypesDelete)
	t.Run("ContentTypes", testContentTypesDelete)
	t.Run("Files", testFilesDelete)
	t.Run("Collections", testCollectionsDelete)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsDelete)
	t.Run("CollectionI18ns", testCollectionI18nsDelete)
	t.Run("ContentUnitI18ns", testContentUnitI18nsDelete)
	t.Run("Tags", testTagsDelete)
	t.Run("Users", testUsersDelete)
	t.Run("TagsI18ns", testTagsI18nsDelete)
	t.Run("ContentUnits", testContentUnitsDelete)
	t.Run("AuthorI18ns", testAuthorI18nsDelete)
	t.Run("SourceTypes", testSourceTypesDelete)
	t.Run("Authors", testAuthorsDelete)
	t.Run("Sources", testSourcesDelete)
	t.Run("SourceI18ns", testSourceI18nsDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Operations", testOperationsQueryDeleteAll)
	t.Run("OperationTypes", testOperationTypesQueryDeleteAll)
	t.Run("ContentTypes", testContentTypesQueryDeleteAll)
	t.Run("Files", testFilesQueryDeleteAll)
	t.Run("Collections", testCollectionsQueryDeleteAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsQueryDeleteAll)
	t.Run("CollectionI18ns", testCollectionI18nsQueryDeleteAll)
	t.Run("ContentUnitI18ns", testContentUnitI18nsQueryDeleteAll)
	t.Run("Tags", testTagsQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
	t.Run("TagsI18ns", testTagsI18nsQueryDeleteAll)
	t.Run("ContentUnits", testContentUnitsQueryDeleteAll)
	t.Run("AuthorI18ns", testAuthorI18nsQueryDeleteAll)
	t.Run("SourceTypes", testSourceTypesQueryDeleteAll)
	t.Run("Authors", testAuthorsQueryDeleteAll)
	t.Run("Sources", testSourcesQueryDeleteAll)
	t.Run("SourceI18ns", testSourceI18nsQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Operations", testOperationsSliceDeleteAll)
	t.Run("OperationTypes", testOperationTypesSliceDeleteAll)
	t.Run("ContentTypes", testContentTypesSliceDeleteAll)
	t.Run("Files", testFilesSliceDeleteAll)
	t.Run("Collections", testCollectionsSliceDeleteAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsSliceDeleteAll)
	t.Run("CollectionI18ns", testCollectionI18nsSliceDeleteAll)
	t.Run("ContentUnitI18ns", testContentUnitI18nsSliceDeleteAll)
	t.Run("Tags", testTagsSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
	t.Run("TagsI18ns", testTagsI18nsSliceDeleteAll)
	t.Run("ContentUnits", testContentUnitsSliceDeleteAll)
	t.Run("AuthorI18ns", testAuthorI18nsSliceDeleteAll)
	t.Run("SourceTypes", testSourceTypesSliceDeleteAll)
	t.Run("Authors", testAuthorsSliceDeleteAll)
	t.Run("Sources", testSourcesSliceDeleteAll)
	t.Run("SourceI18ns", testSourceI18nsSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Operations", testOperationsExists)
	t.Run("OperationTypes", testOperationTypesExists)
	t.Run("ContentTypes", testContentTypesExists)
	t.Run("Files", testFilesExists)
	t.Run("Collections", testCollectionsExists)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsExists)
	t.Run("CollectionI18ns", testCollectionI18nsExists)
	t.Run("ContentUnitI18ns", testContentUnitI18nsExists)
	t.Run("Tags", testTagsExists)
	t.Run("Users", testUsersExists)
	t.Run("TagsI18ns", testTagsI18nsExists)
	t.Run("ContentUnits", testContentUnitsExists)
	t.Run("AuthorI18ns", testAuthorI18nsExists)
	t.Run("SourceTypes", testSourceTypesExists)
	t.Run("Authors", testAuthorsExists)
	t.Run("Sources", testSourcesExists)
	t.Run("SourceI18ns", testSourceI18nsExists)
}

func TestFind(t *testing.T) {
	t.Run("Operations", testOperationsFind)
	t.Run("OperationTypes", testOperationTypesFind)
	t.Run("ContentTypes", testContentTypesFind)
	t.Run("Files", testFilesFind)
	t.Run("Collections", testCollectionsFind)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsFind)
	t.Run("CollectionI18ns", testCollectionI18nsFind)
	t.Run("ContentUnitI18ns", testContentUnitI18nsFind)
	t.Run("Tags", testTagsFind)
	t.Run("Users", testUsersFind)
	t.Run("TagsI18ns", testTagsI18nsFind)
	t.Run("ContentUnits", testContentUnitsFind)
	t.Run("AuthorI18ns", testAuthorI18nsFind)
	t.Run("SourceTypes", testSourceTypesFind)
	t.Run("Authors", testAuthorsFind)
	t.Run("Sources", testSourcesFind)
	t.Run("SourceI18ns", testSourceI18nsFind)
}

func TestBind(t *testing.T) {
	t.Run("Operations", testOperationsBind)
	t.Run("OperationTypes", testOperationTypesBind)
	t.Run("ContentTypes", testContentTypesBind)
	t.Run("Files", testFilesBind)
	t.Run("Collections", testCollectionsBind)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsBind)
	t.Run("CollectionI18ns", testCollectionI18nsBind)
	t.Run("ContentUnitI18ns", testContentUnitI18nsBind)
	t.Run("Tags", testTagsBind)
	t.Run("Users", testUsersBind)
	t.Run("TagsI18ns", testTagsI18nsBind)
	t.Run("ContentUnits", testContentUnitsBind)
	t.Run("AuthorI18ns", testAuthorI18nsBind)
	t.Run("SourceTypes", testSourceTypesBind)
	t.Run("Authors", testAuthorsBind)
	t.Run("Sources", testSourcesBind)
	t.Run("SourceI18ns", testSourceI18nsBind)
}

func TestOne(t *testing.T) {
	t.Run("Operations", testOperationsOne)
	t.Run("OperationTypes", testOperationTypesOne)
	t.Run("ContentTypes", testContentTypesOne)
	t.Run("Files", testFilesOne)
	t.Run("Collections", testCollectionsOne)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsOne)
	t.Run("CollectionI18ns", testCollectionI18nsOne)
	t.Run("ContentUnitI18ns", testContentUnitI18nsOne)
	t.Run("Tags", testTagsOne)
	t.Run("Users", testUsersOne)
	t.Run("TagsI18ns", testTagsI18nsOne)
	t.Run("ContentUnits", testContentUnitsOne)
	t.Run("AuthorI18ns", testAuthorI18nsOne)
	t.Run("SourceTypes", testSourceTypesOne)
	t.Run("Authors", testAuthorsOne)
	t.Run("Sources", testSourcesOne)
	t.Run("SourceI18ns", testSourceI18nsOne)
}

func TestAll(t *testing.T) {
	t.Run("Operations", testOperationsAll)
	t.Run("OperationTypes", testOperationTypesAll)
	t.Run("ContentTypes", testContentTypesAll)
	t.Run("Files", testFilesAll)
	t.Run("Collections", testCollectionsAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsAll)
	t.Run("CollectionI18ns", testCollectionI18nsAll)
	t.Run("ContentUnitI18ns", testContentUnitI18nsAll)
	t.Run("Tags", testTagsAll)
	t.Run("Users", testUsersAll)
	t.Run("TagsI18ns", testTagsI18nsAll)
	t.Run("ContentUnits", testContentUnitsAll)
	t.Run("AuthorI18ns", testAuthorI18nsAll)
	t.Run("SourceTypes", testSourceTypesAll)
	t.Run("Authors", testAuthorsAll)
	t.Run("Sources", testSourcesAll)
	t.Run("SourceI18ns", testSourceI18nsAll)
}

func TestCount(t *testing.T) {
	t.Run("Operations", testOperationsCount)
	t.Run("OperationTypes", testOperationTypesCount)
	t.Run("ContentTypes", testContentTypesCount)
	t.Run("Files", testFilesCount)
	t.Run("Collections", testCollectionsCount)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsCount)
	t.Run("CollectionI18ns", testCollectionI18nsCount)
	t.Run("ContentUnitI18ns", testContentUnitI18nsCount)
	t.Run("Tags", testTagsCount)
	t.Run("Users", testUsersCount)
	t.Run("TagsI18ns", testTagsI18nsCount)
	t.Run("ContentUnits", testContentUnitsCount)
	t.Run("AuthorI18ns", testAuthorI18nsCount)
	t.Run("SourceTypes", testSourceTypesCount)
	t.Run("Authors", testAuthorsCount)
	t.Run("Sources", testSourcesCount)
	t.Run("SourceI18ns", testSourceI18nsCount)
}

func TestInsert(t *testing.T) {
	t.Run("Operations", testOperationsInsert)
	t.Run("Operations", testOperationsInsertWhitelist)
	t.Run("OperationTypes", testOperationTypesInsert)
	t.Run("OperationTypes", testOperationTypesInsertWhitelist)
	t.Run("ContentTypes", testContentTypesInsert)
	t.Run("ContentTypes", testContentTypesInsertWhitelist)
	t.Run("Files", testFilesInsert)
	t.Run("Files", testFilesInsertWhitelist)
	t.Run("Collections", testCollectionsInsert)
	t.Run("Collections", testCollectionsInsertWhitelist)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsInsert)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsInsertWhitelist)
	t.Run("CollectionI18ns", testCollectionI18nsInsert)
	t.Run("CollectionI18ns", testCollectionI18nsInsertWhitelist)
	t.Run("ContentUnitI18ns", testContentUnitI18nsInsert)
	t.Run("ContentUnitI18ns", testContentUnitI18nsInsertWhitelist)
	t.Run("Tags", testTagsInsert)
	t.Run("Tags", testTagsInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
	t.Run("TagsI18ns", testTagsI18nsInsert)
	t.Run("TagsI18ns", testTagsI18nsInsertWhitelist)
	t.Run("ContentUnits", testContentUnitsInsert)
	t.Run("ContentUnits", testContentUnitsInsertWhitelist)
	t.Run("AuthorI18ns", testAuthorI18nsInsert)
	t.Run("AuthorI18ns", testAuthorI18nsInsertWhitelist)
	t.Run("SourceTypes", testSourceTypesInsert)
	t.Run("SourceTypes", testSourceTypesInsertWhitelist)
	t.Run("Authors", testAuthorsInsert)
	t.Run("Authors", testAuthorsInsertWhitelist)
	t.Run("Sources", testSourcesInsert)
	t.Run("Sources", testSourcesInsertWhitelist)
	t.Run("SourceI18ns", testSourceI18nsInsert)
	t.Run("SourceI18ns", testSourceI18nsInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("OperationToOperationTypeUsingType", testOperationToOneOperationTypeUsingType)
	t.Run("OperationToUserUsingUser", testOperationToOneUserUsingUser)
	t.Run("FileToContentUnitUsingContentUnit", testFileToOneContentUnitUsingContentUnit)
	t.Run("FileToFileUsingParent", testFileToOneFileUsingParent)
	t.Run("CollectionToContentTypeUsingType", testCollectionToOneContentTypeUsingType)
	t.Run("CollectionsContentUnitToCollectionUsingCollection", testCollectionsContentUnitToOneCollectionUsingCollection)
	t.Run("CollectionsContentUnitToContentUnitUsingContentUnit", testCollectionsContentUnitToOneContentUnitUsingContentUnit)
	t.Run("CollectionI18nToCollectionUsingCollection", testCollectionI18nToOneCollectionUsingCollection)
	t.Run("CollectionI18nToUserUsingUser", testCollectionI18nToOneUserUsingUser)
	t.Run("ContentUnitI18nToContentUnitUsingContentUnit", testContentUnitI18nToOneContentUnitUsingContentUnit)
	t.Run("ContentUnitI18nToUserUsingUser", testContentUnitI18nToOneUserUsingUser)
	t.Run("TagToTagUsingParent", testTagToOneTagUsingParent)
	t.Run("TagsI18nToTagUsingTag", testTagsI18nToOneTagUsingTag)
	t.Run("TagsI18nToUserUsingUser", testTagsI18nToOneUserUsingUser)
	t.Run("ContentUnitToContentTypeUsingType", testContentUnitToOneContentTypeUsingType)
	t.Run("AuthorI18nToAuthorUsingAuthor", testAuthorI18nToOneAuthorUsingAuthor)
	t.Run("SourceToSourceUsingParent", testSourceToOneSourceUsingParent)
	t.Run("SourceToSourceTypeUsingType", testSourceToOneSourceTypeUsingType)
	t.Run("SourceI18nToSourceUsingSource", testSourceI18nToOneSourceUsingSource)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("OperationToFiles", testOperationToManyFiles)
	t.Run("OperationTypeToTypeOperations", testOperationTypeToManyTypeOperations)
	t.Run("ContentTypeToTypeCollections", testContentTypeToManyTypeCollections)
	t.Run("ContentTypeToTypeContentUnits", testContentTypeToManyTypeContentUnits)
	t.Run("FileToParentFiles", testFileToManyParentFiles)
	t.Run("FileToOperations", testFileToManyOperations)
	t.Run("CollectionToCollectionsContentUnits", testCollectionToManyCollectionsContentUnits)
	t.Run("CollectionToCollectionI18ns", testCollectionToManyCollectionI18ns)
	t.Run("TagToParentTags", testTagToManyParentTags)
	t.Run("TagToTagsI18ns", testTagToManyTagsI18ns)
	t.Run("UserToOperations", testUserToManyOperations)
	t.Run("UserToCollectionI18ns", testUserToManyCollectionI18ns)
	t.Run("UserToContentUnitI18ns", testUserToManyContentUnitI18ns)
	t.Run("UserToTagsI18ns", testUserToManyTagsI18ns)
	t.Run("ContentUnitToFiles", testContentUnitToManyFiles)
	t.Run("ContentUnitToCollectionsContentUnits", testContentUnitToManyCollectionsContentUnits)
	t.Run("ContentUnitToContentUnitI18ns", testContentUnitToManyContentUnitI18ns)
	t.Run("ContentUnitToSources", testContentUnitToManySources)
	t.Run("SourceTypeToTypeSources", testSourceTypeToManyTypeSources)
	t.Run("AuthorToAuthorI18ns", testAuthorToManyAuthorI18ns)
	t.Run("AuthorToSources", testAuthorToManySources)
	t.Run("SourceToParentSources", testSourceToManyParentSources)
	t.Run("SourceToSourceI18ns", testSourceToManySourceI18ns)
	t.Run("SourceToContentUnits", testSourceToManyContentUnits)
	t.Run("SourceToAuthors", testSourceToManyAuthors)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("OperationToOperationTypeUsingType", testOperationToOneSetOpOperationTypeUsingType)
	t.Run("OperationToUserUsingUser", testOperationToOneSetOpUserUsingUser)
	t.Run("FileToContentUnitUsingContentUnit", testFileToOneSetOpContentUnitUsingContentUnit)
	t.Run("FileToFileUsingParent", testFileToOneSetOpFileUsingParent)
	t.Run("CollectionToContentTypeUsingType", testCollectionToOneSetOpContentTypeUsingType)
	t.Run("CollectionsContentUnitToCollectionUsingCollection", testCollectionsContentUnitToOneSetOpCollectionUsingCollection)
	t.Run("CollectionsContentUnitToContentUnitUsingContentUnit", testCollectionsContentUnitToOneSetOpContentUnitUsingContentUnit)
	t.Run("CollectionI18nToCollectionUsingCollection", testCollectionI18nToOneSetOpCollectionUsingCollection)
	t.Run("CollectionI18nToUserUsingUser", testCollectionI18nToOneSetOpUserUsingUser)
	t.Run("ContentUnitI18nToContentUnitUsingContentUnit", testContentUnitI18nToOneSetOpContentUnitUsingContentUnit)
	t.Run("ContentUnitI18nToUserUsingUser", testContentUnitI18nToOneSetOpUserUsingUser)
	t.Run("TagToTagUsingParent", testTagToOneSetOpTagUsingParent)
	t.Run("TagsI18nToTagUsingTag", testTagsI18nToOneSetOpTagUsingTag)
	t.Run("TagsI18nToUserUsingUser", testTagsI18nToOneSetOpUserUsingUser)
	t.Run("ContentUnitToContentTypeUsingType", testContentUnitToOneSetOpContentTypeUsingType)
	t.Run("AuthorI18nToAuthorUsingAuthor", testAuthorI18nToOneSetOpAuthorUsingAuthor)
	t.Run("SourceToSourceUsingParent", testSourceToOneSetOpSourceUsingParent)
	t.Run("SourceToSourceTypeUsingType", testSourceToOneSetOpSourceTypeUsingType)
	t.Run("SourceI18nToSourceUsingSource", testSourceI18nToOneSetOpSourceUsingSource)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("OperationToUserUsingUser", testOperationToOneRemoveOpUserUsingUser)
	t.Run("FileToContentUnitUsingContentUnit", testFileToOneRemoveOpContentUnitUsingContentUnit)
	t.Run("FileToFileUsingParent", testFileToOneRemoveOpFileUsingParent)
	t.Run("CollectionI18nToUserUsingUser", testCollectionI18nToOneRemoveOpUserUsingUser)
	t.Run("ContentUnitI18nToUserUsingUser", testContentUnitI18nToOneRemoveOpUserUsingUser)
	t.Run("TagToTagUsingParent", testTagToOneRemoveOpTagUsingParent)
	t.Run("TagsI18nToUserUsingUser", testTagsI18nToOneRemoveOpUserUsingUser)
	t.Run("SourceToSourceUsingParent", testSourceToOneRemoveOpSourceUsingParent)
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
	t.Run("ContentTypeToTypeCollections", testContentTypeToManyAddOpTypeCollections)
	t.Run("ContentTypeToTypeContentUnits", testContentTypeToManyAddOpTypeContentUnits)
	t.Run("FileToParentFiles", testFileToManyAddOpParentFiles)
	t.Run("FileToOperations", testFileToManyAddOpOperations)
	t.Run("CollectionToCollectionsContentUnits", testCollectionToManyAddOpCollectionsContentUnits)
	t.Run("CollectionToCollectionI18ns", testCollectionToManyAddOpCollectionI18ns)
	t.Run("TagToParentTags", testTagToManyAddOpParentTags)
	t.Run("TagToTagsI18ns", testTagToManyAddOpTagsI18ns)
	t.Run("UserToOperations", testUserToManyAddOpOperations)
	t.Run("UserToCollectionI18ns", testUserToManyAddOpCollectionI18ns)
	t.Run("UserToContentUnitI18ns", testUserToManyAddOpContentUnitI18ns)
	t.Run("UserToTagsI18ns", testUserToManyAddOpTagsI18ns)
	t.Run("ContentUnitToFiles", testContentUnitToManyAddOpFiles)
	t.Run("ContentUnitToCollectionsContentUnits", testContentUnitToManyAddOpCollectionsContentUnits)
	t.Run("ContentUnitToContentUnitI18ns", testContentUnitToManyAddOpContentUnitI18ns)
	t.Run("ContentUnitToSources", testContentUnitToManyAddOpSources)
	t.Run("SourceTypeToTypeSources", testSourceTypeToManyAddOpTypeSources)
	t.Run("AuthorToAuthorI18ns", testAuthorToManyAddOpAuthorI18ns)
	t.Run("AuthorToSources", testAuthorToManyAddOpSources)
	t.Run("SourceToParentSources", testSourceToManyAddOpParentSources)
	t.Run("SourceToSourceI18ns", testSourceToManyAddOpSourceI18ns)
	t.Run("SourceToContentUnits", testSourceToManyAddOpContentUnits)
	t.Run("SourceToAuthors", testSourceToManyAddOpAuthors)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("FileToParentFiles", testFileToManySetOpParentFiles)
	t.Run("TagToParentTags", testTagToManySetOpParentTags)
	t.Run("UserToOperations", testUserToManySetOpOperations)
	t.Run("UserToCollectionI18ns", testUserToManySetOpCollectionI18ns)
	t.Run("UserToContentUnitI18ns", testUserToManySetOpContentUnitI18ns)
	t.Run("UserToTagsI18ns", testUserToManySetOpTagsI18ns)
	t.Run("ContentUnitToFiles", testContentUnitToManySetOpFiles)
	t.Run("SourceToParentSources", testSourceToManySetOpParentSources)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("FileToParentFiles", testFileToManyRemoveOpParentFiles)
	t.Run("TagToParentTags", testTagToManyRemoveOpParentTags)
	t.Run("UserToOperations", testUserToManyRemoveOpOperations)
	t.Run("UserToCollectionI18ns", testUserToManyRemoveOpCollectionI18ns)
	t.Run("UserToContentUnitI18ns", testUserToManyRemoveOpContentUnitI18ns)
	t.Run("UserToTagsI18ns", testUserToManyRemoveOpTagsI18ns)
	t.Run("ContentUnitToFiles", testContentUnitToManyRemoveOpFiles)
	t.Run("SourceToParentSources", testSourceToManyRemoveOpParentSources)
}

func TestReload(t *testing.T) {
	t.Run("Operations", testOperationsReload)
	t.Run("OperationTypes", testOperationTypesReload)
	t.Run("ContentTypes", testContentTypesReload)
	t.Run("Files", testFilesReload)
	t.Run("Collections", testCollectionsReload)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsReload)
	t.Run("CollectionI18ns", testCollectionI18nsReload)
	t.Run("ContentUnitI18ns", testContentUnitI18nsReload)
	t.Run("Tags", testTagsReload)
	t.Run("Users", testUsersReload)
	t.Run("TagsI18ns", testTagsI18nsReload)
	t.Run("ContentUnits", testContentUnitsReload)
	t.Run("AuthorI18ns", testAuthorI18nsReload)
	t.Run("SourceTypes", testSourceTypesReload)
	t.Run("Authors", testAuthorsReload)
	t.Run("Sources", testSourcesReload)
	t.Run("SourceI18ns", testSourceI18nsReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Operations", testOperationsReloadAll)
	t.Run("OperationTypes", testOperationTypesReloadAll)
	t.Run("ContentTypes", testContentTypesReloadAll)
	t.Run("Files", testFilesReloadAll)
	t.Run("Collections", testCollectionsReloadAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsReloadAll)
	t.Run("CollectionI18ns", testCollectionI18nsReloadAll)
	t.Run("ContentUnitI18ns", testContentUnitI18nsReloadAll)
	t.Run("Tags", testTagsReloadAll)
	t.Run("Users", testUsersReloadAll)
	t.Run("TagsI18ns", testTagsI18nsReloadAll)
	t.Run("ContentUnits", testContentUnitsReloadAll)
	t.Run("AuthorI18ns", testAuthorI18nsReloadAll)
	t.Run("SourceTypes", testSourceTypesReloadAll)
	t.Run("Authors", testAuthorsReloadAll)
	t.Run("Sources", testSourcesReloadAll)
	t.Run("SourceI18ns", testSourceI18nsReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Operations", testOperationsSelect)
	t.Run("OperationTypes", testOperationTypesSelect)
	t.Run("ContentTypes", testContentTypesSelect)
	t.Run("Files", testFilesSelect)
	t.Run("Collections", testCollectionsSelect)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsSelect)
	t.Run("CollectionI18ns", testCollectionI18nsSelect)
	t.Run("ContentUnitI18ns", testContentUnitI18nsSelect)
	t.Run("Tags", testTagsSelect)
	t.Run("Users", testUsersSelect)
	t.Run("TagsI18ns", testTagsI18nsSelect)
	t.Run("ContentUnits", testContentUnitsSelect)
	t.Run("AuthorI18ns", testAuthorI18nsSelect)
	t.Run("SourceTypes", testSourceTypesSelect)
	t.Run("Authors", testAuthorsSelect)
	t.Run("Sources", testSourcesSelect)
	t.Run("SourceI18ns", testSourceI18nsSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Operations", testOperationsUpdate)
	t.Run("OperationTypes", testOperationTypesUpdate)
	t.Run("ContentTypes", testContentTypesUpdate)
	t.Run("Files", testFilesUpdate)
	t.Run("Collections", testCollectionsUpdate)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsUpdate)
	t.Run("CollectionI18ns", testCollectionI18nsUpdate)
	t.Run("ContentUnitI18ns", testContentUnitI18nsUpdate)
	t.Run("Tags", testTagsUpdate)
	t.Run("Users", testUsersUpdate)
	t.Run("TagsI18ns", testTagsI18nsUpdate)
	t.Run("ContentUnits", testContentUnitsUpdate)
	t.Run("AuthorI18ns", testAuthorI18nsUpdate)
	t.Run("SourceTypes", testSourceTypesUpdate)
	t.Run("Authors", testAuthorsUpdate)
	t.Run("Sources", testSourcesUpdate)
	t.Run("SourceI18ns", testSourceI18nsUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Operations", testOperationsSliceUpdateAll)
	t.Run("OperationTypes", testOperationTypesSliceUpdateAll)
	t.Run("ContentTypes", testContentTypesSliceUpdateAll)
	t.Run("Files", testFilesSliceUpdateAll)
	t.Run("Collections", testCollectionsSliceUpdateAll)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsSliceUpdateAll)
	t.Run("CollectionI18ns", testCollectionI18nsSliceUpdateAll)
	t.Run("ContentUnitI18ns", testContentUnitI18nsSliceUpdateAll)
	t.Run("Tags", testTagsSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
	t.Run("TagsI18ns", testTagsI18nsSliceUpdateAll)
	t.Run("ContentUnits", testContentUnitsSliceUpdateAll)
	t.Run("AuthorI18ns", testAuthorI18nsSliceUpdateAll)
	t.Run("SourceTypes", testSourceTypesSliceUpdateAll)
	t.Run("Authors", testAuthorsSliceUpdateAll)
	t.Run("Sources", testSourcesSliceUpdateAll)
	t.Run("SourceI18ns", testSourceI18nsSliceUpdateAll)
}

func TestUpsert(t *testing.T) {
	t.Run("Operations", testOperationsUpsert)
	t.Run("OperationTypes", testOperationTypesUpsert)
	t.Run("ContentTypes", testContentTypesUpsert)
	t.Run("Files", testFilesUpsert)
	t.Run("Collections", testCollectionsUpsert)
	t.Run("CollectionsContentUnits", testCollectionsContentUnitsUpsert)
	t.Run("CollectionI18ns", testCollectionI18nsUpsert)
	t.Run("ContentUnitI18ns", testContentUnitI18nsUpsert)
	t.Run("Tags", testTagsUpsert)
	t.Run("Users", testUsersUpsert)
	t.Run("TagsI18ns", testTagsI18nsUpsert)
	t.Run("ContentUnits", testContentUnitsUpsert)
	t.Run("AuthorI18ns", testAuthorI18nsUpsert)
	t.Run("SourceTypes", testSourceTypesUpsert)
	t.Run("Authors", testAuthorsUpsert)
	t.Run("Sources", testSourcesUpsert)
	t.Run("SourceI18ns", testSourceI18nsUpsert)
}
