<template>
  <view class="settings">
    <view class="header">
      <view @click="goBack" class="back-icon">返回</view>
      <text class="title">user_test</text>
      <view class="header-actions">
        <button class="menu-icon"></button>
        <button class="camera-icon"></button>
      </view>
    </view>

    <view class="list">
      <view class="list-item" @click="navigateTo('Idinfo')">
        <text>账号名称</text>
        <text class="numbers right-align">{{ userData.username }}</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="navigateTo('Idinfo')">
        <text>账号ID</text>
        <text class="numbers right-align">{{ userData.accountId }}</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="navigateTo('Idinfo')">
        <text>手机号</text>
        <text class="numbers right-align">{{ userData.phoneNumber }}</text>
        <text class="arrow">></text>
      </view>

      <view class="list-item" @click="navigateTo('newsSetting')">
        <text>邮箱</text>
        <text class="numbers right-align">{{ userData.email }}</text>
        <text class="arrow">></text>
      </view>
      <view class="divider"></view>
      <view class="list-item" @click="navigateTo('recoSetting')">
        <text>账户密码</text>
        <text class="numbers right-align">已设置</text>
        <text class="arrow">></text>
      </view>
      <view class="divider"></view>
      <view class="list-item centered" @click="navigateTo('Bend')">
        <text>安全中心</text>
      </view>
      <view class="list-item centered red-text" @click="navigateTo('DeleteId')">
        <text>冻结账号</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const userData = ref({
  username: '',
  accountId: '',
  phoneNumber: '',
  email: ''
});

function goBack() {
  uni.navigateBack();
}

function navigateTo(link) {
  uni.navigateTo({
    url: `/pagesSetting/${link}/${link}`,
  });
}

// Simulate fetching data from backend
onMounted(() => {
  // Example of an API call to fetch user data
  setTimeout(() => {
    userData.value = {
      username: 'user_test',
      accountId: '00001',
      phoneNumber: '(+86)12345678900',
      email: '2901@mails.com'
    };
  }, 500); // Simulate network delay
});
</script>

<style scoped>
.settings {
  height: 100%;
  background: #f8f8f8;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #ebebeb;
}
.title {
  font-size: 18px;
  font-weight: bold;
}
.header-actions button {
  background: none;
  border: none;
}
.list {
  margin-top: 10px;
}
.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: #fff;
  border-bottom: 1px solid #ebebeb;
}
.list-item.centered {
  justify-content: center;
}
.numbers {
  color: #ccc;
}
.numbers.right-align {
  margin-left: auto;
}
.arrow {
  color: #ccc;
}
.divider {
  height: 1px;
  background-color: #ebebeb;
  margin: 10px 0;
}
.red-text {
  color: red;
}
</style>
