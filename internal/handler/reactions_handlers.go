package handler

import (
	"errors"
	"forum/internal/entities"
	"net/http"
	"strconv"
)

func (app *Application) postReaction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID < 1 {
		app.render(w, http.StatusNotFound, Errorpage, nil)
		return
	}

	comment := r.PostForm.Get("comment_content")
	postIsLike := r.PostForm.Get("post_is_like")
	commentIsLike := r.PostForm.Get("comment_is_like")

	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err = errors.New("get userID in postReaction")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	if comment != "" {
		err = app.Service.Comment.AddComment(postID, userID, comment)
		if err != nil {
			app.Logger.Error("insert comment to database", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return

		}
	} else if postIsLike != "" {
		// Преобразуем isLike в bool
		like := postIsLike == "true"

		var userReaction *entities.PostReaction
		userReaction, err = app.Service.PostReaction.GetUserReaction(userID, postID) // Получите реакцию пользователя
		if err != nil {
			app.Logger.Error("get post user reaction", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}

		if userReaction != nil && userReaction.IsLike == like {
			err = app.Service.PostReaction.RemoveReaction(userID, postID)
			if err != nil {
				app.Logger.Error("remove post reaction", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		} else {
			err = app.Service.PostReaction.AddReaction(userID, postID, like)
			if err != nil {
				app.Logger.Error("add post reaction", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		}
	} else if commentIsLike != "" {
		commentID, err := strconv.Atoi(r.PostForm.Get("comment_id"))
		if err != nil || commentID < 1 {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
			return
		}

		reaction := commentIsLike == "true"
		var commentReaction *entities.CommentReaction

		commentReaction, err = app.Service.CommentReaction.GetUserReaction(userID, commentID)
		if err != nil {
			app.Logger.Error("get user comment reaction", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}

		if commentReaction != nil && reaction == commentReaction.IsLike {
			err = app.Service.CommentReaction.RemoveReaction(userID, commentID)
			if err != nil {
				app.Logger.Error("remove comment reaction", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		} else {
			err = app.Service.CommentReaction.AddReaction(userID, commentID, reaction)
			if err != nil {

				app.Logger.Error("add comment reaction", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		}

	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// /comment/reaction
func (app *Application) commentReaction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	commentID, err := strconv.Atoi(r.PostForm.Get("comment_id"))
	if err != nil || commentID < 1 {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	isLike := r.PostForm.Get("is_like")
	reaction := isLike == "true"
	var commentReaction *entities.CommentReaction
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if ok && userID != 0 {
		commentReaction, err = app.Service.CommentReaction.GetUserReaction(userID, commentID)
		if err != nil {
			app.Logger.Error("get user reaction", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}

	if reaction == commentReaction.IsLike {
		err = app.Service.CommentReaction.RemoveReaction(userID, commentID)
		if err != nil {
			app.Logger.Error("remove reaction", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		} else {
			err = app.Service.CommentReaction.AddReaction(userID, commentID, reaction)
			if err != nil {

				app.Logger.Error("Add Reaction", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
