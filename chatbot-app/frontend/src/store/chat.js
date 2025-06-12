import { defineStore } from 'pinia'
import { createChat, getUserChatList, getChatMessages, sendMessage, updateChatTitle, deleteChat } from '../api/chat'
import { ElMessage } from 'element-plus'

export const useChatStore = defineStore('chat', {
  state: () => ({
    chatList: [],
    currentChatId: null,
    messages: [], // Pinia已经自动处理响应式，不需要手动使用reactive
    loading: false,
    sending: false,
    tempChat: null,
    currentModel: 'deepseek', // 默认模型
    messageUpdateCount: 0 // 添加一个计数器强制触发响应式更新
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
        console.log('Store: 创建聊天响应', response)
        
        // 添加到聊天列表
                  this.chatList.unshift(response.data)
          this.currentChatId = response.data.id
          this.messages = []
          return response.data
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
    
    async sendUserMessage(content, modelId) {
      if (!content.trim()) return
      
      console.log('Store: sendUserMessage开始', {
        content,
        modelId,
        currentChatId: this.currentChatId,
        messagesCount: this.messages.length
      })
      
      try {
        this.sending = true
        this.currentModel = modelId // 保存当前使用的模型
        
        // 创建临时用户消息，使用更精确的ID生成
        const tempUserMessage = {
          id: 'temp_user_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9),
          role: 'user',
          content,
          model_id: modelId,
          created_at: new Date().toISOString()
        }
        
        console.log('Store: 创建临时用户消息', tempUserMessage)
        
        // 如果没有当前会话，创建一个临时会话
        if (!this.currentChatId) {
          this.createTempChat()
        }
        
        // 添加用户消息到当前会话
        this.messages.push(tempUserMessage)
        this.messageUpdateCount++ // 强制触发响应式更新
        
        // 立即添加临时AI消息用于显示流式响应
        const tempBotMessage = {
          id: 'temp_bot_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9),
          role: 'assistant',
          content: '',
          model_id: modelId,
          created_at: new Date().toISOString(),
          isStreaming: true
        }
        
        // 添加AI消息到消息列表
        this.messages.push(tempBotMessage)
        this.messageUpdateCount++ // 强制触发响应式更新
        
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
            // 移除刚添加的消息
            this.messages.splice(-2, 2) // 移除用户消息和AI消息
            return null
          }
          
          chatId = newChat.id
          this.currentChatId = chatId
          this.tempChat = null
          
          this.messages = savedMessages
        }
        
        // 使用SSE发送消息
        return await this.sendMessageWithSSE(chatId, content, modelId, tempUserMessage, tempBotMessage)
        
      } catch (error) {
        ElMessage.error('发送消息失败')
        console.error('发送消息失败:', error)
        
        // 移除临时消息
        this.messages = this.messages.filter(msg => 
          msg.id !== tempUserMessage.id && !msg.isStreaming
        )
        return null
      } finally {
        this.sending = false
      }
    },

    async sendMessageWithSSE(chatId, content, modelId, tempUserMessage, tempBotMessage) {
      // 保存this引用
      const self = this
      
      return new Promise((resolve, reject) => {
        console.log('Store: sendMessageWithSSE开始')
        const token = localStorage.getItem('token')
        const url = `http://localhost:8080/api/chat/${chatId}/message`
        
        // 准备发送SSE请求
        const requestData = {
          content,
          model_id: modelId
        }

        console.log('Store: 发送SSE请求', { url, requestData })

        // 发起请求
        fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
            'Accept': 'text/event-stream',
            'Cache-Control': 'no-cache'
          },
          body: JSON.stringify(requestData)
        })
        .then(response => {
          console.log('Store: fetch响应状态', response.status, response.statusText)
          
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`)
          }

          console.log('Store: SSE请求成功', response)

          const reader = response.body.getReader()
          console.log('Store: 读取器', reader)
          const decoder = new TextDecoder()
          let buffer = ''
          console.log('Store: 准备读取流数据')

          const readStream = async () => {
            try {
              const { done, value } = await reader.read()
              console.log('Store: 读取到数据', { done, valueLength: value })
              
              if (done) {
                console.log('Store: 流读取完成')
                resolve()
                return
              }
              
              buffer += decoder.decode(value, { stream: true })
              console.log('Store: 读取到数据', buffer)
              const lines = buffer.split('\n')
              buffer = lines.pop() // 保留可能不完整的最后一行

              console.log('Store: 处理行数据', lines)

              for (const line of lines) {
                const trimmedLine = line.trim()
                console.log('Store: 原始行数据:', `"${line}"`, '长度:', line.length)
                console.log('Store: 处理后行数据:', `"${trimmedLine}"`, '长度:', trimmedLine.length)
                
                // 处理SSE格式
                // 根据Go后端的SSE格式 c.SSEvent("message", string(data))
                // Gin的SSE格式是: event: message\ndata: jsondata\n\n
                if (trimmedLine.startsWith('event:')) {
                  // 跳过事件行，我们主要关注data行
                  continue
                } else if (trimmedLine.startsWith('data:')) {
                  const data = trimmedLine.substring(5).trim() // 去掉 "data:" 前缀并trim
                  console.log('Store: 接收到SSE数据:', data)
                  
                  if (data.length > 0) {
                    try {
                      const eventData = JSON.parse(data)
                      console.log('Store: 解析SSE事件成功', eventData)
                      self.handleSSEEvent(eventData, tempUserMessage, tempBotMessage, resolve, reject)
                    } catch (error) {
                      console.error('Store: 解析SSE数据失败:', error, data)
                    }
                  }
                } else if (trimmedLine.length > 0) {
                  console.log('Store: 跳过的行数据:', `"${trimmedLine}"`)
                }
              }

              // 继续读取下一块数据
              readStream()
            } catch (error) {
              console.error('Store: 读取流数据失败:', error)
              self.handleStreamError(error, tempUserMessage, tempBotMessage)
              reject(error)
            }
          }

          readStream()
        })
        .catch(error => {
          console.error('Store: 发起SSE请求失败:', error)
          self.handleStreamError(error, tempUserMessage, tempBotMessage)
          reject(error)
        })
      })
    },

    handleSSEEvent(eventData, tempUserMessage, tempBotMessage, resolve, reject) {
      console.log('处理SSE事件:', eventData.type, eventData) // 添加调试日志
      console.log('当前消息数组长度:', this.messages.length)
      
      switch (eventData.type) {
        case 'user_message':
          // 更新用户消息
          const userMessageIndex = this.messages.findIndex(msg => msg.id === tempUserMessage.id)
          console.log('更新用户消息，索引:', userMessageIndex)
          if (userMessageIndex !== -1) {
            // 创建新的消息对象并替换整个消息，确保响应式更新
            const updatedUserMessage = {
              ...eventData.message,
              model_id: tempUserMessage.model_id
            }
            // 使用Vue.set或者直接替换数组元素来触发响应式更新
            this.messages[userMessageIndex] = updatedUserMessage
            this.messageUpdateCount++ // 强制触发响应式更新
            console.log('用户消息更新后，数组长度:', this.messages.length)
          }
          break

        case 'stream_start':
          // 流式响应开始，AI消息已经预先添加，只需要确认状态
          console.log('开始流式响应，AI消息已存在')
          console.log('当前消息数组长度:', this.messages.length)
          console.log('当前所有消息:', this.messages.map(m => ({
            id: m.id, 
            role: m.role, 
            content: m.content?.substring(0, 50) + '...',
            isStreaming: m.isStreaming
          })))
          break

        case 'stream_chunk':
          // 更新AI消息内容
          const botMessageIndex = this.messages.findIndex(msg => msg.id === tempBotMessage.id)
          console.log('接收到文本块:', `"${eventData.text}"`, '消息索引:', botMessageIndex)
          
          if (botMessageIndex !== -1) {
            // 创建新的消息对象来确保响应式更新
            const currentMessage = this.messages[botMessageIndex]
            const updatedMessage = {
              ...currentMessage,
              content: currentMessage.content + eventData.text,
              isStreaming: true
            }
            
            // 直接替换消息对象来触发响应式更新
            this.messages[botMessageIndex] = updatedMessage
            
            // 强制触发响应式更新
            this.messageUpdateCount++
            
            console.log('文本块更新后，消息内容长度:', this.messages[botMessageIndex].content.length)
            console.log('更新后的消息前100字符:', this.messages[botMessageIndex].content.substring(0, 100))
          } else {
            console.error('找不到要更新的AI消息，tempBotMessage.id:', tempBotMessage.id)
            console.log('当前消息ID列表:', this.messages.map(m => ({ id: m.id, role: m.role })))
          }
          break

        case 'stream_end':
          // 流式响应结束
          const finalBotMessageIndex = this.messages.findIndex(msg => msg.id === tempBotMessage.id)
          console.log('流式响应结束，最终消息:', eventData.full_text?.substring(0, 100) + '...')
          console.log('最终消息索引:', finalBotMessageIndex)
          
          if (finalBotMessageIndex !== -1) {
            // 创建最终消息对象
            const finalMessage = {
              id: eventData.message_id || tempBotMessage.id, // 如果有新ID就用新ID，否则保持原ID
              role: 'assistant',
              content: eventData.full_text,
              model_id: tempBotMessage.model_id,
              created_at: new Date().toISOString(),
              isStreaming: false
            }
            // 直接替换消息对象来触发响应式更新
            this.messages[finalBotMessageIndex] = finalMessage
            this.messageUpdateCount++ // 强制触发响应式更新
            console.log('流式响应结束后，消息数组长度:', this.messages.length)
          }
          resolve()
          break

        case 'error':
          // 处理错误
          console.error('收到错误事件:', eventData.error)
          ElMessage.error(eventData.error || '流式响应出错')
          this.handleStreamError(new Error(eventData.error), tempUserMessage, tempBotMessage)
          reject(new Error(eventData.error))
          break

        default:
          console.log('未知SSE事件类型:', eventData.type, eventData)
      }
      
      console.log('事件处理完成后，当前消息数组:', this.messages.map(m => ({
        id: m.id, 
        role: m.role, 
        isStreaming: m.isStreaming,
        contentLength: m.content?.length || 0
      })))
    },

    handleStreamError(error, tempUserMessage, tempBotMessage) {
      // 移除临时消息
      this.messages = this.messages.filter(msg => 
        msg.id !== tempUserMessage.id && msg.id !== tempBotMessage.id
      )
      
      console.error('流式响应错误:', error)
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