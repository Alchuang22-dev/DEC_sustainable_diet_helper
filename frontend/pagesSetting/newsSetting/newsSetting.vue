<template>
  <view class="settings">
    <view class="list">
      <view class="list-item">
        <text>消息通知</text>
        <switch class="toggle-button" @change="toggleImpart" :checked="isImpart" />
      </view>
      <view class="list-item">
        <text>消息提示音</text>
        <switch class="toggle-button" @change="toggleVolume" :checked="isVolume" />
      </view>
      <view class="list-item">
        <text>在微信内发送消息</text>
        <switch class="toggle-button" @change="toggleInWechat" :checked="isInWechat" />
      </view>
      <view class="divider"></view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue';

const isImpart = ref(false);
const isVolume = ref(false);
const isInWechat = ref(false);

function toggleImpart(event) {
  isImpart.value = event.detail.value;
}

function toggleVolume(event) {
  isVolume.value = event.detail.value;
}

function toggleInWechat(event) {
  isInWechat.value = event.detail.value;
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
.toggle-button {
  margin-left: auto;
}
</style>
