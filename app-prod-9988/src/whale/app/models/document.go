//****************************************************/
// Addres: https://github.com/11101171
// Author: ningzhong.zeng
//****************************************************/

package models

import (
	"strings"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

type Document struct {
	Base
	DocumentId    string
	UserId        string
	Info          string
	HtmlContent   string
	MkdownContent string
}

func (m *Document) Validate(v *revel.Validation) {
	m.HtmlContent = strings.TrimSpace(m.HtmlContent)
	m.MkdownContent = strings.TrimSpace(m.MkdownContent)
	m.UserId = strings.TrimSpace(m.UserId)

	v.Check(m.HtmlContent,
		revel.Required{},
		revel.MinSize{Min: 1},
		revel.MaxSize{Max: 30000},
	).Key("document.HtmlContent").Message("网页内容不能为空")
	v.Check(m.MkdownContent,
		revel.MinSize{Min: 1},
		revel.MaxSize{Max: 30000},
	).Key("document.MkdownContent").Message("mkdown内容不能为空")
	v.Check(m.Info,
		revel.MaxSize{Max: 200},
	).Key("document.Info").Message("Info不能超过200个")
}

func (m *Document) ValidateInsert(v *revel.Validation) {
	m.Validate(v)
}

func (m *Document) ValidateUpdate(v *revel.Validation) {
	m.Validate(v)
	v.Check(m.DocumentId,
		revel.Required{},
	).Key("document.DocumentId").Message("更新需要DocumentId")
}

func (m *Document) PreInsert(_ gorp.SqlExecutor) error {
	m.DocumentId = RandomString("document", 6)
	return m.PreI()
}

func SelectDocumentListByUserId(userId string) (documents []Document) {
	CheckErrSQLSelectList(dbmap.Select(
		&documents,
		"select * from o_document where UserId=:UserId order by GmtCreate desc",
		map[string]string{
			"UserId": userId,
		},
	))
	return documents
}

func SelectDocumentOneByDocumentId(documentId string) (document Document) {
	CheckErrSQLSelectOne(dbmap.SelectOne(
		&document,
		"select * from o_document where DocumentId=:DocumentId",
		map[string]string{
			"DocumentId": documentId,
		},
	))
	return document
}

func UpdateOneDocument(document *Document) bool {
	return Exec(
		"update o_document set HtmlContent=:HtmlContent,MkdownContent=:MkdownContent,Info=:Info where DocumentId=:DocumentId",
		map[string]string{
			"DocumentId":    document.DocumentId,
			"HtmlContent":   document.HtmlContent,
			"MkdownContent": document.MkdownContent,
			"Info":          document.Info,
		},
	)
}

func DeleteOneDocumentByDocumentId(documentId string) bool {
	return Exec(
		"delete from o_document where DocumentId=:documentId",
		map[string]string{
			"documentId": documentId,
		},
	)
}
