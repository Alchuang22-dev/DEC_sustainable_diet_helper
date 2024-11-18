<template>
  <view class="settings">
    <view class="header">
      <view @click="goBack" class="back-icon">返回</view>
      <text class="title">用户设置</text>
      <view class="header-actions">
        <button class="menu-icon"></button>
        <button class="camera-icon"></button>
      </view>
    </view>

    <view class="list">
      <view class="list-item">
        <text>账号与安全</text>
        <text class="arrow">></text>
      </view>
      <view class="divider"></view>
      <view class="list-item">
        <text>消息设置</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item">
        <text>个性化推荐设置</text>
        <text class="arrow">></text>
      </view>
      <view class="divider"></view>
      <view class="list-item">
        <text>通信权限</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item">
        <text>个人信息与权限</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item">
        <text>个人信息收集清单</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item">
        <text>第三方信息收集清单</text>
        <text class="arrow">></text>
      </view>
      <view class="divider"></view>
      <view class="list-item centered">
        <text>账户绑定</text>
      </view>
      <view class="list-item centered red-text">
        <text>账户注销</text>
      </view>
    </view>
  </view>
</template>

<script setup>
function goBack() {
  uni.navigateBack();
}

function navigateTo(url) {
  uni.navigateTo({ url });
}
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
