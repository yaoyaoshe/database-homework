package main

import "time"

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

type Provider struct {
    ProviderID    int       `gorm:"column:provider_id;primaryKey" json:"provider_id"`
    LicenseNumber string    `gorm:"column:license_number;size:50;unique" json:"license_number"`
    Name          string    `gorm:"column:name;size:100" json:"name"`
    Specialty     string    `gorm:"column:specialty;size:100" json:"specialty"`
    Qualification string    `gorm:"column:qualification" json:"qualification"`
    IsVerified    bool      `gorm:"column:is_verified" json:"is_verified"`
    ContactInfo   string    `gorm:"column:contact_info" json:"contact_info"`
    CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
    UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}
func (Provider) TableName() string { return "Provider" }

type Email struct {
    EmailID      int        `gorm:"column:email_id;primaryKey" json:"email_id"`
    UserID       int        `gorm:"column:user_id;index" json:"user_id"`
    EmailAddress string     `gorm:"column:email_address;size:100;unique" json:"email_address"`
    IsVerified   bool       `gorm:"column:is_verified" json:"is_verified"`
    IsPrimary    bool       `gorm:"column:is_primary" json:"is_primary"`
    VerificationCode *string `gorm:"column:verification_code" json:"-"`
    VerifiedAt   *time.Time `gorm:"column:verified_at" json:"verified_at"`
    CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
}
func (Email) TableName() string { return "Email" }

type UserPhone struct {
    PhoneID      int        `gorm:"column:phone_id;primaryKey" json:"phone_id"`
    UserID       int        `gorm:"column:user_id;index" json:"user_id"`
    PhoneNumber  string     `gorm:"column:phone_number;size:20" json:"phone_number"`
    PhoneType    string     `gorm:"column:phone_type;size:20" json:"phone_type"`
    IsVerified   bool       `gorm:"column:is_verified" json:"is_verified"`
    IsPrimary    bool       `gorm:"column:is_primary" json:"is_primary"`
    VerificationCode *string`gorm:"column:verification_code" json:"-"`
    VerifiedAt   *time.Time `gorm:"column:verified_at" json:"verified_at"`
    CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
}
func (UserPhone) TableName() string { return "UserPhone" }

type UserProvider struct {
    UserID          int       `gorm:"column:user_id;primaryKey" json:"user_id"`
    ProviderID      int       `gorm:"column:provider_id;primaryKey" json:"provider_id"`
    RelationshipType string   `gorm:"column:relationship_type" json:"relationship_type"`
    LinkDate        time.Time `gorm:"column:link_date" json:"link_date"`
    Notes           *string   `gorm:"column:notes" json:"notes"`
    CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
}
func (UserProvider) TableName() string { return "UserProvider" }

type UserFamily struct {
    UserID        int        `gorm:"column:user_id;primaryKey" json:"user_id"`
    RelatedUserID int        `gorm:"column:related_user_id;primaryKey" json:"related_user_id"`
    Relationship  string     `gorm:"column:relationship;size:30" json:"relationship"`
    IsVerified    bool       `gorm:"column:is_verified" json:"is_verified"`
    CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
    VerifiedAt    *time.Time `gorm:"column:verified_at" json:"verified_at"`
}
func (UserFamily) TableName() string { return "UserFamily" }

type Appointment struct {
    AppointmentID     int        `gorm:"column:appointment_id;primaryKey" json:"appointment_id"`
    UserID            int        `gorm:"column:user_id;index" json:"user_id"`
    ProviderID        int        `gorm:"column:provider_id;index" json:"provider_id"`
    AppointmentDate   time.Time  `gorm:"column:appointment_date;index" json:"appointment_date"`
    DurationMinutes   int        `gorm:"column:duration_minutes" json:"duration_minutes"`
    ConsultationType  string     `gorm:"column:consultation_type" json:"consultation_type"`
    Reason            *string    `gorm:"column:reason" json:"reason"`
    Memo              *string    `gorm:"column:memo" json:"memo"`
    Status            string     `gorm:"column:status;size:20" json:"status"`
    CancellationReason *string   `gorm:"column:cancellation_reason" json:"cancellation_reason"`
    CancelledAt       *time.Time `gorm:"column:cancelled_at" json:"cancelled_at"`
    CreatedAt         time.Time  `gorm:"column:created_at" json:"created_at"`
    UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updated_at"`
}
func (Appointment) TableName() string { return "Appointment" }

// 新增 HealthMetric 模型
type HealthMetric struct {
    MetricID          int        `gorm:"column:metric_id;primaryKey" json:"metric_id"`
    UserID            int        `gorm:"column:user_id" json:"user_id"`
    MetricType        string     `gorm:"column:metric_type" json:"metric_type"`
    MetricValue       float64    `gorm:"column:metric_value" json:"metric_value"`
    Unit              string     `gorm:"column:unit" json:"unit"`
    MeasuredAt        time.Time  `gorm:"column:measured_at" json:"measured_at"`
    MeasurementDevice *string    `gorm:"column:measurement_device" json:"measurement_device"`
    Notes             *string    `gorm:"column:notes" json:"notes"`
    Source            string     `gorm:"column:source" json:"source"`
    CreatedAt         time.Time  `gorm:"column:created_at" json:"created_at"`
}
func (HealthMetric) TableName() string { return "HealthMetric" }

type Challenge struct {
    ChallengeID   int        `gorm:"column:challenge_id;primaryKey" json:"challenge_id"`
    CreatorID     int        `gorm:"column:creator_id" json:"creator_id"`
    ChallengeName string     `gorm:"column:challenge_name" json:"challenge_name"`
    Description   string     `gorm:"column:description" json:"description"`
    ChallengeType string     `gorm:"column:challenge_type" json:"challenge_type"`
    TargetMetric  string     `gorm:"column:target_metric" json:"target_metric"`
    TargetValue   float64    `gorm:"column:target_value" json:"target_value"`
    TargetUnit    string     `gorm:"column:target_unit" json:"target_unit"`
    StartDate     *time.Time `gorm:"column:start_date" json:"start_date"`
    EndDate       *time.Time `gorm:"column:end_date" json:"end_date"`
    MaxParticipants int      `gorm:"column:max_participants" json:"max_participants"`
    IsPublic      bool       `gorm:"column:is_public" json:"is_public"`
    Status        string     `gorm:"column:status" json:"status"`
    CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
    UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
}
func (Challenge) TableName() string { return "Challenge" }

type Participation struct {
    ParticipationID int        `gorm:"column:participation_id;primaryKey" json:"participation_id"`
    ChallengeID     int        `gorm:"column:challenge_id;index" json:"challenge_id"`
    UserID          int        `gorm:"column:user_id;index" json:"user_id"`
    JoinedAt        time.Time  `gorm:"column:joined_at" json:"joined_at"`
    CurrentProgress float64    `gorm:"column:current_progress" json:"current_progress"`
    ProgressUnit    string     `gorm:"column:progress_unit" json:"progress_unit"`
    LastUpdate      *time.Time `gorm:"column:last_update" json:"last_update"`
    Status          string     `gorm:"column:status" json:"status"`
    CompletedAt     *time.Time `gorm:"column:completed_at" json:"completed_at"`
    Notes           *string    `gorm:"column:notes" json:"notes"`
}
func (Participation) TableName() string { return "Participation" }

type Invitation struct {
    InvitationID   int        `gorm:"column:invitation_id;primaryKey" json:"invitation_id"`
    ChallengeID    int        `gorm:"column:challenge_id" json:"challenge_id"`
    SenderID       int        `gorm:"column:sender_id" json:"sender_id"`
    RecipientType  string     `gorm:"column:recipient_type" json:"recipient_type"`
    RecipientValue string     `gorm:"column:recipient_value" json:"recipient_value"`
    InvitationDate time.Time  `gorm:"column:invitation_date" json:"invitation_date"`
    AcceptedAt     *time.Time `gorm:"column:accepted_at" json:"accepted_at"`
    Status         string     `gorm:"column:status" json:"status"`
    Message        *string    `gorm:"column:message" json:"message"`
    CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
}
func (Invitation) TableName() string { return "Invitation" }

type ChallengeDailyProgress struct {
    ProgressID    int       `gorm:"primaryKey;column:progress_id" json:"progress_id"`
    ChallengeID   int       `gorm:"column:challenge_id" json:"challenge_id"`
    UserID        int       `gorm:"column:user_id" json:"user_id"`
    ProgressDate  time.Time `gorm:"column:progress_date" json:"progress_date"`
    ProgressValue float64   `gorm:"column:progress_value" json:"progress_value"`
    ProgressUnit  string    `gorm:"column:progress_unit" json:"progress_unit"`
    IsCompleted   bool      `gorm:"column:is_completed" json:"is_completed"`
    Notes         *string   `gorm:"column:notes" json:"notes"`
    RecordedAt    time.Time `gorm:"column:recorded_at" json:"recorded_at"`
}

func (ChallengeDailyProgress) TableName() string { return "ChallengeDailyProgress" }