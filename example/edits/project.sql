SELECT is_part_of, SUM(edits) as number_of_edits FROM articles GROUP BY is_part_of ORDER BY number_of_edits DESC;
