<template>
  <view class="login-container">
	<image src="/static/images/index/background_img.jpg" class="background-image"></image>
    <view class="header">
      <img src="/static/images/index/logo_wide.png" alt="DEC 可持续饮食助手" class="logo" />
    </view>
    <view class="back-link">
      <span class="back-text"></span>
    </view>
    <view class="welcome">
      <span>欢迎登录！</span>
    </view>
    <view class="welcome-message">
      <span>请输入手机号及密码</span>
    </view>
    <view class="form">
      <input v-model="phoneNumber" type="text" placeholder="请输入手机号" class="input" />
      <input v-model="password" type="password" placeholder="请输入密码" class="input" />
      <button class="login-button" @click="check">注册/登录</button>
      <input v-if="showRepeatPassword" v-model="repeatPassword" type="password" placeholder="请重复输入密码" class="input" />
    </view>
    <view class="links">
      <span class="other-login">用验证码登录</span>
      <span class="forgot-password">忘记密码?</span>
    </view>
    <view class="wechat-login">
      <button class="wechat-button" @click="login">
        <img src="/static/logo.png" alt="WeChat" class="wechat-icon" />
        <span>用微信用户授权登录</span>
      </button>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      phoneNumber: '',
      password: '',
      repeatPassword: '',
      showRepeatPassword: false
    };
  },
  methods: {
	login() {
		uni.switchTab({
		  url: `/pages/my_index/my_index?uid=${this.phoneNumber}`,
		});
	},
    check() {
      uni.request({
        url: `https://122.51.231.155:8080/users/${this.phoneNumber}`,
        method: 'GET',
        success: (response) => {
          if (response.statusCode === 200) {
            const realPassword = response.data.realpassword;
            if (this.password === realPassword) {
              uni.showToast({
                title: '登录成功',
                icon: 'none',
                duration: 2000
              });
			  login();
            } else {
              uni.showToast({
                title: '账号或密码错误',
                icon: 'none',
                duration: 2000
              });
            }
          } else if (response.statusCode === 501) {
            this.registerUser();
          } else {
            uni.showToast({
              title: '发生错误，请稍后再试',
              icon: 'none',
              duration: 2000
            });
          }
        },
        fail: (error) => {
          console.error('请求错误', error);
          uni.showToast({
            title: '发生错误，请稍后再试',
            icon: 'none',
            duration: 2000
          });
        }
      });
    },
    registerUser() {
      uni.request({
        url: 'https://122.51.231.155:8080/users/',
        method: 'POST',
        header: {
          'Content-Type': 'application/json'
        },
        data: {
          phoneNumber: this.phoneNumber,
          password: this.password
        },
        success: (response) => {
          if (response.statusCode === 200 || response.statusCode === 201) {
            this.showRepeatPassword = true;
          } else {
            uni.showToast({
              title: '注册请求失败，请稍后再试',
              icon: 'none',
              duration: 2000
            });
          }
        },
        fail: (error) => {
          console.error('请求错误', error);
          uni.showToast({
            title: '发生错误，请稍后再试',
            icon: 'none',
            duration: 2000
          });
        }
      });
    },
    confirmRegistration() {
      if (!this.showRepeatPassword || this.repeatPassword === '') {
        return;
      }
      uni.request({
        url: `https://122.51.231.155:8080/users/${this.phoneNumber}`,
        method: 'GET',
        success: (response) => {
          if (response.statusCode === 200) {
            const repeatCheck = response.data.repeatcheck;
            if (this.repeatPassword === repeatCheck) {
              uni.showToast({
                title: '注册成功',
                icon: 'none',
                duration: 2000
              });
			  login();
            } else {
              uni.showToast({
                title: '注册账号失败',
                icon: 'none',
                duration: 2000
              });
            }
          } else {
            uni.showToast({
              title: '发生错误，请稍后再试',
              icon: 'none',
              duration: 2000
            });
          }
        },
        fail: (error) => {
          console.error('请求错误', error);
          uni.showToast({
            title: '发生错误，请稍后再试',
            icon: 'none',
            duration: 2000
          });
        }
      });
    }
  },
  watch: {
    repeatPassword: 'confirmRegistration'
  }
};
</script>

<style scoped>
body {
  font-family: 'Arial', sans-serif;
  background: url('/static/images/index/background_img.jpg') no-repeat center center fixed;
  background-size: cover;
  background-color: #f0f4f7;
  margin: 0;
  padding: 0;
}

.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: 0;
  opacity: 0.1;
}
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  background-color: #f8fbf7;
  height: 100vh;
  z-index: 1;
  position: relative;
}
.header {
  margin-bottom: 20px;
  display: flex;
  justify-content: center;
  width: 100%;
}
.logo {
  width: 300px;
  height: 50px;
}
.back-link {
  align-self: flex-start;
  margin-bottom: 10px;
}
.back-text {
  color: #555;
  font-size: 16px;
}
.welcome {
  text-align: left;
  margin-bottom: 20px;
}
.welcome span {
  display: block;
  color: #555;
  font-size: 14px;
  font-weight: bold;
  margin-bottom: 5px;
}
.welcome-message {
  text-align: left;
  margin-bottom: 20px;
}
.welcome-message span {
  display: block;
  color: #333;
  font-size: 20px;
  font-weight: lighter;
  margin-bottom: 5px;
}
.form {
  width: 100%;
  max-width: 300px;
  display: flex;
  flex-direction: column;
  margin-bottom: 20px;
}
.input {
  margin-bottom: 15px;
  padding: 12px;
  border: 1px solid #ccc;
  border-radius: 8px;
  font-size: 16px;
}
.login-button {
  background-color: #48c079;
  color: white;
  border: none;
  padding: 12px;
  border-radius: 8px;
  font-size: 18px;
  cursor: pointer;
  width: 100%;
}
.links {
  display: flex;
  justify-content: space-between;
  width: 100%;
  max-width: 300px;
  margin-bottom: 20px;
}
.other-login, .forgot-password {
  color: #48c079;
  font-size: 16px;
  cursor: pointer;
}
.wechat-login {
  margin-top: 20px;
}
.wechat-button {
  display: flex;
  align-items: center;
  background-color: #f2f2f2;
  padding: 12px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}
.wechat-icon {
  width: 20px;
  height: 20px;
  margin-right: 10px;
}
</style>
