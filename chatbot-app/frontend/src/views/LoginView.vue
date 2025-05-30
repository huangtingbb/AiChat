<template>
  <div class="login-container">
    <!-- 背景动画元素 -->
    <div class="particles">
      <div class="particle" v-for="i in 20" :key="i"></div>
    </div>
    
    <!-- 左侧装饰 -->
    <div class="left-decoration">
      <div class="tech-circle-large"></div>
      <div class="tech-circle-medium"></div>
      <div class="tech-circle-small"></div>
    </div>
    
    <div class="login-box">
      <div class="login-header">
        <div class="logo-circle">
          <div class="logo-pulse"></div>
          <span class="logo-text">AI</span>
        </div>
        
        <h2 class="tech-heading">AI聊天助手</h2>
        <p class="subtitle">登录您的账户，开始智能对话</p>
        <div class="divider"></div>
      </div>
      
      <div class="login-form">
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          label-position="top"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
              class="custom-input"
              :class="{ 'is-error': loginForm.usernameError }"
              @input="loginForm.usernameError = false"
            >
            </el-input>
            <div v-if="loginForm.usernameError" class="error-text">
              请输入用户名
            </div>
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              class="custom-input"
              :class="{ 'is-error': loginForm.passwordError }"
              @input="loginForm.passwordError = false"
              show-password
            >
            </el-input>
            <div v-if="loginForm.passwordError" class="error-text">
              请输入密码
            </div>
          </el-form-item>

          <div class="form-options">
            <el-checkbox v-model="rememberMe" label="记住我" class="remember-me" />
            <a class="forgot-password">忘记密码?</a>
          </div>
          
          <el-button
            type="primary"
            :loading="loading"
            @click="handleLogin"
            class="login-button"
          >
            <span class="button-text">登录</span>
            <span class="button-icon">
              <i class="el-icon-right"></i>
            </span>
          </el-button>
        </el-form>
        
        <div class="register-link">
          还没有账号？ <a @click="$router.push('/register')">立即注册</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '../store/user'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()
const loading = ref(false)
const rememberMe = ref(false)

// 登录表单
const loginForm = reactive({
  username: '',
  password: '',
  usernameError: false,
  passwordError: false
})

// 表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度应为3-20个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度应为6-20个字符', trigger: 'blur' }
  ]
}

// 登录表单引用
const loginFormRef = ref(null)

// 处理登录
const handleLogin = async () => {
  // 重置错误状态
  loginForm.usernameError = !loginForm.username
  loginForm.passwordError = !loginForm.password
  
  if (loginForm.usernameError || loginForm.passwordError) {
    return
  }
  
  loading.value = true
  try {
    const success = await userStore.loginAction(loginForm)
    if (!success) {
      ElMessage.error('登录失败，请检查用户名和密码')
    }
  } finally {
    loading.value = false
  }
}

// 生成随机的背景粒子
onMounted(() => {
  const particles = document.querySelectorAll('.particle')
  particles.forEach(particle => {
    const size = Math.random() * 15 + 5
    const posX = Math.random() * 100
    const posY = Math.random() * 100
    const duration = Math.random() * 20 + 10
    const delay = Math.random() * 5
    
    particle.style.width = `${size}px`
    particle.style.height = `${size}px`
    particle.style.left = `${posX}%`
    particle.style.top = `${posY}%`
    particle.style.animationDuration = `${duration}s`
    particle.style.animationDelay = `${delay}s`
  })
})
</script>

<style scoped>
/* 动画效果 */
@keyframes float {
  0% {
    transform: translateY(0) rotate(0deg);
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
  100% {
    transform: translateY(-100px) rotate(360deg);
    opacity: 0;
  }
}

@keyframes pulse {
  0% {
    transform: scale(1);
    opacity: 0.6;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.2;
  }
  100% {
    transform: scale(1);
    opacity: 0.6;
  }
}

@keyframes rotate {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

@keyframes shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}

.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #121139 0%, #1e1e38 50%, #252547 100%);
  position: relative;
  overflow: hidden;
}

/* 背景粒子 */
.particles {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  z-index: 0;
}

.particle {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(135deg, #3d5afe, #00e5ff);
  opacity: 0.2;
  animation: float linear infinite;
  box-shadow: 0 0 10px rgba(0, 229, 255, 0.5);
}

/* 左侧装饰 */
.left-decoration {
  position: absolute;
  left: 5%;
  top: 50%;
  transform: translateY(-50%);
  z-index: 1;
}

.tech-circle-large {
  width: 250px;
  height: 250px;
  border-radius: 50%;
  border: 2px solid rgba(61, 90, 254, 0.2);
  position: relative;
  animation: rotate 30s linear infinite;
}

.tech-circle-medium {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  border: 2px solid rgba(0, 229, 255, 0.2);
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation: rotate 20s linear infinite reverse;
}

.tech-circle-small {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(0, 229, 255, 0.1), transparent);
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation: pulse 5s ease-in-out infinite;
}

/* 登录框 */
.login-box {
  width: 420px;
  background: rgba(28, 29, 58, 0.8);
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(61, 90, 254, 0.2);
  z-index: 2;
  position: relative;
  display: flex;
  flex-direction: column;
}

.login-header {
  padding: 35px 40px 20px;
  text-align: center;
  position: relative;
}

.login-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 5px;
  background: linear-gradient(90deg, #3d5afe, #00e5ff);
}

.login-form {
  padding: 0 40px 35px;
}

.logo-circle {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  border-radius: 50%;
  background: rgba(20, 20, 40, 0.5);
  border: 1px solid rgba(0, 229, 255, 0.3);
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
}

.logo-pulse {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(0, 229, 255, 0.5), transparent 70%);
  animation: pulse 3s infinite;
}

.logo-text {
  font-size: 28px;
  font-weight: bold;
  color: #00e5ff;
  text-shadow: 0 0 10px rgba(0, 229, 255, 0.7);
  z-index: 1;
}

.tech-heading {
  font-size: 28px;
  background: linear-gradient(90deg, #3d5afe, #00e5ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0 0 10px;
  font-weight: 600;
}

.subtitle {
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  margin: 0 0 15px;
}

.divider {
  height: 1px;
  background: linear-gradient(to right, 
    transparent,
    rgba(0, 229, 255, 0.5),
    rgba(61, 90, 254, 0.5),
    transparent
  );
  margin: 0 auto 30px;
  width: 80%;
  position: relative;
}

.divider::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(to right, 
    transparent, 
    rgba(255, 255, 255, 0.5), 
    transparent
  );
  background-size: 200% 100%;
  animation: shimmer 3s infinite;
}

.form-item {
  margin-bottom: 20px;
  text-align: left;
}

.required-label {
  display: block;
  margin-bottom: 8px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 14px;
  text-align: left;
  position: relative;
  font-weight: 500;
}

.required-label:before {
  content: "*";
  color: #ff4d4f;
  margin-right: 4px;
}

.error-text {
  color: #ff4d4f;
  font-size: 12px;
  margin-top: 5px;
  text-align: left;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
}

.remember-me {
  color: rgba(255, 255, 255, 0.7);
}

.remember-me :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: #00e5ff;
  border-color: #00e5ff;
}

.forgot-password {
  color: rgba(0, 229, 255, 0.8);
  font-size: 14px;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.3s;
}

.forgot-password:hover {
  color: #00e5ff;
  text-shadow: 0 0 8px rgba(0, 229, 255, 0.5);
}

.login-button {
  width: 100%;
  height: 50px;
  background: linear-gradient(90deg, #3d5afe, #00e5ff);
  border: none;
  border-radius: 10px;
  color: white;
  font-size: 16px;
  font-weight: 500;
  margin-top: 10px;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  display: flex;
  justify-content: center;
  align-items: center;
}

.button-text {
  position: relative;
  z-index: 1;
  letter-spacing: 1px;
}

.button-icon {
  opacity: 0;
  margin-left: -20px;
  transition: all 0.3s ease;
  position: relative;
  z-index: 1;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 229, 255, 0.3);
}

.login-button:hover .button-icon {
  opacity: 1;
  margin-left: 10px;
}

.login-button:before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.3) 0%, transparent 70%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.login-button:hover:before {
  opacity: 1;
}

.register-link {
  margin-top: 25px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  text-align: center;
}

.register-link a {
  color: #00e5ff;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.3s;
  font-weight: 500;
}

.register-link a:hover {
  text-shadow: 0 0 8px rgba(0, 229, 255, 0.5);
  color: #56eeff;
}

/* 响应式设计 */
@media screen and (max-width: 768px) {
  .login-box {
    width: 90%;
    max-width: 400px;
  }
  
  .left-decoration {
    display: none;
  }
}
</style> 