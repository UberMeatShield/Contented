package actions

import (
    "errors"
	"fmt"
	"net/http"
	"contented/models"
//    "errors"

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
// Model: Singular (Container)
// DB Table: Plural (containers)
// Resource: Plural (Containers)
// Path: Plural (/containers)
// View Template Folder: Plural (/templates/containers/)

// ContainersResource is the resource for the Container model
type ContainersResource struct {
	buffalo.Resource
}

// List gets all Containers. This function is mapped to the path
// GET /containers
func (v ContainersResource) List(c buffalo.Context) error {
	// Get the DB connection from the context

    man := GetManager(&c)
    containers, err := man.ListContainersContext()
    if err != nil {
        return c.Error(http.StatusBadRequest, err)
    } 
    // TODO figure out how to set the current context
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
    /*
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	containers := &models.Containers{}
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Containers from the DB
	if err := q.All(containers); err != nil {
		return err
	}
    */

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		/*
		   c.Set("pagination", q.Paginator)
		   c.Set("containers", containers)
		   return c.Render(http.StatusOK, r.HTML("/containers/index.plush.html"))
		*/
		return c.Render(200, r.JSON(containers))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(containers))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(containers))
	}).Respond(c)
}

// Show gets the data for one Container. This function is mapped to
// the path GET /containers/{container_id}
func (v ContainersResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context

    c_id, err := uuid.FromString(c.Param("container_id"))
    if err != nil {
        // return c.Error(http.StatusBadRequest, errors.New("Invalid Container ID"))
        return c.Error(http.StatusBadRequest, err)
    } // Hate

    man := GetManager(&c)
    container, err := man.GetContainer(c_id)
    if err != nil {
		return c.Error(http.StatusNotFound, err)
    }
    /*
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Container
	container := &models.Container{}

	// To find the Container the parameter container_id is used.
	if err := tx.Find(container, c.Param("container_id")); err != nil {
	}
    */

	return responder.Wants("html", func(c buffalo.Context) error {
		/*
		   c.Set("container", container)
		   return c.Render(http.StatusOK, r.HTML("/containers/show.plush.html"))
		*/

		return c.Render(200, r.JSON(container))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(container))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(container))
	}).Respond(c)
}

// Create adds a Container to the DB. This function is mapped to the
// path POST /containers
func (v ContainersResource) Create(c buffalo.Context) error {
	// Allocate an empty Container

    // TODO: Reject if it is memory manager
    man := GetManager(&c)
    if man.CanEdit() == false {
        return c.Error(
            http.StatusNotImplemented,
            errors.New("Edit not supported by this manager"),
        )
    }

	// Bind container to the html form elements
	container := &models.Container{}
	if err := c.Bind(container); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(container)
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
			   c.Set("container", container)

			   return c.Render(http.StatusUnprocessableEntity, r.HTML("/containers/new.plush.html"))
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
		//c.Flash().Add("success", T.Translate(c, "container.created.success"))
		//return c.Redirect(http.StatusSeeOther, "/containers/%v", container.ID)

		// and redirect to the show page
		return c.Render(http.StatusCreated, r.JSON(container))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(container))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(container))
	}).Respond(c)
}

// Update changes a Container in the DB. This function is mapped to
// the path PUT /containers/{container_id}
func (v ContainersResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
    man := GetManager(&c)
    if man.CanEdit() == false {
        return c.Error(
            http.StatusNotImplemented,
            errors.New("Edit not supported by this manager"),
        )
    }


	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Container
	container := &models.Container{}

	if err := tx.Find(container, c.Param("container_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Container to the html form elements
	if err := c.Bind(container); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(container)
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
			   c.Set("container", container)

			   return c.Render(http.StatusUnprocessableEntity, r.HTML("/containers/edit.plush.html"))
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
		//c.Flash().Add("success", T.Translate(c, "container.updated.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/containers/%v", container.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(container))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(container))
	}).Respond(c)
}

// Destroy deletes a Container from the DB. This function is mapped
// to the path DELETE /containers/{container_id}
func (v ContainersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
    man := GetManager(&c)
    if man.CanEdit() == false {
        return c.Error(
            http.StatusNotImplemented,
            errors.New("Edit not supported by this manager"),
        )
    }

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Container
	container := &models.Container{}

	// To find the Container the parameter container_id is used.
	if err := tx.Find(container, c.Param("container_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(container); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		// c.Flash().Add("success", T.Translate(c, "container.destroyed.success"))

		// Redirect to the index page
		//return c.Redirect(http.StatusSeeOther, "/containers")
		return c.Render(http.StatusOK, r.JSON(container))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(container))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(container))
	}).Respond(c)
}
