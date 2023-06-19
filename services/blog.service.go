package services

import (
	"resqiar.com-server/entities"
	"resqiar.com-server/inputs"
	"resqiar.com-server/repositories"
	"time"
)

type BlogService interface {
	GetAllBlogs(onlyPublished bool) ([]entities.SafeBlogAuthor, error)
	GetAllBlogsID() ([]string, error)
	GetBlogDetail(blogID string, published bool) (*entities.SafeBlogAuthor, error)
	CreateBlog(payload *inputs.CreateBlogInput, userID string) (*entities.Blog, error)
	EditBlog(payload *inputs.UpdateBlogInput, userID string) error
	GetCurrentUserBlogs(userID string) ([]entities.Blog, error)
	ChangeBlogPublish(payload *inputs.BlogIDInput, userID string, publishState bool) error
}

type BlogServiceImpl struct {
	Repository repositories.BlogRepository
}

// GetAllBlogs retrieves a list of SafeBlogAuthor entities from the database.
// If onlyPublished is true, it retrieves only the published blogs, otherwise everything.
// It returns the list of blogs and any error encountered during the process.
func (service *BlogServiceImpl) GetAllBlogs(onlyPublished bool) ([]entities.SafeBlogAuthor, error) {
	blogs, err := service.Repository.GetBlogs(onlyPublished)
	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func (service *BlogServiceImpl) GetAllBlogsID() ([]string, error) {
	blogs, err := service.GetAllBlogs(true) // get all published blogs
	if err != nil {
		return nil, err
	}

	var result []string

	// only get the id, and append it to array if of string
	for _, blog := range blogs {
		ID := blog.ID
		result = append(result, ID)
	}

	return result, nil
}

// GetPublishedBlogDetail retrieves a single SafeBlogAuthor entity from the database
// based on the provided blogID.
// It returns the retrieved blog and any error encountered during the process.
// If no blog is found or an error occurs, it returns an appropriate error.
func (service *BlogServiceImpl) GetBlogDetail(blogID string, published bool) (*entities.SafeBlogAuthor, error) {
	blog, err := service.Repository.GetBlog(blogID, published)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (service *BlogServiceImpl) CreateBlog(payload *inputs.CreateBlogInput, userID string) (*entities.Blog, error) {
	newBlog := entities.Blog{
		Title:   payload.Title,
		Summary: payload.Summary,
		Content: payload.Content,

		// when creating blog, always set published to false.
		// although the default value in database is false,
		// we still want to ensure the published value here-
		// is NOT coming from the payload, but rather hardcoded.
		Published: false,

		CoverURL: payload.CoverURL, AuthorID: userID,
	}

	result, err := service.Repository.CreateBlog(&newBlog)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *BlogServiceImpl) EditBlog(payload *inputs.UpdateBlogInput, userID string) error {
	blog, err := service.Repository.GetByIDAndAuthor(payload.ID, userID)
	if err != nil {
		return err
	}

	safe := &inputs.SafeUpdateBlogInput{
		Title:    payload.Title,
		Summary:  payload.Summary,
		Content:  payload.Content,
		CoverURL: payload.CoverURL,
	}

	if err := service.Repository.UpdateBlog(blog.ID, safe); err != nil {
		return err
	}

	return nil
}

func (service *BlogServiceImpl) GetCurrentUserBlogs(userID string) ([]entities.Blog, error) {
	blogs, err := service.Repository.GetCurrentUserBlogs(userID)
	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func (service *BlogServiceImpl) ChangeBlogPublish(payload *inputs.BlogIDInput, userID string, publishState bool) error {
	blog, err := service.Repository.GetByIDAndAuthor(payload.ID, userID)
	if err != nil {
		return err
	}

	// update published state based on given param
	blog.Published = publishState

	// if publish state is true
	// then we need to update the PublishedAt field
	if publishState {
		blog.PublishedAt = time.Now()
	} else {
		// otherwise, reset the PublishedAt field to "January 1, year 1, 00:00:00 UTC" (invalid date)
		blog.PublishedAt = time.Time{}
	}

	// change the updated at to newest date
	blog.UpdatedAt = time.Now()

	// save back to the database
	if err := service.Repository.SaveBlog(blog); err != nil {
		return err
	}

	return nil
}
