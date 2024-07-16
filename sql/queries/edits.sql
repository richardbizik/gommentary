-- name: CreateEdit :one
INSERT INTO comment_edit (comment, date, old_text)
    VALUES (@id, CURRENT_TIMESTAMP, @text)
RETURNING
    *;

-- name: GetEdits :many
SELECT
    *
FROM
    comment_edit
WHERE
    comment = @id
ORDER BY
    date DESC
LIMIT @pLimit OFFSET @pOffset;

-- name: CountEdits :one
SELECT
    count(*)
FROM
    comment_edit
WHERE
    comment = @id;
