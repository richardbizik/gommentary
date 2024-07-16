package handlers

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/richardbizik/gommentary/internal/appctx"
	"github.com/richardbizik/gommentary/internal/database/queries"
)

func (o OApiHandlers) CreateComment(ctx context.Context, request CreateCommentRequestObject) (CreateCommentResponseObject, error) {
	tx, err := o.DB.Tx()
	if err != nil {
		return CreateComment500JSONResponse{
			Code:     "SYSTEM_ERROR",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}
	defer tx.Rollback()

	_, err = o.DB.Queries.WithTx(tx).GetSubjectById(ctx, request.Subject)
	if err == sql.ErrNoRows {
		err := o.DB.Queries.WithTx(tx).CreateSubject(ctx, queries.CreateSubjectParams{
			ID: request.Subject,
			Name: sql.NullString{
				String: Get(request.Body.SubjectName),
				Valid:  request.Body.SubjectName != nil,
			},
		})
		if err != nil {
			return CreateComment400JSONResponse{
				Code:     "ERR_CREATE_SUBJECT",
				Message:  err.Error(),
				Severity: "ERROR",
			}, err
		}
	}

	comment, err := o.DB.Queries.WithTx(tx).CreateComment(ctx, queries.CreateCommentParams{
		ID:      uuid.NewString(),
		Subject: request.Subject,
		Author: sql.NullString{
			String: Get(appctx.GetUsername(ctx)),
			Valid:  appctx.GetUsername(ctx) != nil,
		},
		Text: request.Body.Text,
	})
	if err != nil {
		return CreateComment400JSONResponse{
			Code:     "ERR_CREATE_COMMENT",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}

	comResp := CreateComment200JSONResponse{
		Date: comment.Date,
		Id:   comment.ID,
		Text: comment.Text,
	}
	if comment.Author.Valid {
		comResp.Author = &comment.Author.String
	}

	err = tx.Commit()
	if err != nil {
		return CreateComment500JSONResponse{
			Code:     "TRANSACTION_COMMIT_FAILED",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}

	return comResp, nil
}
