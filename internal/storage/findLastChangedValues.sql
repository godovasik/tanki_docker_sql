-- потом буду это юзать, пока не буду. усложняет логику.
WITH last_values AS (
		SELECT
			(SELECT rank FROM datastamps d WHERE d.user_id = $1 AND rank IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_rank,
			(SELECT kills FROM datastamps d WHERE d.user_id = $1 AND kills IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_kills,
			(SELECT deaths FROM datastamps d WHERE d.user_id = $1 AND deaths IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_deaths,
			(SELECT cry FROM datastamps d WHERE d.user_id = $1 AND cry IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_cry
		FROM datastamps ds
		WHERE ds.user_id = $1
		LIMIT 1
	)
	SELECT * FROM last_values;