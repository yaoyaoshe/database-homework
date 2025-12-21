package main

import (
    "net/http"
    "time"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine) {
    api := r.Group("/api")
    {
        // Users
        api.POST("/users", CreateUser)
        api.GET("/users/:id", GetUser)

        // Email / Phone
        api.POST("/users/:id/emails", AddEmail)
        api.DELETE("/users/:id/emails/:emailid", DeleteEmail)
        api.POST("/users/:id/phones", AddPhone)
        api.DELETE("/users/:id/phones/:phoneid", DeletePhone)

        // Providers
        api.POST("/users/:id/providers/:provider_id", LinkProvider)
        api.DELETE("/users/:id/providers/:provider_id", UnlinkProvider)

        // Appointments
        api.POST("/appointments", CreateAppointment)
        api.PUT("/appointments/:id/cancel", CancelAppointment)
        api.GET("/users/:id/appointments", GetUserAppointments)

        // Challenges
        api.POST("/challenges", CreateChallenge)
        api.POST("/challenges/:id/invite", InviteToChallenge)
        api.POST("/challenges/:id/join", JoinChallenge)
        api.GET("/challenges/:id/participants", GetChallengeParticipants)

        // Reports / Stored Procs
        api.POST("/reports/generate", GenerateMonthlyReport)
        api.GET("/stats/popular_challenges", MostPopularChallenges)
        api.GET("/stats/active_users", MostActiveUsers)
    }
}

// CreateUser 简化示例（注意生产需密码哈希）
func CreateUser(c *gin.Context) {
    var in struct {
        HealthID string `json:"health_id" binding:"required"`
        Name     string `json:"name" binding:"required"`
        DOB      string `json:"date_of_birth"`
        Gender   string `json:"gender"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    var dob *time.Time
    if in.DOB != "" {
        t, err := time.Parse("2006-01-02", in.DOB)
        if err == nil { dob = &t }
    }
    u := User{
        HealthID: in.HealthID,
        Name: in.Name,
        DateOfBirth: dob,
        Gender: in.Gender,
        PasswordHash: in.Password, // TODO: hash it in production
    }
    if err := DB.Create(&u).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, u)
}

func GetUser(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var u User
    if err := DB.First(&u, "user_id = ?", id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error":"user not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, u)
}

func AddEmail(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    var in struct { Email string `json:"email" binding:"required,email"` }
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }
    e := Email{UserID: uid, EmailAddress: in.Email, IsVerified: false}
    if err := DB.Create(&e).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
        return
    }
    c.JSON(http.StatusCreated, e)
}

func DeleteEmail(c *gin.Context) {
    emailID, _ := strconv.Atoi(c.Param("emailid"))
    if err := DB.Delete(&Email{}, emailID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}

func AddPhone(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    var in struct { Phone string `json:"phone" binding:"required"` }
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }
    p := UserPhone{UserID: uid, PhoneNumber: in.Phone, PhoneType: "手机", IsVerified: false}
    if err := DB.Create(&p).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()}); return }
    c.JSON(http.StatusCreated, p)
}

func DeletePhone(c *gin.Context) {
    phoneID, _ := strconv.Atoi(c.Param("phoneid"))
    if err := DB.Delete(&UserPhone{}, phoneID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}

func LinkProvider(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    pid, _ := strconv.Atoi(c.Param("provider_id"))
    up := UserProvider{UserID: uid, ProviderID: pid, LinkDate: time.Now()}
    if err := DB.Create(&up).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, up)
}

func UnlinkProvider(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    pid, _ := strconv.Atoi(c.Param("provider_id"))
    if err := DB.Delete(&UserProvider{}, "user_id = ? AND provider_id = ?", uid, pid).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}

func CreateAppointment(c *gin.Context) {
    var in struct {
        UserID          int    `json:"user_id" binding:"required"`
        ProviderID      int    `json:"provider_id" binding:"required"`
        AppointmentDate string `json:"appointment_date" binding:"required"` // RFC3339
        DurationMinutes int    `json:"duration_minutes"`
        ConsultationType string `json:"consultation_type"`
        Reason          string `json:"reason"`
    }
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }
    dt, err := time.Parse(time.RFC3339, in.AppointmentDate)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"bad datetime"}); return }
    appt := Appointment{
        UserID: in.UserID,
        ProviderID: in.ProviderID,
        AppointmentDate: dt,
        DurationMinutes: in.DurationMinutes,
        ConsultationType: in.ConsultationType,
        Reason: &in.Reason,
        Status: "已预约",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    if err := DB.Create(&appt).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, appt)
}

func CancelAppointment(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var in struct { Reason string `json:"reason"` }
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }

    var appt Appointment
    if err := DB.First(&appt, "appointment_id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error":"appointment not found"}); return
    }
    // 应用层检查：必须在 24 小时前取消
    if time.Until(appt.AppointmentDate) < 24*time.Hour {
        c.JSON(http.StatusForbidden, gin.H{"error":"无法在预约时间24小时内取消预约"})
        return
    }
    now := time.Now()
    appt.Status = "已取消"
    appt.CancellationReason = &in.Reason
    appt.CancelledAt = &now
    appt.UpdatedAt = now
    if err := DB.Save(&appt).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
    }
    c.JSON(http.StatusOK, appt)
}

func GetUserAppointments(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    var appts []Appointment
    DB.Where("user_id = ?", uid).Order("appointment_date desc").Find(&appts)
    c.JSON(http.StatusOK, appts)
}

func CreateChallenge(c *gin.Context) {
    var in struct {
        CreatorID int    `json:"creator_id" binding:"required"`
        Name      string `json:"challenge_name" binding:"required"`
        Description string `json:"description"`
        ChallengeType string `json:"challenge_type"`
        TargetMetric  string `json:"target_metric"`
        TargetValue   float64 `json:"target_value"`
        TargetUnit    string  `json:"target_unit"`
        StartDate string `json:"start_date" binding:"required"` // YYYY-MM-DD
        EndDate   string `json:"end_date" binding:"required"`
        MaxParticipants int `json:"max_participants"`
        IsPublic bool `json:"is_public"`
    }
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }
    sd, err := time.Parse("2006-01-02", in.StartDate)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"bad start date"}); return }
    ed, err := time.Parse("2006-01-02", in.EndDate)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"bad end date"}); return }
    ch := Challenge{
        CreatorID: in.CreatorID,
        ChallengeName: in.Name,
        Description: in.Description,
        ChallengeType: in.ChallengeType,
        TargetMetric: in.TargetMetric,
        TargetValue: in.TargetValue,
        TargetUnit: in.TargetUnit,
        StartDate: &sd,
        EndDate: &ed,
        MaxParticipants: in.MaxParticipants,
        IsPublic: in.IsPublic,
        Status: "筹备中",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    if err := DB.Create(&ch).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
    }
    c.JSON(http.StatusCreated, ch)
}

func InviteToChallenge(c *gin.Context) {
    cid, _ := strconv.Atoi(c.Param("id"))
    var in struct {
        SenderID int `json:"sender_id" binding:"required"`
        RecipientType string `json:"recipient_type" binding:"required"` // 邮箱/手机号/用户ID
        RecipientValue string `json:"recipient_value" binding:"required"`
        Message string `json:"message"`
    }
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }
    inv := Invitation{
        ChallengeID: cid,
        SenderID: in.SenderID,
        RecipientType: in.RecipientType,
        RecipientValue: in.RecipientValue,
        InvitationDate: time.Now(),
        Status: "待处理",
        Message: &in.Message,
        CreatedAt: time.Now(),
    }
    if err := DB.Create(&inv).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
    }
    c.JSON(http.StatusCreated, inv)
}

func JoinChallenge(c *gin.Context) {
    cid, _ := strconv.Atoi(c.Param("id"))
    var in struct { UserID int `json:"user_id" binding:"required"` }
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }
    p := Participation{
        ChallengeID: cid,
        UserID: in.UserID,
        JoinedAt: time.Now(),
        CurrentProgress: 0,
        ProgressUnit: "",
        Status: "参与中",
    }
    if err := DB.Create(&p).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }
    c.JSON(http.StatusCreated, p)
}

func GetChallengeParticipants(c *gin.Context) {
    cid, _ := strconv.Atoi(c.Param("id"))
    var parts []Participation
    DB.Where("challenge_id = ?", cid).Find(&parts)
    c.JSON(http.StatusOK, parts)
}

// GenerateMonthlyReport 调用你在 DB 中的存储过程 GenerateMonthlyReport
func GenerateMonthlyReport(c *gin.Context) {
    var in struct { UserID int `json:"user_id" binding:"required"`; Month string `json:"month" binding:"required"` } // month: YYYY-MM-01
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }
    // 直接调用存储过程
    // CALL GenerateMonthlyReport(in_user_id, in_month);
    res := struct{ Message string }{}
    row := DB.Raw("CALL GenerateMonthlyReport(?, ?)", in.UserID, in.Month).Row()
    // 读取返回消息（存储过程用 SELECT 'msg' AS message）
    if err := row.Scan(&res.Message); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": res.Message})
}

func MostPopularChallenges(c *gin.Context) {
    limitStr := c.Query("limit")
    if limitStr == "" { limitStr = "5" }
    limit, _ := strconv.Atoi(limitStr)
    // 调用存储过程 FindMostPopularChallenges
    type Row struct {
        ChallengeID    int     `json:"challenge_id"`
        ChallengeName  string  `json:"challenge_name"`
        ChallengeType  string  `json:"challenge_type"`
        StartDate      *time.Time `json:"start_date"`
        EndDate        *time.Time `json:"end_date"`
        ParticipantCount int   `json:"participant_count"`
        AvgProgress    float64 `json:"avg_progress"`
        CompletedCount int    `json:"completed_count"`
    }
    var rows []Row
    if err := DB.Raw("CALL FindMostPopularChallenges(?)", limit).Scan(&rows).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
    }
    c.JSON(http.StatusOK, rows)
}

func MostActiveUsers(c *gin.Context) {
    limitStr := c.Query("limit")
    if limitStr == "" { limitStr = "10" }
    limit, _ := strconv.Atoi(limitStr)
    type Row struct {
        UserID int `json:"user_id"`
        HealthID string `json:"health_id"`
        Name string `json:"name"`
        HealthRecordCount int64 `json:"health_record_count"`
        CompletedChallenges int64 `json:"completed_challenges"`
        AppointmentCount int64 `json:"appointment_count"`
        ActivityScore float64 `json:"activity_score"`
    }
    var rows []Row
    if err := DB.Raw("CALL FindMostActiveUsers(?)", limit).Scan(&rows).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
    }
    c.JSON(http.StatusOK, rows)
}
