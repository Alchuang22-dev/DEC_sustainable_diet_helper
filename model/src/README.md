# model use Food2K and Ingredient-201
## source
关于Food2K部分的训练方法，请参见https://github.com/JonnyKan/Food2K.git（目前只对论文提及的训练方法进行了复现）
复现的模型在https://cloud.tsinghua.edu.cn/f/46a3fcf759604bacae95/?dl=1
关于已经预训练好的模型，建议通过Food2K官方渠道下载
|  CNN   | link  |
|  ----  | ----  |
| vgg16  | [google](https://drive.google.com/file/d/1r4CQEfCkwLSKz5QdZJGABldercUo5BtF/view?usp=sharing)/[baidu](https://pan.baidu.com/s/1-nI6fodmmzqz9OVqh0yvUw)(Code: puuy)|
| resnet50  | [google](https://drive.google.com/file/d/1h87m392fJIxrADTe8GMH7pibP0rjWu-k/view?usp=sharing)/[baidu](https://pan.baidu.com/s/1WY7VsCBTJt2mL9n3Gdl8Mg)(Code: 5eay) |
| resnet101  | [google](https://drive.google.com/file/d/1_xM2qv1NIjev8voYjXLhfnxDzvFNB85q/view?usp=sharing)/[baidu](https://pan.baidu.com/s/1mEO7KyJFHrkpB5G0Aj6oWw)(Code: yv1o) |
| resnet152  | [google](https://drive.google.com/file/d/1YG_gW6NftjX06-i3bCCYQhlnDo2mUoLn/view?usp=sharing)/[baidu](https://pan.baidu.com/s/1-3LikXkDEvbxQur6n-FUJw)(Code: 22zw) |
| densenet161  | [google](https://drive.google.com/file/d/17PAUHmo1vIM9b4SlbpnLnwp1a5MH9Vem/view?usp=sharing)/[baidu](https://pan.baidu.com/s/1UllqjTJMAQEnGFVgzf6-nQ)(Code: bew5) |
| inception_resnet_v2  | [google](https://drive.google.com/file/d/16PuZRuUB-YFKZT8JWycaay3JdfTlCoVK/view?usp=sharing)/[baidu](https://pan.baidu.com/s/1_974E4eZRzKubemLIQlOHA)(Code: xa8r) |
| senet154  | [google](https://drive.google.com/file/d/1FGs7gH1fYybr3sKB6q4lRl36bX5wiLXw/view?usp=sharing)/[baidu](https://pan.baidu.com/s/1tHpFFSm2AySRjDZ4BTtboQ)(Code: kwzf) |
## working schedule
+ Layer 1: 基于Food2K预训练模型的熟食调整
    + 未来将与本团队自己训练的模型进行交叉训练，以符合本团队数据库中的食品
+ Layer 2: 基于Ingredient-201数据库的生鲜食品识别模型（正在施工）
+ Layer 3: 食谱-食材映射和自主计算模型（正在施工）
