# Shorten

Design and implement (with unit tests) an URL shortener using Go programming language.

## Technical Selection

**postgres**

A RBDMS guarantees ACID, so I select RDBMS as data persistence. MySQL and Postgres are active, distributed-supported open source sofeware. I choose Postgres because I am more familiar with it.

**gorm**

I usually use ORM to handle simple model application, because it can avoid sql injection and connection pool issues. shorten url has a very simple model. I think it can work good.

If want to improve perf, pgx is a option [[ref](https://github.com/efectn/go-orm-benchmarks/blob/master/results.md)]

**gin**

I used two web frameworks before, gin and echo. gin has a more active repository [[ref](https://pkg.go.dev/github.com/mingrammer/go-web-framework-stars#section-readme)]. Also, its api is more developer friendly, just like bind parameter

```go
if err := c.ShouldBindJSON(&req); err != nil {
    _ = c.Error(errors.Wrap(entity.ErrInvalidInput, err.Error()))
    return
}
```

and echo need to implement validator, so I prefer gin. I think it's enough for this application.

**fx**

When developing a application, I am used to split codebase into app and infra. It can bring more flexibility but has more complex deps at the same time. So I use dep injection framework to manage it.

fx can support many scenarios [[ref1](https://medium.com/@ken00535/%E7%94%A8-fx-%E4%BE%86%E6%9B%BF-go-%E4%BE%9D%E8%B3%B4%E6%B3%A8%E5%85%A5%E5%90%A7-d82adcd4d56b)][[ref2](https://speakerdeck.com/ken00535/20220928-golang-meetup-di-fx-release?slide=2)]. I think wire has same effect, but I am not familiar with wire. So I choose fx. 

**zerolog**

It's very import to pick a good logger. It need to have some feature

- high performance (you don't want to waste cpu/mem resources on non-biz place) [[ref](https://github.com/rs/zerolog#benchmarks)]
- field based (easy to integrate with obervability, like ELK)

so zerolog is a good choice, zap is another option.

## Getting start

### Prerequisite



## Requirement

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

## Benchmark

Query when 1M records

| Package             |       Time      |
| :------------------ | :-------------: |
| db without index    | 63010739 ns/op  |
| db with index       | 3486737 ns/op   |