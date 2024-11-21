package handlers

import (
    "Forum/database"
    "Forum/post_logic"
    "context"
    "encoding/json"
    "html/template"
    "fmt"
    "log"
    "net/http"
	"strconv"    
)

// Context key for user
type contextKey string

const userKey contextKey = "user"

// SessionMiddleware ensures that the user is authenticated and authorized
func SessionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session, err := store.Get(r, "session-name")
        if err != nil || session.Values["user_uuid"] == nil || session.IsNew || session.Values["role"] == "GUEST" {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        userUUID, ok := session.Values["user_uuid"].(string)
        if !ok {
            http.Error(w, "Failed to get user ID from session", http.StatusInternalServerError)
            return
        }

        user, err := database.FetchUserByUUID(userUUID)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        ctx := context.WithValue(r.Context(), userKey, user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// ModeratorProfileHandler handles the moderator profile page
func ModeratorProfileHandler(w http.ResponseWriter, r *http.Request) {
    user, ok := r.Context().Value(userKey).(*database.User)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    if user.Role != "moderator" {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    reports, err := database.FetchAllReports()
    if err != nil {
        log.Printf("Error fetching reports: %v", err)
        http.Error(w, "Internal Server Error for reports", http.StatusInternalServerError)
        return
    }
    var data struct {
        User    database.User
        Reports []database.Report
    }
    data.User = *user
    data.Reports = reports
    tmpl := template.Must(template.ParseFiles("templates/profile_moderator.html"))
    if err := tmpl.Execute(w, data); err != nil {
        log.Printf("Error executing template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}

// AdminProfileHandler handles the admin profile page
func AdminProfileHandler(w http.ResponseWriter, r *http.Request) {
    user, ok := r.Context().Value(userKey).(*database.User)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    if user.Role != "ADMIN" {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    posts, err := database.GetPosts()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    var userPosts []post_logic.Post
    for _, post := range posts {
        userPosts = append(userPosts, post_logic.Post(post))
    }
    utilisateurList, err := database.FetchAllUserUtilisateur()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    moderatorList, err := database.FetchAllUserModerator()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    adminList, err := database.FetchAllUserAdmin()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    reports, err := database.FetchAllReports()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    var data struct {
        User         database.User
        Posts        []post_logic.Post
        Utilisateurs []database.User
        Moderators   []database.User
        Admins       []database.User
        Reports      []database.Report
    }
    data.User = *user
    data.Posts = userPosts
    data.Utilisateurs = utilisateurList
    data.Moderators = moderatorList
    data.Admins = adminList
    data.Reports = reports
    tmpl := template.Must(template.ParseFiles("templates/profile_admin.html"))
    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}

// UtilisateurProfileHandler handles the user profile page
func UtilisateurProfileHandler(w http.ResponseWriter, r *http.Request) {
    user, ok := r.Context().Value(userKey).(*database.User)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    if user.Role != "utilisateur" {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    posts, err := database.GetPosts()
    if err != nil {
        http.Error(w, "Internal Server Error: unable to get posts", http.StatusInternalServerError)
        return
    }
    var userPosts []post_logic.Post
    for _, post := range posts {
        userPosts = append(userPosts, post_logic.Post(post))
    }
    data := struct {
        User  database.User
        Posts []post_logic.Post
    }{
        User:  *user,
        Posts: userPosts,
    }
    tmpl, err := template.ParseFiles("templates/profile_utilisateur.html")
    if err != nil {
        http.Error(w, "Internal Server Error: unable to parse template", http.StatusInternalServerError)
        return
    }
    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Internal Server Error: unable to execute template", http.StatusInternalServerError)
        return
    }
}

// UpdateProfileHandler updates the user profile
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    session, err := store.Get(r, "session-name")
    if err != nil || session.Values["user_uuid"] == nil {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    userUUID, ok := session.Values["user_uuid"].(string)
    if !ok {
        http.Error(w, "Failed to get user ID from session", http.StatusInternalServerError)
        return
    }
    username := r.FormValue("username")
    email := r.FormValue("email")
    err = database.UpdateUserProfile(userUUID, username, email)
    if err != nil {
        log.Printf("Failed to update user profile: %v", err)
        http.Error(w, "Failed to update profile", http.StatusInternalServerError)
        return
    }
    session.Values["username"] = username
    session.Save(r, w)
    http.Redirect(w, r, "/profile/utilisateur", http.StatusSeeOther)
}

// CheckAuthorizationHandler checks user authorization
func CheckAuthorizationHandler(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")
    if err != nil || session.IsNew {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    userUUID, ok := session.Values["user_uuid"].(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    role, err := database.GetUserRoleByUUID(userUUID)
    if err != nil {
        log.Printf("Failed to get user role: %v", err)
        http.Error(w, "Failed to validate user role", http.StatusInternalServerError)
        return
    }

    response := map[string]string{
        "role": role,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// DeletePostHandler handles post deletion
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    postID := r.FormValue("postID")
    if postID == "" {
        http.Error(w, "Post ID is required", http.StatusBadRequest)
        return
    }

    postIDInt, err := strconv.Atoi(postID)
    if err != nil {
        http.Error(w, "Invalid Post ID", http.StatusBadRequest)
        return
    }

    err = database.DeletePostByID(postIDInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteCommentHandler handles comment deletion
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    commentID := r.FormValue("commentID")
    if commentID == "" {
        http.Error(w, "Comment ID is required", http.StatusBadRequest)
        return
    }

    commentIDInt, err := strconv.Atoi(commentID)
    if err != nil {
        http.Error(w, "Invalid Comment ID", http.StatusBadRequest)
        return
    }

    err = database.DeleteCommentByID(commentIDInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreateAdminHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.IsNew {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	uuid := r.FormValue("uuid")
	err = database.CreateAdminUser(uuid, username, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/profile/ADMIN", http.StatusSeeOther)
}

func PromoteUserToModeratorHandler(w http.ResponseWriter, r *http.Request) {
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		http.Error(w, "User UUID is required", http.StatusBadRequest)
		return
	}
	// Suppression de la conversion en int et de la vérification d'erreur associée
	err := PromoteUserToModerator(userUUID) // Assurez-vous que PromoteUserToModerator accepte maintenant une chaîne
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/profile/ADMIN", http.StatusSeeOther)
}

func DemoteModeratorToUserHandler(w http.ResponseWriter, r *http.Request) {
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		http.Error(w, "User UUID is required", http.StatusBadRequest)
		return
	}
	// Suppression de la conversion en int et de la vérification d'erreur associée
	err := DemoteModeratorToUser(userUUID) // Assurez-vous que DemoteModeratorToUser accepte maintenant une chaîne
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/profile/ADMIN", http.StatusSeeOther)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userUUID := r.FormValue("userUUID")
	if userUUID == "" {
		http.Error(w, "User UUID is required", http.StatusBadRequest)
		return
	}
	// Suppression de la conversion en int et de la vérification d'erreur associée
	err := database.DeleteUser(userUUID) // Assurez-vous que DeleteUser accepte maintenant une chaîne
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/profile/ADMIN", http.StatusSeeOther)
}

func SendReportHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.IsNew {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	moderatorUUID := r.FormValue("userUUID")
	moderatorName := r.FormValue("username")

	content := r.FormValue("reportContent")
	
	postID, err := strconv.Atoi(r.FormValue("postID"))
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}
	postTitle := r.FormValue("postTitle")
	postContent := r.FormValue("postContent")
	err = ReceiveReport(database.Report{ModeratorUUID: moderatorUUID, ModeratorName: moderatorName, Content: content, Resolved: false, Response: "", PostID: postID, PostTitle: postTitle, PostContent: postContent})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RespondToReportHandler(w http.ResponseWriter, r *http.Request) {
	reportID := r.FormValue("reportID")
	if reportID == "" {
		http.Error(w, "Report ID is required", http.StatusBadRequest)
		return
	}
	reportIDInt, err := strconv.Atoi(reportID)
	if err != nil {
		http.Error(w, "Invalid Report ID", http.StatusBadRequest)
		return
	}
	response := r.FormValue("responseContent")
	err = database.UpdateReportResponse(reportIDInt, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/profile/ADMIN", http.StatusSeeOther)
}

func DeleteReportHandler(w http.ResponseWriter, r *http.Request) {
	reportID := r.FormValue("reportID")
	role := r.FormValue("userRole")
	if reportID == "" {
		http.Error(w, "Report ID is required", http.StatusBadRequest)
		return
	}
	reportIDInt, err := strconv.Atoi(reportID)
	if err != nil {
		http.Error(w, "Invalid Report ID", http.StatusBadRequest)
		return
	}
	err = database.DeleteReportByID(reportIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/profile/%s", role), http.StatusSeeOther)
}

