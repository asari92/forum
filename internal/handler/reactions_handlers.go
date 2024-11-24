package handler

import (
	"errors"
	"net/http"
	"strconv"

	"forum/internal/entities"
)

func (app *Application) postReaction(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err := errors.New("get userID in postReaction")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID < 1 {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	form := app.Service.Reaction.NewReactionForm()
	comment := r.PostForm.Get("comment_content")
	form.Comment = comment
	postIsLike := r.PostForm.Get("post_is_like")
	form.PostIsLike = postIsLike
	commentIsLike := r.PostForm.Get("comment_is_like")
	form.CommentIsLike = commentIsLike
	if commentIsLike != "" {
		commentID, err := strconv.Atoi(r.PostForm.Get("comment_id"))
		if err != nil || commentID < 1 {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
			return
		}
		form.CommentID = commentID
	}

	err = app.Service.Reaction.UpdatePostReaction(userID, postID, &form)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidData) {
			sess.Set(ReactionFormSessionKey, form)
			http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
			return
		}
		app.Logger.Error("update reaction on post in database", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// /comment/reaction
func (app *Application) commentReaction(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err := errors.New("get userID in postReaction")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	form := app.Service.Reaction.NewReactionForm()
	commentID, err := strconv.Atoi(r.PostForm.Get("comment_id"))
	if err != nil || commentID < 1 {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	form.CommentID = commentID
	commentIsLike := r.PostForm.Get("is_like")
	form.CommentIsLike = commentIsLike

	err = app.Service.Reaction.UpdateCommentReaction(userID, &form)
	if err != nil {
		app.Logger.Error("update reaction on post in database", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
