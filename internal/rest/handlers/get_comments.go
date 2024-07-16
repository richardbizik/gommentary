package handlers

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/richardbizik/gommentary/internal/database/queries"
)

func (o OApiHandlers) GetCommentsPage(ctx context.Context, request GetCommentsPageRequestObject) (GetCommentsPageResponseObject, error) {

	size := 20
	page := 0
	if request.Params.Size != nil {
		size = *request.Params.Size
	}
	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	comments, err := o.DB.Queries.GetComments(ctx, queries.GetCommentsParams{
		Subject: request.Subject,
		PLimit:  int64(size),
		POffset: int64(page * size),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return GetCommentsPage200JSONResponse{
			Content: []Comment{},
			First:   true,
			Last:    true,
			Size:    0,
			Subject: Subject{
				Id: request.Subject,
			},
			TotalElements: 0,
			TotalPages:    0,
		}, err
	}
	if err != nil {
		return GetCommentsPage400JSONResponse{
			Code:     "ERR_GET_COMMENTS",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}

	commentsTotal, err := o.DB.Queries.CountComments(ctx, request.Subject)
	if err != nil {
		return GetCommentsPage400JSONResponse{
			Code:     "ERR_COUNT_COMMENTS",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}

	subject, err := o.DB.Queries.GetSubjectById(ctx, request.Subject)
	if err != nil {
		return GetCommentsPage400JSONResponse{
			Code:     "ERR_GET_SUBJECT",
			Message:  err.Error(),
			Severity: "ERROR",
		}, err
	}

	contentResp := make([]Comment, len(comments))
	for i, c := range comments {
		contentResp[i] = Comment{
			Date:     c.Date,
			Id:       c.ID,
			Position: int(c.Position),
			Text:     c.Text,
			Replies:  int(c.Replies),
			Edits:    int(c.Edits),
		}
		if c.Author.Valid {
			contentResp[i].Author = &c.Author.String
		}
	}

	totalPages := int(math.Ceil(float64(commentsTotal) / float64(size)))
	resp := GetCommentsPage200JSONResponse{
		Content: contentResp,
		First:   page == 0,
		Last:    page+1 >= totalPages,
		Size:    len(comments),
		Subject: Subject{
			Id: subject.ID,
		},
		TotalElements: commentsTotal,
		TotalPages:    totalPages,
	}
	if subject.Name.Valid {
		resp.Subject.Name = &subject.Name.String
	}
	return resp, nil
}
