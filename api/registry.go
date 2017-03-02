package api

import "github.com/Bnei-Baruch/mdb/models"

var (
	CONTENT_TYPE_REGISTRY   = &ContentTypeRegistry{}
	OPERATION_TYPE_REGISTRY = &OperationTypeRegistry{}
)

type ContentTypeRegistry struct {
	ByName map[string]*models.ContentType
}

func (r *ContentTypeRegistry) Init() error {
	types, err := models.ContentTypesG().All()
	if err != nil {
		return err
	}

	r.ByName = make(map[string]*models.ContentType)
	for _, t := range types {
		r.ByName[t.Name] = t
	}

	return nil
}

type OperationTypeRegistry struct {
	ByName map[string]*models.OperationType
}

func (r *OperationTypeRegistry) Init() error {
	types, err := models.OperationTypesG().All()
	if err != nil {
		return err
	}

	r.ByName = make(map[string]*models.OperationType)
	for _, t := range types {
		r.ByName[t.Name] = t
	}

	return nil
}
