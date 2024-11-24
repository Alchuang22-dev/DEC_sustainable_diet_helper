<template>
  <view class="settings">
    <view class="list">
      <view class="list-item">
        <text>使用个性化推荐算法</text>
        <switch class="toggle-button" @change="toggleReco" :checked="isReco" />
      </view>
      <view class="divider"></view>
	  <view class="list-item">
	    <text>黑名单</text>
	    <text class="arrow">></text>
	  </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue';

const isReco = ref(false);

function toggleReco(event) {
  isReco.value = event.detail.value;
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
