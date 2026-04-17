-- name: get-accounts-owner-by-id
SELECT 
    id,
    tag, 
    name,
    townhall,
    COALESCE(clan_id, 0) AS clan_id,
    COALESCE(family_id, 0) AS family_id
FROM accounts 
WHERE owner_id = $1




-- name: get-accounts-dashboard-by-owner-id
SELECT
    a.id,
    a.name,
    a.tag,
    a.townhall,

    f.name AS family_name,

    c.tag AS clan_tag,
    c.name AS clan_name,
    c.cwl_league AS clan_league
FROM accounts a
LEFT JOIN families f ON a.family_id = f.id
LEFT JOIN clans c ON a.clan_id = c.id
WHERE a.owner_id = $1




-- name: exists
SELECT EXISTS(
    SELECT 
        1
    FROM accounts
    WHERE tag = $1
)




--  name: save
INSERT INTO accounts (
    tag,
    name,
    townhall,
    owner_id
)
VALUES (
    :tag,
    :name,
    :townhall,
    :owner_id
)