# Ezai Scraper API

This is the documentation for Ezai Scraper API.

## Endpoints

### Create Note

**Method**: `POST`

**URL**: `/notes`

**Headers:**

-   `OAI-KEY`: Your OAI key (required)

**Request Body:**

```json
{
    "url": "https://example.com"
}
```

**Response:**

```json
{
    "status": "success",
    "notes": [
        /* list of notes */
    ],
    "title": "Example Title",
    "url": "https://example.com",
    "message": "Successfully scraped notes from the url"
}
```

**Errors:**

-   Missing `OAI-KEY` header

```json
{
    "status": "failure",
    "message": "OAI-KEY header is required"
}
```

-   Problem getting website

```json
{
    "status": "failure",
    "message": "Problem getting website"
}
```

-   Problem getting content

```json
{
    "status": "failure",
    "message": "Problem getting content"
}
```

-   Problem getting notes

```json
{
    "status": "failure",
    "message": "Problem getting notes"
}
```

### Create Summary

**Method**: `POST`

**URL**: `/summary`

**Headers:**

-   `OAI-KEY`: Your OAI key (required)

**Request Body:**

```json
{
    "url": "https://example.com"
}
```

**Response:**

```json
{
    "status": "success",
    "summary": "Generated summary text",
    "title": "Example Title",
    "url": "https://example.com",
    "message": "Summary generated successfully"
}
```

**Errors:**

-   Problem getting website

```json
{
    "status": "failure",
    "message": "Problem getting website"
}
```

-   Missing `OAI-KEY` header

```json
{
    "status": "failure",
    "message": "OAI-KEY header is required"
}
```

-   Problem getting content

```json
{
    "status": "failure",
    "message": "Problem getting content"
}
```

-   Problem getting summary

```json
{
    "status": "failure",
    "message": "Problem getting summary"
}
```
