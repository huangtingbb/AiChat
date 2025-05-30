import { defineStore } from 'pinia'
import { createChat, getUserChats, getChatMessages, sendMessage, updateChatTitle } from '../api/chat'
import { ElMessage } from 'element-plus'

export const useChatStore = defineStore('chat', {
  state: () => ({
    chatList: [],
    currentChatId: null,
    messages: [],
    loading: false,
    sending: false,
    tempChat: null
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
        this.chatList.unshift(response.chat)
        this.currentChatId = response.chat.id
        this.messages = []
        return response.chat
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
        const response = await getUserChats()
        this.chatList = response.chats
        
        if (this.chatList.length > 0 && !this.currentChatId) {
          this.currentChatId = this.chatList[0].id
          await this.fetchMessages(this.currentChatId)
        }
        
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
        this.messages = response.messages
        return this.messages
      } catch (error) {
        ElMessage.error('获取聊天消息失败')
        console.error('获取聊天消息失败:', error)
        return []
      } finally {
        this.loading = false
      }
    },
    
    async sendUserMessage(content) {
      if (!content.trim()) return
      
      try {
        this.sending = true
        
        const tempUserMessage = {
          id: Date.now(),
          role: 'user',
          content,
          created_at: new Date().toISOString()
        }
        this.messages.push(tempUserMessage)
        
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
          
          this.currentChatId = newChat.id
          this.tempChat = null
          
          this.messages = savedMessages
        }
        
        const response = await sendMessage(this.currentChatId, { content })
        
        const userMessageIndex = this.messages.findIndex(msg => msg.id === tempUserMessage.id)
        if (userMessageIndex !== -1) {
          this.messages[userMessageIndex] = response.user_message
        }
        
        this.messages.push(response.bot_message)
        
        return {
          userMessage: response.user_message,
          botMessage: response.bot_message
        }
      } catch (error) {
        ElMessage.error('发送消息失败')
        console.error('发送消息失败:', error)
        
        this.messages = this.messages.filter(msg => msg.id !== Date.now())
        return null
      } finally {
        this.sending = false
      }
    }
  }
}) 