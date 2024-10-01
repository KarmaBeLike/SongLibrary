# Song Library API

## Introduction
The **Song Library API** provides a set of endpoints to manage songs within a music library. The API allows users to perform CRUD (Create, Read, Update, Delete) operations on songs, as well as retrieve paginated results, fetch song lyrics, and more. It also includes integrated Swagger documentation for easier exploration of the available endpoints.


## Installation

### Prerequisites
- Go 1.18 or higher
- A running instance of the SongService (a service layer to manage song data)

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/KarmaBeLike/SongLibrary.git


### START PROJECT
- **install dependencies:**
```
go mod tidy
```
- **run project:**
```
go run ./cmd
```

- **Listing songs data with filtering and pagination:**
    ```http
    GET /api/songs
    ```
     ```http
    GET /api/songs?group=Muse&page=1
    ```
    - queries for filtering:
        - group
        - song
    - queries for pagination:
        - page
        - limit
    - sample output:
    ```json

    ```
- **Listing verses with pagination:**
    ```http
    GET /api/songs/lyrics?id=41&page=1&limit=1
    ```
    - queries for pagination:
        - page
        - limit
    - sample output:
    ```json
    {
    "id": 41,
    "group": "Muse",
    "song": "Supermassive Black Hole",
    "verses": [
        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?"
    ]   
    }

    ```
    
- **Adding new car data**
    ```http
    POST /api/songs
    ```
    - input body:
    ```json
    {
    "id": "1,message":"Song added successfully" 
    }
    ```
- **Update car info:**
    - required parameter: `id`
     ```http
    PATCH /api/songs?id=2
    ```
    - input body:
    ```json
    
     {
    "id": 2,
    "message": "Song updated successfully."

    }
    ```
- **Delete car info:**
    - required parameter: `id`
     ```http
    DELETE /api/songs?id=39
    ```
    {
    "id": 39,
    "message": "Song deleted successfully."
}
---