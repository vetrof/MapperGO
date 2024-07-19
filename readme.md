# Nearby Attractions Finder

**Nearby Attractions Finder** is a project designed to help tourists find the closest points of interest based on their current location. By leveraging geographical data and efficient querying, this project provides users with a convenient way to discover nearby attractions.

## Features

- **Find Nearby Places**: Given a set of coordinates, the system retrieves the nearest attractions sorted by distance.
- **Distance Calculation**: Utilizes geographical functions to calculate the distance between the user's location and potential attractions.
- **Efficient Queries**: Optimized SQL queries to ensure quick response times even with large datasets.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (1.16 or later)
- [PostgreSQL](https://www.postgresql.org/) with PostGIS extension
- [Git](https://git-scm.com/)


### API Endpoints

- **Find Nearby Places**: `/nearplaces`
    - **Method**: POST
    - **Request Body**:
        ```json
        {
            "lat": 51.5074,
            "lng": -0.1278
        }
        ```
    - **Response**:
        ```json
        [
            {
                "id": 1,
                "name": "Attraction Name",
                "geom": "POINT(51.5074 -0.1278)",
                "distance": 100.0
            }
        ]
        ```



### API Endpoints

- Add addres field
- 