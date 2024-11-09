# Receipt Processor

A web service that processes receipts and calculates points based on specific rules.

## Prerequisites
- Go 1.23
- Docker

---

## Running the Application

Start by cloning the repository.

### Using Docker (Recommended)
1. Build and run the application:
   ```bash
   make docker/run
   ```
   The application will be available at http://localhost:8080
2. Stop the application:
   ```bash
   make docker/stop
   ```

### Running Locally
1. Jump into the server directory and download the required modules:
   ```bash
   cd server
   go mod download
   ``` 
2. Build and run the application:
   ```bash
   make build
   ```
3. Run the application:
   ```bash
   make run
   ```
4. Run the application with live reloading on file changes:
   ```bash
   make watch
   ```
   
---

## API Endpoints

### Process Receipt
- **POST** `/receipts/process`
- Request body:
   ```json
   {
     "retailer": "Target",
     "purchaseDate": "2022-01-01",
     "purchaseTime": "13:01",
     "items": [
       {
         "shortDescription": "Mountain Dew 12PK",
         "price": "6.49"
       },{
         "shortDescription": "Emils Cheese Pizza",
         "price": "12.25"
       },{
         "shortDescription": "Knorr Creamy Chicken",
         "price": "1.26"
       },{
         "shortDescription": "Doritos Nacho Cheese",
         "price": "3.35"
       },{
         "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
         "price": "12.00"
       }
     ],
     "total": "35.35"
   }
   ```
- Response:
   ```json
   { "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
   ```

### Get Points
- **GET** `/receipts/{id}/points`
- Returns points for a processed receipt
- Response:
   ```json
   { "points": 28 }
   ```

---

