package database

import (
	"strings"

	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/internal/username"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func Strptr(str string) *string {
	return &str
}

func generateSeeds() {
	log.Debug("Generating database migrations...")
	DbConnection.AutoMigrate(
		&model.Role{},
		&model.Category{},
		&model.User{},
		&model.FriendRequest{},
		&model.Permission{},
		&model.Post{},
		&model.PostAttachment{},
		&model.Comment{},
	)

	backfillUsernames()

	// Seed the database
	//! Must be removed in production
	seedPermissions()
	seedRoles()
	seedUsers()
	seedCategories()
	// seedPosts()

	log.Debug("Finished generating database migrations.")
}

func backfillUsernames() {
	var users []model.User
	if err := DbConnection.Find(&users).Error; err != nil {
		log.Error(err)
		return
	}

	for _, user := range users {
		if user.Username != nil && strings.TrimSpace(*user.Username) != "" {
			continue
		}

		fullName := "User"
		if user.FullName != nil && strings.TrimSpace(*user.FullName) != "" {
			fullName = *user.FullName
		}

		usernameValue, err := username.Generate(fullName, func(candidate string) (bool, error) {
			var count int64
			query := DbConnection.Model(&model.User{}).Where("username = ?", candidate).Where("id <> ?", user.ID)
			if err := query.Count(&count).Error; err != nil {
				return false, err
			}
			return count > 0, nil
		})
		if err != nil {
			log.Error(err)
			continue
		}

		if err := DbConnection.Model(&user).Update("username", usernameValue).Error; err != nil {
			log.Error(err)
		}
	}
}

func seedPermissions() {
	log.Debug("Seeding permissions...")
	permissionData := []struct {
		name        string
		description string
	}{
		{
			name:        "user:create",
			description: "Create new users",
		},
		{
			name:        "user:read",
			description: "View user information",
		},
		{
			name:        "user:update",
			description: "Update user information",
		},
		{
			name:        "user:delete",
			description: "Update user information",
		},
		{
			name:        "moderator:create",
			description: "Allows moderators to create data",
		},
		{
			name:        "moderator:read",
			description: "Allows moderators to read data",
		},
		{
			name:        "moderator:update",
			description: "Allows moderators to update data",
		},
		{
			name:        "moderator:delete",
			description: "Allows moderators to delete data",
		},
		{
			name:        "admin:access",
			description: "Grants full admin access",
		},
	}

	for _, permission := range permissionData {
		permissionModel := model.Permission{
			Name:        permission.name,
			Description: permission.description,
		}

		if err := DbConnection.FirstOrCreate(&permissionModel, model.Permission{Name: permission.name}).Error; err != nil {
			log.Error(err)
		}
	}
}

func extractPrefixedPermissions(prefix string) []*model.Permission {
	var prefixes []string
	switch prefix {
	case "user":
		prefixes = []string{"user"}
	case "moderator":
		prefixes = []string{"user", "moderator"}
	case "admin":
		prefixes = []string{"user", "moderator", "admin"}
	}

	var prefixedPerms []*model.Permission

	var permissions []model.Permission
	DbConnection.Find(&permissions)

	for _, perm := range permissions {
		for _, currentPrefix := range prefixes {
			if strings.HasPrefix(perm.Name, currentPrefix) {
				prefixedPerms = append(prefixedPerms, &perm)
			}

		}
	}

	return prefixedPerms
}

func seedRoles() {
	log.Debug("Seeding roles...")
	standardRoles := []model.Role{
		{
			Name: "RoleAdmin",
		},
		{
			Name: "RoleModerator",
		},
		{
			Name: "RoleUser",
		},
	}

	for _, role := range standardRoles {
		prefix := strings.ToLower(role.Name[4:])
		perms := extractPrefixedPermissions(prefix)
		role.Permissions = perms

		if err := DbConnection.FirstOrCreate(&role, model.Role{Name: role.Name}).Error; err != nil {
			log.Error(err)
		}
	}
}

func seedCategories() {
	log.Debug("Seeding thematic areas...")
	mainCategories := []model.Category{
		{
			Name: "Philosophy",
		},
		{
			Name: "Literature",
		},
		{
			Name: "History",
		},
		{
			Name: "Art",
		},
		{
			Name: "Humanities and Social Sciences",
		},
		{
			Name: "Antiquity",
		},
		{
			Name: "Middle Ages",
		},
		{
			Name: "Renaissance",
		},
		{
			Name: "Baroque",
		},
		{
			Name: "Enlightenment",
		},
		{
			Name: "XIX century",
		},
		{
			Name: "Modernism",
		},
		{
			Name: "Vanguard",
		},
		{
			Name: "Postmodernity",
		},
		{
			Name: "Modernity",
		},
	}

	for _, mainCategory := range mainCategories {
		if err := DbConnection.FirstOrCreate(&mainCategory, model.Category{Name: mainCategory.Name}).Error; err != nil {
			log.Error(err)
			return
		}
	}

	var categories []model.Category
	if err := DbConnection.Find(&categories).Error; err != nil {
		log.Error(err)
		return
	}

	subCategories := []model.Category{
		//* Philosophy qualifications
		{
			Name:   "A history of ideas",
			Parent: &categories[0],
		},
		{
			Name:   "Ethics",
			Parent: &categories[0],
		},
		{
			Name:   "Metaphysics",
			Parent: &categories[0],
		},
		{
			Name:   "Political philosophy",
			Parent: &categories[0],
		},
		{
			Name:   "Aesthetics",
			Parent: &categories[0],
		},
		//* Literature qualifications
		{
			Name:   "Poetry",
			Parent: &categories[1],
		},
		{
			Name:   "Prose",
			Parent: &categories[1],
		},
		{
			Name:   "Drama",
			Parent: &categories[1],
		},
		{
			Name:   "Essays",
			Parent: &categories[1],
		},
		//* History qualifications
		{
			Name:   "Political history",
			Parent: &categories[2],
		},
		{
			Name:   "Social history",
			Parent: &categories[2],
		},
		{
			Name:   "Cultural history",
			Parent: &categories[2],
		},
		{
			Name:   "Military history",
			Parent: &categories[2],
		},
		//* Art qualifications
		{
			Name:   "Visual arts",
			Parent: &categories[3],
		},
		{
			Name:   "Musics",
			Parent: &categories[3],
		},
		{
			Name:   "Theater",
			Parent: &categories[3],
		},
		{
			Name:   "Cinema",
			Parent: &categories[3],
		},
		{
			Name:   "Architecture and design",
			Parent: &categories[3],
		},
		//* Humanities and Social Sciences qualifications
		{
			Name:   "Anthropology",
			Parent: &categories[4],
		},
		{
			Name:   "Sociology",
			Parent: &categories[4],
		},
		{
			Name:   "Psychology",
			Parent: &categories[4],
		},
		{
			Name:   "Linguistics",
			Parent: &categories[4],
		},
		{
			Name:   "Political science",
			Parent: &categories[4],
		},
	}

	for _, subCategory := range subCategories {
		if err := DbConnection.Create(&subCategory).Error; err != nil {
			log.Error(err)
			return
		}
	}
}

func seedUsers() {
	log.Debug("Seeding users...")

	var roles []model.Role
	if err := DbConnection.Find(&roles).Error; err != nil {
		log.Fatal(err)
		return
	}

	generatePassword := func(password string) string {
		pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err == nil {
			return string(pass)
		}

		log.Panic(err)
		return ""
	}

	users := []model.User{
		{
			FullName: Strptr("Name Name Name"),
			Username: Strptr("name.name.name"),
			Email:    "standard_user@mail.com",
			Password: Strptr(generatePassword("m123#@S")),
			Role:     &roles[2],
		},
		{
			FullName: Strptr("Name Name Name"),
			Username: Strptr("name.name.name.10"),
			Email:    "moderator_user@mail.com",
			Password: Strptr(generatePassword("m123#@S")),
			Role:     &roles[1],
		},
		{
			FullName: Strptr("Name Name Name"),
			Username: Strptr("name.name.name.11"),
			Email:    "admin_user@mail.com",
			Password: Strptr(generatePassword("m123#@S")),
			Role:     &roles[0],
		},
	}

	for _, user := range users {
		if err := DbConnection.FirstOrCreate(&user, model.User{Email: user.Email}).Error; err != nil {
			log.Error(err)
		}
	}
}

func seedPosts() {
	posts := []model.Post{
		model.Post{
			ID:            0,
			Title:         "Antic history first post",
			Content:       "The content",
			AuthorID:      1,
			Tags:          nil,
			TaggedUsersID: nil,
		},
	}
	if err := DbConnection.Create(&posts).Error; err != nil {
		log.Error(err)
	}
}
