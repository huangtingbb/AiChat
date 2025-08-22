-- 使用数据库
USE chatbot;

-- AI模型表
CREATE TABLE `ai_model` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '模型名称，如zhipu-glm-4',
  `display_name` varchar(100) COLLATE utf8mb4_general_ci NOT NULL COMMENT '显示名称，如智谱GLM-4',
  `provider` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '提供商，如zhipu、openai',
  `type` varchar(20) COLLATE utf8mb4_general_ci NOT NULL COMMENT '类型，如chat、image',
  `url` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '请求URL',
  `max_tokens` int(11) DEFAULT '2048' COMMENT '最大Token数',
  `temperature` float DEFAULT '0.7' COMMENT '温度参数',
  `top_p` float DEFAULT '0.9' COMMENT 'Top-P参数',
  `presence_penalty` float DEFAULT '0' COMMENT '重复惩罚',
  `frequency_penalty` float DEFAULT '0' COMMENT '频率惩罚',
  `enabled` tinyint(1) DEFAULT '1' COMMENT '是否启用',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否为默认模型',
  `api_parameters` json DEFAULT NULL COMMENT 'API参数(JSON格式)',
  `description` text COLLATE utf8mb4_general_ci COMMENT '模型描述',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `class` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '大分类 workflow、bot、bigmodal',
  `class_id` varchar(25) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'coze大分类的id,对应workflow_id,bot_id等',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- AI模型使用记录表
CREATE TABLE IF NOT EXISTS ai_model_usage (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,
    model_id INT UNSIGNED NOT NULL,
    message_id INT UNSIGNED,
    prompt TEXT COMMENT '提问内容',
    response TEXT COMMENT '响应内容',
    prompt_tokens INT DEFAULT 0 COMMENT '提问Token数',
    completion_tokens INT DEFAULT 0 COMMENT '回复Token数',
    total_tokens INT DEFAULT 0 COMMENT '总Token数',
    duration INT DEFAULT 0 COMMENT '耗时(毫秒)',
    status VARCHAR(20) NOT NULL COMMENT '状态: success, error',
    error_msg TEXT COMMENT '错误信息',
    cost FLOAT DEFAULT 0 COMMENT '计费金额',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_model_id (model_id),
    INDEX idx_message_id (message_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- 插入默认AI模型数据
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (13, 'glm-4-air-250414', 'glm-4-air-250414', 'zhipu', 'chat', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, '提供了复杂推理、超长上下文、极快推理速度等多款模型，适用于多种应用场景', '2025-05-30 16:28:15', '2025-08-22 11:11:55', NULL, 'bigmodal', 'glm-4-air-250414');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (14, 'glm-4-airx', 'glm-4-airx', 'zhipu', 'chat', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, '提供了复杂推理、超长上下文、极快推理速度等多款模型，适用于多种应用场景', '2025-05-30 16:28:15', '2025-08-22 11:11:55', NULL, 'bigmodal', 'glm-4-airx');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (15, 'glm-4-long', 'glm-4-long', 'zhipu', 'chat', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, '提供了复杂推理、超长上下文、极快推理速度等多款模型，适用于多种应用场景', '2025-05-30 16:28:15', '2025-08-22 11:11:55', NULL, 'bigmodal', 'glm-4-long');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (16, 'glm-4-flashx', 'glm-4-flashx', 'zhipu', 'chat', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, '提供了复杂推理、超长上下文、极快推理速度等多款模型，适用于多种应用场景', '2025-05-30 16:28:15', '2025-08-22 11:11:55', NULL, 'bigmodal', 'glm-4-flashx');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (17, 'glm-4-flash-250414', 'glm-4-flash-250414', 'zhipu', 'chat', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, '提供了复杂推理、超长上下文、极快推理速度等多款模型，适用于多种应用场景', '2025-05-30 16:28:15', '2025-08-22 11:11:55', NULL, 'bigmodal', 'glm-4-flash-250414');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (18, 'glm-z1-air', 'glm-z1-air', 'zhipu', 'chat', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, '提供了复杂推理、超长上下文、极快推理速度等多款模型，适用于多种应用场景', '2025-05-30 16:28:15', '2025-08-22 11:11:55', NULL, 'bigmodal', 'glm-z1-air');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (19, 'glm-z1-airx', 'glm-z1-airx', 'zhipu', 'chat', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, '提供了复杂推理、超长上下文、极快推理速度等多款模型，适用于多种应用场景', '2025-05-30 16:28:15', '2025-08-22 11:11:55', NULL, 'bigmodal', 'glm-z1-airx');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (20, 'glm-z1-flash', 'glm-z1-flash', 'zhipu', 'chat', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, '提供了复杂推理、超长上下文、极快推理速度等多款模型，适用于多种应用场景', '2025-05-30 16:28:15', '2025-08-22 11:11:55', NULL, 'bigmodal', 'glm-z1-flash');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (22, 'glm-4v-plus-0111', 'glm-4v-plus-0111', 'zhipu', 'chat', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, '在不牺牲任何NLP任务性能的情况下，实现了视觉语言特征的深度融合；支持视觉问答、图像字幕、视觉定位、复杂目标检测等各类图像/视频理解任务', '2025-05-30 16:28:15', '2025-08-22 11:11:55', NULL, 'bigmodal', 'glm-4v-plus-0111');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (23, 'glm-4v-flash', 'glm-4v-flash', 'zhipu', 'chat', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, '在不牺牲任何NLP任务性能的情况下，实现了视觉语言特征的深度融合；支持视觉问答、图像字幕、视觉定位、复杂目标检测等各类图像/视频理解任务', '2025-05-30 16:28:15', '2025-08-22 11:11:55', NULL, 'bigmodal', 'glm-4v-flash');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (24, 'glm-4-plus', 'glm-4-plus', 'zhipu', 'chat', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', 2048, 0.7, 0.9, 0, 0, 1, 1, NULL, '提供了复杂推理、超长上下文、极快推理速度等多款模型，适用于多种应用场景', '2025-05-30 16:28:15', '2025-08-22 11:11:55', NULL, 'bigmodal', 'glm-4-plus');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (25, 'ai-customer', 'Ai智能客服', 'coze', 'chat', 'https://api.coze.cn', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, NULL, '2025-08-22 10:54:53', '2025-08-22 11:50:25', NULL, 'bot', '7523118281046458395');
INSERT INTO `chatbot`.`ai_model` (`id`, `name`, `display_name`, `provider`, `type`, `url`, `max_tokens`, `temperature`, `top_p`, `presence_penalty`, `frequency_penalty`, `enabled`, `is_default`, `api_parameters`, `description`, `created_at`, `updated_at`, `deleted_at`, `class`, `class_id`) VALUES (26, 'ai-customer-workflow', 'Ai智能客服工作流', 'coze', 'chat', 'https://api.coze.cn', 2048, 0.7, 0.9, 0, 0, 1, 0, NULL, NULL, '2025-08-22 11:49:24', '2025-08-22 11:50:05', NULL, 'workflow', '7534632837212307495');

