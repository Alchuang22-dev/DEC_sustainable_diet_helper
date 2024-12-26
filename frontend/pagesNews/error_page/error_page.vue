<template>
  <view class="container">
    <view class="icon_holder">
      <img src="/static/images/settings/trash.png" class="delete-icon"></img>
    </view>
    <view class="text_holder" :style="{ fontSize: fontSizeBig + 'px' }">
      <text style="font-weight: bold;">错误</text>
    </view>
    <view class="text_holder" :style="{ fontSize: fontSize + 'px' }">
      <text>该功能已停止运营</text>
    </view>
    <view class="button_holder">
      <button class="confirm-button" @click="toIndex">返回主页</button> 
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue';

// 定义字体大小
const fontSize = ref(16);
const fontSizeBig = fontSize.value + 4;

function toIndex() {
  // TODO：这一部分没找到对应的后端
	 uni.switchTab({
	    url: "/pages/index/index", // 请确保这个路径与您的 TabBar 配置中的路径一致
	    success: () => {
	      console.log("成功返回主页");
	    },
	    fail: (err) => {
	      console.error("返回主页失败:", err);
	    }
	  });
}
</script>

<style>
/* 设置容器为 Flex 布局 */
.container {
  display: flex;
  flex-direction: column;
  height: 100vh; /* 设置页面高度为全屏 */
  justify-content: space-between; /* 自动分布内容，按钮会被推到底部 */
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
  margin-top: auto; /* 将按钮推到底部 */
  padding-bottom: 30px; /* 底部留一些空间 */
}

.confirm-button {
  width: 40%;
  padding: 10px;
  background-color: red;
  color: #fff;
  text-align: center;
  border-radius: 5px;
}
</style>