<template>
  <div class="chat-container">
    <!-- 侧边栏 -->
    <div class="sidebar">
      <div class="sidebar-header">
        <div class="logo-container">
          <div class="small-logo-circle"></div>
          <h3>AI聊天助手</h3>
        </div>
        <el-button type="primary" size="small" @click="createTempChat" :loading="chatStore.loading" class="new-chat-button">
          新建对话
        </el-button>
      </div>
      
      <div class="chat-list">
        <div
          v-if="chatStore.tempChat && chatStore.currentChatId === 'temp'"
          class="chat-item active temp-chat"
        >
          <div class="chat-icon">
            <el-icon><ChatDotRound /></el-icon>
          </div>
          <div class="chat-info">
            <span class="chat-title">{{ chatStore.tempChat.title }}</span>
            <span class="chat-time">{{ formatTime(chatStore.tempChat.created_at) }}</span>
          </div>
        </div>
        <div
          v-for="chat in chatStore.chatList"
          :key="chat.id"
          class="chat-item"
          :class="{ active: chatStore.currentChatId === chat.id }"
          @click="selectChat(chat.id)"
        >
          <div class="chat-icon">
            <el-icon><ChatDotRound /></el-icon>
          </div>
          <div class="chat-info">
            <span class="chat-title">{{ chat.title }}</span>
            <span class="chat-time">{{ formatTime(chat.created_at) }}</span>
          </div>
          <div class="chat-actions">
            <el-button 
              type="link" 
              size="small" 
              @click.stop="editChatTitle(chat)"
              class="edit-action"
            >
              <el-icon><Edit /></el-icon>
            </el-button>
            <el-button 
              type="link" 
              size="small" 
              @click.stop="confirmDeleteChat(chat.id)"
              class="delete-action"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
        <div v-if="chatStore.chatList.length === 0" class="empty-list">
          <div class="empty-icon">
            <el-icon><ChatLineSquare /></el-icon>
          </div>
          <p>没有聊天记录</p>
          <p>点击"新建对话"开始聊天</p>
        </div>
      </div>
      
      <div class="sidebar-footer">
        <div class="user-info">
          <div class="user-avatar">{{ userStore.username.charAt(0).toUpperCase() }}</div>
          <span class="username">{{ userStore.username }}</span>
        </div>
        <el-button type="danger" size="small" @click="logout" class="logout-button">
          <el-icon><Close /></el-icon> 退出登录
        </el-button>
      </div>
    </div>
    
    <!-- 聊天主区域 -->
    <div class="chat-main">
      <!-- 主区域头部 -->
      <div class="chat-header" v-if="chatStore.currentChatId">
        <!-- 聊天标题 (带编辑功能) -->
        <div class="title-container">
          <h3 v-if="!isEditingTitle" @click="startEditTitle" class="editable-title">
            {{ chatStore.currentChat?.title || '新对话' }}
            <el-icon class="edit-icon"><Edit /></el-icon>
          </h3>
          <div v-else class="title-edit-form">
            <el-input
              v-model="editingTitle"
              size="small"
              @keyup.enter="saveTitle"
              @blur="saveTitle"
              ref="titleInput"
              class="title-input"
              placeholder="输入聊天标题"
            ></el-input>
          </div>
        </div>
        <div class="header-actions">
          <el-button type="link" class="more-actions">
            <el-icon><More /></el-icon>
          </el-button>
        </div>
      </div>
      
      <div v-if="!chatStore.currentChatId" class="welcome">
        <div class="tech-circle"></div>
        <h2>欢迎使用AI聊天助手</h2>
        <p>点击左侧"新建对话"或在下方输入框中发送消息开始一段新的对话</p>
      </div>
      
      <template v-else>
        <!-- 聊天消息区域 -->
        <div class="chat-messages" ref="messagesContainer">
          <div v-if="chatStore.messages.length === 0" class="empty-messages">
            <div class="empty-messages-icon">
              <el-icon><Message /></el-icon>
            </div>
            <p>没有消息，在下方输入框中发送消息开始聊天</p>
          </div>
          
          <div
            v-for="(message, index) in chatStore.messages"
            :key="message.id"
            class="message-container"
          >
            <!-- 如果是助手消息，显示头像 -->
            <div v-if="message.role === 'assistant'" class="avatar assistant-avatar">AI</div>
            
            <div
              class="message"
              :class="[
                message.role, 
                { 'continued': index > 0 && chatStore.messages[index - 1].role === message.role }
              ]"
            >
              <div class="message-content" v-html="formatMessage(message.content)"></div>
              <div class="message-time">{{ formatTime(message.created_at) }}</div>
            </div>
            
            <!-- 如果是用户消息，显示头像 -->
            <div v-if="message.role === 'user'" class="avatar user-avatar">
              {{ userStore.username.charAt(0).toUpperCase() }}
            </div>
          </div>
        </div>
      </template>
        
      <!-- 输入区域 - 现在放在主区域的底部，无论是否有活动聊天 -->
      <div class="chat-input-container">
        <div class="chat-input-wrapper">
          <el-input
            v-model="messageInput"
            type="textarea"
            :rows="1"
            :autosize="{ minRows: 1, maxRows: 5 }"
            :placeholder="'在 go 中发消息'"
            @keydown.enter.exact.prevent="sendMessage"
            class="message-textarea"
          />
          
          <!-- 模型选择和操作按钮 -->
          <div class="input-model-selection">
            <div class="model-options">
              <el-dropdown @command="handleModelSelect" trigger="click" placement="top">
                <div class="model-option active dropdown-trigger">
                  {{ selectedModel ? selectedModel.name : 'AI模型' }}
                  <span class="dropdown-label">
                    <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
                  </span>
                </div>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item 
                      v-for="model in modelList" 
                      :key="model.id" 
                      :command="model"
                      :disabled="selectedModel && selectedModel.id === model.id"
                    >
                      {{ model.name }}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
            
            <div class="input-actions">
              <el-button-group class="action-buttons">
                <el-button type="link" class="action-button">
                  <el-icon><FullScreen /></el-icon>
                </el-button>
                <el-button type="link" class="action-button">
                  <el-icon><Upload /></el-icon>
                </el-button>
              </el-button-group>
              <el-button
                type="primary"
                circle
                :loading="chatStore.sending"
                @click="sendMessage"
                :disabled="!messageInput.trim()"
                class="send-button"
              >
                <el-icon><Position /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
        <div class="input-features">
          <span class="feature-hint">内容由AI生成，仅供参考</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useUserStore } from '../store/user'
import { useChatStore } from '../store/chat'
import { ElMessageBox, ElLoading } from 'element-plus'
import { marked } from 'marked'
import { GetAvailableModelList } from '../api/chat'
// 导入Element Plus图标
import {
  Edit,
  Delete,
  ChatDotRound,
  ChatLineSquare,
  Message,
  More,
  Close,
  FullScreen,
  Upload,
  Position,
  ArrowDown
} from '@element-plus/icons-vue'

const userStore = useUserStore()
const chatStore = useChatStore()

const messageInput = ref('')
const messagesContainer = ref(null)
const isEditingTitle = ref(false)
const editingTitle = ref('')
const titleInput = ref(null)

// 模型相关
const availableModels = ref([]) // 不再预设默认模型
const selectedModel = ref(null) // 初始化为null，等待接口返回
const modelList = ref([]) // 从服务器获取的模型列表

// 选择模型的方法
const selectModel = (model) => {
  selectedModel.value = model
}

// 处理下拉菜单选择模型
const handleModelSelect = (model) => {
  selectModel(model)
}

// 获取可用的模型
const fetchAvailableModels = async () => {
  const loading = ElLoading.service({
    lock: true,
    text: '加载模型中...',
    background: 'rgba(255, 255, 255, 0.6)'
  })
  
  try {
    const response = await GetAvailableModelList()
    
    // 处理后端返回的模型数据
    if (response && response.data.models && response.data.models.length > 0) {
      // 将后端模型数据转换为前端需要的格式
      modelList.value = response.data.models.map(model => {
        return {
          id: model.id,
          name: model.name,
          description: model.description
        }
      })
      
      // 只使用接口返回的current_model作为默认选中模型
      if (response.data.current_model) {
        const currentModel = modelList.value.find(m => m.id === response.data.current_model)
        if (currentModel) {
          selectedModel.value = currentModel
        }
      }
    } else {
      console.warn('接口未返回有效的模型数据')
    }
  } catch (error) {
    console.error('获取可用模型失败:', error)
  } finally {
    loading.close()
  }
}

// 格式化消息内容（支持Markdown）
const formatMessage = (content) => {
  return marked(content)
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleString()
}

// 滚动到底部
const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// 选择聊天
const selectChat = async (chatId) => {
  await chatStore.selectChat(chatId)
  scrollToBottom()
}

// 创建新聊天
const createTempChat = () => {
  chatStore.createTempChat()
  scrollToBottom()
}

// 开始编辑标题
const startEditTitle = () => {
  if (!chatStore.currentChat) return
  
  editingTitle.value = chatStore.currentChat.title
  isEditingTitle.value = true
  
  // 等待DOM更新后聚焦输入框
  nextTick(() => {
    if (titleInput.value) {
      titleInput.value.focus()
    }
  })
}

// 保存标题
const saveTitle = async () => {
  if (isEditingTitle.value && chatStore.currentChatId) {
    const newTitle = editingTitle.value.trim()
    if (newTitle && newTitle !== chatStore.currentChat?.title) {
      await chatStore.updateChatTitle(chatStore.currentChatId, newTitle)
    }
    isEditingTitle.value = false
  }
}

// 发送消息
const sendMessage = async () => {
  if (!messageInput.value.trim() || chatStore.sending) return
  
  const content = messageInput.value
  messageInput.value = ''
  
  // 只有在没有当前聊天时，才创建临时聊天
  if (!chatStore.currentChatId) {
    createTempChat()
  }
  
  // 将选中的模型ID传递给store
  const modelId = selectedModel.value ? selectedModel.value.id : null
  await chatStore.sendUserMessage(content, modelId)
  scrollToBottom()
}

// 退出登录
const logout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    userStore.logout()
  }).catch(() => {})
}

// 监听消息变化，自动滚动到底部
watch(() => chatStore.messages.length, scrollToBottom)

// 组件挂载时获取聊天列表和可用模型
onMounted(async () => {
  await chatStore.fetchChatList()
  await fetchAvailableModels()
})
</script>

<style scoped>
.chat-container {
  height: 100vh;
  display: flex;
  overflow: hidden;
  background-color: var(--dark-bg);
  color: var(--text-color);
  font-family: 'Roboto', 'Arial', sans-serif;
}

/* ===== 侧边栏样式 ===== */
.sidebar {
  width: 280px;
  background-color: var(--medium-bg);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 10px rgba(0, 0, 0, 0.2);
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: linear-gradient(to right, var(--medium-bg), var(--light-bg));
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.small-logo-circle {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(61, 90, 254, 0.2), rgba(0, 229, 255, 0.2));
  border: 1px solid rgba(0, 229, 255, 0.3);
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
}

.small-logo-circle::after {
  content: 'AI';
  font-size: 12px;
  font-weight: bold;
  color: var(--accent-color);
  text-shadow: 0 0 5px rgba(0, 229, 255, 0.5);
}

.sidebar-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--accent-color);
  text-shadow: 0 0 5px rgba(0, 229, 255, 0.5);
  letter-spacing: 0.5px;
}

.new-chat-button {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  height: 38px;
  border-radius: 8px;
  font-weight: 500;
}

.chat-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.chat-list::-webkit-scrollbar {
  width: 5px;
}

.chat-list::-webkit-scrollbar-track {
  background: var(--medium-bg);
}

.chat-list::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

.chat-item {
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 12px;
  background-color: var(--light-bg);
  border-left: 3px solid transparent;
}

.chat-icon {
  width: 36px;
  height: 36px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 8px;
  background: rgba(61, 90, 254, 0.15);
  color: var(--accent-color);
  font-size: 16px;
}

.chat-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.chat-item:hover {
  background-color: rgba(61, 90, 254, 0.2);
  transform: translateX(3px);
  box-shadow: var(--glow-effect);
}

.chat-item.active {
  background-color: rgba(61, 90, 254, 0.3);
  border-left: 3px solid var(--accent-color);
  box-shadow: var(--glow-effect);
}

.chat-title {
  font-size: 14px;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-time {
  font-size: 12px;
  color: var(--muted-text);
}

.empty-list {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: var(--muted-text);
  font-size: 14px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  color: var(--accent-color);
  opacity: 0.6;
}

.empty-list p {
  margin: 5px 0;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  gap: 12px;
  background-color: var(--medium-bg);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 5px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary-color), var(--accent-color));
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: bold;
  font-size: 14px;
  box-shadow: 0 0 10px rgba(0, 229, 255, 0.3);
}

.username {
  font-size: 14px;
  color: var(--text-color);
  max-width: 180px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.logout-button {
  padding: 6px 12px;
  border-radius: 4px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 5px;
}

.logout-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(255, 0, 0, 0.2);
}

/* ===== 主区域样式 ===== */
.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: var(--dark-bg);
  position: relative;
}

.chat-main::before {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  width: 100%;
  height: 100%;
  background: 
    radial-gradient(circle at top right, rgba(0, 229, 255, 0.05), transparent 60%),
    radial-gradient(circle at bottom left, rgba(61, 90, 254, 0.05), transparent 60%);
  pointer-events: none;
}

/* 头部 */
.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid var(--border-color);
  background-color: rgba(35, 35, 66, 0.7);
  backdrop-filter: blur(5px);
  z-index: 10;
}

.title-container {
  display: flex;
  align-items: center;
}

.editable-title {
  margin: 0;
  font-size: 16px;
  color: var(--text-color);
  font-weight: 500;
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.editable-title:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.edit-icon {
  font-size: 14px;
  margin-left: 8px;
  color: var(--muted-text);
  opacity: 0;
  transition: all 0.3s ease;
}

.editable-title:hover .edit-icon {
  opacity: 1;
}

.title-edit-form {
  width: 300px;
}

.title-input {
  max-width: 300px;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.more-actions {
  color: var(--muted-text);
  transition: all 0.3s ease;
}

.more-actions:hover {
  color: var(--accent-color);
  transform: translateY(-2px);
}

/* 欢迎页面 */
.welcome {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: var(--muted-text);
  position: relative;
  z-index: 1;
  padding: 0 20px;
}

.tech-circle {
  width: 180px;
  height: 180px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(61, 90, 254, 0.1), rgba(0, 229, 255, 0.1));
  border: 1px solid rgba(0, 229, 255, 0.3);
  margin-bottom: 30px;
  position: relative;
  box-shadow: var(--glow-effect);
  animation: pulse 3s infinite alternate;
}

@keyframes pulse {
  0% { transform: scale(1); }
  100% { transform: scale(1.05); }
}

.tech-circle::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 140px;
  height: 140px;
  transform: translate(-50%, -50%);
  border-radius: 50%;
  border: 1px solid rgba(61, 90, 254, 0.5);
}

.tech-circle::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 100px;
  height: 100px;
  transform: translate(-50%, -50%);
  border-radius: 50%;
  border: 1px solid rgba(0, 229, 255, 0.5);
  animation: spin 20s linear infinite;
}

@keyframes spin {
  0% { transform: translate(-50%, -50%) rotate(0deg); }
  100% { transform: translate(-50%, -50%) rotate(360deg); }
}

.welcome h2 {
  margin-bottom: 16px;
  color: var(--accent-color);
  font-size: 28px;
  text-shadow: 0 0 5px rgba(0, 229, 255, 0.5);
  letter-spacing: 1px;
}

.welcome p {
  text-align: center;
  max-width: 500px;
  line-height: 1.6;
  font-size: 16px;
}

/* 消息区域 */
.chat-messages {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  scroll-behavior: smooth;
}

.chat-messages::-webkit-scrollbar {
  width: 5px;
}

.chat-messages::-webkit-scrollbar-track {
  background: var(--dark-bg);
}

.chat-messages::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

.empty-messages {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: var(--muted-text);
  font-size: 14px;
  text-align: center;
}

.empty-messages-icon {
  font-size: 48px;
  margin-bottom: 16px;
  color: var(--accent-color);
  opacity: 0.6;
}

/* 消息样式 */
.message-container {
  display: flex;
  align-items: flex-start;
  margin-bottom: 24px;
  position: relative;
  width: 100%;
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: bold;
  font-size: 14px;
  flex-shrink: 0;
}

.assistant-avatar {
  background: linear-gradient(135deg, var(--accent-color), var(--primary-color));
  color: white;
  margin-right: 12px;
}

.user-avatar {
  background: linear-gradient(135deg, var(--primary-color), var(--accent-color));
  color: white;
  margin-left: 12px;
}

.message {
  max-width: 70%;
  padding: 14px 18px;
  border-radius: 12px;
  position: relative;
  animation: fadeIn 0.3s ease;
}

.message.continued {
  margin-top: -15px;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.message.user {
  margin-left: auto;
  background: linear-gradient(135deg, rgba(61, 90, 254, 0.2), rgba(61, 90, 254, 0.4));
  color: var(--text-color);
  border-top-right-radius: 4px;
  border: 1px solid rgba(61, 90, 254, 0.5);
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.2);
}

.message.assistant {
  margin-right: auto;
  background: linear-gradient(135deg, rgba(0, 229, 255, 0.1), rgba(0, 229, 255, 0.2));
  color: var(--text-color);
  border-top-left-radius: 4px;
  border: 1px solid rgba(0, 229, 255, 0.3);
  box-shadow: var(--glow-effect);
}

.message-content {
  font-size: 15px;
  line-height: 1.6;
  word-break: break-word;
}

/* 支持代码高亮 */
.message-content :deep(pre) {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 6px;
  padding: 12px;
  overflow-x: auto;
  margin: 8px 0;
}

.message-content :deep(code) {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
}

.message-content :deep(p) {
  margin: 8px 0;
}

.message-time {
  font-size: 11px;
  color: var(--muted-text);
  margin-top: 6px;
  text-align: right;
}

/* 输入区域 */
.chat-input-container {
  padding: 12px 24px 16px;
  background-color: #fff;
  position: relative;
}

.chat-input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
  border-radius: 8px;
  padding: 10px 14px;
  transition: all 0.3s ease;
  position: relative;
  border: 1px solid #e0e0e0;
  background: #fff;
}

.message-textarea {
  width: 100%;
}

.message-textarea :deep(.el-textarea__inner) {
  border: none;
  background-color: transparent;
  resize: none;
  padding: 4px 0;
  font-size: 15px;
  line-height: 1.5;
  color: #333;
  min-height: 24px;
}

.message-textarea :deep(.el-textarea__inner:focus) {
  outline: none;
  box-shadow: none;
}

.message-textarea :deep(.el-textarea__inner::placeholder) {
  color: #999;
}

/* 模型选择区域 */
.input-model-selection {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 4px;
}

.model-options {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 4px;
  scrollbar-width: none; /* Firefox */
}

.model-options::-webkit-scrollbar {
  display: none;
}

.model-option {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: 16px;
  background-color: #f5f5f5;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 14px;
  white-space: nowrap;
  color: #666;
}

.model-option:hover {
  background-color: #eeeeee;
}

.model-option.active {
  background-color: #e6f4ff;
  color: #1677ff;
  border: 1px solid #91caff;
}

/* 下拉标签样式 */
.dropdown-trigger {
  position: relative;
  padding-right: 26px; /* 为下拉图标留出空间 */
}

.dropdown-label {
  display: flex;
  align-items: center;
  justify-content: center;
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 22px;
  background-color: rgba(0, 0, 0, 0.05);
  border-top-right-radius: 16px;
  border-bottom-right-radius: 16px;
  transition: all 0.2s ease;
}

.model-option:hover .dropdown-label {
  background-color: rgba(0, 0, 0, 0.1);
}

.model-option.active .dropdown-label {
  background-color: rgba(22, 119, 255, 0.2);
}

.dropdown-icon {
  font-size: 10px;
  color: inherit;
}

.input-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-buttons {
  display: flex;
}

.action-button {
  color: #999;
  font-size: 16px;
  transition: all 0.2s ease;
  padding: 0 8px;
}

.action-button:hover {
  color: #1677ff;
}

.send-button {
  width: 32px;
  height: 32px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #1677ff;
  border: none;
  transition: all 0.2s ease;
}

.send-button i {
  font-size: 16px;
}

.send-button:hover {
  background: #4096ff;
}

.input-features {
  text-align: center;
  padding-top: 6px;
  font-size: 12px;
  color: #999;
}
</style>