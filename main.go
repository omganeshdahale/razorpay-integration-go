package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/razorpay/razorpay-go"
	"github.com/razorpay/razorpay-go/utils"
)

type Response struct {
	Message string `json:"message"`
}

type WebhookBody struct {
	AccountID string   `json:"account_id"`
	Contains  []string `json:"contains"`
	CreatedAt int      `json:"created_at"`
	Entity    string   `json:"entity"`
	Event     string   `json:"event"`
	Payload   struct {
		Order struct {
			Entity struct {
				Amount     int    `json:"amount"`
				AmountDue  int    `json:"amount_due"`
				AmountPaid int    `json:"amount_paid"`
				Attempts   int    `json:"attempts"`
				CreatedAt  int    `json:"created_at"`
				Currency   string `json:"currency"`
				Entity     string `json:"entity"`
				ID         string `json:"id"`
				Notes      []any  `json:"notes"`
				OfferID    any    `json:"offer_id"`
				Receipt    any    `json:"receipt"`
				Status     string `json:"status"`
			} `json:"entity"`
		} `json:"order"`
		Payment struct {
			Entity struct {
				AcquirerData struct {
					Rrn              string `json:"rrn"`
					UpiTransactionID string `json:"upi_transaction_id"`
				} `json:"acquirer_data"`
				Amount            int    `json:"amount"`
				AmountRefunded    int    `json:"amount_refunded"`
				AmountTransferred int    `json:"amount_transferred"`
				Bank              any    `json:"bank"`
				BaseAmount        int    `json:"base_amount"`
				Captured          bool   `json:"captured"`
				Card              any    `json:"card"`
				CardID            any    `json:"card_id"`
				Contact           string `json:"contact"`
				CreatedAt         int    `json:"created_at"`
				Currency          string `json:"currency"`
				Description       string `json:"description"`
				Email             string `json:"email"`
				Entity            string `json:"entity"`
				ErrorCode         any    `json:"error_code"`
				ErrorDescription  any    `json:"error_description"`
				ErrorReason       any    `json:"error_reason"`
				ErrorSource       any    `json:"error_source"`
				ErrorStep         any    `json:"error_step"`
				Fee               int    `json:"fee"`
				FeeBearer         string `json:"fee_bearer"`
				ID                string `json:"id"`
				International     bool   `json:"international"`
				InvoiceID         any    `json:"invoice_id"`
				Method            string `json:"method"`
				Notes             []any  `json:"notes"`
				OrderID           string `json:"order_id"`
				RefundStatus      any    `json:"refund_status"`
				Status            string `json:"status"`
				Tax               int    `json:"tax"`
				Vpa               string `json:"vpa"`
				Wallet            any    `json:"wallet"`
			} `json:"entity"`
		} `json:"payment"`
		PaymentLink struct {
			Entity struct {
				AcceptPartial  bool   `json:"accept_partial"`
				Amount         int    `json:"amount"`
				AmountPaid     int    `json:"amount_paid"`
				CallbackMethod string `json:"callback_method"`
				CallbackURL    string `json:"callback_url"`
				CancelledAt    int    `json:"cancelled_at"`
				CreatedAt      int    `json:"created_at"`
				Currency       string `json:"currency"`
				Customer       struct {
				} `json:"customer"`
				Description           string `json:"description"`
				ExpireBy              int    `json:"expire_by"`
				ExpiredAt             int    `json:"expired_at"`
				FirstMinPartialAmount int    `json:"first_min_partial_amount"`
				ID                    string `json:"id"`
				Notes                 any    `json:"notes"`
				Notify                struct {
					Email    bool `json:"email"`
					Sms      bool `json:"sms"`
					Whatsapp bool `json:"whatsapp"`
				} `json:"notify"`
				OrderID        string `json:"order_id"`
				ReferenceID    string `json:"reference_id"`
				ReminderEnable bool   `json:"reminder_enable"`
				Reminders      struct {
				} `json:"reminders"`
				ShortURL  string `json:"short_url"`
				Status    string `json:"status"`
				UpdatedAt int    `json:"updated_at"`
				UpiLink   bool   `json:"upi_link"`
				UserID    string `json:"user_id"`
			} `json:"entity"`
		} `json:"payment_link"`
	} `json:"payload"`
}

func main() {
	client := razorpay.NewClient("rzp_test_nihdIl0zFWL4yM", "nsFpwAm0jCt09HmK1qH3fJrx")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Create Razorpay payment link
		linkParams := map[string]interface{}{
			"amount":          100000,
			"currency":        "INR",
			"callback_url":    "http://localhost:8080/success",
			"callback_method": "get",
		}
		paymentLink, err := client.PaymentLink.Create(linkParams, nil)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		// Return the payment link in the API response
		response := Response{Message: paymentLink["short_url"].(string)}
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		response := Response{Message: "Payment successful"}
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		// Verify the signature
		sign := r.Header.Get("X-Razorpay-Signature")
		secret := "F2Tq$NvckR4#YsT"
		result := utils.VerifyWebhookSignature(string(body), sign, secret)

		if result {
			bodyReader := bytes.NewReader(body)
			var wbody WebhookBody
			json.NewDecoder(bodyReader).Decode(&wbody)
			fmt.Printf("Payment successful for payment link: %s\n", wbody.Payload.PaymentLink.Entity.ID)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
