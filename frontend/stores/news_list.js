// stores/news_list.js
import { defineStore } from 'pinia';

export const useNewsStore = defineStore('news', {
  state: () => ({
    allNewsItems: [
    ],
    rawNewsData: [], // 原始后端数据
    //selectedSection: '全部', // 当前选中的分类
    isRefreshing: false, // 是否正在刷新页面
  }),
  getters: {
    filteredNewsItems(state) {
    //  if (state.selectedSection === '全部') {
        return state.allNewsItems;
    //  }
    //  return state.allNewsItems.filter((item) =>
    //    item.categories.includes(state.selectedSection)
    //  );
    },
  },
  actions: {
    //setSection(section) {
	//	console.log("setSection");
    //  this.selectedSection = section;
    //},
	refreshNews() {
	  this.isRefreshing = true;
	  console.log("in refreshing");
	    setTimeout(() => {
	      // 使用 Vue.set 触发响应式更新
	      this.allNewsItems = this.allNewsItems.sort(() => Math.random() - 0.5);
	      this.isRefreshing = false;
	    }, 1000);
	},
    async fetchNews() {
		console.log('正在获取新闻');
	  this.allNewsItems = [];
      try {
        const mockResponse = {
          data: [
            {
				id: '', //唯一识别号
				authoravatar: '', //作者头像
				authorname: '测试作者', //作者姓名
				authorid: '', //作者识别号
				savetime: '2024-12-12', //发布时间
				title: '测试新闻', //文章题目
				description: '测试描述', //文章描述
				components: [], //文章内容
				likeCount: 1001,
				shareCount: 37,
				favoriteCount: 897,
				followCount: 189, //可选
				dislikeCount: 199,
				type: 'main', //可选
            },
          ],
        };

        this.rawNewsData = mockResponse.data;
        this.rawNewsData.forEach(this.convertnewsToItems);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    },
    convertnewsToItems(news) {
	  console.log('正在将新闻放入数组...');
      if (news.type === 'main') {
        //const categories = news.tabs; // 使用 tabs 作为 categories
        let link = '';
        link = 'news_detail';

        this.allNewsItems.push({
          id: news.id,
          link: link,
          title: news.title,
          description: `${news.savetime}`,
          info: `阅读量: ${news.followCount} | 点赞量: ${news.likeCount}`,
          form: news.form,
          //categories: news.tabs, // 添加 categories 字段，在之后可以修改为其他方法
        });
      }
    },
  },
});
