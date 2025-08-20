-- 添加Coze智能体模型
-- 003_coze_models.sql

-- 插入Coze智能体模型
INSERT INTO ai_models (
    name, 
    display_name, 
    provider, 
    description, 
    max_tokens, 
    temperature, 
    top_p, 
    presence_penalty, 
    frequency_penalty, 
    is_enabled, 
    sort_order
) VALUES 
(
    'coze-bot', 
    'Coze智能体', 
    'coze', 
    'Coze平台的智能体，支持多种功能和工作流', 
    4096, 
    0.7, 
    0.9, 
    0.0, 
    0.0, 
    1, 
    30
),
(
    'coze-workflow', 
    'Coze工作流', 
    'coze', 
    'Coze平台的工作流模式，适合复杂的多步骤任务', 
    4096, 
    0.7, 
    0.9, 
    0.0, 
    0.0, 
    1, 
    31
);

-- 添加索引优化查询性能
CREATE INDEX IF NOT EXISTS idx_ai_models_provider ON ai_models(provider);
CREATE INDEX IF NOT EXISTS idx_ai_models_enabled ON ai_models(is_enabled);
