<template>
  <view class="settings">
    <view class="header">
      <view @click="goBack" class="back-icon">返回</view>
      <text class="title">软件设置</text>
      <view class="header-actions">
        <button class="menu-icon"></button>
        <button class="camera-icon"></button>
      </view>
    </view>

    <view class="list">
      <view class="divider">
        <text class="divider-text">显示</text>
      </view>
      <view class="list-item" @click="navigateTo('Fonts')">
        <text>字体大小</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="navigateTo('language')">
        <text>多语言</text>
        <text class="arrow">></text>
      </view>
      <view class="divider">
        <text class="divider-text">存储与数据</text>
      </view>
      <view class="list-item" @click="navigateTo('DeleteData')">
        <text>清除数据</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="navigateTo('Storage')">
        <text>存储空间</text>
        <text class="arrow">></text>
      </view>
      <view class="divider">
        <text class="divider-text">辅助功能</text>
      </view>
      <view class="list-item" @click="navigateTo('DevelopingFunc')">
        <text>开发功能</text>
        <text class="arrow">></text>
      </view>
      <view class="divider">
        <text class="divider-text">关于软件</text>
      </view>
      <view class="list-item" @click="navigateTo('FeedBack')">
        <text>帮助与反馈</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="navigateTo('SoftwareInfo')">
        <text>软件信息</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="navigateTo('ConnectUs')">
        <text>联系我们</text>
        <text class="arrow">></text>
      </view>
    </view>
  </view>
</template>

<script setup>
function goBack() {
  uni.navigateBack();
}

function navigateTo(link) {
  uni.navigateTo({
    url: `/pagesSetting/${link}/${link}`,
  });
}
const dividerText = ""; // Set this to the desired text to add text to dividers
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
  position: relative;
  height: 20px;
  margin: 20px 0;
}
.divider::before {
  content: "";
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background-color: #ebebeb;
  transform: translateY(-50%);
}
.divider-text {
  position: relative;
  z-index: 1;
  background-color: #f8f8f8;
  padding: 0 10px;
  text-align: center;
  font-size: 14px;
  color: #888;
}
.red-text {
  color: red;
}
</style>

