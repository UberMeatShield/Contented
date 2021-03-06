package actions
/**
*  These are helpers for use in grifts, but we want them compiling in the dev service in case of breaks.
*
* Bad code in a grift is harder to notice and the compilation with tests also seems a little broken. ie
* you break the grift via main package changes and never notice.  If You break the test in a grift directory
* and then the compilation just failed with no error messages...
*/

import (
    "os"
    "log"
    // "time"
    "strings"
    "path/filepath"
    "contented/models"
    "contented/utils"
    "github.com/pkg/errors"
    "github.com/gofrs/uuid"
    "github.com/gobuffalo/nulls"
)


// Process all the directories and get a valid setup into the DB
// Probably should return a count of everything
func CreateInitialStructure(cfg *utils.DirConfigEntry) error {
    dirs := utils.FindContainers(cfg.Dir)
    log.Printf("Found %d sub-directories.\n", len(dirs))
    if len(dirs) == 0 {
        return errors.New("No subdirectories found under path: " + cfg.Dir)
    }

    // Optional?  Some sort of crazy merge for later?
    err := models.DB.TruncateAll()
    if err != nil {
        return errors.WithStack(err)
    }
    // This should be initialized

    // TODO: Need to do this in a single transaction vs partial
    for _, dir := range dirs {
        log.Printf("Adding directory %s with id %s\n", dir.Name, dir.ID)

        // A more sensible limit on the absolute max lookup?
        media := utils.FindMediaMatcher(dir, 90001, 0, cfg.IncFiles, cfg.ExcFiles) 
        log.Printf("Adding Media to %s with total media %d \n", dir.Name, len(media))

        // Use the database version of uuid generation (minimize the miniscule conflict)
        unset_uuid, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
        dir.ID = unset_uuid
        dir.Total = len(media)
        models.DB.Create(&dir)
        log.Printf("Created %s with id %s\n", dir.Name, dir.ID)

        for _, mc := range media {
            mc.ContainerID = nulls.NewUUID(dir.ID) 
            c_err := models.DB.Create(&mc)
            if c_err != nil {
                log.Fatal(c_err)
            }
        }
    }
    return nil
}

// TODO: Move this code into manager (likely?)
func ClearContainerPreviews(c *models.Container) error {
    dst := GetContainerPreviewDst(c) 
    if _, err := os.Stat(dst); os.IsNotExist(err) {
        return nil
    }
    r_err := os.RemoveAll(dst)
    if r_err != nil {
        log.Fatal(r_err)
        return r_err
    }
    return nil
}

// TODO: Move to utils or make it wrapped for some reason?
func GetContainerPreviewDst(c *models.Container) string {
    cfg := utils.GetCfg()
    return filepath.Join(cfg.Dir, c.Name, "container_previews")
}

// Init a manager and pass it in or just do this via config value instead of a pass in
func CreateAllPreviews(cm ContentManager) error {
    cnts, c_err := cm.ListContainers(0, 9001)
    if c_err != nil {
        return c_err
    }
    if len(*cnts) == 0 {
        return errors.New("No Containers were found in the database")
    }
    for _, cnt := range *cnts {
        err := CreateContainerPreviews(&cnt, cm)    
        if err != nil {
            return err
        }
    }
    return nil
}

// TODO: Should this return a total of previews created or something?
func CreateContainerPreviews(c *models.Container, cm ContentManager) error {
    log.Printf("About to try and create previews for %s:%s\n", c.Name, c.ID.String())
    // Reset the preview directory, then create it fresh (update tests if this changes)
    c_err := ClearContainerPreviews(c)
    if c_err == nil {
        err := utils.MakePreviewPath(GetContainerPreviewDst(c))
        if err != nil {
            log.Fatal(err)
        }
    }

    // TODO: It should fix up the total count there (-1 for unlimited?)
    media, q_err := cm.ListMedia(c.ID, 0, 90000)
    if q_err != nil {
        log.Fatal(q_err)
        return q_err
    }
    log.Printf("Found a set of media to make previews for %d", len(*media))

    update_list, err := CreateMediaPreviews(c, *media)
    if err != nil {
        return err 
    }
    log.Printf("Finished creating previews, updating the database %d", len(update_list))
    for _, mc := range update_list {
        if mc.Preview != "" {
            log.Printf("Created a preview %s for mc %s", mc.Preview, mc.ID.String())
            cm.UpdateMedia(&mc)
        }
    }
    return nil
}

func CreateMediaPreviews(c *models.Container, media models.MediaContainers) (models.MediaContainers, error) {
    if len(media) == 0 {
        return models.MediaContainers{}, nil 
    }
    cfg := utils.GetCfg()
    processors := cfg.CoreCount
    if processors <= 0 {
        processors = 1 // Without at least one processor this will hang forever
    }
    log.Printf("Creating %d listeners for processing previews", processors)

    // We expect a result for every message so can create the channels in a way that they have a length
    expected_total := len(media)
    reply := make(chan utils.PreviewResult, expected_total)
    input := make(chan utils.PreviewRequest, expected_total)

    // Starts the workers
    for i := 0; i < processors; i++ {
        pw := utils.PreviewWorker{Id: i, In: input}
        go StartWorker(pw)
    }

    // Queue up a bunch of preview work
    mediaMap := models.MediaMap{}
    for _, mc := range media {
        mediaMap[mc.ID] = mc

        ref_mc := mc
        input <- utils.PreviewRequest{
            C: c,
            Mc: &ref_mc,
            Out: reply,
        }
    }

    // Exception handling should close the input and output probably
    total := 0
    previews := models.MediaContainers{}

    error_list := ""
    for result := range reply {
        total++
        if total == expected_total {
            close(input)  // Do I close this immediately
            close(reply)
        } 

        // Get a list of just the preview items?  Or just update by reference?
        if mc_update, ok := mediaMap[result.MC_ID]; ok {
            if (result.Preview != "") {
                log.Printf("We found a reply around this %s id was %s \n", result.Preview, result.MC_ID)
                mc_update.Preview = result.Preview
                previews = append(previews, mc_update)
            } else {
                log.Printf("No preview was needed for media %s", result.MC_ID)
            }
        } else {
            log.Printf("Missing Response ID, something went wrong %s\n", result.MC_ID)
        }
        if result.Err != nil {
            log.Printf("ERROR: Failed to create a preview %s\n", result.Err)
            error_list += "" + result.Err.Error()
        }
        log.Printf("Found a result for %s\n", result.MC_ID.String())
    }
    if error_list != "" {
        return previews, errors.New(error_list)
    }
    return previews, nil
}

func StartWorker(w utils.PreviewWorker) {
    // sleepTime := time.Duration(w.Id) * time.Millisecond
    // log.Printf("Worker %d with sleep %d\n", w.Id, sleepTime)
    // Sleep before kicking off?  Kinda don't need to
    for pr := range w.In {
        c := pr.C
        mc := pr.Mc
        log.Printf("Worker %d Doing a preview for %s\n", w.Id, mc.ID.String())
        preview, err :=  CreateMediaPreview(c, mc)
        pr.Out <- utils.PreviewResult{
            C_ID: c.ID,
            MC_ID: mc.ID,
            Preview: preview,
            Err: err,
        }
    }
}


// This might not need to be a fatal on an error, but is nice for debugging now
func CreateMediaPreview(c *models.Container, mc *models.MediaContainer) (string, error) {
    cfg := utils.GetCfg()
    cntPath := filepath.Join(cfg.Dir, c.Name)
    dstPath := GetContainerPreviewDst(c)

    // TODO: ensure that other content does not explode...
    dstFqPath, err := utils.GetImagePreview(cntPath, mc.Src, dstPath, cfg.PreviewOverSize)
    if err != nil {
        log.Fatal(err)
        return "", err
    }
    return strings.ReplaceAll(dstFqPath, cntPath, ""), err
}
