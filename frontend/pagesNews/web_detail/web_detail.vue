<template>
  <view>
	<text>小程序端暂不支持外部跳转，请移至PC端获取更多内容</text>
	<text>需要投稿者请联系DEC公众号</text>
  </view>
</template>

<script>
export default {
  data() {
    return {
      webUrl: '',
    };
  },
  onLoad(options) {
    if (options.url) {
      this.webUrl = decodeURIComponent(options.url);
    }
  },
};
</script>

<style scoped>
web-view {
  width: 100%;
  height: 100vh;
}
</style>
