import os
import torch
import torch.nn as nn
from torch.utils.data import Dataset, DataLoader, ConcatDataset
from torchvision import models, transforms
from PIL import Image
import time
import copy
import numpy as np
from sklearn.metrics import accuracy_score, precision_score, recall_score, f1_score
import itertools
from sklearn.model_selection import ParameterGrid
import json

# 为Food-101数据集创建自定义数据集类
class Food101Dataset(Dataset):
    def __init__(self, data_dir, transform=None, split='train'):
        """
        Food-101数据集加载器
        
        参数:
            data_dir: Food-101数据集根目录
            transform: 图像转换
            split: 'train'或'test'
        """
        self.data_dir = data_dir
        self.transform = transform
        self.split = split
        
        # 加载类别名称
        self.class_names = []
        with open(os.path.join(data_dir, 'meta', 'classes.txt'), 'r') as f:
            for line in f:
                self.class_names.append(line.strip())
        
        # 创建类别到索引的映射
        self.class_to_idx = {cls: i for i, cls in enumerate(self.class_names)}
        
        # 加载图像路径和标签
        self.image_paths = []
        self.labels = []
        
        with open(os.path.join(data_dir, 'meta', f'{split}.txt'), 'r') as f:
            for line in f:
                path = line.strip()
                cls = path.split('/')[0]
                img_path = os.path.join(data_dir, 'images', path + '.jpg')
                
                self.image_paths.append(img_path)
                self.labels.append(self.class_to_idx[cls])
        
        print(f"Food-101 {split}: 加载了 {len(self.image_paths)} 张图像，共 {len(self.class_names)} 个类别")
    
    def __len__(self):
        return len(self.image_paths)
    
    def __getitem__(self, idx):
        img_path = self.image_paths[idx]
        label = self.labels[idx]
        
        try:
            image = Image.open(img_path).convert('RGB')
        except Exception as e:
            print(f"无法加载图像 {img_path}: {e}")
            image = Image.new('RGB', (224, 224))
        
        if self.transform:
            image = self.transform(image)
            
        return image, 0  # 0表示来自Food-101数据集

# 为本地数据库创建自定义数据集类
class CustomFoodDataset(Dataset):
    def __init__(self, root_dir, list_file, name_file, transform=None):
        """
        自定义食物数据集
        
        参数:
            root_dir: 数据集根目录
            list_file: 包含图片路径和标签的文件
            name_file: 包含类别名称的文件
            transform: 图像转换操作
        """
        self.root_dir = root_dir
        self.transform = transform
        
        # 加载类别名称
        self.class_names = []
        self.class_dict = {}  # 用于存储类别ID到索引的映射
        
        with open(os.path.join(root_dir, name_file), 'r', encoding='utf-8') as f:
            for line in f:
                parts = line.strip().split('\t')
                if len(parts) >= 3:
                    class_id = int(parts[0])
                    chinese_name = parts[1]
                    english_name = parts[2]
                    self.class_names.append(f"{chinese_name} ({english_name})")
                    self.class_dict[class_id] = len(self.class_names) - 1
        
        # 加载图片路径和标签
        self.image_paths = []
        self.labels = []
        
        split_type = list_file.split('_')[0]  # 获取train, val或test
        
        with open(os.path.join(root_dir, list_file), 'r') as f:
            for line in f:
                parts = line.strip().split()
                if len(parts) >= 2:
                    img_path = parts[0]
                    label = int(parts[1])
                    
                    # 构建完整的图像路径
                    full_img_path = os.path.join(root_dir, split_type, img_path)
                    
                    if label in self.class_dict:
                        self.image_paths.append(full_img_path)
                        self.labels.append(self.class_dict[label])
                    else:
                        print(f"警告: 忽略标签 {label}，因为它不在类别映射中")
        
        print(f"本地数据库 {split_type}: 加载了 {len(self.image_paths)} 张图像，共 {len(self.class_names)} 个类别")
    
    def __len__(self):
        return len(self.image_paths)
    
    def __getitem__(self, idx):
        img_path = self.image_paths[idx]
        label = self.labels[idx]
        
        try:
            image = Image.open(img_path).convert('RGB')
        except Exception as e:
            print(f"无法加载图像 {img_path}: {e}")
            image = Image.new('RGB', (224, 224))
        
        if self.transform:
            image = self.transform(image)
            
        return image, 1  # 1表示来自本地数据集
    
# 为本地数据库测试集创建自定义数据集类
class CustomFoodTestDataset(Dataset):
    def __init__(self, root_dir, test_list_file, name_file, transform=None):
        """
        自定义食物测试数据集
        
        参数:
            root_dir: 数据集根目录
            test_list_file: 包含测试图片路径和真实标签的文件 (test_truth_list.txt)
            name_file: 包含类别名称的文件
            transform: 图像转换操作
        """
        self.root_dir = root_dir
        self.transform = transform
        
        # 加载类别名称
        self.class_names = []
        self.class_dict = {}  # 用于存储类别ID到索引的映射
        
        with open(os.path.join(root_dir, name_file), 'r', encoding='utf-8') as f:
            for line in f:
                parts = line.strip().split('\t')
                if len(parts) >= 3:
                    class_id = int(parts[0])
                    chinese_name = parts[1]
                    english_name = parts[2]
                    self.class_names.append(f"{chinese_name} ({english_name})")
                    self.class_dict[class_id] = len(self.class_names) - 1
        
        # 加载图片路径和标签
        self.image_paths = []
        self.labels = []
        
        # 从test_truth_list.txt加载测试集标签
        with open(os.path.join(root_dir, test_list_file), 'r') as f:
            for line in f:
                parts = line.strip().split()
                if len(parts) >= 2:
                    img_name = parts[0]  # 例如 000000.jpg
                    label = int(parts[1])
                    
                    # 构建完整的图像路径
                    full_img_path = os.path.join(root_dir, 'test', img_name)
                    
                    if label in self.class_dict:
                        self.image_paths.append(full_img_path)
                        self.labels.append(self.class_dict[label])
                    else:
                        print(f"警告: 忽略标签 {label}，因为它不在类别映射中")
        
        print(f"本地数据库测试集: 加载了 {len(self.image_paths)} 张图像，共 {len(self.class_names)} 个类别")
    
    def __len__(self):
        return len(self.image_paths)
    
    def __getitem__(self, idx):
        img_path = self.image_paths[idx]
        label = self.labels[idx]
        
        try:
            image = Image.open(img_path).convert('RGB')
        except Exception as e:
            print(f"无法加载图像 {img_path}: {e}")
            image = Image.new('RGB', (224, 224))
        
        if self.transform:
            image = self.transform(image)
            
        return image, 1  # 1表示来自本地数据库

# 特征提取器类
class FeatureExtractor(nn.Module):
    def __init__(self, model_path, model_name='resnet50', num_classes=101):
        """
        从预训练模型创建特征提取器
        
        参数:
            model_path: 预训练模型的路径
            model_name: 模型架构
            num_classes: 原模型的类别数
        """
        super(FeatureExtractor, self).__init__()
        
        # 创建基础模型
        if model_name == 'resnet18':
            self.model = models.resnet18(weights=None)
            num_features = self.model.fc.in_features
            self.model.fc = nn.Linear(num_features, num_classes)
        elif model_name == 'resnet50':
            self.model = models.resnet50(weights=None)
            num_features = self.model.fc.in_features
            self.model.fc = nn.Linear(num_features, num_classes)
        elif model_name == 'resnet101':
            self.model = models.resnet101(weights=None)
            num_features = self.model.fc.in_features
            self.model.fc = nn.Linear(num_features, num_classes)
        else:
            raise ValueError(f"不支持的模型: {model_name}")
        
        # 加载预训练权重
        checkpoint = torch.load(model_path, map_location='cpu')
        if 'model_state_dict' in checkpoint:
            self.model.load_state_dict(checkpoint['model_state_dict'])
        else:
            self.model.load_state_dict(checkpoint)
        
        # 提取特征的模型（去掉分类层）
        self.features = nn.Sequential(*list(self.model.children())[:-1])
        
        # 冻结所有参数
        for param in self.features.parameters():
            param.requires_grad = False
    
    def forward(self, x):
        """提取特征"""
        x = self.features(x)
        x = torch.flatten(x, 1)
        return x
    
# 无标签测试集    
class CustomFoodUnlabeledTestDataset(Dataset):
    def __init__(self, root_dir, test_dir='test', transform=None):
        """
        自定义食物无标签测试数据集
        
        参数:
            root_dir: 数据集根目录
            test_dir: 测试集文件夹名称
            transform: 图像转换操作
        """
        self.root_dir = root_dir
        self.test_dir = test_dir
        self.transform = transform
        
        # 加载所有测试图片
        test_path = os.path.join(root_dir, test_dir)
        image_files = [f for f in os.listdir(test_path) if os.path.isfile(os.path.join(test_path, f)) and 
                       f.lower().endswith(('.png', '.jpg', '.jpeg'))]
        
        self.image_paths = [os.path.join(test_path, f) for f in image_files]
        self.image_names = image_files
        
        print(f"无标签测试集: 加载了 {len(self.image_paths)} 张图像")
    
    def __len__(self):
        return len(self.image_paths)
    
    def __getitem__(self, idx):
        img_path = self.image_paths[idx]
        img_name = self.image_names[idx]
        
        try:
            image = Image.open(img_path).convert('RGB')
        except Exception as e:
            print(f"无法加载图像 {img_path}: {e}")
            image = Image.new('RGB', (224, 224))
        
        if self.transform:
            image = self.transform(image)
            
        # 返回图像、虚拟标签和图像名称
        return image, 0, img_name  # 0是虚拟标签

# 二分类器模型
class BinaryClassifier(nn.Module):
    def __init__(self, food101_model_path, custom_model_path, feature_dim=2048):
        """
        使用两个预训练模型作为特征提取器的二分类模型
        
        参数:
            food101_model_path: Food-101模型路径
            custom_model_path: 本地数据库模型路径
            feature_dim: 特征维度
        """
        super(BinaryClassifier, self).__init__()
        
        # 创建特征提取器
        self.food101_extractor = FeatureExtractor(
            model_path=food101_model_path,
            model_name='resnet50',
            num_classes=101
        )
        
        self.custom_extractor = FeatureExtractor(
            model_path=custom_model_path,
            model_name='resnet50',
            num_classes=208  # 假设本地数据库有208个类别
        )
        
        # 二分类器
        self.classifier = nn.Sequential(
            nn.Linear(feature_dim * 2, 512),  # 连接两个特征提取器的输出
            nn.ReLU(),
            nn.Dropout(0.5),
            nn.Linear(512, 2)
        )
    
    def forward(self, x):
        """前向传播"""
        # 从两个特征提取器获取特征
        food101_features = self.food101_extractor(x)
        custom_features = self.custom_extractor(x)
        
        # 连接特征
        combined_features = torch.cat((food101_features, custom_features), dim=1)
        
        # 分类
        output = self.classifier(combined_features)
        
        return output
    
def grid_search_hyperparameters():
    """使用网格搜索寻找最佳超参数组合"""
    
    # 定义超参数搜索空间
    param_grid = {
        'base_model': ['resnet50'],  # 可以添加'resnet18', 'resnet101'等
        'learning_rate': [0.0001 * 0.5, 0.0001, 0.0001 * 2],
        'dropout_rate': [0.3 * 0.8, 0.3, 0.3 * 1.2],
        'batch_size': [128],
        'hidden_size': [256, 512, 1024],
        'optimizer_type': ['adam', 'adamw'],
        'learning_rate': [0.001, 0.0005, 0.0001],
        'weight_decay': [1e-4, 1e-5],
        'scheduler_type': ['step', 'plateau', 'cosine'],
        'batch_size': [16, 32, 64],
        'stage1_epochs': [3, 5],
        'stage2_epochs': [3, 5],
        'data_augment_level': ['standard', 'strong']
    }
    
    # 为了减少搜索空间，这里只选择一部分组合
    # 实际使用时，可以根据需要进行修改
    selected_params = [
        {'base_model': 'resnet50', 'optimizer_type': 'adam', 'scheduler_type': 'step'},
        {'base_model': 'resnet50', 'optimizer_type': 'adamw', 'scheduler_type': 'cosine'},
        {'base_model': 'resnet50', 'optimizer_type': 'adam', 'scheduler_type': 'plateau'}
    ]
    
    # 从参数网格中生成配置
    configs = []
    for base_config in selected_params:
        grid = {k: v for k, v in param_grid.items() if k not in base_config}
        for params in ParameterGrid(grid):
            full_config = base_config.copy()
            full_config.update(params)
            configs.append(full_config)
    
    # 保存所有配置
    with open('hyperparam_configs.json', 'w') as f:
        json.dump(configs, f, indent=2)
    
    print(f"将搜索 {len(configs)} 种超参数组合")
    
    # 搜索最佳超参数
    best_config = None
    best_score = 0
    results = []
    
    for i, config in enumerate(configs):
        print(f"\n搜索超参数组合 {i+1}/{len(configs)}")
        print(f"配置: {config}")
        
        try:
            # 训练和评估模型
            metrics = main_efficient_binary_classifier(config)
            score = metrics['f1']  # 使用F1分数作为评价指标
            
            # 记录结果
            result = {
                'config': config,
                'accuracy': metrics['accuracy'],
                'precision': metrics['precision'],
                'recall': metrics['recall'],
                'f1': metrics['f1']
            }
            results.append(result)
            
            # 更新最佳配置
            if score > best_score:
                best_score = score
                best_config = config
                print(f"找到新的最佳配置，F1分数: {best_score:.4f}")
            
            # 保存中间结果
            with open('hyperparam_results.json', 'w') as f:
                json.dump(results, f, indent=2)
        
        except Exception as e:
            print(f"配置 {config} 训练失败: {e}")
    
    print("\n网格搜索完成!")
    print(f"最佳配置: {best_config}")
    print(f"最佳F1分数: {best_score:.4f}")
    
    return best_config, results

def analyze_parameter_sensitivity(param_name, param_values):
    """分析特定超参数的敏感性"""
    import matplotlib.pyplot as plt
    
    # 基础配置
    base_config = {
        'base_model': 'resnet50',
        'dropout_rate': 0.5,
        'hidden_size': 512,
        'optimizer_type': 'adam',
        'learning_rate': 0.001,
        'weight_decay': 1e-4,
        'scheduler_type': 'step',
        'batch_size': 32,
        'stage1_epochs': 5,
        'stage2_epochs': 5,
        'data_augment_level': 'standard'
    }
    
    results = []
    
    for value in param_values:
        # 创建新配置
        config = base_config.copy()
        config[param_name] = value
        
        print(f"\n测试 {param_name} = {value}")
        
        # 训练和评估模型
        metrics = main_efficient_binary_classifier(config)
        results.append(metrics['accuracy'])
    
    # 绘制验证曲线
    plt.figure(figsize=(10, 6))
    plt.plot(param_values, results, 'o-', label='测试准确率')
    plt.xlabel(param_name)
    plt.ylabel('准确率')
    plt.title(f'{param_name} 的验证曲线')
    plt.grid(True)
    plt.legend()
    plt.savefig(f'validation_curve_{param_name}.png')
    plt.show()
    
    return results

def prepare_binary_datasets(food101_dir, local_data_dir, batch_size=32, augment_level='standard'):
    """准备二分类数据集，支持不同级别的数据增强"""
    # 标准数据增强
    if augment_level == 'standard':
        train_transform = transforms.Compose([
            transforms.RandomResizedCrop(224),
            transforms.RandomHorizontalFlip(),
            transforms.ColorJitter(brightness=0.2, contrast=0.2, saturation=0.2),
            transforms.ToTensor(),
            transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
        ])
    # 强数据增强
    elif augment_level == 'strong':
        train_transform = transforms.Compose([
            transforms.RandomResizedCrop(224),
            transforms.RandomHorizontalFlip(),
            transforms.RandomVerticalFlip(p=0.1),
            transforms.RandomRotation(15),
            transforms.ColorJitter(brightness=0.3, contrast=0.3, saturation=0.3, hue=0.1),
            transforms.RandomGrayscale(p=0.1),
            transforms.ToTensor(),
            transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225]),
            transforms.RandomErasing(p=0.1)
        ])
    
    # 测试和验证转换保持不变
    test_transform = transforms.Compose([
        transforms.Resize(256),
        transforms.CenterCrop(224),
        transforms.ToTensor(),
        transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
    ])
    
    # 数据转换字典
    data_transforms = {
        'train': train_transform,
        'val': test_transform,
        'test': test_transform
    }
    
    # 创建Food-101数据集
    food101_train = Food101Dataset(
        data_dir=food101_dir,
        transform=data_transforms['train'],
        split='train'
    )
    
    # 划分Food-101训练集，保留一小部分作为验证和测试
    food101_dataset_size = len(food101_train)
    train_size = int(food101_dataset_size * 0.7)
    val_size = int(food101_dataset_size * 0.15)
    test_size = food101_dataset_size - train_size - val_size
    
    food101_train_subset, food101_val, food101_test = torch.utils.data.random_split(
        food101_train, [train_size, val_size, test_size]
    )
    
    # 创建本地数据集（训练和验证集）
    local_train = CustomFoodDataset(
        root_dir=local_data_dir,
        list_file='train_list.txt',
        name_file='name.txt',
        transform=data_transforms['train']
    )
    
    local_val = CustomFoodDataset(
        root_dir=local_data_dir,
        list_file='val_list.txt',
        name_file='name.txt',
        transform=data_transforms['val']
    )
    
    # 使用新的加载器创建本地测试集
    local_test = CustomFoodTestDataset(
        root_dir=local_data_dir,
        test_list_file='test_truth_list.txt',  # 使用真实标签文件
        name_file='name.txt',
        transform=data_transforms['test']
    )
    
    # 为了保持数据平衡，我们从较大的数据集中采样与较小的数据集相同数量的样本
    # 这里假设local_train比food101_train_subset小
    if len(local_train) < len(food101_train_subset):
        food101_train_indices = torch.randperm(len(food101_train_subset))[:len(local_train)]
        food101_train_subset = torch.utils.data.Subset(food101_train_subset, food101_train_indices)
    elif len(food101_train_subset) < len(local_train):
        local_train_indices = torch.randperm(len(local_train))[:len(food101_train_subset)]
        local_train = torch.utils.data.Subset(local_train, local_train_indices)
    
    # 同样平衡验证和测试集
    if len(local_val) < len(food101_val):
        food101_val_indices = torch.randperm(len(food101_val))[:len(local_val)]
        food101_val = torch.utils.data.Subset(food101_val, food101_val_indices)
    elif len(food101_val) < len(local_val):
        local_val_indices = torch.randperm(len(local_val))[:len(food101_val)]
        local_val = torch.utils.data.Subset(local_val, local_val_indices)
    
    if len(local_test) < len(food101_test):
        food101_test_indices = torch.randperm(len(food101_test))[:len(local_test)]
        food101_test = torch.utils.data.Subset(food101_test, food101_test_indices)
    elif len(food101_test) < len(local_test):
        local_test_indices = torch.randperm(len(local_test))[:len(food101_test)]
        local_test = torch.utils.data.Subset(local_test, local_test_indices)
    
    # 组合数据集
    train_dataset = ConcatDataset([food101_train_subset, local_train])
    val_dataset = ConcatDataset([food101_val, local_val])
    test_dataset = ConcatDataset([food101_test, local_test])
    
    # 创建数据加载器
    dataloaders = {
        'train': DataLoader(train_dataset, batch_size=batch_size, shuffle=True, num_workers=4),
        'val': DataLoader(val_dataset, batch_size=batch_size, shuffle=False, num_workers=4),
        'test': DataLoader(test_dataset, batch_size=batch_size, shuffle=False, num_workers=4)
    }
    
    dataset_sizes = {
        'train': len(train_dataset),
        'val': len(val_dataset),
        'test': len(test_dataset)
    }
    
    return dataloaders, dataset_sizes

# 训练二分类器
def train_binary_classifier(model, dataloaders, dataset_sizes, criterion, optimizer, scheduler, num_epochs=5, device='cpu'):
    """训练二分类模型"""
    since = time.time()
    
    best_model_wts = copy.deepcopy(model.state_dict())
    best_acc = 0.0
    
    for epoch in range(num_epochs):
        print(f'Epoch {epoch+1}/{num_epochs}')
        print('-' * 10)
        
        # 每个epoch有训练和验证阶段
        for phase in ['train', 'val']:
            if phase == 'train':
                model.train()  # 设置模型为训练模式
            else:
                model.eval()   # 设置模型为评估模式
            
            running_loss = 0.0
            running_corrects = 0
            
            # 迭代数据
            for inputs, labels in dataloaders[phase]:
                inputs = inputs.to(device)
                labels = labels.to(device)
                
                # 梯度清零
                optimizer.zero_grad()
                
                # 前向传播
                with torch.set_grad_enabled(phase == 'train'):
                    outputs = model(inputs)
                    _, preds = torch.max(outputs, 1)
                    loss = criterion(outputs, labels)
                    
                    # 在训练阶段，反向传播和优化
                    if phase == 'train':
                        loss.backward()
                        optimizer.step()
                
                # 统计
                running_loss += loss.item() * inputs.size(0)
                running_corrects += torch.sum(preds == labels)
            
            if phase == 'train' and scheduler is not None and not isinstance(scheduler, torch.optim.lr_scheduler.ReduceLROnPlateau):
                scheduler.step()
            
            epoch_loss = running_loss / dataset_sizes[phase]
            epoch_acc = running_corrects.double() / dataset_sizes[phase]
            
            print(f'{phase} Loss: {epoch_loss:.4f} Acc: {epoch_acc:.4f}')
            
            # 如果是验证阶段并且有更好的准确率，保存模型
            if phase == 'val' and epoch_acc > best_acc:
                best_acc = epoch_acc
                best_model_wts = copy.deepcopy(model.state_dict())
                # 保存最佳模型
                torch.save({
                    'model_state_dict': best_model_wts,
                    'optimizer_state_dict': optimizer.state_dict(),
                    'epoch': epoch,
                    'accuracy': best_acc
                }, 'best_binary_classifier.pth')
                print(f'保存新的最佳模型，精度: {best_acc:.4f}')
            
            # 如果使用ReduceLROnPlateau，在验证阶段后调整学习率
            if phase == 'val' and scheduler is not None and isinstance(scheduler, torch.optim.lr_scheduler.ReduceLROnPlateau):
                scheduler.step(epoch_loss)
        
        print()
    
    time_elapsed = time.time() - since
    print(f'训练完成，用时 {time_elapsed // 60:.0f}m {time_elapsed % 60:.0f}s')
    print(f'最佳验证精度: {best_acc:.4f}')
    
    # 加载最佳模型权重
    model.load_state_dict(best_model_wts)
    return model

class EfficientBinaryClassifier(nn.Module):
    def __init__(self, base_model_name='resnet50', dropout_rate=0.5, hidden_size=512):
        """
        使用单一骨干网络的二分类器
        
        参数:
            base_model_name: 基础模型架构
            dropout_rate: Dropout比率
            hidden_size: 隐藏层大小
        """
        super(EfficientBinaryClassifier, self).__init__()
        
        # 创建基础模型（使用ImageNet预训练权重）
        if base_model_name == 'resnet18':
            base_model = models.resnet18(weights='IMAGENET1K_V1')
            num_features = base_model.fc.in_features
        elif base_model_name == 'resnet50':
            base_model = models.resnet50(weights='IMAGENET1K_V1')
            num_features = base_model.fc.in_features
        elif base_model_name == 'resnet101':
            base_model = models.resnet101(weights='IMAGENET1K_V1')
            num_features = base_model.fc.in_features
        else:
            raise ValueError(f"不支持的模型: {base_model_name}")
        
        # 去掉分类层
        self.features = nn.Sequential(*list(base_model.children())[:-1])
        
        # 二分类器
        self.classifier = nn.Sequential(
            nn.Flatten(),
            nn.Linear(num_features, hidden_size),
            nn.ReLU(),
            nn.Dropout(dropout_rate),
            nn.Linear(hidden_size, 2)
        )
    
    def forward(self, x):
        """前向传播"""
        # 特征提取
        features = self.features(x)
        
        # 分类
        output = self.classifier(features)
        
        return output

def main_efficient_binary_classifier(config=None):
    """
    带有超参数配置的二分类器训练函数
    
    参数:
        config: 超参数配置字典
    """
    # 默认配置
    default_config = {
        'base_model': 'resnet50',
        'dropout_rate': 0.5,
        'hidden_size': 512,
        'optimizer_type': 'adam',
        'learning_rate': 0.001,
        'weight_decay': 1e-4,
        'scheduler_type': 'step',
        'batch_size': 32,
        'stage1_epochs': 5,
        'stage2_epochs': 5,
        'data_augment_level': 'standard'  # 'standard' or 'strong'
    }
    
    # 更新配置
    if config is not None:
        default_config.update(config)
    
    config = default_config
    
    # 设置随机种子
    torch.manual_seed(42)
    np.random.seed(42)
    
    # 设置设备
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    print(f"使用设备: {device}")
    
    # 数据集路径
    food101_dir = 'data/food-101'  # Food-101数据集路径
    local_data_dir = 'data/dataset_release/release_data'  # 本地数据集路径
    
    # 准备数据集
    dataloaders, dataset_sizes = prepare_binary_datasets(
        food101_dir, 
        local_data_dir,
        batch_size=config['batch_size'],
        augment_level=config['data_augment_level']
    )
    print(f"数据集大小: 训练={dataset_sizes['train']}, 验证={dataset_sizes['val']}, 测试={dataset_sizes['test']}")
    
    # 创建二分类模型
    model = EfficientBinaryClassifier(
        base_model_name=config['base_model'],
        dropout_rate=config['dropout_rate'],
        hidden_size=config['hidden_size']
    )
    model = model.to(device)
    
    # 冻结特征提取部分的参数
    for param in model.features.parameters():
        param.requires_grad = False
    
    # 选择优化器
    if config['optimizer_type'].lower() == 'adam':
        optimizer = torch.optim.Adam(
            model.classifier.parameters(), 
            lr=config['learning_rate'],
            weight_decay=config['weight_decay']
        )
    elif config['optimizer_type'].lower() == 'sgd':
        optimizer = torch.optim.SGD(
            model.classifier.parameters(), 
            lr=config['learning_rate'],
            momentum=0.9,
            weight_decay=config['weight_decay']
        )
    elif config['optimizer_type'].lower() == 'adamw':
        optimizer = torch.optim.AdamW(
            model.classifier.parameters(), 
            lr=config['learning_rate'],
            weight_decay=config['weight_decay']
        )
    
    # 损失函数
    criterion = nn.CrossEntropyLoss()
    
    # 学习率调度器
    if config['scheduler_type'].lower() == 'step':
        scheduler = torch.optim.lr_scheduler.StepLR(optimizer, step_size=3, gamma=0.1)
    elif config['scheduler_type'].lower() == 'plateau':
        scheduler = torch.optim.lr_scheduler.ReduceLROnPlateau(
            optimizer, mode='min', factor=0.1, patience=2, verbose=True
        )
    elif config['scheduler_type'].lower() == 'cosine':
        scheduler = torch.optim.lr_scheduler.CosineAnnealingLR(
            optimizer, T_max=config['stage1_epochs']
        )
    
    # 训练模型的第一阶段（只训练分类器）
    print("训练第一阶段（只训练分类器）...")
    model = train_binary_classifier(
        model=model,
        dataloaders=dataloaders,
        dataset_sizes=dataset_sizes,
        criterion=criterion,
        optimizer=optimizer,
        scheduler=scheduler,
        num_epochs=config['stage1_epochs'],
        device=device
    )
    
    # 解冻特征提取部分的最后一层
    print("训练第二阶段（微调特征提取器的最后一层）...")
    for param in model.features[-1].parameters():
        param.requires_grad = True
    
    # 使用较小的学习率训练整个模型
    if config['optimizer_type'].lower() == 'adam':
        optimizer = torch.optim.Adam([
            {'params': model.classifier.parameters(), 'lr': config['learning_rate'] * 0.1},
            {'params': model.features[-1].parameters(), 'lr': config['learning_rate'] * 0.01}
        ], weight_decay=config['weight_decay'])
    elif config['optimizer_type'].lower() == 'sgd':
        optimizer = torch.optim.SGD([
            {'params': model.classifier.parameters(), 'lr': config['learning_rate'] * 0.1},
            {'params': model.features[-1].parameters(), 'lr': config['learning_rate'] * 0.01}
        ], momentum=0.9, weight_decay=config['weight_decay'])
    elif config['optimizer_type'].lower() == 'adamw':
        optimizer = torch.optim.AdamW([
            {'params': model.classifier.parameters(), 'lr': config['learning_rate'] * 0.1},
            {'params': model.features[-1].parameters(), 'lr': config['learning_rate'] * 0.01}
        ], weight_decay=config['weight_decay'])
    
    # 学习率调度器
    if config['scheduler_type'].lower() == 'step':
        scheduler = torch.optim.lr_scheduler.StepLR(optimizer, step_size=3, gamma=0.1)
    elif config['scheduler_type'].lower() == 'plateau':
        scheduler = torch.optim.lr_scheduler.ReduceLROnPlateau(
            optimizer, mode='min', factor=0.1, patience=2, verbose=True
        )
    elif config['scheduler_type'].lower() == 'cosine':
        scheduler = torch.optim.lr_scheduler.CosineAnnealingLR(
            optimizer, T_max=config['stage2_epochs']
        )
    
    # 训练模型的第二阶段
    model = train_binary_classifier(
        model=model,
        dataloaders=dataloaders,
        dataset_sizes=dataset_sizes,
        criterion=criterion,
        optimizer=optimizer,
        scheduler=scheduler,
        num_epochs=config['stage2_epochs'],
        device=device
    )
    
    # 保存完整模型
    model_save_path = f"binary_classifier_{config['base_model']}_{config['optimizer_type']}.pth"
    torch.save({
        'model_state_dict': model.state_dict(),
        'config': config
    }, model_save_path)
    print(f"完整模型已保存为 '{model_save_path}'")
    
    # 评估最终模型
    model.eval()
    test_metrics = evaluate_model(model, dataloaders['test'], device=device)
    print(f"测试准确率: {test_metrics['accuracy']:.4f}")
    print(f"F1分数: {test_metrics['f1']:.4f}")
    
    # 返回评估指标
    return test_metrics

# 评估函数
def evaluate_model(model, dataloader, device='cpu'):
    """评估模型性能并返回多个指标"""
    
    model.eval()
    all_labels = []
    all_preds = []
    
    with torch.no_grad():
        for inputs, labels in dataloader:
            inputs = inputs.to(device)
            labels = labels.to(device)
            
            outputs = model(inputs)
            _, preds = torch.max(outputs, 1)
            
            all_labels.extend(labels.cpu().numpy())
            all_preds.extend(preds.cpu().numpy())
    
    # 计算指标
    accuracy = accuracy_score(all_labels, all_preds)
    precision = precision_score(all_labels, all_preds, average='weighted')
    recall = recall_score(all_labels, all_preds, average='weighted')
    f1 = f1_score(all_labels, all_preds, average='weighted')
        
    return {
        'accuracy': accuracy,
        'precision': precision,
        'recall': recall,
        'f1': f1
    }
            
if __name__ == "__main__":
    # # 分析学习率敏感性
    # learning_rates = [0.0001, 0.0005, 0.001, 0.005, 0.01]
    # lr_results = analyze_parameter_sensitivity('learning_rate', learning_rates)
    
    # # 分析Dropout率敏感性
    # dropout_rates = [0.1, 0.3, 0.5, 0.7, 0.9]
    # dropout_results = analyze_parameter_sensitivity('dropout_rate', dropout_rates)

    # # 分析批量大小敏感性
    # batch_sizes = [8, 16, 32, 64, 128]
    # batch_results = analyze_parameter_sensitivity('batch_size', batch_sizes)
    
    # 基于之前的敏感性分析调整搜索空间
    param_grid = {
        'learning_rate': [0.0001 * 0.5, 0.0001, 0.0001 * 2],
        'dropout_rate': [0.3 * 0.8, 0.3, 0.3 * 1.2],
        'batch_size': [128]
        # 其他超参数
    }

    # 执行小规模网格搜索
    best_config, results = grid_search_hyperparameters()
    
    
    # 使用找到的最佳配置
    final_metrics = main_efficient_binary_classifier(best_config)
    print(f"最终模型准确率: {final_metrics['accuracy']:.4f}")
    print(f"最终模型F1分数: {final_metrics['f1']:.4f}")