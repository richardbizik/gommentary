-- name: GetSubjectById :one
select * from subject where id = @id;

-- name: CreateSubject :exec
insert into subject (id, name, created) values(@id, @name, CURRENT_TIMESTAMP);
