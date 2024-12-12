<template>
  <view class="container">
    <view v-if="notshowVerification" class="icon_holder">
      <img src="/static/images/settings/mail.png" class="delete-icon"></img>
    </view>
	<view v-if="showVerification" class="icon_holder">
	  <img src="/static/images/settings/mail_open.png" class="delete-icon"></img>
	</view>
    <view class="text_holder" :style="{ fontSize: fontSizeBig + 'px' }">
      <text style="font-weight: bold;">{{ $t('settings_account_binding') }}</text>
    </view>
    <view class="text_holder" :style="{ fontSize: fontSize + 'px' }">
      <text>{{ $t('change_binding_phone') }}</text>
    </view>
    <view class="text_holder" :style="{ fontSize: fontSize + 'px' }">
      <text>{{ $t('unbind_account_notice') }}</text>
    </view>
    <view class="form">
      <input v-model="phoneNumber" type="text" :placeholder="$t('enterPhoneNumber')" class="input" />
    </view> 
    <view class="button_holder">
      <button class="confirm-button" :disabled="isButtonDisabled" @click="sendVerificationCode">
        {{ buttonText }}
      </button> 
    </view>
    <view v-if="showVerification" class="form">
      <input v-model="verificationCode" type="text" :placeholder="$t('input_verification_code')" class="input" />
      <button class="confirm-button" @click="verifyCode">{{ $t('confirm_binding') }}</button>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue';

// 定义字体大小
const fontSize = ref(16);
const fontSizeBig = fontSize.value + 4;

const phoneNumber = ref('');
const verificationCode = ref('');
const showVerification = ref(false);
const notshowVerification = ref(true);
const isButtonDisabled = ref(false);
const buttonText = ref('发送验证码');
const countdown = ref(60);
let timer = null;

function sendVerificationCode() {
  // 模拟发送验证码的过程
  uni.showToast({
    title: '已发送验证码',
    icon: 'success'
  });
  showVerification.value = true;
  notshowVerification.value = false;
  startCountdown();
}

function startCountdown() {
  isButtonDisabled.value = true;
  buttonText.value = `请等待 ${countdown.value}s`;
  timer = setInterval(() => {
    countdown.value--;
    buttonText.value = `请等待 ${countdown.value}s`;
    if (countdown.value === 0) {
      clearInterval(timer);
      isButtonDisabled.value = false;
      buttonText.value = '发送验证码';
      countdown.value = 60;
    }
  }, 1000);
}

function verifyCode() {
  // 模拟验证码校验
  if (verificationCode.value === '1234') { // 假设正确的验证码是1234
    const userId = uni.getStorageSync('UserId');
    uni.request({
      url: `http://122.51.231.155:8080/users/${userId}`,
      method: 'PUT',
      header: {
        'Content-Type': 'application/json'
      },
      data: {
        phonenumber: phoneNumber.value
      },
      success: (res) => {
        if (res.statusCode === 200) {
          uni.showToast({
            title: '绑定成功',
            icon: 'success'
          });
          showVerification.value = false;
		  notshowVerification.value = true;
        } else {
          uni.showToast({
            title: '绑定失败',
            icon: 'error'
          });
        }
      },
      fail: () => {
        uni.showToast({
          title: '绑定失败',
          icon: 'error'
        });
      }
    });
  } else {
    uni.showToast({
      title: '不支持此操作',
      icon: 'error'
    });
  }
}
</script>

<style>
/* 设置容器为 Flex 布局 */
.container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  justify-content: space-between;
}

.icon_holder {
  text-align: center;
  padding-top: 100px;
}

.delete-icon {
  width: 200px;
  height: 200px;
}

.text_holder {
  text-align: center;
  padding-top: 10px;
}

.button_holder {
  text-align: center;
  margin-top: auto;
  padding-bottom: 30px;
}

.confirm-button {
  width: 40%;
  padding: 10px;
  background-color: #48c079;
  color: #fff;
  text-align: center;
  border-radius: 5px;
  margin-top: 10px;
}

.form {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  max-width: 300px;
  margin: 0 auto;
  text-align: center;
  margin-top: 30px;
}

.input {
  margin-bottom: 15px;
  padding: 12px;
  text-align: center;
  border: 1px solid #ccc;
  border-radius: 8px;
  font-size: 16px;
}
</style>
