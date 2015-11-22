package assgn3Models

//TripPostReq to accept input request for POST operation
type TripPostReq struct {
	StartLocationID string   `json:"starting_from_location_id"`
	DestLocationID  []string `json:"location_ids"`
}

//CountID structure to keep the track of "_id"
type CountID struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

//Counter structure to keep track of PUT pointer
type Counter struct {
	ID    int `bson:"_id"`
	Count int `bson:"counter"`
}

//TripPostGetResp to receive the response for POST operation
type TripPostGetResp struct {
	ID                  int      `json:"id" bson:"_id"`
	Status              string   `json:"status" bson:"status"`
	StartLocationID     string   `json:"starting_from_location_id" bson:"starting_from_location_id"`
	BestRouteLocationID []string `json:"best_route_location_ids" bson:"best_route_location_ids"`
	TotalCost           int      `json:"total_uber_costs" bson:"total_uber_costs"`
	TotalDuration       int      `json:"total_uber_duration" bson:"total_uber_duration"`
	TotalDistance       float64  `json:"total_distance" bson:"total_distance"`
}

//NextDestination struct
type NextDestination struct {
	ID                        int      `json:"id" bson:"_id"`
	Status                    string   `json:"status" bson:"status"`
	StartLocationID           string   `json:"starting_from_location_id" bson:"starting_from_location_id"`
	NextDestinationLocationID string   `json:"next_destination_location_id" bson:"next_destination_location_id"`
	BestRouteLocationID       []string `json:"best_route_location_ids" bson:"best_route_location_ids"`
	TotalCost                 int      `json:"total_uber_costs" bson:"total_uber_costs"`
	TotalDuration             int      `json:"total_uber_duration" bson:"total_uber_duration"`
	TotalDistance             float64  `json:"total_distance" bson:"total_distance"`
	ETA                       int      `json:"uber_wait_time_eta" bson:"uber_wait_time_ets"`
}

//POSTReq struct
type POSTReq struct {
	StartLatitude  string `json:"start_latitude"`
	StartLongitude string `json:"start_longitude"`
	EndLatitude    string `json:"end_latitude"`
	EndLongitude   string `json:"end_longitude"`
	ProdID         string `json:"product_id"`
}

//CoordinatesStruct structure to store data in mongodb and json structure for displaying in POSTMAN
type CoordinatesStruct struct {
	ID          int    `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Address     string `bson:"address" json:"address"`
	City        string `bson:"city" json:"city"`
	State       string `bson:"state" json:"state"`
	Zip         string `bson:"zip" json:"zip"`
	Coordinates `json:"coordinate"`
}

//Coordinates struct
type Coordinates struct {
	Latitude  float64 `bson:"lat" json:"lat"`
	Longitude float64 `bson:"long" json:"long"`
}

//Respsandbox struct
type Respsandbox struct {
	Status   string `json:"status"`
	ReqID    string `json:"request_id"`
	Driver   string `json:"driver"`
	Eta      int    `json:"eta"`
	Location string `json:"location"`
	Vehicle  string `json:"vehicle"`
	Surge    int    `json:"surge_multiplier"`
}

//Response struct
type Response struct {
	Prices []struct {
		CurrencyCode         string  `json:"currency_code"`
		DisplayName          string  `json:"display_name"`
		Distance             float64 `json:"distance"`
		Duration             int     `json:"duration"`
		Estimate             string  `json:"estimate"`
		HighEstimate         int     `json:"high_estimate"`
		LocalizedDisplayName string  `json:"localized_display_name"`
		LowEstimate          int     `json:"low_estimate"`
		Minimum              int     `json:"minimum"`
		ProductID            string  `json:"product_id"`
		SurgeMultiplier      int     `json:"surge_multiplier"`
	} `json:"prices"`
}

//ProductResp struct
type ProductResp struct {
	Products []struct {
		Capacity     int    `json:"capacity"`
		Description  string `json:"description"`
		DisplayName  string `json:"display_name"`
		Image        string `json:"image"`
		PriceDetails struct {
			Base            float64 `json:"base"`
			CancellationFee int     `json:"cancellation_fee"`
			CostPerDistance float64 `json:"cost_per_distance"`
			CostPerMinute   float64 `json:"cost_per_minute"`
			CurrencyCode    string  `json:"currency_code"`
			DistanceUnit    string  `json:"distance_unit"`
			Minimum         float64 `json:"minimum"`
			ServiceFees     []struct {
				Fee  float64 `json:"fee"`
				Name string  `json:"name"`
			} `json:"service_fees"`
		} `json:"price_details"`
		ProductID string `json:"product_id"`
	} `json:"products"`
}
