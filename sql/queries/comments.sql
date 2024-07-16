-- name: CreateComment :one
INSERT INTO comment (id, subject, author, text, date, position)
    VALUES (@id, @subject, @author, @text, CURRENT_TIMESTAMP, (
            SELECT
                coalesce(max(position) + 1, 1)
            FROM
                comment
            WHERE
                subject = @subject))
RETURNING
    *;

-- name: GetComments :many
SELECT
    c.*,
    (
        SELECT
            count(*)
        FROM
            comment AS i
        WHERE
            i.subject = c.id) AS replies,
    (
        SELECT
            count(*)
        FROM
            comment_edit
        WHERE
            comment_edit.comment = c.id) AS edits
FROM
    comment AS c
WHERE
    c.subject = @subject
ORDER BY
    position ASC
LIMIT @pLimit OFFSET @pOffset;

-- name: CountComments :one
SELECT
    count(*)
FROM
    comment
WHERE
    subject = @subject;

-- name: GetComment :one
SELECT
    *
FROM
    comment
WHERE
    id = @id;

-- name: UpdateComment :one
UPDATE
    comment
SET
    text = @new_text
WHERE
    id = @id
RETURNING
    *;
