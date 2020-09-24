package grifts

import (
//    "path/filepath"
    "fmt"
    "contented/models"
    "contented/utils"
	"github.com/markbates/grift/grift"
    "github.com/pkg/errors"
    "github.com/gobuffalo/envy"
    "github.com/gobuffalo/nulls"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Populate the DB with a set of directory content.")
	grift.Add("seed", func(c *grift.Context) error {

        // Clean out the current DB setup
        err := models.DB.TruncateAll()
        if err != nil {
            return errors.WithStack(err)
        }

        // Grab the directory which we want to process
        dir_name, d_err := envy.MustGet("DIR")
        if d_err != nil {
            return errors.WithStack(d_err)
        }

        // Process all the directories and get a valid setup
        dirs := utils.ListDirs(dir_name, 4)
        fmt.Printf("Found %d directories.", len(dirs))

        // TODO: Need to do this in a single transaction
        for _, dir := range dirs {
            fmt.Printf("Adding directory %s\n", dir.Name)
            dirObj := &models.Container{
              Path: dir_name,
              Name: dir.Name,
            }
            models.DB.Create(dirObj)

            for _, fi := range dir.Contents {
              mc := &models.MediaContainer{
                  Src: fi.Src,
                  Type: fi.Type,
                  Preview: "TBD",
                  ContainerID: nulls.NewUUID(dirObj.ID),
              }
              models.DB.Create(mc)
            }
        }

		// Add DB seeding stuff here
		return nil
	})

    // Need to do a stanard lookup

    // Then add the content for the entire directory structure

    // Then add in linkage to the related models.

})
