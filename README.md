<p align="center">
  <a href="https://github.com/PeteGabriel/yamda_go/actions?query=workflow%3ATest">
    <img src="https://img.shields.io/github/workflow/status/PeteGabriel/yamda_go/Test?label=%F0%9F%A7%AA%20tests&style=flat&color=75C46B">
  </a>    
</p>


## Yet Another Movie Database API 


## Overview of API

| Method | URL Pattern               | Action                                          |
|--------|---------------------------|-------------------------------------------------|
| GET    | /v1/healthcheck           | Show application health and version information |
| GET    | /v1/movies                | Show the details of all movies                  |
| POST   | /v1/movies                | Create a new movie                              |
| GET    | /v1/movies/:id            | Show the details of a specific movie            |
| PATCH  | /v1/movies/:id            | Update the details of a specific movie          |
| DELETE | /v1/movies/:id            | Delete a specific movie                         |
| POST   | /v1/users                 | Register a new user                             |
| PUT    | /v1/users/activated       | Activate a specific user                        |
| PUT    | /v1/users/password        | Update the password for a specific user         |
| POST   | /v1/tokens/authentication | Generate a new authentication token             |
| POST   | /v1/tokens/password-reset | Generate a new password-reset token             |
| GET    | /debug/vars               | Display application metrics                     |


### Error response

This API tries to make use of a standardized mediatype called `application/problem+json`. You should expect this for all
the errors in the range 400-4xx.

HTTP/1.1 401 Unauthorized
Content-Type: application/problem+json; charset=utf-8
Date: Wed, 07 Aug 2019 10:10:06 GMT
```json
{
    "type": "https://example.com/v1/movies/78",
    "title": "Not authorized to view movie details",
    "status": 401,
    "detail": "Due to privacy concerns you are not allowed to view account details of others.",
    "instance": "/error/123456/details"
}
```
