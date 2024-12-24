// stores/food_list.js
import { defineStore } from 'pinia';
import { reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n'; // 导入 vue-i18n

export const useFoodListStore = defineStore('foodList', () => {
    // 定义食物列表状态（始终使用英文名称）
    const foodList = reactive([
        // 初始食物项已移除
        // 可以预先添加更多食物项，名称为英文
    ]);

    // 定义一个加载状态
    const loaded = ref(false);

    // 存储所有后端的食物名和对应的ID，包含中英文名称和 image_url
    const availableFoods = reactive([]);

    // 获取语言设置
    const { locale } = useI18n(); // 使用 vue-i18n 获取语言

    // 从后端获取食物ID和姓名的函数（获取中英文名称和 image_url）
    const fetchAvailableFoods = async () => {
        try {
            // 分别请求英文和中文名称
            const [resEn, resZh] = await Promise.all([
                new Promise((resolve, reject) => {
                    uni.request({
                        url: `http://122.51.231.155:8080/foods/names`,
                        method: 'GET',
                        data: { lang: 'en' },
                        success: resolve,
                        fail: reject
                    });
                }),
                new Promise((resolve, reject) => {
                    uni.request({
                        url: `http://122.51.231.155:8080/foods/names`,
                        method: 'GET',
                        data: { lang: 'zh' },
                        success: resolve,
                        fail: reject
                    });
                })
            ]);

            if (resEn.statusCode === 200 && resZh.statusCode === 200) {
                // 创建一个映射，根据ID合并中英文名称和 image_url
                const enFoods = resEn.data; // [{ id: 1, name: 'Apple', image_url: '...' }, ...]
                const zhFoods = resZh.data; // [{ id: 1, name: '苹果', image_url: '...' }, ...]

                const foodMap = {};

                // 合并英文
                enFoods.forEach(food => {
                    foodMap[food.id] = {
                        id: food.id,
                        name_en: food.name,
                        name_zh: '',
                        image_url: food.image_url
                    };
                });

                // 合并中文
                zhFoods.forEach(food => {
                    if (foodMap[food.id]) {
                        foodMap[food.id].name_zh = food.name;
                        // 如果中文返回也有 image_url，并且你想以中文优先，可以取消下行注释
                        // foodMap[food.id].image_url = food.image_url;
                    } else {
                        foodMap[food.id] = {
                            id: food.id,
                            name_en: '',
                            name_zh: food.name,
                            image_url: food.image_url
                        };
                    }
                });

                // 将映射转换为数组，并过滤掉没有名称的食物
                const combinedFoods = Object.values(foodMap).filter(food => food.name_en || food.name_zh);

                // 使用 splice 方法更新 reactive 数组
                availableFoods.splice(0, availableFoods.length, ...combinedFoods);
            } else {
                console.error('获取食物列表失败:', resEn.data.error, resZh.data.error);
                uni.showToast({
                    title: '获取食物列表失败',
                    icon: 'none',
                    duration: 2000,
                });
            }
        } catch (err) {
            console.error('网络错误:', err);
            // 提供降级数据，确保包含 image_url
            availableFoods.splice(0, availableFoods.length,
                {
                    id: 1,
                    name_en: 'Apple',
                    name_zh: '苹果',
                    image_url: 'https://example.com/apple.jpg'
                },
                {
                    id: 2,
                    name_en: 'Banana',
                    name_zh: '香蕉',
                    image_url: 'https://example.com/banana.jpg'
                }
            );
            uni.showToast({
                title: '网络错误，无法获取食物列表',
                icon: 'none',
                duration: 2000,
            });
        }
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

    // 根据id获取食物名
    const getFoodName = (id) => {
        const food = availableFoods.find(food => food.id === id);
        if (food) {
            return locale.value === 'zh-Hans' ? food.name_zh : food.name_en;
        } else {
            return '';
        }
    };

    // 根据id获取食物的 image_url
    const getFoodImageUrl = (id) => {
        const food = availableFoods.find(food => food.id === id);
        return food ? food.image_url : '';
    };

    // 辅助函数：从字符串中提取数字
    function extractNumber(value) {
        if (typeof value === 'string') {
            const num = parseFloat(value.replace(/[^\d.]/g, ''));
            return isNaN(num) ? 0 : num;
        }
        return Number(value) || 0;
    }

    // 新增函数：向后端请求计算营养和碳排放的数据
    const calculateNutritionAndEmission = async () => {
        try {
            // 准备请求数据
            const requestData = foodList.map(food => ({
                id: Number(food.id),
                price: extractNumber(food.price),
                weight: extractNumber(food.weight)
            }));

            if (requestData.length === 0) {
                uni.showToast({
                    title: '请先添加食物',
                    icon: 'none',
                    duration: 2000,
                });
                return;
            }

            // 发送POST请求
            const response = await new Promise((resolve, reject) => {
                uni.request({
                    url: 'http://122.51.231.155:8080/foods/calculate',
                    method: 'POST',
                    data: requestData,
                    header: {
                        'Content-Type': 'application/json'
                    },
                    success: resolve,
                    fail: reject
                });
            });

            if (response.statusCode === 200) {
                const responseData = response.data;
                // 更新 foodList 中的每个食物项，添加返回的数据
                responseData.forEach(item => {
                    const food = foodList.find(f => Number(f.id) === item.id);
                    if (food) {
                        food.emission = item.emission;
                        food.calories = item.calories;
                        food.protein = item.protein;
                        food.fat = item.fat;
                        food.carbohydrates = item.carbohydrates;
                        food.sodium = item.sodium;
                        console.log('更新食物数据:', food.emission);
                    }
                });
                // 保存更新后的数据
                saveFoodList();
            } else {
                console.error('计算失败:', response.data.error);
                uni.showToast({
                    title: '计算失败',
                    icon: 'none',
                    duration: 2000,
                });
            }
        } catch (err) {
            console.error('请求失败', err);
            uni.showToast({
                title: '请求失败',
                icon: 'none',
                duration: 2000,
            });
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
        availableFoods,
        fetchAvailableFoods,
        getFoodName,
        getFoodImageUrl, // 导出新添加的函数
        calculateNutritionAndEmission // 导出新添加的函数
    };
});