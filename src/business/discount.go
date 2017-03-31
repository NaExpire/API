package business

// func DiscountCreateHandler(writer http.ResponseWriter, request *http.Request) {
// 	x := &businessRegistrationCredentials{}
// 	err := decodeJSON(request.Body, x)
// 	fmt.Printf("Got %s request to DiscountCreateHandler\n", request.Method)
// 	if request.Method != "POST" {
// 		writer.WriteHeader(http.StatusMethodNotAllowed)
// 	} else if err != nil {
// 		io.WriteString(writer, err.Error()+"\n")
// 	} else {
// 		io.WriteString(writer, x.Username+"\n")
// 		io.WriteString(writer, x.Password+"\n")
// 	}
// }

// func DiscountsRetrieveHandler(writer http.ResponseWriter, request *http.Request) {
// 	x := &businessRegistrationCredentials{}
// 	err := decodeJSON(request.Body, x)
// 	fmt.Printf("Got %s request to DiscountCreateHandler\n", request.Method)
// 	if request.Method != "POST" {
// 		writer.WriteHeader(http.StatusMethodNotAllowed)
// 	} else if err != nil {
// 		io.WriteString(writer, err.Error()+"\n")
// 	} else {
// 		io.WriteString(writer, x.Username+"\n")
// 		io.WriteString(writer, x.Password+"\n")
// 	}
// }
