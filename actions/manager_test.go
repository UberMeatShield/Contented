package actions

import (
	//"fmt"
	"contented/utils"
    "net/http"
    /*
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite"
	"github.com/gofrs/uuid"
	"github.com/gobuffalo/envy"
    "github.com/gobuffalo/buffalo"
    */
    "context"
    //"sync"
    //"github.com/gobuffalo/logger"
    "github.com/gobuffalo/envy"
    "github.com/gobuffalo/buffalo"
)

var expect_len = map[string]int{
    "dir1": 12,
    "dir2": 2,
    "dir3": 6,
    "screens": 4,
}

func basicContext() buffalo.DefaultContext {
	return buffalo.DefaultContext{
		Context: context.Background(),
		//logger:  buffalo.logger.New(logger.DebugLevel),
		//data:    &sync.Map{},
		//flash:   &Flash{data: make(map[string][]string)},
	}
}

func (as *ActionSuite) Test_ManagerContainers() {
    init_fake_app(false)
    ctx := getContext(app)
    man := GetManager(&ctx)
    containers, err := man.ListContainersContext()
    as.NoError(err)

    for _, c := range *containers {
        c_mem, err := man.FindDirRef(c.ID)
        if err != nil {
            as.Fail("It should not have an issue finding valid containers")
        }
        as.Equal(c_mem.ID, c.ID)
    }
}


func (as *ActionSuite) Test_ManagerMediaContainer() {
    init_fake_app(false)
    ctx := getContext(app)
    man := GetManager(&ctx)
    mcs, err := man.ListAllMedia(1, 9001)
    as.NoError(err)

    for _, mc := range *mcs {
        cm, err := man.FindFileRef(mc.ID)
        if err != nil {
            as.Fail("It should not have an issue finding valid containers")
        }
        as.Equal(cm.ID, mc.ID)
    }
}

func (as *ActionSuite) Test_AssignManager() {
    dir, _ := envy.MustGet("DIR")
    cfg := GetCfg()
    cfg.UseDatabase = false
    utils.InitConfig(dir, cfg)

    mem := ContentManagerMemory{}
    mem.validate = "Memory"
    mem.SetCfg(cfg)
    mem.Initialize()

    memCfg := mem.GetCfg()
    as.NotNil(memCfg, "It should be defined")
    mcs, err := mem.ListAllMedia(1, 9001)
    as.NoError(err)
    as.Greater(len(*mcs), 0, "It should have valid files in the manager")

    appCfg.UseDatabase = false
    ctx := getContext(app)
    man := GetManager(&ctx)  // New Reference but should have the same count of media
    mcs_2, _ := man.ListAllMedia(1, 9000)

    as.Equal(len(*mcs), len(*mcs_2), "A new instance should use the same storage")
}


func (as *ActionSuite) Test_MemoryDenyEdit() {
    cfg := init_fake_app(false)
    cfg.UseDatabase = false
    ctx := getContext(app)
    man := GetManager(&ctx)

    containers, err := man.ListContainersContext()
    as.NoError(err, "It should list containers")

    as.Greater(len(*containers), 0, "There should be containers")

    for _, c := range *containers {
        c.Name = "Update Should fail"
        res := as.JSON("/containers/" + c.ID.String()).Put(&c)
        as.Equal(http.StatusNotImplemented, res.Code)
    }
}

func (as *ActionSuite) Test_MemoryManagerPaginate() {
    cfg := init_fake_app(false)
    cfg.UseDatabase = false

    ctx := getContextParams(app, "/containers", "1", "2")
    man := GetManager(&ctx)
    as.Equal(man.CanEdit(), false, "Memory manager should not be editing")

    // Hate
    containers, err := man.ListContainers(1, 1)
    as.NoError(err, "It should list with pagination")
    as.Equal(1, len(*containers), "It should respect paging")

    cnt := (*containers)[0]
    as.NotNil(cnt, "There should be a container with 12 entries")
    as.Equal(cnt.Total, 12, "There should be 12 test images in the first ORDERED containers")
    as.NoError(err)
    media_page_1, _ := man.ListMedia(cnt.ID, 1, 4)
    as.Equal(len(*media_page_1), 4, "It should respect page size")

    media_page_3, _ := man.ListMedia(cnt.ID, 3, 4)
    as.Equal(len(*media_page_3), 4, "It should respect page size and get the last page")

    as.NotEqual((*media_page_3)[3].ID, (*media_page_1)[3].ID, "Ensure it actually paged")

    // Last container pagination check
    l_cnts, _ := man.ListContainers(4, 1)
    as.Equal(1, len(*l_cnts), "It should still return only as we are on the last page")
    l_cnt := (*l_cnts)[0]
    as.Equal(l_cnt.Total, expect_len[l_cnt.Name], "There are 3 entries in the ordered test data last container")

}

func (as *ActionSuite) Test_ManagerInitialize() {
    cfg := init_fake_app(false)
    cfg.UseDatabase = false

    ctx := getContext(app)
    man := GetManager(&ctx)
    as.NotNil(man, "It should have a manager defined after init")

    containers, err := man.ListContainersContext()
    as.NoError(err, "It should list all containers")
    as.NotNil(containers, "It should have containers")
    as.Equal(len(*containers), 4, "It should have 4 of them")

    // Memory test working
    for _, c := range *containers {
        // fmt.Printf("Searching for this container %s with name %s\n", c.ID, c.Name)
        media, err := man.ListMediaContext(c.ID)
        as.NoError(err)
        as.NotNil(media)

        media_len := len(*media)
        // fmt.Printf("Media length was %d\n", media_len)
        as.Greater(media_len, 0, "There should be a number of media")
        as.Equal(expect_len[c.Name], media_len, "It should have this many instances: " + c.Name )
        as.Greater(c.Total, 0, "All of them should have a total assigned")
    }
}
