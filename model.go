package main

type User struct {
	ID          int                 `json:"id"`
	Level       int                 `json:"level"`
	Username    string              `json:"username"`
	FirstName   string              `json:"firstname"`
	LastName    string              `json:"lastname"`
	Email       string              `json:"email"`
	PhoneNumber string              `json:"phonenumber"`
	Password    string              `json:"password"`
	Requests    []RestaurantRequest `json:"requests"`
	Restaurant  BusinessResponse    `json:"restaurant"`
	Token       string              `json:"token"`
}
type RestaurantRequest struct {
	ID             int    `json:"id"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	RestaurantName string `json:"restaurantname"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phonenumber"`
	City           string `json:"city"`
	State          string `json:"state"`
	ZipCode        string `json:"zipcode"`
	Country        string `json:"country"`
	Active         bool   `json:"active"`
	Enabled        bool   `json:"enabled"`
}
type BusinessResponse struct {
	Categories  []Category `json:"categories"`
	ID          string     `json:"id"`
	ImageURL    string     `json:"image_url"`
	Phone       string     `json:"phone"`
	Price       string     `json:"price"`
	Rating      float32    `json:"rating"`
	Name        string     `json:"name"`
	Location    Location   `json:"location"`
	ReviewCount int32      `json:"review_count"`
	Photos      []string   `json:"photos"`
	Hours       []Hour     `json:"hours"`
	URL         string     `json:"url"`
}
type Hour struct {
	HoursType string `json:"hours_type"`
	IsOpenNow bool   `json:"is_open_now"`
	Open      []Open `json:"open"`
}
type Open struct {
	End         string `json:"end"`
	IsOvernight bool   `json:"is_overnight"`
	Day         int    `json:"day"`
	Start       string `json:"start"`
}
type Category struct {
	Title string `json:"title"`
	Alias string `json:"alias"`
}
type Location struct {
	DisplayAddress []string `json:"display_address"`
	Address1       string   `json:"address1"`
	Address2       string   `json:"address2"`
	Address3       string   `json:"address3"`
	City           string   `json:"city"`
	State          string   `json:"state"`
	ZipCode        string   `json:"zip_code"`
	Country        string   `json:"country"`
}
type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}
type ID struct {
	Value string `json:"value"`
}
