SELECT *,
    (
        SELECT TOTAL(reaction)
        FROM posts_reactions
        WHERE post_id_fkey = p.id
    ) AS rating,
    IFNULL (
        (
            SELECT reaction
            FROM posts_reactions
            WHERE user_id_fkey = $1
                AND post_id_fkey = p.id
        ),
        0
    ) AS yor_reaction,
    FROM comments
WHERE post_id_fkey = ?
ORDER BY created DESC