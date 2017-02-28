package kmodels

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("FileTypes", testFileTypes)
	t.Run("LecturerDescriptions", testLecturerDescriptions)
	t.Run("Lecturers", testLecturers)
	t.Run("ContentTypes", testContentTypes)
	t.Run("Roles", testRoles)
	t.Run("CatalogDescriptions", testCatalogDescriptions)
	t.Run("Catalogs", testCatalogs)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatterns)
	t.Run("ContainerDescriptions", testContainerDescriptions)
	t.Run("Labels", testLabels)
	t.Run("RolesUsers", testRolesUsers)
	t.Run("Containers", testContainers)
	t.Run("FileAssetDescriptions", testFileAssetDescriptions)
	t.Run("Servers", testServers)
	t.Run("VirtualLessons", testVirtualLessons)
	t.Run("Languages", testLanguages)
	t.Run("FileAssets", testFileAssets)
	t.Run("Users", testUsers)
}

func TestDelete(t *testing.T) {
	t.Run("FileTypes", testFileTypesDelete)
	t.Run("LecturerDescriptions", testLecturerDescriptionsDelete)
	t.Run("Lecturers", testLecturersDelete)
	t.Run("ContentTypes", testContentTypesDelete)
	t.Run("Roles", testRolesDelete)
	t.Run("CatalogDescriptions", testCatalogDescriptionsDelete)
	t.Run("Catalogs", testCatalogsDelete)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsDelete)
	t.Run("ContainerDescriptions", testContainerDescriptionsDelete)
	t.Run("Labels", testLabelsDelete)
	t.Run("RolesUsers", testRolesUsersDelete)
	t.Run("Containers", testContainersDelete)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsDelete)
	t.Run("Servers", testServersDelete)
	t.Run("VirtualLessons", testVirtualLessonsDelete)
	t.Run("Languages", testLanguagesDelete)
	t.Run("FileAssets", testFileAssetsDelete)
	t.Run("Users", testUsersDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("FileTypes", testFileTypesQueryDeleteAll)
	t.Run("LecturerDescriptions", testLecturerDescriptionsQueryDeleteAll)
	t.Run("Lecturers", testLecturersQueryDeleteAll)
	t.Run("ContentTypes", testContentTypesQueryDeleteAll)
	t.Run("Roles", testRolesQueryDeleteAll)
	t.Run("CatalogDescriptions", testCatalogDescriptionsQueryDeleteAll)
	t.Run("Catalogs", testCatalogsQueryDeleteAll)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsQueryDeleteAll)
	t.Run("ContainerDescriptions", testContainerDescriptionsQueryDeleteAll)
	t.Run("Labels", testLabelsQueryDeleteAll)
	t.Run("RolesUsers", testRolesUsersQueryDeleteAll)
	t.Run("Containers", testContainersQueryDeleteAll)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsQueryDeleteAll)
	t.Run("Servers", testServersQueryDeleteAll)
	t.Run("VirtualLessons", testVirtualLessonsQueryDeleteAll)
	t.Run("Languages", testLanguagesQueryDeleteAll)
	t.Run("FileAssets", testFileAssetsQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("FileTypes", testFileTypesSliceDeleteAll)
	t.Run("LecturerDescriptions", testLecturerDescriptionsSliceDeleteAll)
	t.Run("Lecturers", testLecturersSliceDeleteAll)
	t.Run("ContentTypes", testContentTypesSliceDeleteAll)
	t.Run("Roles", testRolesSliceDeleteAll)
	t.Run("CatalogDescriptions", testCatalogDescriptionsSliceDeleteAll)
	t.Run("Catalogs", testCatalogsSliceDeleteAll)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsSliceDeleteAll)
	t.Run("ContainerDescriptions", testContainerDescriptionsSliceDeleteAll)
	t.Run("Labels", testLabelsSliceDeleteAll)
	t.Run("RolesUsers", testRolesUsersSliceDeleteAll)
	t.Run("Containers", testContainersSliceDeleteAll)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsSliceDeleteAll)
	t.Run("Servers", testServersSliceDeleteAll)
	t.Run("VirtualLessons", testVirtualLessonsSliceDeleteAll)
	t.Run("Languages", testLanguagesSliceDeleteAll)
	t.Run("FileAssets", testFileAssetsSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("FileTypes", testFileTypesExists)
	t.Run("LecturerDescriptions", testLecturerDescriptionsExists)
	t.Run("Lecturers", testLecturersExists)
	t.Run("ContentTypes", testContentTypesExists)
	t.Run("Roles", testRolesExists)
	t.Run("CatalogDescriptions", testCatalogDescriptionsExists)
	t.Run("Catalogs", testCatalogsExists)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsExists)
	t.Run("ContainerDescriptions", testContainerDescriptionsExists)
	t.Run("Labels", testLabelsExists)
	t.Run("RolesUsers", testRolesUsersExists)
	t.Run("Containers", testContainersExists)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsExists)
	t.Run("Servers", testServersExists)
	t.Run("VirtualLessons", testVirtualLessonsExists)
	t.Run("Languages", testLanguagesExists)
	t.Run("FileAssets", testFileAssetsExists)
	t.Run("Users", testUsersExists)
}

func TestFind(t *testing.T) {
	t.Run("FileTypes", testFileTypesFind)
	t.Run("LecturerDescriptions", testLecturerDescriptionsFind)
	t.Run("Lecturers", testLecturersFind)
	t.Run("ContentTypes", testContentTypesFind)
	t.Run("Roles", testRolesFind)
	t.Run("CatalogDescriptions", testCatalogDescriptionsFind)
	t.Run("Catalogs", testCatalogsFind)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsFind)
	t.Run("ContainerDescriptions", testContainerDescriptionsFind)
	t.Run("Labels", testLabelsFind)
	t.Run("RolesUsers", testRolesUsersFind)
	t.Run("Containers", testContainersFind)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsFind)
	t.Run("Servers", testServersFind)
	t.Run("VirtualLessons", testVirtualLessonsFind)
	t.Run("Languages", testLanguagesFind)
	t.Run("FileAssets", testFileAssetsFind)
	t.Run("Users", testUsersFind)
}

func TestBind(t *testing.T) {
	t.Run("FileTypes", testFileTypesBind)
	t.Run("LecturerDescriptions", testLecturerDescriptionsBind)
	t.Run("Lecturers", testLecturersBind)
	t.Run("ContentTypes", testContentTypesBind)
	t.Run("Roles", testRolesBind)
	t.Run("CatalogDescriptions", testCatalogDescriptionsBind)
	t.Run("Catalogs", testCatalogsBind)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsBind)
	t.Run("ContainerDescriptions", testContainerDescriptionsBind)
	t.Run("Labels", testLabelsBind)
	t.Run("RolesUsers", testRolesUsersBind)
	t.Run("Containers", testContainersBind)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsBind)
	t.Run("Servers", testServersBind)
	t.Run("VirtualLessons", testVirtualLessonsBind)
	t.Run("Languages", testLanguagesBind)
	t.Run("FileAssets", testFileAssetsBind)
	t.Run("Users", testUsersBind)
}

func TestOne(t *testing.T) {
	t.Run("FileTypes", testFileTypesOne)
	t.Run("LecturerDescriptions", testLecturerDescriptionsOne)
	t.Run("Lecturers", testLecturersOne)
	t.Run("ContentTypes", testContentTypesOne)
	t.Run("Roles", testRolesOne)
	t.Run("CatalogDescriptions", testCatalogDescriptionsOne)
	t.Run("Catalogs", testCatalogsOne)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsOne)
	t.Run("ContainerDescriptions", testContainerDescriptionsOne)
	t.Run("Labels", testLabelsOne)
	t.Run("RolesUsers", testRolesUsersOne)
	t.Run("Containers", testContainersOne)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsOne)
	t.Run("Servers", testServersOne)
	t.Run("VirtualLessons", testVirtualLessonsOne)
	t.Run("Languages", testLanguagesOne)
	t.Run("FileAssets", testFileAssetsOne)
	t.Run("Users", testUsersOne)
}

func TestAll(t *testing.T) {
	t.Run("FileTypes", testFileTypesAll)
	t.Run("LecturerDescriptions", testLecturerDescriptionsAll)
	t.Run("Lecturers", testLecturersAll)
	t.Run("ContentTypes", testContentTypesAll)
	t.Run("Roles", testRolesAll)
	t.Run("CatalogDescriptions", testCatalogDescriptionsAll)
	t.Run("Catalogs", testCatalogsAll)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsAll)
	t.Run("ContainerDescriptions", testContainerDescriptionsAll)
	t.Run("Labels", testLabelsAll)
	t.Run("RolesUsers", testRolesUsersAll)
	t.Run("Containers", testContainersAll)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsAll)
	t.Run("Servers", testServersAll)
	t.Run("VirtualLessons", testVirtualLessonsAll)
	t.Run("Languages", testLanguagesAll)
	t.Run("FileAssets", testFileAssetsAll)
	t.Run("Users", testUsersAll)
}

func TestCount(t *testing.T) {
	t.Run("FileTypes", testFileTypesCount)
	t.Run("LecturerDescriptions", testLecturerDescriptionsCount)
	t.Run("Lecturers", testLecturersCount)
	t.Run("ContentTypes", testContentTypesCount)
	t.Run("Roles", testRolesCount)
	t.Run("CatalogDescriptions", testCatalogDescriptionsCount)
	t.Run("Catalogs", testCatalogsCount)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsCount)
	t.Run("ContainerDescriptions", testContainerDescriptionsCount)
	t.Run("Labels", testLabelsCount)
	t.Run("RolesUsers", testRolesUsersCount)
	t.Run("Containers", testContainersCount)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsCount)
	t.Run("Servers", testServersCount)
	t.Run("VirtualLessons", testVirtualLessonsCount)
	t.Run("Languages", testLanguagesCount)
	t.Run("FileAssets", testFileAssetsCount)
	t.Run("Users", testUsersCount)
}

func TestInsert(t *testing.T) {
	t.Run("FileTypes", testFileTypesInsert)
	t.Run("FileTypes", testFileTypesInsertWhitelist)
	t.Run("LecturerDescriptions", testLecturerDescriptionsInsert)
	t.Run("LecturerDescriptions", testLecturerDescriptionsInsertWhitelist)
	t.Run("Lecturers", testLecturersInsert)
	t.Run("Lecturers", testLecturersInsertWhitelist)
	t.Run("ContentTypes", testContentTypesInsert)
	t.Run("ContentTypes", testContentTypesInsertWhitelist)
	t.Run("Roles", testRolesInsert)
	t.Run("Roles", testRolesInsertWhitelist)
	t.Run("CatalogDescriptions", testCatalogDescriptionsInsert)
	t.Run("CatalogDescriptions", testCatalogDescriptionsInsertWhitelist)
	t.Run("Catalogs", testCatalogsInsert)
	t.Run("Catalogs", testCatalogsInsertWhitelist)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsInsert)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsInsertWhitelist)
	t.Run("ContainerDescriptions", testContainerDescriptionsInsert)
	t.Run("ContainerDescriptions", testContainerDescriptionsInsertWhitelist)
	t.Run("Labels", testLabelsInsert)
	t.Run("Labels", testLabelsInsertWhitelist)
	t.Run("RolesUsers", testRolesUsersInsert)
	t.Run("RolesUsers", testRolesUsersInsertWhitelist)
	t.Run("Containers", testContainersInsert)
	t.Run("Containers", testContainersInsertWhitelist)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsInsert)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsInsertWhitelist)
	t.Run("Servers", testServersInsert)
	t.Run("Servers", testServersInsertWhitelist)
	t.Run("VirtualLessons", testVirtualLessonsInsert)
	t.Run("VirtualLessons", testVirtualLessonsInsertWhitelist)
	t.Run("Languages", testLanguagesInsert)
	t.Run("Languages", testLanguagesInsertWhitelist)
	t.Run("FileAssets", testFileAssetsInsert)
	t.Run("FileAssets", testFileAssetsInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("CatalogDescriptionToCatalogUsingCatalog", testCatalogDescriptionToOneCatalogUsingCatalog)
	t.Run("CatalogDescriptionToLanguageUsingLang", testCatalogDescriptionToOneLanguageUsingLang)
	t.Run("CatalogToCatalogUsingParent", testCatalogToOneCatalogUsingParent)
	t.Run("CatalogToUserUsingUser", testCatalogToOneUserUsingUser)
	t.Run("ContainerDescriptionToLanguageUsingLang", testContainerDescriptionToOneLanguageUsingLang)
	t.Run("ContainerDescriptionToContainerUsingContainer", testContainerDescriptionToOneContainerUsingContainer)
	t.Run("RolesUserToUserUsingUser", testRolesUserToOneUserUsingUser)
	t.Run("ContainerToLanguageUsingLang", testContainerToOneLanguageUsingLang)
	t.Run("ContainerToContentTypeUsingContentType", testContainerToOneContentTypeUsingContentType)
	t.Run("ContainerToVirtualLessonUsingVirtualLesson", testContainerToOneVirtualLessonUsingVirtualLesson)
	t.Run("FileAssetDescriptionToFileAssetUsingFile", testFileAssetDescriptionToOneFileAssetUsingFile)
	t.Run("FileAssetToLanguageUsingLang", testFileAssetToOneLanguageUsingLang)
	t.Run("FileAssetToUserUsingUser", testFileAssetToOneUserUsingUser)
	t.Run("FileAssetToServerUsingServername", testFileAssetToOneServerUsingServername)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("ContentTypeToContainers", testContentTypeToManyContainers)
	t.Run("CatalogToContainers", testCatalogToManyContainers)
	t.Run("CatalogToCatalogDescriptions", testCatalogToManyCatalogDescriptions)
	t.Run("CatalogToParentCatalogs", testCatalogToManyParentCatalogs)
	t.Run("CatalogToContainerDescriptionPatterns", testCatalogToManyContainerDescriptionPatterns)
	t.Run("ContainerDescriptionPatternToCatalogs", testContainerDescriptionPatternToManyCatalogs)
	t.Run("LabelToContainers", testLabelToManyContainers)
	t.Run("ContainerToCatalogs", testContainerToManyCatalogs)
	t.Run("ContainerToFileAssets", testContainerToManyFileAssets)
	t.Run("ContainerToContainerDescriptions", testContainerToManyContainerDescriptions)
	t.Run("ContainerToLabels", testContainerToManyLabels)
	t.Run("ServerToServernameFileAssets", testServerToManyServernameFileAssets)
	t.Run("VirtualLessonToContainers", testVirtualLessonToManyContainers)
	t.Run("LanguageToLangCatalogDescriptions", testLanguageToManyLangCatalogDescriptions)
	t.Run("LanguageToLangContainerDescriptions", testLanguageToManyLangContainerDescriptions)
	t.Run("LanguageToLangContainers", testLanguageToManyLangContainers)
	t.Run("LanguageToLangFileAssets", testLanguageToManyLangFileAssets)
	t.Run("FileAssetToContainers", testFileAssetToManyContainers)
	t.Run("FileAssetToFileFileAssetDescriptions", testFileAssetToManyFileFileAssetDescriptions)
	t.Run("UserToCatalogs", testUserToManyCatalogs)
	t.Run("UserToRolesUsers", testUserToManyRolesUsers)
	t.Run("UserToFileAssets", testUserToManyFileAssets)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("CatalogDescriptionToCatalogUsingCatalog", testCatalogDescriptionToOneSetOpCatalogUsingCatalog)
	t.Run("CatalogDescriptionToLanguageUsingLang", testCatalogDescriptionToOneSetOpLanguageUsingLang)
	t.Run("CatalogToCatalogUsingParent", testCatalogToOneSetOpCatalogUsingParent)
	t.Run("CatalogToUserUsingUser", testCatalogToOneSetOpUserUsingUser)
	t.Run("ContainerDescriptionToLanguageUsingLang", testContainerDescriptionToOneSetOpLanguageUsingLang)
	t.Run("ContainerDescriptionToContainerUsingContainer", testContainerDescriptionToOneSetOpContainerUsingContainer)
	t.Run("RolesUserToUserUsingUser", testRolesUserToOneSetOpUserUsingUser)
	t.Run("ContainerToLanguageUsingLang", testContainerToOneSetOpLanguageUsingLang)
	t.Run("ContainerToContentTypeUsingContentType", testContainerToOneSetOpContentTypeUsingContentType)
	t.Run("ContainerToVirtualLessonUsingVirtualLesson", testContainerToOneSetOpVirtualLessonUsingVirtualLesson)
	t.Run("FileAssetDescriptionToFileAssetUsingFile", testFileAssetDescriptionToOneSetOpFileAssetUsingFile)
	t.Run("FileAssetToLanguageUsingLang", testFileAssetToOneSetOpLanguageUsingLang)
	t.Run("FileAssetToUserUsingUser", testFileAssetToOneSetOpUserUsingUser)
	t.Run("FileAssetToServerUsingServername", testFileAssetToOneSetOpServerUsingServername)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("CatalogDescriptionToLanguageUsingLang", testCatalogDescriptionToOneRemoveOpLanguageUsingLang)
	t.Run("CatalogToCatalogUsingParent", testCatalogToOneRemoveOpCatalogUsingParent)
	t.Run("CatalogToUserUsingUser", testCatalogToOneRemoveOpUserUsingUser)
	t.Run("ContainerDescriptionToLanguageUsingLang", testContainerDescriptionToOneRemoveOpLanguageUsingLang)
	t.Run("ContainerToLanguageUsingLang", testContainerToOneRemoveOpLanguageUsingLang)
	t.Run("ContainerToContentTypeUsingContentType", testContainerToOneRemoveOpContentTypeUsingContentType)
	t.Run("ContainerToVirtualLessonUsingVirtualLesson", testContainerToOneRemoveOpVirtualLessonUsingVirtualLesson)
	t.Run("FileAssetToLanguageUsingLang", testFileAssetToOneRemoveOpLanguageUsingLang)
	t.Run("FileAssetToUserUsingUser", testFileAssetToOneRemoveOpUserUsingUser)
	t.Run("FileAssetToServerUsingServername", testFileAssetToOneRemoveOpServerUsingServername)
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
	t.Run("ContentTypeToContainers", testContentTypeToManyAddOpContainers)
	t.Run("CatalogToContainers", testCatalogToManyAddOpContainers)
	t.Run("CatalogToCatalogDescriptions", testCatalogToManyAddOpCatalogDescriptions)
	t.Run("CatalogToParentCatalogs", testCatalogToManyAddOpParentCatalogs)
	t.Run("CatalogToContainerDescriptionPatterns", testCatalogToManyAddOpContainerDescriptionPatterns)
	t.Run("ContainerDescriptionPatternToCatalogs", testContainerDescriptionPatternToManyAddOpCatalogs)
	t.Run("LabelToContainers", testLabelToManyAddOpContainers)
	t.Run("ContainerToCatalogs", testContainerToManyAddOpCatalogs)
	t.Run("ContainerToFileAssets", testContainerToManyAddOpFileAssets)
	t.Run("ContainerToContainerDescriptions", testContainerToManyAddOpContainerDescriptions)
	t.Run("ContainerToLabels", testContainerToManyAddOpLabels)
	t.Run("ServerToServernameFileAssets", testServerToManyAddOpServernameFileAssets)
	t.Run("VirtualLessonToContainers", testVirtualLessonToManyAddOpContainers)
	t.Run("LanguageToLangCatalogDescriptions", testLanguageToManyAddOpLangCatalogDescriptions)
	t.Run("LanguageToLangContainerDescriptions", testLanguageToManyAddOpLangContainerDescriptions)
	t.Run("LanguageToLangContainers", testLanguageToManyAddOpLangContainers)
	t.Run("LanguageToLangFileAssets", testLanguageToManyAddOpLangFileAssets)
	t.Run("FileAssetToContainers", testFileAssetToManyAddOpContainers)
	t.Run("FileAssetToFileFileAssetDescriptions", testFileAssetToManyAddOpFileFileAssetDescriptions)
	t.Run("UserToCatalogs", testUserToManyAddOpCatalogs)
	t.Run("UserToRolesUsers", testUserToManyAddOpRolesUsers)
	t.Run("UserToFileAssets", testUserToManyAddOpFileAssets)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("ContentTypeToContainers", testContentTypeToManySetOpContainers)
	t.Run("CatalogToParentCatalogs", testCatalogToManySetOpParentCatalogs)
	t.Run("ServerToServernameFileAssets", testServerToManySetOpServernameFileAssets)
	t.Run("VirtualLessonToContainers", testVirtualLessonToManySetOpContainers)
	t.Run("LanguageToLangCatalogDescriptions", testLanguageToManySetOpLangCatalogDescriptions)
	t.Run("LanguageToLangContainerDescriptions", testLanguageToManySetOpLangContainerDescriptions)
	t.Run("LanguageToLangContainers", testLanguageToManySetOpLangContainers)
	t.Run("LanguageToLangFileAssets", testLanguageToManySetOpLangFileAssets)
	t.Run("UserToCatalogs", testUserToManySetOpCatalogs)
	t.Run("UserToFileAssets", testUserToManySetOpFileAssets)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("ContentTypeToContainers", testContentTypeToManyRemoveOpContainers)
	t.Run("CatalogToParentCatalogs", testCatalogToManyRemoveOpParentCatalogs)
	t.Run("ServerToServernameFileAssets", testServerToManyRemoveOpServernameFileAssets)
	t.Run("VirtualLessonToContainers", testVirtualLessonToManyRemoveOpContainers)
	t.Run("LanguageToLangCatalogDescriptions", testLanguageToManyRemoveOpLangCatalogDescriptions)
	t.Run("LanguageToLangContainerDescriptions", testLanguageToManyRemoveOpLangContainerDescriptions)
	t.Run("LanguageToLangContainers", testLanguageToManyRemoveOpLangContainers)
	t.Run("LanguageToLangFileAssets", testLanguageToManyRemoveOpLangFileAssets)
	t.Run("UserToCatalogs", testUserToManyRemoveOpCatalogs)
	t.Run("UserToFileAssets", testUserToManyRemoveOpFileAssets)
}

func TestReload(t *testing.T) {
	t.Run("FileTypes", testFileTypesReload)
	t.Run("LecturerDescriptions", testLecturerDescriptionsReload)
	t.Run("Lecturers", testLecturersReload)
	t.Run("ContentTypes", testContentTypesReload)
	t.Run("Roles", testRolesReload)
	t.Run("CatalogDescriptions", testCatalogDescriptionsReload)
	t.Run("Catalogs", testCatalogsReload)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsReload)
	t.Run("ContainerDescriptions", testContainerDescriptionsReload)
	t.Run("Labels", testLabelsReload)
	t.Run("RolesUsers", testRolesUsersReload)
	t.Run("Containers", testContainersReload)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsReload)
	t.Run("Servers", testServersReload)
	t.Run("VirtualLessons", testVirtualLessonsReload)
	t.Run("Languages", testLanguagesReload)
	t.Run("FileAssets", testFileAssetsReload)
	t.Run("Users", testUsersReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("FileTypes", testFileTypesReloadAll)
	t.Run("LecturerDescriptions", testLecturerDescriptionsReloadAll)
	t.Run("Lecturers", testLecturersReloadAll)
	t.Run("ContentTypes", testContentTypesReloadAll)
	t.Run("Roles", testRolesReloadAll)
	t.Run("CatalogDescriptions", testCatalogDescriptionsReloadAll)
	t.Run("Catalogs", testCatalogsReloadAll)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsReloadAll)
	t.Run("ContainerDescriptions", testContainerDescriptionsReloadAll)
	t.Run("Labels", testLabelsReloadAll)
	t.Run("RolesUsers", testRolesUsersReloadAll)
	t.Run("Containers", testContainersReloadAll)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsReloadAll)
	t.Run("Servers", testServersReloadAll)
	t.Run("VirtualLessons", testVirtualLessonsReloadAll)
	t.Run("Languages", testLanguagesReloadAll)
	t.Run("FileAssets", testFileAssetsReloadAll)
	t.Run("Users", testUsersReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("FileTypes", testFileTypesSelect)
	t.Run("LecturerDescriptions", testLecturerDescriptionsSelect)
	t.Run("Lecturers", testLecturersSelect)
	t.Run("ContentTypes", testContentTypesSelect)
	t.Run("Roles", testRolesSelect)
	t.Run("CatalogDescriptions", testCatalogDescriptionsSelect)
	t.Run("Catalogs", testCatalogsSelect)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsSelect)
	t.Run("ContainerDescriptions", testContainerDescriptionsSelect)
	t.Run("Labels", testLabelsSelect)
	t.Run("RolesUsers", testRolesUsersSelect)
	t.Run("Containers", testContainersSelect)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsSelect)
	t.Run("Servers", testServersSelect)
	t.Run("VirtualLessons", testVirtualLessonsSelect)
	t.Run("Languages", testLanguagesSelect)
	t.Run("FileAssets", testFileAssetsSelect)
	t.Run("Users", testUsersSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("FileTypes", testFileTypesUpdate)
	t.Run("LecturerDescriptions", testLecturerDescriptionsUpdate)
	t.Run("Lecturers", testLecturersUpdate)
	t.Run("ContentTypes", testContentTypesUpdate)
	t.Run("Roles", testRolesUpdate)
	t.Run("CatalogDescriptions", testCatalogDescriptionsUpdate)
	t.Run("Catalogs", testCatalogsUpdate)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsUpdate)
	t.Run("ContainerDescriptions", testContainerDescriptionsUpdate)
	t.Run("Labels", testLabelsUpdate)
	t.Run("RolesUsers", testRolesUsersUpdate)
	t.Run("Containers", testContainersUpdate)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsUpdate)
	t.Run("Servers", testServersUpdate)
	t.Run("VirtualLessons", testVirtualLessonsUpdate)
	t.Run("Languages", testLanguagesUpdate)
	t.Run("FileAssets", testFileAssetsUpdate)
	t.Run("Users", testUsersUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("FileTypes", testFileTypesSliceUpdateAll)
	t.Run("LecturerDescriptions", testLecturerDescriptionsSliceUpdateAll)
	t.Run("Lecturers", testLecturersSliceUpdateAll)
	t.Run("ContentTypes", testContentTypesSliceUpdateAll)
	t.Run("Roles", testRolesSliceUpdateAll)
	t.Run("CatalogDescriptions", testCatalogDescriptionsSliceUpdateAll)
	t.Run("Catalogs", testCatalogsSliceUpdateAll)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsSliceUpdateAll)
	t.Run("ContainerDescriptions", testContainerDescriptionsSliceUpdateAll)
	t.Run("Labels", testLabelsSliceUpdateAll)
	t.Run("RolesUsers", testRolesUsersSliceUpdateAll)
	t.Run("Containers", testContainersSliceUpdateAll)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsSliceUpdateAll)
	t.Run("Servers", testServersSliceUpdateAll)
	t.Run("VirtualLessons", testVirtualLessonsSliceUpdateAll)
	t.Run("Languages", testLanguagesSliceUpdateAll)
	t.Run("FileAssets", testFileAssetsSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
}

func TestUpsert(t *testing.T) {
	t.Run("FileTypes", testFileTypesUpsert)
	t.Run("LecturerDescriptions", testLecturerDescriptionsUpsert)
	t.Run("Lecturers", testLecturersUpsert)
	t.Run("ContentTypes", testContentTypesUpsert)
	t.Run("Roles", testRolesUpsert)
	t.Run("CatalogDescriptions", testCatalogDescriptionsUpsert)
	t.Run("Catalogs", testCatalogsUpsert)
	t.Run("ContainerDescriptionPatterns", testContainerDescriptionPatternsUpsert)
	t.Run("ContainerDescriptions", testContainerDescriptionsUpsert)
	t.Run("Labels", testLabelsUpsert)
	t.Run("RolesUsers", testRolesUsersUpsert)
	t.Run("Containers", testContainersUpsert)
	t.Run("FileAssetDescriptions", testFileAssetDescriptionsUpsert)
	t.Run("Servers", testServersUpsert)
	t.Run("VirtualLessons", testVirtualLessonsUpsert)
	t.Run("Languages", testLanguagesUpsert)
	t.Run("FileAssets", testFileAssetsUpsert)
	t.Run("Users", testUsersUpsert)
}
