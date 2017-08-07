package kmedia

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

type AtomicInt32 struct {
	value int32
}

func (a *AtomicInt32) Inc(delta int32) {
	atomic.AddInt32(&a.value, delta)
}

func (a *AtomicInt32) Get() int32 {
	return atomic.LoadInt32(&a.value)
}

type AtomicHistogram struct {
	value        map[string]int
	sync.RWMutex // Read Write mutex, guards access to internal map
}

func NewAtomicHistogram() *AtomicHistogram {
	return &AtomicHistogram{value: make(map[string]int)}
}

func (h *AtomicHistogram) Inc(key string, delta int) {
	h.Lock()
	h.value[key] += delta
	h.Unlock()
}

func (h *AtomicHistogram) Get() map[string]int {
	h.RLock()
	r := make(map[string]int, len(h.value))
	for k, v := range h.value {
		r[k] = v
	}
	h.RUnlock()
	return r
}

type ImportStatistics struct {
	LessonsProcessed       AtomicInt32
	ValidLessons           AtomicInt32
	InvalidLessons         AtomicInt32
	CatalogsProcessed       AtomicInt32
	ContainersProcessed    AtomicInt32
	ContainersVisited      AtomicInt32
	ContainersWithFiles    AtomicInt32
	ContainersWithoutFiles AtomicInt32
	FileAssetsProcessed    AtomicInt32
	FileAssetsMissingSHA1  AtomicInt32
	FileAssetsWInvalidMT   AtomicInt32
	FileAssetsMissingType  AtomicInt32

	OperationsCreated   AtomicInt32
	CollectionsCreated  AtomicInt32
	CollectionsUpdated  AtomicInt32
	ContentUnitsCreated AtomicInt32
	ContentUnitsUpdated AtomicInt32
	FilesCreated        AtomicInt32
	FilesUpdated        AtomicInt32

	TxCommitted  AtomicInt32
	TxRolledBack AtomicInt32

	UnkownCatalogs AtomicHistogram
}

func NewImportStatistics() *ImportStatistics {
	return &ImportStatistics{UnkownCatalogs: *NewAtomicHistogram()}
}

func (s *ImportStatistics) dump() {
	fmt.Println("Here comes import statistics:")

	fmt.Println("Kmedia:")
	fmt.Printf("ValidLessons            		%d\n", s.ValidLessons.Get())
	fmt.Printf("InvalidLessons          		%d\n", s.InvalidLessons.Get())
	fmt.Printf("LessonsProcessed        		%d\n", s.LessonsProcessed.Get())
	fmt.Printf("CatalogsProcessed        		%d\n", s.CatalogsProcessed.Get())
	fmt.Printf("ContainersVisited       		%d\n", s.ContainersVisited.Get())
	fmt.Printf("ContainersWithFiles     		%d\n", s.ContainersWithFiles.Get())
	fmt.Printf("ContainersWithoutFiles  		%d\n", s.ContainersWithoutFiles.Get())
	fmt.Printf("ContainersProcessed     		%d\n", s.ContainersProcessed.Get())
	fmt.Printf("FileAssetsProcessed     		%d\n", s.FileAssetsProcessed.Get())
	fmt.Printf("FileAssetsMissingSHA1   		%d\n", s.FileAssetsMissingSHA1.Get())
	fmt.Printf("FileAssetsWInvalidMT    		%d\n", s.FileAssetsWInvalidMT.Get())
	fmt.Printf("FileAssetsMissingType   		%d\n", s.FileAssetsMissingType.Get())

	fmt.Println("MDB:")
	fmt.Printf("OperationsCreated       		%d\n", s.OperationsCreated.Get())
	fmt.Printf("CollectionsCreated      		%d\n", s.CollectionsCreated.Get())
	fmt.Printf("CollectionsUpdated      		%d\n", s.CollectionsUpdated.Get())
	fmt.Printf("ContentUnitsCreated     		%d\n", s.ContentUnitsCreated.Get())
	fmt.Printf("ContentUnitsUpdated     		%d\n", s.ContentUnitsUpdated.Get())
	fmt.Printf("FilesCreated            		%d\n", s.FilesCreated.Get())
	fmt.Printf("FilesUpdated            		%d\n", s.FilesUpdated.Get())

	fmt.Printf("TxCommitted             		%d\n", s.TxCommitted.Get())
	fmt.Printf("TxRolledBack            		%d\n", s.TxRolledBack.Get())

	uc := s.UnkownCatalogs.Get()
	fmt.Printf("UnkownCatalogs            		%d\n", len(uc))
	keys := make([]string, len(uc))
	i := 0
	for k := range uc {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s\t%d\n", k, uc[k])
	}
}
