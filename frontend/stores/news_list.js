// stores/news_list.js
import { defineStore } from 'pinia';

export const useNewsStore = defineStore('news', {
  state: () => ({
    allNewsItems: [
      {
        title: "国际氢能联盟和麦肯锡联合发布《氢能洞察2024》",
        description: "环保要闻  |  双碳指挥  |  刚刚",
        link: "news_detail",
        categories: ["环保要闻", "双碳指挥"],
      },
      {
        title: "把自然讲给你听 | 什么是森林？",
        description: "环保科普  |  环保科普365  |  1小时前",
        image: "",
        link: "https://mp.weixin.qq.com/s/mzFR2d-17men_Lm297fweQ",
        categories: ["环保科普", "环保科普365"],
      },
      {
        title: "视频 | 垃圾分类",
        description: "环保科普  |  环保科普365  |  2024-10-14",
        video: true,
        link: "video_detail",
        categories: ["环保科普", "环保科普365"],
      },
      {
        title: "联合国发布2024气候计划",
        description: "环保要闻  |  环保科普365  |  2024-10-14",
        video: true,
        link: "news_detail",
        categories: ["环保要闻", "环保科普365"],
      },
      {
        title: "专栏 | 寂静的春天",
        description: "环保专栏  |  爱读夜  |  2024-10-14",
        video: true,
        link: "news_detail",
        categories: ["环保专栏", "爱读夜"],
      },
      {
        title: "社团招新 | 根与芽2025",
        description: "环保专栏  |  公益事业  |  2024-6-16",
        video: true,
        link: "news_detail",
        categories: ["环保专栏", "公益事业"],
      },
    ],
    rawNewsData: [], // 原始后端数据
    selectedSection: '全部', // 当前选中的分类
    isRefreshing: false, // 是否正在刷新页面
  }),
  getters: {
    filteredNewsItems(state) {
      if (state.selectedSection === '全部') {
        return state.allNewsItems;
      }
      return state.allNewsItems.filter((item) =>
        item.categories.includes(state.selectedSection)
      );
    },
  },
  actions: {
    setSection(section) {
		console.log("setSection");
      this.selectedSection = section;
    },
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
      try {
        const mockResponse = {
          data: [
            {
              id: 1,
              form: 'video',
              newsSrc: 'http://vjs.zencdn.net/v/oceans.mp4',
              imgsSrc: '',
              tabs: ['环境保护', '环保要闻'],
              time: '2024-4-17',
              newsName: '测试视频',
              authorName: 'user_test',
              authorAvatar: '',
              newsinfo: '测试测试测试测试测试',
              newsbody: '测试内容',
              likeCount: 1001,
              shareCount: 37,
              favoriteCount: 897,
              followCount: 189,
              dislikeCount: 100,
              type: 'main',
            },
          ],
        };

        this.rawNewsData = mockResponse.data;
        this.rawNewsData.forEach(this.convertVideoToItems);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    },
    convertVideoToItems(video) {
      if (video.type === 'main') {
        const categories = video.tabs; // 使用 tabs 作为 categories
        let link = '';
        if (video.form === 'web') {
          link = video.newsSrc;
        } else if (video.form === 'news') {
          link = 'news_detail';
        } else if (video.form === 'video') {
          link = 'video_detail';
        }

        this.allNewsItems.push({
          id: video.id,
          link: link,
          image: video.imgsSrc || '',
          title: video.newsName,
          description: `${video.tabs.join(' | ')} | ${video.time}`,
          info: `阅读量: ${video.followCount} | 点赞量: ${video.likeCount}`,
          form: video.form,
          categories: video.tabs, // 添加 categories 字段，在之后可以修改为其他方法
        });
      }
    },
  },
});
