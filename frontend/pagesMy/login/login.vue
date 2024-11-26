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
      <button class="wechat-button" @click="testLogin">
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

const testLogin = () => {
  uni.showLoading({
    title: t('loggingIn') || '正在登录...',
    mask: true
  });

  // 获取可用的OAuth服务提供商
  uni.getProvider({
    service: 'oauth',
    success: function(res) {
      // 检查是否支持微信登录
      if (res.provider && res.provider.includes('weixin')) {
        // 调用微信登录接口
        uni.login({
          provider: 'weixin',
          success: function(loginRes) {
            if (loginRes.code) {
              const code = loginRes.code;
              console.log('微信登录码:', code);

              // 获取用户信息
              uni.getUserProfile({
                desc: t('authorizeLogin') || '用于完善会员资料', // 必须填写描述
                success: function(infoRes) {
                  const userInfo = infoRes.userInfo;
                  console.log('用户信息:', userInfo);

                  // 准备发送到后端的数据
                  const data = {
                    code: code,
                    userInfo: {
                      nickName: userInfo.nickName,
                      avatarUrl: userInfo.avatarUrl,
                      gender: userInfo.gender,
                      province: userInfo.province,
                      city: userInfo.city,
                      country: userInfo.country
                    }
                  };

                  // 调用 Vuex 的 Login action
                  store.dispatch('Login', {
                    type: 'weixin',
                    url: 'https://122.51.231.155:8080/wechat-login', // 后端登录接口
                    data
                  }).then(res => {
                    uni.hideLoading();
                    if (res === 'ok') {
                      uni.showToast({
                        title: t('loginSuccess') || '登录成功',
                        icon: 'success',
                        duration: 2000
                      });
                      login();
                    } else {
                      uni.showToast({
                        icon: 'none',
                        title: res || (t('loginFailed') || '登录失败')
                      });
                      console.error('后端登录失败:', res);
                    }
                  }).catch(err => {
                    uni.hideLoading();
                    uni.showToast({
                      icon: 'none',
                      title: err || (t('loginFailed') || '登录失败')
                    });
                    console.error('登录过程中发生错误:', err);
                  });
                },
                fail: function(err) {
                  uni.hideLoading();
                  uni.showToast({
                    icon: 'none',
                    title: t('getUserProfileFailed') || '获取用户信息失败'
                  });
                  console.error('获取用户信息失败:', err);
                }
              });
            } else {
              uni.hideLoading();
              uni.showToast({
                title: t('loginFailed') || '登录失败',
                icon: 'none',
                duration: 2000
              });
              console.error('微信登录失败:', loginRes.errMsg);
            }
          },
          fail: function(err) {
            uni.hideLoading();
            uni.showToast({
              icon: 'none',
              title: t('loginFailed') || '登录失败',
              duration: 2000
            });
            console.error('微信登录接口调用失败:', err);
          }
        });
      } else {
        uni.hideLoading();
        uni.showToast({
          icon: 'none',
          title: t('weixinNotSupported') || '当前设备不支持微信登录',
          duration: 2000
        });
        console.warn('当前设备不支持微信登录:', res.provider);
      }
    },
    fail: function(err) {
      uni.hideLoading();
      uni.showToast({
        icon: 'none',
        title: t('getProviderFailed') || '获取服务提供商失败',
        duration: 2000
      });
      console.error('获取OAuth服务提供商失败:', err);
    }
  });
};

const performWeixinLogin = () => {
  uni.login({
    provider: 'weixin',
    success: function(loginRes) {
      if (loginRes.authResult && loginRes.authResult.code) {
        const authResult = loginRes.authResult;
        
        // 获取用户信息
        uni.getUserInfo({
          provider: 'weixin',
          success: function(infoRes) {
            const userInfo = infoRes.userInfo;
            const data = {
              ...authResult,
              userInfo
            };

            // 调用 Vuex 的 Login action
            store.dispatch('Login', {
              type: 'weixin',
              url: 'https://122.51.231.155:8080/wechat-login', // 后端登录接口
              data
            }).then(res => {
              if (res === 'ok') {
                uni.hideLoading();
                uni.showToast({
                  title: t('loginSuccess'),
                  icon: 'success',
                  duration: 2000
                });
                login();
              } else {
                uni.hideLoading();
                uni.showToast({
                  icon: 'none',
                  title: res || t('loginFailed')
                });
                console.error('后端登录失败:', res);
              }
            }).catch(err => {
              uni.hideLoading();
              uni.showToast({
                icon: 'none',
                title: err || t('loginFailed')
              });
              console.error('登录过程中发生错误:', err);
            });
          },
          fail: function(err) {
            uni.hideLoading();
            uni.showToast({
              icon: 'none',
              title: t('getUserInfoFailed') || '获取用户信息失败',
              duration: 2000
            });
            console.error('获取用户信息失败:', err);
          }
        });
      } else {
        uni.hideLoading();
        uni.showToast({
          title: t('loginFailed'),
          icon: 'none',
          duration: 2000
        });
        console.error('微信登录失败:', loginRes.errMsg);
      }
    },
    fail: function(err) {
      uni.hideLoading();
      uni.showToast({
        icon: 'none',
        title: t('loginFailed'),
        duration: 2000
      });
      console.error('微信登录接口调用失败:', err);
    }
  });
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
