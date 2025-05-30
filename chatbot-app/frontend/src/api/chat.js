import request from '../utils/request'

// 创建聊天会话
export function createChat(data) {
  return request({
    url: '/chats',
    method: 'post',
    data
  })
}

// 获取用户的聊天会话列表
export function getUserChats() {
  return request({
    url: '/chats',
    method: 'get'
  })
}

// 获取聊天会话的消息列表
export function getChatMessages(chatId) {
  return request({
    url: `/chats/${chatId}/messages`,
    method: 'get'
  })
}

// 发送消息
export function sendMessage(chatId, data) {
  return request({
    url: `/chats/${chatId}/messages`,
    method: 'post',
    data
  })
}

// 更新聊天会话标题
export function updateChatTitle(chatId, data) {
  return request({
    url: `/chats/${chatId}`,
    method: 'patch',
    data
  })
} 