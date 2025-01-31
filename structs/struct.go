package structs

import (
	"time"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	Contact   string `json:"contact"`
	Is2FA		  bool   `json:"is_2fa"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Disaster struct {
	ID          int    		`json:"id"`
	Type        string 		`json:"type"`
	Location    string 		`json:"location"`
	Description string 		`json:"description"`
	Status      string 		`json:"status"`
	ReportedBy  int    		`json:"reported_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Shelter struct {
	ID                int    		`json:"id"`
	Name              string 		`json:"name"`
	Location          string 		`json:"location"`
	CapacityTotal     int    		`json:"capacity_total"`
	CapacityRemaining int    		`json:"capacity_remaining"`
	EmergencyNeeds    string 		`json:"emergency_needs"`
	DisasterID        *int    	`json:"disaster_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Refugee struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Condition  string `json:"condition"`
	Needs      string `json:"needs"`
	ShelterID  *int   `json:"shelter_id,omitempty"`
	DisasterID *int   `json:"disaster_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Logistic struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Quantity   int    `json:"quantity"`
	Status     string `json:"status"`
	DisasterID *int   `json:"disaster_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DistributionLog struct {
	ID            int    `json:"id"`
	LogisticID    *int   `json:"logistic_id,omitempty"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	Distance      float64 `json:"distance"`
	SenderName    string `json:"sender_name"`
	RecipientName string `json:"recipient_name"`
	QuantitySent  int    `json:"quantity_sent"`
	SentAt        string `json:"sent_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type EvacuationRoute struct {
	ID            int    		`json:"id"`
	DisasterID    *int   		`json:"disaster_id,omitempty"`
	Origin        string 		`json:"origin"`
	Destination   string 		`json:"destination"`
	Distance      float64 	`json:"distance"`
	Route         string 		`json:"route"`
	Status        string 		`json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type EmergencyReport struct {
	ID          int    `json:"id"`
	UserID      *int   `json:"user_id,omitempty"`
	DisasterID  *int   `json:"disaster_id,omitempty"`
	Description string `json:"description"`
	Location    string `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Donation struct {
	ID         int     `json:"id"`
	DonorID		 *int    `json:"user_id,omitempty"`
	DisasterID *int    `json:"disaster_id,omitempty"`
	Amount     float64 `json:"amount"`
	ItemName	 string  `json:"item_name"`
	Status     string  `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Volunteer struct {
	ID         int    `json:"id"`
	DonorID    *int   `json:"donor_id,omitempty"`
	DisasterID *int   `json:"disaster_id,omitempty"`
	Skill      string `json:"skill"`
	Location   string `json:"location"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type Login struct {
	Email		 string `json:"email"`
	Password string `json:"password"`
}

type VerifyOTP struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type Enable2FA struct {
	Is2FA bool `json:"is_2fa"`
}

type UserInput struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	Contact   string `json:"contact"`
}

type UpdateUserInfoWithoutEmail struct {
	Name    string `json:"name"`
	Role    string `json:"role"`
	Contact string `json:"contact"`
}

type ChangeUserRole struct {
	Role string `json:"role"`
}

type DisasterInput struct {
	Type        string `json:"type"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ReportedBy  int    `json:"reported_by"`
}

type ShelterInput struct {
	Name              string `json:"name"`
	Location          string `json:"location"`
	CapacityTotal     int    `json:"capacity_total"`
	CapacityRemaining int    `json:"capacity_remaining"`
	EmergencyNeeds    string `json:"emergency_needs"`
	DisasterID        *int   `json:"disaster_id,omitempty"`
}

type RefugeeInput struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Condition  string `json:"condition"`
	Needs      string `json:"needs"`
	ShelterID  *int   `json:"shelter_id,omitempty"`
	DisasterID *int   `json:"disaster_id,omitempty"`
}

type LogisticInput struct {
	Type       string `json:"type"`
	Quantity   int    `json:"quantity"`
	Status     string `json:"status"`
	DisasterID *int   `json:"disaster_id,omitempty"`
}

type DistributionLogInput struct {
	LogisticID    *int   `json:"logistic_id,omitempty"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	Distance      float64 `json:"distance"`
	SenderName    string `json:"sender_name"`
	RecipientName string `json:"recipient_name"`
	QuantitySent  int    `json:"quantity_sent"`
	SentAt        string `json:"sent_at"`
}

type EvacuationRouteInput struct {
	DisasterID    *int   `json:"disaster_id,omitempty"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	Distance      float64 `json:"distance"`
	Route         string `json:"route"`
	Status        string `json:"status"`
}

type EmergencyReportInput struct {
	UserID      *int   `json:"user_id,omitempty"`
	DisasterID  *int   `json:"disaster_id,omitempty"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

type DonationInput struct {
	DonorID    *int    `json:"user_id,omitempty"`
	DisasterID *int    `json:"disaster_id,omitempty"`
	Amount     float64 `json:"amount"`
	ItemName   string  `json:"item_name"`
	Status     string  `json:"status"`
}

type VolunteerInput struct {
	DonorID    *int   `json:"donor_id,omitempty"`
	DisasterID *int   `json:"disaster_id,omitempty"`
	Skill      string `json:"skill"`
	Location   string `json:"location"`
	Status     string `json:"status"`
}