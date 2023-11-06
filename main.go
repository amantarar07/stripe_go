// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/gorilla/mux"
// 	"github.com/stripe/stripe-go/v72"
// 	"github.com/stripe/stripe-go/v72/checkout/session"
// 	"github.com/stripe/stripe-go/v72/customer"
// 	"github.com/stripe/stripe-go/v72/paymentintent"
// 	"github.com/stripe/stripe-go/v72/price"
// 	"github.com/stripe/stripe-go/v72/product"
// 	"github.com/stripe/stripe-go/webhook"
// )

// func init() {
// 	stripe.Key = "sk_test_51NpPR2SFiSFsxEsdUPKl5WEuM7awa4ZweK3CxwJ71axEdsosQVKecXz0EefkVKM3OivbuwDcaf6PB7ytnZwozKWE00SRo2OJcB"
// }

// func CORSCheck(handler http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, req *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Access-Control-Allow-Headers, Origin, Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
// 		w.Header().Set("Access-Control-Allow-Credentials", "true")
// 		w.Header().Set("Access-Control-Allow-Methods", "*")
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Max-Age", "3600")
// 		if req.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}
// 		handler(w, req)
// 	}
// }

// func CheckoutCreator(w http.ResponseWriter, req *http.Request) {

// 	fmt.Println("request body:", req.Body)
// 	input := &EmailInput{}
// 	err := json.NewDecoder(req.Body).Decode(input)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	stripeSession, err := checkout(input.Email)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = json.NewEncoder(w).Encode(&SessionOutput{Id: stripeSession.ID})

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func HandleEvent(w http.ResponseWriter, req *http.Request) {

// 	fmt.Println("inside handle event")
// 	event, err := getEvent(w, req)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("EVENT ", event)

// 	log.Println(event.Type)

// 	if event.Type == "customer.subscription.created" {
// 		c, err := customer.Get(event.Data.Object["customer"].(string), nil)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		email := c.Metadata["FinalEmail"]
// 		log.Println("Subscription created by", email)
// 	}
// }

// func getEvent(w http.ResponseWriter, req *http.Request) (eventRes *stripe.Event, err error) {

// 	fmt.Println("Inside get event ")
// 	const MaxBodyBytes = int64(65536)
// 	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
// 	payload, err := io.ReadAll(req.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	event := stripe.Event{}
// 	err = json.Unmarshal(payload, &event)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &event, nil
// }

// func checkout(email string) (*stripe.CheckoutSession, error) {

// 	// Create a new Product
// 	params := &stripe.ProductParams{
// 		Name: stripe.String("XYZ"),
// 		// Other optional attributes:
// 		// Description: stripe.String("Product Description"),
// 		// Active:      stripe.Bool(true),
// 		// Images:      []string{"https://example.com/product-image.jpg"},
// 		// ...

// 		// Metadata: map[string]string{
// 		//     "key1": "value1",
// 		//     "key2": "value2",
// 		// },
// 	}

// 	newProduct, err := product.New(params)
// 	if err != nil {
// 		log.Fatal(err)

// 	}

// 	fmt.Printf("Created Product ID: %s\n", newProduct.ID)

// 	paramsPrice := &stripe.PriceParams{
// 		Product:    stripe.String(newProduct.ID), // Replace with your product ID
// 		UnitAmount: stripe.Int64(1000),           // Price amount in cents (e.g., $10.00)
// 		Currency:   stripe.String("usd"),         // Currency code
// 		Recurring: &stripe.PriceRecurringParams{ // Optional: Set up recurring pricing if needed
// 			Interval: stripe.String("month"),
// 		},
// 	}

// 	newPrice, err := price.New(paramsPrice)
// 	if err != nil {
// 		log.Fatal(err)

// 	}

// 	fmt.Printf("Created Price ID %s\n", newPrice.ID)

// 	paramsCheckout := &stripe.CheckoutSessionParams{
// 		CustomerEmail: stripe.String(email),
// 		SuccessURL:    stripe.String("https://www.youtube.com/channel/UCzgn3FvGR1UK_0M0B6GiLug"),
// 		CancelURL:     stripe.String("https://www.youtube.com/channel/UCzgn3FvGR1UK_0M0B6GiLug"),
// 		PaymentMethodTypes: stripe.StringSlice([]string{
// 			"card",
// 		}),
// 		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
// 		LineItems: []*stripe.CheckoutSessionLineItemParams{
// 			&stripe.CheckoutSessionLineItemParams{
// 				Price:    stripe.String(newPrice.ID),
// 				Quantity: stripe.Int64(1),
// 			},
// 		},
// 		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
// 			TrialPeriodDays: stripe.Int64(7),
// 			Metadata: map[string]string{
// 				"FinalEmail": email,
// 			},
// 		},
// 	}
// 	return session.New(paramsCheckout)
// }

// type EmailInput struct {
// 	Email string `json:"email"`
// }

// type SessionOutput struct {
// 	Id string `json:"id"`
// }

// func CreateCustomer(w http.ResponseWriter, r *http.Request) {

// 	input := &EmailInput{}
// 	err := json.NewDecoder(r.Body).Decode(&input)
// 	if err != nil {
// 		fmt.Println("error in decoding body: ", err.Error())
// 	}
// 	params := &stripe.CustomerParams{
// 		Email: stripe.String(input.Email),
// 		// You can add more fields here as needed, such as Name, Description, etc.
// 	}
// 	c, err := customer.New(params)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	json.NewEncoder(w).Encode(c.ID)
// }

// func Payment(w http.ResponseWriter, r *http.Request) {
// 	params := &stripe.PaymentIntentParams{
// 		Amount:   stripe.Int64(1000), // Amount in cents (e.g., $10.00)
// 		Currency: stripe.String(string(stripe.CurrencyUSD)),
// 		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
// 			Enabled: stripe.Bool(true),
// 		},
// 	}

// 	pi, err := paymentintent.New(params)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("pi client secret: ", pi.ClientSecret)

// 	// Confirm the Payment Intent
// 	confirmParams := &stripe.PaymentIntentConfirmParams{
// 		PaymentMethod: stripe.String("pm_card_visa"), // Replace with the actual payment method ID or token
// 	}

// 	_, err = paymentintent.Confirm(pi.ID, confirmParams)
// 	if err != nil {
// 		// If the payment requires action (e.g., 3D Secure), handle it here
// 		if stripeErr, ok := err.(*stripe.Error); ok && stripeErr.Code == stripe.ErrorCodePaymentIntentActionRequired {
// 			// You can extract the client secret from the Payment Intent
// 			clientSecret := stripeErr.PaymentIntent.ClientSecret
// 			fmt.Printf("Payment requires action. Client Secret: %s\n", clientSecret)

// 			// You should send the client secret to your frontend to complete the authentication
// 			// In the frontend, you can use Stripe.js to handle the 3D Secure authentication

// 			// Respond to the client with the client secret or other information needed for frontend handling
// 			http.Error(w, "Payment requires action", http.StatusPaymentRequired)
// 			return
// 		}

// 		log.Fatal(err)
// 	}

// 	// Check the payment status
// 	if pi.Status == "succeeded" {
// 		// Payment succeeded, you can perform additional actions here
// 		fmt.Fprint(w,pi.ClientSecret)
// 		fmt.Println("Payment Succeeded!")
// 	} else {
// 		// Payment failed, handle the failure
// 		fmt.Println("Payment Failed!")
// 		fmt.Printf("Payment Error: %s\n", pi.LastPaymentError)
// 	}
// }

// func HandleConfig(w http.ResponseWriter,r *http.Request){

// 	Publishable_Key:=os.Getenv("PUBLISHABLE_KEY")
// 	fmt.Fprint(w,Publishable_Key)

// }

// func WebhookResponse(w http.ResponseWriter, req *http.Request) {

// 	fmt.Println("inside webhook!!!")
// 	const MaxBodyBytes = int64(65536)
// 	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
// 	payload, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
// 		w.WriteHeader(http.StatusServiceUnavailable)
// 		return
// 	}

// 	// This is your Stripe CLI webhook secret for testing your endpoint locally.
// 	endpointSecret := "whsec_mdEK7Dm3uQxzhgleNzhaK5hDD9QhLo30"
// 	// Pass the request body and Stripe-Signature header to ConstructEvent, along
// 	// with the webhook signing key.
// 	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"),
// 		endpointSecret)

// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
// 		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
// 		return
// 	}
// 	if event.Type == "payment_intent.created" {
// 		var paymentintent stripe.PaymentIntent

// 		err:=json.Unmarshal(event.Data.Raw,&paymentintent)
// 		if err!=nil{
// 			fmt.Println("error in unmarshalling event data")
// 		}
// 		// err := json.NewDecoder(event.Data.Raw).Decode(&paymentintent)
// 		fmt.Println("payment intent created!!",paymentintent)
// 	}

// 	// Unmarshal the event data into an appropriate struct depending on its Type
// 	// fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)

// 	switch event.Type {
// 	case "account.updated":
// 		// Then define and call a function to handle the event account.updated
// 	case "account.external_account.created":
// 		// Then define and call a function to handle the event account.external_account.created
// 	case "balance.available":
// 		// Then define and call a function to handle the event balance.available
// 	case "cash_balance.funds_available":
// 		// Then define and call a function to handle the event cash_balance.funds_available
// 	case "customer.created":
// 		// Then define and call a function to handle the event customer.created
// 	case "customer.updated":
// 		// Then define and call a function to handle the event customer.updated
// 	case "invoice.created":
// 		// Then define and call a function to handle the event invoice.created
// 	case "payment_method.attached":
// 		// Then define and call a function to handle the event payment_method.attached
// 	case "payment_method.automatically_updated":
// 		// Then define and call a function to handle the event payment_method.automatically_updated
// 	case "payment_method.detached":
// 		// Then define and call a function to handle the event payment_method.detached
// 	case "payment_method.updated":
// 		// Then define and call a function to handle the event payment_method.updated
// 	// ... handle other event types
// 	default:
// 		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
// 	}

// 	w.WriteHeader(http.StatusOK)

// }

// func main() {

// 	fmt.Println("started session")
// 	r := mux.NewRouter()
// 	r.HandleFunc("/checkout", CORSCheck(CheckoutCreator))
// 	r.HandleFunc("/event", CORSCheck(HandleEvent))
// 	r.HandleFunc("/create-customer", CORSCheck(CreateCustomer))
// 	r.HandleFunc("/payment", CORSCheck(Payment))
// 	r.HandleFunc("/webhook-response", CORSCheck(WebhookResponse))
// 	r.HandleFunc("/config", CORSCheck(HandleConfig))


// 	fmt.Println("server started")
// 	http.ListenAndServe(":8080", r)
// }
package main 



type CarRes struct {
    CarId    string `json:"carId,omitempty"`
    CarName  string `json:"carName,omitempty"`
    Rarity   int64  `json:"rarity,omitempty"`
    Defaults struct {
        Stats         Stats `json:"stats,omitempty"`
        Customization struct {
        } `json:"customization,omitempty"`
        Purchase struct {
            CurrencyType string `json:"currencyType,omitempty"` // 1 for coins and 2 fro cash
            Amount       int64  `json:"amount,omitempty"`
            PremimumBuy  int64  `json:"premimumBuy,omitempty"`
        } `json:"price,omitempty"`
    } `json:"defaults,omitempty"`
    Current struct {
        Stats         Stats `json:"stats,omitempty"`
        Customization struct {
        } `json:"carLooks,omitempty"`
    } `json:"current,omitempty"`
    Status struct {
        Purchasable bool `json:"purchasable"`
        Owned       bool `json:"owned"`
    } `json:"status,omitempty"`
}

type Stats struct {
    Power      int64   "json:\"power,omitempty\""
    Grip       int64   "json:\"grip,omitempty\""
    Weight     int64   "json:\"weight,omitempty\""
    ShiftTime  float64 "json:\"shiftTime,omitempty\""
    OVR        float64 "json:\"or,omitempty\""
    Durability int64   "json:\"durability,omitempty\""
}


