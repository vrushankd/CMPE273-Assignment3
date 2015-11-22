# CMPE273-Assignment3
Create RESTful API's for finding the shortest route using Google Maps and UBER APIs in Golang.

There are 3 files:
1. structs.go - Consist of all models(struct) used in the assignemnts.
2. tripserver.go - The main server file which is to be executed.
3. assgn3Controllers.go - Consist of all the controllers (GET, POST and PUT) used in the tripserver.go.

Sample POST call:

POST localhost:3000/trips
request
{
  "starting_from_location_id": "999999",
  "location_ids" : [ "10000", "10001", "20004", "30003" ] 
}

Note: To make this call make sure you have the above id's available in your mongoDB.
