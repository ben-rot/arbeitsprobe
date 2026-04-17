-- name: clan-stats-for-dashboard
SELECT 
    c.tag,
    c.name,
    c.cwl_league,
    c.is_locked,
    COUNT(a.id) AS acc_count
FROM clans c
LEFT JOIN accounts a ON c.id = a.clan_id
WHERE c.family_id = $1
GROUP BY c.id
ORDER BY c.cwl_league DESC, acc_count DESC, c.name ASC