# Assignment

Design and implement (with unit tests) an URL shortener using Go programming language.

## Getting start



## Criteria

1. URL shortener has 2 APIs, please follow API example to implement:
    1. A RESTful API to upload a URL with its expired date and response with a shorten URL.
    2. An API to serve shorten URLs responded by upload API, and redirect to original URL. If URL is expired, please response with status 404.
2. Please feel free to use any external libs if needed.
3. It is also free to use following external storage including:
    1. Relational database (MySQL, PostgreSQL, SQLite)
    2. Cache storage (Redis, Memcached)
4. Please implement reasonable constrains and error handling of these 2 APIs.
5. You do not need to consider auth.
6. Many clients might access shorten URL simultaneously or try to access with non-existent shorten URL, please take
performance into account.

## API Example

```bash
# Upload URL API
curl -X POST -H "Content-Type:application/json" http://localhost/api/v1/urls -d '{
    "url": "<original_url>",
    "expireAt": "2021-02-08T09:20:41Z"
}'

# Response
{
    "id": "<url_id>",
    "shortUrl": "http://localhost/<url_id>"
}

# ------------------

# Redirect URL API
curl -L -X GET http://localhost/<url_id> => REDIRECT to original URL
```