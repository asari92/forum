package handler

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

const sessionContextKey = contextKey("session")

const csrfTokenContextKey = contextKey("csrfToken")
