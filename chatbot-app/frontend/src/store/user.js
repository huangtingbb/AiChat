import { defineStore } from 'pinia'
import { login, register, getUserInfo, logout } from '../api/user'
import { ElMessage } from 'element-plus'
import router from '../router'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: JSON.parse(localStorage.getItem('user') || '{}')
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token,
    username: (state) => state.userInfo.username || ''
  },
  
  actions: {
    // 用户登录
    async loginAction(loginData) {
      try {
        const response = await login(loginData)
        this.token = response.token
        this.userInfo = response.user
        
        // 保存到本地存储
        localStorage.setItem('token', this.token)
        localStorage.setItem('user', JSON.stringify(this.userInfo))
        
        ElMessage.success('登录成功')
        router.push('/chat')
        return true
      } catch (error) {
        console.error('登录失败:', error)
        return false
      }
    },
    
    // 用户注册
    async registerAction(registerData) {
      try {
        await register(registerData)
        ElMessage.success('注册成功，请登录')
        router.push('/login')
        return true
      } catch (error) {
        console.error('注册失败:', error)
        return false
      }
    },
    
    // 获取用户信息
    async getUserInfoAction() {
      if (!this.token) return
      
      try {
        const response = await getUserInfo()
        this.userInfo = response.user
        localStorage.setItem('user', JSON.stringify(this.userInfo))
        return this.userInfo
      } catch (error) {
        console.error('获取用户信息失败:', error)
      }
    },
    
    // 退出登录
    async logout() {
      try {
        // 调用后端登出API
        if (this.token) {
          await logout()
        }
      } catch (error) {
        console.error('退出登录API调用失败:', error)
        // 即使API调用失败，也要清除本地状态
      } finally {
        // 清除本地状态
        this.token = ''
        this.userInfo = {}
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        router.push('/login')
        ElMessage.success('已退出登录')
      }
    }
  }
}) 