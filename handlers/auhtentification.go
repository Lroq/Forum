package handlers

import (
	"Forum/auth"
	"Forum/database"
	"encoding/json"

	"golang.org/x/oauth2"

	"html/template"
	"net/http"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		hashedPassword, err := database.HashPassword(password)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		role := "utilisateur"
		err = database.InsertUser(username, email, hashedPassword, role)
		if err != nil {
			http.Error(w, "Email or Username already taken", http.StatusConflict)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/signup.html"))
	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		user, err := database.GetUserByEmail(email)
		if err != nil || !database.CheckPasswordHash(password, user.HashedPassword) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		if !database.CheckPasswordHash(password, user.HashedPassword) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, "Failed to get session", http.StatusInternalServerError)
			return
		}
		session.Values["user_uuid"] = user.UUID
		session.Values["username"] = user.Username
		session.Values["email"] = user.Email
		session.Values["password"] = user.HashedPassword
		userRole, err := database.GetUserRoleByUUID(user.UUID)
		if err != nil {
			http.Error(w, "Failed to get user role", http.StatusInternalServerError)
			return
		}
		session.Values["role"] = userRole
		if err = session.Save(r, w); err != nil {
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
	}
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := auth.GenerateStateOauthCookie(w)
	url := auth.GoogleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	token, err := auth.GoogleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	client := auth.GoogleOauthConfig.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := database.GetUserByEmail(userInfo.Email)
	if err != nil {
		// Si l'utilisateur n'existe pas, créez un nouveau compte
		hashedPassword, _ := database.HashPassword("unMotDePasseGénéréAléatoirement") // Générez un mot de passe sécurisé
		role := "utilisateur"
		err = database.InsertUser(userInfo.Name, userInfo.Email, hashedPassword, role)
		if err != nil {
			http.Error(w, "Failed to create user account: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Récupérez l'utilisateur nouvellement créé
		user, err = database.GetUserByEmail(userInfo.Email)
		if err != nil {
			http.Error(w, "Failed to retrieve user after creation: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Configurez la session pour l'utilisateur
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	session.Values["user_uuid"] = user.UUID
	session.Values["username"] = user.Username
	session.Values["email"] = user.Email
	session.Values["role"] = user.Role // Assurez-vous que le champ 'Role' existe dans votre structure utilisateur
	if err = session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := auth.GenerateStateOauthCookie(w)
	url := auth.GithubOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	token, err := auth.GithubOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	client := auth.GithubOauthConfig.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var userInfo struct {
		Login string `json:"login"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Vérifiez si l'utilisateur existe déjà
	user, err := database.GetUserByEmail(userInfo.Email)
	if err != nil {
		// Si l'utilisateur n'existe pas, créez un nouveau compte
		hashedPassword, _ := database.HashPassword("unMotDePasseGénéréAléatoirement")
		role := "utilisateur"
		err = database.InsertUser(userInfo.Login, userInfo.Email, hashedPassword, role)
		if err != nil {
			http.Error(w, "Failed to create user account: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Récupérez l'utilisateur nouvellement créé
		user, err = database.GetUserByEmail(userInfo.Email)
		if err != nil {
			http.Error(w, "Failed to retrieve user after creation: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Configurez la session pour l'utilisateur
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	session.Values["user_uuid"] = user.UUID
	session.Values["username"] = user.Username
	session.Values["email"] = user.Email
	session.Values["role"] = "utilisateur"
	if err = session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func FacebookLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := auth.GenerateStateOauthCookie(w)
	url := auth.FacebookOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func FacebookCallbackHandler(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	token, err := auth.FacebookOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	client := auth.FacebookOauthConfig.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var userInfo struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Vérifiez si l'utilisateur existe déjà
	user, err := database.GetUserByEmail(userInfo.Email)
	if err != nil {
		// Si l'utilisateur n'existe pas, créez un nouveau compte
		hashedPassword, _ := database.HashPassword("unMotDePasseGénéréAléatoirement")
		role := "utilisateur"
		err = database.InsertUser(userInfo.Name, userInfo.Email, hashedPassword, role)
		if err != nil {
			http.Error(w, "Failed to create user account: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Récupérez l'utilisateur nouvellement créé
		user, err = database.GetUserByEmail(userInfo.Email)
		if err != nil {
			http.Error(w, "Failed to retrieve user after creation: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Configurez la session pour l'utilisateur
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	session.Values["user_uuid"] = user.UUID
	session.Values["username"] = user.Username
	session.Values["email"] = user.Email
	session.Values["role"] = "utilisateur"
	if err = session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
