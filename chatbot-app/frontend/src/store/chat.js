import { defineStore } from 'pinia'
import { createChat, getUserChatList, getChatMessages, sendMessage, updateChatTitle, deleteChat } from '../api/chat'
import { ElMessage } from 'element-plus'

export const useChatStore = defineStore('chat', {
  state: () => ({
    chatList: [],
    currentChatId: null,
    messages: [],
    loading: false,
    sending: false,
    tempChat: null,
    currentModel: 'deepseek' // 默认模型
  }),
  
  getters: {
    currentChat: (state) => {
      if (state.tempChat && state.currentChatId === 'temp') {
        return state.tempChat
      }
      return state.chatList.find(chat => chat.id === state.currentChatId) || null
    },
    isFirstMessage: (state) => state.messages.length === 0,
    isTempChat: (state) => state.currentChatId === 'temp'
  },
  
  actions: {
    createTempChat() {
      this.tempChat = {
        id: 'temp',
        title: '新对话',
        created_at: new Date().toISOString()
      }
      this.currentChatId = 'temp'
      this.messages = []
      return this.tempChat
    },
    
    async createNewChat(title = '新对话') {
      try {
        this.loading = true
        const response = await createChat({ title })
        this.chatList.unshift(response.data.chat)
        this.currentChatId = response.data.chat.id
        this.messages = []
        return response.data.chat
      } catch (error) {
        ElMessage.error('创建聊天失败')
        console.error('创建聊天失败:', error)
        return null
      } finally {
        this.loading = false
      }
    },
    
    async updateChatTitle(chatId, title) {
      if (!title.trim()) return false
      
      try {
        this.loading = true
        await updateChatTitle(chatId, { title })
        
        const chatIndex = this.chatList.findIndex(chat => chat.id === chatId)
        if (chatIndex !== -1) {
          this.chatList[chatIndex].title = title
        }
        
        return true
      } catch (error) {
        ElMessage.error('更新聊天标题失败')
        console.error('更新聊天标题失败:', error)
        return false
      } finally {
        this.loading = false
      }
    },
    
    async updateTitleFromMessage(content) {
      if (!this.currentChatId || !this.currentChat) return
      
      if (this.messages.length <= 1) {
        let newTitle = content.trim()
        
        if (newTitle.length > 20) {
          newTitle = newTitle.substring(0, 20) + '...'
        }
        
        await this.updateChatTitle(this.currentChatId, newTitle)
      }
    },
    
    async fetchChatList() {
      try {
        this.loading = true
        const response = await getUserChatList()
        this.chatList = response.data.chats
        
        this.currentChatId = null
        this.messages = []
        
        return this.chatList
      } catch (error) {
        ElMessage.error('获取聊天列表失败')
        console.error('获取聊天列表失败:', error)
        return []
      } finally {
        this.loading = false
      }
    },
    
    async selectChat(chatId) {
      if (this.currentChatId === chatId) return
      
      if (this.currentChatId === 'temp') {
        this.tempChat = null
      }
      
      this.currentChatId = chatId
      this.messages = []
      
      if (chatId !== 'temp') {
        await this.fetchMessages(chatId)
      }
    },
    
    async fetchMessages(chatId) {
      if (chatId === 'temp') return []
      
      try {
        this.loading = true
        const response = await getChatMessages(chatId)
        this.messages = response.data.messages
        return this.messages
      } catch (error) {
        ElMessage.error('获取聊天消息失败')
        console.error('获取聊天消息失败:', error)
        return []
      } finally {
        this.loading = false
      }
    },
    
    async sendUserMessage(content, modelId = 'deepseek') {
      if (!content.trim()) return
      
      try {
        this.sending = true
        this.currentModel = modelId // 保存当前使用的模型
        
        // 创建临时用户消息
        const tempUserMessage = {
          id: Date.now(),
          role: 'user',
          content,
          model_id: modelId,
          created_at: new Date().toISOString()
        }
        
        // 如果没有当前会话，创建一个临时会话
        if (!this.currentChatId) {
          this.createTempChat()
        }
        
        // 添加用户消息到当前会话
        this.messages.push(tempUserMessage)
        
        // 只有当当前是临时会话时，才创建新的聊天
        let chatId = this.currentChatId
        if (this.currentChatId === 'temp') {
          const savedMessages = [...this.messages]
          
          let title = content.trim()
          if (title.length > 20) {
            title = title.substring(0, 20) + '...'
          }
          
          const newChat = await this.createNewChat(title)
          if (!newChat) {
            this.messages.pop()
            return null
          }
          
          chatId = newChat.id
          this.currentChatId = chatId
          this.tempChat = null
          
          this.messages = savedMessages
        }
        
        // 发送消息时附带模型信息
        const response = await sendMessage(chatId, { 
          content,
          model_id: modelId 
        })
        
        const userMessageIndex = this.messages.findIndex(msg => msg.id === tempUserMessage.id)
        if (userMessageIndex !== -1) {
          this.messages[userMessageIndex] = response.data.user_message
        }
        
        // 确保响应的bot消息中包含模型信息
        const botMessage = {
          ...response.data.bot_message,
          model_id: modelId
        }
        
        this.messages.push(botMessage)
        
        return {
          userMessage: response.data.user_message,
          botMessage
        }
      } catch (error) {
        ElMessage.error('发送消息失败')
        console.error('发送消息失败:', error)
        
        this.messages = this.messages.filter(msg => msg.id !== tempUserMessage.id)
        return null
      } finally {
        this.sending = false
      }
    },
    
    async deleteChat(chatId) {
      try {
        this.loading = true
        await deleteChat(chatId)
        
        const chatIndex = this.chatList.findIndex(chat => chat.id === chatId)
        if (chatIndex !== -1) {
          this.chatList.splice(chatIndex, 1)
        }
        
        if (this.currentChatId === chatId) {
          this.currentChatId = null
          this.messages = []
        }
        
        return true
      } catch (error) {
        ElMessage.error('删除聊天失败')
        console.error('删除聊天失败:', error)
        return false
      } finally {
        this.loading = false
      }
    }
  }
}) 