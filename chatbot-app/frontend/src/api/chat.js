import request from './request'

// 创建聊天会话
export function createChat(data) {
  return request({
    url: '/chat',
    method: 'post',
    data
  })
}

// 获取用户的聊天会话列表
export function getUserChatList() {
  return request({
    url: '/chat',
    method: 'get'
  })
}

// 获取聊天会话的消息列表
export function getChatMessages(chatId) {
  return request({
    url: `/chat/${chatId}/message`,
    method: 'get'
  })
}

// 发送消息
export function sendMessage(chatId, data) {
  return request({
    url: `/chat/${chatId}/message`,
    method: 'post',
    data
  })
}

// 更新聊天会话标题
export function updateChatTitle(chatId, data) {
  return request({
    url: `/chat/${chatId}`,
    method: 'patch',
    data
  })
}

// 删除聊天会话
export function deleteChat(chatId) {
  return request({
    url: `/chat/${chatId}`,
    method: 'delete'
  })
}

// 获取可用的AI模型
export function GetAvailableModelList() {
  return request({
    url: '/ai/model',
    method: 'get'
  })
}

// 设置当前AI模型
export function SelectModel(data) {
  return request({
    url: '/ai/model/set',
    method: 'post',
    data
  })
} 