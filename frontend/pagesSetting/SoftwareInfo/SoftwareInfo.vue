<template>
  <view class="settings">
    <view class="icon_holder">
      <img src="/static/images/settings/logo.png" class="delete-icon" />
    </view>
    <view class="text_holder" :style="{ fontSize: fontSizeBig + 'px' }">
      <text style="font-weight: bold;">DEC 可持续饮食助手</text>
    </view>
    <view class="text_holder" :style="{ fontSize: fontSize + 'px' }">
      <text>version test 0.0.1</text>
    </view>

    <view class="list">
      <view class="list-item">
        <text>评分</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item">
        <text>进入门户</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item">
        <text>更新记录</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="switchtoBase">
        <text>开源库</text>
        <text class="arrow">></text>
      </view>
      <view v-if="showBase" class="list-item" @click="switchtoBase">
        <text>https://github.com/Alchuang22-dev/DEC_sustainable_diet_helper</text>
      </view>
      <view class="divider"></view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { ref } from 'vue'

/* ----------------- Setup ----------------- */
const fontSize = ref(16)
const fontSizeBig = fontSize.value + 4
const showBase = ref(false)

/* ----------------- Methods ----------------- */
function switchtoBase() {
  showBase.value = !showBase.value
}
</script>

<style scoped>
.icon_holder {
  text-align: center;
  padding-top: 50px;
}

.delete-icon {
  width: 100px;
  height: 100px;
}

.text_holder {
  text-align: center;
  padding-top: 10px;
}
.settings {
  height: 100%;
  background: #f8f8f8;
}

.header-actions button {
  background: none;
  border: none;
}
.list {
  margin-top: 30px;
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