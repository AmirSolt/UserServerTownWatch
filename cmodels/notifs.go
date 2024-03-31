package cmodels

import (
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

const notifs string = "notifs"

var _ models.Model = (*Customer)(nil)

type Notif struct {
	models.BaseModel
	User     string `db:"user" json:"user"`
	Subject  string `db:"subject" json:"subject"`
	BodyHTML string `db:"body_html" json:"body_html"`
}

func (m *Notif) TableName() string {
	return notifs // the name of your collection
}

// =======================================

// func (c *Customer) Find(app core.App, ctx echo.Context, skipNoRowsErr bool) *base.CError {
// 	return FindOne(app, ctx, c, "")
// }

// func (c *Customer) Save(app core.App, ctx echo.Context) *base.CError {

// }

// =======================================

func createNotifsCollection(app core.App) {

	collectionName := notifs

	existingCollection, _ := app.Dao().FindCollectionByNameOrId(collectionName)
	if existingCollection != nil {
		return
	}

	users, err := app.Dao().FindCollectionByNameOrId(users)
	if err != nil {
		log.Fatalf("users table not found: %+v", err)
	}

	collection := &models.Collection{
		Name:       collectionName,
		Type:       models.CollectionTypeBase,
		ListRule:   nil,
		ViewRule:   nil,
		CreateRule: nil,
		UpdateRule: nil,
		DeleteRule: nil,
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "user",
				Type:     schema.FieldTypeRelation,
				Required: true,
				Options: &schema.RelationOptions{
					MaxSelect:     types.Pointer(1),
					CollectionId:  users.Id,
					CascadeDelete: true,
				},
			},
			&schema.SchemaField{
				Name:     "subject",
				Type:     schema.FieldTypeText,
				Required: true,
				Options:  &schema.TextOptions{},
			},
			&schema.SchemaField{
				Name:     "body_html",
				Type:     schema.FieldTypeText,
				Required: true,
				Options:  &schema.TextOptions{},
			},
		),
		Indexes: types.JsonArray[string]{
			fmt.Sprintf("CREATE UNIQUE INDEX idx_user ON %s (user)", collectionName),
		},
	}

	if err := app.Dao().SaveCollection(collection); err != nil {
		log.Fatalf("%s collection failed: %+v", collectionName, err)
	}
}
