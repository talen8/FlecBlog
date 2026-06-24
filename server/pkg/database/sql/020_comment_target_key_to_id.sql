UPDATE comments
SET target_key = a.id::text
FROM articles a
WHERE comments.target_type = 'article'
  AND comments.target_key = a.slug;
