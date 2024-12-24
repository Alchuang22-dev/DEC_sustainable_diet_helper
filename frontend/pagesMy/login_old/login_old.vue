<template>
  <view class="login-container">
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>
    <view class="header">
      <img src="/static/images/index/logo_wide.png" :alt="$t('appTitle')" class="logo" />
    </view>
    <view class="back-link">
      <span class="back-text"></span>
    </view>
    <view class="welcome">
      <span>{{ $t('welcomeLogin') }}</span>
    </view>
    <view class="welcome-message">
      <span>{{ $t('enterPhoneAndPassword') }}</span>
    </view>
    <view class="form">
      <input v-model="phoneNumber" type="text" :placeholder="$t('enterPhoneNumber')" class="input" />
      <input v-model="password" type="password" :placeholder="$t('enterPassword')" class="input" />
      <button class="login-button" @click="check">{{ $t('registerOrLogin') }}</button>
      <input v-if="showRepeatPassword" v-model="repeatPassword" type="password" :placeholder="$t('repeatPassword')" class="input" />
    </view>
    <view class="links">
      <span class="other-login">{{ $t('loginWithCode') }}</span>
      <span class="forgot-password">{{ $t('forgotPassword') }}</span>
    </view>
    <view class="wechat-login">
      <button class="wechat-button" @tap="testLogin">
        <img src="/static/logo.png" alt="WeChat" class="wechat-icon" />
        <span>{{ $t('loginWithWeChat') }}</span>
      </button>
    </view>
    <view class="wechat-login">
      <button class="wechat-button" @click="autoLogin">
        <img src="/static/logo.png" alt="WeChat" class="wechat-icon" />
        <span>测试登录</span>
      </button>
    </view>
  </view>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
//import uni from '@dcloudio/uni-app';
import { useStore } from 'vuex'; // 引入 Vuex 的 useStore

const { t, locale, messages } = useI18n();
const store = useStore(); // 获取 store 实例

console.log('Current locale:', locale.value);
console.log('Available messages:', messages);

const phoneNumber = ref('');
const password = ref('');
const repeatPassword = ref('');
const showRepeatPassword = ref(false);
const avatarUrl = ref(''); // 用户头像
const nickName = ref('');  // 用户昵称


const switchLanguage = (lang) => {
  locale.value = lang;
};

const login = () => {
  uni.switchTab({
    url: '/pages/my_index/my_index',
  });
};

const autoLogin = () => {
  uni.setStorageSync('uid', 'test');
  uni.switchTab({
    url: '/pages/my_index/my_index',
  });
};

// 更新后的 testLogin 函数
const testLogin = () => {
 uni.navigateTo({
 	url: "/pagesMy/wechatLogin/wechatLogin",
 })
};

// 选择头像并输入昵称
const selectAvatarAndInputNickName = () => {
  if (uni.getSystemInfoSync().platform !== 'h5') {
    // 小程序平台使用 chooseAvatar
    uni.chooseAvatar({
      success: (avatarRes) => {
        avatarUrl.value = avatarRes.avatarUrl;
        promptForNickName();
      },
      fail: (err) => {
        uni.showToast({
          title: '头像选择失败',
          icon: 'none',
        });
      }
    });
  } else {
    // H5 使用文件选择
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*';
    input.onchange = (e) => {
      const file = e.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (event) => {
          avatarUrl.value = event.target.result; // 设置头像为选择的图片路径
          promptForNickName(); // 弹出输入昵称的提示框
        };
        reader.readAsDataURL(file); // 将文件转为Base64
      }
    };
    input.click();
  }
};

const promptForNickName = () => {
  uni.showModal({
    title: '输入昵称',
    content: '请输入您的昵称',
    showCancel: true,
    confirmText: '确认',
    cancelText: '取消',
    success: (inputRes) => {
      if (inputRes.confirm) {
        nickName.value = inputRes.content;
        saveUserInfo();
      } else {
        uni.showToast({
          title: '您可以稍后再输入昵称',
          icon: 'none',
        });
      }
    }
  });
};

// 保存用户信息并跳转
const saveUserInfo = () => {
  // 保存头像和昵称到本地存储
  uni.setStorageSync('userInfo', {
    nickName: nickName.value,
    avatarUrl: avatarUrl.value
  });
  uni.setStorageSync('uid', nickName.value);

  // 显示登录成功提示
  uni.showToast({
    title: '登录成功',
    icon: 'success',
    duration: 2000
  });

  // 跳转到首页
  login();
};

const check = () => {
  uni.request({
    url: `https://122.51.231.155:8080/users/${phoneNumber.value}`,
    method: 'GET',
    success: (response) => {
      if (response.statusCode === 200) {
        const realPassword = response.data.realpassword;
        if (password.value === realPassword) {
          uni.showToast({
            title: t('loginSuccess'),
            icon: 'none',
            duration: 2000
          });
          login();
        } else {
          uni.showToast({
            title: t('incorrectCredentials'),
            icon: 'none',
            duration: 2000
          });
        }
      } else if (response.statusCode === 501) {
        registerUser();
      } else {
        uni.showToast({
          title: t('errorTryAgain'),
          icon: 'none',
          duration: 2000
        });
      }
    },
    fail: (error) => {
      console.error('Request error', error);
      uni.showToast({
        title: t('errorTryAgain'),
        icon: 'none',
        duration: 2000
      });
    }
  });
};

const registerUser = () => {
  uni.request({
    url: 'https://122.51.231.155:8080/users/',
    method: 'POST',
    header: {
      'Content-Type': 'application/json'
    },
    data: {
      phoneNumber: phoneNumber.value,
      password: password.value
    },
    success: (response) => {
      if (response.statusCode === 200 || response.statusCode === 201) {
        showRepeatPassword.value = true;
      } else {
        uni.showToast({
          title: t('registerFailed'),
          icon: 'none',
          duration: 2000
        });
      }
    },
    fail: (error) => {
      console.error('Request error', error);
      uni.showToast({
        title: t('errorTryAgain'),
        icon: 'none',
        duration: 2000
      });
    }
  });
};

const confirmRegistration = () => {
  if (!showRepeatPassword.value || repeatPassword.value === '') {
    return;
  }
  uni.request({
    url: `https://122.51.231.155:8080/users/${phoneNumber.value}`,
    method: 'GET',
    success: (response) => {
      if (response.statusCode === 200) {
        const repeatCheck = response.data.repeatcheck;
        if (repeatPassword.value === repeatCheck) {
          uni.showToast({
            title: t('registerSuccess'),
            icon: 'none',
            duration: 2000
          });
          login();
        } else {
          uni.showToast({
            title: t('registerFailed'),
            icon: 'none',
            duration: 2000
          });
        }
      } else {
        uni.showToast({
          title: t('errorTryAgain'),
          icon: 'none',
          duration: 2000
        });
      }
    },
    fail: (error) => {
      console.error('Request error', error);
      uni.showToast({
        title: t('errorTryAgain'),
        icon: 'none',
        duration: 2000
      });
    }
  });
};

watch(repeatPassword, confirmRegistration);
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
