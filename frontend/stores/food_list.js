// stores/food_list.js
import { defineStore } from 'pinia';
import { reactive, ref } from 'vue';

export const useFoodListStore = defineStore('foodList', () => {
    // 定义食物列表状态
    const foodList = reactive([
        {
            name: "西红柿",
            id: "1",
            weight: "1kg",
            price: "5元",
            transportMethod: "陆运",
            foodSource: "本地",
            image: "",
            isAnimating: false,
            emission: 0,
            calories: 100,
            protein: 200,
            fat: 400,
            carbohydrates: 300,
            sodium: 500
        },
        // 可以预先添加更多食物项
    ]);

    // 定义一个加载状态
    const loaded = ref(false);

    // 存储所有后端的食物名和对应的ID，改为 reactive 数组
    const availableFoods = reactive([]);

    // 从后端获取食物ID和姓名的函数
    const fetchAvailableFoods = () => {
        const lang = 'zh'; // 根据应用的语言设置动态获取
        uni.request({
            url: `http://122.51.231.155:8080/foods/names`,
            method: 'GET',
            data: {
                lang: lang
            },
            success: (res) => {
                if (res.statusCode === 200) {
                    // 使用 splice 方法更新 reactive 数组
                    availableFoods.splice(0, availableFoods.length, ...res.data);
                } else {
                    console.error('获取食物列表失败:', res.data.error);
                    uni.showToast({
                        title: '获取食物列表失败',
                        icon: 'none',
                        duration: 2000,
                    });
                }
            },
            fail: (err) => {
                console.error('网络错误:', err);
                // 使用 splice 方法更新 reactive 数组
                availableFoods.splice(0, availableFoods.length,
                    {
                        "id": 1,
                        "name": "苹果"  // 或 "Apple"，取决于 lang 参数
                    },
                    {
                        "id": 2,
                        "name": "香蕉"  // 或 "Banana"
                    }
                );
                console.log('availableFoods:', availableFoods);
                uni.showToast({
                    title: '网络错误，无法获取食物列表',
                    icon: 'none',
                    duration: 2000,
                });
            }
        });
    };

    // 加载食物列表从本地存储
    const loadFoodList = () => {
        const storedFoodList = uni.getStorageSync('foodDetails');
        if (storedFoodList && storedFoodList.length > 0) {
            foodList.splice(0, foodList.length, ...storedFoodList.map(food => ({
                ...food,
                isAnimating: false
            })));
        }
        loaded.value = true;
    };

    // 保存食物列表到本地存储
    const saveFoodList = () => {
        uni.setStorageSync('foodDetails', foodList);
    };

    // 添加一个新的食物项
    const addFood = (newFood) => {
        foodList.push(newFood);
        saveFoodList();
    };

    // 删除指定索引的食物项
    const deleteFood = (index) => {
        foodList.splice(index, 1);
        saveFoodList();
    };

    // 更新指定索引的食物项
    const updateFood = (index, updatedFood) => {
        if (foodList[index]) {
            Object.assign(foodList[index], updatedFood);
            saveFoodList();
        }
    };

    // 初始化加载
    if (!loaded.value) {
        loadFoodList();
    }

    return {
        foodList,
        addFood,
        deleteFood,
        updateFood,
        saveFoodList,
        loadFoodList,
        loaded,
        availableFoods, // 现在是 reactive 数组
        fetchAvailableFoods
    };
});
