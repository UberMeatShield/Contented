package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
    "github.com/gofrs/uuid"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite"
	"contented/utils"
	"contented/models"
)

type ActionSuite struct {
	*suite.Action
}

func TestMain(m *testing.M) {
	dir, err := envy.MustGet("DIR")
	cfg.Dir = dir // Absolute path to our known test data
	cfg.ValidDirs = utils.GetDirectoriesLookup(cfg.Dir)

	if err != nil {
		fmt.Println("DIR ENV REQUIRED$ export=DIR=`pwd`/mocks/content/ && buffalo test")
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

func (as *ActionSuite) Test_ContentList() {
	res := as.JSON("/content/").Get()
	as.Equal(http.StatusOK, res.Code)

	resObj := PreviewResults{}
	json.NewDecoder(res.Body).Decode(&resObj)
	as.Equal(4, len(resObj.Results), "We should have this many directories present")

	// Check all the directories have contents
	for _, dir := range resObj.Results {
		as.Greater(dir.Total, 0, "There should be multiple results")
		as.NotNil(dir.Id, "There should be a directory ID")
		as.Greater(len(dir.Contents), 0, "And all our test mocks have content")
	}
}

func (as *ActionSuite) Test_ContentDirLoad() {
    lookup := utils.GetDirectoriesLookup(cfg.Dir)

    as.Equal(len(lookup), 4, "There should be 4 test directories")

    for id, f := range lookup {
        res := as.JSON("/content/" + id).Get()
        as.Equal(http.StatusOK, res.Code)

        resObj := utils.DirContents{}
        json.NewDecoder(res.Body).Decode(&resObj)

        if f.Name() == "dir1" {
            as.Equal(resObj.Total, 12, "It should have a known number of images")
        }
    }
}

func (as *ActionSuite) Test_ViewRef() {
    // Oof, that is rough... need a better way to select the file not by index but ID
    lookup := utils.GetDirectoriesLookup(cfg.Dir)

    // TODO: Make it better about the type checking
    // TODO: Make it always pass in the file ID
    for id, _ := range lookup {
        res := as.HTML("/view/" + id + "/0").Get()
        as.Equal(http.StatusOK, res.Code)
        header := res.Header()

        as.Equal("image/png", header.Get("Content-Type"))
    }
}

func (as *ActionSuite) Test_ContentDirDownload() {

    // Oof, that is rough... need a better way to select the file not by index but ID
    lookup := utils.GetDirectoriesLookup(cfg.Dir)

    valid := map[string]bool{"image/png": true, "image/jpeg": true}
    
    for id, f := range lookup {
        res := as.HTML("/download/" + id + "/0").Get()
        as.Equal(http.StatusOK, res.Code)
        if f.Name() == "dir3" {
            header := res.Header()
	        as.Contains(valid, header.Get("Content-Type"))
        }
    }

    // Make it so we know that there is only png in one dir, but iterate over all content

	dir_id_url := "/download/9d553cdef482947b97b5beda2dc594c7c818a69a49e04f044f4505bc223a3535/1"
	res1 := as.HTML(dir_id_url).Get()
	as.Equal(http.StatusOK, res1.Code)
	header1 := res1.Header()
	as.Equal("image/png", header1.Get("Content-Type"))
}

func (as *ActionSuite) Test_FindAndCache() {
    mc1_id, _ := uuid.NewV4()
    mc2_id, _ := uuid.NewV4()
    mc1 := models.MediaContainer{ID: mc1_id, Src: "mc1"}
    mc2 := models.MediaContainer{ID: mc2_id, Src: "mc2"}

    CacheFile(mc1)
    CacheFile(mc2)

    as.Equal(len(cfg.ValidFiles), 2)

    name1, _ := FindFileRef(mc1_id)
    as.Equal(name1.Src, mc1.Src)
    name2, _ := FindFileRef(mc2_id)
    as.Equal(name2.Src, mc2.Src)

    invalid, _ := uuid.NewV4()
    _, err := FindFileRef(invalid)
    as.Error(err)
}


func Test_ActionSuite(t *testing.T) {
	action, err := suite.NewActionWithFixtures(App(), packr.New("Test_ActionSuite", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}
