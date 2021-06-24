package grifts

import (
    "fmt"
    "net/url"
    "contented/models"
    "contented/actions"
    "contented/utils"
    "github.com/gobuffalo/pop/v5"
	"github.com/markbates/grift/grift"
    "github.com/pkg/errors"
    "github.com/gobuffalo/envy"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Populate the DB with a set of directory content.")
	grift.Add("seed", func(c *grift.Context) error {
        cfg := utils.GetCfg()
        utils.InitConfigEnvy(cfg)

        // Require the directory which we want to process (maybe just trust it is set)
        _, d_err := envy.MustGet("DIR")
        if d_err != nil {
            return errors.WithStack(d_err)
        }
        // Clean out the current DB setup then builds a new one
        fmt.Printf("Configuration is loaded %s Starting import", cfg.Dir)
        return actions.CreateInitialStructure(cfg)
	})
    // Then add the content for the entire directory structure

	grift.Add("preview", func(c *grift.Context) error {
        cfg := utils.GetCfg()
        utils.InitConfigEnvy(cfg)
        fmt.Printf("Configuration is loaded %s doing preview creation", cfg.Dir)

        get_params := func() *url.Values {
            vals := url.Values{}  // TODO: Maybe set this via something sensible
            return &vals
        }
        if cfg.UseDatabase {
            // The scope of transactions is a bit odd.  Seems like this is handled in
            // buffalo via the magical buffalo middleware.
            return models.DB.Transaction(func(tx *pop.Connection) error {
                get_connection := func() *pop.Connection {
                    return tx
                } 
                man := actions.CreateManager(cfg, get_connection, get_params)
                cnts, c_err := man.ListContainers(1, 90001)
                if c_err != nil {
                    fmt.Printf("Error loading containers %s", c_err)
                }
                fmt.Printf("Manager containers (%d)\n", len(*cnts))
                return c_err
            })
        } else {
            get_connection := func() *pop.Connection {
                return nil // Do not do anything with the DB
            }
            man := actions.CreateManager(cfg, get_connection, get_params)
            fmt.Printf("Use memory manager %s", man.CanEdit())
        }
        return nil
        //return actions.CreateAllPreviews(cfg.PreviewOverSize)
        /*
        man := actions.GetManager(c)
        cnts, err := man.ListContainers(0, 90000)
        if err != nil {
            fmt.Printf("Error listing containers %s", err)
        } else {
            fmt.Printf("Listed all containers, wooo %d", len(*cnts))
        }
        */
	})
})
