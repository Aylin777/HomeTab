package model2

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Tag struct {
	Id  int    `json:"id" gorm:"AUTO_INCREMENT" json:"id"`
	Tag string `gorm:"column:tag"`

	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func AddTagToFile(db *gorm.DB, tag string, fileId int) {
	// get tag if exists already
	// FIXME use transactions
	var tagInDb Tag
	db.Where("tag = ?", tag).First(&tagInDb)

	// tag is not present in DB, create new
	if tagInDb.Id < 1 {
		tagInDb = Tag{Tag: tag}
		db.Create(&tagInDb)
	}

	// tag is now present in DB
	//log.Printf("Tag in DB: %+v", tagInDb)

	// check if link file <-> tag is present too
	var fileTagInDb FileTag
	db.Where("tag_id = ? AND file_id = ?", tagInDb.Id, fileId).First(&fileTagInDb)

	// link is not in DB, create new
	//log.Printf("%v", fileTagInDb)
	if fileTagInDb.Id < 1 {
		fileTagInDb = FileTag{
			FileId: fileId,
			TagId:  tagInDb.Id,
		}
		db.Create(&fileTagInDb)
	}

}

/*

### Return all tags for a specific file

SELECT tags.tag
FROM tags
   INNER JOIN file_tags on tags.id = file_tags.tag_id
   INNER JOIN files on file_tags.file_id = files.id
WHERE files.sha256 = "CHANGE_ME"

### Return all files having tag “beer”

SELECT DISTINCT files.id, files.sha256, file_name
FROM files
   INNER JOIN file_tags
      ON file_tags.file_id = files.id
   INNER JOIN tags
      ON tags.id = file_tags.tag_id
   WHERE tags.tag IN ('beer')


### Return top tags (most used)

SELECT COUNT(tags.id) AS rank, tags.*
FROM tags
   LEFT JOIN file_tags ON tags.id = file_tags.tag_id
     GROUP BY tags.id
ORDER BY COUNT(tags.id)
DESC LIMIT 50

 */
