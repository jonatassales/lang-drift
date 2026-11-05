package judge

// routeRefusal is the body the judge route answers a refused request with:
// the code, and the field, the offered models or the cause when the code
// needs one.
type routeRefusal struct {
	Code    string   `json:"code"`
	Field   string   `json:"field"`
	Offered []string `json:"offered"`
	Reason  string   `json:"reason"`
}
