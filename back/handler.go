package main

import (
    "net/http"
    "time"
    "strconv"
    // "fmt" // 如果需要打印日志可取消注释
    "crypto/sha256"
    "encoding/hex"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "encoding/json"
)

func RegisterRoutes(r *gin.Engine) {
    api := r.Group("/api")
    {

        // Login
        api.POST("/auth/login",Login)

        // Users
        api.POST("/users", CreateUser)
        api.GET("/users/:id", GetUser)

        // Email / Phone
        api.POST("/users/:id/emails", AddEmail)
        api.GET("/users/:id/emails", GetUserEmails)       // 🆕 新增：获取邮箱列表
        api.DELETE("/users/:id/emails/:emailid", DeleteEmail)
        
        api.POST("/users/:id/phones", AddPhone)
        api.GET("/users/:id/phones", GetUserPhones)       // 🆕 新增：获取电话列表
        api.DELETE("/users/:id/phones/:phoneid", DeletePhone)

        // Providers
        api.GET("/providers", ListProviders)
        api.POST("/users/:id/providers", LinkProvider)
        api.GET("/users/:id/providers", GetUserProviders) // 🆕 新增：获取已关联医生
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
        api.GET("/users/:id/challenges", GetUserChallenges)

        // Reports / Stored Procs
        api.POST("/reports/generate", GenerateMonthlyReport)
        api.GET("/stats/popular_challenges", MostPopularChallenges)
        api.GET("/stats/active_users", MostActiveUsers)

        // 每日挑战记录
        api.POST("/challenges/:id/checkin", CheckinChallenge)
        api.GET("/challenges/:id/progress", GetChallengeDailyProgress)
        api.GET("/challenges/:id/summary", GetChallengeSummary)

        api.POST("/users/:id/family", AddFamilyMember)
        api.GET("/users/:id/family", GetFamilyList)
        api.DELETE("/users/:id/family/:related_id", DeleteFamilyMember)

        api.GET("/users/:id/appointments/search", SearchAppointments)
        api.GET("/providers/search", SearchProviders)
        api.GET("/users/:id/metrics/summary", MetricSummary)

        api.POST("/invitations/accept", AcceptInvitation)
        api.GET("/users/:id/pending-invites", PendingInvites)

        api.PUT("/users/:id/primary-provider", SetPrimaryProvider)
        api.GET("/users/:id/primary-provider", GetPrimaryProvider)

        api.POST("/users/:id/verify", VerifyAccount)

    }
}

//家庭组
type AddFamilyInput struct {
    TargetHealthID string `json:"target_health_id" binding:"required"`
    Relationship   string `json:"relationship" binding:"required"`
}

// 家庭组：修改为通过 Health ID 添加
func AddFamilyMember(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))

    var in AddFamilyInput
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 1. 根据 Health ID 查找被关联用户
    var relatedUser User
    if err := DB.Where("health_id = ?", in.TargetHealthID).First(&relatedUser).Error; err != nil {
        c.JSON(404, gin.H{"error": "未找到该 Health ID 对应的用户"})
        return
    }

    if uid == relatedUser.UserID {
        c.JSON(400, gin.H{"error": "不能关联自己"})
        return
    }

    // 2. 检查是否已经关联
    var cnt int64
    DB.Model(&UserFamily{}).Where("user_id = ? AND related_user_id = ?", uid, relatedUser.UserID).Count(&cnt)
    if cnt > 0 {
         c.JSON(400, gin.H{"error": "该用户已在您的家庭成员列表中"})
         return
    }

    // 3. 创建关联
    fam := UserFamily{
        UserID:        uid,
        RelatedUserID: relatedUser.UserID,
        Relationship:  in.Relationship,
        IsVerified:    false, // 默认未验证
    }

    if err := DB.Create(&fam).Error; err != nil {
        // 捕获可能的枚举错误
        c.JSON(500, gin.H{"error": "添加失败，请检查关系类型是否正确: " + err.Error()})
        return
    }

    c.JSON(201, fam)
}

func GetFamilyList(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))

    var res []struct {
        RelatedUserID int    `json:"user_id"`
        Name          string `json:"name"`
        Relationship  string `json:"relationship"`
        IsVerified    bool   `json:"is_verified"`
    }

    DB.Table("UserFamily uf").
        Select("u.user_id as related_user_id, u.name, uf.relationship, uf.is_verified").
        Joins("JOIN User u ON uf.related_user_id = u.user_id").
        Where("uf.user_id = ?", uid).
        Scan(&res)

    c.JSON(200, res)
}

func DeleteFamilyMember(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    rid, _ := strconv.Atoi(c.Param("related_id"))

    DB.Delete(&UserFamily{}, "user_id = ? AND related_user_id = ?", uid, rid)
    c.Status(204)
}

//查询搜索
func SearchAppointments(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    q := DB.Where("user_id = ?", uid)

    if v := c.Query("status"); v != "" {
        q = q.Where("status = ?", v)
    }
    if s := c.Query("start_date"); s != "" {
        q = q.Where("appointment_date >= ?", s)
    }
    if e := c.Query("end_date"); e != "" {
        q = q.Where("appointment_date <= ?", e)
    }

    var appts []Appointment
    q.Order("appointment_date desc").Find(&appts)
    c.JSON(200, appts)
}

func SearchProviders(c *gin.Context) {
    name := c.Query("name")

    var ps []Provider
    DB.Where("name LIKE ?", "%"+name+"%").Find(&ps)

    c.JSON(200, ps)
}

func MetricSummary(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    metric := c.Query("metric_type")
    month := c.Query("month") // YYYY-MM

    var res struct {
        Total float64 `json:"total"`
        Avg   float64 `json:"avg"`
    }

    DB.Raw(`
        SELECT 
          SUM(metric_value) AS total,
          AVG(metric_value) AS avg
        FROM HealthMetric
        WHERE user_id = ?
          AND metric_type = ?
          AND DATE_FORMAT(measured_at, '%Y-%m') = ?
    `, uid, metric, month).Scan(&res)

    c.JSON(200, res)
}

//接受邀请
func AcceptInvitation(c *gin.Context) {
    var in struct {
        InvitationID int `json:"invitation_id" binding:"required"`
        UserID       int `json:"user_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var inv Invitation
    if err := DB.First(&inv, in.InvitationID).Error; err != nil {
        c.JSON(404, gin.H{"error": "邀请不存在"})
        return
    }

    now := time.Now()
    DB.Model(&inv).Updates(map[string]interface{}{
        "status": "已接受",
        "accepted_at": now,
    })

    p := Participation{
        ChallengeID: inv.ChallengeID,
        UserID: in.UserID,
        JoinedAt: now,
        Status: "参与中",
    }
    DB.Create(&p)

    c.JSON(200, gin.H{"message": "已接受邀请"})
}

func PendingInvites(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))

    var invs []Invitation
    DB.Where("recipient_type = '用户ID' AND recipient_value = ? AND status = '待处理'",
        strconv.Itoa(uid)).
        Find(&invs)

    c.JSON(200, invs)
}

//设置主治
func SetPrimaryProvider(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    var in struct {
        ProviderID int `json:"provider_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var cnt int64
    DB.Model(&UserProvider{}).
        Where("user_id = ? AND provider_id = ?", uid, in.ProviderID).
        Count(&cnt)

    if cnt == 0 {
        c.JSON(400, gin.H{"error": "该医生尚未关联"})
        return
    }

    DB.Model(&User{}).
        Where("user_id = ?", uid).
        Update("primary_provider_id", in.ProviderID)

    c.JSON(200, gin.H{"message": "主治医生设置成功"})
}

func GetPrimaryProvider(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))

    var res Provider
    DB.Table("User u").
        Joins("JOIN Provider p ON u.primary_provider_id = p.provider_id").
        Where("u.user_id = ?", uid).
        Scan(&res)

    c.JSON(200, res)
}

func VerifyAccount(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    var in struct {
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var u User
    if err := DB.First(&u, uid).Error; err != nil {
        c.JSON(404, gin.H{"error": "用户不存在"})
        return
    }

    if in.Password != u.PasswordHash {
        c.JSON(401, gin.H{"error": "密码错误"})
        return
    }

    c.JSON(200, gin.H{"message": "账户验证通过"})
}


func Login(c *gin.Context) {
    var in struct {
        Identifier string `json:"identifier" binding:"required"`
        Password   string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var user User
    var err error

    // 1️⃣ 先按 health_id 查
    err = DB.Where("health_id = ?", in.Identifier).First(&user).Error

    // 2️⃣ 如果没找到，用已验证邮箱查
    if err != nil {
        err = DB.Joins("JOIN Email e ON e.user_id = User.user_id").
            Where("e.email_address = ?", in.Identifier).
            First(&user).Error
    }

    // 3️⃣ 如果还没找到，用已验证手机号查
    if err != nil {
        err = DB.Joins("JOIN UserPhone p ON p.user_id = User.user_id").
            Where("p.phone_number = ?", in.Identifier).
            First(&user).Error
    }

    if err != nil {
        c.JSON(401, gin.H{"error": "用户不存在或联系方式未验证"})
        return
    }
    
    // 密码校验 (兼容前端传来的 SHA256 哈希)
    if user.PasswordHash != in.Password {
         hash := sha256.Sum256([]byte(in.Password))
         hashedInput := hex.EncodeToString(hash[:])
         if user.PasswordHash != hashedInput {
             c.JSON(401, gin.H{"error": "密码错误"})
             return
         }
    }

    c.JSON(200, gin.H{
        "user_id": user.UserID,
        "health_id": user.HealthID,
        "name": user.Name,
    })
}

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
    
    hash := sha256.Sum256([]byte(in.Password))
    hashedPassword := hex.EncodeToString(hash[:])

    u := User{
        HealthID: in.HealthID,
        Name: in.Name,
        DateOfBirth: dob,
        Gender: in.Gender,
        PasswordHash: hashedPassword, 
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

// 🆕 获取用户的邮箱列表
func GetUserEmails(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    var emails []Email
    if err := DB.Where("user_id = ?", uid).Find(&emails).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, emails)
}

func AddEmail(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    var in struct { Email string `json:"email" binding:"required,email"` }
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }
    e := Email{UserID: uid, EmailAddress: in.Email, IsVerified: true}
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

// 🆕 获取用户的电话列表
func GetUserPhones(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    var phones []UserPhone
    if err := DB.Where("user_id = ?", uid).Find(&phones).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, phones)
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

func ListProviders(c *gin.Context) {
    var providers []Provider
    q := DB.Model(&Provider{})
    if s := c.Query("specialty"); s != "" {
        q = q.Where("specialty = ?", s)
    }
    if c.Query("verified") == "true" {
        q = q.Where("is_verified = true")
    }
    q.Find(&providers)
    c.JSON(200, providers)
}

// 🆕 获取用户关联的医生列表
func GetUserProviders(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    
    // 定义返回结构
    type LinkedProvider struct {
        ProviderID       int       `json:"provider_id"`
        Name             string    `json:"name"`
        Specialty        string    `json:"specialty"`
        ContactInfo      string    `json:"contact_info"`
        RelationshipType string    `json:"relationship_type"`
        LinkDate         time.Time `json:"link_date"`
    }

    var result []LinkedProvider

    // 联表查询
    err := DB.Table("Provider").
        Select("Provider.provider_id, Provider.name, Provider.specialty, Provider.contact_info, UserProvider.relationship_type, UserProvider.link_date").
        Joins("JOIN UserProvider ON UserProvider.provider_id = Provider.provider_id").
        Where("UserProvider.user_id = ?", uid).
        Scan(&result).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}

func LinkProvider(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))

    var in struct {
        ProviderID int    `json:"provider_id" binding:"required"`
        Relation   string `json:"relationship_type" binding:"required"`
    }
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    up := UserProvider{
        UserID: uid,
        ProviderID: in.ProviderID,
        RelationshipType: in.Relation,
        LinkDate: time.Now(),
    }

    if err := DB.Create(&up).Error; err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(201, up)
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
        AppointmentDate string `json:"appointment_date" binding:"required"`
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

func GetUserChallenges(c *gin.Context) {
    uid, _ := strconv.Atoi(c.Param("id"))
    type UserChallengeResponse struct {
        ChallengeID     int     `json:"challenge_id"`
        ChallengeName   string  `json:"challenge_name"`
        Description     string  `json:"description"`
        Status          string  `json:"status"`
        CurrentProgress float64 `json:"current_progress"`
        TargetValue     float64 `json:"target_value"`
        TargetUnit      string  `json:"target_unit"`
    }
    var results []UserChallengeResponse
    err := DB.Table("Participation").
        Select("Participation.challenge_id, Challenge.challenge_name, Challenge.description, Participation.status, Participation.current_progress, Challenge.target_value, Challenge.target_unit").
        Joins("JOIN Challenge ON Challenge.challenge_id = Participation.challenge_id").
        Where("Participation.user_id = ?", uid).
        Scan(&results).Error
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, results)
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
        StartDate string `json:"start_date" binding:"required"`
        EndDate   string `json:"end_date" binding:"required"`
        MaxParticipants int `json:"max_participants"`
        IsPublic bool `json:"is_public"`
    }
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }
    sd, err := time.Parse("2006-01-02", in.StartDate)
    if err != nil { 
         t, e2 := time.Parse(time.RFC3339, in.StartDate)
         if e2 != nil {
             c.JSON(http.StatusBadRequest, gin.H{"error":"bad start date"}); return 
         }
         sd = t
    }
    ed, err := time.Parse("2006-01-02", in.EndDate)
    if err != nil { 
         t, e2 := time.Parse(time.RFC3339, in.EndDate)
         if e2 != nil {
             c.JSON(http.StatusBadRequest, gin.H{"error":"bad end date"}); return 
         }
         ed = t
    }

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
        RecipientType string `json:"recipient_type" binding:"required"`
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

type MonthlyReportResponse struct {
    UserID                int             `json:"user_id"`
    ReportMonth           string          `json:"report_month"`
    TotalAppointments     int             `json:"total_appointments"`
    CompletedAppointments int             `json:"completed_appointments"`
    CancelledAppointments int             `json:"cancelled_appointments"`
    TotalChallenges       int             `json:"total_challenges"`
    CompletedChallenges   int             `json:"completed_challenges"`
    HealthSummary         json.RawMessage `json:"health_summary"`
    Recommendations       string          `json:"recommendations"`
    GeneratedAt           time.Time       `json:"generated_at"`
}

func GenerateMonthlyReport(c *gin.Context) {
    var in struct {
        UserID int    `json:"user_id" binding:"required"`
        Month  string `json:"month" binding:"required"`
    }

    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := DB.Exec("CALL GenerateMonthlyReport(?, ?)", in.UserID, in.Month).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "生成月度报表失败: " + err.Error()})
        return
    }

    var report MonthlyReportResponse
    err := DB.Raw(`
        SELECT 
            user_id,
            DATE_FORMAT(report_month, '%Y-%m-%d') AS report_month,
            total_appointments,
            completed_appointments,
            cancelled_appointments,
            total_challenges,
            completed_challenges,
            health_summary,
            recommendations,
            generated_at
        FROM MonthlyReport
        WHERE user_id = ? AND report_month = ?
    `, in.UserID, in.Month).Scan(&report).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "读取月度报表失败: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, report)
}

func MostPopularChallenges(c *gin.Context) {
    limitStr := c.Query("limit")
    if limitStr == "" { limitStr = "5" }
    limit, _ := strconv.Atoi(limitStr)
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

type ChallengeDailyProgress struct {
    ProgressID    int       `gorm:"primaryKey;column:progress_id"`
    ChallengeID   int       `gorm:"column:challenge_id"`
    UserID        int       `gorm:"column:user_id"`
    ProgressDate  time.Time `gorm:"column:progress_date"`
    ProgressValue float64   `gorm:"column:progress_value"`
    ProgressUnit  string    `gorm:"column:progress_unit"`
    IsCompleted   bool      `gorm:"column:is_completed"`
    Notes         *string   `gorm:"column:notes"`
    RecordedAt    time.Time `gorm:"column:recorded_at"`
}

func (ChallengeDailyProgress) TableName() string { return "ChallengeDailyProgress" }

func CheckinChallenge(c *gin.Context) {
    challengeID, _ := strconv.Atoi(c.Param("id"))

    var in struct {
        UserID        int     `json:"user_id" binding:"required"`
        Date          string  `json:"date" binding:"required"`
        ProgressValue float64 `json:"progress_value" binding:"required"`
        ProgressUnit  string  `json:"progress_unit"`
        IsCompleted   bool    `json:"is_completed"`
        Notes         string  `json:"notes"`
    }

    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    progressDate, err := time.Parse("2006-01-02", in.Date)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式应为 YYYY-MM-DD"})
        return
    }

    record := ChallengeDailyProgress{
        ChallengeID:   challengeID,
        UserID:        in.UserID,
        ProgressDate:  progressDate,
        ProgressValue: in.ProgressValue,
        ProgressUnit:  in.ProgressUnit,
        IsCompleted:   in.IsCompleted,
        RecordedAt:    time.Now(),
    }
    if in.Notes != "" { record.Notes = &in.Notes }

    if err := DB.Create(&record).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "今日已打卡或数据错误：" + err.Error()})
        return
    }

    // 同步更新 Participation
    if err := DB.Exec(`
        UPDATE Participation
        SET 
            current_progress = current_progress + ?,
            last_update = NOW(),
            status = CASE 
                WHEN current_progress + ? >= (
                    SELECT target_value FROM Challenge WHERE challenge_id = ?
                ) THEN '已完成'
                ELSE status
            END
        WHERE challenge_id = ? AND user_id = ?
    `, in.ProgressValue, in.ProgressValue, challengeID, challengeID, in.UserID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "更新累计进度失败：" + err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "打卡成功", "data": record})
}

func GetChallengeDailyProgress(c *gin.Context) {
    challengeID, _ := strconv.Atoi(c.Param("id"))
    userID, _ := strconv.Atoi(c.Query("user_id"))
    var records []ChallengeDailyProgress
    DB.Where("challenge_id = ? AND user_id = ?", challengeID, userID).Order("progress_date asc").Find(&records)
    c.JSON(http.StatusOK, records)
}

func GetChallengeSummary(c *gin.Context) {
    challengeID, _ := strconv.Atoi(c.Param("id"))
    userID, _ := strconv.Atoi(c.Query("user_id"))
    var p Participation
    if err := DB.First(&p, "challenge_id = ? AND user_id = ?", challengeID, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "未找到参与记录"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "challenge_id":    challengeID,
        "user_id":         userID,
        "current_progress": p.CurrentProgress,
        "status":          p.Status,
        "last_update":     p.LastUpdate,
    })
}