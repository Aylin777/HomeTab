package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/gotag/model"
	"gitlab.com/systemz/gotag/model2"
	"log"
	"strconv"
)

func init() {
	rootCmd.AddCommand(migrate)
}

var migrate = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate to new DB",
	Run:   migrateExec,
}

func migrateExec(cmd *cobra.Command, args []string) {
	// SQLite stuff
	sqlite := model.DbInit()
	allFiles := model.CountAllFiles(sqlite)
	log.Printf("all files: %v", allFiles)

	// MySQL stuff
	mysql := model2.InitMysql()

	// scan all entries in DB
	imgs := model.ListAll(sqlite)
	for _, img := range imgs {

		log.Printf("%v %v %v", img.Fid, img.Name, img.Sha256)

		//time.Sleep(time.Millisecond * 50)
		//continue

		// upgrade pHash storage
		pHashA := 0
		pHashB := 0
		pHashC := 0
		pHashD := 0
		if len(img.Phash) > 1 {
			pHashA, _ = strconv.Atoi(img.Phash[0:16])
			pHashB, _ = strconv.Atoi(img.Phash[16:32])
			pHashC, _ = strconv.Atoi(img.Phash[32:48])
			pHashD, _ = strconv.Atoi(img.Phash[48:64])
		}

		// save file to DB
		file := &model2.File{
			Filename: img.Name,
			FilePath: img.Path,
			SizeB:    img.Size,
			// mime
			Sha256: img.Sha256,
			PhashA: pHashA,
			PhashB: pHashB,
			PhashC: pHashC,
			PhashD: pHashD,
		}
		mysql.Save(&file)

		//TODO mime
		// SELECT  HAMMINGDISTANCE(a1,a2,a3,a4,b1,b2,b3,b4) AS res FROM `files` WHERE sha256 = "changeme"

		// add tags to DB
		found, tags := model.TagList(sqlite, img.Fid)
		if !found {
			// finish work if no tags for this file
			continue
		}
		for _, tag := range tags {
			model2.AddTagToFile(mysql, tag.Name, file.Id)
		}
	}
}
