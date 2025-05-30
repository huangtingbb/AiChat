<template>
  <div class="app">
    <router-view />
  </div>
</template>

<style>
:root {
  /* Tech theme color palette */
  --primary-color: #3d5afe;
  --primary-light: #536dfe;
  --primary-dark: #304ffe;
  --accent-color: #00e5ff;
  --accent-light: #18ffff;
  --accent-dark: #00b8d4;
  
  /* Background colors */
  --dark-bg: #1a1a2e;
  --medium-bg: #232342;
  --light-bg: #2a2a4a;
  
  /* Text colors */
  --text-color: #e6e6fa;
  --muted-text: #ababc8;
  --light-text: #ffffff;
  
  /* Border and shadow */
  --border-color: #303060;
  --shadow-color: rgba(0, 229, 255, 0.1);
  --glow-effect: 0 0 10px rgba(0, 229, 255, 0.3);
  
  /* Element Plus theme overrides */
  --el-color-primary: var(--primary-color);
  --el-color-primary-light-3: var(--primary-light);
  --el-color-primary-dark-2: var(--primary-dark);
  --el-color-success: #42d392;
  --el-color-warning: #ffab2b;
  --el-color-danger: #ff4c51;
  --el-color-error: #ff4c51;
  --el-color-info: #5b6278;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Roboto', 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--text-color);
  background-color: var(--dark-bg);
  background-image: 
    radial-gradient(circle at top right, rgba(0, 229, 255, 0.05), transparent 60%),
    radial-gradient(circle at bottom left, rgba(61, 90, 254, 0.05), transparent 60%);
  background-attachment: fixed;
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 5px;
  height: 5px;
}

::-webkit-scrollbar-track {
  background: var(--medium-bg);
}

::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--accent-color);
}

/* Customize Element Plus components */
.el-button--primary {
  background: linear-gradient(135deg, var(--primary-color), var(--accent-color));
  border: none;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
}

.el-button--primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 229, 255, 0.3);
  background: linear-gradient(135deg, var(--accent-color), var(--primary-color));
}

.el-button--primary:active {
  transform: translateY(1px);
}

.el-form-item__label {
  color: var(--text-color) !important;
}

/* Basic animations */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes glow {
  0% { box-shadow: 0 0 5px rgba(0, 229, 255, 0.3); }
  50% { box-shadow: 0 0 15px rgba(0, 229, 255, 0.5); }
  100% { box-shadow: 0 0 5px rgba(0, 229, 255, 0.3); }
}

/* Tech-themed components */
.tech-panel {
  background-color: var(--medium-bg);
  border-radius: 8px;
  border: 1px solid var(--border-color);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
  padding: 20px;
  position: relative;
  overflow: hidden;
  animation: fadeIn 0.5s ease;
}

.tech-panel::before {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  height: 2px;
  background: linear-gradient(to right, 
    transparent, 
    var(--accent-color), 
    var(--primary-color), 
    transparent
  );
}

.tech-heading {
  color: var(--accent-color);
  font-size: 28px;
  text-shadow: 0 0 5px rgba(0, 229, 255, 0.5);
  letter-spacing: 1px;
  margin-bottom: 20px;
}

.app {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
}
</style> 