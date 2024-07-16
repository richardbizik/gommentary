package handlers

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/richardbizik/gommentary/internal/database/queries"
)

func (o OApiHandlers) GetCommentHistory(ctx context.Context, request GetCommentHistoryRequestObject) (GetCommentHistoryResponseObject, error) {
	size := 20
	page := 0
	if request.Params.Size != nil {
		size = *request.Params.Size
	}
	if request.Params.Page != nil {
		page = *request.Params.Page
	}

	comment, err := o.DB.Queries.GetComment(ctx, request.Comment)
	if errors.Is(err, sql.ErrNoRows) {
		return GetCommentHistory200JSONResponse{
			Comment:       CommentSimple{},
			Content:       []CommentHistory{},
			First:         true,
			Last:          true,
			Size:          0,
			TotalElements: 0,
			TotalPages:    0,
		}, err
	} else if err != nil {
		return GetCommentHistory400JSONResponse{
			Code:     "NOT_FOUND",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}

	commentResp := CommentSimple{
		Date: comment.Date,
		Id:   comment.ID,
		Text: comment.Text,
	}
	if comment.Author.Valid {
		commentResp.Author = &comment.Author.String
	}

	edits, err := o.DB.Queries.GetEdits(ctx, queries.GetEditsParams{
		ID:      request.Comment,
		PLimit:  int64(size),
		POffset: int64(page * size),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return GetCommentHistory200JSONResponse{
			Comment:       commentResp,
			Content:       []CommentHistory{},
			First:         true,
			Last:          true,
			Size:          0,
			TotalElements: 0,
			TotalPages:    0,
		}, err
	} else if err != nil {
		return GetCommentHistory400JSONResponse{
			Code:     "NOT_FOUND",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}

	totalCount, err := o.DB.Queries.CountEdits(ctx, request.Comment)
	if errors.Is(err, sql.ErrNoRows) {
		return GetCommentHistory404JSONResponse{
			Code:     "NOT_FOUND",
			Message:  "Subject has no comments",
			Severity: "INFO",
		}, err
	}

	commentHistory := make([]CommentHistory, len(edits))
	for i, e := range edits {
		commentHistory[i] = CommentHistory{
			Date: e.Date,
			Text: e.OldText,
		}
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(size)))
	resp := GetCommentHistory200JSONResponse{
		Comment:       commentResp,
		Content:       commentHistory,
		First:         page == 0,
		Last:          page+1 >= totalPages,
		Size:          len(edits),
		TotalElements: totalCount,
		TotalPages:    totalPages,
	}

	return resp, err
}
