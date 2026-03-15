package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/internal/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func parseUint64SliceField(raw string) ([]uint64, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	var values []uint64
	if strings.HasPrefix(strings.TrimSpace(raw), "[") {
		if err := json.Unmarshal([]byte(raw), &values); err != nil {
			return nil, err
		}
		return values, nil
	}

	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		value, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func parseInt64SliceField(raw string) ([]int64, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	var values []int64
	if strings.HasPrefix(strings.TrimSpace(raw), "[") {
		if err := json.Unmarshal([]byte(raw), &values); err != nil {
			return nil, err
		}
		return values, nil
	}

	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		value, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func parseStringSliceField(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	var values []string
	if strings.HasPrefix(strings.TrimSpace(raw), "[") {
		if err := json.Unmarshal([]byte(raw), &values); err == nil {
			return values
		}
	}

	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			values = append(values, part)
		}
	}

	return values
}

func detectAttachmentKind(contentType string) string {
	if strings.HasPrefix(contentType, "image/") {
		return "image"
	}
	if strings.HasPrefix(contentType, "video/") {
		return "video"
	}
	if contentType == "application/pdf" ||
		strings.Contains(contentType, "msword") ||
		strings.Contains(contentType, "officedocument") ||
		strings.HasPrefix(contentType, "text/") {
		return "document"
	}
	return "file"
}

// CreatePost creates a new post
// @Summary Create a new post
// @Description Creates a new post with title, content, categories, tags, and tagged users
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body model.CreatePost true "Post creation data"
// @Success 201 {object} map[string]interface{} "Created post"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/posts [post]
func CreatePost(c *gin.Context) {
	var post model.CreatePost
	var files []*multipart.FileHeader

	contentType := c.ContentType()

	if strings.Contains(contentType, "multipart/form-data") {
		post.Title = strings.TrimSpace(c.PostForm("title"))
		post.Content = strings.TrimSpace(c.PostForm("content"))

		categoryIDs, err := parseUint64SliceField(c.PostForm("category_ids"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_ids"})
			return
		}
		post.CategoryIDs = categoryIDs

		post.Tags = parseStringSliceField(c.PostForm("tags"))

		taggedUsers, err := parseInt64SliceField(c.PostForm("tagged_users"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tagged_users"})
			return
		}
		post.TaggedUsersID = taggedUsers

		if rawPrivate := strings.TrimSpace(c.PostForm("is_private")); rawPrivate != "" {
			value, err := strconv.ParseBool(rawPrivate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid is_private"})
				return
			}
			post.IsPrivate = &value
		}

		form, err := c.MultipartForm()
		if err != nil && err != http.ErrNotMultipart {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
			return
		}
		if form != nil {
			files = form.File["attachments"]
		}
	} else {
		if err := c.BindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
			return
		}
	}

	user, _ := c.Get("user")
	currentUser := user.(model.User)

	if strings.TrimSpace(post.Title) == "" && strings.TrimSpace(post.Content) != "" {
		runes := []rune(post.Content)
		if len(runes) > 80 {
			post.Title = strings.TrimSpace(string(runes[:80]))
		} else {
			post.Title = strings.TrimSpace(post.Content)
		}
	}

	newPost := model.Post{
		Title:         post.Title,
		Content:       post.Content,
		AuthorID:      currentUser.ID,
		CategoryIDs:   append([]uint64(nil), post.CategoryIDs...),
		Tags:          append([]string(nil), post.Tags...),
		TaggedUsersID: append([]int64(nil), post.TaggedUsersID...),
	}
	if post.IsPrivate != nil {
		newPost.IsPrivate = *post.IsPrivate
	}

	if err := service.SavePost(&newPost); err != nil {
		log.Error("Error while creating post:", err.Error())
		if strings.Contains(err.Error(), "external") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()[10:],
			})
			return
		} else if strings.Contains(err.Error(), "internal") {
			log.Debugf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()[10:],
			})
			return
		}
	}

	if len(files) > 0 {
		uploadDir := filepath.Join("uploads", "posts", strconv.FormatUint(currentUser.ID, 10), strconv.FormatUint(newPost.ID, 10))
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare upload directory"})
			return
		}

		for _, file := range files {
			filename := filepath.Base(file.Filename)
			targetPath := filepath.Join(uploadDir, filename)
			if err := c.SaveUploadedFile(file, targetPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save attachment"})
				return
			}

			attachment := model.PostAttachment{
				PostID:   newPost.ID,
				FileName: filename,
				FilePath: targetPath,
				FileURL:  "/" + filepath.ToSlash(targetPath),
				MimeType: file.Header.Get("Content-Type"),
				Size:     file.Size,
				Kind:     detectAttachmentKind(file.Header.Get("Content-Type")),
			}

			if err := database.DbConnection.Create(&attachment).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist attachment"})
				return
			}
			newPost.Attachments = append(newPost.Attachments, attachment)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": newPost,
	})
}

// GetPosts retrieves a paginated list of posts
// @Summary Get all posts
// @Description Returns a paginated list of posts with optional filtering for authenticated users
// @Tags posts
// @Accept json
// @Produce json
// @Param limit query int false "Number of posts to return (default: 20, max: 200)" default(20)
// @Param offset query int false "Number of posts to skip (default: 0)" default(0)
// @Success 200 {object} map[string]interface{} "List of posts"
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Router /api/posts [get]
func GetPosts(c *gin.Context) {
	u, exists := c.Get("user")
	var user model.User
	if exists == true {
		user, _ = u.(model.User)
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if limit <= 0 || err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0")) // default: 0
	if err != nil || offset < 0 {
		offset = 0
	}

	var (
		posts []model.Post
	)
	posts, err = service.GetPosts(limit, offset, &user)

	if err != nil {
		if strings.Contains(err.Error(), "external") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()[10:],
			})
			return
		} else if strings.Contains(err.Error(), "internal") {
			log.Debugf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()[10:],
			})
			return
		}
	}

	if len(posts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": "category has no current posts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": posts,
	})
}

// UpdatePost updates an existing post
// @Summary Update a post
// @Description Updates the content, title, categories, tags, or privacy settings of a post
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Param input body model.UpdatePost true "Updated post data"
// @Success 200 {object} map[string]interface{} "Updated post"
// @Failure 400 {object} map[string]string "Invalid input or post ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Post not found"
// @Router /api/posts/{id} [put]
func UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid path param",
		})
		return
	}
	var input model.UpdatePost
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|invalid data",
		})

		return
	}
	post := model.Post{
		ID: id,
	}
	updatedPost := model.Post{
		ID:      id,
		Title:   input.Title,
		Content: input.Content,
	}
	copy(updatedPost.Tags, input.Tags)
	copy(updatedPost.TaggedUsersID, input.TaggedUsersID)
	copy(updatedPost.CategoryIDs, input.CategoryIDs)
	if input.IsPrivate != nil {
		post.IsPrivate = *input.IsPrivate
	}

	updatedPost.UpdatedAt = time.Now()
	if err := service.UpdatePost(&post, &updatedPost); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

// DeletePost deletes a post
// @Summary Delete a post
// @Description Permanently deletes a post from the system
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]string "Deletion confirmation"
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Post not found"
// @Router /api/posts/{id} [delete]
func DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|invalid path parameter",
		})
		return
	}
	post := model.Post{
		ID: id,
	}
	if err := service.DeletePost(&post); err != nil {
		if strings.Contains(err.Error(), "external") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()[10:],
			})
			return
		} else if strings.Contains(err.Error(), "internal") {
			log.Debugf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()[10:],
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "deleted category and it's related usages",
	})
}

// GetPost retrieves a specific post by ID
// @Summary Get post by ID
// @Description Returns detailed information about a specific post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{} "Post details"
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 404 {object} map[string]string "Post not found"
// @Router /api/posts/{id} [get]
func GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid path param",
		})
		return
	}

	post := model.Post{
		ID: id,
	}
	if err := service.GetPost(&post); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	// Track view if user is authenticated
	if u, exists := c.Get("user"); exists {
		user, ok := u.(model.User)
		if ok && user.ID > 0 {
			// Track view asynchronously (don't fail if tracking fails)
			if err := service.TrackPostView(&user, &post); err != nil {
				log.Debug("Failed to track post view: ", err)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  post,
	})
}

// LikePost likes or unlikes a post
// @Summary Like/unlike a post
// @Description Toggles the like status for a post. If already liked, removes the like.
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{} "Like count"
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/posts/{id}/like [post]
func LikePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid path param",
		})
		return
	}

	u, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"authenticated": false,
		})
		return
	}

	post := model.Post{ID: id}
	user, _ := u.(model.User)

	likes, err := service.LikePost(&post, &user)
	if err != nil {
		if strings.Contains(err.Error(), "external") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()[10:],
			})
			return
		} else if strings.Contains(err.Error(), "internal") {
			log.Debugf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()[10:],
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"likes": likes,
		},
	})
}
