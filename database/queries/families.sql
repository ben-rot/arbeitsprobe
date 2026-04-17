-- name: family-stats-for-dashboard
SELECT 
    f.discord_server_id,
    f.name,
    f.server_icon,
    COUNT(a.id) AS total_accounts
FROM managers m
INNER JOIN families f ON m.family_id = f.id
LEFT JOIN accounts a ON f.id = a.family_id
WHERE user_id = $1
GROUP BY f.discord_server_id, f.name, f.server_icon