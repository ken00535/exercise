INSERT INTO shorten.url (
    url, short_url, expire_at
)
SELECT
    md5(random()::text),
    md5(random()::text),
    '2026-12-19 16:39:57'
FROM generate_series(1, 1000000)