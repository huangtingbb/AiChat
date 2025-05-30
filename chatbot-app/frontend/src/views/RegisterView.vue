<template>
  <div class="register-container">
    <!-- 背景动画元素 -->
    <div class="particles">
      <div class="particle" v-for="i in 20" :key="i"></div>
    </div>
    
    <!-- 右侧装饰 -->
    <div class="right-decoration">
      <div class="tech-circle-large"></div>
      <div class="tech-circle-medium"></div>
      <div class="tech-circle-small"></div>
    </div>
    
    <div class="register-box">
      <div class="register-header">
        <div class="logo-circle">
          <div class="logo-pulse"></div>
          <span class="logo-text">AI</span>
        </div>
        
        <h2 class="tech-heading">注册账号</h2>
        <p class="subtitle">创建您的账户，体验智能对话</p>
        <div class="divider"></div>
      </div>
      
      <div class="register-form">
        <el-form
          ref="registerFormRef"
          :model="registerForm"
          :rules="registerRules"
          label-position="top"
        >
          <div class="form-item">
            <label class="required-label">用户名</label>
            <el-input
              v-model="registerForm.username"
              placeholder="请输入用户名(3-20个字符)"
              class="custom-input"
            >
            </el-input>
            <div class="error-text" v-if="registerForm.usernameError">请输入用户名</div>
          </div>
          
          <div class="form-item">
            <label class="required-label">密码</label>
            <el-input
              v-model="registerForm.password"
              type="password"
              placeholder="请输入密码(至少6个字符)"
              show-password
              class="custom-input"
            >
            </el-input>
            <div class="error-text" v-if="registerForm.passwordError">请输入密码</div>
          </div>
          
          <div class="form-item">
            <label class="required-label">确认密码</label>
            <el-input
              v-model="registerForm.confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              show-password
              class="custom-input"
            >
            </el-input>
            <div class="error-text" v-if="registerForm.confirmPasswordError">
              {{ registerForm.confirmPasswordError === true ? '请确认密码' : registerForm.confirmPasswordError }}
            </div>
          </div>
          
          <div class="form-item">
            <label class="required-label">邮箱</label>
            <el-input
              v-model="registerForm.email"
              placeholder="请输入邮箱"
              class="custom-input"
            >

            </el-input>
            <div class="error-text" v-if="registerForm.emailError">
              {{ registerForm.emailError === true ? '请输入邮箱' : registerForm.emailError }}
            </div>
          </div>
          
          <div class="form-options">
            <el-checkbox v-model="agreeTerms" label="我已阅读并同意服务条款" class="agree-terms" />
          </div>
          
          <el-button
            type="primary"
            :loading="loading"
            @click="handleRegister"
            class="register-button"
          >
            <span class="button-text">注册</span>
            <span class="button-icon">
              <i class="el-icon-right"></i>
            </span>
          </el-button>
        </el-form>
        
        <div class="login-link">
          已有账号？ <a @click="$router.push('/login')">返回登录</a>
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
const registerFormRef = ref(null)
const agreeTerms = ref(false)

// 注册表单
const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  email: '',
  usernameError: false,
  passwordError: false,
  confirmPasswordError: false,
  emailError: false
})

// 验证密码是否一致
const validateConfirmPassword = (value) => {
  if (!value) return true
  if (value !== registerForm.password) {
    return '两次输入的密码不一致'
  }
  return false
}

// 验证邮箱格式
const validateEmail = (value) => {
  if (!value) return true
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return !emailRegex.test(value) ? '请输入正确的邮箱格式' : false
}

// 表单验证规则
const registerRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度应为3-20个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度应为6-20个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: (rule, value, callback) => {
        const result = validateConfirmPassword(value)
        if (result !== false) {
          callback(new Error(result))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

// 处理注册
const handleRegister = async () => {
  // 重置错误状态
  registerForm.usernameError = !registerForm.username
  registerForm.passwordError = !registerForm.password
  registerForm.confirmPasswordError = !registerForm.confirmPassword ? true : validateConfirmPassword(registerForm.confirmPassword)
  registerForm.emailError = !registerForm.email ? true : validateEmail(registerForm.email)
  
  if (
    registerForm.usernameError || 
    registerForm.passwordError || 
    registerForm.confirmPasswordError || 
    registerForm.emailError
  ) {
    return
  }
  
  if (!agreeTerms.value) {
    ElMessage.warning('请阅读并同意服务条款')
    return
  }
  
  loading.value = true
  try {
    // 剔除确认密码字段和错误状态字段
    const { 
      confirmPassword, 
      usernameError, 
      passwordError, 
      confirmPasswordError, 
      emailError, 
      ...registerData 
    } = registerForm
    
    const success = await userStore.registerAction(registerData)
    if (!success) {
      ElMessage.error('注册失败，请稍后重试')
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

.register-container {
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

/* 右侧装饰 */
.right-decoration {
  position: absolute;
  right: 5%;
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

/* 注册框 */
.register-box {
  width: 450px;
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

.register-header {
  padding: 35px 40px 20px;
  text-align: center;
  position: relative;
}

.register-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 5px;
  background: linear-gradient(90deg, #3d5afe, #00e5ff);
}

.register-form {
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

.agree-terms {
  color: rgba(255, 255, 255, 0.7);
}

.agree-terms :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: #00e5ff;
  border-color: #00e5ff;
}

.register-button {
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

.register-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 229, 255, 0.3);
}

.register-button:hover .button-icon {
  opacity: 1;
  margin-left: 10px;
}

.register-button:before {
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

.register-button:hover:before {
  opacity: 1;
}

.login-link {
  margin-top: 25px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  text-align: center;
}

.login-link a {
  color: #00e5ff;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.3s;
  font-weight: 500;
}

.login-link a:hover {
  text-shadow: 0 0 8px rgba(0, 229, 255, 0.5);
  color: #56eeff;
}

/* 响应式设计 */
@media screen and (max-width: 768px) {
  .register-box {
    width: 90%;
    max-width: 430px;
  }
  
  .right-decoration {
    display: none;
  }
}
</style> 