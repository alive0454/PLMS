-- 为 person 表添加新字段
-- 是否党员、入党日、民族、学历、党员备注

-- 添加是否党员字段（0=否，1=是）
ALTER TABLE person 
ADD COLUMN is_cp TINYINT(1) DEFAULT 0 COMMENT '是否党员：0否，1是' 
AFTER other_info;

-- 添加入党日字段（日期类型，可为空）
ALTER TABLE person 
ADD COLUMN cp_joining_day DATE NULL COMMENT '入党日' 
AFTER is_cp;

-- 添加民族字段
ALTER TABLE person 
ADD COLUMN nationality VARCHAR(50) NULL COMMENT '民族' 
AFTER cp_joining_day;

-- 添加学历字段
ALTER TABLE person 
ADD COLUMN education VARCHAR(50) NULL COMMENT '学历' 
AFTER nationality;

-- 添加党员备注字段
ALTER TABLE person 
ADD COLUMN cp_remark TEXT NULL COMMENT '党员备注' 
AFTER education;

-- 验证字段添加结果
DESCRIBE person;
