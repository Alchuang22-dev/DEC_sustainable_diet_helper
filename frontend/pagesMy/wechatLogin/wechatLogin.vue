<template>
  <view class="containar">
    <view class="avatarUrl">
      <button type="balanced" open-type="chooseAvatar" @chooseavatar="onChooseavatar">
        <image :src="avatarUrl || '/static/images/index/background_img.jpg'" class="refreshIcon"></image> <!-- 使用占位符 -->
      </button>
    </view>
    <view class="nickname">
      <text>{{$t('Nickname')}}</text>
      <input type="nickname" class="weui-input" :value="nickName" @blur="bindblur"
        :placeholder= "$t('enterNickname')" @input="bindinput"/>
    </view>

    <view class="btn">
      <view class="btn-sub" @click="onSubmit">{{$t('save')}}</view>
    </view>
  </view>
</template>


<script setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
//import uni from '@dcloudio/uni-app';
import { useStore } from 'vuex'; // 引入 Vuex 的 useStore

const { t, locale, messages } = useI18n();

// 使用 ref 创建响应式数据
const avatarUrl = ref('');  // 初始值为空，头像上传成功后更新
const nickName = ref('');

// 绑定昵称输入框失去焦点时的事件
function bindblur(e) {
  nickName.value = e.detail.value; // 获取微信昵称
}

// 绑定昵称输入框输入时的事件
function bindinput(e) {
  nickName.value = e.detail.value; // 如果只用 blur 方法，用户在输入昵称后点击保存按钮，昵称才会更新
}

// 选择头像
function onChooseavatar(e) {
  let { avatarUrl: newAvatarUrl } = e.detail;
  uni.showLoading({
    title: '加载中',
  });

  // 判断 avatarUrl 是否为空
  if (!newAvatarUrl) {
    uni.showToast({
      title: '请选择头像',
      icon: 'none',
      duration: 2000,
    });
    return;
  }

  // 上传头像
  uni.uploadFile({
    url: '后台uploadFile接口',  // 替换为实际接口
    filePath: newAvatarUrl,
    name: 'file',
    header: {
      token: '自己的token',
    },
    success: (uploadFileRes) => {
      let data = JSON.parse(uploadFileRes.data);
      if (data.code === 0 && data.data) {
        avatarUrl.value = data.data;  // 更新头像 URL
      } else {
        uni.showToast({
          title: '头像上传失败',
          icon: 'none',
          duration: 2000,
        });
      }
    },
    fail: (error) => {
      uni.showToast({
        title: error,
        duration: 2000,
      });
    },
    complete: () => {
      uni.hideLoading();
    },
  });
}

// 提交表单
function onSubmit() {
  if (nickName.value === '') {
    uni.showToast({
      icon: 'none',
      title: t('needNickname'),
    });
    return false;
  }

  // 调用 uni.login 获取微信登录凭证
  uni.login({
    success: (loginRes) => {
      if (loginRes.code) {
        // 发送登录凭证到服务器，获取用户信息和 token
        uni.showLoading({
          title: '登录中',
        });
        uni.request({
          url: 'http://122.51.231.155:8080/users/wechatlogin', // 替换为实际的登录接口
          method: 'POST',
          data: {
            code: loginRes.code,  // 发送微信登录凭证
          },
          success: (res) => {
            if (res.data.code === 0) {
              // 假设服务器返回用户信息和 token
              const { userInfo, token } = res.data.data;
              uni.setStorageSync('token', token); // 保存 token

              // 显示登录成功提示
              uni.showToast({
                title: '登录成功',
                icon: 'success',
                duration: 2000,
              });

              // 跳转到首页
              login();
            } else {
              uni.showToast({
                title: '登录失败，请重试',
                icon: 'none',
                duration: 2000,
              });
            }
          },
          fail: (error) => {
            uni.showToast({
              title: '请求失败，请重试',
              icon: 'none',
              duration: 2000,
            });
          },
          complete: () => {
            uni.hideLoading();
          }
        });
      } else {
        uni.showToast({
          title: '微信登录失败，请重试',
          icon: 'none',
          duration: 2000,
        });
      }
    },
    fail: (error) => {
      uni.showToast({
        title: '登录失败，请重试',
        icon: 'none',
        duration: 2000,
      });
    }
  });
  // 保存用户信息和 token 到本地存储
  uni.setStorageSync('userInfo', {
    nickName: nickName.value,
    avatarUrl: avatarUrl.value || '/static/images/index/background_img.jpg',  // 使用默认头像如果头像为空
  });
  uni.setStorageSync('uid', nickName.value);
  login();
}

// 跳转到首页
function login() {
  uni.switchTab({
    url: '/pages/my_index/my_index',
  });
}

</script>

<style scoped>
.containar {
  padding: 0;
}

.avatarUrl {
  padding: 80rpx 0 40rpx;
  background: #fff;
}

.avatarUrl button {
  background: #fff;
  line-height: 80rpx;
  height: auto;
  width: auto;
  padding: 20rpx 30rpx;
  margin: 0;
  display: flex;
  justify-content: center;
  align-items: center;
}

.refreshIcon {
  width: 160rpx;
  height: 160rpx;
  border-radius: 50%;
  object-fit: cover; /* 保证图片不会变形 */
}

.nickname {
  background: #fff;
  padding: 20rpx 30rpx 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.weui-input {
  padding-left: 60rpx;
}

.btn {
  width: 100%;
}

.btn-sub {
  width: 670rpx;
  margin: 80rpx auto 0;
  height: 90rpx;
  background: #48c079;
  border-radius: 45rpx;
  line-height: 90rpx;
  text-align: center;
  font-size: 36rpx;
  color: #fff;
}
</style>

