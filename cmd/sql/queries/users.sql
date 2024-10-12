-- name: CreateNormalUser :exec
INSERT INTO users (role_id, first_name, last_name, email, password, subscription_id)
VALUES (2,?, ?, ?, ?, 2);

-- name: CreateSuperUser :exec
INSERT INTO users (role_id, first_name, last_name, email, password, subscription_id)
VALUES (1,?, ?, ?, ?, 1);

-- name: GetUserByID :one
SELECT 
	users.user_id,
    roles.name AS 'role',
    users.first_name,
    users.last_name,
    users.email,
    users.password,
    subscriptions.status AS 'subscription status',
    users.created_at,
    users.updated_at
FROM users users
JOIN roles roles USING(role_id)
JOIN subscriptions subscriptions USING (subscription_id)
WHERE user_id = ?;

-- name: GetUserByEmail :one 
SELECT 
	users.user_id,
    roles.name AS 'role',
    users.first_name,
    users.last_name,
    users.email,
    users.password,
    subscriptions.status AS 'subscription status',
    users.created_at,
    users.updated_at
FROM users users
JOIN roles roles USING(role_id)
JOIN subscriptions subscriptions USING (subscription_id)
WHERE email = ?;

-- name: UpdateUserToSuperUser :exec
UPDATE users SET role_id = 1 WHERE user_id = ?;

-- name: GetUserRoleByUserID :one
SELECT roles.name
FROM users users
JOIN roles roles using(role_id)
WHERE user_id = ?;

-- name: UpdateUserToNormalUser :exec
UPDATE users SET role_id = 2 WHERE user_id = ?;
