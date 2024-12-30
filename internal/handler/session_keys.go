package handler

// Ключи для значений, хранимых в сессиях
const (
	FlashSessionKey                  = "flash"
	AuthUserIDSessionKey             = "authenticated_userID"
	UserRoleSessionKey               = "user_role"
	CsrfTokenSessionKey              = "token"
	RedirectPathAfterLoginSessionKey = "redirect_path_after_login"
	ReactionFormSessionKey           = "reaction_form"
	GoogleOAuthStateSessionKey       = "google_oauth_state"
	GitHubOAuthStateSessionKey       = "github_oauth_state"
)
