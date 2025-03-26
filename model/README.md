# SimpleTorch

一个使用简单 CNN 架构实现 food101 图像分类的 PyTorch 项目。
在我们的项目中，它将用于食物识别的需求模块。

## 项目结构

```
SimpleTorch/
├── README.md          # 项目文档
├── requirements.txt   # 项目依赖
├── main.py           # 主训练脚本
├── model.py          # CNN 模型定义
├── model-v2.py       # ResNet 微调模型定义
├── test.py           # CNN 测试脚本
├── test-v2.py        # ResNet 测试脚本
├── configs/          # 配置文件
│   ├── __init__.py
│   └── config.py     # 配置类
├── utils/            # 工具函数
│   ├── __init__.py
│   ├── data_utils.py # 数据加载和处理
│   ├── train_utils.py # 训练工具
│   └── visualization.py # 可视化工具
├── tests/            # 测试文件
│   └── __init__.py
├── checkpoints/      # 模型检查点
└── data/            # 数据集目录
```

## 环境配置

1. 创建新的 conda 环境：
```bash
conda create -n simpletorch python=3.9
conda activate simpletorch
```
或者(linux推荐)
```bash
python3 -m venv simpletorch
source simpletorch\\bin\\activate
```

2. 安装依赖：
```bash
pip install -r requirements.txt
```

## 模型架构

### CNN 模型结构（model.py）：
- 输入：3x32x32 RGB 图像
- 卷积层1：3->16 通道，3x3 卷积核
- 卷积层2：16->32 通道，3x3 卷积核
- 全连接层：32*8*8 -> 101 个类别

### ResNet 模型结构（model-v2.py）
- 正在调试

## 配置系统

项目使用模块化的配置系统：
- `TrainingConfig`：训练参数（批次大小、学习率等）
- `ModelConfig`：模型架构参数
- `DataConfig`：数据集和数据加载参数

配置示例：
```python
from configs.config import Config

config = Config()
config.training.batch_size = 32
config.training.learning_rate = 0.01
```

## 数据集和数据加载

### 数据预处理
```python
from utils.data_utils import get_data_transforms, load_dataset

transform = get_data_transforms(train=True)
dataset = load_dataset(config.data)
```

### 数据集可视化
运行以下命令查看数据集样本和分布：
```bash
python show_dataset.py
```
这将生成两个可视化文件：
- `food101_samples.png`：数据集样本图像
- `food101_distribution.png`：类别分布可视化

## 训练

1. 运行训练：
```bash
python main.py
```

2. 脚本将执行以下操作：
   - 下载 FOOD-101 数据集
   - 训练模型
   - 显示训练进度
   - 保存模型检查点
   - 生成训练可视化

## 可视化工具

项目包含多个可视化工具：
- 数据集样本可视化
- 类别分布绘图
- 训练历史绘图
- 模型预测可视化

## 环境要求

- Python 3.9
- PyTorch 2.6.0
- torchvision
- numpy
- matplotlib

## 测试
测试CNN模型
```bash
python test.py
```
测试ResNet模型（使用单张图片测试）
```bash
python test_v2.py
```


## 许可证

MIT 许可证 