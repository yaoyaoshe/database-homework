package main

import (
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// 用户与账户
		api.POST("/users", CreateUser)
		api.GET("/users/:id", GetUser)
		api.PUT("/users/:id", UpdateUser)
		api.GET("/users/:id/emails", GetUserEmails)
		api.POST("/users/:id/emails", AddEmail)
		api.DELETE("/users/:id/emails/:emailid", DeleteEmail)
		api.GET("/users/:id/phones", GetUserPhones)
		api.POST("/users/:id/phones", AddPhone)
		api.DELETE("/users/:id/phones/:phoneid", DeletePhone)
		api.GET("/users/:id/providers", GetUserProviders)
		api.POST("/users/:id/providers/:provider_id", LinkProvider)
		api.DELETE("/users/:id/providers/:provider_id", UnlinkProvider)

		// 预约
		api.POST("/appointments", CreateAppointment)
		api.PUT("/appointments/:id/cancel", CancelAppointment)
		api.GET("/users/:id/appointments", GetUserAppointments)

		// 挑战
		api.POST("/challenges", CreateChallenge)
		api.POST("/challenges/:id/invite", InviteToChallenge)
		api.POST("/challenges/:id/join", JoinChallenge)
		api.GET("/users/:id/challenges", GetUserChallenges)
		
		// 数据与报表
		api.POST("/health-data", CreateHealthData)
		api.GET("/health-data/:user_id", GetHealthData)
		api.POST("/reports/generate", GenerateMonthlyReport)
		api.GET("/search/health", SearchHealthData)

		// 统计
		api.GET("/stats/popular_challenges", MostPopularChallenges)
		api.GET("/stats/active_users", MostActiveUsers)
	}
}

// --- 用户 ---
func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.CreatedAt, user.UpdatedAt = time.Now(), time.Now()
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func GetUser(c *gin.Context) {
	var user User
	if err := DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	var req UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Model(&User{}).Where("user_id = ?", c.Param("id")).Updates(map[string]interface{}{
		"name": req.Name, "gender": req.Gender, "date_of_birth": req.DateOfBirth,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}

// --- 账户信息 ---
func GetUserEmails(c *gin.Context) {
	var l []Email
	DB.Where("user_id=?", c.Param("id")).Find(&l)
	c.JSON(http.StatusOK, l)
}

func AddEmail(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))
	var req struct { Email string `json:"email"` }
	c.ShouldBindJSON(&req)
	DB.Create(&Email{UserID: uid, EmailAddress: req.Email})
	c.JSON(http.StatusOK, gin.H{"message": "Added"})
}

func DeleteEmail(c *gin.Context) {
	DB.Delete(&Email{}, c.Param("emailid"))
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func GetUserPhones(c *gin.Context) {
	var l []PhoneNumber
	DB.Where("user_id=?", c.Param("id")).Find(&l)
	c.JSON(http.StatusOK, l)
}

func AddPhone(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))
	var req struct { Phone string `json:"phone"` }
	c.ShouldBindJSON(&req)
	DB.Create(&PhoneNumber{UserID: uid, PhoneNumber: req.Phone})
	c.JSON(http.StatusOK, gin.H{"message": "Added"})
}

func DeletePhone(c *gin.Context) {
	DB.Delete(&PhoneNumber{}, c.Param("phoneid"))
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func GetUserProviders(c *gin.Context) {
	var p []Provider
	DB.Raw("SELECT p.* FROM Provider p JOIN UserProviders up ON p.provider_id=up.provider_id WHERE up.user_id=?", c.Param("id")).Scan(&p)
	c.JSON(http.StatusOK, p)
}

func LinkProvider(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))
	pid, _ := strconv.Atoi(c.Param("provider_id"))
	DB.Create(&UserProvider{UserID: uid, ProviderID: pid, LinkDate: time.Now()})
	c.JSON(http.StatusOK, gin.H{"message": "Linked"})
}

func UnlinkProvider(c *gin.Context) {
	DB.Where("user_id=? AND provider_id=?", c.Param("id"), c.Param("provider_id")).Delete(&UserProvider{})
	c.JSON(http.StatusOK, gin.H{"message": "Unlinked"})
}

// --- 预约 (含 24h 限制) ---
func CreateAppointment(c *gin.Context) {
	var a Appointment
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.Status = "已预约"
	DB.Create(&a)
	c.JSON(http.StatusCreated, a)
}

func CancelAppointment(c *gin.Context) {
	var a Appointment
	if err := DB.First(&a, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "预约不存在"})
		return
	}

	// 24小时限制
	if a.AppointmentDate.Sub(time.Now()) < 24*time.Hour {
		c.JSON(http.StatusForbidden, gin.H{"error": "距离预约时间不足24小时，无法取消"})
		return
	}

	DB.Model(&a).Update("status", "已取消")
	c.JSON(http.StatusOK, gin.H{"message": "已取消"})
}

func GetUserAppointments(c *gin.Context) {
	var list []Appointment
	DB.Where("user_id=?", c.Param("id")).Order("appointment_date desc").Find(&list)
	c.JSON(http.StatusOK, list)
}

// --- 挑战与邀请 ---
func CreateChallenge(c *gin.Context) {
	var ch Challenge
	c.ShouldBindJSON(&ch)
	ch.Status = "进行中"
	DB.Create(&ch)
	c.JSON(http.StatusCreated, ch)
}

func JoinChallenge(c *gin.Context) {
	cid, _ := strconv.Atoi(c.Param("id"))
	var req struct { UserID int `json:"user_id"` }
	c.ShouldBindJSON(&req)
	DB.Create(&Participation{ChallengeID: cid, UserID: req.UserID, JoinedAt: time.Now(), Status: "进行中"})
	c.JSON(http.StatusOK, gin.H{"message": "Joined"})
}

func InviteToChallenge(c *gin.Context) {
	cid, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		SenderID       int    `json:"sender_id"`
		RecipientType  string `json:"recipient_type"`
		RecipientValue string `json:"recipient_value"`
		Message        string `json:"message"`
	}
	c.ShouldBindJSON(&req)

	inv := Invitation{
		ChallengeID:    cid,
		SenderID:       req.SenderID,
		RecipientType:  req.RecipientType,
		RecipientValue: req.RecipientValue,
		Message:        req.Message,
		InvitationDate: time.Now(),
		Status:         "Pending",
	}
	DB.Create(&inv)
	c.JSON(http.StatusCreated, gin.H{"message": "邀请已发送"})
}

func GetUserChallenges(c *gin.Context) {
	var list []struct {
		ChallengeID     int     `json:"challenge_id"`
		ChallengeName   string  `json:"challenge_name"`
		ChallengeType   string  `json:"challenge_type"`
		TargetValue     float64 `json:"target_value"`
		TargetUnit      string  `json:"target_unit"`
		CurrentProgress float64 `json:"current_progress"`
		Status          string  `json:"status"`
	}
	DB.Raw(`
		SELECT c.challenge_id, c.challenge_name, c.challenge_type, c.target_value, c.target_unit, 
		       p.current_progress, p.status 
		FROM Participation p JOIN Challenge c ON p.challenge_id = c.challenge_id 
		WHERE p.user_id = ?`, c.Param("id")).Scan(&list)
	c.JSON(http.StatusOK, list)
}

func GetChallengeParticipants(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// --- 健康数据与高级搜索 ---
func CreateHealthData(c *gin.Context) {
	var d HealthData
	c.ShouldBindJSON(&d)
	DB.Create(&d)
	c.JSON(http.StatusCreated, d)
}

func GetHealthData(c *gin.Context) {
	uid := c.Param("user_id")
	dtype := c.Query("type")
	q := DB.Where("user_id=?", uid).Order("recorded_at desc")
	if dtype != "" {
		q = q.Where("data_type=?", dtype)
	}
	var list []HealthData
	q.Find(&list)
	c.JSON(http.StatusOK, list)
}

func SearchHealthData(c *gin.Context) {
	uid := c.Query("user_id")
	dtype := c.Query("type")
	month := c.Query("month")

	q := DB.Model(&HealthData{}).Where("user_id=? AND data_type=?", uid, dtype)
	if month != "" {
		t, _ := time.Parse("2006-01", month)
		q = q.Where("recorded_at >= ? AND recorded_at < ?", t, t.AddDate(0, 1, 0))
	}

	var list []HealthData
	q.Order("recorded_at desc").Find(&list)
	var total float64
	q.Select("COALESCE(SUM(data_value), 0)").Scan(&total)

	c.JSON(http.StatusOK, gin.H{"total_value": total, "records": list})
}

// --- 报表 (真实计算) ---
func GenerateMonthlyReport(c *gin.Context) {
	var req struct {
		UserID int    `json:"user_id"`
		Month  string `json:"month"`
	}
	c.ShouldBindJSON(&req)

	t, _ := time.Parse("2006-01-02", req.Month)
	start, end := t, t.AddDate(0, 1, 0)

	var wStats struct {
		Avg float64 `json:"avg"`
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	}
	DB.Model(&HealthData{}).Select("AVG(data_value) as avg, MIN(data_value) as min, MAX(data_value) as max").
		Where("user_id=? AND data_type='Weight' AND recorded_at >= ? AND recorded_at < ?", req.UserID, start, end).Scan(&wStats)

	var stepsTotal int
	DB.Model(&HealthData{}).Select("COALESCE(SUM(data_value), 0)").
		Where("user_id=? AND data_type='Steps' AND recorded_at >= ? AND recorded_at < ?", req.UserID, start, end).Scan(&stepsTotal)

	c.JSON(http.StatusOK, gin.H{"month": req.Month, "weight_stats": wStats, "total_steps": stepsTotal, "message": "Report Generated"})
}

// --- 统计 ---
func MostPopularChallenges(c *gin.Context) {
	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = "5"
	}
	limit, _ := strconv.Atoi(limitStr)

	type Row struct {
		ChallengeID      int     `json:"challenge_id"`
		ChallengeName    string  `json:"challenge_name"`
		ParticipantCount int     `json:"participant_count"`
		AvgProgress      float64 `json:"avg_progress"`
	}
	var rows []Row
	// 使用 Raw SQL 避免存储过程调用失败
	err := DB.Raw(`
		SELECT c.challenge_id, c.challenge_name, COUNT(p.user_id) as participant_count, AVG(p.current_progress) as avg_progress 
		FROM Challenge c 
		LEFT JOIN Participation p ON c.challenge_id = p.challenge_id 
		GROUP BY c.challenge_id 
		ORDER BY participant_count DESC 
		LIMIT ?`, limit).Scan(&rows).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func MostActiveUsers(c *gin.Context) {
	type Row struct {
		Name          string  `json:"name"`
		ActivityScore float64 `json:"activity_score"`
	}
	var rows []Row
	// 简单的 SQL 替代存储过程
	DB.Raw(`
		SELECT u.name, COUNT(p.participation_id)*10 + COUNT(a.appointment_id)*5 as activity_score
		FROM User u
		LEFT JOIN Participation p ON u.user_id = p.user_id
		LEFT JOIN Appointment a ON u.user_id = a.user_id
		GROUP BY u.user_id
		ORDER BY activity_score DESC
		LIMIT 10`).Scan(&rows)
	c.JSON(http.StatusOK, rows)
}