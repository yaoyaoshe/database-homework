-- ============================================================
-- HealthTrack 演示数据填充脚本 (Demo Data)
-- ============================================================

USE HealthTrackDB;

-- 为了重新填充，先清空表（注意外键顺序）
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE ChallengeDailyProgress;
TRUNCATE TABLE MonthlyReport;
TRUNCATE TABLE Invitation;
TRUNCATE TABLE Participation;
TRUNCATE TABLE Challenge;
TRUNCATE TABLE HealthMetric;
TRUNCATE TABLE Appointment;
TRUNCATE TABLE UserFamily;
TRUNCATE TABLE UserProvider;
TRUNCATE TABLE UserPhone;
TRUNCATE TABLE Email;
TRUNCATE TABLE Provider;
TRUNCATE TABLE User;
SET FOREIGN_KEY_CHECKS = 1;

-- ==================== 1. 用户数据 ====================
-- 密码统一为 '123456' 的 SHA256 哈希值: 8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92

-- 用户 1: 主演示账号
INSERT INTO User (user_id, health_id, name, date_of_birth, gender, password_hash) VALUES 
(1, 'HID-DEMO', '演示用户', '1995-05-20', '男', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92');

-- 用户 2: 家人账号（用于演示家庭组）
INSERT INTO User (user_id, health_id, name, date_of_birth, gender, password_hash) VALUES 
(2, 'HID-MOM', '妈妈', '1968-03-15', '女', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92');

-- 用户 3-5: 其他活跃用户（用于演示“活跃用户排行榜”）
INSERT INTO User (user_id, health_id, name, date_of_birth, gender, password_hash) VALUES 
(3, 'HID-ALICE', 'Alice', '1990-01-01', '女', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92'),
(4, 'HID-BOB', 'Bob', '1988-11-11', '男', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92'),
(5, 'HID-CHARLIE', 'Charlie', '2000-07-07', '男', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92');

-- ==================== 2. 联系方式 ====================
INSERT INTO Email (user_id, email_address, is_verified, is_primary) VALUES 
(1, 'demo@healthtrack.com', TRUE, TRUE),
(2, 'mom@healthtrack.com', TRUE, TRUE);

INSERT INTO UserPhone (user_id, phone_number, phone_type, is_verified, is_primary) VALUES 
(1, '13800138000', '手机', TRUE, TRUE);

-- ==================== 3. 医疗提供者 ====================
INSERT INTO Provider (provider_id, license_number, name, specialty, qualification, is_verified, contact_info) VALUES 
(1, 'DOC001', '张医生', '内科', '主任医师', TRUE, '{"科室": "心血管内科", "医院": "第一人民医院"}'),
(2, 'DOC002', '李医生', '外科', '副主任医师', TRUE, '{"科室": "骨科", "医院": "第二附属医院"}'),
(3, 'DOC003', '王心理', '心理咨询', '注册心理师', TRUE, '{"科室": "心理健康中心", "医院": "社区医院"}');

-- 关联: 演示用户 关联 张医生 (家庭医生)
INSERT INTO UserProvider (user_id, provider_id, relationship_type, link_date) VALUES 
(1, 1, '家庭医生', DATE_SUB(NOW(), INTERVAL 1 YEAR));

-- 设置主治医生
UPDATE User SET primary_provider_id = 1 WHERE user_id = 1;

-- ==================== 4. 家庭关系 ====================
-- 演示用户 <-> 妈妈 (已验证)
INSERT INTO UserFamily (user_id, related_user_id, relationship, is_verified, verified_at) VALUES 
(1, 2, '父母', TRUE, NOW()),
(2, 1, '子女', TRUE, NOW());

-- Alice 请求添加 演示用户 为朋友 (未验证，用于演示“待处理请求”)
INSERT INTO UserFamily (user_id, related_user_id, relationship, is_verified) VALUES 
(3, 1, '朋友', FALSE);

-- ==================== 5. 预约记录 ====================
-- 过去: 已完成的预约 (用于月度报表统计)
INSERT INTO Appointment (user_id, provider_id, appointment_date, status, consultation_type, reason) VALUES 
(1, 1, DATE_SUB(NOW(), INTERVAL 5 DAY), '已完成', '线下就诊', '定期体检'),
(1, 2, DATE_SUB(NOW(), INTERVAL 10 DAY), '已完成', '电话咨询', '腿部疼痛咨询');

-- 未来: 待参加的预约 (用于“我的预约”展示)
INSERT INTO Appointment (user_id, provider_id, appointment_date, status, consultation_type, reason) VALUES 
(1, 1, DATE_ADD(NOW(), INTERVAL 2 DAY), '已预约', '线下就诊', '复查血压'),
(1, 3, DATE_ADD(NOW(), INTERVAL 1 WEEK), '已预约', '线上咨询', '睡眠问题咨询');

-- ==================== 6. 健康指标 (本月数据) ====================
-- 必须包含本月数据，以便 MonthlySummary.vue 默认显示内容
INSERT INTO HealthMetric (user_id, metric_type, metric_value, unit, measured_at, source) VALUES 
(1, '体重', 70.5, 'kg', DATE_SUB(NOW(), INTERVAL 1 DAY), '健康设备'),
(1, '体重', 71.0, 'kg', DATE_SUB(NOW(), INTERVAL 10 DAY), '手动录入'),
(1, '步数', 8500, '步', DATE_SUB(NOW(), INTERVAL 0 DAY), '手机APP'),
(1, '步数', 12000, '步', DATE_SUB(NOW(), INTERVAL 1 DAY), '手机APP'),
(1, '步数', 6000, '步', DATE_SUB(NOW(), INTERVAL 2 DAY), '手机APP'),
(1, '血压收缩压', 120, 'mmHg', NOW(), '手动录入'),
(1, '血压舒张压', 80, 'mmHg', NOW(), '手动录入');

-- ==================== 7. 健康挑战与参与 ====================
-- 挑战 A: 热门挑战 (大家都在参加)
INSERT INTO Challenge (challenge_id, creator_id, challenge_name, description, challenge_type, target_metric, target_value, target_unit, start_date, end_date, status, max_participants) VALUES 
(1, 1, '30天万步挑战', '每天坚持走一万步，保持健康活力！', '运动', '步数', 10000, '步', DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_ADD(NOW(), INTERVAL 25 DAY), '进行中', 50);

-- 参与情况
INSERT INTO Participation (challenge_id, user_id, status, current_progress) VALUES 
(1, 1, '参与中', 25000), -- 演示用户
(1, 3, '参与中', 30000), -- Alice
(1, 4, '参与中', 10000), -- Bob
(1, 5, '参与中', 5000);  -- Charlie

-- 挑战 B: 演示用户已完成的挑战 (用于报表统计)
INSERT INTO Challenge (challenge_id, creator_id, challenge_name, challenge_type, target_metric, target_value, target_unit, start_date, end_date, status) VALUES 
(2, 2, '早睡早起周', '睡眠', '睡眠时长', 8, '小时', DATE_SUB(NOW(), INTERVAL 10 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), '已结束');

INSERT INTO Participation (challenge_id, user_id, status, current_progress) VALUES 
(1, 2, '参与中', 20000),
(2, 1, '已完成', 56); -- 7天 * 8小时 = 56

-- ==================== 8. 挑战邀请 ====================
-- Bob 邀请 演示用户 参加一个新挑战
INSERT INTO Challenge (challenge_id, creator_id, challenge_name, challenge_type, target_metric, target_value, target_unit, start_date, end_date, status) VALUES 
(3, 4, '减脂突击营', '减重', '体重', 65, 'kg', NOW(), DATE_ADD(NOW(), INTERVAL 1 MONTH), '进行中');

-- 发送邀请 (Health ID 方式)
INSERT INTO Invitation (challenge_id, sender_id, recipient_type, recipient_value, status, message) VALUES 
(3, 4, '用户ID', '1', '待处理', '哥们，一起来减脂吧！');

-- ==================== 9. 用户活跃记录 (模拟) ====================
-- 确保活跃用户排行榜有数据
INSERT INTO UserActivity (user_id, activity_type, performed_at) VALUES (1, '登录', NOW());
INSERT INTO UserActivity (user_id, activity_type, performed_at) VALUES (3, '登录', NOW());
INSERT INTO UserActivity (user_id, activity_type, performed_at) VALUES (4, '登录', NOW());

SELECT '演示数据加载完成！' AS Status;