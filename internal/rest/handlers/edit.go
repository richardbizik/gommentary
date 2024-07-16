package handlers

import (
	"context"

	"github.com/richardbizik/gommentary/internal/appctx"
	"github.com/richardbizik/gommentary/internal/database/queries"
)

func (o OApiHandlers) EditComment(ctx context.Context, request EditCommentRequestObject) (EditCommentResponseObject, error) {
	tx, err := o.DB.Tx()
	if err != nil {
		return EditComment500JSONResponse{
			Code:     "SYSTEM_ERROR",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}
	defer tx.Rollback()

	comment, err := o.DB.Queries.WithTx(tx).GetComment(ctx, request.Comment)
	if err != nil {
		return EditComment404JSONResponse{
			Code:     "NOT_FOUND",
			Message:  err.Error(),
			Severity: "INFO",
		}, err
	}
	if o.ValidateUsers {
		// if the comment is not created by a user return 404
		if comment.Author.String != Get(appctx.GetUsername(ctx)) {
			return EditComment404JSONResponse{
				Code:     "NOT_FOUND",
				Message:  err.Error(),
				Severity: "INFO",
			}, err
		}
	}
	updatedComment, err := o.DB.Queries.WithTx(tx).UpdateComment(ctx, queries.UpdateCommentParams{
		NewText: request.Body.Text,
		ID:      comment.ID,
	})
	if err != nil {
		return EditComment404JSONResponse{
			Code:     "UPDATE_FAILED",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}
	_, err = o.DB.Queries.WithTx(tx).CreateEdit(ctx, queries.CreateEditParams{
		ID:   comment.ID,
		Text: comment.Text,
	})
	if err != nil {
		return EditComment400JSONResponse{
			Code:     "UPDATE_FAILED",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}
	err = tx.Commit()
	if err != nil {
		return EditComment500JSONResponse{
			Code:     "TRANSACTION_COMMIT_FAILED",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}

	resp := EditComment200JSONResponse{
		Date: updatedComment.Date,
		Id:   updatedComment.ID,
		Text: updatedComment.Text,
	}
	if updatedComment.Author.Valid {
		resp.Author = &updatedComment.Author.String
	}
	return resp, nil
}
