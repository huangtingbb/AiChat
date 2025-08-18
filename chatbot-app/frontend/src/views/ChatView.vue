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
              link 
              size="small" 
              @click.stop="editChatTitle(chat)"
              class="edit-action"
            >
              <el-icon><Edit /></el-icon>
            </el-button>
            <el-button 
              link 
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
          <el-button link class="more-actions">
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
        <div class="chat-messages" ref="messagesContainer" @scroll="detectUserScroll">
          <div v-if="messagesForRender.length === 0" class="empty-messages">
            <div class="empty-messages-icon">
              <el-icon><Message /></el-icon>
            </div>
            <p>没有消息，在下方输入框中发送消息开始聊天</p>
          </div>
          
          <div
            v-for="(message, index) in messagesForRender"
            :key="message._renderKey"
            class="message-container"
          >
            <!-- 如果是助手消息，显示头像 -->
            <div v-if="message.role === 'assistant'" class="avatar assistant-avatar">AI</div>
            
            <div
              class="message"
              :class="[
                message.role, 
                { 
                  'continued': index > 0 && chatStore.messages[index - 1].role === message.role,
                  'streaming': message.isStreaming
                }
              ]"
            >
              <div class="message-content">
                <div v-html="formatMessage(message.content)"></div>
                <!-- 流式响应指示器 -->
                <div v-if="message.isStreaming" class="streaming-indicator">
                  <div class="typing-dots">
                    <span></span>
                    <span></span>
                    <span></span>
                  </div>
                </div>
              </div>
              <div class="message-time">
                {{ formatTime(message.created_at) }}
                <span v-if="message.isStreaming" class="streaming-label">正在输入...</span>
              </div>
            </div>
            
            <!-- 如果是用户消息，显示头像 -->
            <div v-if="message.role === 'user'" class="avatar user-avatar">
              {{ userStore.username.charAt(0).toUpperCase() }}
            </div>
          </div>
        </div>
        
        <!-- 滚动到底部按钮 -->
        <div v-if="isUserScrolling && !messagesForRender.some(msg => msg.isStreaming)" class="scroll-to-bottom-button" @click="forceScrollToBottom(true)">
          <el-button type="primary" circle size="small">
            <el-icon><ArrowDown /></el-icon>
          </el-button>
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
                <el-button link class="action-button">
                  <el-icon><FullScreen /></el-icon>
                </el-button>
                <el-button link class="action-button">
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

// 计算属性：优化消息渲染key生成，避免闪烁
const messagesForRender = computed(() => {
  // 只为内容真正变化的消息生成新的渲染key
  return chatStore.messages.map((msg, index) => {
    // 为每个消息生成稳定的标识符
    // 只有在消息内容确实发生变化时才改变key
    const contentHash = msg.content ? msg.content.length : 0
    const statusFlag = msg.isStreaming ? 'streaming' : 'complete'
    
    return {
      ...msg,
      // 使用更稳定的key生成策略，避免不必要的重新渲染
      _renderKey: `${msg.id}_${contentHash}_${statusFlag}`
    }
  })
})

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

// 配置marked选项
marked.setOptions({
  breaks: true, // 支持换行
  gfm: true, // 启用GitHub风格的Markdown
  sanitize: false, // 允许HTML（注意安全性）
  smartLists: true,
  smartypants: true
})

// 格式化消息内容（支持Markdown）
const formatMessage = (content) => {
  if (!content) return ''
  
  try {
    return marked(content)
  } catch (error) {
    console.error('Markdown渲染错误:', error)
    return content // 如果渲染失败，返回原始内容
  }
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleString()
}

// 滚动到底部 - 优化滚动逻辑和用户体验
let scrollTimer = null
let isUserScrolling = false
let userScrollTimer = null
let isProgrammaticScroll = false // 标识是否为程序触发的滚动

// 检测用户是否主动滚动
const detectUserScroll = () => {
  // 如果是程序触发的滚动，忽略
  if (isProgrammaticScroll) {
    return
  }
  
  // 检查是否有流式消息正在输出
  const hasStreamingMessages = messagesForRender.value.some(msg => msg.isStreaming)
  
  // 如果有流式消息在输出，减少用户滚动检测的敏感度
  if (hasStreamingMessages) {
    console.log('检测到流式消息输出中，降低滚动检测敏感度')
    // 流式输出时，只有用户明显向上滚动才认为是手动滚动
    const container = messagesContainer.value
    if (container) {
      const currentScroll = container.scrollTop
      const maxScroll = container.scrollHeight - container.clientHeight
      // 只有滚动到远离底部的位置才认为是用户手动滚动
      if (currentScroll < maxScroll - 200) {
        console.log('用户明显向上滚动，暂停自动滚动')
        isUserScrolling = true
        if (userScrollTimer) {
          clearTimeout(userScrollTimer)
        }
        // 流式输出时缩短恢复时间
        userScrollTimer = setTimeout(() => {
          isUserScrolling = false
          console.log('恢复自动滚动')
        }, 1000)
      }
    }
  } else {
    // 没有流式消息时使用正常的检测逻辑
    isUserScrolling = true
    if (userScrollTimer) {
      clearTimeout(userScrollTimer)
    }
    // 2秒后重新启用自动滚动
    userScrollTimer = setTimeout(() => {
      isUserScrolling = false
    }, 2000)
  }
}

const scrollToBottom = async (forceScroll = false) => {
  // 如果用户正在手动滚动且不是强制滚动，则不自动滚动
  if (isUserScrolling && !forceScroll) {
    console.log('用户正在滚动，跳过自动滚动')
    return
  }
  
  // 清除之前的定时器，避免频繁滚动
  if (scrollTimer) {
    clearTimeout(scrollTimer)
  }
  
  scrollTimer = setTimeout(async () => {
    await nextTick()
    if (messagesContainer.value) {
      const container = messagesContainer.value
      const currentScroll = container.scrollTop
      const maxScroll = container.scrollHeight - container.clientHeight
      const isNearBottom = currentScroll >= maxScroll - 100
      
      console.log('滚动信息:', {
        currentScroll,
        maxScroll,
        isNearBottom,
        forceScroll,
        scrollHeight: container.scrollHeight,
        clientHeight: container.clientHeight
      })
      
      // 只有在接近底部时才自动滚动，避免打断用户阅读
      if (isNearBottom || forceScroll) {
        // 标记为程序触发的滚动
        isProgrammaticScroll = true
        
        console.log('执行滚动到底部')
        container.scrollTo({
          top: container.scrollHeight,
          behavior: forceScroll ? 'auto' : 'smooth'
        })
        
        // 滚动完成后重置标记
        setTimeout(() => {
          isProgrammaticScroll = false
        }, 500) // 增加延时，确保滚动动画完成
      } else {
        console.log('不满足滚动条件，跳过滚动')
      }
    }
  }, 50) // 50ms的节流延迟
}

// 选择聊天
const selectChat = async (chatId) => {
  await chatStore.selectChat(chatId)
  await nextTick() // 确保DOM更新
  forceScrollToBottom() // 选择聊天后强制滚动到底部
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

// 强制滚动函数 - 支持流式输出时的实时滚动
const forceScrollToBottom = (useSmooth = false) => {
  if (messagesContainer.value) {
    const container = messagesContainer.value
    
    // 检查是否需要滚动
    const currentScroll = container.scrollTop
    const maxScroll = container.scrollHeight - container.clientHeight
    const shouldScroll = currentScroll < maxScroll - 50 // 给一些缓冲空间
    
    if (shouldScroll) {
      console.log('强制滚动前:', {
        scrollTop: container.scrollTop,
        scrollHeight: container.scrollHeight,
        clientHeight: container.clientHeight,
        maxScroll,
        shouldScroll
      })
      
      if (useSmooth) {
        // 使用smooth滚动
        container.scrollTo({
          top: container.scrollHeight,
          behavior: 'smooth'
        })
      } else {
        // 直接设置scrollTop，适用于流式输出
        container.scrollTop = container.scrollHeight
      }
      
      console.log('强制滚动后:', {
        scrollTop: container.scrollTop,
        scrollHeight: container.scrollHeight
      })
    }
  }
}

// 发送消息
const sendMessage = async () => {
  if (!messageInput.value.trim() || chatStore.sending) return
  
  const content = messageInput.value
  messageInput.value = ''
  
  console.log('发送消息开始:', {
    content,
    currentChatId: chatStore.currentChatId,
    messagesLength: messagesForRender.value.length,
    selectedModel: selectedModel.value
  })
  
  // 只有在没有当前聊天时，才创建临时聊天
  if (!chatStore.currentChatId) {
    console.log('创建临时聊天')
    createTempChat()
  }
  
  // 将选中的模型Id传递给store
  const modelId = selectedModel.value ? selectedModel.value.id : null
  console.log('调用sendUserMessage，模型Id:', modelId)
  
  try {
    await chatStore.sendUserMessage(content, modelId)
    console.log('sendUserMessage完成，消息数量:', messagesForRender.value.length)
    
    // 使用nextTick确保DOM更新后再滚动
    await nextTick()
    forceScrollToBottom() // 使用简单滚动函数测试
  } catch (error) {
    console.error('发送消息失败:', error)
  }
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

// 合并监听器，减少重复触发 - 优化性能
let scrollDebounceTimer = null
const debouncedScrollToBottom = (forceScroll = false) => {
  if (scrollDebounceTimer) {
    clearTimeout(scrollDebounceTimer)
  }
  scrollDebounceTimer = setTimeout(() => {
    // 暂时使用简单滚动函数进行测试
    if (forceScroll) {
      forceScrollToBottom()
    } else {
      scrollToBottom(forceScroll)
    }
  }, 100) // 100ms防抖
}

// 监听消息数量变化
watch(() => messagesForRender.value.length, (newLength, oldLength) => {
  if (newLength > oldLength) {
    console.log('消息数量增加:', oldLength, '->', newLength)
    debouncedScrollToBottom(true) // 新消息时强制滚动
  }
}, { flush: 'post' })

// 专门监听流式消息的内容长度变化，确保实时滚动
watch(() => {
  // 获取所有流式消息的总长度
  const streamingMessages = messagesForRender.value.filter(msg => msg.isStreaming)
  return streamingMessages.reduce((total, msg) => total + (msg.content?.length || 0), 0)
}, (newTotalLength, oldTotalLength) => {
  // 只要流式消息内容增加就立即滚动
  if (newTotalLength > oldTotalLength && newTotalLength > 0) {
    console.log('流式内容长度变化:', oldTotalLength, '->', newTotalLength)
    // 立即滚动，不使用防抖，确保跟随性
    forceScrollToBottom()
  }
}, { flush: 'post' })

// 监听流式消息内容变化（确保实时滚动）
let lastStreamingInfo = { count: 0, totalLength: 0 }
watch(() => {
  const streamingMessages = messagesForRender.value.filter(m => m.isStreaming)
  return {
    count: streamingMessages.length,
    totalLength: streamingMessages.reduce((total, msg) => total + (msg.content?.length || 0), 0),
    // 只关注正在流式响应的消息，避免影响其他消息
    streamingIds: streamingMessages.map(m => m.id)
  }
}, (newInfo) => {
  console.log('流式消息变化:', {
    oldLength: lastStreamingInfo.totalLength,
    newLength: newInfo.totalLength,
    oldCount: lastStreamingInfo.count,
    newCount: newInfo.count
  })
  
  // 只有在流式消息确实增加内容时才滚动
  if (newInfo.totalLength > lastStreamingInfo.totalLength || 
      newInfo.count > lastStreamingInfo.count) {
    lastStreamingInfo = newInfo
    console.log('流式消息内容增加，触发滚动')
    
    // 流式输出时使用强制滚动，确保跟随内容
    setTimeout(() => {
      forceScrollToBottom()
    }, 50)
  }
}, { flush: 'post' })

// 监听消息更新计数器（减少日志输出）
// 注释掉这个监听器，因为我们已经不再依赖messageUpdateCount来触发重新渲染
// watch(() => chatStore.messageUpdateCount, (newCount, oldCount) => {
//   if (newCount % 5 === 0) {
//     console.log('消息更新计数器:', newCount, '消息数量:', chatStore.messages.length)
//   }
// }, { flush: 'post' })

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
  /* 优化渲染性能，减少闪屏 */
  will-change: scroll-position;
  transform: translateZ(0); /* 开启硬件加速 */
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
  /* 优化渲染性能，避免闪烁 */
  will-change: auto;
  contain: layout style;
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
  /* 优化文本渲染，减少重排重绘 */
  will-change: auto;
  contain: layout style;
}

/* 支持代码高亮和markdown样式 */
.message-content :deep(pre) {
  background-color: rgba(0, 0, 0, 0.3);
  border-radius: 6px;
  padding: 12px;
  overflow-x: auto;
  margin: 8px 0;
  border-left: 3px solid var(--accent-color);
}

.message-content :deep(code) {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  background-color: rgba(0, 0, 0, 0.1);
  padding: 2px 4px;
  border-radius: 3px;
}

.message-content :deep(pre code) {
  background-color: transparent;
  padding: 0;
  border-radius: 0;
  color: #e6e6e6;
}

.message-content :deep(p) {
  margin: 8px 0;
  line-height: 1.6;
}

.message-content :deep(ul),
.message-content :deep(ol) {
  margin: 8px 0;
  padding-left: 20px;
}

.message-content :deep(li) {
  margin: 4px 0;
  line-height: 1.5;
}

.message-content :deep(h1),
.message-content :deep(h2),
.message-content :deep(h3),
.message-content :deep(h4),
.message-content :deep(h5),
.message-content :deep(h6) {
  margin: 12px 0 8px 0;
  color: var(--accent-color);
  font-weight: 600;
}

.message-content :deep(blockquote) {
  margin: 8px 0;
  padding: 8px 12px;
  border-left: 3px solid var(--accent-color);
  background-color: rgba(0, 229, 255, 0.05);
  font-style: italic;
}

.message-content :deep(table) {
  border-collapse: collapse;
  margin: 8px 0;
  width: 100%;
}

.message-content :deep(th),
.message-content :deep(td) {
  border: 1px solid var(--border-color);
  padding: 6px 8px;
  text-align: left;
}

.message-content :deep(th) {
  background-color: rgba(0, 229, 255, 0.1);
  font-weight: 600;
}

.message-content :deep(a) {
  color: var(--accent-color);
  text-decoration: none;
}

.message-content :deep(a:hover) {
  text-decoration: underline;
}

.message-time {
  font-size: 11px;
  color: var(--muted-text);
  margin-top: 6px;
  text-align: right;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.streaming-label {
  font-size: 11px;
  color: var(--accent-color);
  font-style: italic;
}

/* 流式响应样式 - 优化动画性能 */
.message.streaming {
  position: relative;
  overflow: visible;
  /* 优化动画性能 */
  will-change: box-shadow;
  transform: translateZ(0);
}

.message.streaming.assistant {
  /* 使用更温和的动画，减少闪屏 */
  animation: streamingGlow 3s ease-in-out infinite alternate;
}

@keyframes streamingGlow {
  0% { 
    box-shadow: var(--glow-effect);
  }
  100% { 
    box-shadow: 0 0 15px rgba(0, 229, 255, 0.3);
  }
}

.streaming-indicator {
  display: flex;
  align-items: center;
  margin-top: 8px;
  padding: 4px 0;
}

.typing-dots {
  display: flex;
  gap: 4px;
  align-items: center;
}

.typing-dots span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: var(--accent-color);
  animation: typing 1.6s infinite ease-in-out;
  opacity: 0.4;
  /* 优化动画性能 */
  will-change: transform, opacity;
  transform: translateZ(0);
}

.typing-dots span:nth-child(1) {
  animation-delay: 0s;
}

.typing-dots span:nth-child(2) {
  animation-delay: 0.3s;
}

.typing-dots span:nth-child(3) {
  animation-delay: 0.6s;
}

@keyframes typing {
  0%, 80%, 100% {
    transform: scale(0.8) translateZ(0);
    opacity: 0.4;
  }
  40% {
    transform: scale(1.1) translateZ(0);
    opacity: 0.8;
  }
}

/* 滚动到底部按钮 */
.scroll-to-bottom-button {
  position: absolute;
  bottom: 120px;
  right: 30px;
  z-index: 1000;
  animation: fadeIn 0.3s ease-in-out;
}

.scroll-to-bottom-button .el-button {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  border: none;
  background: linear-gradient(135deg, var(--primary-color), var(--accent-color));
  transition: all 0.3s ease;
}

.scroll-to-bottom-button .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.4);
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