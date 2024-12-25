-- name: GetEmailTemplateByID :one
SELECT *
FROM email_template
WHERE id = ?;
-- name: GetEmailTemplateByEventID :one
SELECT *
FROM email_template
WHERE event_id = ?;
-- name: CreateEmailTemplate :exec
INSERT INTO email_template (
        event_id,
        header_image,
        content,
        footer_image,
        subject,
        bg_color,
        message
    )
VALUES (?, ?, ?, ?, ?, ?, ?);
-- name: UpdateEmailTemplateByID :exec
UPDATE email_template
SET event_id = ?,
    header_image = ?,
    content = ?,
    footer_image = ?,
    subject = ?,
    bg_color = ?,
    message = ?
WHERE id = ?;