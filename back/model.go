package main

import "time"

// User 用户表
type User struct {
	UserID            int        `gorm:"column:user_id;primaryKey" json:"user_id"`
	HealthID          string     `gorm:"column:health_id;unique;size:50" json:"health_id"`
	Name              string     `gorm:"column:name;size:100" json:"name"`
	DateOfBirth       *time.Time `gorm:"column:date_of_birth" json:"date_of_birth"`
	Gender            string     `gorm:"column:gender;size:10" json:"gender"`
	PasswordHash      string     `gorm:"column:password_hash" json:"-"`
	PrimaryProviderID *int       `gorm:"column:primary_provider_id" json:"primary_provider_id"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updated_at"`
}
func (User) TableName() string { return "User" }

// Provider 医生表
type Provider struct {
	ProviderID          int     `gorm:"column:provider_id;primaryKey" json:"provider_id"`
	LicenseNumber       string  `gorm:"column:license_number;size:50;unique" json:"license_number"`
	Name                string  `gorm:"column:name;size:100" json:"name"`
	Specialty           string  `gorm:"column:specialty;size:100" json:"specialty"`
	Qualification       string  `gorm:"column:qualification" json:"qualification"`
	IsVerified          bool    `gorm:"column:is_verified" json:"is_verified"`
	ContactInfo         *string `gorm:"column:contact_info" json:"contact_info"`
	HospitalAffiliation string  `gorm:"column:hospital_affiliation" json:"hospital_affiliation"`
}
func (Provider) TableName() string { return "Provider" }

// Email 邮箱表
type Email struct {
	EmailID      int    `gorm:"column:email_id;primaryKey" json:"email_id"`
	UserID       int    `gorm:"column:user_id" json:"user_id"`
	EmailAddress string `gorm:"column:email_address" json:"email_address"`
	IsVerified   bool   `gorm:"column:is_verified" json:"is_verified"`
}
func (Email) TableName() string { return "Email" }

// PhoneNumber 电话表
type PhoneNumber struct {
	PhoneID     int    `gorm:"column:phone_id;primaryKey" json:"phone_id"`
	UserID      int    `gorm:"column:user_id" json:"user_id"`
	PhoneNumber string `gorm:"column:phone_number" json:"phone_number"`
	PhoneType   string `gorm:"column:phone_type" json:"phone_type"`
	IsVerified  bool   `gorm:"column:is_verified" json:"is_verified"`
}
func (PhoneNumber) TableName() string { return "UserPhone" }

// UserProvider 关联表
type UserProvider struct {
	UserID     int       `gorm:"column:user_id;primaryKey"`
	ProviderID int       `gorm:"column:provider_id;primaryKey"`
	LinkDate   time.Time `gorm:"column:link_date"`
}
func (UserProvider) TableName() string { return "UserProvider" }

// Appointment 预约表
type Appointment struct {
	AppointmentID    int       `gorm:"column:appointment_id;primaryKey" json:"appointment_id"`
	UserID           int       `gorm:"column:user_id" json:"user_id"`
	ProviderID       int       `gorm:"column:provider_id" json:"provider_id"`
	AppointmentDate  time.Time `gorm:"column:appointment_date" json:"appointment_date"`
	DurationMinutes  int       `gorm:"column:duration_minutes" json:"duration_minutes"`
	ConsultationType string    `gorm:"column:consultation_type" json:"consultation_type"`
	Reason           string    `gorm:"column:reason" json:"reason"`
	Status           string    `gorm:"column:status" json:"status"`
	Notes            *string   `gorm:"column:notes" json:"notes"`
}
func (Appointment) TableName() string { return "Appointment" }

// Challenge 挑战表
type Challenge struct {
	ChallengeID     int        `gorm:"column:challenge_id;primaryKey" json:"challenge_id"`
	CreatorID       int        `gorm:"column:creator_id" json:"creator_id"`
	ChallengeName   string     `gorm:"column:challenge_name" json:"challenge_name"`
	Description     string     `gorm:"column:description" json:"description"`
	ChallengeType   string     `gorm:"column:challenge_type" json:"challenge_type"`
	TargetMetric    string     `gorm:"column:target_metric" json:"target_metric"`
	TargetValue     float64    `gorm:"column:target_value" json:"target_value"`
	TargetUnit      string     `gorm:"column:target_unit" json:"target_unit"`
	StartDate       time.Time  `gorm:"column:start_date" json:"start_date"`
	EndDate         time.Time  `gorm:"column:end_date" json:"end_date"`
	MaxParticipants int        `gorm:"column:max_participants" json:"max_participants"`
	IsPublic        bool       `gorm:"column:is_public" json:"is_public"`
	Status          string     `gorm:"column:status" json:"status"`
}
func (Challenge) TableName() string { return "Challenge" }

// Participation 参与表
type Participation struct {
	ParticipationID int        `gorm:"column:participation_id;primaryKey" json:"participation_id"`
	ChallengeID     int        `gorm:"column:challenge_id" json:"challenge_id"`
	UserID          int        `gorm:"column:user_id" json:"user_id"`
	JoinedAt        time.Time  `gorm:"column:joined_at" json:"joined_at"`
	CurrentProgress float64    `gorm:"column:current_progress" json:"current_progress"`
	Status          string     `gorm:"column:status" json:"status"`
}
func (Participation) TableName() string { return "Participation" }

// Invitation 邀请表 (Gap Fix)
type Invitation struct {
	InvitationID   int       `gorm:"column:invitation_id;primaryKey" json:"invitation_id"`
	ChallengeID    int       `gorm:"column:challenge_id" json:"challenge_id"`
	SenderID       int       `gorm:"column:sender_id" json:"sender_id"`
	RecipientType  string    `gorm:"column:recipient_type" json:"recipient_type"`
	RecipientValue string    `gorm:"column:recipient_value" json:"recipient_value"`
	InvitationDate time.Time `gorm:"column:invitation_date" json:"invitation_date"`
	Status         string    `gorm:"column:status" json:"status"` // Pending, Accepted, Expired
	Message        string    `gorm:"column:message" json:"message"`
}
func (Invitation) TableName() string { return "Invitation" }

// HealthData 健康数据表 (Gap Fix)
type HealthData struct {
	DataID     int     `gorm:"column:data_id;primaryKey" json:"data_id"`
	UserID     int     `gorm:"column:user_id" json:"user_id"`
	DataType   string  `gorm:"column:data_type" json:"data_type"`
	DataValue  float64 `gorm:"column:data_value" json:"data_value"`
	Unit       string  `gorm:"column:unit" json:"unit"`
	RecordedAt string  `gorm:"column:recorded_at" json:"recorded_at"` // 前端传 ISO 字符串
}
func (HealthData) TableName() string { return "HealthData" }

// API 请求/响应辅助结构体
type UserUpdateReq struct {
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}