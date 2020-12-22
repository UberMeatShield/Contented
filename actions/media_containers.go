package actions

import (
    "log"
	"contented/models"
	"fmt"
	"net/http"

    "github.com/gofrs/uuid"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/responder"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (MediaContainer)
// DB Table: Plural (media_containers)
// Resource: Plural (MediaContainers)
// Path: Plural (/media_containers)
// View Template Folder: Plural (/templates/media_containers/)

// MediaContainersResource is the resource for the MediaContainer model
type MediaContainersResource struct {
	buffalo.Resource
}

// List gets all MediaContainers. This function is mapped to the path
// GET /media_containers
func (v MediaContainersResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	mediaContainers := &models.MediaContainers{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

    c_id := c.Param("container_id")
    if c_id != "" {
        log.Printf("Attempting to get media using %s", c_id)
        container_id, err := uuid.FromString(c_id)
        if err != nil {
            return err
        }
        q_conn := q.Where("container_id = ?", container_id)
        if q_err := q_conn.All(mediaContainers); q_err != nil {
            return q_err
        }
    } else {
        // Retrieve all MediaContainers from the DB
        if err := q.All(mediaContainers); err != nil {
            return err
        }
    }

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		/*
		   c.Set("pagination", q.Paginator)
		   c.Set("mediaContainers", mediaContainers)
		   return c.Render(http.StatusOK, r.HTML("/media_containers/index.plush.html"))
		*/
		return c.Render(200, r.JSON(mediaContainers))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(mediaContainers))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(mediaContainers))
	}).Respond(c)
}

// Show gets the data for one MediaContainer. This function is mapped to
// the path GET /media_containers/{media_container_id}
func (v MediaContainersResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty MediaContainer
	mediaContainer := &models.MediaContainer{}

	// To find the MediaContainer the parameter media_container_id is used.
	if err := tx.Find(mediaContainer, c.Param("media_container_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		/*
		   c.Set("mediaContainer", mediaContainer)
		   return c.Render(http.StatusOK, r.HTML("/media_containers/show.plush.html"))
		*/
		return c.Render(200, r.JSON(mediaContainer))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(mediaContainer))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(mediaContainer))
	}).Respond(c)
}

// Create adds a MediaContainer to the DB. This function is mapped to the
// path POST /media_containers
func (v MediaContainersResource) Create(c buffalo.Context) error {
	// Allocate an empty MediaContainer
	mediaContainer := &models.MediaContainer{}

	// Bind mediaContainer to the html form elements
	if err := c.Bind(mediaContainer); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(mediaContainer)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			/*
			   c.Set("errors", verrs)

			   // Render again the new.html template that the user can
			   // correct the input.
			   c.Set("mediaContainer", mediaContainer)

			   return c.Render(http.StatusUnprocessableEntity, r.HTML("/media_containers/new.plush.html"))
			*/
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		/*
		   c.Flash().Add("success", T.Translate(c, "mediaContainer.created.success"))

		   // and redirect to the show page
		   return c.Redirect(http.StatusSeeOther, "/media_containers/%v", mediaContainer.ID)
		*/
		return c.Render(http.StatusCreated, r.JSON(mediaContainer))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(mediaContainer))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(mediaContainer))
	}).Respond(c)
}

// Update changes a MediaContainer in the DB. This function is mapped to
// the path PUT /media_containers/{media_container_id}
func (v MediaContainersResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty MediaContainer
	mediaContainer := &models.MediaContainer{}

	if err := tx.Find(mediaContainer, c.Param("media_container_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind MediaContainer to the html form elements
	if err := c.Bind(mediaContainer); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(mediaContainer)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			/*
			   c.Set("errors", verrs)

			   // Render again the edit.html template that the user can
			   // correct the input.
			   c.Set("mediaContainer", mediaContainer)

			   return c.Render(http.StatusUnprocessableEntity, r.HTML("/media_containers/edit.plush.html"))
			*/
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		/*
		   c.Flash().Add("success", T.Translate(c, "mediaContainer.updated.success"))

		   // and redirect to the show page
		   return c.Redirect(http.StatusSeeOther, "/media_containers/%v", mediaContainer.ID)
		*/
		return c.Render(http.StatusOK, r.JSON(mediaContainer))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(mediaContainer))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(mediaContainer))
	}).Respond(c)
}

// Destroy deletes a MediaContainer from the DB. This function is mapped
// to the path DELETE /media_containers/{media_container_id}
func (v MediaContainersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty MediaContainer
	mediaContainer := &models.MediaContainer{}

	// To find the MediaContainer the parameter media_container_id is used.
	if err := tx.Find(mediaContainer, c.Param("media_container_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(mediaContainer); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		/*
		   c.Flash().Add("success", T.Translate(c, "mediaContainer.destroyed.success"))

		   // Redirect to the index page
		   return c.Redirect(http.StatusSeeOther, "/media_containers")
		*/
		return c.Render(http.StatusOK, r.JSON(mediaContainer))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(mediaContainer))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(mediaContainer))
	}).Respond(c)
}
