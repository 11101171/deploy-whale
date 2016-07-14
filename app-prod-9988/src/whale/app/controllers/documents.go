//****************************************************/
// Addres: https://github.com/11101171
// Author: ningzhong.zeng
//****************************************************/

package controllers

import (
	"whale/app/models"
	"whale/app/routes"

	"github.com/revel/revel"
)

type Documents struct {
	Super
}

func (c Documents) List() revel.Result {
	user := c.SessionGetUser()
	documents := models.SelectDocumentListByUserId(user.UserId)
	return c.Render(documents)
}

func (c Documents) Operate(id string, document models.Document) revel.Result {
	if c.IsGet() {
		if id != "" {
			if document = models.SelectDocumentOneByDocumentId(id); document.DocumentId == "" {
				return c.ErrDataBase(MsgSeleteError)
			}
			return c.Render(document)
		}
		return c.Render()
	} else {
		if document.DocumentId != "" {
			document.ValidateUpdate(c.Validation)
			if c.Validation.HasErrors() {
				return c.Render(document)
			}
			if !models.UpdateOneDocument(&document) {
				return c.ErrDataBase(MsgUpdateError)
			}
		} else {
			document.UserId = c.SessionGetUser().UserId
			document.ValidateInsert(c.Validation)
			if c.Validation.HasErrors() {
				return c.Render(document)
			}
			if err := models.DBMap().Insert(&document); err != nil {
				return c.ErrDataBase(MsgInsertError)
			}
		}
		return c.Redirect(routes.Documents.List())
	}
}

func (c Documents) Del(id string) revel.Result {
	if id != "" && models.DeleteOneDocumentByDocumentId(id) {
		return c.Redirect(routes.Documents.List())
	}
	return c.RenderJsonErr()
}
