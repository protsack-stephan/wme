SELECT in_language, SUM(edits) as number_of_edits FROM articles GROUP BY in_language ORDER BY number_of_edits DESC;
