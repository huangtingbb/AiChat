<template>
  <div class="not-found">
    <div class="tech-circle-large"></div>
    <div class="tech-glitch">404</div>
    <h2 class="tech-heading">页面不存在</h2>
    <div class="tech-line"></div>
    <p>您访问的页面不存在或已被移除</p>
    <router-link to="/" class="back-home">返回首页</router-link>
  </div>
</template>

<style scoped>
.not-found {
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  background-color: var(--dark-bg);
  position: relative;
  overflow: hidden;
  padding: 0 20px;
}

.not-found::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    radial-gradient(circle at 20% 20%, rgba(61, 90, 254, 0.1), transparent 25%),
    radial-gradient(circle at 80% 80%, rgba(0, 229, 255, 0.1), transparent 25%);
}

.tech-circle-large {
  width: 240px;
  height: 240px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(61, 90, 254, 0.05), rgba(0, 229, 255, 0.05));
  border: 1px solid rgba(0, 229, 255, 0.3);
  position: absolute;
  z-index: 0;
  animation: rotateSlow 20s linear infinite;
}

.tech-circle-large::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 200px;
  height: 200px;
  border-radius: 50%;
  border: 1px solid rgba(61, 90, 254, 0.3);
}

.tech-circle-large::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 160px;
  height: 160px;
  border-radius: 50%;
  border: 1px solid rgba(0, 229, 255, 0.3);
}

@keyframes rotateSlow {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.tech-glitch {
  font-size: 120px;
  font-weight: 900;
  color: var(--accent-color);
  text-shadow: 
    0 0 10px rgba(0, 229, 255, 0.7),
    0 0 20px rgba(0, 229, 255, 0.5),
    0 0 30px rgba(0, 229, 255, 0.3);
  position: relative;
  z-index: 1;
  animation: glitch 3s infinite alternate-reverse;
  letter-spacing: 3px;
}

@keyframes glitch {
  0% {
    text-shadow: 
      0 0 10px rgba(0, 229, 255, 0.7),
      0 0 20px rgba(0, 229, 255, 0.5),
      0 0 30px rgba(0, 229, 255, 0.3);
    transform: skewX(0deg);
  }
  5% {
    text-shadow: 
      -2px 0 rgba(61, 90, 254, 0.7),
      2px 0 rgba(0, 229, 255, 0.7);
    transform: skewX(10deg);
  }
  10% {
    text-shadow: 
      2px 0 rgba(61, 90, 254, 0.7),
      -2px 0 rgba(0, 229, 255, 0.7);
    transform: skewX(-10deg);
  }
  15%, 100% {
    text-shadow: 
      0 0 10px rgba(0, 229, 255, 0.7),
      0 0 20px rgba(0, 229, 255, 0.5),
      0 0 30px rgba(0, 229, 255, 0.3);
    transform: skewX(0deg);
  }
}

h2 {
  font-size: 30px;
  color: var(--text-color);
  margin-top: 10px;
  margin-bottom: 10px;
  position: relative;
  z-index: 1;
}

.tech-line {
  height: 1px;
  width: 200px;
  background: linear-gradient(to right, 
    transparent,
    var(--accent-color),
    var(--primary-color),
    transparent
  );
  margin: 10px 0 20px;
  position: relative;
  z-index: 1;
}

p {
  font-size: 16px;
  color: var(--muted-text);
  margin-bottom: 30px;
  position: relative;
  z-index: 1;
}

.back-home {
  display: inline-block;
  padding: 12px 30px;
  background: linear-gradient(135deg, var(--primary-color), var(--accent-color));
  color: white;
  text-decoration: none;
  border-radius: 30px;
  transition: all 0.3s ease;
  position: relative;
  z-index: 1;
  overflow: hidden;
  font-weight: 500;
  letter-spacing: 1px;
  box-shadow: 0 5px 15px rgba(0, 229, 255, 0.3);
}

.back-home:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(0, 229, 255, 0.5);
  background: linear-gradient(135deg, var(--accent-color), var(--primary-color));
}

.back-home:active {
  transform: translateY(1px);
  box-shadow: 0 2px 10px rgba(0, 229, 255, 0.3);
}

.back-home::before {
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

.back-home:hover::before {
  opacity: 1;
}
</style> 