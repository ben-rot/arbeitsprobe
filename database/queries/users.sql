-- name: upsert-user
INSERT INTO users (
    discord_id,
    username,
    nickname,
    avatar_hash
)
VALUES (
    :discord_id,
    :username,
    :nickname,
    :avatar_hash
)
ON CONFLICT (discord_id)
DO UPDATE SET
    username = EXCLUDED.username,
    nickname = EXCLUDED.nickname,
    avatar_hash = EXCLUDED.avatar_hash
RETURNING id




-- name: user-by-id
SELECT 
    discord_id,
    username,
    nickname,
    avatar_hash
FROM users
WHERE id = $1




-- name: check-manager
SELECT 
    family_id
FROM managers
WHERE user_id = $1