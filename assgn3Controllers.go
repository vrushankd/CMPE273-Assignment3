package assgn3Controllers

import (
	"assgn3Models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Pair struct
type Pair struct {
	Key   int
	Value int
}

//PairList struct
type PairList []Pair

//UserController structure
type UserController struct{}

//NewUserController function
func NewUserController() *UserController {
	return &UserController{}
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	return p[i].Value < p[j].Value
}

var (
	iterationCounter = 0
	status           string
	nextDestCount    int
)

//PlanTrip to perform POST operation
func (uc UserController) PlanTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("******************************POST***************************")
	var (
		jsonResp         assgn3Models.Response
		destcordsarr     []string
		tmpBestRouteMap  = make(map[string]string)
		totalcost        = 0
		totalduration    = 0
		totaldistance    = 0.0
		finalcostmap     = make(map[int]int)
		finaldurationmap = make(map[int]int)
		finaldistancemap = make(map[int]float64)
		tmpDestMap       = make(map[int]int)
	)
	//connect to mongodb in cloud using mongolab
	uri := "mongodb://vrushankd:Vrushank90@ds045628.mongolab.com:45628/gomongodb"
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("gomongodb").C("cmpe273Assgn3")

	req := assgn3Models.TripPostReq{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&req)

	strtointcords, _ := strconv.Atoi(req.StartLocationID)
	destcordsarr = req.DestLocationID

	destresultcords := make([]string, len(destcordsarr))

	for i := 0; i < len(destcordsarr); i++ {
		tmpDestMap[i], _ = strconv.Atoi(destcordsarr[i])
	}

	for outer := 0; outer < len(destcordsarr); outer++ {

		var (
			costmap     = make(map[int]int)
			durationmap = make(map[int]int)
			distancemap = make(map[int]float64)
		)

		startcoords := getCoordinates(strtointcords)
		startlatlong := strings.Split(startcoords, ",")

		for it := 0; it < len(destcordsarr); it++ {
			if strtointcords == tmpDestMap[it] || tmpBestRouteMap[destcordsarr[it]] != "" {
				continue
			} else {
				strtointdestcords, _ := strconv.Atoi(destcordsarr[it])
				destcords := getCoordinates(strtointdestcords)
				destlatlong := strings.Split(destcords, ",")
				response, err := http.Get("https://sandbox-api.uber.com/v1/estimates/price?start_latitude=" + startlatlong[0] + "&start_longitude=" + startlatlong[1] + "&end_latitude=" + destlatlong[0] + "&end_longitude=" + destlatlong[1] + "&server_token=SKXhsZoJQI-TYb52bFo8_SeGHhk2B2bwQSHvmM8g")

				if err != nil {
					fmt.Printf("%s", err)
					os.Exit(1)
				} else {
					defer response.Body.Close()
					contents, err := ioutil.ReadAll(response.Body)
					if err != nil {
						fmt.Printf("%s", err)
						os.Exit(1)
					}
					json.Unmarshal([]byte(contents), &jsonResp)

					if strings.Contains(string(contents), "start_longitude") {
						fmt.Println("Start Cordinate: ", strtointcords)
						fmt.Println("Destination Loc: ", strtointdestcords)
						fmt.Println("Distance exceeded 100 miles..Check your cordinates")
						os.Exit(1)
					}
					costmap[strtointdestcords] = jsonResp.Prices[0].LowEstimate
					durationmap[strtointdestcords] = jsonResp.Prices[0].Duration
					distancemap[strtointdestcords] = jsonResp.Prices[0].Distance
				}
			}
		}

		var sortcostmap PairList
		sortcostmap = sortMapByValue(costmap)

		if len(sortcostmap) > 1 {
			if sortcostmap[0].Value < sortcostmap[1].Value {
				destresultcords[outer] = strconv.Itoa(sortcostmap[0].Key)
			} else {
				var sortdurationmap PairList
				var duration1, duration2, locid1, locid2 int
				sortdurationmap = sortMapByValue(durationmap)

				for j := 0; j < len(sortdurationmap); j++ {
					if sortdurationmap[j].Key == sortcostmap[0].Key {
						duration1 = sortdurationmap[j].Value
						locid1 = j
					} else if sortdurationmap[j].Key == sortcostmap[1].Key {
						duration2 = sortdurationmap[j].Value
						locid2 = j
					}
				}
				if duration1 < duration2 {
					destresultcords[outer] = strconv.Itoa(sortdurationmap[locid1].Key)
				} else if duration1 > duration2 {
					destresultcords[outer] = strconv.Itoa(sortdurationmap[locid2].Key)
				} else {
					destresultcords[outer] = strconv.Itoa(sortcostmap[0].Key)
				}
			}
		} else {
			destresultcords[outer] = strconv.Itoa(sortcostmap[0].Key)
		}

		tmpBestRouteMap[strconv.Itoa(sortcostmap[0].Key)] = "tempValue"
		finalcostmap[outer] = sortcostmap[0].Value
		finaldurationmap[outer] = durationmap[sortcostmap[0].Key]
		finaldistancemap[outer] = distancemap[sortcostmap[0].Key]
		strtointcords, _ = strconv.Atoi(destresultcords[outer])
	}

	for _, value := range finalcostmap {
		totalcost = totalcost + value
	}

	for _, value := range finaldurationmap {
		totalduration = totalduration + value
	}

	for _, value := range finaldistancemap {
		totaldistance = totaldistance + value
	}

	dd, _ := strconv.ParseFloat((strconv.FormatFloat(totaldistance, 'f', 3, 64)), 64)

	newID := getNextSequence()

	finalresp := assgn3Models.TripPostGetResp{
		ID:                  newID,
		Status:              "planning",
		StartLocationID:     req.StartLocationID,
		BestRouteLocationID: destresultcords,
		TotalCost:           totalcost,
		TotalDuration:       totalduration,
		TotalDistance:       dd,
	}

	err = collection.Insert(finalresp)
	if err != nil {
		fmt.Printf("Can't insert document: %v\n", err)
		os.Exit(1)
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(finalresp)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

//CheckTrip to perform Get operation
func (uc UserController) CheckTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("******************************GET***************************")
	//connect to mongodb in cloud using mongolab
	uri := "mongodb://vrushankd:Vrushank90@ds045628.mongolab.com:45628/gomongodb"
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("gomongodb").C("cmpe273Assgn3")
	collection2 := sess.DB("gomongodb").C("cmpe273Assgn3PUT")
	collection3 := sess.DB("gomongodb").C("cmpe273Assgn3PutResp")

	id := p.ByName("id")

	u := assgn3Models.TripPostGetResp{}
	v := assgn3Models.NextDestination{}
	put := assgn3Models.Counter{}

	intID, _ := strconv.Atoi(id)

	//Fetch data
	finderr1 := collection2.FindId(intID).One(&put)
	if finderr1 != nil {
		finderr := collection.FindId(intID).One(&u)
		if finderr != nil {
			w.WriteHeader(404)
			return
		}
		// Marshal provided interface into JSON structure
		disp, _ := json.Marshal(u)
		// Write content-type, statuscode, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", disp)
	} else {
		finderr2 := collection3.FindId(intID).One(&v)
		if finderr2 != nil {
			w.WriteHeader(404)
			return
		}
		// Marshal provided interface into JSON structure
		disp, _ := json.Marshal(v)
		// Write content-type, statuscode, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", disp)
	}

}

//CheckNextDestination for PUT operation
func (uc UserController) CheckNextDestination(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("******************************PUT***************************")
	var allcordinatesmap = make(map[int]string)
	var respsandbox assgn3Models.Respsandbox
	//connect to mongodb in cloud using mongolab
	uri := "mongodb://vrushankd:Vrushank90@ds045628.mongolab.com:45628/gomongodb"
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("gomongodb").C("cmpe273Assgn3")
	collection2 := sess.DB("gomongodb").C("cmpe273Assgn3PUT")
	collection3 := sess.DB("gomongodb").C("cmpe273Assgn3PutResp")

	id := p.ByName("id")

	u := assgn3Models.TripPostGetResp{}
	v := assgn3Models.Counter{}

	intID, _ := strconv.Atoi(id)

	//Fetch data
	finderr := collection.FindId(intID).One(&u)
	finderr1 := collection2.FindId(intID).One(&v)

	if finderr != nil {
		fmt.Println("error")
		w.WriteHeader(404)
		return
	}

	if finderr1 != nil {
		insertdoc := assgn3Models.Counter{
			ID:    intID,
			Count: 1,
		}
		err := collection2.Insert(insertdoc)
		iterationCounter = 1
		if err != nil {
			fmt.Printf("Can't insert document: %v\n", err)
			os.Exit(1)
		}
	} else {
		iterationCounter = getCounter(intID)
	}

	routeLocID := u.BestRouteLocationID

	strLoc := u.StartLocationID
	allcordinatesmap[1] = strLoc

	for i := 0; i < len(routeLocID); i++ {
		allcordinatesmap[i+2] = routeLocID[i]
	}

	if iterationCounter < len(allcordinatesmap) {
		nextDestCount = iterationCounter + 1
		status = "requesting"
	} else if iterationCounter == len(allcordinatesmap) {
		nextDestCount = 1
		status = "finished"
	} else {
		fmt.Println("Since you already reached home, this id is killed..End")
		os.Exit(1)
	}

	intstartLoc, _ := strconv.Atoi(allcordinatesmap[iterationCounter])
	strstartcords := getCoordinates(intstartLoc)
	spltstartCords := strings.Split(strstartcords, ",")

	nextDestLocID, _ := strconv.Atoi(allcordinatesmap[nextDestCount])
	nextDestCords := getCoordinates(nextDestLocID)

	spltnextDestCords := strings.Split(nextDestCords, ",")

	var jsonStr = []byte(`{"start_latitude":"` + spltstartCords[0] + `","start_longitude":"` + spltstartCords[1] + `","end_latitude":"` + spltnextDestCords[0] + `","end_longitude":"` + spltnextDestCords[1] + `","product_id":"04a497f5-380d-47f2-bf1b-ad4cfdcb51f2"}`)

	urlreq := "https://sandbox-api.uber.com/v1/requests"
	req, err := http.NewRequest("POST", urlreq, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsiaGlzdG9yeV9saXRlIiwiaGlzdG9yeSIsInByb2ZpbGUiLCJkZWxpdmVyeSIsInJlcXVlc3RfcmVjZWlwdCIsImRlbGl2ZXJ5X3NhbmRib3giLCJyZXF1ZXN0Il0sInN1YiI6IjAxMjhiYjE1LTIxZTgtNDQwMi1hOWU3LTgwNDI1MGVlNjgzOSIsImlzcyI6InViZXItdXMxIiwianRpIjoiOTMzMzM5ZjItZGJkOC00YWRjLThiOTQtZGU2NDY1ODdkOWU4IiwiZXhwIjoxNDQ5NjMxODAzLCJpYXQiOjE0NDcwMzk4MDMsInVhY3QiOiJORU1uS3ZVVkpCWjc1RG1lTmRYb2Yzcm5qQjhmYUsiLCJuYmYiOjE0NDcwMzk3MTMsImF1ZCI6IjBwV0Jpa1pBNE8tVGRzR1dESXQ5V284NXdwV19uelllIn0.lrlPNUa5hIucgSGduNqkXoPUATg_ePKK4iCw8X7ZEFz85HFzFvgvqtDDYiwIlkRPz0bFO0RxJYwm620aA8WOxe4jmweD0j3g7IaenaLpKD5Q8DLqma1C1SrKUc8yehDIYVe4bSYa1Y8luoo-4F56c5prrHuseoIr-asWxmtmASw1GoOQW0Ae7n1sMD-HIXiv2EPlUXN0c3Ir0tqnUNL7f61ptqBm1e9EKUnodFWNy7W0CWiY0aRtCO2LNyuAaDpdK7S_LAQbgjtDVQ_CFrs2qa6TB0s_fJ50IlW9NUfMX4ttV7y8mpULCRbryXTiSxiXv8UnklBgw-5NjOKgay4hDw")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, errmesg := ioutil.ReadAll(resp.Body)

	if errmesg != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	json.Unmarshal([]byte(body), &respsandbox)

	etaWait := respsandbox.Eta
	requestID := respsandbox.ReqID

	putURL := "https://sandbox-api.uber.com/v1/sandbox/requests/" + requestID
	var putmesg = []byte(`{"status":"accepted"}`)
	reqmesg, _ := http.NewRequest("PUT", putURL, bytes.NewBuffer(putmesg))
	reqmesg.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsiaGlzdG9yeV9saXRlIiwiaGlzdG9yeSIsInByb2ZpbGUiLCJkZWxpdmVyeSIsInJlcXVlc3RfcmVjZWlwdCIsImRlbGl2ZXJ5X3NhbmRib3giLCJyZXF1ZXN0Il0sInN1YiI6IjAxMjhiYjE1LTIxZTgtNDQwMi1hOWU3LTgwNDI1MGVlNjgzOSIsImlzcyI6InViZXItdXMxIiwianRpIjoiOTMzMzM5ZjItZGJkOC00YWRjLThiOTQtZGU2NDY1ODdkOWU4IiwiZXhwIjoxNDQ5NjMxODAzLCJpYXQiOjE0NDcwMzk4MDMsInVhY3QiOiJORU1uS3ZVVkpCWjc1RG1lTmRYb2Yzcm5qQjhmYUsiLCJuYmYiOjE0NDcwMzk3MTMsImF1ZCI6IjBwV0Jpa1pBNE8tVGRzR1dESXQ5V284NXdwV19uelllIn0.lrlPNUa5hIucgSGduNqkXoPUATg_ePKK4iCw8X7ZEFz85HFzFvgvqtDDYiwIlkRPz0bFO0RxJYwm620aA8WOxe4jmweD0j3g7IaenaLpKD5Q8DLqma1C1SrKUc8yehDIYVe4bSYa1Y8luoo-4F56c5prrHuseoIr-asWxmtmASw1GoOQW0Ae7n1sMD-HIXiv2EPlUXN0c3Ir0tqnUNL7f61ptqBm1e9EKUnodFWNy7W0CWiY0aRtCO2LNyuAaDpdK7S_LAQbgjtDVQ_CFrs2qa6TB0s_fJ50IlW9NUfMX4ttV7y8mpULCRbryXTiSxiXv8UnklBgw-5NjOKgay4hDw")
	reqmesg.Header.Set("Content-Type", "application/json")

	clientPUT := &http.Client{}
	Resp, err := clientPUT.Do(reqmesg)
	if err != nil {
		panic(err)
	}
	defer Resp.Body.Close()

	responsePut := assgn3Models.NextDestination{}

	disp := assgn3Models.NextDestination{
		ID:                        u.ID,
		Status:                    status,
		StartLocationID:           strLoc,
		NextDestinationLocationID: allcordinatesmap[nextDestCount],
		BestRouteLocationID:       u.BestRouteLocationID,
		TotalCost:                 u.TotalCost,
		TotalDuration:             u.TotalDuration,
		TotalDistance:             u.TotalDistance,
		ETA:                       etaWait,
	}

	finderrput := collection3.FindId(intID).One(&responsePut)

	if finderrput != nil {
		collection3.Insert(disp)
	} else {
		err = collection3.Update(bson.M{"_id": intID}, disp)
		if err != nil {
			fmt.Printf("Can't update document %v\n", err)
			os.Exit(1)
		}
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(disp)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

//getCoordinates to fetch cordinates from mongodb
func getCoordinates(locationID int) string {
	//connect to mongodb in cloud using mongolab
	uri := "mongodb://vrushankd:Vrushank90@ds045628.mongolab.com:45628/gomongodb"
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	u := assgn3Models.CoordinatesStruct{}

	//Fetch data
	finderr := sess.DB("gomongodb").C("cmpe273Assgn2").FindId(locationID).One(&u)

	if finderr != nil {
		fmt.Println("Id not found")
		os.Exit(1)
	}

	lat := strconv.FormatFloat(u.Latitude, 'f', -1, 64)
	return lat + "," + strconv.FormatFloat(u.Longitude, 'f', -1, 64)
}

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[int]int) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

//getNextSequence to track and auto increment the "_id" field each time user performs the POST operation.
func getNextSequence() int {
	var doc assgn3Models.CountID
	sess, err := mgo.Dial("mongodb://vrushankd:Vrushank90@ds045628.mongolab.com:45628/gomongodb")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("gomongodb").C("Assgn3CountId")

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		ReturnNew: true,
	}

	_, err1 := collection.Find(bson.M{"_id": "docid"}).Apply(change, &doc)
	if err1 != nil {
		fmt.Println("got an error finding a doc")
		os.Exit(1)
	}
	return doc.Seq
}

//getCounter to track and auto increment the "_id" field each time user performs the POST operation.
func getCounter(id int) int {
	var doc assgn3Models.Counter
	sess, err := mgo.Dial("mongodb://vrushankd:Vrushank90@ds045628.mongolab.com:45628/gomongodb")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("gomongodb").C("cmpe273Assgn3PUT")

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"counter": 1}},
		ReturnNew: true,
	}

	_, err1 := collection.Find(bson.M{"_id": id}).Apply(change, &doc)
	if err1 != nil {
		fmt.Println("got an error finding a doc")
		os.Exit(1)
	}
	return doc.Count
}
