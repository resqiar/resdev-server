package routes

import (
	"resqiar.com-server/handlers"
	"resqiar.com-server/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitBlogRoute(server *fiber.App, handler handlers.BlogHandler) {
	blog := server.Group("/blog")

	// ONLY SEND PUBLISHED BLOGS FOR PUBLIC
	// drafted/unpublished blogs must only
	// be available to its author scope.
	blog.Get("/list", handler.SendPublishedBlogs)
	blog.Get("/list/id", handler.SendPublishedBlogsID)
	blog.Get("/get/:id", handler.SendPublishedBlog)

	blog.Get("/list/current", middlewares.ProtectedRoute, handler.SendCurrentUserBlogs)
	blog.Post("/get/preview", middlewares.ProtectedRoute, handler.SendCurrentUserBlog)
	blog.Post("/get/my", middlewares.ProtectedRoute, handler.SendMyBlog)

	blog.Post("/create", middlewares.ProtectedRoute, middlewares.TesterRoute, handler.SendBlogCreate)
	blog.Post("/publish", middlewares.ProtectedRoute, middlewares.TesterRoute, handler.SendPublishBlog)
	blog.Post("/unpublish", middlewares.ProtectedRoute, middlewares.TesterRoute, handler.SendUnpublishBlog)
	blog.Post("/update", middlewares.ProtectedRoute, middlewares.TesterRoute, handler.SendUpdateBlog)

	// =========== SPECIAL ROUTES FOR ADM ONLY ===========
	blogADM := server.Group("/blog/adm", middlewares.ProtectedRoute, middlewares.AdminRoute)
	blogADM.Get("/list", handler.SendBlogList)
}
