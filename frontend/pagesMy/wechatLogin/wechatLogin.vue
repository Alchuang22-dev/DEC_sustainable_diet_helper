<template>
  <view class="container">
    <view class="avatarUrl">
      <button type="balanced" open-type="chooseAvatar" @chooseavatar="onChooseAvatar">
        <image :src="avatarUrl || '/static/images/index/background_img.jpg'" class="refreshIcon"></image> <!-- 使用占位符 -->
      </button>
    </view>
    <view class="nickname">
      <text>{{ $t('Nickname') }}</text>
      <input type="nickname" class="weui-input" :value="nickName" @blur="bindblur"
        :placeholder= "$t('enterNickname')" @input="bindinput"/>
    </view>

    <view class="btn">
      <view class="btn-sub" @click="onSubmit">{{ $t('save') }}</view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '@/stores/user'; // 引入用户存储
import { storeToRefs } from 'pinia';

// 国际化
const { t } = useI18n();

// Pinia用户存储
const userStore = useUserStore();
const { uid, isLoggedIn } = storeToRefs(userStore);

// 响应式数据
const avatarUrl = ref(''); // 初始值为空，头像选择后更新
const nickName = ref('');

// 页面加载时从本地存储读取头像和昵称
onMounted(() => {
  const storedAvatarUrl = uni.getStorageSync('avatarUrl');
  const storedNickName = uni.getStorageSync('userInfo')?.nickName;
  if (storedAvatarUrl) {
    avatarUrl.value = storedAvatarUrl;
  }
  if (storedNickName) {
    nickName.value = storedNickName;
  }
});

// 绑定昵称输入框失去焦点时的事件
function bindBlur(e) {
  nickName.value = e.detail.value; // 获取微信昵称
}

// 绑定昵称输入框输入时的事件
function bindInput(e) {
  nickName.value = e.detail.value; // 实时更新昵称
}

// 选择头像
function onChooseAvatar(e) {
  let { avatarUrl: newAvatarUrl } = e.detail;

  if (!newAvatarUrl) {
    uni.showToast({
      title: '请选择头像',
      icon: 'none',
      duration: 2000,
    });
    return;
  }

  // 更新头像 URL
  avatarUrl.value = newAvatarUrl;

  // 保存头像路径到本地存储
  uni.setStorageSync('avatarUrl', newAvatarUrl);

  uni.showToast({
    title: '头像已更新',
    icon: 'success',
    duration: 2000,
  });
}

// 提交表单
async function onSubmit() {
  if (nickName.value === '') {
    uni.showToast({
      icon: 'none',
      title: t('needNickname'),
    });
    return false;
  }

  try {
    // 调用 uni.login 获取微信登录凭证
    const loginRes = await new Promise((resolve, reject) => {
      uni.login({
        success: resolve,
        fail: reject,
      });
    });

    if (!loginRes.code) {
      throw new Error('微信登录失败，请重试');
    }
	console.log(loginRes.code);
	console.log(nickName.value);
    // 发送登录凭证和昵称到后端进行微信认证
    uni.showLoading({
      title: '登录中',
    });

    const authRes = await new Promise((resolve, reject) => {
      uni.request({
        url: 'http://122.51.231.155:8080/users/auth', // 替换为实际的认证接口
        method: 'POST',
        data: {
          code: loginRes.code,
          nickname: nickName.value,
        },
        success: resolve,
        fail: reject,
      });
    });

    if (authRes.data.code !== 200) {
      throw new Error(authRes.data.message || '登录失败，请重试');
    }

    const { token, user } = authRes.data.data;

    // 保存 token 和用户信息到本地存储
    uni.setStorageSync('token', token);
    uni.setStorageSync('userInfo', {
      nickName: user.nickname,
      avatarUrl: user.avatar_url || '/static/images/index/background_img.jpg',
    });

    // 更新 Pinia 用户存储
    userStore.setUid(user.id);
    userStore.setIsLoggedIn(true);
    //userStore.setUserInfo({
    //  id: user.id,
    //  nickname: user.nickname,
    //  avatar_url: user.avatar_url,
    //});

    uni.showToast({
      title: '登录成功',
      icon: 'success',
      duration: 2000,
    });

    // 如果用户选择了头像，并且头像不是默认头像，则上传头像
    if (avatarUrl.value && avatarUrl.value !== '/static/images/index/background_img.jpg') {
      await uploadAvatar(token, avatarUrl.value);
    }

    // 跳转到首页
    login();
  } catch (error) {
    uni.showToast({
      title: error.message || '登录失败，请重试',
      icon: 'none',
      duration: 2000,
    });
  } finally {
    uni.hideLoading();
  }
}

// 上传头像
function uploadAvatar(token, filePath) {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: 'http://122.51.231.155:8080/users/avatar', // 替换为实际的头像上传接口
      filePath: filePath,
      name: 'avatar', // 根据后端要求，字段名为 "avatar"
      header: {
        Authorization: `Bearer ${token}`, // 使用 Bearer Token 认证
      },
      formData: {
        // 如果需要其他表单数据，可以在这里添加
      },
      success: (uploadFileRes) => {
        let data;
        try {
          data = JSON.parse(uploadFileRes.data);
        } catch (e) {
          return reject(new Error('上传响应格式错误'));
        }

        if (data.code === 200 && data.avatar_url) {
          // 更新本地存储中的头像URL
          uni.setStorageSync('userInfo', {
            nickName: nickName.value,
            avatarUrl: data.avatar_url,
          });

          // 更新 Pinia 用户存储中的头像URL
          //userStore.setUserInfo({
          //  ...userStore.userInfo,
          //  avatar_url: data.avatar_url,
          //});

          uni.showToast({
            title: '头像上传成功',
            icon: 'success',
            duration: 2000,
          });

          resolve();
        } else {
          reject(new Error(data.message || '头像上传失败'));
        }
      },
      fail: (error) => {
        reject(new Error(error.errMsg || '头像上传失败'));
      },
    });
  });
}

// 跳转到首页
function login() {
  uni.switchTab({
    url: '/pages/my_index/my_index',
  });
}
</script>

<style scoped>
.container {
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
