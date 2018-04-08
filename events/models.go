package events

import (
	"github.com/Bnei-Baruch/mdb/models"
)

type Event struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func makeEvent(Type string, Payload map[string]interface{}) Event {
	return Event{Type: Type, Payload: Payload}
}

func CollectionCreateEvent(c *models.Collection) Event {
	return makeEvent(E_COLLECTION_CREATE, map[string]interface{}{
		"id":  c.ID,
		"uid": c.UID,
		//"type_id": c.TypeID,
	})
}

func CollectionUpdateEvent(c *models.Collection) Event {
	return makeEvent(E_COLLECTION_UPDATE, map[string]interface{}{
		"id":  c.ID,
		"uid": c.UID,
	})
}

func CollectionDeleteEvent(c *models.Collection) Event {
	return makeEvent(E_COLLECTION_DELETE, map[string]interface{}{
		"id":  c.ID,
		"uid": c.UID,
	})
}

func CollectionPublishedChangeEvent(c *models.Collection) Event {
	return makeEvent(E_COLLECTION_PUBLISHED_CHANGE, map[string]interface{}{
		"id":  c.ID,
		"uid": c.UID,
		//"published": c.Published,
	})
}

func CollectionContentUnitsChangeEvent(c *models.Collection) Event {
	return makeEvent(E_COLLECTION_CONTENT_UNITS_CHANGE, map[string]interface{}{
		"id":  c.ID,
		"uid": c.UID,
	})
}

func ContentUnitCreateEvent(cu *models.ContentUnit) Event {
	return makeEvent(E_CONTENT_UNIT_CREATE, map[string]interface{}{
		"id":  cu.ID,
		"uid": cu.UID,
		//"type_id": cu.TypeID,
	})
}

func ContentUnitUpdateEvent(cu *models.ContentUnit) Event {
	return makeEvent(E_CONTENT_UNIT_UPDATE, map[string]interface{}{
		"id":  cu.ID,
		"uid": cu.UID,
	})
}

func ContentUnitDeleteEvent(cu *models.ContentUnit) Event {
	return makeEvent(E_CONTENT_UNIT_DELETE, map[string]interface{}{
		"id":  cu.ID,
		"uid": cu.UID,
	})
}

func ContentUnitPublishedChangeEvent(cu *models.ContentUnit) Event {
	return makeEvent(E_CONTENT_UNIT_PUBLISHED_CHANGE, map[string]interface{}{
		"id":  cu.ID,
		"uid": cu.UID,
		//"published": cu.Published,
	})
}

func ContentUnitDerivativesChangeEvent(cu *models.ContentUnit) Event {
	return makeEvent(E_CONTENT_UNIT_DERIVATIVES_CHANGE, map[string]interface{}{
		"id":  cu.ID,
		"uid": cu.UID,
	})
}

func ContentUnitSourcesChangeEvent(cu *models.ContentUnit) Event {
	return makeEvent(E_CONTENT_UNIT_SOURCES_CHANGE, map[string]interface{}{
		"id":  cu.ID,
		"uid": cu.UID,
	})
}

func ContentUnitTagsChangeEvent(cu *models.ContentUnit) Event {
	return makeEvent(E_CONTENT_UNIT_TAGS_CHANGE, map[string]interface{}{
		"id":  cu.ID,
		"uid": cu.UID,
	})
}

func ContentUnitPersonsChangeEvent(cu *models.ContentUnit) Event {
	return makeEvent(E_CONTENT_UNIT_PERSONS_CHANGE, map[string]interface{}{
		"id":  cu.ID,
		"uid": cu.UID,
	})
}

func ContentUnitPublishersChangeEvent(cu *models.ContentUnit) Event {
	return makeEvent(E_CONTENT_UNIT_PUBLISHERS_CHANGE, map[string]interface{}{
		"id":  cu.ID,
		"uid": cu.UID,
	})
}

func FileUpdateEvent(f *models.File) Event {
	return makeEvent(E_FILE_UPDATE, map[string]interface{}{
		"id":  f.ID,
		"uid": f.UID,
		//"name": f.Name,
		//"size": f.Size,
		//"sha1": hex.EncodeToString(f.Sha1.Bytes),
	})
}

func FileInsertEvent(f *models.File, insertType string) Event {
	return makeEvent(E_FILE_INSERT, map[string]interface{}{
		"id":  f.ID,
		"uid": f.UID,
		//"name": f.Name,
		//"size": f.Size,
		//"sha1": hex.EncodeToString(f.Sha1.Bytes),
		"insert_type": insertType,
	})
}

func FileReplaceEvent(oldFile *models.File, newFile *models.File, insertType string) Event {
	return makeEvent(E_FILE_REPLACE, map[string]interface{}{
		"old": map[string]interface{}{
			"id":  oldFile.ID,
			"uid": oldFile.UID,
			//"name": oldFile.Name,
			//"size": oldFile.Size,
			//"sha1": hex.EncodeToString(oldFile.Sha1.Bytes),
		},
		"new": map[string]interface{}{
			"id":  newFile.ID,
			"uid": newFile.UID,
			//"name": newFile.Name,
			//"size": newFile.Size,
			//"sha1": hex.EncodeToString(oldFile.Sha1.Bytes),
		},
		"insert_type": insertType,
	})
}

func FilePublishedEvent(f *models.File) Event {
	return makeEvent(E_FILE_PUBLISHED, map[string]interface{}{
		"id":  f.ID,
		"uid": f.UID,
		//"name": f.Name,
		//"size": f.Size,
		//"sha1": hex.EncodeToString(f.Sha1.Bytes),
	})
}

func FileRemoveEvent(f *models.File) Event {
	return makeEvent(E_FILE_REMOVE, map[string]interface{}{
		"id":  f.ID,
		"uid": f.UID,
		//"name": f.Name,
		//"size": f.Size,
		//"sha1": hex.EncodeToString(f.Sha1.Bytes),
	})
}

func SourceCreateEvent(s *models.Source) Event {
	return makeEvent(E_SOURCE_CREATE, map[string]interface{}{
		"id":  s.ID,
		"uid": s.UID,
	})
}

func SourceUpdateEvent(s *models.Source) Event {
	return makeEvent(E_SOURCE_UPDATE, map[string]interface{}{
		"id":  s.ID,
		"uid": s.UID,
	})
}

func TagCreateEvent(t *models.Tag) Event {
	return makeEvent(E_TAG_CREATE, map[string]interface{}{
		"id":  t.ID,
		"uid": t.UID,
	})
}

func TagUpdateEvent(t *models.Tag) Event {
	return makeEvent(E_TAG_UPDATE, map[string]interface{}{
		"id":  t.ID,
		"uid": t.UID,
	})
}

func PersonCreateEvent(p *models.Person) Event {
	return makeEvent(E_PERSON_CREATE, map[string]interface{}{
		"id":  p.ID,
		"uid": p.UID,
	})
}

func PersonUpdateEvent(p *models.Person) Event {
	return makeEvent(E_PERSON_UPDATE, map[string]interface{}{
		"id":  p.ID,
		"uid": p.UID,
	})
}

func PersonDeleteEvent(p *models.Person) Event {
	return makeEvent(E_PERSON_DELETE, map[string]interface{}{
		"id":  p.ID,
		"uid": p.UID,
	})
}

func PublisherCreateEvent(p *models.Publisher) Event {
	return makeEvent(E_PUBLISHER_CREATE, map[string]interface{}{
		"id":  p.ID,
		"uid": p.UID,
	})
}

func PublisherUpdateEvent(p *models.Publisher) Event {
	return makeEvent(E_PUBLISHER_UPDATE, map[string]interface{}{
		"id":  p.ID,
		"uid": p.UID,
	})
}
