<template>
  <view class="container">
    <view class="avatarUrl">
      <button type="balanced" open-type="chooseAvatar" @chooseavatar="onChooseAvatar">
        <image :src="avatarUrl || '/static/images/index/background_img.jpg'" class="refreshIcon"></image>
        <!-- 使用占位符 -->
      </button>
    </view>
    <view class="nickname">
      <text>{{ $t('Nickname') }}</text>
      <input type="nickname" class="weui-input" :value="nickName" @blur="bindBlur"
             :placeholder="$t('enterNickname')" @input="bindInput" />
    </view>


    <view class="btn">
      <view class="btn-sub" @click="onSubmit">{{ $t('save') }}</view>
    </view>
  </view>
</template>

<script setup>
import {
  ref,
  computed
} from 'vue';
import {
  useI18n
} from 'vue-i18n';
import {
  useUserStore
} from '../../stores/user'; // 引入用户存储

// 国际化
const {
  t
} = useI18n();

// Pinia用户存储
const userStore = useUserStore();

// 响应式数据
const avatarUrl = ref('/static/images/index/background_img.jpg');
const nickName = ref('');

// 绑定昵称输入框失去焦点时的事件
function bindBlur(e) {
  nickName.value = e.detail.value; // 获取微信昵称
}

// 绑定昵称输入框输入时的事件
function bindInput(e) {
  console.log("绑定昵称");
  console.log(e.detail.value);
  nickName.value = e.detail.value; // 实时更新昵称
}

// 选择头像
function onChooseAvatar(e) {
  let {
    avatarUrl: newAvatarUrl
  } = e.detail;

  if (!newAvatarUrl) {
    uni.showToast({
      title: '请选择头像',
      icon: 'none',
      duration: 2000,
    });
    return;
  }
  avatarUrl.value = newAvatarUrl;

  uni.showToast({
    title: '头像已更新',
    icon: 'success',
    duration: 2000,
  });
}

async function onSubmit() {
  if (nickName.value === '') {
    uni.showToast({
      icon: 'none',
      title: t('needNickname'),
    });
    return false;
  }

  try {
    // 显示加载提示
    uni.showLoading({
      title: '登录中',
    });

    // 调用 user.js 中的 login 函数
    await userStore.login(nickName.value, avatarUrl.value);

    // 显示成功提示
    uni.showToast({
      title: '登录成功',
      icon: 'success',
      duration: 2000,
    });

    // 跳转到首页
    login();
  } catch (error) {
    // 显示错误提示
    uni.showToast({
      title: error.message || '登录失败，请重试',
      icon: 'none',
      duration: 2000,
    });
  } finally {
    // 隐藏加载提示
    uni.hideLoading();
  }
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
  object-fit: cover;
  /* 保证图片不会变形 */
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
