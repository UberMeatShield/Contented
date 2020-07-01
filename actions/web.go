package actions

import (
    //"io"
    //"fmt"
//    "url"
  "os"
  "strings"
  "strconv"
  "log"
//    "encoding/json"
//    "net/http"
  "net/url"
  "net/http"
//    "io/ioutil"
  "contented/utils"
  "github.com/gobuffalo/buffalo"
)

type PreviewResults struct{
    Success bool `json:"success"`
    Results []utils.DirContents `json:"results"`
}

type HttpError struct{
    Error string `json:"error"`
    Debug string `json:"debug"`
}

type DirConfigEntry struct{
  Dir string
  ValidDirs map[string]bool
  PreviewCount int
  Limit int
}

// HomeHandler is a default handler to serve up
var DefaultLimit int = 10000  // The max limit set by environment variable
var DefaultPreviewCount int = 8
var cfg = DirConfigEntry{
    Dir: "",
    PreviewCount: DefaultPreviewCount,
    Limit: DefaultLimit,
}

func SetupContented(app *buffalo.App, contentDir string, numToPreview int, limit int) {
    if !strings.HasSuffix(contentDir, "/") {
         contentDir = contentDir + "/"
    }
    log.Printf("Setting up the content directory with %s", contentDir)

    cfg.Dir = contentDir
    cfg.ValidDirs = utils.GetDirectoriesLookup(cfg.Dir)
    cfg.PreviewCount = numToPreview
    cfg.Limit = limit

    // TODO: Somehow need to move the dir into App, but first we want to validate the dir...
    app.ServeFiles("/static", http.Dir(cfg.Dir))
}

func ListDefaultHandler(c buffalo.Context) error {
    path, _ := os.Executable()
    log.Printf("Calling into ListDefault run_dir: %s looking at dir: %s", path, cfg.Dir)
    response := PreviewResults{
        Success: true,
        Results: utils.ListDirs(cfg.Dir, cfg.PreviewCount),
    }
    return c.Render(200, r.JSON(response))
}

func DownloadHandler(c buffalo.Context) error {

    dir_to_list := c.Param("dir_to_list")
    filename := c.Param("filename")

    log.Printf("Calling into download handler with filename %s under %s", dir_to_list, filename)
    if cfg.ValidDirs[dir_to_list] {
        fileref := utils.GetFileContents(cfg.Dir + dir_to_list, filename)
        if fileref != nil {
            return c.Render(200, r.Download(c, filename, fileref))
        } 
    } 
    return c.Render(403, r.JSON(invalidDirMsg(dir_to_list, filename)))
}


// Provide a full listing of a specific directory, not just the preview
func ListSpecificHandler(c buffalo.Context) error {
    argument := c.Param("dir_to_list")

    // Pull out the limit and offset queries, provides pagination
    limit := DefaultLimit
    offset := 0

    limit, _ = strconv.Atoi(GetKeyVal(c, "limit", strconv.Itoa(DefaultLimit)))
    if limit <= 0 || limit > DefaultLimit {
        limit = DefaultLimit // Still cannot ask for more than the startup specified
    }
    offset, _ = strconv.Atoi(GetKeyVal(c, "offset", "0"))

    log.Printf("Limit %d with offset %d in dir %s", limit, offset, cfg.Dir)

    // Now actually return the results for a valid directory
    if cfg.ValidDirs[argument] {
        return c.Render(200, r.JSON(getDirectory(cfg.Dir, argument, limit, offset)))
    } 
    return c.Render(403, r.JSON(invalidDirMsg(argument, "")))
}

func GetKeyVal(c buffalo.Context, key string, defaultVal string) string {
  if m, ok := c.Params().(url.Values); ok {
    for k, v := range m {
      if k == key && v != nil {
          return v[0]
      }
    }
  }
  return defaultVal
}

/**
 * Get the response for a single specific directory
 */
func getDirectory(dir string, argument string, limit int, offset int) utils.DirContents {
    path := dir + argument
    return utils.GetDirContents(path, limit, offset, dir)
}

// TODO: Make this a method that does the writting & just takes debug data
func invalidDirMsg(directory string, filename string) HttpError {
    err := HttpError{
        Error: "This is not a valid directory: " + directory + " " + filename,
        Debug: "Not in valid dirs",
    }
    return err
}
