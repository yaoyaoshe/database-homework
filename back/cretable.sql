-- ============================================================
-- HealthTrack Personal Wellness Platform 数据库建表脚本
-- ============================================================

-- 创建数据库
DROP DATABASE IF EXISTS HealthTrackDB;
CREATE DATABASE HealthTrackDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE HealthTrackDB;

-- ==================== 基础表 ====================

-- 1. 用户表（User）
CREATE TABLE User (
    user_id INT PRIMARY KEY AUTO_INCREMENT,
    health_id VARCHAR(50) UNIQUE NOT NULL COMMENT '用户的唯一健康标识号',
    name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL COMMENT '出生日期',
    gender ENUM('男', '女', '其他') DEFAULT '其他',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希值',
    primary_provider_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_health_id (health_id),
    INDEX idx_name (name)
) COMMENT '用户基本信息表';

-- 2. 医疗提供者表（Provider）
CREATE TABLE Provider (
    provider_id INT PRIMARY KEY AUTO_INCREMENT,
    license_number VARCHAR(50) UNIQUE NOT NULL COMMENT '医疗执照号，唯一识别',
    name VARCHAR(100) NOT NULL,
    specialty VARCHAR(100) NOT NULL COMMENT '专业领域',
    qualification VARCHAR(200) COMMENT '资格认证',
    is_verified BOOLEAN DEFAULT FALSE COMMENT '是否已验证',
    contact_info TEXT COMMENT '联系方式信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_license_number (license_number),
    INDEX idx_specialty (specialty)
) COMMENT '医疗服务提供者表';

-- ==================== 联系方式表 ====================

-- 3. 邮箱表（Email）
CREATE TABLE Email (
    email_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    email_address VARCHAR(100) NOT NULL UNIQUE,
    is_verified BOOLEAN DEFAULT FALSE,
    is_primary BOOLEAN DEFAULT FALSE COMMENT '是否为主要邮箱',
    verification_code VARCHAR(10) COMMENT '验证码',
    verified_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_email_address (email_address),
    INDEX idx_is_primary (is_primary)
) COMMENT '用户邮箱表';

-- 4. 用户电话表（UserPhone）
CREATE TABLE UserPhone (
    phone_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    phone_type ENUM('手机', '家庭', '工作', '紧急联系人') DEFAULT '手机',
    is_verified BOOLEAN DEFAULT FALSE,
    is_primary BOOLEAN DEFAULT FALSE COMMENT '是否为主要电话',
    verification_code VARCHAR(10) COMMENT '验证码',
    verified_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    INDEX idx_user_phone (user_id, is_primary),
    INDEX idx_phone_verified (phone_number, is_verified)
) COMMENT '用户电话表';

-- ==================== 关系表 ====================

-- 5. 用户-提供者关联表（UserProvider）
CREATE TABLE UserProvider (
    user_id INT NOT NULL,
    provider_id INT NOT NULL,
    relationship_type ENUM('主治医生', '专科医生', '家庭医生', '咨询医生') DEFAULT '主治医生',
    link_date DATE NOT NULL,
    notes TEXT COMMENT '备注信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, provider_id),
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES Provider(provider_id) ON DELETE CASCADE,
    INDEX idx_link_date (link_date)
) COMMENT '用户与医疗服务提供者关联表';

-- 6. 家庭组关系表（UserFamily）
CREATE TABLE UserFamily (
    user_id INT NOT NULL,
    related_user_id INT NOT NULL,
    relationship ENUM('父母', '子女', '配偶','朋友', '兄弟姐妹', '监护人', '亲戚') NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE COMMENT '关系是否已验证',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified_at TIMESTAMP NULL,
    PRIMARY KEY (user_id, related_user_id),
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    FOREIGN KEY (related_user_id) REFERENCES User(user_id) ON DELETE CASCADE
) COMMENT '用户家庭关系表';

-- ==================== 核心业务表 ====================

-- 7. 预约表（Appointment）
CREATE TABLE Appointment (
    appointment_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    provider_id INT NOT NULL,
    appointment_date TIMESTAMP NOT NULL COMMENT '预约时间',
    duration_minutes INT DEFAULT 30 COMMENT '预约时长（分钟）',
    consultation_type ENUM('线下就诊', '线上咨询', '电话咨询') DEFAULT '线下就诊',
    reason TEXT COMMENT '预约原因/症状描述',
    memo TEXT COMMENT '备注',
    status ENUM('已预约', '已取消', '已完成', '未到诊') DEFAULT '已预约',
    cancellation_reason VARCHAR(255) COMMENT '取消原因',
    cancelled_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES Provider(provider_id) ON DELETE CASCADE,
    INDEX idx_appointment_date (appointment_date),
    INDEX idx_user_date (user_id, appointment_date),
    INDEX idx_status_date (status, appointment_date),
    INDEX idx_provider_date (provider_id, appointment_date)
) COMMENT '医疗预约表';

-- 8. 健康指标表（HealthMetric）
CREATE TABLE HealthMetric (
    metric_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    metric_type ENUM(
        '体重', '身高', 'BMI', '体脂率',
        '血压收缩压', '血压舒张压', '心率', 
        '血糖空腹', '血糖餐后', '糖化血红蛋白',
        '总胆固醇', '低密度脂蛋白', '高密度脂蛋白', '甘油三酯',
        '步数', '睡眠时长', '运动时长', '卡路里消耗'
    ) NOT NULL,
    metric_value DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20) NOT NULL COMMENT '单位',
    measured_at TIMESTAMP NOT NULL COMMENT '测量时间',
    measurement_device VARCHAR(100) COMMENT '测量设备',
    notes TEXT COMMENT '备注',
    source ENUM('手动录入', '健康设备', '手机APP', '医院记录') DEFAULT '手动录入',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    INDEX idx_user_metric_date (user_id, metric_type, measured_at),
    INDEX idx_metric_type (metric_type)
) COMMENT '健康指标记录表';

-- ==================== 健康挑战表 ====================

-- 9. 健康挑战表（Challenge）
CREATE TABLE Challenge (
    challenge_id INT PRIMARY KEY AUTO_INCREMENT,
    creator_id INT NOT NULL,
    challenge_name VARCHAR(100) NOT NULL COMMENT '挑战名称',
    description TEXT NOT NULL COMMENT '挑战描述',
    challenge_type ENUM('减重', '运动', '饮食', '睡眠', '健康习惯', '综合') NOT NULL,
    target_metric ENUM('体重', '步数', '睡眠时长', '运动时长', '卡路里消耗') NOT NULL,
    target_value DECIMAL(10,2) NOT NULL COMMENT '目标值',
    target_unit VARCHAR(20) NOT NULL COMMENT '目标单位',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    max_participants INT DEFAULT 100 COMMENT '最大参与人数',
    is_public BOOLEAN DEFAULT TRUE COMMENT '是否公开挑战',
    status ENUM('筹备中', '进行中', '已结束', '已取消') DEFAULT '筹备中',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES User(user_id) ON DELETE CASCADE,
    INDEX idx_dates (start_date, end_date),
    INDEX idx_status (status),
    INDEX idx_type (challenge_type)
) COMMENT '健康挑战表';

-- 10. 挑战参与表（Participation）
CREATE TABLE Participation (
    participation_id INT PRIMARY KEY AUTO_INCREMENT,
    challenge_id INT NOT NULL,
    user_id INT NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    current_progress DECIMAL(10,2) DEFAULT 0.0 COMMENT '当前进度',
    progress_unit VARCHAR(20) COMMENT '进度单位',
    last_update TIMESTAMP NULL COMMENT '最后更新进度时间',
    status ENUM('参与中', '已完成', '已放弃') DEFAULT '参与中',
    completed_at TIMESTAMP NULL,
    notes TEXT COMMENT '参与备注',
    FOREIGN KEY (challenge_id) REFERENCES Challenge(challenge_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    UNIQUE KEY unique_challenge_user (challenge_id, user_id),
    INDEX idx_challenge_status (challenge_id, status),
    INDEX idx_user_challenges (user_id, status),
    INDEX idx_progress (challenge_id, current_progress DESC)
) COMMENT '挑战参与记录表';

-- 11. 挑战邀请表（Invitation）
CREATE TABLE Invitation (
    invitation_id INT PRIMARY KEY AUTO_INCREMENT,
    challenge_id INT NOT NULL,
    sender_id INT NOT NULL COMMENT '发送者用户ID',
    recipient_type ENUM('邮箱', '手机号', '用户ID') NOT NULL,
    recipient_value VARCHAR(100) NOT NULL COMMENT '收件人标识（邮箱/手机号/用户ID）',
    invitation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP NULL,
    status ENUM('待处理', '已接受', '已拒绝', '已过期', '已取消') DEFAULT '待处理',
    message TEXT COMMENT '邀请消息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (challenge_id) REFERENCES Challenge(challenge_id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES User(user_id) ON DELETE CASCADE,
    INDEX idx_status_date (status, invitation_date),
    INDEX idx_recipient (recipient_type, recipient_value)
) COMMENT '挑战邀请表';

-- 12. 健康月报表（MonthlyReport）
CREATE TABLE MonthlyReport (
    report_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    report_month DATE NOT NULL COMMENT '报告月份，格式为YYYY-MM-01',
    total_appointments INT DEFAULT 0 COMMENT '总预约次数',
    completed_appointments INT DEFAULT 0 COMMENT '已完成预约',
    cancelled_appointments INT DEFAULT 0 COMMENT '已取消预约',
    total_challenges INT DEFAULT 0 COMMENT '参与挑战总数',
    completed_challenges INT DEFAULT 0 COMMENT '完成挑战数',
    health_summary JSON COMMENT '健康数据摘要（JSON格式）',
    recommendations TEXT COMMENT '健康建议',
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_viewed_at TIMESTAMP NULL COMMENT '最后查看时间',
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_month (user_id, report_month),
    INDEX idx_report_month (report_month),
    INDEX idx_user_reports (user_id, report_month DESC)
) COMMENT '健康月度报告表';

CREATE TABLE ChallengeDailyProgress (
    progress_id INT PRIMARY KEY AUTO_INCREMENT,
    challenge_id INT NOT NULL,
    user_id INT NOT NULL,
    progress_date DATE NOT NULL COMMENT '记录日期',
    progress_value DECIMAL(10,2) NOT NULL COMMENT '当日进度值',
    progress_unit VARCHAR(20) COMMENT '单位，如步、千卡',
    is_completed BOOLEAN DEFAULT FALSE COMMENT '当日是否完成',
    notes TEXT COMMENT '备注说明',
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (challenge_id) REFERENCES Challenge(challenge_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,

    UNIQUE KEY uniq_user_challenge_date (challenge_id, user_id, progress_date),
    INDEX idx_challenge_date (challenge_id, progress_date),
    INDEX idx_user_date (user_id, progress_date)
) COMMENT='挑战每日进度记录表';


-- ==================== 系统功能表 ====================

-- 13. 用户活动记录表（用于统计活跃用户）
CREATE TABLE UserActivity (
    activity_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    activity_type ENUM('登录', '记录健康数据', '创建预约', '完成挑战', '创建挑战', '查看报告') NOT NULL,
    activity_details JSON COMMENT '活动详情',
    performed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    INDEX idx_user_activity (user_id, performed_at),
    INDEX idx_activity_type_date (activity_type, performed_at)
) COMMENT '用户活动记录表';

-- ==================== 添加外键约束 ====================

-- 为User表添加主治医生外键
ALTER TABLE User 
ADD CONSTRAINT fk_user_primary_provider 
FOREIGN KEY (primary_provider_id) 
REFERENCES Provider(provider_id) 
ON DELETE SET NULL;

-- ==================== 创建触发器 ====================

-- 触发器1：确保只能在预约24小时前取消
DELIMITER //
CREATE TRIGGER before_appointment_cancellation
BEFORE UPDATE ON Appointment
FOR EACH ROW
BEGIN
    IF NEW.status = '已取消' AND OLD.status != '已取消' THEN
        IF TIMESTAMPDIFF(HOUR, NOW(), OLD.appointment_date) < 24 THEN
            SIGNAL SQLSTATE '45000' 
            SET MESSAGE_TEXT = '无法在预约时间24小时内取消预约。';
        END IF;
        SET NEW.cancelled_at = NOW();
    END IF;
END;
//
DELIMITER ;


-- 触发器4：检查不能自己关联自己（替代 CHECK 约束）
DELIMITER //
CREATE TRIGGER before_userfamily_insert
BEFORE INSERT ON UserFamily
FOR EACH ROW
BEGIN
    IF NEW.user_id = NEW.related_user_id THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = '不能自己关联自己。';
    END IF;
END;
//
DELIMITER ;

DELIMITER //
CREATE TRIGGER before_userfamily_update
BEFORE UPDATE ON UserFamily
FOR EACH ROW
BEGIN
    IF NEW.user_id = NEW.related_user_id THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = '不能自己关联自己。';
    END IF;
END;
//
DELIMITER ;

-- ==================== 创建存储过程 ====================

-- 存储过程1：生成月度健康报告
-- 修改：即使存在也重新生成（确保数据最新）
DELIMITER //
CREATE PROCEDURE GenerateMonthlyReport(
    IN p_user_id INT,
    IN p_month DATE
)
BEGIN
    DECLARE v_total_appointments INT;
    DECLARE v_completed_appointments INT;
    DECLARE v_cancelled_appointments INT;
    DECLARE v_total_challenges INT;
    DECLARE v_completed_challenges INT;
    DECLARE v_health_summary JSON;
    
    -- 1. 删除已存在的旧报告（确保重新计算）
    DELETE FROM MonthlyReport 
    WHERE user_id = p_user_id AND report_month = p_month;
    
    -- 2. 计算当月预约统计
    SELECT 
        COUNT(*) AS total,
        SUM(CASE WHEN status = '已完成' THEN 1 ELSE 0 END) AS completed,
        SUM(CASE WHEN status = '已取消' THEN 1 ELSE 0 END) AS cancelled
    INTO v_total_appointments, v_completed_appointments, v_cancelled_appointments
    FROM Appointment
    WHERE user_id = p_user_id
        AND YEAR(appointment_date) = YEAR(p_month)
        AND MONTH(appointment_date) = MONTH(p_month);
    
    -- 3. 计算挑战参与统计
    SELECT 
        COUNT(*) AS total,
        SUM(CASE WHEN p.status = '已完成' THEN 1 ELSE 0 END) AS completed
    INTO v_total_challenges, v_completed_challenges
    FROM Participation p
    JOIN Challenge c ON p.challenge_id = c.challenge_id
    WHERE p.user_id = p_user_id
        AND YEAR(c.start_date) = YEAR(p_month)
        AND MONTH(c.start_date) = MONTH(p_month);
    
    -- 4. 获取健康指标统计
    SET v_health_summary = (
        SELECT JSON_OBJECT(
            'weight_stats', (
                SELECT JSON_OBJECT(
                    'avg', IFNULL(ROUND(AVG(metric_value), 2), 0),
                    'min', IFNULL(ROUND(MIN(metric_value), 2), 0),
                    'max', IFNULL(ROUND(MAX(metric_value), 2), 0)
                )
                FROM HealthMetric
                WHERE user_id = p_user_id
                    AND metric_type = '体重'
                    AND YEAR(measured_at) = YEAR(p_month)
                    AND MONTH(measured_at) = MONTH(p_month)
            ),
            'blood_pressure_stats', (
                SELECT JSON_OBJECT(
                    'systolic_avg', IFNULL(ROUND(AVG(CASE WHEN metric_type = '血压收缩压' THEN metric_value END), 2), 0),
                    'diastolic_avg', IFNULL(ROUND(AVG(CASE WHEN metric_type = '血压舒张压' THEN metric_value END), 2), 0)
                )
                FROM HealthMetric
                WHERE user_id = p_user_id
                    AND metric_type IN ('血压收缩压', '血压舒张压')
                    AND YEAR(measured_at) = YEAR(p_month)
                    AND MONTH(measured_at) = MONTH(p_month)
            ),
            'step_stats', (
                SELECT JSON_OBJECT(
                    'total', IFNULL(SUM(metric_value), 0),
                    'daily_avg', IFNULL(ROUND(AVG(metric_value), 2), 0)
                )
                FROM HealthMetric
                WHERE user_id = p_user_id
                    AND metric_type = '步数'
                    AND YEAR(measured_at) = YEAR(p_month)
                    AND MONTH(measured_at) = MONTH(p_month)
            ),
            'generated_at', NOW()
        )
    );
    
    -- 5. 插入月度报告
    INSERT INTO MonthlyReport (
        user_id, report_month, 
        total_appointments, completed_appointments, cancelled_appointments,
        total_challenges, completed_challenges,
        health_summary, recommendations
    ) VALUES (
        p_user_id, p_month,
        COALESCE(v_total_appointments, 0), 
        COALESCE(v_completed_appointments, 0),
        COALESCE(v_cancelled_appointments, 0),
        COALESCE(v_total_challenges, 0),
        COALESCE(v_completed_challenges, 0),
        v_health_summary,
        '请继续保持健康生活习惯，定期监测健康指标。'
    );
    
    SELECT '月度健康报告生成成功' AS message;
END;
//
DELIMITER ;

-- 存储过程2：查找最活跃用户
DELIMITER //
CREATE PROCEDURE FindMostActiveUsers(
    IN p_limit INT
)
BEGIN
    -- 基于健康记录数量、挑战完成情况和预约数量综合评分
    SELECT 
        u.user_id,
        u.health_id,
        u.name,
        COUNT(DISTINCT hm.metric_id) AS health_record_count,
        COUNT(DISTINCT CASE WHEN p.status = '已完成' THEN p.participation_id END) AS completed_challenges,
        COUNT(DISTINCT a.appointment_id) AS appointment_count,
        (COUNT(DISTINCT hm.metric_id) * 0.4 + 
         COUNT(DISTINCT CASE WHEN p.status = '已完成' THEN p.participation_id END) * 0.3 + 
         COUNT(DISTINCT a.appointment_id) * 0.3) AS activity_score
    FROM User u
    LEFT JOIN HealthMetric hm ON u.user_id = hm.user_id
    LEFT JOIN Participation p ON u.user_id = p.user_id
    LEFT JOIN Appointment a ON u.user_id = a.user_id
    WHERE hm.measured_at >= DATE_SUB(NOW(), INTERVAL 3 MONTH)  -- 最近3个月的数据
       OR p.joined_at >= DATE_SUB(NOW(), INTERVAL 3 MONTH)
       OR a.appointment_date >= DATE_SUB(NOW(), INTERVAL 3 MONTH)
    GROUP BY u.user_id, u.health_id, u.name
    ORDER BY activity_score DESC
    LIMIT p_limit;
END;
//
DELIMITER ;

-- 存储过程3：查找参与人数最多的健康挑战
DELIMITER //
CREATE PROCEDURE FindMostPopularChallenges(
    IN p_limit INT
)
BEGIN
    SELECT 
        c.challenge_id,
        c.challenge_name,
        c.challenge_type,
        c.start_date,
        c.end_date,
        COUNT(p.user_id) AS participant_count,
        ROUND(AVG(p.current_progress), 2) AS avg_progress,
        SUM(CASE WHEN p.status = '已完成' THEN 1 ELSE 0 END) AS completed_count
    FROM Challenge c
    LEFT JOIN Participation p ON c.challenge_id = p.challenge_id
    WHERE c.status IN ('进行中', '已结束')
    GROUP BY c.challenge_id, c.challenge_name, c.challenge_type, c.start_date, c.end_date
    ORDER BY participant_count DESC, avg_progress DESC
    LIMIT p_limit;
END;
//
DELIMITER ;

-- 输出确认信息
SELECT '==========================================' AS '';
SELECT 'HealthTrack数据库创建完成' AS 消息;
SELECT '==========================================' AS '';
SELECT '创建的表数量:' AS 统计项, COUNT(*) AS 数量 FROM information_schema.tables WHERE table_schema = 'HealthTrackDB';
SELECT '创建的存储过程数量:' AS 统计项, COUNT(*) AS 数量 FROM information_schema.routines WHERE routine_schema = 'HealthTrackDB' AND routine_type = 'PROCEDURE';
SELECT '创建的触发器数量:' AS 统计项, COUNT(*) AS 数量 FROM information_schema.triggers WHERE trigger_schema = 'HealthTrackDB';
SELECT '==========================================' AS '';